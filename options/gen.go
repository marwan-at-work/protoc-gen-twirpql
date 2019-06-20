package options

//go:generate protoc -I . -I /usr/local/include --go_out=paths=source_relative:. options.proto
