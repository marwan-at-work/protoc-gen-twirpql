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
		"query": "query q($req: HelloReq) {\n  hello(req: $req) {\n    text  }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	require.Equal(t, "twirpql", s.helloReq.GetName(), "Expected GraphQL request to populate Twirp Object")

	expected := `{"data":{"hello":{"text":"hello"}}}`

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
				"color": "GREEN",
				"trafficLights": ["YELLOW", "RED"]
			}
		},
		"query": "query q($req: TrafficJamReq) {\n  trafficJam(req: $req) {\n next }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	require.Equal(t, e2e.TrafficLight_GREEN, s.trafficJamReq.GetColor(), "Expected GraphQL request to populate Twirp Object")
	require.Equal(t, []e2e.TrafficLight{e2e.TrafficLight_YELLOW, e2e.TrafficLight_RED}, s.trafficJamReq.GetTrafficLights(), "Expected repeated enums to be equal")

	expected := `{"data":{"trafficJam":{"next":"YELLOW"}}}`

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
		"query": "query q {\n  getPainters {\n bestPainter {\n name }\n allPainters }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	expected := `{"data":{"getPainters":{"bestPainter":{"name":"picasso"},"allPainters":["one","two"]}}}`

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
		"query": "query q($req: TranslateReq) {\n  translate(req: $req) {\n translations }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	require.Equal(t, s.translateReq.GetWords(), map[string]*e2e.Word{"english": {Word: "hello"}})

	expected := `{"data":{"translate":{"translations":{"english":{"word":"hello"}}}}}`

	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")
}

func TestBread(t *testing.T) {
	s := &service{breadResp: &e2e.BreadResp{
		Answer: &e2e.BreadResp_Toasted{Toasted: true},
	}}
	h := twirpql.Handler(s, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{
		"operationName": "q",
		"variables": {
			"req": {
				"count": 3
			}
		},
		"query": "query q($req: BreadReq) {\n  bread(req: $req) {\n answer\n {\n __typename } \n }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	expected := `{"data":{"bread":{"answer":{"__typename":"BreadRespAnswerToasted"}}}}`
	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")

	require.Equal(t, s.breadReq.GetCount(), int64(3))
}

func TestMutations(t *testing.T) {
	s := &service{changeResp: &e2e.ChangeMeResp{
		Name: "james",
		Previous: map[string]*e2e.ChangeMeResp{
			"john": &e2e.ChangeMeResp{
				Name: "john",
				Previous: map[string]*e2e.ChangeMeResp{
					"jack": &e2e.ChangeMeResp{Name: "jack"},
				},
			},
		},
		Answer: &e2e.ChangeMeResp_Changed{
			Changed: true,
		},
	}}
	h := twirpql.Handler(s, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{
		"operationName": "m",
		"variables": {
			"req": {
				"name": "john",
				"previous": "{\"jack\": {\"name\": \"jack\"}}"
			}
		},
		"query": "mutation m($req: ChangeMeReq) {\n  changeMe(req: $req) {\n    name\n    previous\n    answer {\n      __typename\n      ... on ChangeMeRespAnswerChanged {\n        changed\n      }\n      ... on ChangeMeRespAnswerNewName {\n        newName\n      }\n    }\n  }\n}\n"
	}`))
	req.Header.Add("Content-Type", "application/json")
	h.ServeHTTP(w, req)

	expected := `{"data":{"changeMe":{"name":"james","previous":{"john":{"name":"john","Answer":null,"previous":{"jack":{"name":"jack","Answer":null}}}},"answer":{"__typename":"ChangeMeRespAnswerChanged","changed":true}}}}`
	require.Equal(t, expected, w.Body.String(), "Expected GraphQL query to return valid json")
	require.Equal(t, s.changeReq.GetName(), "john")
	require.Equal(t, s.changeReq.GetPrevious()["jack"].GetName(), "jack")
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
	breadReq       *e2e.BreadReq
	breadResp      *e2e.BreadResp
	changeReq      *e2e.ChangeMeReq
	changeResp     *e2e.ChangeMeResp
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

func (s *service) Bread(ctx context.Context, req *e2e.BreadReq) (*e2e.BreadResp, error) {
	s.breadReq = req
	return s.breadResp, s.err
}

func (s *service) ChangeMe(ctx context.Context, req *e2e.ChangeMeReq) (*e2e.ChangeMeResp, error) {
	s.changeReq = req
	return s.changeResp, s.err
}
