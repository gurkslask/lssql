package lssql

import (
	"gopkg.in/yaml.v2"
)

type Config_yml struct {
	Table  string `yaml:"table"`
	Limit  int    `yaml:"limit"`
	Offset int    `yaml:"offset"`
	Dbtype string `yaml:"dbtype"`
	Path   string `yaml:"path"`
}

func (c Config_yml) ReadConfig(in []byte) (*ConfigT, error) {
	err := yaml.Unmarshal(in, &c)
	if err != nil {
		return nil, err
	}

	cc := new(ConfigT)
	cc.Table = c.Table
	cc.Limit = c.Limit
	cc.Offset = c.Offset
	cc.Dbtype = c.Dbtype
	cc.Path = c.Path
	return cc, nil
}

func (c Config_yml) MakeConfig() []byte {
	return []byte("table: TABLENAME\nlimit: 0\noffset: 0\ndbtype: sqlite\npath: /tmp/foo")
}
