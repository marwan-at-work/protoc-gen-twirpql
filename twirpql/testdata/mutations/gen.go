package gen

//go:generate protoc -I . -I ../../../options -I /usr/local/include --debug_out=.:. mutations.proto
