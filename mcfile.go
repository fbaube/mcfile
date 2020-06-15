package mcfile

// link: tag, att, text raw, ref abspath

import (
	"fmt"
	"os"

	"github.com/fbaube/db"
	FU "github.com/fbaube/fileutils"
	MU "github.com/fbaube/miscutils"
	XM "github.com/fbaube/xmlmodels"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	// "github.com/pkg/errors"
	_ "github.com/sanity-io/litter"
)

// MCFileProcessor is here for reference, altho the name is not really used.
// It signals an error by setting the embedded field `error`; see
// `(*CheckedPath)` funcs `Error() string` and `GetError() error`.
type MCFileProcessor = func(*MCFile) *MCFile

// CCTnode is a node in a CCT, which is a parse tree of HTML or Markdown-
type CCTnode interface{}

// The data structure food chain:
// - CLI arg
// - RelFilePath
// - AbsFilePath
// - BasicPath
// - CheckedContent
// - MCFile
// - - TypeXml OR TypeHtml OR TypeMkdn
// - CCT (Concrete Content Tree)
// - ACT (Abstract Content Tree) (GTree embedded in MCFile)
// - ForesTree (with x-refs etc) (TBD)
// - Grove (?)

// NOTE We always create an MCFile for every input file, so
// it is a logical place to store a GTokenization and a GTree.

// MCFile contains an input file (i.e. its raw UTF-8 content),
// its deduced (i.e. guessed) properties (MIME, XML, DITA),
// its markup element tree, its output paths & files, and its
// references to external files, URLs, and other content files.
// NOTE We always create an MCFile for every input file, so
// this is a logical place to store a GTokenization and a GTree.
type MCFile struct {
	MU.GCtx
	db.Times
	Idx         int
	Idx_Inbatch int

	// These THREE fields contain the file contents.
	// CheckedContent.Raw == Meta_raw + Text_raw
	FU.CheckedContent // Field `Raw` has the raw content of the entire file
	// Header (i.e. Metadata) -- see comments below, after this struct def
	Meta_raw   string
	MetaFormat string // "yaml", "dita", "html", etc.
	MetaProps  map[string]string
	Text_raw   string

	// File-format-specific ptr to additional data - XML, HTML, MKDN
	// FFSdataP interface{}

	// Data stuctures and conversions ("FFS" = file-format-specific):
	// 1) CCT = Concrete Content Tree = FFS-nodes [not available for XML]
	// 2) CFL = Concrete Flat List = FFS-nodes,tokens [incl. depth info]
	// 3) AFL = Abstract Flat List = GTokens > GTags
	// 4) ACT = Abstract Content Tree = GTags [assembled using depth info]

	// CCT is the root of the CTT (tree); note that the type of the root
	// node is normally the same as the type of the other nodes in the tree.

	// CPR is ConcreteParseResults (i.e. parseutils.ParseResults_ffs)
	// and it includes CCT and CFL.
	CPR interface{}
	// AFL
	GTokens []*gtoken.GToken
	// AFL
	GTags []*gtree.GTag
	// ACT
	*gtree.GTree // maybe not need GRootTag or RootOfASTp

	// []GTokens are Generalized tokens, converted from the precursor
	// tokens (or, nodes) emitted by format-specific parsers.
	// - XML : https://golang.org/pkg/encoding/xml/#Token
	// - HTML: https://godoc.org/golang.org/x/net/html#Token
	// - MKDN: https://godoc.org/github.com/yuin/goldmark/ast#Node
	// Each `gparse.GToken` wraps its precursor token.

	// A  `gtree.GTag` wraps its corresponding
	// `gparse.GToken` wraps its corresponding
	// FFS-token

	TagTally StringTally
	AttTally StringTally

	// The article about go types for functions
	// MAKE BLOCK LIST
	// SORT OUT RESOLUTION OF GLinks
	// GATHER ToC ELEMENTS
	// Sep. XML types into a map of callback functions ?

	// FU.OutputFiles // NOTE Does this belong here ? Not sure.

	*GLinks

	XM.XmlInfo
	TypeXml
	// DitaInfo is two enums (so far): Markup language & Content type.
	// ML: "1.2", "1.3", "XDITA", "HDITA", "MDATA".
	// CT: "Map", "Bookmap", "Topic", "Task", "Concept", "Reference",
  // "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"
	XM.DitaInfo
}

// The terms "header" and "metadata" are used interchangeably.
// In default usage, this is metadata stored *in* the file, e.g.
// YAML in LwDITA Markdown-XP, or `head/meta` tags in [X]HTML.
// (Obv we want to store as much metadata as possible in-file rather
// than externally, and we will need to map select terms btwn formats.)
// We store it at the same level as the file's "content", aka "Text".
// Metadata stored this way is easier to manage in a format-independent
// manner, and it is easier to add to it and modify it at runtime, and
// (TODO) when it is stored as JSON K/V pairs, it can be accessed from
// the command line using Sqlite (and other nifty) tools.

func (p *MCFile) LogIt(s string) {
	if p.Log != nil {
		p.Log.Printf(s)
	} else {
		p.SsnLog.Printf(s)
	}
}

// Blare is used for errors that will stop processing.
func (p *MCFile) Blare(s string) {
	p.LogIt(s)
	fmt.Fprintf(os.Stderr, s)
	// println("Bogus SU ref \n", SU.GetIndent(2))
}

// Whine is used for non-fatal errors, i.e. strong warnings.
func (p *MCFile) Whine(s string) {
	p.LogIt(s)
	fmt.Fprintf(os.Stdout, "--> " + s)
}


// NewMCFile // also sets `MCFile.MType[..]`.
func NewMCFile(pCC *FU.CheckedContent) *MCFile {
	pMF := new(MCFile)
	pMF.CheckedContent = *pCC
	if pCC.GetError() != nil {
		pCC.SetError(fmt.Errorf("NewMCFile <%s>: %w", pCC.AbsFilePath, pCC.GetError()))
		return pMF
	}
	pMF.GLinks = new(GLinks)
	println("NewMCFile:", pCC.MType, pCC.AbsFilePath)
	return pMF
}


// NewMCFileFromPath checks that the path is actually a file,
// and also sets `MCFile.MType[..]`.
func NewMCFileFromPath(path string) *MCFile {
	// FIRST we work with a new CheckedPath
	pBP := FU.NewPathInfo(path)
	if pBP.GetError() != nil || !pBP.IsOkayFile() {
		pBP.SetError(fmt.Errorf("NewMCFileFromPath.BP <%s>: %w", path, pBP.GetError()))
		return nil
	}
	/*
		println("====")
		litter.Dump(*pCP)
		println("====")
	*/
	pCC := FU.NewCheckedContent(pBP) // new(FU.CheckedContent)
	/*
	pCC.BasicPath    =  pBP
	pCC.BasicContent = *pBP.FetchContent()
	pCC.BasicAnalysis.AnalyzeFileContent(pCC.Raw)
	*/
	println("--> MType:", pCC.MType) // Mstring())
	/*
		println("====")
		litter.Dump(*pCP)
		println("====")
	*/
	// Create the MCFile
	pMF := new(MCFile)
	pMF.CheckedContent = *pCC
	if pCC.GetError() != nil {
		pCC.SetError(fmt.Errorf("NewMCFileFromPath.CC <%s>: %w", path, pCC.GetError()))
		return pMF
	}
	return pMF
}

/*

// TheXml is a convenience function.
func (p *MCFile) TheXml() *TypeXml {
	switch ptr := p.FFSdataP.(type) {
	case *TypeXml:
		return ptr // (p.FFSdataP).(*TypeXml)
	case *TypeHtml:
		return &(ptr.TypeXml) // (p.FFSdataP).(*TypeXml)
	}
	// return (p.FFSdataP).(*TypeXml)
	panic("mcfile.TheXml")
}

// TheMkdn is a convenience function.
func (p *MCFile) TheMkdn() *TypeMkdn {
	return (p.FFSdataP).(*TypeMkdn)
}

// TheHtml is a convenience function.
func (p *MCFile) TheHtml() *TypeHtml {
	return (p.FFSdataP).(*TypeHtml)
}

*/

// At the top level (i.e. in main()), we don't wrap errors
// and return them. We just complain and die. Simple!
func (p *MCFile) Errorbarf(e error, s string) bool {
	if e == nil {
		return false
	}
	if e.Error() == "" {
		return false
	}
	p.SetError(e)
	// elog.Printf("%s failed: %s \n", myAppName, e.Error())
	fmt.Fprintf(os.Stderr, "%s failed: %s \n\t error was: %s \n",
		p.PathInfo.AbsFP(), s, e.Error())
	// os.Exit(1)
	println("==> DUMP OF FAILING MCFILE:")
	println(p.String())
	return true
}

func (p *MCFile) Lengths() string {
	return fmt.Sprintf("len.raw.file<%d> len.meta.hdr.props<%d> len.text.body.content<%d>",
		len(p.Raw), len(p.Meta_raw), len(p.Text_raw))
}

// String is developer output. Hafta dump:
// FU.InputFile, FU.OutputFiles, GTree,
// GRefs, *XmlFileMeta, *XmlItems, *DitaInfo
func (p MCFile) String() string {
	var BF BigFields = p.PushBigFields()

	// s := fmt.Sprintf("[len:%d]", p.Size())
	s := fmt.Sprintf("(DD:GFILE)||%s||OtFiles|ss||GTree|%s||OutbKeyLinks|%+v|KeyLinkTgts|%+v|OutbUriLinks|%+v|UriLinkTgts|%+v||",
		p.PathInfo.AbsFP(), /* p.OutputFiles.String(), */ p.GTree.String(),
		p.OutgoingKeys, p.IncomableKeys, p.OutgoingURIs, p.IncomableURIs)
	/*
			if p.XmlFileMeta != nil {
				s += fmt.Sprintf("XmlFileMeta|%s||", p.XmlFileMeta.String())
			}
		* /
		if p.IDinfo != nil {
			s += fmt.Sprintf("xf.IDinfo|%s||", p.IDinfo.String())
		}
	*/
	/* ==
	if p.GEnts != nil {
		// FIXME s += fmt.Sprintf("GEnts|%s||", p.GEnts.String())
		 * 	}
		 * 	if p.DElms != nil {
		// FIXME s += fmt.Sprintf("DElms|%s||", p.DElms.String())
	}
	== */
	// if p.DitaInfo != nil {
		s += fmt.Sprintf("DitaInfo|%s||", p.DitaInfo.String())
	// }

	p.PopBigFields(BF)
	return s
}

// ConfigureOutputFiles might do nothing, depending on the dirSuffix
// (it can be "") and the GFile's InputFile.
func (p *MCFile) ConfigureOutputFiles(dirSuffix string) error {
	println("mcfile.CfgOutputFiles OMITTED")
	/*
	pOF, e := p.CheckedContent.NewOutputFiles(dirSuffix)
	if e != nil {
		return errors.Wrapf(e,
			"mcfile.ConfigureOutputFiles<%s>", p.BasicPath.AbsFilePathParts.String())
	}
	p.OutputFiles = *pOF
	*/
	return nil
}
