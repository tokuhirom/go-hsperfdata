package fields

import (
	"fmt"

	"github.com/tokuhirom/go-hsperfdata/hstop/core"
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

func (*RssField) Render(state *core.State) string {
	return fmt.Sprintf("%8s", support.Size(uint64(state.Process.Rss)))
}
