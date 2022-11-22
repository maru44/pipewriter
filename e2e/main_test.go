package main_test

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	data, err := os.ReadFile("rwbytest")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "RubyWeissBreakYangPenny", string(data))

	f, err := os.Open("rwbytest.csv")
	if err != nil {
		t.Fatal(err)
	}
	r := csv.NewReader(f)
	wants := [][]string{
		{"\ufeffname", "age", "color"},
		{"Ruby", "15", "red"},
		{"Weiss", "17", "white"},
		{"Break", "17", "black"},
		{"Yang", "17", "yellow"},
		{"Penny", "0", "light green"},
	}
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wants, rows)
}
