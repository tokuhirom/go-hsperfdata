package fields

import (
	"fmt"

	"github.com/tokuhirom/go-hsperfdata/hstop/state"
	"github.com/tokuhirom/go-hsperfdata/hstop/support"
)

type RssField struct {
}

func (*RssField) GetTitle() string {
	return "RSS"
}

func (*RssField) GetWidth() int {
	return 8
}

func (*RssField) Render(state *state.State) string {
	return fmt.Sprintf("%8s", support.Size(uint64(state.GetRss())))
}
