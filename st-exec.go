package mcfile

import (
       // "log"
	"runtime/debug"

	// "github.com/fbaube/must"
	SU "github.com/fbaube/stringutils"
)

// ExecuteStages processes a Contentity to completion in an isolated
// thread, and can eaily be converted to run as a goroutine. Summary:
//   - st0_Init()
//   - st1_Read()
//   - st2_Tree()
//   - st3_Refs()
//   - st4_Done() (not currently called, but
//     will work on all input files at once !)
//
// An interesting question is, how can we indicate an error and
// terminate a thread prematurely ? The method currently chosen
// is to use interface [github.com/fbaube/miscutils/Errer].
// This has to be checked for at the start of a func.
//
// We could also pass in a `Context` and use its cancellation
// capability. Yet another way might be simply to `panic`,
// and so this function already has code to catch a panic.
// .
func (p *Contentity) ExecuteStages() *Contentity {
	// The E family of functions all remove a final error return,
	// panicking if non-nil.
	// Handle converts such a panic to a returnable error value.
	// Other panics are not recovered.
	// defer must.F(log.Fatal)

	if p.MarkupType() == "UNK" {
		panic("UNK MarkupType in ExecuteStages")
	}
	if p.FSItem.Raw == "" {
		p.L(LWarning, "ExecuteStages :: ZERO-len Raw")
		return p
	}
	if p.HasError() {
		p.L(LInfo, "Has error: skipping")
		return p
	}
	if p.MarkupType() == SU.MU_type_BIN {
		p.L(LWarning, "Skipping ALL stages for binary file")
		return p
	}
	// p.L(LInfo, "LENGTH %d SIZE %d", len(p.FSItem.Raw), p.Size())
	if len(p.FSItem.Raw) == 0 { // p.Size() == 0 {
		p.L(LWarning, "Skipping ALL stages for empty file")
		return p
	}
	if p.IsDirlike() {
		p.L(LInfo, "Is dir or similar: skipping content processing")
		return p
	}
	p.logStg = "--"
	defer func() {
		if r := recover(); r != nil {
			p.L(LPanic, "= = = = = = = = = = = = = = = = = = = =")
			p.L(LPanic, " ** PANIC caught in ExecuteStages ** ")
			var sRecovered string
			var eRecovered error
			p.L(LInfo, "recover() got a: %T", r)
			switch r.(type) {
			case string:
				sRecovered = r.(string)
			case error:
				eRecovered = r.(error)
			}
			if sRecovered == "" {
				sRecovered = eRecovered.Error()
			}
			p.L(LError, "defer'd-recover()-string: "+sRecovered)
			debug.PrintStack()
			/* Barfed!
			bb := debug.Stack()
			ss := string(bb)
			p.L(LError, "stdlib:runtime/debug.Stack(): "+
				ss) // string(debug.Stack()))
			*/
			p.L(LError, "= = = = = = = = = = = = = = = = = = = =")
		}
	}()
	// Execute stages/steps
	return p.
		st0_Init().
		st1_Read().
		st2_Tree().
		st3_Refs()
	// Stage/Step 4 works on all input files at once !
	// st4_Done()
}
