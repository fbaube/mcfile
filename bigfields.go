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

func (p *MCFile) PushBigFields() BigFields {
	bf := new(BigFields)
	bf.hed = p.Meta_raw
	bf.bod = p.Text_raw 
	// bf.gtox = p.GTokens
	// bf.gtags = p.GTags
	// p.Header = new(Header)
	p.Meta_raw = "[raw.header]"
	p.Text_raw = "[raw.text]"
	// ?? p.Raw = "[et.raw]"
	// p.GTags = nil
	return *bf
}

func (p *MCFile) PopBigFields(BF BigFields) {
	p.Meta_raw = BF.hed
	p.Text_raw = BF.bod
	// p.gparses = BF.gtox
	// p.GTags = BF.gtags
}
