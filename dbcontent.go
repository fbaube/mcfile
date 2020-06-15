package mcfile

import (
  "github.com/fbaube/db"
)

/*
type Content struct {
	Idx int         // `db:"idx_content"`
	Idx_Inbatch int // `db:"idx_inbatch"`
	BaseInfo
	Meta_raw string
	Text_raw string
	Analysis
}
type BaseInfo struct {
	RelFilePath string
	AbsFilePath FU.AbsFilePath `db:"absfilepath"` // necessary ceremony
	Creatime string // ISO-8601 / RFC 3339
}
type Analysis struct {
	MimeType    string
	Mtype       string
	RootTag     string
	RootAtts    string // e.g. <html lang="en">
	XmlContype  string
	XmlDoctype  string
	DitaContype string
}
*/

// AsDBContent adds a content item (i.e. a file) to the DB.
func (p *MCFile) AsDBContent() (pC *db.Content, e error) {
  pC = new(db.Content)
  pC.Idx = p.Idx
  pC.Idx_Inbatch = p.Idx_Inbatch
  pC.Times = p.Times
  // BaseInfo
  pC.RelFilePath = p.CheckedContent.RelFilePath
  pC.AbsFilePath = p.CheckedContent.AbsFilePath
  // pC.Creatime =
  // As-is
  pC.Meta_raw = p.Meta_raw
  pC.Text_raw = p.Text_raw
  // Analysis
  // pC.MimeType    =
  // pC.Mtype          = p.MType.String()
  // pC.RootTag     =
  // pC.RootAtts    =
  // pC.XmlContype  =
  // pC.XmlDoctype  =
  // pC.DitaContype = p.DitaInfo.

  return pC, nil
}