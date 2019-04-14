package genenums

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "rewrite all golden files")

func TestGenEnums(t *testing.T) {
	d := &Data{
		ImportPath: "pkg.go/enums",
		Pkg:        "enums",
		Enums:      []string{"one"},
	}

	var b bytes.Buffer
	err := Render(d, &b)
	require.NoError(t, err)

	if *update {
		ioutil.WriteFile("testdata/enums.golden", b.Bytes(), 0660)
		return
	}

	expected, err := ioutil.ReadFile("testdata/enums.golden")
	require.NoError(t, err)
	require.Equal(t, string(expected), b.String())
}
