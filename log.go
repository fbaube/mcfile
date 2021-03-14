package mcfile

import (
	"fmt"

	L "github.com/fbaube/mlog"
)

func (p *Contentity) L(level LL, format string, a ...interface{}) {
	// L.L.Log(level, format, a...)
	L.L.LogWithString(L.Level(level), format, fmt.Sprintf("%02d", p.logIdx)+","+p.logStg, a...)
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
