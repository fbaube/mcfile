package mcfile

import (
	// "errors"
	"io/fs"
	// "fmt"
	FP "path/filepath"
	// S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	CT "github.com/fbaube/ctoken"
	// FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
)

// wfnBuildContentityTree is
// type fs.WalkDirFunc func(path string, d DirEntry, err error) error
//
// ACCURATE?? It can NOT return a [*fs.PathError], because of 
// the problems of an interface that is both nil and not nil.
//
// The basic procedure is:
//  - check validity of path argument (and reject if it is a file) 
//  - filter out unwanted values (and if dir, return [os.SkipDir])
//  - add to tree - whether dir or file 
//
// Here the variable [CntyFS] is used as a global singleton,
// which is very dodgy and will cause problems if used in a
// re-entrant way or with concurrency.
//
// Note that symlinks are probably not handled correctly. 
//
// It filters out several file types:
//  - hidden (esp'ly .git directory)
//  - leading underbars ("_")
//  - emacs backup ("myfile~")
//  - this app's info  files: "*gtk,*gtr"
//  - this app's debug files: "*_(echo,tkns,tree)"
//  - filenames without a dot (indicating no file extension)
//  - NOTE that zero-length files (no content to analyse)
//    should not be filtered out 
// TODO Convert the following into arguments to this func:
// filter-out prefixes, midfixes, suffixes. 
//
// Note that as path separator, "/" is usually assumed, not [os.PathSep]. 
// .
func wfnBuildContentityTree(inPath string, inDE fs.DirEntry, inErr error) error {
     	// So what's the deal with the error input argument ? Well...
	// https://pkg.go.dev/io/fs#WalkDirFunc
	// The err argument reports an error related to path,
	// such that WalkDir will not walk into that directory.
	// The function can decide how to handle that error; as
	// described earlier, returning the error causes WalkDir
	// to stop walking the entire tree. WalkDir calls the
	// function with a non-nil err argument in two cases:
	//  - First, if the initial Stat on the root directory fails,
	//    WalkDir calls the function with path set to root, d set
	//    to nil, and err set to the error from fs.Stat.
	//  - Second, if a directory's ReadDir method (see
	//    https://pkg.go.dev/io/fs#ReadDirFile) fails, WalkDir
	//    calls the function with path set to the directory's path,
	//    d set to an DirEntry describing the directory, and err 
	//    set to the error from ReadDir. In this second case, the 
	//    function is called twice with the path of the directory: 
	//    the first call is before the directory read is attempted 
	//    and has err set to nil, giving the function a chance to 
	//    return SkipDir or SkipAll and avoid the ReadDir entirely. 
	//    The second call is after a failed ReadDir and reports the
	//    error from ReadDir. (If ReadDir succeeds, there is no 
	//    second call.)

	// --------------------------
	//  Were we passed an error?
	// --------------------------
	if inErr != nil {
	   	 return CntyFS.handleWalkerErrorArgument(inPath, inDE, inErr)
	}
	// --------------------
	//  Set some variables 
	// --------------------
	var needInit = CntyFS.mustInitRoot() // first call ? 
        var inName   = inDE.Name()
	var inIsDir  = inDE.IsDir()
	// func [filepath.Abs] needs more than just a Base 
	// file name, because it does only lexical processing. 
	// absfp,_ = FP.Abs(path)

	// --------------------
	//  root must be a dir 
	// --------------------
	if needInit && !inIsDir {
	   	return &fs.PathError { Path:inPath,
		   	Op:"cntyfswalker.isdir", Err:inErr }
	}
	// If it's a directory, make sure it has a trailing slash.
	// We also used to check for existence, but this is silly
	// cos the [fs.DirEntry] argument passed in has the info.
	if inIsDir {
	   inPath = FU.EnsureTrailingPathSep(inPath)
	   inName = FU.EnsureTrailingPathSep(inName)
	   }
	L.L.Debug("wfnBuildContentityTree: path: %s / %s", inName, inPath)
	L.L.Debug("wfnBuildContentityTree: dir: %+v", inDE)
	
	// var p *Contentity
	var e error 
	// ==================
	//  HANDLE ROOT NODE 
	// (without filtering)
	// ==================
	if needInit {
	   	e = CntyFS.doInitRoot()
		if e != nil {
		     	return &fs.PathError { Op:"newrootcnty.doinitroot",
				 Err:e,Path:CntyFS.RootAbsPath() }
		}
		return nil 
	}
	// ---------------------------
	//  Filter out unwanted stuff 
	// ---------------------------
	var bP, bN bool
	var rP, rN string
	bP, rP = excludeFilenamepath(inPath)
	bN, rN = excludeFilenamepath(inName)
	if bP || bN {
	        if (bP) { rP = "(path) " + rP }
	        if (bN) { rN = "(name) " + rN }
		L.L.Debug("Rejecting (%s): %s%s", inPath, rP, rN)
		// continue
		if inIsDir { return fs.SkipDir } 
		return nil
	}
	// -----------------------------------------------
	//  Now at this point, even if it's a directory,
	//  it's OK ! So let's go ahead and form the path
	//  of the file-or-dir and make the Contentity
	// -----------------------------------------------
	var absPathToUse string
	absPathToUse = FP.Join(CntyFS.RootAbsPath(), inPath)

	var pCty *Contentity
	// println("ABS.PATH TO USE IS: " + absPathToUse)
	pCty, e = NewContentity(absPathToUse)
	if pCty == nil { 
		L.L.Warning("Rejecting (new Contentity(%s) failed): %T %+v",
			absPathToUse, e, e)
		return nil
	}
	// This is where bugs have appeared when it's a directory,
	// because other code was assuming a Contentity.
	// TODO: Not sure what happens with symlinks
	if inIsDir {
	   	if pCty.FSItem.TypedRaw == nil {
		   pCty.FSItem.TypedRaw = new(CT.TypedRaw)
		   } 
	        pCty.FSItem.Raw_type = SU.Raw_type_DIRLIKE
		CntyFS.nDirs++
		// println("================ DIR ========")
		// These next two stmts should barf, cos
		// they should not be allocated for a dir !
		// p.MimeType = "dir"
		// p.MType = SU.MU_type_DIRLIKE
		L.L.Okay("Item (DIR) OK; CntyPtr nil") // : MType<%s>", pCty.MType)
	} else { // non-dir 
		CntyFS.nFiles++ // just a simple counter 
		L.L.Okay("Item OK: MType<%s> MarkupType<%s>",
			pCty.MType, pCty.RawType())
	}
	// L.L.Info("Directory traverser: MarkupType: " + string(p.MarkupType()))

	// -------------------------
	//   Also add it to the
	//  arena-slice and the map
	// -------------------------
	CntyFS.asSlice = append(CntyFS.asSlice, pCty)
	// p.logIdx = nxtIdx // oops NPE
	CntyFS.asMap[absPathToUse] = pCty
	// L.L.Warning("ADDED TO MAP L204: " + pathToUse)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
