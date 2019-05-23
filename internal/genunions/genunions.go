package genunions

import (
	"bytes"
	"go/format"
	"io"
	"text/template"
)

var tmpl = template.Must(template.New("genunions").Parse(tmplStr))

// Render renders a union
func Render(w io.Writer) error {
	var bts bytes.Buffer
	err := tmpl.Execute(&bts, nil)
	if err != nil {
		return err
	}
	formatted, err := format.Source(bts.Bytes())
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewReader(formatted))
	return err
}

const tmplStr = `package twirpql

type unionMask interface {}
`
