package state

import linuxproc "github.com/c9s/goprocinfo/linux"
import "github.com/tokuhirom/go-hsperfdata/hsperfdata"
import (
	"fmt"
)

type State struct {
	pid     string
	process *linuxproc.ProcessStat
	Result  *hsperfdata.Result
}

func New(pid string, result *hsperfdata.Result) (*State, error) {
	processStat, err := linuxproc.ReadProcessStat("/proc/" + pid + "/stat")
	if err != nil {
		return nil, fmt.Errorf("Cannot read /proc/%d/stat: %v", pid, err)
	}
	return &State{pid, processStat, result}, nil
}

func (self *State) GetPid() string {
	return self.pid
}

func (self *State) GetVsize() uint64 {
	return self.process.Vsize
}

func (self *State) GetNice() int64 {
	return self.process.Nice
}

func (self *State) GetRss() int64 {
	return self.process.Rss
}

func (self *State) GetVMInfo() string {
	vendor := self.Result.GetString("java.property.java.vm.vendor")
	version := self.Result.GetString("java.property.java.version")
	return vendor[:1] + version[2:]
}
