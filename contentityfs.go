package mcfile

import (
	"io/fs"
	L "github.com/fbaube/mlog"
)

// ContentityFS is an instance of an [fs.FS] where every
// node is an [mcfile.Contentity].
//
// Note that directories ARE included in the tree, because
// the instances of [orderednodes.Nord] in each [Contentity]
// must properly interconnect in forming a complete tree.
//
// Note that the file system is stored as a tree AND as a
// slice AND as a map. If any of these is modified without 
// also modifying the others to match, there WILL be problems.
// For that reason, we use unexported instance variables that
// are accessible only via getters. It ain't bulletproof tho. 
// .
type ContentityFS struct {
	// FS will be an [os.DirFS]
	fs.FS
	rootAbsPath string
	rootNord    *RootContentity
	asSlice     []*Contentity
	// The string is the relative filepath w.r.t. the rootAbsPath
	asMap         map[string]*Contentity
	nFiles, nDirs int
}

func (p *ContentityFS) Size() int {
	// /* Not init'z ?
	if p.asMap != nil && len(p.asSlice) != len(p.asMap) {
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
     L.L.Warning("mcm.ContentityFS.DoForEvery: not implemented")
}

