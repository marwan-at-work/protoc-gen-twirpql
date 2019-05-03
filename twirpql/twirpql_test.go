package twirpql

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "rewrite all golden files")

func TestSchema(t *testing.T) {
	dirs, err := ioutil.ReadDir("testdata")
	require.NoError(t, err)
	for _, dir := range dirs {
		t.Run(dir.Name(), func(t *testing.T) {
			m, f := getModule(t, dir.Name())
			var bts bytes.Buffer
			m.generateSchema(f, &bts)
			if *update {
				writeGoldenSchema(t, bts.Bytes(), dir.Name())
				return
			}
			given := readGoldenSchema(t, dir.Name())
			require.Equal(t, given, bts.String())
		})
	}
}

func TestGQLConfig(t *testing.T) {
	dirs, err := ioutil.ReadDir("testdata")
	require.NoError(t, err)
	for _, dir := range dirs {
		t.Run(dir.Name(), func(t *testing.T) {
			m, f := getModule(t, dir.Name())
			m.generateSchema(f, ioutil.Discard)
			var bts bytes.Buffer
			m.touchConfig(&bts)
			if *update {
				writeGoldenConfig(t, bts.Bytes(), dir.Name())
				return
			}
			given := readGoldenConfig(t, dir.Name())
			require.Equal(t, given, bts.String())
		})
	}
}

func getModule(t *testing.T, dirName string) (*twirpql, pgs.File) {
	t.Helper()
	ast := buildGraph(t, dirName)
	f := ast.Targets()[dirName+".proto"]
	ctx := pgsgo.InitContext(pgs.ParseParameters(""))
	m := New(dirName).(*twirpql)
	m.ctx = ctx
	m.svc = f.Services()[0]
	m.protopkg = f.Package()
	return m, f
}

func writeGoldenConfig(t *testing.T, bts []byte, dir ...string) {
	t.Helper()
	dirs := append(append([]string{"testdata"}, dir...), "gqlgen.yml.golden")
	filename := filepath.Join(dirs...)
	err := ioutil.WriteFile(filename, bts, 0660)
	require.NoError(t, err)
}

func readGoldenConfig(t *testing.T, dir ...string) string {
	t.Helper()
	dirs := append(append([]string{"testdata"}, dir...), "gqlgen.yml.golden")
	filename := filepath.Join(dirs...)

	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err, "unable to read CDR at %q", filename)

	return string(data)
}

func writeGoldenSchema(t *testing.T, bts []byte, dir ...string) {
	t.Helper()
	dirs := append(append([]string{"testdata"}, dir...), "schema.graphql.golden")
	filename := filepath.Join(dirs...)
	err := ioutil.WriteFile(filename, bts, 0660)
	require.NoError(t, err)
}

func readGoldenSchema(t *testing.T, dir ...string) string {
	t.Helper()
	dirs := append(append([]string{"testdata"}, dir...), "schema.graphql.golden")
	filename := filepath.Join(dirs...)

	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err, "unable to read CDR at %q", filename)

	return string(data)
}

func readCodeGenReq(t *testing.T, dir ...string) *plugin_go.CodeGeneratorRequest {
	t.Helper()
	dirs := append(append([]string{"testdata"}, dir...), "code_generator_request.pb.bin")
	filename := filepath.Join(dirs...)

	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err, "unable to read CDR at %q", filename)

	req := &plugin_go.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, req)
	require.NoError(t, err, "unable to unmarshal CDR data at %q", filename)

	return req
}

func buildGraph(t *testing.T, dir ...string) pgs.AST {
	t.Helper()
	d := pgs.InitMockDebugger()
	ast := pgs.ProcessCodeGeneratorRequest(d, readCodeGenReq(t, dir...))
	require.False(t, d.Failed(), "failed to build graph (see previous log statements)")
	return ast
}
