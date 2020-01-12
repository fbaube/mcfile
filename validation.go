package mcfile

import (
	S "strings"

	MU "github.com/fbaube/miscutils"
)

// TODO If no DOCTYPE, make a guess based on Filext but it can't be fatal.
func (p *MCFile) DoValidation() (dtdS string, docS string, errS string) {
	errS = ""
	pX := p.TheXml()
	if !p.IsXML {
		panic("DoValidation !IsXML")
	}
	/*
		if pX.XmlFileMeta == nil {
			return "No_XFM!", "valid??", ""
		} */
	if pX.XmlDoctype == nil {
		return "No_DcTp", "valid??", ""
	}
	var ppid = pX.XmlDoctype.XmlPublicID
	// print("\t" + ppid.PTDesc + " --> ")
	val := CA.XmlCatalog.GetByPublicID(ppid.String())
	if val == nil {
		// println("DTD NOT FOUND")
		return "DTD_Unk", "valid??", ""
	}
	// print("DTD OK; ")
	// NOTE At CLI, can use :: alias validate-lw-topic =
	// 'xmllint --noout --dtdvalid file:///LwDTD/lw-topic.dtd'
	// NOTE that if we don't specify the DTD, then xmllint only checks
	// for well-formedness, but we've already done this ourselves when
	// building the GTree. So, we should specify the DTD when we invoke
	// xmllint, or else not even bother.
	// NOTE that we have to pass the flag "--nowarning", or else
	// validation will fail if the SYSTEM ID can't be found.
	stdOut, stdErr, err := MU.RunCommand(
		"xmllint", "--noout", "--nowarning", "--nonet", "--dtdvalid",
		"file://"+string(val.AbsFilePath), p.AbsFilePathParts.String())
	// NOTE that the return value "err" is dumb:
	// it contains stuff like "exit status 3".
	if S.TrimSpace(stdErr) == "" {
		// print("Document is valid \n")
		return "isFound", "isValid", ""
	}
	// print("Validation failed: ")
	if S.TrimSpace(stdOut) != "" {
		errS += stdOut + "\n"
	}
	errS += stdErr
	errS += "==> End of validation errors <== \n"
	if err != nil {
		errS += "==> xmllint: " + err.Error() + "\n"
	}
	return "isFound", "Errors!", errS
}
