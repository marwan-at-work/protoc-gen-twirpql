package twirpql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/twitchtv/twirp"
	"marwan.io/protoc-gen-twirpql/e2e"
)

// Playground is a proxy to github.com/99designs/gqlgen/handler.Playground
// All you need to do is provide a title and the URL Path to the GraphQL handler
func Playground(title, endpoint string) http.Handler {
	return handler.Playground(title, endpoint)
}

// Handler returns a handler to the GraphQL API.
// Server Hooks are optional but if present, they will
// be injected as GraphQL middleware.
func Handler(service e2e.Service, hooks *twirp.ServerHooks) http.Handler {
	if hooks == nil {
		return handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{service}}))
	}
	h := &middlewareHooks{hooks}
	return handler.GraphQL(
		NewExecutableSchema(
			Config{Resolvers: &Resolver{service}},
		),
		handler.ResolverMiddleware(h.withErr),
	)
}

type middlewareHooks struct {
	hooks *twirp.ServerHooks
}

func (h *middlewareHooks) withErr(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	res, err = next(ctx)
	if err != nil && h.hooks.Error != nil {
		terr, ok := err.(twirp.Error)
		if ok {
			h.hooks.Error(ctx, terr)
		} else {
			fmt.Println("Twirp err does not implement twirp.Error:", err)
		}
	}
	return res, err
}
