package mcfile

import (
	SU "github.com/fbaube/stringutils"
	"runtime/debug"
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
// is to use interface [github.com/fbaube/utils/miscutils/Errer].
// This has to be checked for at the start of a func.
//
// We could also pass in a `Context` and use its cancellation
// capability. Yet another way might be simply to `panic(..)“,
// and so this function already has code to catch a panic.
// .
func (p *Contentity) ExecuteStages() *Contentity {
	if p.PathProps.Raw == "" {
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
	// p.L(LInfo, "LENGTH %d SIZE %d", len(p.PathProps.Raw), p.Size())
	if len(p.PathProps.Raw) == 0 { // p.Size() == 0 {
		p.L(LWarning, "Skipping ALL stages for empty file")
		return p
	}
	if p.IsDir() {
		p.L(LInfo, "Is a dir: skipping content processing")
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
			p.L(LError, "stdlib:runtime/debug.Stack(): "+string(debug.Stack()))
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
