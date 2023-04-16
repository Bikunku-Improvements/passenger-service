protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--python-out=. --grpc_python_out=source_relative \
		grpc/pb/*.proto

