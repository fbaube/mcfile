package mcfile

import (
	"log"
	"os"
)

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
