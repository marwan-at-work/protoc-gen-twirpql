package genresolver

import (
	"strings"
	"text/template"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

func New(
	serviceName,
	pkgName string,
	emptys []string,
	scalars map[string]string,
	unions map[string]bool,
	responseUnions map[string]string,
) plugin.Plugin {
	return &Plugin{
		ServiceName:    serviceName,
		PackageName:    pkgName,
		Emptys:         emptys,
		Scalars:        scalars,
		Unions:         unions,
		ResponseUnions: responseUnions,
	}
}

type Plugin struct {
	ServiceName    string
	PackageName    string
	Emptys         []string
	Scalars        map[string]string
	Unions         map[string]bool
	ResponseUnions map[string]string
}

func (m *Plugin) isEmpty(f *codegen.Field) bool {
	name := templates.CurrentImports.LookupType(f.TypeReference.GO)
	for _, e := range m.Emptys {
		if name == "*"+e {
			return true
		}
	}
	return false
}

var _ plugin.CodeGenerator = &Plugin{}

func (m *Plugin) Name() string {
	return "resovleroverride"
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	if !data.Config.Resolver.IsDefined() {
		return nil
	}

	resolverBuild := &ResolverBuild{
		Data:               data,
		PackageName:        data.Config.Resolver.Package,
		ResolverType:       data.Config.Resolver.Type,
		ServiceName:        m.ServiceName,
		ServicePackageName: m.PackageName,
	}

	return templates.Render(templates.Options{
		GeneratedHeader: true,
		Template:        tmpl,
		PackageName:     data.Config.Resolver.Package,
		Filename:        data.Config.Resolver.Filename,
		Data:            resolverBuild,
		Funcs: template.FuncMap{
			"hasPrefix": hasPrefix,
			"isEmpty":   m.isEmpty,
			"getType": func(f *codegen.Field) string {
				return strings.Replace(templates.CurrentImports.LookupType(f.TypeReference.GO), "*", "&", 1)
			},
			"isScalar": func(s string) bool {
				_, ok := m.Scalars[s]
				return ok
			},
			"isUnion": func(s string) bool {
				_, ok := m.Unions[s]
				return ok
			},
			"isResponseUnion": func(s string) bool {
				_, ok := m.ResponseUnions[s]
				return ok
			},
			"responseUnionName": func(s string) string {
				return m.ResponseUnions[s]
			},
		},
	})
}

type ResolverBuild struct {
	*codegen.Data

	PackageName        string
	ResolverType       string
	ServiceName        string
	ServicePackageName string
}

func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
