rm *.go
protoc --go_out=plugins=grpc:. *.proto

