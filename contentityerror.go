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
	PE fs.PathError
	*Contentity
}

func (p *Contentity) SetErrMsg(s string) {
	p.Err = fmt.Errorf("[%02d:%s] %s", p.logIdx, p.logStg, s)
}

func (p *Contentity) SetErrWrap(s string, e error) {
	p.Err = fmt.Errorf("[%02d:%s] %s: %w", p.logIdx, p.logStg, s, e)
}

func WrapAsContentityError(e error, op string, cty *Contentity) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = e
	ce.PE.Op = op
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.PathProps.AbsFP.S()
	}
	return ce
}

func NewContentityError(ermsg string, op string, cty *Contentity) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = errors.New(ermsg)
	ce.PE.Op = op
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
	return s
}
