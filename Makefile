.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:internal/app/genproto/customer -I api/protobuf api/protobuf/customer.proto
	protoc --go_out=plugins=grpc:internal/app/genproto/trainer -I api/protobuf api/protobuf/trainer.proto