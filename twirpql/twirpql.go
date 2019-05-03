package twirpql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/99designs/gqlgen/api"
	gqlconfig "github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"gopkg.in/yaml.v2"
	"marwan.io/protoc-gen-twirpql/internal/genenums"
	"marwan.io/protoc-gen-twirpql/internal/genresolver"
	"marwan.io/protoc-gen-twirpql/internal/genscalar"
	"marwan.io/protoc-gen-twirpql/internal/genserver"
	"marwan.io/protoc-gen-twirpql/internal/gqlfmt"
)

// twirpql creates a report of all the target messages generated by the
// protoc run, writing the file into the /tmp directory.
type twirpql struct {
	*pgs.ModuleBase
	// an input is a protobuf "message" that is
	// found inside an RPC's Request so that GraphQL
	// interprets it as an Input declaration.
	// Note that if the same input is also found
	// in an rpc's "returns" value, then the name
	// will be suffixed with the word "Input"
	// because GraphQL does not allow types and
	// inputs with matching names.
	inputs map[string]*serviceType
	// a "type" is a protobuf "message" that is
	// found inside an RPC's Return so that GraphQL
	// interprets it as a "Type" declaration.
	types map[string]*serviceType
	// an empty type keeps track of empty returns
	// because GraphQL Types can't be empty
	// and therefore we need to inject a dummy
	// field.
	emptys map[string]bool
	// Enums are integers in protobuf but strings
	// in GraphQL. Therefore, we need to keep track
	// of declared enums in the proto file so that
	// we create proper conversion for the GraphQL queries.
	enums map[string][]string
	maps  map[string]string
	// gqlTypes are specific for the gqlgen config file
	// so that we make all the input/output GraphQL
	// types point to the generated .pb.go types.
	gqlTypes  gqlconfig.TypeMap
	tmpl      *template.Template
	ctx       pgsgo.Context
	modname   string
	gopkgname string
	svcname   string
	// destpkgname is the directory path
	// where the GraphQL generated code will
	// live. It defaults to a "twirpql".
	destpkgname string
	svc         pgs.Service
	protopkg    pgs.Package
}

// New configures the module with an instance of ModuleBase
func New(importPath string) pgs.Module {
	return &twirpql{
		ModuleBase:  &pgs.ModuleBase{},
		inputs:      map[string]*serviceType{},
		types:       map[string]*serviceType{},
		emptys:      map[string]bool{},
		enums:       map[string][]string{},
		maps:        map[string]string{},
		gqlTypes:    gqlconfig.TypeMap{},
		tmpl:        template.Must(template.New("").Parse(schemaTemplate)),
		modname:     importPath,
		ctx:         pgsgo.InitContext(pgs.ParseParameters("")),
		destpkgname: "./twirpql",
	}
}

// Name is the identifier used to identify the module. This value is
// automatically attached to the BuildContext associated with the ModuleBase.
func (tql *twirpql) Name() string { return "twirpql" }

func (tql *twirpql) InitContext(c pgs.BuildContext) {
	tql.ModuleBase.InitContext(c)
	tql.ctx = pgsgo.InitContext(c.Parameters())
}

// Execute is passed the target files as well as its dependencies in the pkgs
// map. The implementation should return a slice of Artifacts that represent
// the files to be generated. In this case, "/tmp/report.txt" will be created
// outside of the normal protoc flow.
func (tql *twirpql) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	tql.destpkgname = tql.Parameters().StrDefault("dest", tql.destpkgname)
	os.MkdirAll(tql.destpkgname, 0777)

	if len(targets) != 1 {
		panic("only one proto file is supported at this moment; see https://twirpql.dev/docs/multiple-services")
	}

	for fileName, targetFile := range targets {
		tql.svc = tql.pickServiceFromFile(tql.Parameters().Str("service"), targetFile)
		tql.protopkg = targetFile.Package()
		serviceDir := filepath.Dir(fileName)
		tql.setImportPath(serviceDir)
		f, err := os.Create(tql.schemaPath())
		must(err)
		defer f.Close()
		tql.generateSchema(targetFile, f)
	}

	if len(tql.maps) > 0 {
		f, err := os.Create(filepath.Join(tql.destpkgname, "scalars.go"))
		must(err)
		genscalar.Render(tql.maps, f)
		f.Close()
	}

	f, err := os.Create(filepath.Join(tql.destpkgname, "gqlgen.yml"))
	must(err)
	defer f.Close()
	cfg := tql.touchConfig(f)
	if len(tql.enums) > 0 {
		tql.bridgeEnums()
	}
	tql.initGql(cfg, tql.svcname)

	return tql.Artifacts()
}

func (tql *twirpql) pickServiceFromFile(svc string, f pgs.File) pgs.Service {
	switch len(f.Services()) {
	case 0:
		panic("proto file must have at least one service")
	case 1:
		return f.Services()[0]
	}
	if svc == "" {
		panic("service name must be provided if proto file has multiple services; see https://twirpql.dev/docs/multiple-services")
	}
	for _, service := range f.Services() {
		if svc == service.Name().String() {
			return service
		}
	}
	panic("protofile does not have the given service: " + svc)
}

func (tql *twirpql) setImportPath(serviceDir string) {
	cmd := exec.Command("go", "list")
	cmd.Dir = serviceDir
	cmd.Env = os.Environ()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	pkgpath, err := cmd.Output()
	if err != nil {
		msg := fmt.Sprintf("go list failed: %v - stdout: %v - stderr: %v", err, string(pkgpath), stderr.String())
		if strings.Contains(stderr.String(), "cannot find module providing package") {
			msg = "go list failed. Make sure you have .go files where your .proto file is." +
				"Also make sure to run the --go_out=. --twirp_out=. plugins on a separate command before you run --twirpql_out"
		}
		panic(msg)
	}
	modname := strings.TrimSpace(string(pkgpath))
	tql.modname = tql.Parameters().StrDefault("importpath", modname)
	if tql.modname == "" {
		panic("import path must be provided by `go list` in the .proto directory or through the importpath plugin parameter")
	}
}

func (tql *twirpql) generateSchema(f pgs.File, out io.Writer) {
	tql.svcname = tql.svc.Name().String()
	tql.gopkgname = tql.ctx.PackageName(f).String()
	gqlFile := &file{}
	gqlFile.Service = tql.getService(tql.svc)
	for _, v := range tql.inputs {
		gqlFile.Inputs = append(gqlFile.Inputs, v)
	}
	for _, v := range tql.types {
		gqlFile.Types = append(gqlFile.Types, v)
	}
	for k, v := range tql.enums {
		gqlFile.Enums = append(gqlFile.Enums, &enums{Name: k, Fields: v})
	}
	for k := range tql.maps {
		gqlFile.Scalars = append(gqlFile.Scalars, k)
	}

	var buf bytes.Buffer

	err := tql.tmpl.Execute(&buf, gqlFile)
	must(err)
	err = gqlfmt.Print(buf.String(), out)
	must(err)
}

// bridgeEnums creates a type conversion between
// protobuf's enums (int32) and GraphQL's enums (string).
func (tql *twirpql) bridgeEnums() {
	f, err := os.Create(filepath.Join(tql.destpkgname, "enums.gen.go"))
	must(err)
	defer f.Close()
	ed := &genenums.Data{
		ImportPath: tql.modname,
		Pkg:        tql.gopkgname,
		Enums:      []string{},
	}
	for k := range tql.enums {
		ed.Enums = append(ed.Enums, k)
	}
	must(genenums.Render(ed, f))
}

func (tql *twirpql) touchConfig(out io.Writer) *gqlconfig.Config {
	var cfg gqlconfig.Config
	cfg.SchemaFilename = gqlconfig.StringList{tql.schemaPath()}
	cfg.Exec = gqlconfig.PackageConfig{Filename: filepath.Join(tql.destpkgname, "generated.go")}
	cfg.Resolver = gqlconfig.PackageConfig{Filename: filepath.Join(tql.destpkgname, "resolver.go"), Type: "Resolver"}
	cfg.Models = tql.gqlTypes
	cfg.Model = gqlconfig.PackageConfig{Filename: filepath.Join(tql.destpkgname, "models_gen.go")}
	must(yaml.NewEncoder(out).Encode(&cfg))
	return &cfg
}

func (tql *twirpql) initGql(cfg2 *gqlconfig.Config, svcName string) {
	cfg, err := gqlconfig.LoadConfig(filepath.Join(tql.destpkgname, "gqlgen.yml"))
	must(err)
	emptys := []string{}
	for k := range tql.emptys {
		emptys = append(emptys, k)
	}

	err = api.Generate(
		cfg,
		api.NoPlugins(),
		api.AddPlugin(modelgen.New()),
		api.AddPlugin(genresolver.New(svcName, tql.gopkgname, emptys, tql.maps)),
		api.AddPlugin(genserver.New(filepath.Join(tql.destpkgname, "server.go"), tql.modname, svcName)),
	)
	must(err)
}

func (tql *twirpql) getService(svc pgs.Service) *service {
	var s service
	s.Methods = tql.getMethods(svc.Methods())
	return &s
}

func (tql *twirpql) getMethods(protoMethods []pgs.Method) []*method {
	methods := []*method{}

	// collect all types first, so that we de-dupe mixed
	// inputs && types
	for _, pm := range protoMethods {
		tql.setType(pm.Output())
	}

	for _, pm := range protoMethods {
		var m method
		m.Name = pm.Name().String()
		emptyInput := len(pm.Input().Fields()) == 0
		if !emptyInput {
			tql.setInput(pm.Input())
			m.Request = tql.formatQueryInput(pm.Input())
		}
		m.Response = tql.getQualifiedName(pm.Output())
		methods = append(methods, &m)
	}
	return methods
}

func (tql *twirpql) setType(msg pgs.Message) {
	typeName := tql.getQualifiedName(msg)
	if _, ok := tql.types[typeName]; ok {
		return
	}
	var i serviceType
	i.Name = typeName
	tql.types[i.Name] = &i
	tql.setGraphQLType(i.Name, msg)
	i.Fields = tql.getFields(msg.Fields(), true)
}

// getQualifiedName returns the name that will be defined inside the GraphQL Schema File.
// For messgae declarations that are part of the target .proto file, they will stay the same
// but if it's part of an import like "google.protobuf.Timestamp" then we combine the package name
// with the Message namd to ensure we have no clashes so it becomes: "google_protobuf_Timestamp"
func (tql *twirpql) getQualifiedName(msg pgs.Message) string {
	if msg.Package() == tql.protopkg {
		return msg.Name().String()
	}
	pkgName := strings.ReplaceAll(msg.Package().ProtoName().String(), ".", "_")
	return pkgName + "_" + msg.Name().String()
}

func (tql *twirpql) setInput(msg pgs.Message) {
	if _, ok := tql.inputs[tql.getInputName(msg)]; ok {
		return
	}
	var i serviceType
	i.Name = tql.getInputName(msg)
	tql.inputs[i.Name] = &i
	tql.setGraphQLType(i.Name, msg)
	i.Fields = tql.getFields(msg.Fields(), false)
}

// getInputName returns exactly the name of the message declaration:
// message SomeMessage {
//   ... fields
// }
// would return SomeMessage. However, if SomeMessage was also
// used as an Output and not just Input, then GraphQL will
// not allow an Input and a Type to be the same name, therefore
// we will append an "Input" so that it becomes SomeMessageInput.
func (tql *twirpql) getInputName(msg pgs.Message) string {
	msgName := tql.getQualifiedName(msg)
	if _, ok := tql.types[msgName]; ok {
		return msgName + "Input"
	}
	return msgName
}

func (tql *twirpql) setGraphQLType(name string, msg pgs.Message) {
	if len(msg.Fields()) == 0 {
		tql.emptys[name] = true
		return
	}
	importpath := tql.ctx.ImportPath(msg.File()).String()
	if importpath == "." {
		importpath = tql.modname
	}
	tql.gqlTypes[name] = gqlconfig.TypeMapEntry{
		Model: gqlconfig.StringList{importpath + "." + msg.Name().String()},
	}
}

func (tql *twirpql) setEnum(protoEnum pgs.Enum) {
	vals := []string{}
	for _, v := range protoEnum.Values() {
		vals = append(vals, v.Name().String())
	}
	tql.enums[protoEnum.Name().String()] = vals
	tql.setGraphQLEnum(protoEnum.Name().String())
}

func (tql *twirpql) setGraphQLEnum(name string) {
	tql.gqlTypes[name] = gqlconfig.TypeMapEntry{
		Model: gqlconfig.StringList{tql.modname + "." + name},
	}
}

func (tql *twirpql) setMap(fieldName string, f pgs.Field) {
	upField := strings.Title(fieldName)
	tql.maps[upField] = tql.ctx.Type(f).Value().String()
	tql.gqlTypes[upField] = gqlconfig.TypeMapEntry{
		Model: gqlconfig.StringList{tql.modname + "/twirpql." + upField},
	}
}

func (tql *twirpql) getFields(protoFields []pgs.Field, isType bool) []*serviceField {
	fields := []*serviceField{}
	for _, pf := range protoFields {
		var f serviceField
		f.Name = pf.Name().String()
		pt := pf.Type().ProtoType().Proto()
		var tmp string
		switch pt {
		case 11:
			if pf.Type().IsMap() {
				tql.setMap(f.Name, pf)
				tmp = strings.Title(f.Name)
			} else {
				if isType {
					tmp = tql.getQualifiedName(pf.Type().Embed())
					tql.setType(pf.Type().Embed())
				} else {
					tmp = tql.getInputName(pf.Type().Embed())
					tql.setInput(pf.Type().Embed())
				}
			}
		case 14:
			tql.setEnum(pf.Type().Enum())
			tmp = tql.ctx.Type(pf).Value().String()
		default:
			tmp = protoTypesToGqlTypes[pt.String()]
			if tmp == "" {
				panic("unsupported type: " + pt.String())
			}
		}
		if pf.Type().IsRepeated() {
			tmp = fmt.Sprintf("[%v]", tmp)
		}
		f.Type = tmp
		fields = append(fields, &f)
	}
	return fields
}

// formatQueryInput returns a template-formatted representation
// of a query input. In GraphQL a query looks like this:
// `someQuery(req: Request): Response`
// However, if we don't want to have an input at all in a query,
// the query will now have to look like this:
// `someQuery: Response`
func (tql *twirpql) formatQueryInput(msg pgs.Message) string {
	return fmt.Sprintf("(req: %v)", tql.getInputName(msg))
}

func (tql *twirpql) schemaPath() string {
	return filepath.Join(tql.destpkgname, "schema.graphql")
}

var protoTypesToGqlTypes = map[string]string{
	"TYPE_DOUBLE":  "Float",
	"TYPE_FLOAT":   "Float",
	"TYPE_INT64":   "Int",
	"TYPE_UINT64":  "Int",
	"TYPE_INT32":   "Int",
	"TYPE_FIXED64": "Float",
	"TYPE_FIXED32": "Float",
	"TYPE_BOOL":    "Boolean",
	"TYPE_STRING":  "String",
	// "TYPE_GROUP": "",
	// "TYPE_MESSAGE": "", // must be mapped to its sibling type
	"TYPE_BYTES":  "String",
	"TYPE_UINT32": "Int",
	// "TYPE_ENUM": "", // mapped to its sibling type
	// "TYPE_SFIXED32": "",
	// "TYPE_SFIXED64": "",
	// "TYPE_SINT32": "",
	// "TYPE_SINT64": "",
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
