package cap

import "github.com/libs4go/errors"

//go:generate protoc --proto_path=./proto-cap --go_out=plugins=grpc,paths=source_relative:. cap.proto
//go:generate protoc --proto_path=./proto-cap --go_out=plugins=grpc,paths=source_relative:./blockchain blockchain.proto

// ScopeOfAPIError .
const errVendor = "go-cap"

// errors
var (
	ErrInternal = errors.New("the internal error", errors.WithVendor(errVendor), errors.WithCode(-1))
	ErrConfig   = errors.New("the config error", errors.WithVendor(errVendor), errors.WithCode(-2))
)
