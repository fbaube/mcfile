package mcfile

import (
	"errors"
	"fmt"

	SU "github.com/fbaube/stringutils"
)

// GetYamlHeader extracts YAML frontmatter in MDITA-XP.
// IMPORTANT NOTE: If a metadata block it found,
// it is trimmed from the content in field `Raw`.
// The metadata is unmarshalled into a map (i.e. a `PropSet`),
// so variables can be freely added, but there can be no error
// checking or checking for required fields.
//
func (p *MCFile) GetYamlHeader() *MCFile {
	if p.GetError() != nil {
		return p
	}
	if p.FileType() != "MKDN" {
		panic("Tried to get YAML metadata header for non-MKDN")
	}
	var yps SU.PropSet // *SU.YamlMeta
	var rem string
	var e error
	yps, rem, e = SU.GetYamlMetadataAsPropSet(p.Raw)
	fmt.Printf("(D) mcfl.yaml: nProps<%d> Rem'ng:< %s >\n", len(yps), rem)
	if e != nil {
		println("--> YAML error:", e.Error())
		return p
	}
	// if p.Header != nil {
	println("--> YAML header reallocated!")
	p.SetError(errors.New("YAML header reallocated"))
	return p
	// }
	// p.Header = new(Header)
	p.MetaFormat = "yaml"
	p.MetaProps = yps
	if rem != "" {
		println("OK!: Yaml caused re-assignment of MCFile Raw")
		p.Raw = rem
	} else {
		println("ERR: Yaml left no MCFile Raw")
	}
	/*
			m := ymb.AsMap()
				if p.Props == nil {
					p.Props = m
				} else {
					for k, v := range m {
		    	p.Props[k] = v
				}

				println("mcfile got YAML. TODO:80 Add it to MCFile.")
		* 		fmt.Printf("YAML: %+v \n", *ymb)
		* 	}
		* 	}
	*/
	return p
}
