package mcfile

import (
	"os"
	"fmt"
	"github.com/fbaube/gtree"
	"github.com/fbaube/gparse"
)

// st2_Tree takes the output of st1_Read - which at a minimum
// is a complete set of `GToken`s - and creates a `GTree`.
// Note that st1_Read might have generated an CST (MKDN and
// HHTML do this) but in such cases, st1_Read also prepared
// the list of corresponding `GToken`s.
// - PrepareToTree() // e.g. GTags
// - ParseIntoTree()
// - PostTreeMeta()
// - NormalizeTree()
func (p *MCFile) st2_Tree() *MCFile {
	if p.GetError() != nil {
		return p
	}
	println("--> (2) Tree")
	return p.st2a_PrepareToTree().ParseIntoTree().PostTreeMeta().NormalizeTree()
}

// PrepareToTree is Step 2a. <br/>
// This is used when there is some preparation specific to building the tree.
func (p *MCFile) st2a_PrepareToTree() *MCFile {
	if p.GetError() != nil {
		return p
	}
	var e error
	switch p.FileType() {
	case "XML":
		// G-Tag-ify PRE-TREE
		// pX := p.TheXml()
		p.GTags, e = gtree.MakeGTagsFromGTokens(p.GTokens)
		if e != nil {
			p.SetError(e)
			return p // errors.Wrap(e, "MakeGTagsFromGTokens")
		}
	case "MKDN":
		p.GTags, e = gtree.MakeGTagsFromGTokens(p.GTokens)
		if e != nil {
			p.SetError(e)
			return p // errors.Wrap(e, "MakeGTagsFromGTokens")
		}
	case "HTML":
		p.GTags, e = gtree.MakeGTagsFromGTokens(p.GTokens)
		if e != nil {
			p.SetError(e)
			return p // errors.Wrap(e, "MakeGTagsFromGTokens")
		}
	}
	return p
}

// ParseIntoTree is Step 2b.
func (p *MCFile) ParseIntoTree() *MCFile {
	if p.GetError() != nil {
		return p
	}
	var e error
	fmt.Printf("==> mcfl.st2b: FileType<%s> nGTags<%d> \n",
		p.FileType(),len(p.GTags))
	switch p.FileType() {
	case "XML":
		// TREE
		// ProcessXml(p) // This does the work
		// G-Tree-ify DO-TREE
		p.GTree, e = gtree.NewGTreeFromGTags(p.GTags)
		if e != nil {
			p.SetError(e)
			println("==> mcfl.st2b: Error!:", e)
			return p // errors.Wrap(e, "NewGTreeFromGTags")
		}
	case "MKDN", "HTML":
		// TREE
		p.GTree, e = gtree.NewGTreeFromGTags(p.GTags)
		if e != nil {
			p.SetError(e)
			println("==> mcfl.st2b: Error!:", e)
			return p // errors.Wrap(e, "NewGTreeFromGTags")
		}
	default:
		println("==> mcfl.st2b: bad FileType:", p.FileType)
	}
	if p.GTree == nil {
		println("==> mcfl.st2b: NIL Gtree !!")
	}
	if p.GTree != nil {
		// println(p.GTree.String())
		// gparse.DumpTo(p.GTree, os.Stdout)
		gparse.DumpTo(p.GTokens, os.Stdout)
	}
	return p
}

// PostTreeMeta is Step 2c. <br/>
// This is used to process metadata that is best handled
// only after tree-building. For example, XML `DOCTYPE`.
func (p *MCFile) PostTreeMeta() *MCFile {
	if p.GetError() != nil {
		return p
	}
	switch p.FileType() {
	case "XML":
		println("TODO> 2c. NormalizeTree XML ==> BIG!")
		/*
			// =========================
			//   XML ANALYSIS and also
			//    BTW get the DOCTYPE
			// =========================
			pX.GTokens = gparse.XmlCheckForPreambleToken(pX.GTokens)
			e = p.ProcessMetaGetDoctype()
			if e != nil {
				return errors.Wrap(e, "ProcessMetaGetDoctype")
			}

				// println("\t Got XFM/DT:", pGF.XmlFileMeta.String())
				if pX.DoctypeIsDeclared {
					e = SetMtypeUsingDeclaredDoctype(p)
					if e != nil {
						return errors.Wrap(e, "SetMtypeUsingDeclaredDoctype")
					}
					// FIXME!! println("    --> Declared DOCTYPE:", pX.String())
				}
				e = p.RefineDirectives()
				if e != nil {
					return errors.Wrap(e, "RefineDirectives")
				}
		*/
	case "MKDN":
		println("TODO> 2c. NormalizeTree MKDN")
	case "HTML":
		println("TODO> 2c. NormalizeTree HTML")
	}
	return p
}

// NormalizeTree is Step 2d. <br/>
// This is used when we have a tree that is generated by a library,
// and which therefore must still be converted into a `GTree`.
func (p *MCFile) NormalizeTree() *MCFile {
	if p.GetError() != nil {
		return p
	}
	switch p.FileType() {
	case "XML":
		println("TODO> 2d. NormalizeTree XML ==> ENTs, etc.!")
		/*
			e = p.DoEntitiesList()
			if e != nil {
				return errors.Wrap(e, "DoEntitiesList")
			}
			e = p.SubstituteEntities()
			if e != nil {
				return errors.Wrap(e, "SubstituteEntities")
			}
		*/
	case "MKDN":
		println("TODO> 2d. NormalizeTree MKDN")
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
