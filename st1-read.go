package mcfile

import (
	"fmt"
	"errors"
	"github.com/fbaube/gparse"
	PU "github.com/fbaube/parseutils"
)

// - "XML"
// - - (§1) Use stdlib `encoding/xml` to get `[]xml.Token`
// - - (§1) Convert `[]xml.Token` to `[]gparse.GToken`
// - "MKDN"
// - - (§1) Use `yuin/goldmark` to get tree of `yuin/goldmark/ast/Node`
// - - (§1) From each Node make a `MkdnToken` (in a list?) incl. `GToken` and `GTag`
// - "HTML"
// - - (§1) Use `golang.org/x/net/html` to get a tree of `html.Node`
// - - (§1) From each Node make a `HtmlToken` (in a list?) incl. `GToken` and `GTag`

// GetParseTokenization_Xml v GetParseCST_nonXml
// GetNodelistFromCST_NonXml
// GetGTokensFromParseTokenization_Xml v
// GetGTokensFromNodelist_NonXml

// st1_Read reads in the file and does what is
// needed to end up with a list of `GToken`s.
// - DoPreMeta()
// - DoTokenize()
// // - DoPostMeta()
func (p *MCFile) st1_Read() *MCFile {
	if p.GetError() != nil {
		return p
	}
	println("--> (1) Read")
	fmt.Printf("--> FileType<%s> MType<%v> \n", p.FileType(), p.MType)
	return p.
		st1a_PreMeta().
		st1b_GetCPR().
		st1c_MakeAFLfromCFL().
		st1d_PostMeta_notmkdn() // XML per format; HTML <head>
}

// st1a_PreMeta is Step 1a: used when metadata can easily be
// separated from content, e.g. YAML frontmatter in MDITA-XP.
//
func (p *MCFile) st1a_PreMeta() *MCFile {
	if p.GetError() != nil {
		return p
	}
	switch p.FileType() {
	case "XML", "HTML":
		return p.TryXmlPreamble()
	case "MKDN":
		return p.GetYamlHeader()
	}
	return p
}

// st1b_GetCPR is Step 1b: Get ContreteParseResults
func (p *MCFile) st1b_GetCPR() *MCFile {
	if p.GetError() != nil {
		return p
	}
	if len(p.CheckedContent.Raw) == 0 {
		p.Whine(p.OLP + "st[1b] " + "Zero-length content")
		return p
	}
	var e error
	switch p.FileType() {
	case "MKDN":
		var pPR *PU.ConcreteParseResults_mkdn
		pPR, e = PU.GetParseResults_mkdn(p.CheckedContent.Raw)
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			p.Blare(p.OLP + e.Error())
			p.SetError(e)
			return p
		}
		p.CPR = pPR
		fmt.Printf("==> MKDNtokens: got %d \n", len(pPR.NodeList))
		// p.TallyTags()
		return p
	case "HTML":
		var pPR *PU.ConcreteParseResults_html
		pPR, e = PU.GetParseResults_html(p.CheckedContent.Raw)
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			p.Blare(p.OLP + e.Error())
			p.SetError(e)
			return p
		}
		p.CPR = pPR
		fmt.Printf("==> HTMLtokens: got %d \n", len(pPR.NodeList))
		// p.TallyTags()
		return p
	case "XML":
		var pPR *PU.ConcreteParseResults_xml
		pPR, e := PU.GetParseResults_xml(p.CheckedContent.Raw)
		if e != nil {
			e = fmt.Errorf("XML tokenization failed: %w", e)
		}
		p.CPR = pPR
		fmt.Printf("==> XMLtokens: got %d \n", len(pPR.NodeList))
		return p
	default:
		panic("st1b_GetCPR: bad file type: " + p.FileType())
	}
}


// st1c_MakeAFLfromCFL is Step 1d:
// Make Abstract Flat List from Concrete Flat List
func (p *MCFile) st1c_MakeAFLfromCFL() *MCFile {
	if p.GetError() != nil {
		return p
	}
	var e error
	var errmsg string
	var GTs []*gparse.GToken

		switch p.FileType() {
		case "MKDN":
			GTs, e = gparse.DoGTokens_mkdn(p.CPR.(*PU.ConcreteParseResults_mkdn))
			if e != nil {
				p.SetError(fmt.Errorf("st1d: mkdn.GTs: %w", e))
			}
			p.GTokens = GTs
			return p
		case "HTML":
			GTs, e = gparse.DoGTokens_html(p.CPR.(*PU.ConcreteParseResults_html))
			if e != nil {
				p.SetError(fmt.Errorf("st1d: html.GTs: %w", e))
			}
			p.GTokens = GTs
			return p
		case "XML":
			GTs, e = gparse.DoGTokens_xml(p.CPR.(*PU.ConcreteParseResults_xml))
			if e != nil {
				e = fmt.Errorf("GToken-ization failed: %w", e)
			}
			if e != nil {
				errmsg = "st[1f] " + e.Error()
				p.Blare(p.OLP + errmsg)
				p.SetError(e)
				return p
			}
			p.TallyTags()
			// fmt.Printf("==> Tags: %v \n", pGF.TagTally)
			// fmt.Printf("==> Atts: %v \n", pGF.AttTally)
			p.GTokens = GTs
		}
		return p
	}

// st1d_PostMeta_notmkdn is Step 1c (XML,HTML): XML per format; HTML <head>
func (p *MCFile) st1d_PostMeta_notmkdn() *MCFile {
	switch p.FileType() {
	case "MKDN":
		// Markdown YAML metadata was processed in step st1a
		return p
	case "HTML": /*
		var pPR *PU.ConcreteParseResults_html
		pPR = p.CPR.(*PU.ConcreteParseResults_html)
		z := pPR. */
		// Inside <head>: <meta> <title> <base> <link> <style>
		// See also: https://gist.github.com/lancejpollard/1978404
		return p
	case "XML":
		// [Lw]DITA stuff, ?DublinCore
	}
	return p
}
