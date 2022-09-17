package mcfile

import (
	"fmt"
)

func (p *Contentity) SetErrMsg(s string) {
	p.Err = fmt.Errorf("[%02d:%s] %s", p.logIdx, p.logStg, s)
}

func (p *Contentity) WrapError(s string, e error) {
	p.Err = fmt.Errorf("[%02d:%s] %s: %w", p.logIdx, p.logStg, s, e)
}
