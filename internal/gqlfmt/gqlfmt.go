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
func Print(input string, out io.Writer) error {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: input,
	})
	if err != nil {
		return err
	}
	f := &formatter{schema: schema, out: out}
	f.printSchema()

	return nil
}

// PrintSchema formats a given schema and returns
// the output as a string
func PrintSchema(s *ast.Schema) (string, error) {
	var out strings.Builder
	f := &formatter{schema: s, out: &out}
	f.printSchema()
	return out.String(), nil
}

type formatter struct {
	schema     *ast.Schema
	out        io.Writer
	types      []string
	inputs     []string
	scalars    []string
	enums      []string
	unions     []string
	directives []string
}

func (f *formatter) printSchema() {
	f.sortDeclarations()
	f.printQuery()
	f.printMutation()
	f.printTypes()
	f.printInputs()
	f.printEnums()
	f.printScalars()
	f.printUnions()
	f.printDirectiveDefs()
}

func (f *formatter) sortDeclarations() {
	for k, def := range f.schema.Types {
		if k == "Query" || k == "Mutation" || def.BuiltIn {
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
		case ast.Union:
			f.unions = append(f.unions, k)
		}
	}
	sort.Strings(f.types)
	sort.Strings(f.inputs)
	sort.Strings(f.enums)
	sort.Strings(f.scalars)
	sort.Strings(f.unions)
	for _, d := range f.schema.Directives {
		if d.Position.Src.BuiltIn {
			continue
		}
		f.directives = append(f.directives, d.Name)
	}
	sort.Strings(f.directives)
}

func (f *formatter) printQuery() {
	f.printDoc(f.schema.Query.Description, 0)
	f.print("type Query")
	f.printDirectives(f.schema.Query.Directives)
	f.print(" {\n")
	for _, field := range f.schema.Query.Fields {
		if strings.HasPrefix(field.Name, "__") {
			continue
		}
		f.printDoc(field.Description, 1)
		f.printf("\t%v", field.Name)
		f.printArgs(field.Arguments)
		f.printf(": %v\n", field.Type.String())
	}
	f.print("}\n")
}

func (f *formatter) printArgs(aa ast.ArgumentDefinitionList) {
	if len(aa) == 0 {
		return
	}
	args := []string{}
	for _, a := range aa {
		arg := a.Name
		arg += ": "
		arg += a.Type.String()
		args = append(args, arg)
	}

	f.printf("(%v)", strings.Join(args, ", "))
}

func (f *formatter) fmtDirectiveArgument(a *ast.Argument) string {
	var b strings.Builder
	b.WriteString(a.Name)
	b.WriteString(": ")
	b.WriteString(a.Value.String())

	return b.String()
}

func (f *formatter) printMutation() {
	if f.schema.Mutation == nil || len(f.schema.Mutation.Fields) == 0 {
		return
	}
	f.print("\ntype Mutation {\n")
	for _, field := range f.schema.Mutation.Fields {
		if strings.HasPrefix(field.Name, "__") {
			continue
		}
		doc := strings.TrimSpace(field.Description)
		if doc != "" {
			f.printDoc(doc, 1)
		}
		f.printf("\t%v", field.Name)
		f.printArgs(field.Arguments)
		f.printf(": %v\n", field.Type.String())
	}
	f.print("}\n")
}

func (f *formatter) printTypes() {
	for _, t := range f.types {
		f.print("\n")
		typeDecl := f.schema.Types[t]
		f.printDoc(typeDecl.Description, 0)
		f.printf("type %v", typeDecl.Name)
		f.printDirectives(typeDecl.Directives)
		f.print(" {\n")
		for _, field := range typeDecl.Fields {
			f.printDoc(field.Description, 1)
			f.printf("\t%v: %v\n", field.Name, field.Type.String())
		}
		f.print("}\n")
	}
}

func (f *formatter) printDirectives(dirs []*ast.Directive) {
	if len(dirs) > 0 {
		f.print(" ")
	}
	for _, dir := range dirs {
		f.printDirective(dir)
	}
}

func (f *formatter) printDirective(d *ast.Directive) {
	f.printf("@%v", d.Name)
	if len(d.Arguments) > 0 {
		f.print(`(`)
		args := []string{}
		for _, a := range d.Arguments {
			args = append(args, f.fmtDirectiveArgument(a))
		}
		f.print(strings.Join(args, ", "))
		f.print(`)`)
	}
}

func (f *formatter) printInputs() {
	for _, t := range f.inputs {
		f.println()
		typeDecl := f.schema.Types[t]
		f.printDoc(typeDecl.Description, 0)
		f.printf("input %v {\n", typeDecl.Name)
		for _, field := range typeDecl.Fields {
			f.printDoc(field.Description, 1)
			f.printf("\t%v: %v\n", field.Name, field.Type.String())
		}
		f.println("}")
	}
}

func (f *formatter) printEnums() {
	for _, t := range f.enums {
		f.println()
		typeDecl := f.schema.Types[t]
		f.printDoc(typeDecl.Description, 0)
		f.printf("enum %v {\n", typeDecl.Name)
		for _, field := range typeDecl.EnumValues {
			f.printDoc(field.Description, 1)
			f.printf("\t%v\n", field.Name)
		}
		f.println("}")
	}
}

func (f *formatter) printScalars() {
	for _, t := range f.scalars {
		f.println()
		typeDecl := f.schema.Types[t]
		f.printf("scalar %v\n", typeDecl.Name)
	}
}

func (f *formatter) printUnions() {
	if len(f.unions) > 0 {
		f.println()
	}
	for _, t := range f.unions {
		decl := f.schema.Types[t]
		sort.Strings(decl.Types)
		f.printf("union %v = %v\n", decl.Name, strings.Join(decl.Types, " | "))
	}
}

func (f *formatter) printDirectiveDefs() {
	if len(f.directives) > 0 {
		f.println()
	}
	for _, t := range f.directives {
		decl := f.schema.Directives[t]
		locs := []string{}
		for _, l := range decl.Locations {
			locs = append(locs, string(l))
		}
		sort.Strings(locs)
		args := ""
		if len(decl.Arguments) > 0 {
			args += "("
			argList := []string{}
			for _, a := range decl.Arguments {
				arg := a.Name
				arg += ": "
				arg += a.Type.String()
				if a.DefaultValue != nil {
					arg += " = "
					arg += a.DefaultValue.String()
				}
				argList = append(argList, arg)
			}
			args += strings.Join(argList, ", ")
			args += ")"
		}
		f.printf("directive @%v%v on %v\n", decl.Name, args, strings.Join(locs, " | "))
	}
}

func (f *formatter) printDoc(doc string, indent int) {
	doc = strings.TrimSpace(doc)
	if doc == "" {
		return
	}
	tab := strings.Repeat("\t", indent)
	f.print(tab)
	f.print(`"""`)
	f.println()
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		f.print(tab)
		f.print(line)
		f.println()
	}
	f.print(tab)
	f.print(`"""`)
	f.println()
}

func (f *formatter) print(a ...interface{}) {
	fmt.Fprint(f.out, a...)
}

func (f *formatter) println(a ...interface{}) {
	fmt.Fprintln(f.out, a...)
}

func (f *formatter) printf(s string, a ...interface{}) {
	fmt.Fprintf(f.out, s, a...)
}
