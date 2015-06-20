package fields

import (
	"fmt"

	"github.com/tokuhirom/go-hsperfdata/hstop/state"
)

type NiceField struct {
}

func (*NiceField) GetTitle() string {
	return "NI"
}

func (*NiceField) GetWidth() int {
	return 3
}

func (*NiceField) Render(state *state.State) string {
	return fmt.Sprintf("%3d", state.GetNice())
}
