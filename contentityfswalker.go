package mcfile

import (
	// "errors"
	"io/fs"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	// FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
	assert "github.com/lainio/err2/assert"
)

// wfnBuildContentityTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
//
// It filters out several file types:
//  - (TODO:) zero-length file (no content to analyse)
//  - hidden (esp'ly .git directory)
//  - leading underbars ("_")
//  - emacs backup ("myfile~")
//  - this app's debug files: "*_(echo,tkns,tree)"
//  - filenames without a dot (indicating no file extension)
//
// As path separator, "/" is assumed, not os.Sep
// .
func wfnBuildContentityTree(path string, d fs.DirEntry, err error) error {
        var name, absfp string 
        name = d.Name()
	
	// This absfp is UNreliable !!  WTF.
	// absfp,_ = FP.Abs(path)
	
	// If it's a directory, make sure it has a trailing slash
	if d.IsDir() {
	   if !S.HasSuffix(path, "/") {
	      path += "/"; L.L.Progress("Dir path gets trlg slash") }
	   if !S.HasSuffix(name, "/") {
	      name += "/"; L.L.Progress("Dir name gets trlg slash") }
	   if !S.HasSuffix(absfp,"/") {
	      absfp+= "/"; L.L.Progress("Dir abs.path gets trlg slash") }
	// L.L.Warning("ASSERTING DIR EXISTS: " + absfp)
	// if !FU.IsDirAndExists(absfp) {
	   if !FU.IsDirAndExists(FP.Join(CntyFS.rootAbsPath, path)) {
	      L.L.Error("OOPS?? rootPath<%s>; path<%s>, absfp<%s>; " +
	      	"err<%+v> dirEntry<%#v>; TOGETHER<%s>",
	      	CntyFS.rootAbsPath, path, absfp, err, d, FP.Join(absfp, path))
		}
	   }
	L.L.Progress("wfnBuildContentityTree: path: %s / %s", name, path)
	L.L.Info("wfnBuildContentityTree: dir: %+v", d)
	if err != nil {
		L.L.Error("wfnBuildContentityTree: "+
			"UNHANDLED in-arg-err: %w", err)
	}
	var p *Contentity
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
		pRC.FSItem.MarkupType = SU.MU_type_DIRLIKE
		// println("wfnBuildContentityTree: root node abs.FP:\n\t", p.AbsFP())
		var pC *Contentity
		pC = ((*Contentity)(pRC))
		CntyFS.asSlice = append(CntyFS.asSlice, pC)
		CntyFS.asMap[path] = pC
		// println("ADDED TO MAP:", path)
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
		L.L.Warning("Rejecting (%s): %s%s", path, r1, r2)
		// continue 
		return nil
	}
	// Now at this point, if it's a directory, it's OK !
	sfp := FP.Join(CntyFS.RootAbsPath(), path)
	// afp := FU.AbsFilePath(sfp)
	var isDir = FU.IsDirAndExists(sfp) // afp.DirExists()
	assert.Equal(isDir, d.IsDir())

	{
		var reason string 
		var ln = len(path)
		if S.HasSuffix(path, "gtk") || S.HasSuffix(path, "gtr") {
			reason = "gtk/gtr"
		} else if ln >= 5 && path[ln-5] == '_' { // debug file via "-t" flag
			reason = "_echo,_tkns,_tree"
		} 
		if reason != "" {
			L.L.Warning("Rejecting (%s): %s", path, reason)
			// continue 
			return nil
		}
	}

	// println("PATH TO TRY IS: " + path)
	tryPath := FP.Join(CntyFS.RootAbsPath(), path)
	println("PATH TO TRY IS: " + tryPath)
	p, e = NewContentity(tryPath) // path
	if p == nil { // || e != nil {
		L.L.Warning("Rejecting (new Contentity(%s) failed): %T %+v",
			tryPath, e, e)
		return nil
	}
	// And so following code applies only to files, not to directories
	// TODO: Not sure what happens with symlinks
	if isDir {
	        p.FSItem.MarkupType = SU.MU_type_DIRLIKE
		CntyFS.nDirs++
		// println("================ DIR ========")
		p.MimeType = "dir"
		p.MType = "dir"
	} else {
		CntyFS.nFiles++
	}
	/*
	if p.PathAnalysis == nil || p.FSItem.Raw == "" && !isDir {
		L.L.Warning("Rejecting (%s): zero length", path)
		L.L.Error("Skipping this item")
		return nil
	}
	*/
	L.L.Dbg("Directory traverser: MarkupType: " + string(p.MarkupType()))
	nxtIdx := len(CntyFS.asSlice)
	CntyFS.asSlice = append(CntyFS.asSlice, p)
	p.logIdx = nxtIdx
	CntyFS.asMap[path] = p
	// println("ADDED TO MAP:", path)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
