package mcfile

import (
	"fmt"
	"io"

	"github.com/fbaube/gtoken"
	PU "github.com/fbaube/parseutils"
	SU "github.com/fbaube/stringutils"
L "github.com/fbaube/mlog"
	XU "github.com/fbaube/xmlutils"
	// L "github.com/fbaube/mlog"
)

// GetParseTokenization_Xml v GetParseCST_nonXml
// GetNodelistFromCST_NonXml
// GetGTokensFromParseTokenization_Xml v
// GetGTokensFromNodelist_NonXml

// st1_Read reads in the file and does what is
// needed to end up with a list of `GToken`s.
//
// Summary of processing per Contentity type:
// "XML"
//   - (§1) Use stdlib [encoding/xml] to get slice of [XU.XToken]
//   - (§1) Convert [XU.XToken] to [gparse.GToken]
//
// "MKDN"
//   - (§1) Use [yuin/goldmark] to get tree of [yuin/goldmark/ast/Node]
//   - (§1) From each Node make a [MkdnToken] (in a list?) incl. [GToken] and [GTag]
//
// "HTML"
//   - (§1) Use [golang.org/x/net/html] to get a tree of [html.Node]
//   - (§1) From each Node make a [HtmlToken] (in a list?) incl. [GToken] and [GTag]
//
// .
func (p *Contentity) st1_Read() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "11"
	p.L(LProgress, "=== 11:Read ===")
	p.L(LInfo, "@entry: MarkupType<%s> MType<%s>",
		p.MarkupType(), p.MType)
	return p.
		st1a_ProcessMetadata().
		st1b_GetCPR().
		st1c_MakeAFLfromCFL().
		st1d_PostMeta_notmkdn() // XML per format; HTML <head>
}

// st1a_ProcessMetadata processes metadata.
// Note that for Markdown, YAML metadata parsing is
// already done during initial file content analysis.
// .
func (p *Contentity) st1a_ProcessMetadata() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1a"
	metaRaw := p.Meta.GetSpanOfString(p.FSItem.TypedRaw.S())
	if metaRaw == "" {
		p.L(LInfo, "No metadata found")
		return p
	}
	switch mut := p.MarkupType(); mut {
	case SU.MU_type_XML, SU.MU_type_HTML:
		p.L(LDbg, "MetaPos:%d MetaRaw(): %s",
			p.Meta.Beg.Pos, metaRaw)
		if p.Meta.Beg.Pos != 0 {
			var e error
			var ct int
			p.L(LProgress, "Doing "+string(mut))
			if mut == SU.MU_type_HTML {
				var pPR *PU.ParserResults_html
				pPR, e = PU.GenerateParserResults_html(metaRaw)
				ct = len(pPR.NodeSlice)
				p.ParserResults = pPR
			}
			if mut == SU.MU_type_XML {
				var pPR *XU.ParserResults_xml
				pPR, e = XU.GenerateParserResults_xml(metaRaw)
				ct = len(pPR.NodeSlice)
				p.ParserResults = pPR
			}
			if e != nil {
				p.L(LError, "%s tokenization failed: %w", mut, e)
				p.ParserResults = nil
			}
			p.L(LOkay, "%s tokens: got %d", mut, ct)
			p.L(LWarning, "TODO: Do sthg with XML/HTML metadata")
			return p
		}
	case SU.MU_type_MKDN:
		p.L(LWarning, "TODO: Do sthg with YAML metadata")
	}
	return p
}

// st1b_GetCPR generates Concrete Parser Results
// .
func (p *Contentity) st1b_GetCPR() *Contentity {
	if p.HasError() {
		return p
	}
	textRaw := p.Text.GetSpanOfString(p.FSItem.TypedRaw.S())
	if textRaw == "" {
		p.L(LWarning, "Lame hack in st1-read L105")
		textRaw = p.FSItem.TypedRaw.S()
	}
	p.logStg = "1b"
	if len(textRaw) == 0 {
		p.L(LWarning, "Zero-length content")
		p.SetError("no content")
		return p
	}
	var e error
	switch p.MarkupType() {
	case SU.MU_type_MKDN:
		var pPR *PU.ParserResults_mkdn
		pPR, e = PU.GenerateParserResults_mkdn(textRaw)
		if e != nil {
			p.WrapError("GenerateParserResults_mkdn", e)
			return p
		}
		if pPR == nil {
			p.SetError("nil ParserResults_mkdn")
		}
		p.ParserResults = pPR
		p.L(LOkay, "MKDN tokens: got %d", len(pPR.NodeSlice))
		// p.TallyTags()
		return p
	case SU.MU_type_HTML:
		var pPR *PU.ParserResults_html
		pPR, e = PU.GenerateParserResults_html(textRaw)
		if e != nil {
			p.WrapError("GenerateParserResults_html", e)
			return p
		}
		p.ParserResults = pPR
		p.L(LOkay, "HTML tokens: got %d", len(pPR.NodeSlice))
		/* DBG
		for i, pp := range pPR.NodeSlice {
			p.L(LWarning, "[%d] %+v \n", i, *pp)
		}
		*/
		// p.TallyTags()
		return p
	case SU.MU_type_XML:
		var pPR *XU.ParserResults_xml
		pPR, e := XU.GenerateParserResults_xml(textRaw)
		if e != nil {
			p.WrapError("GenerateParserResults_xml", e)
		}
		p.ParserResults = pPR
		p.L(LOkay, "XML tokens: got %d \n", len(pPR.NodeSlice))
		return p
	default:
		p.SetError("bad file markup type: " +
			string(p.MarkupType()))
	}
	return p
}

// st1c_MakeAFLfromCFL is Step 1c:
// Make Abstract Flat List (GToken's)
// from Concrete Flat List (File-Format-Specific tokens).
// .
func (p *Contentity) st1c_MakeAFLfromCFL() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1c"
	var e error
	var GTs []*gtoken.GToken
	var common XU.CommonCPR
	var NSer NodeStringser

	fmt.Fprintln(p.GTknsWriter, "=== Input file:", p.AbsFP())

	switch p.MarkupType() {
	case SU.MU_type_MKDN:
		var pCPR_M *PU.ParserResults_mkdn
		if nil == p.ParserResults {
			p.L(LError, "ParserResults are nil")
		}
		pCPR_M = p.ParserResults.(*PU.ParserResults_mkdn)
		common = pCPR_M.CommonCPR
		NSer = pCPR_M
		// Do this
		pCPR_M.Writer = io.Discard
		// instead of this
		// pCPR_M.Writer = p.GTknsWriter
		GTs, e = gtoken.DoGTokens_mkdn(pCPR_M)
		if e != nil {
			p.WrapError("mkdn.gtokens", e)
		}
		// Use this
		p.GTokens = GTs
		/* instead of this
		// Compress out nil GTokens
		p.GTokens = make([]*gtoken.GToken, 0)
		for _, GT := range GTs {
			if GT != nil {
				p.GTokens = append(p.GTokens, GT)
			}
		}
		*/
	case SU.MU_type_HTML:
		var pCPR_H *PU.ParserResults_html
		pCPR_H = p.ParserResults.(*PU.ParserResults_html)
		common = pCPR_H.CommonCPR
		NSer = pCPR_H
		// Do this
		pCPR_H.Writer = io.Discard
		// instead of this
		// pCPR_H.Writer = p.GTknsWriter
		GTs, e = gtoken.DoGTokens_html(pCPR_H)
		if e != nil {
			p.WrapError("html.gtokens", e)
		}
		p.GTokens = GTs
	case SU.MU_type_XML:
		var pCPR_X *XU.ParserResults_xml
		pCPR_X = p.ParserResults.(*XU.ParserResults_xml)
		common = pCPR_X.CommonCPR
		NSer = pCPR_X
		// Do this
		pCPR_X.Writer = io.Discard
		// instead of this
		// pCPR_X.Writer = p.GTknsWriter
		GTs, e = gtoken.DoGTokens_xml(pCPR_X)
		if e != nil {
			p.WrapError("GToken-ization", e)
		}
		p.TallyTags()
		// fmt.Printf("==> Tags: %v \n", pGF.TagTally)
		// fmt.Printf("==> Atts: %v \n", pGF.AttTally)
		p.GTokens = GTs
	}
	ndL := len(common.NodeDepths)
	fpL := len(common.FilePosns)
	gtL := len(p.GTokens)
	nsL := NSer.NodeCount()
	fmt.Fprintln(p.GTknsWriter, "=== Output:")
	fmt.Fprintf(p.GTknsWriter, "=== array lengths: " +
		"cmn.NodeDepths<%d> cmn.FilePosns<%d> GTokens<%d> " +
		"NSer.NodeCount<%d> \n", ndL, fpL, gtL, nsL)
	/* fmt.Printf("st1-read: array lengths: " +
		"cmn.NodeDepths<%d> cmn.FilePosns<%d> GTokens<%d> " +
		"NSer.NodeCount<%d> \n", ndL, fpL, gtL, nsL) */
	count := (ndL + fpL + gtL + nsL + 2) / 4
	// For every GToken, we should print:
	//  - original node's original text
	//  - ditto, but as rendered by node's class
	//  - original node's NodeEcho(int)
	//  - original node's NodeInfo(int) and/or NodeDebug(int)
	//  - the GToken
	for i := 0; i < count; i++ {
	        if p.GTokens[i] == nil {
		   L.L.Warning("NIL at GTokens[%d]", i)
		   continue
		   }
		tkn := *(p.GTokens[i])
		fmt.Fprintf(p.GTknsWriter, "[%d]\n", i)
		fmt.Fprintf(p.GTknsWriter, "echo: %s \n", NSer.NodeEcho(i))
		fmt.Fprintf(p.GTknsWriter, "info: %s \n", NSer.NodeInfo(i))
		fmt.Fprintf(p.GTknsWriter, "dbug: %s \n", NSer.NodeDebug(i))
		// Dump the GToken
		// fmt.Fprintln(p.GTknsWriter, (*pGtkn).String())
		// fmt.Fprintf(p.GTknsWriter, "<%+v>\n", *pGtkn)
		fmt.Fprintf(p.GTknsWriter, "<%s>\n", tkn.Echo())
	}
	// fmt.Printf("st1c_MakeAFLfromCFL: nGTokens: %d %d \n", len(p.GTokens), len(GTs))
	return p
}

// st1d_PostMeta_notmkdn is Step 1d (XML,HTML): XML per format; HTML <head>
func (p *Contentity) st1d_PostMeta_notmkdn() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "1d"
	switch p.MarkupType() {
	case SU.MU_type_MKDN:
		// Markdown YAML metadata was processed in step st1a
		return p
	case SU.MU_type_HTML:
		/* code to use!
		var pPR *PU.ParserResults_html
		pPR = p.CPR.(*PU.ParserResults_html)
		z := pPR. */
		// Inside <head>: <meta> <title> <base> <link> <style>
		// See also: https://gist.github.com/lancejpollard/1978404
		return p
	case SU.MU_type_XML:
		// [Lw]DITA stuff, ?DublinCore
		p.L(LWarning, "cty.st1.TODO: SetMTypePerDoctypeFields:")
		p.L(LWarning, "     \\ "+p.PathAnalysis.String())

	}
	return p
}
