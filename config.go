package lssql

type Config interface {
	ReadConfig([]byte) (*ConfigT, error)
	MakeConfig() []byte
}

type ConfigT struct {
	Table  string
	Limit  int
	Offset int
	Dbtype string
	Path   string
}
