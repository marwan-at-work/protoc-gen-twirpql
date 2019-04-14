package e2e_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"marwan.io/protoc-gen-twirpql/e2e"
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

type service struct {
	e2e.Service
	helloReq       *e2e.HelloReq
	helloResp      *e2e.HelloResp
	trafficJamReq  *e2e.TrafficJamReq
	trafficJamResp *e2e.TrafficJamResp
	err            error
}

func (s *service) Hello(ctx context.Context, req *e2e.HelloReq) (*e2e.HelloResp, error) {
	s.helloReq = req
	return s.helloResp, s.err
}

func (s *service) TrafficJam(ctx context.Context, req *e2e.TrafficJamReq) (*e2e.TrafficJamResp, error) {
	s.trafficJamReq = req
	return s.trafficJamResp, nil
}
