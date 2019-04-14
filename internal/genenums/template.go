package genenums

import "text/template"

var enumTemplate = template.Must(template.New("").Parse(`package twirpql

import (
	"context"
	"errors"

	"{{.ImportPath}}"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/ast"
)
{{- $pkg := .Pkg -}}
{{ range .Enums }}
func (ec *executionContext) _{{ . }}(ctx context.Context, sel ast.SelectionSet, v *{{$pkg}}.{{.}}) graphql.Marshaler {
	return graphql.MarshalString((*v).String())
}

func (ec *executionContext) unmarshalInput{{.}}(ctx context.Context, v interface{}) ({{$pkg}}.{{.}}, error) {
	switch v := v.(type) {
	case string:
		intValue, ok := {{$pkg}}.{{.}}_value[v]
		if !ok {
			return 0, errors.New("unknown value: " + v)
		}
		return {{$pkg}}.{{.}}(intValue), nil
	}
	return 0, errors.New("wrong type")
}
{{ end }}
`))
