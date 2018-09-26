package main

type Config interface {
	ReadConfig([]byte) (*ConfigT, error)
	MakeConfig() []byte
}

type ConfigT struct {
	table  string
	limit  int
	offset int
	dbtype string
	path   string
}
