package genscalar

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "rewrite all golden files")

func TestGenScalar(t *testing.T) {
	d := map[string]string{
		"MyMap":  "map[string]string",
		"MyInts": "map[int64]string",
	}

	var b bytes.Buffer
	err := Render(d, nil, &b)
	require.NoError(t, err)

	if *update {
		ioutil.WriteFile("testdata/scalars.golden", b.Bytes(), 0660)
		return
	}

	expected, err := ioutil.ReadFile("testdata/scalars.golden")
	require.NoError(t, err)
	require.Equal(t, string(expected), b.String())
}
