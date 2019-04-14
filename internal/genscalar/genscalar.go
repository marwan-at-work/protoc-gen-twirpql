package genscalar

import (
	"bytes"
	"go/format"
	"io"
	"text/template"
)

var tmpl = template.Must(template.New("").Parse(templateText))

// Render renders a scalar
// implementation.
func Render(mp map[string]string, out io.Writer) error {
	var b bytes.Buffer
	err := tmpl.Execute(&b, mp)
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

const templateText = `package twirpql

import (
	"encoding/json"
	"io"
)

{{range $key, $val := .}}
type {{$key}} {{$val}}

func (scalar *{{$key}}) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return nil
	}
	return json.Unmarshal([]byte(str), scalar)
}

func (scalar {{$key}}) MarshalGQL(w io.Writer) {
	json.NewEncoder(w).Encode(scalar)
}
{{end}}`
