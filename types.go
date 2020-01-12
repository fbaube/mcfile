package mcfile

import (
	"log"
	"os"
)

type SystemID string
type PublicID string

type StringTally map[string]int

/*
func (tt TagTally) String() string {
	return "tt"
}
*/
var logerr *log.Logger

func init() {
	logerr = log.New(os.Stderr, "ERR:mcfile> ", log.Lshortfile)
}

type XmlContype string

// XmlContypes, maybe DTDmod should be DTDelms.
var XmlContypes = []XmlContype{"Unknown", "DTD", "DTDmod", "DTDent",
	"RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
