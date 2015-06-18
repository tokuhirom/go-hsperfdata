package fields

import (
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hstop/core"
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

func (*VsizeField) Render(state *core.State) string {
	return fmt.Sprintf("%8s", support.Size(state.Process.Vsize))
}
