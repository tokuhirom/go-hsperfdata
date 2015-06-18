package fields

import (
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hstop/core"
)

type MainClassField struct {
}

func (*MainClassField) GetTitle() string {
	return "MAIN-CLASS"
}

func (*MainClassField) GetWidth() int {
	return 15
}

func shorten_main_class_name(class_name string, length int) string {
	if len(class_name) > length {
		return class_name[len(class_name)-length:]
	} else {
		return class_name
	}
}

func (*MainClassField) Render(state *core.State) string {
	return fmt.Sprintf("%15s", shorten_main_class_name(state.Result.GetProcName(), 15))
}
