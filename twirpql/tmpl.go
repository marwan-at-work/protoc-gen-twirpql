package twirpql

const schemaTemplate = `schema {
    query: Query
}

type Query { {{range .Service.Methods}}
    {{.Name}}{{.Request}}: {{.Response}}!{{end}}
}
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
`
