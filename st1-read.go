package mcfile

import (
	"fmt"
	// "encoding/xml"
	S "strings"
	// "log"
	"errors"
	"golang.org/x/net/html"
	"github.com/fbaube/gparse"
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

// GetParseTokenization_Xml v GetParseTree_nonXml
// GetNodelistFromParseTree_NonXml
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
	return p.
		st1a_PreMeta().
		st1b_AST_notXml().
		st1c_Tokenize().
		st1d_PostMeta_notMkdn(). // XML per format; HTML <head>
		st1e_GTokenize()
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

// st1b_AST_notXml is Step 1b (MKDN,HTML): TBS...
func (p *MCFile) st1b_AST_notXml() *MCFile {
	if p.GetError() != nil {
		return p
	}
	var e error
	switch p.FileType() {
	case "XML":
		return p
	case "MKDN":
		var pTheMTokzn *gparse.MkdnTokenization
		pTheMTokzn, e = gparse.MkdnTokenizeBuffer(p.CheckedContent.Raw)
		p.RootOfASTp = pTheMTokzn.TreeRootNode
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			p.Blare(p.OLP + e.Error())
			p.SetError(e)
			return p
		}
		p.TypeSpecificTokenizationP = pTheMTokzn
		// p.TallyTags()
	case "HTML":
		// Parse returns the parse tree for the HTML from the given Reader.
		// The input is assumed to be UTF-8 encoded.
		var pTheHTokzn *gparse.HtmlTokenization
		pTheHTokzn.TreeRootNodeP, e = html.Parse(S.NewReader(p.CheckedContent.Raw))
		p.RootOfASTp = pTheHTokzn.TreeRootNodeP // gparse.HtmlAST{Node:*pHN}
		if e != nil {
			e = errors.New("st[1b] " + e.Error())
			p.Blare(p.OLP + e.Error())
			p.SetError(e)
			return p
		}
		p.TypeSpecificTokenizationP = pTheHTokzn
		// p.TallyTags()
		return p
	}
	return p
}

// st1c_Tokenize is Step 1c: processing continues as far as it can whilst
// still being done at the level of individual tokens, so that (for example)
// every XML token is up-converted into a `GToken`.
//
// NOTE that if we have no separate tokenization step, and we go straight to
// an AST (for example, Markdown and HTML), then h/ere we pause to make a list
// of tokens, and up-convert them to `GToken`s, but we do not otherwise
// process the AST at all yet.
//
func (p *MCFile) st1c_Tokenize() *MCFile {
	if p.GetError() != nil {
		return p
	}
	if len(p.CheckedContent.Raw) == 0 {
		p.Whine(p.OLP + "st[1c] " + "Zero-length content")
		return p
	}
	var e error
	var errmsg string

	switch p.FileType() {
	case "MKDN":
		var pMT *gparse.MkdnTokenization
		pMT = (p.TypeSpecificTokenizationP).(*gparse.MkdnTokenization)
		// pMT = (&p.TreeRootNode).(*gparse.MkdnTokenization)
		pMT.MkdnNodeListFromAST()
		return p
	case "HTML":
		var pHT *gparse.HtmlTokenization
		pHT = (p.TypeSpecificTokenizationP).(*gparse.HtmlTokenization)
		// pMT = (&p.TreeRootNode).(*gparse.MkdnTokenization)
		pHT.HtmlNodeListFromAST()
		return p
	case "XML":
		var TheXTokzn gparse.XmlTokenization
		TheXTokzn = *new(gparse.XmlTokenization)
		// TheXTokzn.Tokens = make([]xml.Token, 0)
		TheXTokzn.Tokens, e = gparse.XmlTokenizeBuffer(p.CheckedContent.Raw)
		if e != nil {
			e = fmt.Errorf("XML tokenization failed: %w", e)
		}
		// TypeSpecificTokens []gparse.MarkupStringer
		p.TypeSpecificTokenizationP = TheXTokzn
	}
	if e != nil {
		errmsg = "st[1c] " + e.Error()
		p.Blare(p.OLP + errmsg)
		p.SetError(e)
		return p
	}
	p.TallyTags()
	// fmt.Printf("==> Tags: %v \n", pGF.TagTally)
	// fmt.Printf("==> Atts: %v \n", pGF.AttTally)
	return p
}

// st1d_PostMeta_notMkdn is Step 1d (XML,HTML): XML per format; HTML <head>
func (p *MCFile) st1d_PostMeta_notMkdn() *MCFile {
	switch p.FileType() {
	case "MKDN":
		return p
	case "HTML":
		// Inside <head>: <meta> <title> <base> <link> <style>
		// See also: https://gist.github.com/lancejpollard/1978404
		return p
	}
	return p
}

// st1e_GTokenize is Step 1e: TBS...
func (p *MCFile) st1e_GTokenize() *MCFile {
	if p.GetError() != nil {
		return p
	}
	var e error
	var errmsg string

	switch p.FileType() {
	case "MKDN":
		// TBS: CONVERT THE FLAT LIST
		return p
	case "HTML":
		// TBS: CONVERT THE FLAT LIST
		return p
	case "XML":
		// XML-Tokenize
		var XT gparse.XmlTokenization
		XT = p.TypeSpecificTokenizationP.(gparse.XmlTokenization)
		p.GTokens, e = gparse.GTokznFromXmlTokens(XT.Tokens)
		if e != nil {
			e = fmt.Errorf("GToken-ization failed: %w", e)
		}
		if e != nil {
			errmsg = "st[1b] " + e.Error()
			p.Blare(p.OLP + errmsg)
			p.SetError(e)
			return p
		}
		p.TallyTags()
		// fmt.Printf("==> Tags: %v \n", pGF.TagTally)
		// fmt.Printf("==> Atts: %v \n", pGF.AttTally)
	}
	return p
}
