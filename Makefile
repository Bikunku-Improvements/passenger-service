protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpc/pb/*.proto

load-test:
	ghz --config=benchmark/config.json

test:
	pwd