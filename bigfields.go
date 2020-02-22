package mcfile

import (
	"github.com/fbaube/gtoken"
	// "github.com/dimchansky/utfbom"
)

type BigFields struct {
	hed  *Header
	bod  string
	gtox []*gtoken.GToken
	// gtags []*gtree.GTag
}

func (p *MCFile) PushBigFields() BigFields {
	bf := new(BigFields)
	bf.hed = p.Header
	bf.bod = p.Body
	// bf.gtox = p.GTokens
	// bf.gtags = p.GTags
	p.Header = new(Header)
	p.Header.HedRaw = "[raw.header]"
	p.Body = "[raw.body]"
	// ?? p.Raw = "[et.raw]"
	// p.GTags = nil
	return *bf
}

func (p *MCFile) PopBigFields(BF BigFields) {
	p.Header = BF.hed
	p.Body = BF.bod
	// p.gparses = BF.gtox
	// p.GTags = BF.gtags
}
