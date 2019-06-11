package e2e

//go:generate protoc -I . -I /usr/local/include -I ../twirpql --go_out=. --twirp_out=. service.proto
//go:generate protoc -I . -I /usr/local/include -I ../twirpql --twirpql_out=. service.proto
