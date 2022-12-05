package mcfile

import (
	"fmt"

	LU "github.com/fbaube/logutils"
	L "github.com/fbaube/mlog"
)

func (p *Contentity) LogPrefix(mid string) string {
	return fmt.Sprintf("%02d%sst%s", p.logIdx, mid, p.logStg)
}

func (p *Contentity) L(level LL, format string, a ...interface{}) {
	// L.L.Log(level, format, a...)
	L.L.LogWithString(LU.Level(level), format,
		fmt.Sprintf("F%02d", p.logIdx)+"|stg"+p.logStg, a...)
}

func (p *Contentity) LogTextQuote(level LL, textquote string, format string, a ...interface{}) {
	// L.L.Log(level, format, a...)
	L.L.LogWithString(LU.Level(level), format,
		fmt.Sprintf("%02d", p.logIdx)+","+p.logStg, a...)
	panic("FIXME")
	// L.L.LogMultilineAsIs(SU.IndentWith("   |  ", textquote))
}

type LL LU.Level

var LDbg, LProgress, LInfo, LOkay, LWarning, LError, LPanic LL

func init() {
	LDbg = LL(LU.LevelDbg)
	LProgress = LL(LU.LevelProgress)
	LInfo = LL(LU.LevelInfo)
	LOkay = LL(LU.LevelOkay)
	LWarning = LL(LU.LevelWarning)
	LError = LL(LU.LevelError)
	LPanic = LL(LU.LevelPanic)
}
