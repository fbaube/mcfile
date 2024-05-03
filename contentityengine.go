package mcfile

// ContentityEngine tracks the (oops, global) state of a
// ContentityFS tree being assembled, for example when a
// directory is specified for recursive analysis.
// 
// FIXME: ID assignment should be offloaded to the DB ?
// .
type ContentityEngine struct {
	// nexSeqID should be reset to 0 when starting another tree ?
	// No, because every single entity (dir/file) gets one,
	// even if it is listed on the CLI as an individual file.
	nexSeqID      int
	rootPath      string
	// summaryString StringFunc // for ON.NordEngine but not this
}

// CntyEng is a package global, which is dodgy and not re-entrant.
// The solution probably involves currying.
// 
// NOTE: Is the call to new(..) unnecessary? This variable
// should NOT be reinitialized for every new ContentityFS. 
var CntyEng *ContentityEngine = new(ContentityEngine)

