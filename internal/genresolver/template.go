package genresolver

var tmpl = `
{{ reserveImport "context"  }}
{{ reserveImport "fmt"  }}
{{ reserveImport "io"  }}
{{ reserveImport "strconv"  }}
{{ reserveImport "time"  }}
{{ reserveImport "sync"  }}
{{ reserveImport "errors"  }}
{{ reserveImport "bytes"  }}

{{ reserveImport "github.com/99designs/gqlgen/handler" }}
{{ reserveImport "github.com/vektah/gqlparser" }}
{{ reserveImport "github.com/vektah/gqlparser/ast" }}
{{ reserveImport "github.com/99designs/gqlgen/graphql" }}
{{ reserveImport "github.com/99designs/gqlgen/graphql/introspection" }}
{{ $serviceName := .ServiceName }}
{{ $servicePackageName := .ServicePackageName }}
type {{.ResolverType}} struct {
    {{$servicePackageName}}.{{$serviceName}}
}

{{ range $object := .Objects -}}
	{{- if $object.HasResolvers -}}
		func (r *{{$.ResolverType}}) {{$object.Name}}() {{ $object.ResolverInterface | ref }} {
			return &{{lcFirst $object.Name}}Resolver{r}
		}
	{{ end -}}
{{ end }}

{{ range $object := .Objects -}}
	{{- if $object.HasResolvers -}}
		type {{lcFirst $object.Name}}Resolver struct { *Resolver }

		{{ range $field := $object.Fields -}}
			{{- if $field.IsResolver -}}
			func (r *{{lcFirst $object.Name}}Resolver) {{$field.GoFieldName}}{{ $field.ShortResolverDeclaration }} {
				{{- $reqArg := "req" -}}
				{{- if (hasPrefix ($field.ShortResolverDeclaration) "(ctx context.Context)") -}}
					{{ $reqArg = "nil" }}
				{{ end -}}
				{{- if (isEmpty $field) -}}
				_, err := r.{{$serviceName}}.{{$field.GoFieldName}}(ctx, {{$reqArg}})
				if err != nil {
					return nil, err
				}
				return {{getType $field}}{}, nil
				{{ else if (isScalar ($field.GoFieldName)) }}
					return obj.Get{{$field.GoFieldName}}(), nil
				{{ else if (isUnion ($field.GoFieldName)) }}
					return obj.Get{{$field.GoFieldName}}(), nil
				{{ else if (isResponseUnion ($field.GoFieldName)) }}
				resp, err := r.{{$serviceName}}.{{$field.GoFieldName}}(ctx, {{$reqArg}})
				if err != nil {
					{{ $errorTypeName := (responseUnionName ($field.GoFieldName)) }}
					if errval, ok := err.(interface {
						{{$errorTypeName}}() *{{$servicePackageName}}.{{$errorTypeName}}
					}); ok {
						newresp := errval.{{$errorTypeName}}()
						if newresp != nil {
							return newresp, nil
						}
					}
				}
				return resp, err
				{{- else -}}
				return r.{{$serviceName}}.{{$field.GoFieldName}}(ctx, {{$reqArg}})
				{{ end -}}
			}
			{{ end }}
		{{ end -}}
	{{ end -}}
{{ end }}
`
