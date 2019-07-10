package twirpql

type file struct {
	Service *service // TODO: multiple services
	Types   []*serviceType
	Inputs  []*serviceType
	Enums   []*enums
	Scalars []string
	Unions  []*union // TODO:
}

type service struct {
	Methods   []*method
	Mutations []*method
}

type enums struct {
	Name   string
	Fields []*serviceField
	Doc    string
}

type serviceType struct {
	Name   string
	Fields []*serviceField
	Doc    string
}

type serviceField struct {
	Name string
	Type string
	Doc  string
}

type method struct {
	Name, Request, Response string
	Doc                     string
}

type union struct {
	Name  string
	Types []string
}

type oneOf struct {
	Name   string
	Type   string
	Fields []*serviceField
}
