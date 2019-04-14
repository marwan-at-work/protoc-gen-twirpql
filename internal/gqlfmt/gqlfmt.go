package gqlfmt

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

// Print parses the input as a graphql schema
// and prints to the given io.Writer.
// TODO: preserve Description
func Print(input string, out io.Writer) error {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: input,
	})
	if err != nil {
		return err
	}
	f := &formatter{schema: schema, out: out}
	f.print()

	return nil
}

type formatter struct {
	schema  *ast.Schema
	out     io.Writer
	types   []string
	inputs  []string
	scalars []string
	enums   []string
}

func (f *formatter) print() {
	f.sortDeclarations()
	f.printHeader()
	f.printQuery()
	f.printTypes()
	f.printInputs()
	f.printEnums()
	f.printScalars()
}

func (f *formatter) sortDeclarations() {
	for k, def := range f.schema.Types {
		// TODO: maybe include query and refactor above.
		if k == "Query" || def.BuiltIn {
			continue
		}
		switch def.Kind {
		case ast.Object:
			f.types = append(f.types, k)
		case ast.InputObject:
			f.inputs = append(f.inputs, k)
		case ast.Enum:
			f.enums = append(f.enums, k)
		case ast.Scalar:
			f.scalars = append(f.scalars, k)
		}

		sort.Strings(f.types)
		sort.Strings(f.inputs)
		sort.Strings(f.enums)
		sort.Strings(f.scalars)
	}
}

func (f *formatter) printHeader() {
	f.out.Write([]byte(`schema {
	query: Query
}

`))
}

func (f *formatter) printQuery() {
	f.out.Write([]byte("type Query {\n"))
	for _, field := range f.schema.Query.Fields {
		if strings.HasPrefix(field.Name, "__") {
			continue
		}
		f.out.Write([]byte{'\t'})
		f.out.Write([]byte(field.Name))
		if len(field.Arguments) != 0 {
			f.out.Write([]byte("(req: " + field.Arguments[0].Type.Name() + ")"))
		}
		f.out.Write([]byte(": "))
		f.out.Write([]byte(field.Type.Name()))
		f.out.Write([]byte{'!', '\n'})
	}
	f.out.Write([]byte{'}', '\n'})
}

func (f *formatter) printTypes() {
	for _, t := range f.types {
		f.out.Write([]byte{'\n'})
		typeDecl := f.schema.Types[t]
		f.out.Write([]byte("type " + typeDecl.Name + " {\n"))
		for _, field := range typeDecl.Fields {
			f.out.Write([]byte{'\t'})
			fmt.Fprintf(f.out, "%v: %v!\n", field.Name, field.Type.Name())
		}
		f.out.Write([]byte{'}', '\n'})
	}
}

func (f *formatter) printInputs() {
	for _, t := range f.inputs {
		f.out.Write([]byte{'\n'})
		typeDecl := f.schema.Types[t]
		f.out.Write([]byte("input " + typeDecl.Name + " {\n"))
		for _, field := range typeDecl.Fields {
			f.out.Write([]byte{'\t'})
			fmt.Fprintf(f.out, "%v: %v!\n", field.Name, field.Type.Name())
		}
		f.out.Write([]byte{'}', '\n'})
	}
}

func (f *formatter) printScalars() {
	for _, t := range f.scalars {
		f.out.Write([]byte{'\n'})
		typeDecl := f.schema.Types[t]
		f.out.Write([]byte("scalar " + typeDecl.Name + "\n"))
	}
}

func (f *formatter) printEnums() {
	for _, t := range f.enums {
		f.out.Write([]byte{'\n'})
		typeDecl := f.schema.Types[t]
		f.out.Write([]byte("enum " + typeDecl.Name + " {\n"))
		for _, field := range typeDecl.EnumValues {
			f.out.Write([]byte{'\t'})
			fmt.Fprintf(f.out, "%v\n", field.Name)
		}
		f.out.Write([]byte{'}', '\n'})
	}
}
