package hsperfdata

import (
	"fmt"
	"strings"
	"time"
)

type Result struct {
	data         map[string]interface{}
	modTimestamp int64
}

func (self *Result) GetMap() map[string]interface{} {
	return self.data
}

func (self *Result) GetTimestamp() (time.Time, error) {
	if createVmBeginTime, ok := self.data["sun.rt.createVmBeginTime"]; ok {
		if now, ok := createVmBeginTime.(int64); ok {
			return time.Unix(0, (now+self.modTimestamp)*int64(time.Millisecond)), nil
		} else {
			return time.Time{}, fmt.Errorf("sun.rt.createVmBeginTime wasn't an int64: %#v", createVmBeginTime)
		}
	} else {
		return time.Time{}, fmt.Errorf("couldn't find sun.rt.createVmBeginTime to add prologue ModTimestamp %v to", self.modTimestamp)
	}
}

func (self *Result) GetProcName() string {
	if javaCommand, ok := self.data["sun.rt.javaCommand"]; ok {
		if str, ok := javaCommand.(string); ok {
			splitted := strings.SplitN(str, " ", 2)
			return splitted[0]
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (self *Result) GetString(key string) string {
	v := self.data[key]
	if v != nil {
		if str, ok := v.(string); ok {
			return str
		} else {
			return ""
		}
	} else {
		return ""
	}
}
