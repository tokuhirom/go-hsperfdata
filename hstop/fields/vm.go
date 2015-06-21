package fields

import (
	"fmt"

	"github.com/tokuhirom/go-hsperfdata/hstop/state"
)

type VMField struct {
}

func (*VMField) GetTitle() string {
	return "VM"
}

func (*VMField) GetWidth() int {
	return 7
}

func (*VMField) Render(state *state.State) string {
	return fmt.Sprintf("%7s", state.GetVMInfo())
}
