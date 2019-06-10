package gen

//go:generate protoc -I . -I ../.. -I /usr/local/include --debug_out=.:. mutations.proto
