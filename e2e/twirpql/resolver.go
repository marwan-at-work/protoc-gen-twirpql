// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package twirpql

import (
	"context"

	"marwan.io/protoc-gen-twirpql/e2e"
)

type Resolver struct {
	e2e.Service
}

func (r *Resolver) BreadResp() BreadRespResolver {
	return &breadRespResolver{r}
}
func (r *Resolver) ChangeMeResp() ChangeMeRespResolver {
	return &changeMeRespResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) TranslateResp() TranslateRespResolver {
	return &translateRespResolver{r}
}

type breadRespResolver struct{ *Resolver }

func (r *breadRespResolver) Answer(ctx context.Context, obj *e2e.BreadResp) (unionMask, error) {
	return obj.GetAnswer(), nil
}

type changeMeRespResolver struct{ *Resolver }

func (r *changeMeRespResolver) Previous(ctx context.Context, obj *e2e.ChangeMeResp) (Previous, error) {
	return obj.GetPrevious(), nil
}

func (r *changeMeRespResolver) Answer(ctx context.Context, obj *e2e.ChangeMeResp) (unionMask, error) {
	return obj.GetAnswer(), nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) ChangeMe(ctx context.Context, req *e2e.ChangeMeReq) (*e2e.ChangeMeResp, error) {
	return r.Service.ChangeMe(ctx, req)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Hello(ctx context.Context, req *e2e.HelloReq) (*e2e.HelloResp, error) {
	return r.Service.Hello(ctx, req)
}

func (r *queryResolver) TrafficJam(ctx context.Context, req *e2e.TrafficJamReq) (*e2e.TrafficJamResp, error) {
	return r.Service.TrafficJam(ctx, req)
}

func (r *queryResolver) GetPainters(ctx context.Context) (*e2e.PaintersResp, error) {
	return r.Service.GetPainters(ctx, nil)
}

func (r *queryResolver) Translate(ctx context.Context, req *e2e.TranslateReq) (*e2e.TranslateResp, error) {
	return r.Service.Translate(ctx, req)
}

func (r *queryResolver) Bread(ctx context.Context, req *e2e.BreadReq) (*e2e.BreadResp, error) {
	return r.Service.Bread(ctx, req)
}

type translateRespResolver struct{ *Resolver }

func (r *translateRespResolver) Translations(ctx context.Context, obj *e2e.TranslateResp) (Translations, error) {
	return obj.GetTranslations(), nil
}
