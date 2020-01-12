package mcfile

import (
	"fmt"

	// "github.com/dimchansky/utfbom"
	// BF "github.com/fbaube/blackfriday"
	// BF "github.com/yuin/goldmark"
	AST "github.com/yuin/goldmark/ast"
)

/*

func NewGElmTreeFromBFtree(p *BF.Node) *GRootElm {
	println("NewGElmTreeFromBFtree ENTRY")
	var pp *GElm
	pp = NewGElmFromBFnode(p)
	var Kids []*BF.Node
	Kids = KidsAsSlice(p)
	for i, c := range Kids {
		fmt.Printf("[%d]INN ", i)
		gtoken.DumpBFnode(c, 0)
		fmt.Printf("[%d]OUT ", i)
		cc := NewGElmFromBFnode(c)
		if pp == nil {
			panic("nil Parent")
		}
		if cc == nil {
			panic("nil Kid")
		}
		pp.AddKid(cc)
	}
	var pRoot *GRootElm
	pRoot = (*GRootElm)(pp)
	return pRoot
}

// NewGElmFromBFnode basically just assigns to this field:
// - gtoken.GToken
// which comprises:
// - GElmTokType
// - GName
// - GAttList
func NewGElmFromBFnode(p *BF.Node) *GElm {
	var NT BF.NodeType
	NT = p.Type
	var pp *GElm
	pp = new(GElm)
	pp.GElmTokType = "SE"
	sKids := ListKids(p)
	println("New GElm ::", NT.String(), "::", sKids)

	switch NT.String() {
	case "Document":
		println("START OF DOCUMENT")
		pp.GName = *gtoken.NewGName("", "markdown")
		return pp
	case "List", "Item":
		var lst = p.ListData
		println(gtoken.DumpList(lst))
		return pp
	case "Heading":
		var hdg = p.HeadingData
		println(gtoken.DumpHdg(hdg))
		pp.GName = *gtoken.NewGName("", fmt.Sprintf("H%d", hdg.Level))
		return pp
	case "Link", "Image":
		var lnk = p.LinkData
		println(gtoken.DumpLink(lnk))
		return pp
	case "CodeBlock":
		var cbd = p.CodeBlockData
		println(gtoken.DumpCdBlk(cbd))
		return pp
	case "Table", "TableCell", "TableHead", "TableBody", "TableRow":
		var tcd = p.TableCellData
		println(gtoken.DumpTableCell(tcd))
		return pp

	case "BlockQuote":
		return pp
	case "Paragraph":
		return pp
	case "HorizontalRule":
		return pp
	case "Emph":
		return pp
	case "Strong":
		return pp
	case "Del":
		return pp
	case "Text":
		return pp
	case "HTMLBlock":
		return pp
	case "Softbreak":
		return pp
	case "Hardbreak":
		return pp
	case "Code":
		return pp
	case "HTMLSpan":
		return pp
	}
	return pp
}

*/

// =================
// ==== BF node ====
// == (BF library) =
// =================

func KidsAsSlice(p AST.Node) []AST.Node {
	var pp []AST.Node
	c := p.FirstChild()
	for c != nil {
		pp = append(pp, c)
		c = c.NextSibling()
	}
	return pp
}

func ListKids(p AST.Node) string {
	var pp []AST.Node
	pp = KidsAsSlice(p)
	if len(pp) == 0 {
		return "<0-kids>"
	}
	var s string
	for i, c := range pp {
		s += fmt.Sprintf("[%d:%s]", i, c.Type)
	}
	return s
}

/*

// type myBFnodeVisitor func(node *Node, entering bool) WalkStatus
func myBFnodeVisitor(N *BF.Node, entering bool) BF.WalkStatus {
	if !entering {
		return BF.GoToNext
	}
	// fmt.Printf("%s \n" /* SU.GetIndent(lvl), * /, DumpBFnode(N, 0))
	return BF.GoToNext
}

*/

// ==================
// ==== AST node ====
// = (found object) =
// ==================

/*
Type     string
Literal  string      `json:",omitempty"`
Attr     interface{} `json:"-"`
Children []*ASTNode  `json:",omitempty"`
*/

/*
func (p *ASTNode) DumpASTnode() string {

for child := range p.Children { // p.FirstChild; child != nil; child = child.Next {
	a.Children = append(a.Children, NewASTNode(child))
}
}

// type myASTnodeVisitor func(node *Node, entering bool) WalkStatus
func myASTnodeVisitor(N *ASTnode, entering bool) BF.WalkStatus {
	if !entering {
		return BF.GoToNext
	}
	fmt.Printf("%s \n" /* SU.GetIndent(lvl), * /, DumpASTnode(N))
	return BF.GoToNext
}
*/

// ================
// ===== GElm =====
// ================

func DumpGElm(p AST.Node) string {
	var s string
	/*
		s = fmt.Sprintf("BFnode<%s>: ", gparse.MDnodeType[p.Type])
		s += "<|" + SU.NormalizeWhitespace(string(p.Literal)) + "|> "
		s += gparse.DumpHdg(p.HeadingData)
		s += gparse.DumpList(p.ListData)
		s += gparse.DumpCdBlk(p.CodeBlockData)
		s += gparse.DumpLink(p.LinkData)
		s += gparse.DumpTableCell(p.TableCellData)
	*/
	return s
}

// type myGElmVisitor func(node *Node, entering bool) WalkStatus
func myGElmVisitor(N AST.Node, entering bool) AST.WalkStatus {
	if !entering {
		return AST.WalkContinue // GoToNext
	}
	fmt.Printf("%s \n" /* SU.GetIndent(lvl), */, DumpGElm(N))
	return AST.WalkContinue // GoToNext
}

/*

// =======================
// ==== BF node stuff ====
// =======================

// Level        int    // This holds the heading level number
// HeadingID    string // This might hold heading ID, if present
// IsTitleblock bool   // Specifies whether it's a title block
func DumpHdg(h BF.HeadingData) string {
	var ttl string
	if h.IsTitleblock {
		ttl = "IsTtl:"
	}
	if h.Level == 0 && h.HeadingID == "" {
		return ""
	}
	return fmt.Sprintf("<H%d::%s%s> ", h.Level, ttl, h.HeadingID)
}

// ListFlags   ListType
// Tight       bool   // Skip <p>s around list item data if true
// BulletChar  byte   // '*', '+' or '-' in bullet lists
// Delimiter   byte   // '.' or ')' after the number in ordered lists
// RefLink   []byte   // If not nil, turns this list item into a footnote item and triggers different rendering
// IsFootnotesList bool   // This is a list of footnotes
func DumpList(L BF.ListData) string {
	var tight, ftnts string
	if L.Tight {
		tight = "IsTight:"
	}
	if L.IsFootnotesList {
		ftnts = "IsFtnts:"
	}
	if L.ListFlags == 0 && L.BulletChar == 0 &&
		L.Delimiter == 0 && len(L.RefLink) == 0 {
		return ""
	}
	return fmt.Sprintf("<List:%s%sBult:%c:Delim:%c:RefLink:%s> ",
		tight, ftnts, L.BulletChar, L.Delimiter, L.RefLink)
}

// IsFenced   bool  // Fenced code block, or else an indented one
// Info     []byte  // This holds the info string
// FenceChar  byte
// FenceLength int
// FenceOffset int
func DumpCdBlk(cb BF.CodeBlockData) string {
	var fenced string
	if cb.IsFenced {
		fenced = "IsFenced"
	}
	if len(cb.Info) == 0 && cb.FenceChar == 0 &&
		cb.FenceLength == 0 && cb.FenceOffset == 0 {
		return ""
	}
	return fmt.Sprintf("<CdBlk:%s:ch:%c:len:%d:ofs:%d:Info:%s> ",
		fenced, cb.FenceChar, cb.FenceLength, cb.FenceOffset, string(cb.Info))
}

// Destination []byte // Destination is what goes into a href
// Title       []byte // The tooltip thing that goes in a title attribute
// NoteID      int    // The S/N of a footnote, or 0 if not a footnote
// Footnote    *Node  // If footnote, a direct link to the FN Node, else nil.
func DumpLink(L BF.LinkData) string {
	var isFN bool
	isFN = (L.NoteID != 0) && (L.Footnote == nil)
	if len(L.Destination) == 0 && len(L.Title) == 0 &&
		L.NoteID == 0 && L.Footnote == nil {
		return ""
	}
	if !isFN {
		return fmt.Sprintf("<Link:Ttl:%s:Dest:%s> ",
			string(L.Title), string(L.Destination))
	}
	return fmt.Sprintf("<FN-link:#%d:Ttl:%s:Dest:%s> ",
		L.NoteID, string(L.Title), string(L.Destination))
}

// IsHeader  bool       // This tells if it's under the header row
// Align CellAlignFlags // This holds the value for align attribute
func DumpTableCell(tc BF.TableCellData) string {
	if tc.Align == 0 && !tc.IsHeader {
		return ""
	}
	if tc.IsHeader {
		return "<TblCell:IsHdr> "
	}
	return "<TblCell:notHdr> "
}
*/
