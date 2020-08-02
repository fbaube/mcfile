package mcfile

import (
	"errors"
	"fmt"

	"github.com/fbaube/gtoken"
	PU "github.com/fbaube/parseutils"
	SU "github.com/fbaube/stringutils"
	XM "github.com/fbaube/xmlmodels"
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
	fmt.Printf("    --> FileType<%s> MType<%v> \n", p.FileType(), p.MType)
	return p.
		st1a_Split().
		st1b_ProcessMetadata().
		st1c_GetCPR().
		st1d_MakeAFLfromCFL().
		st1e_PostMeta_notmkdn() // XML per format; HTML <head>
}

/*
type ContentitySections struct {
	Raw string // The entire input file
	// Text_raw + Meta_raw = Raw (maybe plus surrounding tags)
	Text_raw   string
	Meta_raw   string
	MetaFormat string
	MetaProps  SU.PropSet
}
*/
// st1a_Split is Step 1a: used to split the file into two parts -
// (header) meta and (body) text.
//
func (p *MCFile) st1a_Split() *MCFile {
	if p.HasError() {
		return p
	}
	switch p.FileType() {
	case "MKDN":
		i, e := SU.YamlMetadataHeaderLength(p.Raw)
		if e != nil {
			p.SetError(fmt.Errorf("yaml metadata header: %w", e))
			return p
		}
		if i == 0 {
			p.Text_raw = p.Raw
		} else {
			p.Meta_raw = p.Raw[:i]
			p.Text_raw = p.Raw[i:]
			/*
				println(
					"D=> === META ===\n", p.Meta_raw,
					"D=> === TEXT === \n", p.Text_raw,
					"D=> === End ===")
			*/
		}
	case "XML", "HTML":
		println("st1aa_PreMeta: XML/HTML...")
		// HTML, XHTML: Look for <html>, <head>, <body>
		//  XML (DITA): Look for...
		// topic: (title, shortdesc?, prolog?, body?)
		//   map: (topicmeta?, (topicref|keydef)*)
	}
	return p
}

// st1b_ProcessMetadata is Step 1b: used to process metadata.
//
func (p *MCFile) st1b_ProcessMetadata() *MCFile {
	if p.HasError() {
		return p
	}
	if p.Meta_raw == "" {
		return p
	}
	switch p.FileType() {
	case "XML", "HTML":
		// return p.TryXmlPreamble()
		println("st1a_PreMeta: XML/HTML TBS")
	case "MKDN":
		ps, e := SU.GetYamlMetadataAsPropSet(SU.TrimYamlMetadataDelimiters(p.Meta_raw))
		if e != nil {
			p.SetError(fmt.Errorf("yaml metadata: %w", e))
			return p
		}
		p.MetaProps = ps
	}
	return p
}

// st1c_GetCPR is Step 1c: Get ConcreteParseResults
func (p *MCFile) st1c_GetCPR() *MCFile {
	if p.HasError() {
		return p
	}
	if len(p.Raw) == 0 {
		p.Whine(p.OwnLogPfx + "st[1b] " + "Zero-length content")
		return p
	}
	var e error
	switch p.FileType() {
	case "MKDN":
		var pPR *PU.ConcreteParseResults_mkdn
		pPR, e = PU.GetConcreteParseResults_mkdn(p.Raw)
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			p.Blare(p.OwnLogPfx + e.Error())
			p.SetError(e)
			return p
		}
		p.CPR = pPR
		fmt.Printf("==> MKDNtokens: got %d \n", len(pPR.NodeList))
		// p.TallyTags()
		return p
	case "HTML":
		var pPR *PU.ConcreteParseResults_html
		pPR, e = PU.GetConcreteParseResults_html(p.Raw)
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			p.Blare(p.OwnLogPfx + e.Error())
			p.SetError(e)
			return p
		}
		p.CPR = pPR
		fmt.Printf("==> HTMLtokens: got %d \n", len(pPR.NodeList))
		// p.TallyTags()
		return p
	case "XML":
		var pPR *XM.ConcreteParseResults_xml
		pPR, e := XM.GetConcreteParseResults_xml(p.Raw)
		if e != nil {
			e = fmt.Errorf("XML tokenization failed: %w", e)
		}
		p.CPR = pPR
		fmt.Printf("==> XMLtokens: got %d \n", len(pPR.NodeList))
		return p
	default:
		println("ERROR st1b_GetCPR: bad file type:", p.FileType())
	}
	return p
}

// st1d_MakeAFLfromCFL is Step 1d:
// Make Abstract Flat List from Concrete Flat List
func (p *MCFile) st1d_MakeAFLfromCFL() *MCFile {
	if p.GetError() != nil {
		return p
	}
	var e error
	var errmsg string
	var GTs []*gtoken.GToken

	switch p.FileType() {
	case "MKDN":
		GTs, e = gtoken.DoGTokens_mkdn(p.CPR.(*PU.ConcreteParseResults_mkdn))
		if e != nil {
			p.SetError(fmt.Errorf("st1d: mkdn.GTs: %w", e))
		}
		p.GTokens = GTs
	case "HTML":
		GTs, e = gtoken.DoGTokens_html(p.CPR.(*PU.ConcreteParseResults_html))
		if e != nil {
			p.SetError(fmt.Errorf("st1d: html.GTs: %w", e))
		}
		p.GTokens = GTs
	case "XML":
		GTs, e = gtoken.DoGTokens_xml(p.CPR.(*XM.ConcreteParseResults_xml))
		if e != nil {
			e = fmt.Errorf("GToken-ization failed: %w", e)
		}
		if e != nil {
			errmsg = "st[1f] " + e.Error()
			p.Blare(p.OwnLogPfx + errmsg)
			p.SetError(e)
			return p
		}
		p.TallyTags()
		// fmt.Printf("==> Tags: %v \n", pGF.TagTally)
		// fmt.Printf("==> Atts: %v \n", pGF.AttTally)
		p.GTokens = GTs
	}
	// fmt.Printf("st1c_MakeAFLfromCFL: nGTokens: %d %d \n", len(p.GTokens), len(GTs))
	return p
}

// st1e_PostMeta_notmkdn is Step 1e (XML,HTML): XML per format; HTML <head>
func (p *MCFile) st1e_PostMeta_notmkdn() *MCFile {
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
		println("mcfl.st1.todo: SetMTypePerDoctypeFields")
	}
	return p
}
