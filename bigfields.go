package mcfile

import (
	"github.com/fbaube/gtoken"
	// "github.com/dimchansky/utfbom"
)

type BigFields struct {
	hed  string // Header
	bod  string
	gtox []*gtoken.GToken
	// gtags []*gtree.GTag
}

/*
func (p *MCFile) PushBigFields() BigFields {
	bf := new(BigFields)
	bf.hed = p.MetaRaw()
	bf.bod = p.TextRaw()
	// bf.gtox = p.GTokens
	// bf.gtags = p.GTags
	// p.Header = new(Header)
	p.MetaRaw() = "[raw.header]"
	p.TextRaw() = "[raw.text]"
	// ?? p.Raw = "[et.raw]"
	// p.GTags = nil
	return *bf
}

func (p *MCFile) PopBigFields(BF BigFields) {
	p.MetaRaw() = BF.hed
	p.TextRaw() = BF.bod
	// p.gparses = BF.gtox
	// p.GTags = BF.gtags
}
*/
