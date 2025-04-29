package mcfile

import (
	"io/fs"
	"os"
	"fmt"
	"errors"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"

	CT "github.com/fbaube/ctoken"
)

// NewContentityFS proceeds as follows:
//  - initialize
//  - create an [os.DirFS]
//  - FIXME: an os.Root
//  - walk the DirFS, creating Contentities
//    and appending them to a slice
//  - process the list to identify and make
//    parent/child links
//
// The path argument should probably be an absolute filepath, 
// because a relative filepath might cause problems. Altho
// this is the opposite of the advice for lower-level items.
//
// It uses the global [CntyFS], which precludes
// re-entrancy and concurrency.
//
// Note that when we use [os.DirFS], it appears
// to make no difference whether path 
//  - is relative or absolute
//  - ends with a trailing slash or not
//  - is a directory or a symlink to a directory
//
// The only error returns for this func are:
//  - a bad path, rejected by func FU.NewFilepaths
//  - the path is not a directory (altho it can be
//    a symlnk to a directory ?)
//  - TBD: WHat happens of [os.Root] barfs on something ? 
// ContentityFS does not embed Errer and cannot
// itself return an error. FIXME: change this ? 
//
// TODO: Maybe it needs two boolean arguments:
//  - One to say whether to be strict about security (using [os.Root]
//    and Valid/Local, and
//  - One to say whether to follow symlinks.
// These two flags might have some interesting interactions.
// OTOH since this func can (and does?) use [os.Root], it can 
// easily (and should probably) also default to higher security 
// using funcs [io/fs.ValidPath] and [path/filepath.IsLocal].
//
// Accumulated NewContentity errors are counted
// in the field CotentityFS.nErrors 
// .
func NewContentityFS(aPath string, okayFilexts []string) (*ContentityFS, error){
     	var CntyFS *ContentityFS
     	var Empty  *ContentityFS
	var e error 
	var pFPs *FU.Filepaths
	Empty = new(ContentityFS) 

	// --------------------
	//  Prepare filepath(s)
	// --------------------
	L.L.Info("Making NewContentityFS: " + aPath)
	pFPs, e = FU.NewFilepaths(aPath)
	if e != nil {
	     	L.L.Error("NewCntyFS: bad path: %s", aPath)
		ee := &fs.PathError { Path:aPath, Err:e,
		       Op:"newcntyfs: bad root path" }
		// Empty.SetError(ee) 
		return Empty, ee
	}
	pathToUse := FU.EnsureTrailingPathSep(aPath)
	// os.DirFS(..) does not check or report problems
	// with the path argument, so we DIY here 
	if !FU.IsDirAndExists(pathToUse) {
		L.L.Error("NewCntyFS: Not a directory: %s", aPath)
		ee := &fs.PathError { Path:aPath, Err:errors.New(
		       "not a valid directory"), Op:"newcntyfs.root" }
		// Empty.SetError(ee)
		return Empty, ee
	}
	CntyFS = new(ContentityFS)
	// 2025.01 Change from path to AbsFP
	CntyFS.rootAbsPath = pFPs.AbsFP // path 
	L.L.Info("Path for new os.DirFS: " + SU.Tildotted(aPath))
	
	var osRoot *os.Root 
	osRoot, e = os.OpenRoot(pathToUse)
	CntyFS.FS = osRoot.FS()
	println("contentityfs_new L87 FIXED? os.Root")
	// CntyFS.FS = os.DirFS(pathToUse)
	
	// Initialize slice & map. Their length 
	// 0 will be detected by func [mustInit]
	CntyFS.asSlice = make([]*Contentity, 0)
	CntyFS.asMapOfAbsFP = make(map[string]*Contentity)

	// ==================
	//    FIRST PASS
	//  Load slice & map
	// ==================
	// NOTE that rel.path "." seems to be necessary 
	// here or else really weird errors occur.
	// Note that this is the place where [CntyFS]
	// being a global singleton can cause problems.
// ======================================================================
//	e = fs.WalkDir(CntyFS.FS, ".", wfnBuildContentityTree)
//
// wfnBuildContentityTree is a [fs.WalkDirFunc]
// 
// Altho it returns [*fs.PathError], this has to be declared as error
// because of the problems of an interface that is both nil and not nil.
//
// The basic procedure of the func is is:
//  - handle an error passed in 
//  - check validity of path argument (and reject if it is a file)
//  - handle first time thru (i.e. root node) 
//  - filter out unwanted values (and if unwanted dir, return [fs.SkipDir])
//  - add to slice - and also map - whether dir or file
//  - use materialised paths in slice to form links to build a tree 
//
// FIXME: Note that symlinks might not be handled securely, 
// not until [os.Root] is used. And even then, they might
// not be handled correctly. 
//
// This func filters out several file types:
//  - hidden (esp'ly .git directory)
//  - leading underbars ("_")
//  - emacs backup ("myfile~")
//  - this app's info  files: "*gtk,*gtr"
//  - this app's debug files: "*_(echo,tkns,tree)"
//  - filenames without a dot (indicating no file extension)
//  - NOTE that zero-length files (no content to analyse)
//    should NOT be filtered out 
//
// Note that as path separator, "/" is usually assumed, not [os.PathSep]. 
// .
// wfnBuildContentityTree(inPath string, inDE fs.DirEntry, inErr error) error {
// ======================================================================
e = fs.WalkDir(CntyFS.FS, ".",
func(inPath string, inDE fs.DirEntry, inErr error) error { // fs.WalkDirFunc
	// --------------------------
	//  Were we passed an error?
	// --------------------------
	if inErr != nil {
	   	 return CntyFS.handleWalkerErrorArgument(inPath, &inDE, inErr)
	}
	// --------------------
	//  Set some variables 
	// --------------------
	var isFirst = CntyFS.mustInitRoot() // first call ?
        var inName  = inDE.Name()
	var inDEisDir = inDE.IsDir()
	// If it's a directory, make sure it has a trailing slash.
	if inDEisDir {
	   inPath = FU.EnsureTrailingPathSep(inPath)
	   inName = FU.EnsureTrailingPathSep(inName)
	}
	// func [filepath.Abs] fails here cos it needs more than just 
	// a Base file name cos it does only lexical processing. 
	// absfp,_ = FP.Abs(path)

	// var p *Contentity
	var e error 
	// ==================
	//  HANDLE ROOT NODE 
	// (without filtering)
	// ==================
	if isFirst {
		L.L.Info("CntyFSWalker: inPath: " + inPath)
	   	if !inDEisDir { return &fs.PathError { Path:inPath,
		   	Op:"cntyfswalker.root", Err:errors.New("not a dir") } }
		L.L.Debug("cntyfswalker.root: path: %s / %s", inName, inPath)
		L.L.Debug("cntyfswalker.root: dirEntry: %+v", inDE)
	   	e = CntyFS.doInitRoot()
		if e == nil { return nil }
		return &fs.PathError { Err:e, Path:inPath,
		       Op:"newrootcnty.doinitroot" }
	}
	// ---------------------------
	//  Filter out unwanted stuff 
	// ---------------------------
	bad, rsn := excludeFilenamepath(inPath)
	if bad {
		L.L.Debug("Rejecting (%s): %s", inPath, rsn)
		if inDEisDir { return fs.SkipDir } 
		return nil
	}
	// -----------------------------------------------
	//  Now at this point, even if it's a directory,
	//  it's OK ! So let's go ahead and form the path
	//  of the file-or-dir and make the Contentity
	// -----------------------------------------------
	absPathToUse := FU.EnsureTrailingPathSep(
		        FP.Join(CntyFS.RootAbsPath(), inPath))
	var pCty *Contentity
	pCty = NewContentity(absPathToUse)
	if pCty.HasError() { 
		L.L.Warning("Rejecting (new Contentity(%s) failed): %T %+v",
			absPathToUse, e, e)
		CntyFS.nErrors++
		return nil
	}
	// ================================
	//  Now we do things based on just 
	//  what exactly the input DirEntry
	//  (inDE) is - file, dir, wotevs.
	// ================================
	// This is where bugs appeared when it's a directory,
	// cos calling code was assuming a valid Contentity.
	// TODO: Not sure what happens with symlinks
	if inDEisDir {
	   	if pCty.FSItem.TypedRaw == nil {
		   pCty.FSItem.TypedRaw = new(CT.TypedRaw)
		   } 
	        pCty.FSItem.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE
		CntyFS.nDirs++ // just a simple counter
		// println("================ DIR ========")
		// These next two stmts should barf, cos
		// they should not be allocated for a dir !
		// p.MimeType = "dir"
		// p.MType = SU.MU_type_DIRLIKE
		L.L.Okay("Item (DIR/) OK; no MmeType or MType") 
	} else if inDE.Type() == 0 { // regular file
		CntyFS.nFiles++ // just a simple counter 
		L.L.Okay("Item (FILE) OK: MType<%s> MarkupType<%s>",
			pCty.MType, pCty.RawType())
	} else if (inDE.Type() & fs.ModeSymlink) != 0 { // Symlink
	       	if pCty.FSItem.TypedRaw == nil {
		   pCty.FSItem.TypedRaw = new(CT.TypedRaw)
		}
	        pCty.FSItem.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE // OOPS
		CntyFS.nMiscs++ // just a simple counter 
		L.L.Okay("Item (SYML) OK: what to do ?!")
	} else { // Some weirdness in the Mode bits 
	        pCty.FSItem.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE
             // CntyFS.nMiscs++ // just a simple counter
		CntyFS.nErrors++
                L.L.Error("Item (WTF) BAD: what to do ?!")
	}
	// -------------------------
	//   Also add it to the
	//  arena-slice and the map
	// -------------------------
	CntyFS.asSlice = append(CntyFS.asSlice, pCty)
	CntyFS.asMapOfAbsFP[absPathToUse] = pCty
	// L.L.Info("ADDED TO MAP L241: " + pathToUse)
	return nil
})
// ======================================================================
	if e != nil {
		// L.L.Panic("NewCntyFS.WalkDir: " + e.Error())
		return nil, &fs.PathError { Op:"newcntyfs.walkdir",
		       Err:e, Path:aPath } 
	}
	L.L.Okay("NewCntyFS: walked OK %d nords from path %s",
		 len(CntyFS.asSlice), pathToUse)

	// Debuggery 
	for ii, cc := range CntyFS.asSlice {
	    if cc == nil {
	       L.L.Error ("OOPS, CntyFS.asSlice[%02d] is NIL", ii)
	       continue
	    }
	    /* if cc.FSItem == nil || cc.FSItem.FileMeta == nil {
	       L.L.Error("WTF, man!")
	       continue } */
	    if cc.FSItem.IsDirlike() {
	        L.L.Debug("[%02d] isDIRLIKE: AbsFP: %s",
			ii, cc.FSItem.FPs.AbsFP)
	    } else {
		L.L.Debug("[%02d] MarkupType: %s", ii, cc.RawType())
	    }
	}

	// ==============================================
	//      SECOND PASS
	//  Range over the slice, using the materialised
	//  paths in asMapToAbsFS to identify parent/kid 
	//  Nord relationships and link together
	// ==============================================
	// TODO This needs to be in some generalized 
	// form, such as TreeFromMaterializedPaths
	// =========================================
	var i int
	var pC *Contentity
	for i, pC = range CntyFS.asSlice {
		if i == 0 { // skip over root 
			continue
		}
		// ---------------------------
		//  Shortcut if child of root
		// ---------------------------
		if !S.Contains(pC.RelFP(), FU.PathSep) {
			CntyFS.rootNord.AddKid(pC)
			continue
		}
		// --------------------------
		//   Get dir portion of path
		// --------------------------
		itsDir := FP.Dir(pC.RelFP())
		itsDir = FU.EnsureTrailingPathSep(itsDir)
		// println(n.Path, "|cnex2|", itsDir)
		// L.L.Warning("itsDir: " + itsDir)
		// L.L.Warning("theMap: %+v", CntyFS.asMap)
		var pPar *Contentity
		var ok bool
		// PROBLEMS HERE ?
		// The parent directory should be in the map.
		// If it's not, then possibly we have messed
		// up with trailing separators. 
		if pPar, ok = CntyFS.asMapOfAbsFP[itsDir]; !ok {
			L.L.Error("findParentInMap: failed for: " +
				itsDir + " of " + pC.AbsFP())
			println(fmt.Sprintf("%+v", CntyFS.asMapOfAbsFP))
			panic(pC.AbsFP())
		}
		/*
		if itsDir != par.AbsFP() { // or, Rel? 
			panic(itsDir + " != " + par.AbsFP())
		}
		*/
		pPar.AddKid(pC)
	}
	// TODO Look for entries that do not have a parent assigned !
	
	/* more debugging
	println("DUMP LIST")
	for _, n := range pFTFS.asSlice {
		println(n.LinePrefixString(), n.LineSummaryString())
	}
	println("DUMP MAP")
	for k, v := range pFTFS.asMap {
		fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
	}
	*/
	// println(SU.Gbg("=== TREE ==="))
	// CntyFS.rootNord.PrintAll(os.Stdout)
	return CntyFS, nil
}
