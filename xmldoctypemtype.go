package mcfile

import (
	S "strings"
)

func (p *MCFile) SetMtypeUsingDeclaredDoctype() *MCFile {
	pX := p.TheXml()
	// println("    --> SetMtypeUsingDeclaredDoctype... init...")
	// println("      MMCtype: ", p.MMCstring())
	println("    --> in-Mtype:", p.Mstring())
	println("    --> Sniffed:", p.SniftMimeType)
	// println("MagicMimeType:", p.MagicMimeType)
	println("    --> DDT-xfm:", pX.TopTag, ",", pX.PublicTextDesc)

	var PubDescIsLwDita bool
	var PDU = S.ToUpper(pX.PublicTextDesc)
	PubDescIsLwDita = S.Contains(PDU, "DITA") &&
		(S.Contains(PDU, "LW") || S.Contains(PDU, "LIGHTWEIGHT"))
	if PubDescIsLwDita {
		// xml/dita/topic
		p.CheckedContent.MType[0] = "xml"
		p.CheckedContent.MType[1] = "lwdita"
		p.CheckedContent.MType[2] = S.ToLower(pX.TopTag)
	}
	println("    --> outMtype:", p.Mstring())
	return p
}
