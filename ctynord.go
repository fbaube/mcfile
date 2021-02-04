package mcfile

import (
	ON "github.com/fbaube/orderednodes"
)

// Ignore https://godoc.org/golang.org/x/net/html#Node

// FilePropsNord is an Ordered Propertied Path node:
// NOT ONLY the child nodes have a specific specified order
// BUT ALSO each node has a filepath plus the file properties.
// This means Pthat every Parent node is a directory.
//
// It also means we can use the redundancy to do a lot of error checking.
// Also we can use fields of seqId's to store parent and kid seqId's,
// adding yet another layer of error checking and simplified access.
//
type ContentityNord struct {
	ON.Nord
	// FU.PathProps
}

// Available to ensure that assignments to/from root node are explicit.
type RootContentityNord ContentityNord

type norderCreationState struct {
	nexSeqID int // reset to 0 when doing another tree ?
	rootPath string
	// summaryString StringFunc
}

var pNCS *norderCreationState = new(norderCreationState)

func NewRootContentityNord(rootPath string /*,smryFunc StringFunc*/) *ContentityNord {
	p := new(ContentityNord)
	p.Nord = *ON.NewRootNord(rootPath, nil)

	println("NewRootContentityNord:", p.AbsFP())
	return p
}

func NewContentityNord(aPath string) *ContentityNord {
	p := new(ContentityNord)
	p.Nord = *ON.NewNord(aPath)
	if aPath == "" {
		println("newNord: missing path")
	}

	return p
}
