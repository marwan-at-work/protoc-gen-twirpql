package genenums

import (
	"bytes"
	"go/format"
	"io"
)

// Data is the Data that's needed
// to bridge protocol buffer enums
// and gql enums.
type Data struct {
	ImportPath string
	Pkg        string
	Name       string
	GoName     string
}

type final struct {
	Imports []string
	Enums   []*Data
}

// Render extends gqlgen's exectuionContext
// to map protobuf enums to gql enums
func Render(data []*Data, out io.Writer) error {
	var b bytes.Buffer
	final := &final{}
	mp := map[string]struct{}{}
	for _, d := range data {
		mp[d.ImportPath] = struct{}{}
		final.Enums = append(final.Enums, d)
	}
	for k := range mp {
		final.Imports = append(final.Imports, k)
	}
	err := enumTemplate.Execute(&b, final)
	if err != nil {
		return err
	}
	bts, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	_, err = io.Copy(out, bytes.NewReader(bts))
	return err
}
