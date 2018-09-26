package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	in := []byte("table: TEST\nlimit: 0\noffset: 0\ndbtype: sqlite\npath: /tmp/foo")
	//var c Config_yml
	c := new(Config_yml)
	var want ConfigT
	var got *ConfigT
	want.path = "/tmp/foo"
	want.table = "TEST"
	want.limit = 0
	want.offset = 0
	want.dbtype = "sqlite"
	got, err := c.ReadConfig(in)
	if err != nil {
		t.Fatal(err)
	}
	if want.table != got.table {
		t.Errorf("Wrong table, want: %s, got: %s", want.table, got.table)
	}
	if want.path != got.path {
		t.Errorf("Wrong path, want: %s, got: %s", want.path, got.path)
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
