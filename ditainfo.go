package mcfile

// DitaInfo is TBS.
type DitaInfo struct {
	// These next two are "" IFF the file is not DITA/LwDITA.
	ditaML      DitaML
	ditaContype DitaContype
}

type DitaML string
type DitaContype string

var DitaMLs = []DitaML{"1.2", "1.3", "XDITA", "HDITA", "MDATA"}
var DitaContypes = []DitaContype{"Map", "Bookmap", "Topic", "Task", "Concept",
	"Reference", "Dita", "Glossary", "Conrefs", "LwMap", "LwTopic"}

func (di DitaInfo) String() string { return "" }

func (di DitaInfo) DString() string {
	return "<-- DITA " + string(di.ditaML) +
		" " + string(di.ditaContype) + " -->\n"
}
