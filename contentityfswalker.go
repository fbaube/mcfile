package mcfile

import (
	// "errors"
	"io/fs"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	CT "github.com/fbaube/ctoken"
	// FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
	assert "github.com/lainio/err2/assert"
)

// wfnBuildContentityTree is
// type fs.WalkDirFunc func(path string, d DirEntry, err error) error
//
// Therefore it can NOT be declared as returns *[fs.PathError].
// However it can and does return them, which has the problems
// of an interface that is both nil and not nil.
//
// Here the variable [CntyFS] is used as a global singleton,
// which is very dodgy and will cause problems if used in a
// re-entrant way or with concurrency. 
//
// It filters out several file types:
//  - (TODO:) zero-length file (no content to analyse)
//  - hidden (esp'ly .git directory)
//  - leading underbars ("_")
//  - emacs backup ("myfile~")
//  - this app's debug files: "*_(echo,tkns,tree)"
//  - filenames without a dot (indicating no file extension)
// TODO Make these arguments: filter-out prefixes, midfixes, suffixes. 
//
// As path separator, "/" is usually assumed, not [os.PathSep]. 
// .
func wfnBuildContentityTree(path string, d fs.DirEntry, err error) error {
        var name string 
        name = d.Name()
	
	// This is UNreliable. Func [filepath.Abs]
	// needs more than just a Base file name. 
	// absfp,_ = FP.Abs(path)
	
	// If it's a directory, make sure it has a trailing slash.
	if d.IsDir() {
	   path = FU.EnsureTrailingPathSep(path)
	   name = FU.EnsureTrailingPathSep(name)
	   if !FU.IsDirAndExists(FP.Join(CntyFS.rootAbsPath, path)) {
	      L.L.Error("IsDir/NonDir: rootPath<%s>; " + 
	      	"path<%s>; err<%+v> dirEntry<%#v>",
	      	CntyFS.rootAbsPath, path, err, d)
		}
	   }
	L.L.Debug("wfnBuildContentityTree: path: %s / %s", name, path)
	L.L.Debug("wfnBuildContentityTree: dir: %+v", d)
	if err != nil {
		L.L.Error("wfnBuildContentityTree: "+
			"UNHANDLED in-arg-err: %w", err)
	}
	// var p *Contentity
	var e error
	// ==================
	//  HANDLE ROOT NODE 
	// ==================
	if mustInitRoot() {
		var pRC *RootContentity
		assert.NotEmpty(CntyFS.RootAbsPath())
		/* if CntyFS.RootAbsPath() == "" {
			panic("wfnBuildContentityTree: no ROOT")
		} */
		pRC, e = NewRootContentity(CntyFS.RootAbsPath())
		if e != nil || pRC == nil {
			return &fs.PathError{Op:"WalkFn.NewRootContentity",
			Err:e,Path:CntyFS.RootAbsPath()}
		}
		// assert.That(pRC.IsDir()) SHOULD NOT FAIL, BUT DID
		// Assign to globals (i.e. package vars)
		CntyFS.rootNord = pRC
		// These next two get NPE cos no such struct for a dir 
		// pRC.MimeType = "dir"
		// pRC.MType = "dir"
		if pRC.FSItem.TypedRaw == nil {
		   println("Oops, contentityfswalker, newRoot has no TypedRaw")
		   pRC.FSItem.TypedRaw = new(CT.TypedRaw)
		}
		pRC.FSItem.Raw_type = SU.Raw_type_DIRLIKE
		// println("wfnBuildContentityTree: root node abs.FP:\n\t", p.AbsFP())
		var pC *Contentity
		pC = ((*Contentity)(pRC))
		CntyFS.asSlice = append(CntyFS.asSlice, pC)
		CntyFS.asMap[CntyFS.RootAbsPath()] = pC
		// L.L.Warning("ADDED TO MAP L84: " + CntyFS.RootAbsPath())
		CntyFS.nDirs = 1
		CntyFS.nFiles = 0
		return nil // NOT pRC! This is a walker func 
	}
	// =====================
	//  FILTER OUT UNWANTED
	// =====================
	// Filter out .git/ and other dot-directories ASAP. 
	// Note that "/" is assumed, not os.Sep .
	// Give no reason, cos .git/* et al. can be frickin' huge.
	if S.HasSuffix (path, "/.git") ||
	   S.Contains (path, "/.git/") ||
	  (S.Contains(path, "/.") && d.IsDir()) {
		return fs.SkipDir
	}
	var b1, b2 bool
	var r1, r2 string
	b1, r1 = excludeFilenamepath(path)
	b2, r2 = excludeFilenamepath(name)
	if b1 || b2 {
	      if (b1) { r1 = "(file path) " + r1 }
	      if (b2) { r2 = "(file name) " + r2 }
		L.L.Debug("Rejecting (%s): %s%s", path, r1, r2)
		// continue 
		return nil
	}
	// =======================================
	// Now at this point, if it's a directory,
	// it's OK ! So let's go ahead and form 
	// the path and make the Contentity
	// =======================================
	var pathToUse string
	var pathIsDir bool 
	pathToUse = FP.Join(CntyFS.RootAbsPath(), path)
	pathIsDir = FU.IsDirAndExists(pathToUse)
	if pathIsDir {
	   pathToUse = FU.EnsureTrailingPathSep(pathToUse)
	}
	assert.Equal(pathIsDir, d.IsDir())

	{
		var reason string 
		var ln = len(pathToUse)
		if S.HasSuffix(pathToUse, "gtk") ||
		   S.HasSuffix(pathToUse, "gtr") {
			reason = "gtk/gtr"
		} else if ln >= 5 && pathToUse[ln-5] == '_' {
		     // debug file via "-t" flag
			reason = "_echo,_tkns,_tree"
		} 
		if reason != "" {
			L.L.Debug("Rejecting (%s): %s", pathToUse, reason)
			// continue 
			return nil
		}
	}

	var pCty *Contentity
	// println("PATH TO TRY IS: " + pathToUse)
	pCty, e = NewContentity(pathToUse)
	if pCty == nil { 
		L.L.Warning("Rejecting (new Contentity(%s) failed): %T %+v",
			pathToUse, e, e)
		return nil
	}
	// And so following code applies only to files, not to directories
	// TODO: Not sure what happens with symlinks
	if pathIsDir {
	        pCty.FSItem.Raw_type = SU.Raw_type_DIRLIKE
		CntyFS.nDirs++
		// println("================ DIR ========")
		// These next two stmts should barf, cos
		// they should not be allocated for a dir !
		// p.MimeType = "dir"
		// p.MType = SU.MU_type_DIRLIKE
		L.L.Okay("Item (DIR) OK; CntyPtr nil") // : MType<%s>", pCty.MType)
	} else {
		CntyFS.nFiles++
		L.L.Okay("Item OK: MType<%s> MarkupType<%s>",
			pCty.MType, pCty.RawType())
	}
	// L.L.Info("Directory traverser: MarkupType: " + string(p.MarkupType()))
	// nxtIdx := len(CntyFS.asSlice)
	CntyFS.asSlice = append(CntyFS.asSlice, pCty)
	// p.logIdx = nxtIdx // NPE
	CntyFS.asMap[pathToUse] = pCty
	// L.L.Warning("ADDED TO MAP L168: " +pathToUse)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
