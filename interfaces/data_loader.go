package interfaces

type DataLoader interface {
	LoadData(string) []string
	Name() string
}
