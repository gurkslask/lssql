package lssql

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	in := []byte("table: TEST\nlimit: 0\noffset: 0\ndbtype: sqlite\npath: /tmp/foo")
	//var c Config_yml
	c := new(Config_yml)
	var want ConfigT
	var got *ConfigT
	want.Path = "/tmp/foo"
	want.Table = "TEST"
	want.Limit = 0
	want.Offset = 0
	want.Dbtype = "sqlite"
	got, err := c.ReadConfig(in)
	if err != nil {
		t.Fatal(err)
	}
	if want.Table != got.Table {
		t.Errorf("Wrong table, want: %s, got: %s", want.Table, got.Table)
	}
	if want.Path != got.Path {
		t.Errorf("Wrong path, want: %s, got: %s", want.Path, got.Path)
	}
}

func TestMakeConfig(t *testing.T) {
	want := []byte("table: TABLENAME\nlimit: 0\noffset: 0\ndbtype: sqlite\npath: /tmp/foo")
	var c Config_yml

	got := c.MakeConfig()
	if string(got) != string(want) {
		t.Errorf("Got: %s, Want: %s", string(got), string(want))
	}
}
