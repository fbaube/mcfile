package mcfile

import (
	"fmt"
	L "github.com/fbaube/mlog"
)

func (p *Contentity) SetError(s string) {
     	var e error
	e = fmt.Errorf("[F%02d:S%s] %s", p.logIdx, p.logStg, s)
	p.Errer.Err = e
	L.L.Error(e.Error())
}

func (p *Contentity) WrapError(s string, e error) {
     	var e2 error
	e2 = fmt.Errorf("[F%02d:S%s] %s: %w", p.logIdx, p.logStg, s, e)
	p.Errer.Err = e2
	L.L.Error(e2.Error())
}
