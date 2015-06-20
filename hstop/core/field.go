package core

import state "github.com/tokuhirom/go-hsperfdata/hstop/state"

type Field interface {
	GetTitle() string
	GetWidth() int
	Render(*state.State) string
}
