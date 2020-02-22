package mcfile

import (
	"encoding/xml"
	"fmt"
	"sort"
	S "strings"
	"github.com/fbaube/gtoken"
)

var GlobalTagTally StringTally
var GlobalAttTally StringTally
var GlobalTagCount int
var GlobalAttCount int

func init() {
	GlobalTagTally = make(map[string]int)
	GlobalAttTally = make(map[string]int)
}

func (p *MCFile) TallyTags() {
	p.TagTally = make(map[string]int)
	p.AttTally = make(map[string]int)
	for _, pGT := range p.GTokens {
		if pGT.GName.Local == "" {
			continue
		}
		AddInGName(p.TagTally, p.AttTally, pGT)
		AddInGName(GlobalTagTally, GlobalAttTally, pGT)
		GlobalTagCount += 1
		GlobalAttCount += len(pGT.GAtts)
	}
}

func AddInGName(ElmT StringTally, AttT StringTally, gT *gtoken.GToken) {
	// Is there an entry for the tag yet ?
	// if val, ok := dict["foo"]; ok {
	var ok bool
	var n int
	if n, ok = ElmT[gT.GName.Local]; ok {
		ElmT[gT.GName.Local] = n + 1
	} else {
		ElmT[gT.GName.Local] = 1
	}
	// Now process the attributes
	for _, A := range gT.GAtts { // A is a *GAtt i.e. *xml.Attr
		var gat gtoken.GAtt = A // *A
		var xat = xml.Attr(gat)
		var sat = xat.Name.Local
		if n, ok = AttT[sat]; ok {
			AttT[sat] = n + 1
		} else {
			AttT[sat] = 1
		}
	}
}

func (st StringTally) StringSortedValues() string {
	n := map[int][]string{}
	var a []int
	for k, v := range st {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	var sb S.Builder
	for _, k := range a {
		for _, s := range n[k] {
			sb.WriteString(fmt.Sprintf("%s:%d,", s, k))
		}
	}
	return sb.String()
}
