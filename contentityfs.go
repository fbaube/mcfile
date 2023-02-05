package mcfile

import (
	"errors"
	"io/fs"
	FP "path/filepath"
	S "strings"

	FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
)

type ContentityFS struct {
	FSU.BaseFS
	rootNord      *RootContentity
	asSlice       []*Contentity
	asMap         map[string]*Contentity // string is Rel.Path
	nFiles, nDirs int
}

/* Open is a dummy function, just here to satisfy an interface.
func (p *ContentityFS) Open(path string) (fs.File, error) {
	return p. /*inputFS. * / Open(path)
} */

func (p *ContentityFS) Size() int {
	return len(p.asSlice)
}

func (p *ContentityFS) DirCount() int {
	return p.nDirs
}

func (p *ContentityFS) FileCount() int {
	return p.nFiles
}

func (p *ContentityFS) RootContentity() *RootContentity {
	return p.rootNord
}

func (p *ContentityFS) AsSlice() []*Contentity {
	var z []*Contentity
	z = p.asSlice
	return z // p.AsSlice
}

func mustInitRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(pCFS.asSlice) == 0 && len(pCFS.asMap) == 0)
	didDoInit = (len(pCFS.asSlice) > 0 && len(pCFS.asMap) > 0)
	if !(needsInit || didDoInit) {
		panic("mustInitRoot: illegal state")
	}
	return needsInit
}

func (p *ContentityFS) DoForEvery(stgprocsr ContentityStage) {

}

// wfnBuildContentityTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
//
// It filters out several file types:
// - (TODO:) zero-length file (no content to analyse)
// - hidden (esp'ly .git directory)
// - leading underbars ("_")
// - emacs backup ("myfile~")
// - this app's debug files: "*_(echo,tkns,tree)"
// - filenames without a dot (indicating no file extension)
//
// As path separator, "/" is assumed, not os.Sep
// .
func wfnBuildContentityTree(path string, d fs.DirEntry, err error) error {
	var p *Contentity
	var e error
	// ROOT NODE ?
	if mustInitRoot() {
		var r *RootContentity
		if path != "." {
			println("wfnBuildContentityTree: "+
				"root path is not dot but instead:", path)
		}
		if pCFS.RootAbsPath() == "" {
			panic("wfnBuildContentityTree: nil ROOT")
		}
		r, e = NewRootContentity(pCFS.RootAbsPath())
		if e != nil || r == nil {
			// panic("wfnBuildContentityTree mustInitRoot NewRootContentityNord FAILED")
			return errors.New("wfnBuildContentityTree mustInitRoot NewRootContentityNord L77")
		}
		pCFS.rootNord = r
		r.MimeType = "dir"
		r.MType = "dir"
		// println("wfnBuildContentityTree: root node abs.FP:\n\t", p.AbsFP())
		var p *Contentity
		p = ((*Contentity)(r))
		pCFS.asSlice = append(pCFS.asSlice, p)
		pCFS.asMap[path] = p
		// println("ADDED TO MAP:", path)
		pCFS.nDirs = 1
		pCFS.nFiles = 0
		return nil
	}
	// Filter out hidden, esp'ly .git & .git* ;
	// note that "/" is assumed, not os.Sep ;
	// we issue no message, cos .git/* can be huge.
	if S.Contains(path, "/.git") {
		return nil
	}
	var reasonToReject string
	var ln = len(path)

	if S.HasPrefix(path, ".") || S.HasPrefix(path, "_") ||
		S.Contains(path, "/.") || S.Contains(path, "/_") {
		reasonToReject = "leading . or _ "
	} else if S.HasSuffix(path, "gtk") || S.HasSuffix(path, "gtr") {
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
	p, e = NewContentity(path) // FP.Join(pCFS.RootAbsPath(), path))
	if p == nil || e != nil {
		L.L.Warning("Rejecting (no Contentity): " + path)
		L.L.Error("Skipping this item")
		return nil
	}
	if /* p.PathAnalysis == nil || */ p.PathProps.Raw == "" {
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
		// println("================DIR ========")
		p.MimeType = "dir"
		p.MType = "dir"
	} else {
		pCFS.nFiles++
	}
	// println("ADDED TO MAP:", path)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}
