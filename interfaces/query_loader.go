package interfaces

type QueryLoader interface {
	LoadQuery(string) []string
	Name() string
}
