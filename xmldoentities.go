package mcfile

func (p *MCFile) ProcessEntities_() error {
	pX := p.TheXml()
	// p := new(XmlItems)
	// var pXI *XmlItems
	// pXI = new(XmlItems)
	// ## pGFile.IDinfo = new(mcfile.IDinfo)
	// pGFile.XmlItems = pXI
	var e error

	// PASS 5
	// Fully process ENTITY declarations and build lists
	pX.GEnts, e = p.NewEntitiesList()
	if e != nil {
		panic("Pass 5: ENTITY processing failed")
	}
	println("==> Pass 5 OK: Entities are collected")
	/*
		if len(pRTx.DEnts) > 0 {
			for e, ent := range pRTx.DEnts {
				ilog.Printf("RTx' EntDef: || %s || %s ||", e, ent)
			}
		}
	*/

	// PASS 6
	// Entity substitutions, recursively.
	// We do this before ELEMENTs and ATTLISTs are processed.
	// TODO 
	e = p.SubstituteEntities()
	if e != nil {
		panic("Pass 6: Recursive ENTITY substitution failed: " + e.Error())
	}
	println("==> Pass 6 ??: TODO: Entities are not substituted")

	// PASS 7
	// If there's a root element, then build a DOM tree that is
	// block/inline-aware, mixed content.
	/*
		if pET.GotRootTag {
			println("NewBlockTree()... TODO!!")
			r = S.NewReader(string(pET.raw)) // (inString)
			// rootNode, err = NewNodetree(r, nil, os.Stdout)
			if e != nil {
				// panic(inFileName + ": Attempt to build DOM tree failed")
			}
		}
	*/

	return nil
}
