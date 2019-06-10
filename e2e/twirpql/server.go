package twirpql

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/twitchtv/twirp"
	"github.com/twitchtv/twirp/ctxsetters"
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
func Handler(service e2e.Service, hooks *twirp.ServerHooks, opts ...handler.Option) http.Handler {
	if hooks == nil {
		return handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{service}}), opts...)
	}
	h := &middlewareHooks{hooks}
	opts = append([]handler.Option{handler.ResolverMiddleware(h.hook)}, opts...)
	return handler.GraphQL(
		NewExecutableSchema(
			Config{Resolvers: &Resolver{service}},
		),
		opts...,
	)
}

type middlewareHooks struct {
	hooks *twirp.ServerHooks
}

func (h *middlewareHooks) hook(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	ifc := graphql.GetResolverContext(ctx).Path()
	if len(ifc) > 0 {
		queryName, ok := ifc[0].(string)
		if ok {
			ctx = ctxsetters.WithMethodName(ctx, strings.Title(queryName))
		}
	}
	if h.hooks.RequestRouted != nil {
		ctx, err = h.hooks.RequestRouted(ctx)
		if err != nil {
			return nil, err
		}
	}
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
