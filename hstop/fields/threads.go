package fields

import (
	"fmt"
	"github.com/tokuhirom/go-hsperfdata/hstop/core"
)

type ThreadsField struct {
}

func (*ThreadsField) GetTitle() string {
	return "T#"
}

func (*ThreadsField) GetWidth() int {
	return 3
}

func (*ThreadsField) Render(state *core.State) string {
	return fmt.Sprintf("%3s", state.Result.GetString("java.threads.live"))
}
