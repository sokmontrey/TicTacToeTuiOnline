package page

type PageMsg any

type PageCmd int

const (
	ProgramQuit PageCmd = iota
)

type Page interface {
	Init() PageCmd
	Update(msg PageMsg) PageCmd
	View() string
}
