package mcfile

import (
	"fmt"
	"runtime/debug"

	SU "github.com/fbaube/stringutils"
)

// ExecuteStages runs an `MCFile` to completion in an isolated thread, and
// can eaily be converted to run as a goroutine. An interesting question
// is, how is an error indicated and a thread terminated prematurely ?
// The current method is to set the field `MCFile.CheckedPath.error` to
// non-`nil`, which has to be checked for at the start of functions. Another
// way might be to pass in a `Context` and use its cancellation capability.
// Yet another way might be to `panic(..)``, and so this function already
// has code to catch a panic.
// func (p *MCFile) ExecuteStages() *MCFile {
func (p *Contentity) ExecuteStages() *Contentity {
	if p.GetError() != nil {
		return p
	}
	defer func() {
		if r := recover(); r != nil {
			/*
					fmt.Fprintf(w, SU.Rfg(SU.Ybg(" ** ERROR ** ")))
					color.Set(color.FgHiRed)
					fmt.Fprintf(w, "\n" + e.Error() + "\n")
					color.Unset()
				}
			*/
			var sRecovered string
			var eRecovered error
			fmt.Printf("recover got: %T \n", r)
			switch r.(type) {
			case string:
				sRecovered = r.(string)
			case error:
				eRecovered = r.(error)
			}
			if sRecovered == "" {
				sRecovered = eRecovered.Error()
			}
			println(SU.Rfg(SU.Ybg("\n=== PANICKED ===")))
			println("Recovered in MCFile.ExecuteStages():", sRecovered)
			println("Stacktrace from panic: \n", string(debug.Stack()))
			/* // //
			p.LogIt("\n\t=== PANICKED ===")
			color.Set(color.FgHiRed)
			p.Blare("\nRecovered in MCFile.ExecuteStages(): " + sRecovered + "\n")
			p.Blare("Stacktrace from panic: \n" + string(debug.Stack()))
			color.Unset()
			*/
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
