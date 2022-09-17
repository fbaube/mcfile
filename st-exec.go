package mcfile

import (
	"runtime/debug"
)

// ExecuteStages processes a Contentity to completion in an isolated
// thread, and can eaily be converted to run as a goroutine.
//
// Package mlog has been added. HOWEVER an interesting question is,
// how is an error indicated and a thread terminated prematurely ?
// One method was to set the field `Contentity.Err`, which has to
// be checked for at the start of functions. Another way might be
// to pass in a `Context` and use its cancellation capability. Yet
// another way might be to `panic(..)“, and so this function already
// has code to catch a panic.
func (p *Contentity) ExecuteStages() *Contentity {
	if p.PathProps == nil {
		panic("ExecuteStages :: nil PP")
	}
	if len(p.PathProps.Raw) == 0 {
		panic("ExecuteStages :: ZERO")
	}
	if p.HasError() {
		p.L(LInfo, "Has error: skipping")
		return p
	}
	if p.FileType() == "BIN" {
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
		st3_Refs() /* .
		// Stage/Step 4 works on all input files at once !
		st4_Done() */
}
