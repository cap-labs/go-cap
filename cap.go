package cap

//go:generate protoc --proto_path=./proto-cap --go_out=plugins=grpc,paths=source_relative:. cap.proto
