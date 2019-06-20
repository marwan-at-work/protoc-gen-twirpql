package e2e

//go:generate protoc -I . -I /usr/local/include -I ../options --go_out=. --twirp_out=. service.proto
//go:generate protoc -I . -I /usr/local/include -I ../options --twirpql_out=. service.proto
