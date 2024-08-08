.PHONY: proto
proto:
	protoc \
		-I api/protobuf api/protobuf/customer.proto \
		--go_opt=paths=source_relative --go_out=internal/app/genproto/customer \
		--go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
		--go-grpc_out=internal/app/genproto/customer
	protoc \
		-I api/protobuf api/protobuf/trainer.proto \
		--go_opt=paths=source_relative --go_out=internal/app/genproto/trainer \
		--go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
		--go-grpc_out=internal/app/genproto/trainer