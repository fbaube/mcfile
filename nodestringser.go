package mcfile

type NodeStringser interface {
	NodeEcho(int) string
	NodeInfo(int) string
	NodeDebug(int) string
	NodeCount() int
}
