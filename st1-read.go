package mcfile

import (
	"errors"
	"fmt"

	"github.com/fbaube/gtoken"
	PU "github.com/fbaube/parseutils"
	XU "github.com/fbaube/xmlutils"
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
//
func (p *Contentity) st1_Read() *Contentity {
	if p.GetError() != nil {
		return p
	}
	p.logStg = "1:"
	// p.L(LProgress, "Read")
	p.L(LInfo, "At entry: FileType<%s> MType<%v>", p.FileType(), p.MType)
	return p.
		st1a_ProcessMetadata().
		st1b_GetCPR().
		st1c_MakeAFLfromCFL().
		st1d_PostMeta_notmkdn() // XML per format; HTML <head>
}

// st1a_ProcessMetadata is Step 1a: used to process metadata.
// Note that for Markdown, YAML metadata parsing is currently
// done during initial file content analysis.
//
func (p *Contentity) st1a_ProcessMetadata() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1a"
	metaRaw := p.GetSpan(p.Meta)
	// textRaw := p.GetSpan(p.Text)
	if metaRaw == "" {
		p.L(LInfo, "No metadata encountered")
		return p
	}
	switch p.FileType() {
	case "XML", "HTML":
		ft := p.FileType()
		p.L(LDbg, "MetaPos:%d MetaRaw(): %s",
			p.Meta.Beg.Pos, metaRaw)
		if p.Meta.Beg.Pos != 0 {
			var e error
			var ct int
			p.L(LProgress, "Doing "+ft)
			if ft == "HTML" {
				var pPR *PU.ParserResults_html
				pPR, e = PU.GenerateParserResults_html(metaRaw)
				ct = len(pPR.NodeSlice)
				p.ParserResults = pPR
			}
			if ft == "XML" {
				var pPR *XU.ParserResults_xml
				pPR, e = XU.GenerateParserResults_xml(metaRaw)
				ct = len(pPR.NodeSlice)
				p.ParserResults = pPR
			}
			if e != nil {
				p.L(LError, "%s tokenization failed: %w", ft, e)
				p.ParserResults = nil
			}
			p.L(LOkay, "%s tokens: got %d \n", ft, ct)
			p.L(LWarning, "TODO: Do something with XML/HTML metadata")
			return p
		}
	case "MKDN":
		p.L(LWarning, "TODO: Do something with YAML metadata")
	}
	return p
}

// st1b_GetCPR is Step 1b: Generate ParserResults
//
func (p *Contentity) st1b_GetCPR() *Contentity {
	if p.HasError() {
		return p
	}
	textRaw := p.GetSpan(p.Text)
	p.logStg = "1b"
	if len(textRaw) == 0 {
		p.L(LWarning, "Zero-length content")
		p.SetError(errors.New("no content"))
		return p
	}
	var e error
	switch p.FileType() {
	case "MKDN":
		var pPR *PU.ParserResults_mkdn
		pPR, e = PU.GenerateParserResults_mkdn(textRaw)
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
		pPR, e = PU.GenerateParserResults_html(textRaw)
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
		var pPR *XU.ParserResults_xml
		pPR, e := XU.GenerateParserResults_xml(textRaw)
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

// st1c_MakeAFLfromCFL is Step 1c:
// Make Abstract Flat List from Concrete Flat List
//
func (p *Contentity) st1c_MakeAFLfromCFL() *Contentity {
	if p.GetError() != nil {
		return p
	}
	p.logStg = "1c"
	var e error
	// var errmsg string
	var GTs []*gtoken.GToken

	fmt.Fprintln(p.GTokensOutput, "=== Input:", p.AbsFP())

	switch p.FileType() {
	case "MKDN":
		var pCPR_M *PU.ParserResults_mkdn
		if nil == p.ParserResults {
			p.L(LError, "ParserResults are nil")
		}
		pCPR_M = p.ParserResults.(*PU.ParserResults_mkdn)
		pCPR_M.DiagDest = p.GTokensOutput
		GTs, e = gtoken.DoGTokens_mkdn(pCPR_M)
		if e != nil {
			p.SetError(fmt.Errorf("st1d: mkdn.GTs: %w", e))
		}
		// p.GTokens = GTs
		// Compress out nil GTokens
		p.GTokens = make([]*gtoken.GToken, 0)
		for _, GT := range GTs {
			if GT != nil {
				p.GTokens = append(p.GTokens, GT)
			}
		}
	case "HTML":
		var pCPR_H *PU.ParserResults_html
		pCPR_H = p.ParserResults.(*PU.ParserResults_html)
		pCPR_H.DiagDest = p.GTokensOutput
		GTs, e = gtoken.DoGTokens_html(pCPR_H)
		if e != nil {
			p.SetError(fmt.Errorf("st1d: html.GTs: %w", e))
		}
		p.GTokens = GTs
	case "XML":
		var pCPR_X *XU.ParserResults_xml
		pCPR_X = p.ParserResults.(*XU.ParserResults_xml)
		pCPR_X.DiagDest = p.GTokensOutput
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
	fmt.Fprintln(p.GTokensOutput, "=== Output:")
	for i, pGtkn := range p.GTokens {
		if pGtkn != nil {
			fmt.Fprintf(p.GTokensOutput, "[%02d:L%d] %s \n", i, p.Level(), pGtkn.String())
		}
	}
	// fmt.Printf("st1c_MakeAFLfromCFL: nGTokens: %d %d \n", len(p.GTokens), len(GTs))
	return p
}

// st1d_PostMeta_notmkdn is Step 1d (XML,HTML): XML per format; HTML <head>
//
func (p *Contentity) st1d_PostMeta_notmkdn() *Contentity {
	p.logStg = "1d"
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
		println("mcfl.st1.todo: SetMTypePerDoctypeFields:")
		fmt.Println(p.AnalysisRecord.String())

	}
	return p
}
