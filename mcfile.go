package mcfile

// link: tag, att, text raw, ref abspath

import (
	"fmt"
	"io"
	"os"

	CFU "github.com/fbaube/cliflagutils"
	"github.com/fbaube/db"
	FU "github.com/fbaube/fileutils"

	// XM "github.com/fbaube/xmlmodels"
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"

	// "github.com/pkg/errors"
	_ "github.com/sanity-io/litter"
)

// MCFileProcessor is here for reference, altho the name is not really used.
// It signals an error by setting the embedded field `error`; see
// `(*CheckedPath)` funcs `Error() string` and `GetError() error`.
type MCFileProcessor func(*MCFile) *MCFile

type ContentityStage func(*Contentity) *Contentity

// The data structure food chain (TODO update this):

// - CLI arg
// - RelFilePath
// - AbsFilePath
// - BasicPath
// - CheckedContent
// - MCFile (TODO ContentityNord)
// - - TypeXml OR TypeHtml OR TypeMkdn
// - CCT (Concrete Content Tree)
// - ACT (Abstract Content Tree) (GTree embedded in MCFile)
// - ForesTree (with x-refs etc) (TBD)
// - Grove (?)

// MCFile contains an input file (i.e. its raw UTF-8 content),
// its deduced (i.e. guessed) properties (MIME, XML, DITA),
// its markup element tree, its output paths & files, and its
// references to external files, URLs, and other content files.
// NOTE We always create an MCFile for every input file, so
// this is a logical place to store a GTokenization and a GTree.
type MCFile struct {
	CFU.GCtx
	// logIdx int
	// logStg string
	// db.Times
	// FU.PathProps
	db.ContentRecord // embeds FU.AnalysisRecord

	// Data stuctures and conversions ("FFS" = file-format-specific):
	// 1) CCT = Concrete Content Tree = FFS-nodes [not available for XML]
	// 2) CFL = Concrete Flat List = FFS-nodes,tokens [incl. depth info]
	// 3) AFL = Abstract Flat List = GTokens > GTags
	// 4) ACT = Abstract Content Tree = GTags [assembled using depth info]

	// CCT is the root of the CTT (tree); note that the type of the root
	// node is normally the same as the type of the other nodes in the tree.

	// ParserResults is parseutils.ParserResults_ffs and includes CCT and CFL.
	ParserResults interface{}
	// AFL
	GTokens []*gtoken.GToken
	// AFL
	GTags []*gtree.GTag
	// ACT
	*gtree.GTree // maybe not need GRootTag or RootOfASTp

	GTokensOutput, GTreeOutput io.Writer

	// []GTokens are Generalized tokens, converted from the precursor
	// tokens (or, nodes) emitted by format-specific parsers.
	// - XML : https://golang.org/pkg/encoding/xml/#Token
	// - HTML: https://godoc.org/golang.org/x/net/html#Token
	// - MKDN: https://godoc.org/github.com/yuin/goldmark/ast#Node
	// Each `gparse.GToken` wraps its precursor token.

	// A `gtree.GTag`   wraps its corresponding
	//  `gparse.GToken` wraps its corresponding
	//     FFS-token

	TagTally StringTally
	AttTally StringTally

	// The article about go types for functions
	// MAKE BLOCK LIST
	// SORT OUT RESOLUTION OF GLinks
	// GATHER ToC ELEMENTS
	// Sep. XML types into a map of callback functions ?

	// FU.OutputFiles // NOTE Does this belong here ? Not sure.

	*GLinks

	// GEnts is "ENTITY"" directives (both with "%" and without).
	GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
	DElms map[string]*gtree.GTag
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

/*
func (p *MCFile) LogIt(s string) {
	if p.OwnLog != nil {
		p.OwnLog.Printf(s)
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
	fmt.Fprintf(os.Stdout, "--> "+s+"\n")
}
*/

// NewMCFile // also sets `MCFile.MType[..]`.
func NewMCFile(pCR *db.ContentRecord) *MCFile {
	p := new(MCFile)
	p.ContentRecord = *pCR
	// println("ZORK3:", pCR.Error())
	// println("ZORK2:", pCR.GetError().Error())
	if pCR.GetError() != nil {
		pCR.SetError(fmt.Errorf("newMCF<%s> failed: %w",
			pCR.AbsFilePath, pCR.GetError()))
		return p
	}
	p.GLinks = new(GLinks)
	println("D=> (C:NewMCF)", p.String()) // p.MType, p.AbsFP())
	return p
}

// Errorbarf just complains and dies. At the top level
// (i.e. in main()), we don't wrap errors and return them.
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
		p.PathProps.AbsFP(), s, e.Error())
	// os.Exit(1)
	println("==> DUMP OF FAILING MCFILE:")
	println(p.String())
	return true
}

func (p *MCFile) Lengths() string {
	return fmt.Sprintf("len.raw.file<%d> len.meta.hdr.props<%d> len.text.body.content<%d>",
		len(p.Raw), len(p.MetaRaw()), len(p.TextRaw()))
}

// String is developer output. Hafta dump:
// FU.InputFile, FU.OutputFiles, GTree,
// GRefs, *XmlFileMeta, *XmlItems, *DitaInfo
func (p MCFile) String() string {
	// !! // var BF BigFields = p.PushBigFields()

	var sGTree string
	if p.GTree != nil {
		sGTree = p.GTree.String()
	}
	// s := fmt.Sprintf("[len:%d]", p.Size())
	s := fmt.Sprintf("||%s||GTree|%s||OutbKeyLinks|%+v|KeyLinkTgts|%+v|OutbUriLinks|%+v|UriLinkTgts|%+v||",
		FU.Tildotted(p.PathProps.AbsFP()) /* p.OutputFiles.String(), */, sGTree,
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
	s += fmt.Sprintf("DitaInfo|ML:%s|Cntp:%s|", p.DitaMarkupLg, p.DitaContype)
	// }

	// !! // p.PopBigFields(BF)
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

func (pMCF *MCFile) GatherGLinksInto(pGL *GLinks) {
	var aGL *GLink

	if len(pMCF.OutgoingKeys) > 0 || len(pMCF.IncomableKeys) > 0 {

		// Dir := pCF.CheckedPath.AbsFilePathParts.DirPath
		// println("### KEY GLinks:", Dir)
		if len(pMCF.OutgoingKeys) > 0 {
			pGL.OutgoingKeys = append(pGL.OutgoingKeys, pMCF.OutgoingKeys...)
			for _, aGL = range pMCF.OutgoingKeys {
				fmt.Printf("Outgng: %s@%s: %s \n", aGL.Tag, aGL.Att, aGL.Link_raw)
			}
		}
		if len(pMCF.IncomableKeys) > 0 {
			pGL.IncomableKeys = append(pGL.IncomableKeys, pMCF.IncomableKeys...)
			for _, aGL = range pMCF.IncomableKeys {
				fmt.Printf("Incmbl: %s@%s: %s \n", aGL.Tag, aGL.Att, aGL.Link_raw)
			}
		}
	}
	if len(pMCF.OutgoingURIs) > 0 || len(pMCF.IncomableURIs) > 0 {
		// !! println("### URI GLinks:", Dir)
		if len(pMCF.OutgoingURIs) > 0 {
			pGL.OutgoingURIs = append(pGL.OutgoingURIs, pMCF.OutgoingURIs...)
			for _, aGL = range pMCF.OutgoingURIs {
				fmt.Printf("###  # Outgng: %s@%s: %s \n", aGL.Tag, aGL.Att, aGL.Link_raw)
			}
		}
		if len(pMCF.IncomableURIs) > 0 {
			pGL.IncomableURIs = append(pGL.IncomableURIs, pMCF.IncomableURIs...)
			for _, aGL = range pMCF.IncomableURIs {
				fmt.Printf("Incmbl: %s@%s: %s \n", aGL.Tag, aGL.Att, aGL.Link_raw)
			}
		}
	}
}

// === Implement interface Errable

func (p *MCFile) HasError() bool {
	return p.ContentRecord.HasError() || p.PathProps.HasError()
}

// GetError is necessary cos "Error()"" dusnt tell you whether "error"
// is "nil", which is the indication of no error. Therefore we need
// this function, which can actually return the telltale "nil".
func (p *MCFile) GetError() error {
	if p.PathProps.HasError() {
		return p.PathProps.GetError()
	}
	return p.ContentRecord.GetError()
}

// Error satisfies interface "error", but the
// weird thing is that "error" can be nil.
func (p *MCFile) Error() string {
	if p.PathProps.HasError() {
		return p.PathProps.Error()
	}
	return p.ContentRecord.Error()
}

func (p *MCFile) SetError(e error) {
	p.ContentRecord.SetError(e)
}
