package genenums

import (
	"bytes"
	"go/format"
	"io"
)

// Data is the data that's needed
// to bridge protocol buffer enums
// and gql enums.
type Data struct {
	ImportPath string
	Pkg        string
	Enums      []string
}

// Render extends gqlgen's exectuionContext
// to map protobuf enums to gql enums
func Render(data *Data, out io.Writer) error {
	var b bytes.Buffer
	err := enumTemplate.Execute(&b, data)
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
