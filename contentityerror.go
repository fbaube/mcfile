package mcfile

import (
	"errors"
	"fmt"
	"io/fs"
)

// ContentityError is
// Contentity + Func (in source code) +
// PathError struct { Op, Path string; Err error }
//
type ContentityError struct {
	PE   fs.PathError
	Func string
	*Contentity
}

func WrapContentityError(e error, op string, cty *Contentity, fnc string) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = e
	ce.PE.Op = op
	ce.Func = fnc
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.PathProps.AbsFP.S()
	}
	return ce
}

func NewContentityError(ermsg string, op string, cty *Contentity, fnc string) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = errors.New(ermsg)
	ce.PE.Op = op
	ce.Func = fnc
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.PathProps.AbsFP.S()
	}
	return ce
}

func (ce *ContentityError) String() string {
	var s string
	s = fmt.Sprintf("%s(%s): %s", ce.PE.Op, ce.PE.Path, ce.PE.Err.Error())
	if ce.Func != "" {
		s += fmt.Sprintf("(in %s)", ce.Func)
	}
	return s
}
