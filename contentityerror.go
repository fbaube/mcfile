package mcfile

import (
	"errors"
	"fmt"
	"io/fs"
)

// ContentityError is
// Contentity + SrcLoc (in source code) + 
// PathError struct { Op, Path string; Err error }
//
// Maybe use the format pkg.filename.methodname.Lnn
//
// In code where package `mcfile` is not available,
// try a fileutils.PathPropsError 
//
type ContentityError struct {
	PE     fs.PathError
	SrcLoc string
	*Contentity
}

func WrapAsContentityError(e error, op string, cty *Contentity, srcLoc string) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = e
	ce.PE.Op  = op
	ce.SrcLoc = srcLoc 
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.PathProps.AbsFP.S()
	}
	return ce
}

func NewContentityError(ermsg string, op string, cty *Contentity, srcLoc string) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = errors.New(ermsg)
	ce.PE.Op  = op
	ce.SrcLoc = srcLoc
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.PathProps.AbsFP.S()
	}
	return ce
}

func (ce ContentityError) Error() string {
	return ce.String()
}

func (ce *ContentityError) String() string {
	var s string
	s = fmt.Sprintf("%s(%s): %s", ce.PE.Op, ce.PE.Path, ce.PE.Err.Error())
	if ce.SrcLoc != "" {
		s += fmt.Sprintf(" (in %s)", ce.SrcLoc)
	}
	return s
}
