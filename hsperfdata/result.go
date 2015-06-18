package hsperfdata

import "strings"

type Result struct {
	data map[string]interface{}
}

func (self *Result) GetMap() map[string]interface{} {
	return self.data
}

func (self *Result) GetProcName() string {
	javaCommand := self.data["sun.rt.javaCommand"]
	if javaCommand != nil {
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
