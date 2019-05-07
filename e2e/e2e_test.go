package e2e_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"marwan.io/protoc-gen-twirpql/e2e"
	"marwan.io/protoc-gen-twirpql/e2e/painters"
	"marwan.io/protoc-gen-twirpql/e2e/twirpql"
)

func TestHello(t *testing.T) {
	s := &service{helloResp: &e2e.HelloResp{Text: "hello"}}
	h := twirpql.Handler(s, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{
		"operationName": "q",
		"variables": {
			"req": {
				"name": "twirpql"
			}
		},
		"query": "query q($req: HelloReq) {\n  Hello(req: $req) {\n    text  }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	require.Equal(t, "twirpql", s.helloReq.GetName(), "Expected GraphQL request to populate Twirp Object")

	expected := `{"data":{"Hello":{"text":"hello"}}}`

	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")
}

func TestTrafficJam(t *testing.T) {
	s := &service{trafficJamResp: &e2e.TrafficJamResp{Next: e2e.TrafficLight_YELLOW}}
	h := twirpql.Handler(s, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{
		"operationName": "q",
		"variables": {
			"req": {
				"color": "GREEN"
			}
		},
		"query": "query q($req: TrafficJamReq) {\n  TrafficJam(req: $req) {\n next }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	require.Equal(t, e2e.TrafficLight_GREEN, s.trafficJamReq.GetColor(), "Expected GraphQL request to populate Twirp Object")

	expected := `{"data":{"TrafficJam":{"next":"YELLOW"}}}`

	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")
}

func TestPainters(t *testing.T) {
	s := &service{paintersResp: &e2e.PaintersResp{
		BestPainter: &painters.Painter{Name: "picasso"},
		AllPainters: []string{"one", "two"},
	}}
	h := twirpql.Handler(s, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{
		"operationName": "q",
		"variables": {},
		"query": "query q {\n  GetPainters {\n bestPainter {\n name }\n allPainters }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	expected := `{"data":{"GetPainters":{"bestPainter":{"name":"picasso"},"allPainters":["one","two"]}}}`

	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")
}

func TestTranslate(t *testing.T) {
	s := &service{translateResp: &e2e.TranslateResp{
		Translations: map[string]*e2e.Word{
			"english": &e2e.Word{
				Word: "hello",
			},
		},
	}}
	h := twirpql.Handler(s, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{
		"operationName": "q",
		"variables": {
			"req": {
				"words": "{\"english\": {\"word\": \"hello\"}}"
			}
		},
		"query": "query q($req: TranslateReq) {\n  Translate(req: $req) {\n translations }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	require.Equal(t, s.translateReq.GetWords(), map[string]*e2e.Word{"english": {Word: "hello"}})

	expected := `{"data":{"Translate":{"translations":{"english":{"word":"hello"}}}}}`

	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")
}

type service struct {
	e2e.Service
	helloReq       *e2e.HelloReq
	helloResp      *e2e.HelloResp
	trafficJamReq  *e2e.TrafficJamReq
	trafficJamResp *e2e.TrafficJamResp
	paintersReq    *e2e.PaintersReq
	paintersResp   *e2e.PaintersResp
	translateReq   *e2e.TranslateReq
	translateResp  *e2e.TranslateResp
	err            error
}

func (s *service) Hello(ctx context.Context, req *e2e.HelloReq) (*e2e.HelloResp, error) {
	s.helloReq = req
	return s.helloResp, s.err
}

func (s *service) TrafficJam(ctx context.Context, req *e2e.TrafficJamReq) (*e2e.TrafficJamResp, error) {
	s.trafficJamReq = req
	return s.trafficJamResp, s.err
}

func (s *service) GetPainters(ctx context.Context, req *e2e.PaintersReq) (*e2e.PaintersResp, error) {
	s.paintersReq = req
	return s.paintersResp, s.err
}

func (s *service) Translate(ctx context.Context, req *e2e.TranslateReq) (*e2e.TranslateResp, error) {
	s.translateReq = req
	return s.translateResp, s.err
}
