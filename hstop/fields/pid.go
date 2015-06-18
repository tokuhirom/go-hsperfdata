package fields

import (
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hstop/core"
)

type PidField struct {
}

func (*PidField) GetTitle() string {
	return "PID"
}

func (*PidField) GetWidth() int {
	return 5
}

func (*PidField) Render(state *core.State) string {
	return fmt.Sprintf("%v", state.Pid)
}
