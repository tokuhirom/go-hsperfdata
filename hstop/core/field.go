package core

type Field interface {
	GetTitle() string
	GetWidth() int
	Render(*State) string
}
