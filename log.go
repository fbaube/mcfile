package mcfile

import (
	"fmt"

	L "github.com/fbaube/mlog"
)

func (p *Contentity) LogPrefix(mid string) string {
	return fmt.Sprintf("%02d%sst%s", p.logIdx, mid, p.logStg)
}

func (p *Contentity) L(level LL, format string, a ...interface{}) {
	// L.L.Log(level, format, a...)
	L.L.LogWithString(L.Level(level), format,
		fmt.Sprintf("F%02d", p.logIdx)+"|stg"+p.logStg, a...)
}

func (p *Contentity) LogTextQuote(level LL, textquote string, format string, a ...interface{}) {
	// L.L.Log(level, format, a...)
	L.L.LogWithString(L.Level(level), format,
		fmt.Sprintf("%02d", p.logIdx)+","+p.logStg, a...)
	panic("FIXME")
	// L.L.LogMultilineAsIs(SU.IndentWith("   |  ", textquote))
}

type LL L.Level

var LDbg, LProgress, LInfo, LOkay, LWarning, LError, LPanic LL

func init() {
	LDbg = LL(L.LevelDbg)
	LProgress = LL(L.LevelProgress)
	LInfo = LL(L.LevelInfo)
	LOkay = LL(L.LevelOkay)
	LWarning = LL(L.LevelWarning)
	LError = LL(L.LevelError)
	LPanic = LL(L.LevelPanic)
}
