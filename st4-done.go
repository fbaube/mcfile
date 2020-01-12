package mcfile

// st4_Done does final cleanup and beautification. 
func (p *MCFile) st4_Done() *MCFile {
	if p.GetError() != nil {
		return p
	}
	println("--> (4) Done")
	switch p.FileType() {
	case "XML":
		println("TODO> 4. Done XML")
	case "MKDN":
		println("TODO> 4. Done MKDN")
	}
	return p
}

/*
	switch p.TypeSpecific.(type) {
	case TypeMkdn:
		if p.MType[0] != "md" {
			s := "Not markdown (Mtype[0]!=\"md\") ?!"
			logerr.Println(s)
			p.SetError(errors.New(s))
		}
	case TypeXml:
		if !(p.IsXML && p.MType[0] == "xml") {
			s := "Not XML (Mtype[0]!=\"xml\") ?!"
			logerr.Println(s)
			p.SetError(errors.New(s))
		}
	default:
		panic("SanityCheck failed")
	}
	return p
}
*/
