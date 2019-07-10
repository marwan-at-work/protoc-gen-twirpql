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
    """
    {{.Doc}}
    """
    {{.Name}}{{.Request}}: {{.Response}}!{{end}}
}

{{ end }}

{{ if (gt (len .Service.Mutations) 0) }}

type Mutation { {{range .Service.Mutations}}
    {{.Name}}{{.Request}}: {{.Response}}!{{end}}
}

{{ end }}

{{range .Types}}
"""
{{ .Doc }}
"""
type {{.Name}} { {{- range .Fields}}
    """
    {{ .Doc }}
    """
    {{.Name}}: {{.Type}}!{{end}}
    {{- if (eq (len .Fields) 0) }}
    responseMessage: String!
    {{- end }}
}
{{end}}
{{range .Inputs}}
"""
{{ .Doc }}
"""
input {{.Name}} { {{range .Fields}}
    """
    {{ .Doc }}
    """
    {{.Name}}: {{.Type}}!{{end}}
}
{{end}}
{{range .Enums}}
"""
{{ .Doc }}
"""
enum {{.Name}} { {{range .Fields}}
    """
    {{ .Doc }}
    """
    {{.Name}}{{end}}
}{{end}}
{{range .Scalars}}
scalar {{.}}
{{end}}
{{range .Unions}}
union {{.Name}} = {{fmtUnions .Types }}
{{ end }}
`
