package mcfile

import (
	"errors"
	"io/fs"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	// FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
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
	L.L.Dbg("wfnBuildContentityTree: path: %s", path)
	// L.L.Warning("wfnBuildContentityTree: dir: %+v", d)
	if err != nil {
		L.L.Error("wfnBuildContentityTree: "+
			"UNHANDLED err: %s", err.Error())
	}
	var p *Contentity
	var e error
	// ROOT NODE ?
	if mustInitRoot() {
		var pRC *RootContentity
		if pCFS.RootAbsPath() == "" {
			panic("wfnBuildContentityTree: no ROOT")
		}
		pRC, e = NewRootContentity(pCFS.RootAbsPath())
		if e != nil || pRC == nil {
			return errors.New("wfnBuildContentityTree UNHANDLED" +
				" mustInitRoot NewRootContentityNord L101")
		}
		// Assign to globals (i.e. package vars)
		pCFS.rootNord = pRC
		pRC.MimeType = "dir"
		pRC.MType = "dir"
		// println("wfnBuildContentityTree: root node abs.FP:\n\t", p.AbsFP())
		var p *Contentity
		p = ((*Contentity)(pRC))
		pCFS.asSlice = append(pCFS.asSlice, p)
		pCFS.asMap[path] = p
		// println("ADDED TO MAP:", path)
		pCFS.nDirs = 1
		pCFS.nFiles = 0
		return nil // NOT pRC
	}
	// Filter out hidden, esp'ly .git & .git* ;
	// note that "/" is assumed, not os.Sep ;
	// we issue no message, cos .git/* can be huge.
	if S.Contains(path, "/.git") {
		return nil
	}
	/* Old gnarly debuggging
	// IS USING CWD NOT CLI ARGUMENT ?!
	L.L.Warning("ROOT: %s", pCFS.RootAbsPath())
	L.L.Warning("PATH: %s", path)
	var sfp string
	var afp FU.AbsFilePath
	sfp = FP.Join(pCFS.RootAbsPath(), path)
	L.L.Warning("BOTH: %s", sfp)
	afp = FU.AbsFilePath(sfp)
	*/
	var reasonToReject string
	if S.HasPrefix(path, ".") || S.HasPrefix(path, "_") ||
		S.Contains(path, "/.") || S.Contains(path, "/_") {
		reasonToReject = "leading . or _ "
		L.L.Warning("Rejecting (%s): %s", reasonToReject, path)
		return nil
	}
	// Now at this point, if it's a directory, it's OK !
	sfp := FP.Join(pCFS.RootAbsPath(), path)
	afp := FU.AbsFilePath(sfp)
	var isDir = afp.DirExists()

	// And so following code applies only to files, not to directories
	// TODO: Not sure what happens with symlinks
	if !isDir {
		var ln = len(path)
		if S.HasSuffix(path, "gtk") || S.HasSuffix(path, "gtr") {
			reasonToReject = "gtk/gtr"
		} else if ln >= 5 && path[ln-5] == '_' { // debug file via "-t" flag
			reasonToReject = "_echo,_tkns,_tree"
		} else if S.Index(FP.Base(path), ".") == -1 { // untyped file
			reasonToReject = "no dot, untyped"
		} else if S.HasSuffix(path, "~") {
			reasonToReject = "emacs"
		}
		if reasonToReject != "" {
			L.L.Warning("Rejecting (%s): %s", reasonToReject, path)
			return nil
		}
	}
	p, e = NewContentity(path) // FP.Join(pCFS.RootAbsPath(), path))
	if p == nil || e != nil {
		L.L.Warning("Rejecting (no Contentity): " + path)
		L.L.Error("Skipping this item")
		return nil
	}
	if /* p.PathAnalysis == nil || */ p.PathProps.Raw == "" && !isDir {
		L.L.Warning("Rejecting (len 0 or too short): " + path)
		L.L.Error("Skipping this item")
		return nil
	}
	L.L.Dbg("Directory traverser: MarkupType: " + string(p.MarkupType()))
	nxtIdx := len(pCFS.asSlice)
	pCFS.asSlice = append(pCFS.asSlice, p)
	p.logIdx = nxtIdx
	pCFS.asMap[path] = p
	if p.IsDir() {
		pCFS.nDirs++
		// println("================ DIR ========")
		p.MimeType = "dir"
		p.MType = "dir"
	} else {
		pCFS.nFiles++
	}
	// println("ADDED TO MAP:", path)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
