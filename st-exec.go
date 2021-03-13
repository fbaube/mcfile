package mcfile

import (
	"runtime/debug"

	"github.com/fatih/color"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// ExecuteStages processes an Contentity to completion in an isolated thread,
// and can eaily be converted to run as a goroutine.
//
// Package mlog has been added. HOWEVER an interesting question is,
// how is an error indicated and a thread terminated prematurely ?
// One method was to set the field `MCFile.CheckedPath.error` to
// non-`nil`, which has to be checked for at the start of functions.
// Another way might be to pass in a `Context` and use its
// cancellation capability. Yet another way might be to `panic(..)``,
// and so this function already has code to catch a panic.
//
func (p *Contentity) ExecuteStages() *Contentity {
	if p.GetError() != nil {
		return p
	}
	p.logStg = "--"
	defer func() {
		if r := recover(); r != nil {
			L.L.Panic(SU.Rfg(SU.Ybg(" ** PANIC caught in ExecuteStages ** ")))
			color.Set(color.FgHiRed)
			// L.L.Panic(e.Error())
			color.Unset()
			var sRecovered string
			var eRecovered error
			L.L.Info("recover got a: %T", r)
			switch r.(type) {
			case string:
				sRecovered = r.(string)
			case error:
				eRecovered = r.(error)
			}
			if sRecovered == "" {
				sRecovered = eRecovered.Error()
			}
			L.L.Panic(SU.Rfg(SU.Ybg("=== PANICKED ===")))
			L.L.Error("Recovered in MCFile.ExecuteStages(): " + sRecovered)
			L.L.Error("Stacktrace from panic: " + string(debug.Stack()))
		}
	}()
	// Execute stages/steps
	if p.IsDir() {
		// L.L.Info("Is a dir: skipping all stages")
		p.L(LInfo, "Is a dir: skipping content processing")
		return p
	}
	// println("--> DOING STAGES FOR:", p.AbsFP())
	return p.
		st0_Init().
		st1_Read().
		st2_Tree().
		st3_Refs() /* .
		// Stage/Step 4 works on all input files at once !
		st4_Done() */
}
