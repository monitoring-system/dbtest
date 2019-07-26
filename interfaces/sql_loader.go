package interfaces

type SqlLoader interface {
	LoadSql(string) []string
	Name() string
}