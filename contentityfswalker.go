package mcfile

import (
	"errors"
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
	absfp,_ = FP.Abs(path)
	// If it's a directory, make sure it has a trailing slash
	if d.IsDir() {
	   if !S.HasSuffix(path, "/") {
	      path += "/"; L.L.Progress("Dir path gets trlg slash") }
	   if !S.HasSuffix(name, "/") {
	      name += "/"; L.L.Progress("Dir name gets trlg slash") }
	   if !S.HasSuffix(absfp,"/") {
	      absfp+= "/"; L.L.Progress("Dir abs.path gets trlg slash") }
	   assert.That(FU.AbsFilePath(absfp).DirExists())
	   }
	L.L.Progress("wfnBuildContentityTree: path: %s / %s", name, path)
	L.L.Info("wfnBuildContentityTree: dir: %+v", d)
	if err != nil {
		L.L.Error("wfnBuildContentityTree: "+ "UNHANDLED err: %w", err)
	}
	var p *Contentity
	var e error
	// ==================
	//  HANDLE ROOT NODE 
	// ==================
	if mustInitRoot() {
		var pRC *RootContentity
		assert.NotEmpty(pCFS.RootAbsPath())
		/* if pCFS.RootAbsPath() == "" {
			panic("wfnBuildContentityTree: no ROOT")
		} */
		pRC, e = NewRootContentity(pCFS.RootAbsPath())
		if e != nil || pRC == nil {
			return &fs.PathError{Op:"NewRootContentity",
			Err:errors.New("wfnBuildContentityTree UNHANDLED" +
			" mustInitRoot NewRootContentity L62"),
			Path:pCFS.RootAbsPath()}
		}
		assert.That(false)
		assert.That(pRC.IsDir())
		// Assign to globals (i.e. package vars)
		pCFS.rootNord = pRC
		pRC.MimeType = "dir"
		pRC.MType = "dir"
		pRC.FSItem.MarkupType = SU.MU_type_DIRLIKE
		// println("wfnBuildContentityTree: root node abs.FP:\n\t", p.AbsFP())
		var pC *Contentity
		pC = ((*Contentity)(pRC))
		pCFS.asSlice = append(pCFS.asSlice, pC)
		pCFS.asMap[path] = pC
		// println("ADDED TO MAP:", path)
		pCFS.nDirs = 1
		pCFS.nFiles = 0
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
	sfp := FP.Join(pCFS.RootAbsPath(), path)
	afp := FU.AbsFilePath(sfp)
	var isDir = afp.DirExists()
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

	p, e = NewContentity(path) // FP.Join(pCFS.RootAbsPath(), path))
	if p == nil || e != nil {
		L.L.Warning("Rejecting (new Contentity failed): " + path)
		return nil
	}
	// And so following code applies only to files, not to directories
	// TODO: Not sure what happens with symlinks
	if isDir {
	        p.FSItem.MarkupType = SU.MU_type_DIRLIKE
		pCFS.nDirs++
		// println("================ DIR ========")
		p.MimeType = "dir"
		p.MType = "dir"
	} else {
		pCFS.nFiles++
	}
	/*
	if p.PathAnalysis == nil || p.FSItem.Raw == "" && !isDir {
		L.L.Warning("Rejecting (%s): zero length", path)
		L.L.Error("Skipping this item")
		return nil
	}
	*/
	L.L.Dbg("Directory traverser: MarkupType: " + string(p.MarkupType()))
	nxtIdx := len(pCFS.asSlice)
	pCFS.asSlice = append(pCFS.asSlice, p)
	p.logIdx = nxtIdx
	pCFS.asMap[path] = p
	// println("ADDED TO MAP:", path)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
