package mcfile

import (
	"io/fs"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
        CT "github.com/fbaube/ctoken"
	assert "github.com/lainio/err2/assert")

// ContentityFS is an instance of an [fs.FS] where every
// node is an [mcfile.Contentity].
//
// Note that directories ARE included in the tree, because
// the instances of [orderednodes.Nord] in each [Contentity]
// must properly interconnect in forming a complete tree.
//
// Note that the file system is stored as a tree AND as a slice AND as a
// map. If any of these is modified without also modifying the others to
// match, there WILL be problems. For that reason, [asSlice] and [asMap]
// are unexported instance variables that are accessible only via getters.
//
// It ain't bulletproof tho. In any case, users of a ContentityFS should
// feel free to use the functions on the embedded [Nord] ordered nodes.
// .
type ContentityFS struct {
	// FS will be set from func [os.DirFS]
	fs.FS
	rootAbsPath string
	rootNord    *RootContentity
	asSlice     []*Contentity
	// The string is the relative filepath w.r.t. 
	// the rootAbsPath. But does this index into 
	// the tree of Nord's or into the slice ?
	asMap         map[string]*Contentity
	nItems, nFiles, nDirs int
}

func (p *ContentityFS) ItemCount() int {
     return p.Size()
}

func (p *ContentityFS) Size() int {
	// /* Not init'lzd ?
	if p.asMap != nil &&  p.asSlice != nil &&
	   len(p.asSlice) != len(p.asMap) {
		L.L.Error("contentityfs size mismatch (slice &d, map %d)",
			len(p.asSlice), len(p.asMap))
	}
	// */
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

func (p *ContentityFS) RootAbsPath() string {
	return p.rootAbsPath
}

func (p *ContentityFS) AsSlice() []*Contentity {
	var z []*Contentity
	z = p.asSlice
	return z // p.AsSlice
}

func (p *ContentityFS) DoForEvery(stgprocsr ContentityStage) {
     L.L.Warning("mcm.ContentityFS.DoForEvery: not implemented")
}

func (p *ContentityFS) mustInitRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(p.asSlice) == 0 && len(p.asMap) == 0)
	didDoInit = (len(p.asSlice) > 0  && len(p.asMap) > 0)
	if !(needsInit || didDoInit) {
		panic("mustInitRoot: illegal state")
	}
	return needsInit
}

func (p *ContentityFS) doInitRoot() error {
	var pRC *RootContentity
	var e error 
	assert.NotEmpty(p.RootAbsPath())
	/* if p.RootAbsPath() == "" {
		panic("wfnBuildContentityTree: no ROOT")
	} */
	pRC, e = NewRootContentity(p.RootAbsPath())
	if e != nil || pRC == nil {
		return &fs.PathError{Op:"WalkFn.NewRootContentity",
		Err:e,Path:p.RootAbsPath()}
	}
	// assert.That(pRC.IsDir()) SHOULD NOT FAIL, BUT DID
	// Assign to globals (i.e. package vars)
	p.rootNord = pRC
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
	p.asSlice = append(p.asSlice, pC)
	p.asMap[p.RootAbsPath()] = pC
	// L.L.Warning("ADDED TO MAP L84: " + p.RootAbsPath())
	p.nDirs = 1
	p.nFiles = 0
	return nil // NOT pRC! This is a walker func 
}
