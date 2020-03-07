package genserver

var tmpl = `{{ reserveImport "context" }}
{{ reserveImport "net/http" }}

{{ reserveImport "github.com/99designs/gqlgen/graphql/handler" }}
{{ reserveImport "github.com/twitchtv/twirp" }}
{{ reserveImport "github.com/twitchtv/twirp/ctxsetters" }}
{{ reserveImport "github.com/99designs/gqlgen/graphql" }}
{{ reserveImport "github.com/99designs/gqlgen/graphql/handler/extension" }}
{{ reserveImport "github.com/99designs/gqlgen/graphql/handler/transport" }}
{{ reserveImport "github.com/99designs/gqlgen/graphql/playground" }}

// Playground is a proxy to github.com/99designs/gqlgen/handler.Playground
// All you need to do is provide a title and the URL Path to the GraphQL handler
func Playground(title, endpoint string) http.Handler {
	return playground.Handler(title, endpoint)
}

// Handler returns a handler to the GraphQL API.
// Server Hooks are optional but if present, they will
// be injected as GraphQL middleware.
func Handler(service {{lookupImport .ModPath}}.{{.ServiceName}}, hooks *twirp.ServerHooks) *handler.Server {
	es := NewExecutableSchema(Config{Resolvers: &Resolver{service}})
	srv := handler.New(es)
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	if hooks == nil {
		return srv
	}
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		f := graphql.GetFieldContext(ctx)
		parent := f.Parent.Path().String()
		if parent != "" {
			return next(ctx)
		}
		ctx = ctxsetters.WithMethodName(ctx, f.Field.Name)
		if hooks.RequestRouted != nil {
			ctx, err = hooks.RequestRouted(ctx)
			if err != nil {
				if terr, ok := err.(twirp.Error); ok && hooks.Error != nil {
					ctx = hooks.Error(ctx, terr)
				}
				return nil, err
			}
		}
		res, err = next(ctx)
		if terr, ok := err.(twirp.Error); ok && hooks.Error != nil {
			ctx = hooks.Error(ctx, terr)
		}
		return res, err
	})
	return srv
}
`
