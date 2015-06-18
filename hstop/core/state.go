package core

import linuxproc "github.com/c9s/goprocinfo/linux"
import "github.com/tokuhirom/go-hsperfdata/hsperfdata"

type State struct {
	Pid     string
	Process *linuxproc.ProcessStat
	Result  *hsperfdata.Result
}
