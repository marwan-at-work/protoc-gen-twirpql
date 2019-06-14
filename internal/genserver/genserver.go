package genserver

import (
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

func New(filename, modPath, serviceName string) plugin.Plugin {
	return &Plugin{filename, modPath, serviceName}
}

type Plugin struct {
	filename    string
	modPath     string
	serviceName string
}

var _ plugin.CodeGenerator = &Plugin{}

func (m *Plugin) Name() string {
	return "servergen"
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	serverBuild := &ServerBuild{
		Data:                data,
		ExecPackageName:     data.Config.Exec.ImportPath(),
		ResolverPackageName: data.Config.Resolver.ImportPath(),
		ModPath:             m.modPath,
		ServiceName:         m.serviceName,
	}

	return templates.Render(templates.Options{
		GeneratedHeader: true,
		Template:        tmpl,
		PackageName:     "twirpql", // TODO: dynamic package name
		Filename:        m.filename,
		Data:            serverBuild,
	})
}

type ServerBuild struct {
	*codegen.Data

	ExecPackageName     string
	ResolverPackageName string
	ModPath             string
	ServiceName         string
}
