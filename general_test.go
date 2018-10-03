package lssql

import (
	//"strings"
	//"fmt"
	"testing"
)

func TestPadString(t *testing.T) {
	data := "testdata"
	length := 0
	dest2 := ""
	dest := &dest2

	PadString(data, length, dest)
	if *dest != "testdata\t" {
		t.Fail()
	}
	PadString(data, length, dest)
	if *dest != "testdata\ttestdata\t" {
		t.Fail()
	}
	*dest = ""
	PadString(data, 19, dest)
	want := "testdata           \t"
	if *dest != want {
		t.Errorf("Got :%s, want:%s,", *dest, want)
	}
}

func TestMaxColumnLength(t *testing.T) {
	data := [][]string{
		{"hej", "test", "dnsajkldsjakdnsöad"},
		{"hejdsadsa", "test", "dnsajkldsjakdnsöad"},
		{"hejdsadsa", "thisshouldbe19", "dnsajkldsjakdnsöad"},
	}
	want := []int{14, 19, 24}
	result := MaxColumnLength(data)
	for i, _ := range result {
		if want[i] != result[i] {
			t.Errorf("Got :%d, want:%d,", result, want)
		}
	}
}
