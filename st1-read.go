package mcfile

import (
	"errors"
	"fmt"
	"os"

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
//
func (p *Contentity) st1_Read() *Contentity {
	if p.GetError() != nil {
		return p
	}
	p.logStg = "1:"
	// p.L(LProgress, "Read")
	p.L(LInfo, "At entry: FileType<%s> MType<%v>", p.FileType(), p.MType)
	return p.
		st1a_Split_mkdn().
		st1b_ProcessMetadata().
		st1c_GetCPR().
		st1d_MakeAFLfromCFL().
		st1e_PostMeta_notmkdn() // XML per format; HTML <head>
}

/*
type ContentitySections is struct {
	Raw string // The entire input file
	// TextRaw() + MetaRaw() = Raw (maybe plus surrounding tags)
	TextRaw()   string
	MetaRaw()   string
	MetaFormat string
	MetaProps  SU.PropSet
}
*/

// st1a_Split is Step 1a: used to split the file into two parts:
// meta (i.e. header) and text (i.e. body) However for XML & HTML,
// this has already been done in Peek, so this stage is for Markdown only.
//
func (p *Contentity) st1a_Split_mkdn() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1a"
	switch p.FileType() {
	case "MKDN":
		p.L(LDbg, "MkdnSxns: len<%d> root<%s> meta<%s> text<%s> \n",
			len(p.ContentityStructure.Raw), p.Root.String(), p.Meta.String(), p.Text.String())
		ln, e := SU.YamlMetadataHeaderLength(p.Raw)
		if e != nil {
			p.SetError(fmt.Errorf("yaml metadata header: %w", e))
			return p
		}
		p.L(LWarning, "Split_mkdn() hdr-len %d FIXME", ln)
		/*
			p.TextRaw() = p.Raw
			if i != 0 {
				p.MetaRaw() = p.Raw[:i]
				p.TextRaw() = p.Raw[i:]
				println(
					"D=> === BODY ===\n", p.Raw,
					"D=> === META ===\n", p.MetaRaw(),
					"D=> === TEXT === \n", p.TextRaw(),
					"D=> === End ===")
			}
		*/
	}
	return p
}

// st1b_ProcessMetadata is Step 1b: used to process metadata.
//
func (p *Contentity) st1b_ProcessMetadata() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1b"
	if p.MetaRaw() == "" && p.Meta.Beg.Pos == 0 {
		p.L(LInfo, "No metadata encountered")
		return p
	}
	switch p.FileType() {
	case "XML", "HTML":
		ft := p.FileType()
		p.L(LDbg, "MetaPos:%d MetaRaw(): %s",
			p.Meta.Beg.Pos, p.MetaRaw())
		if p.Meta.Beg.Pos != 0 {
			var e error
			var ct int
			p.L(LProgress, "Doing "+ft)
			if ft == "HTML" {
				var pPR *PU.ParserResults_html
				pPR, e = PU.GenerateParserResults_html(p.MetaRaw())
				ct = len(pPR.NodeSlice)
				p.ParserResults = pPR
			}
			if ft == "XML" {
				var pPR *XM.ParserResults_xml
				pPR, e = XM.GenerateParserResults_xml(p.MetaRaw())
				ct = len(pPR.NodeSlice)
				p.ParserResults = pPR
			}
			if e != nil {
				p.L(LError, "%s tokenization failed: %w", ft, e)
				p.ParserResults = nil
			}
			p.L(LOkay, "%s tokens: got %d \n", ft, ct)
			return p
		}
	case "MKDN":
		ps, e := SU.GetYamlMetadataAsPropSet(
			SU.TrimYamlMetadataDelimiters(p.MetaRaw()))
		if e != nil {
			p.SetError(fmt.Errorf("yaml metadata: %w", e))
			return p
		}
		if len(p.TextRaw()) == 0 {
			p.L(LWarning, "NO MKDN in st1b_ProcessMetadata")
		}
		p.MetaProps = ps
	}
	return p
}

// st1c_GetCPR is Step 1c: Generate ParserResults
//
func (p *Contentity) st1c_GetCPR() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1c"
	if len(p.TextRaw()) == 0 {
		p.L(LWarning, "Zero-length content")
		return p
	}
	var e error
	switch p.FileType() {
	case "MKDN":
		var pPR *PU.ParserResults_mkdn
		pPR, e = PU.GenerateParserResults_mkdn(p.TextRaw())
		if e != nil {
			e = errors.New("st[1c] " + e.Error())
			// // p.Blare(p.OwnLogPfx + e.Error())
			p.SetError(e)
			p.L(LError, "Failure in GenerateParserResults_mkdn")
			return p
		}
		if pPR == nil {
			p.L(LError, "No Error but nil ParserResults")
		}
		p.ParserResults = pPR
		p.L(LOkay, "MKDN tokens: got %d", len(pPR.NodeSlice))
		// p.TallyTags()
		return p
	case "HTML":
		var pPR *PU.ParserResults_html
		pPR, e = PU.GenerateParserResults_html(p.TextRaw())
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			// p.Blare(p.OwnLogPfx + e.Error())
			p.SetError(e)
			p.L(LError, "Failure in GenerateParserResults_html")
			return p
		}
		p.ParserResults = pPR
		p.L(LOkay, "HTML tokens: got %d", len(pPR.NodeSlice))
		// p.TallyTags()
		return p
	case "XML":
		var pPR *XM.ParserResults_xml
		pPR, e := XM.GenerateParserResults_xml(p.TextRaw())
		if e != nil {
			e = fmt.Errorf("XML tokenization failed: %w", e)
			p.L(LError, "Failure in GenerateParserResults_xml")
		}
		p.ParserResults = pPR
		p.L(LOkay, "XML tokens: got %d \n", len(pPR.NodeSlice))
		return p
	default:
		p.L(LError, "st1b_GetCPR: bad file type: "+p.FileType())
	}
	return p
}

// st1d_MakeAFLfromCFL is Step 1d:
// Make Abstract Flat List from Concrete Flat List
//
func (p *Contentity) st1d_MakeAFLfromCFL() *Contentity {
	if p.GetError() != nil {
		return p
	}
	p.logStg = "1d"
	var e error
	// var errmsg string
	var GTs []*gtoken.GToken

	p.L(LDbg, "ParserResults: %T", p.ParserResults)

	switch p.FileType() {
	case "MKDN":
		var pCPR_M *PU.ParserResults_mkdn
		p.L(LWarning, "BEFORE BARF")
		if nil == p.ParserResults {
			p.L(LError, "BARF ON NIL ParserResults")
		}
		pCPR_M = p.ParserResults.(*PU.ParserResults_mkdn)
		p.L(LOkay, "!AFTER BARF")
		if p.GTokensOutput != nil {
			pCPR_M.DumpDest = p.GTokensOutput
		} else {
			pCPR_M.DumpDest = os.Stdout
		}
		GTs, e = gtoken.DoGTokens_mkdn(pCPR_M)
		if e != nil {
			p.SetError(fmt.Errorf("st1d: mkdn.GTs: %w", e))
		}
		// p.GTokens = GTs
		// Compress out nil GTokens ?
		p.GTokens = make([]*gtoken.GToken, 0)
		for _, GT := range GTs {
			if GT != nil {
				p.GTokens = append(p.GTokens, GT)
			}
		}
	case "HTML":
		var pCPR_H *PU.ParserResults_html
		pCPR_H = p.ParserResults.(*PU.ParserResults_html)
		if p.GTokensOutput != nil {
			pCPR_H.DumpDest = p.GTokensOutput
		} else {
			pCPR_H.DumpDest = os.Stdout
		}
		GTs, e = gtoken.DoGTokens_html(pCPR_H)
		if e != nil {
			p.SetError(fmt.Errorf("st1d: html.GTs: %w", e))
		}
		p.GTokens = GTs
	case "XML":
		var pCPR_X *XM.ParserResults_xml
		pCPR_X = p.ParserResults.(*XM.ParserResults_xml)
		if p.GTokensOutput != nil {
			pCPR_X.DumpDest = p.GTokensOutput
		} else {
			pCPR_X.DumpDest = os.Stdout
		}
		GTs, e = gtoken.DoGTokens_xml(pCPR_X)
		if e != nil {
			e = fmt.Errorf("GToken-ization failed: %w", e)
		}
		if e != nil {
			// errmsg = "st[1f] " + e.Error()
			// p.Blare(p.OwnLogPfx + errmsg)
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
//
func (p *Contentity) st1e_PostMeta_notmkdn() *Contentity {
	switch p.FileType() {
	case "MKDN":
		// Markdown YAML metadata was processed in step st1a
		return p
	case "HTML": /*
			var pPR *PU.ParserResults_html
			pPR = p.CPR.(*PU.ParserResults_html)
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
