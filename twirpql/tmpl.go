package twirpql

import (
	"strings"
	"text/template"
)

var schemaFuncs = template.FuncMap{
	"fmtUnions": func(types []string) string {
		return strings.Join(types, " | ")
	},
}

const schemaTemplate = `
{{ if (gt (len .Service.Methods) 0) }}

type Query { {{range .Service.Methods}}
    {{.Name}}{{.Request}}: {{.Response}}!{{end}}
}

{{ end }}

{{ if (gt (len .Service.Mutations) 0) }}

type Mutation { {{range .Service.Mutations}}
    {{.Name}}{{.Request}}: {{.Response}}!{{end}}
}

{{ end }}

{{range .Types}}
type {{.Name}} { {{- range .Fields}}
    {{.Name}}: {{.Type}}!{{end}}
    {{- if (eq (len .Fields) 0) }}
    responseMessage: String!
    {{- end }}
}
{{end}}
{{range .Inputs}}
input {{.Name}} { {{range .Fields}}
    {{.Name}}: {{.Type}}!{{end}}
}
{{end}}
{{range .Enums}}
enum {{.Name}} { {{range .Fields}}
    {{.}}{{end}}
}{{end}}
{{range .Scalars}}
scalar {{.}}
{{end}}
{{range .Unions}}
union {{.Name}} = {{fmtUnions .Types }}
{{ end }}
`
