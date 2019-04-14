package gqlfmt

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update expected file to the given results")

func TestPrint(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/given.graphql")
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	err = Print(string(input), &b)
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		ioutil.WriteFile("testdata/expected.graphql", b.Bytes(), 0660)
		return
	}

	expected, err := ioutil.ReadFile("testdata/expected.graphql")
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, string(expected), b.String())
}
