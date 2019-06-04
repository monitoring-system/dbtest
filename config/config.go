package config

type Config struct {
	Loop         int
	LoopInterval int

	DataLoaders  string
	QueryLoaders string
	Comparor     string
	CellFilter   string

	StandardDB string
	TestDB     string
}
