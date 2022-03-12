package mcfile

import (
	// "github.com/pkg/errors"
	"encoding/xml"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// TYPES ARE; external+http, file, pure-ID

var LwDitaAttsForGLinks = []string{
	"name",        // TARGET, in <a>, acts like @id
	"href",        // SOURCE, in many tags
	"id",          // TARGET, in all(?) tags
	"idref",       // SOURCE
	"idrefs",      // SOURCE
	"conref",      // SOURCE, xdita
	"data-conref", // SOURCE, hdita
	"keys",        // TARGET, xdita
	"data-keys",   // TARGET, hdita
	"keyref",      // SOURCE, xdita
	"data-keyref", // SOURCE, hdita
}

// LinkInfos is:
// @conref to reuse block-level content,
// @keyref to reuse phrase-level content.
// TODO Each type of link (i.e. elm/att where it occurs) has to be categorised.
// TODO Each format of link target has to be categorised.
// Cross ref : <xref> : <a href> : [link](/URI "title")
// Key def : <keydef> : <div data-class="keydef"> : <div data- class="keydef"> in HDITA syntax
// Map : <map> : <nav> : See Example of an MDITA map (20)
// Topic ref : <topicref> : <a href> inside a <li> : [link](/URI "title") inside a list item
// TODO Stuff to get:
// XDITA map
// - topicref @href (w @format)
// - task @id
// HDITA
// - article @id
// - span @data-keyref
// - p @data-conref
// MDITA
// - has YAML "id"
// - uses <p @data-conref>
// - uses <span @data-keyref>
// - uses MD [link_text](link_target.dita)
// - uses ![The remote](../images/remote-control-callouts.png "The remote")
// XDITA
// - topic @id
// - ph @keyref
// - image @href
// - p @id
// - video/source @value
// - section @id
// - p @conref
//
// In GFile:
// LinkInfos:
type LinkInfos struct {
	xmlIDs    []LinkInfo
	xmlIDrefs []LinkInfo
	Conrefs   []LinkInfo
	Keyrefs   []LinkInfo
	Datarefs  []LinkInfo
}

type LinkInfo struct {
	// link: tag, att, text raw, ref abspath
}

// XmlFileMeta is non-nil IFF it is an XML file (incl. DTD types).
// // *gparse.XmlFileMeta

// GatherXmlGLinks is:
// XmlItems is (DOCS) IDs & IDREFs, (DTDs) Elm defs (incl. Att defs) & Ent defs
// *xmlfile.XmlItems
// // *IDinfo
//
// func (p *MCFile) GatherXmlGLinks() *MCFile {
func (p *Contentity) GatherXmlGLinks() *Contentity {
	// pX := p.TheXml()
	// Establish the directory of the GFile.
	// ## // ## // var Dir FU.AbsFilePath = pGF.InputFile.FileFullName.DirPath
	// Iterate over all GTokens.
	for _, GT := range p.GTokens {
		if GT == nil {
			continue
		}
		// If it's not a Start Element, skip it
		if GT.TTType != "Elm" {
			continue
		}
		GN := GT.GName
		XN := xml.Name(GN)
		var theTag string = XN.Local
		// Iterate over all attributes
		for _, GA := range GT.GAtts {
			XA := xml.Attr(GA) // (*GA)
			if !SU.IsInSliceIgnoreCase(XA.Name.Local, LwDitaAttsForGLinks) {
				continue
			}
			pGL := new(GLink)
			/*
				type GLink struct {
					Att    string         // id, idref, href, xref, keyref, etc.
					Tag    string
					Raw    string         // as redd in during parsing
					RelFP  FU.RelFilePath // "" if target is in same file
					AbsFP  FU.AbsFilePath // "" if target is in same file
					Anchor string  // peeled off from Raw
					Resolved bool  // used to narrow in on difficult cases
					IsHref   bool  // resolution is simpler if is href
					Source *gtree.GTag    // the tag where the GLink is defined
					Target *gtree.GTag    // can be nil: the tag where the GLink is resolved
					to, quite possibly in another file, which we hope we also have availablein memory.
			*/
			pGL.Att = XA.Name.Local
			pGL.Tag = theTag
			pGL.Link_raw = XA.Value
			// Is it HTTP, FTP, etc. ?
			if i := S.Index(pGL.Link_raw, "://"); i > 0 {
				pGL.AddressMode = pGL.Link_raw[:i]
				pGL.Resolved = true
				if pGL.Att != "href" {
					panic("Non-@href http:://!")
				}
				// This is a hack !!
				pGL.AbsFP = "/"
			} else if S.Contains(pGL.Att, "key") {
				pGL.AddressMode = "key"
				if i := S.Index(pGL.Link_raw, "#"); i != -1 {
					pGL.FragID = pGL.Link_raw[i+1:]
					pGL.RelFP = pGL.Link_raw[:i]
				} else {
					pGL.RelFP = pGL.Link_raw
				}
				p.L(LDbg, "KEY:", pGL.RelFP, "#", pGL.FragID)
				// p.AbsFP = FU.RelFilePath(FP.Join(
				// 	pGF.InputFile.FileFullName.Echo(), p.RelFP.S())).AbsFP()
				s, _ := FP.Abs(FP.Join(p.PathProps.AbsFP.S(), pGL.RelFP))
				pGL.AbsFP = FU.AbsFilePath(s)
				p.L(LDbg, "2. AbsFP:", pGL.AbsFP)
			} else if S.HasPrefix(pGL.Att, "idref") {
				pGL.AddressMode = "idref"
				if i := S.Index(pGL.Link_raw, "#"); i != -1 {
					panic("IDREF has fragment ID")
				}
				println("IDREF:", pGL.Link_raw)
			} else {
				pGL.AddressMode = "uri"
				if i := S.Index(pGL.Link_raw, "#"); i != -1 {
					pGL.FragID = pGL.Link_raw[i+1:]
					pGL.RelFP = pGL.Link_raw[:i]
				} else {
					pGL.RelFP = pGL.Link_raw
				}
				L.L.Dbg("URI: " + pGL.RelFP + "#" + pGL.FragID)
				// p.AbsFP = FU.RelFilePath(FP.Join(
				// 	pGF.InputFile.FileFullName.Echo(), p.RelFP.S())).AbsFP()
				s, _ := FP.Abs(FP.Join(p.PathProps.AbsFP.S(), pGL.RelFP))
				pGL.AbsFP = FU.AbsFilePath(s)
				L.L.Dbg("URI AbsFP: " + FU.Enhomed(pGL.AbsFP.S()))
			}
			switch pGL.Att {
			// Link SOURCES
			case "idref", "idrefs", "href", "conref",
				"keyref", "data-conref", "data-keyref":
				pGL.IsRef = true
				if pGL.AddressMode == "key" {
					p.OutgoingKeys = append(p.OutgoingKeys, pGL)
				} else {
					p.OutgoingURIs = append(p.OutgoingURIs, pGL)
				}
				// Link TARGETS
			case "id", "keys", "data-keys":
				if pGL.AddressMode == "key" {
					p.IncomableKeys = append(p.IncomableKeys, pGL)
				} else {
					p.IncomableURIs = append(p.IncomableURIs, pGL)
				}
			}
		}
	}
	return p
}
