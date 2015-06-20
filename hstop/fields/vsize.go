package fields

import (
	"fmt"

	"github.com/tokuhirom/go-hsperfdata/hstop/state"
	"github.com/tokuhirom/go-hsperfdata/hstop/support"
)

type VsizeField struct {
}

func (*VsizeField) GetTitle() string {
	return "VSIZE"
}

func (*VsizeField) GetWidth() int {
	return 8
}

func (*VsizeField) Render(state *state.State) string {
	return fmt.Sprintf("%8s", support.Size(state.GetVsize()))
}
