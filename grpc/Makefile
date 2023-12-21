LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/golang/protobuf/ptypes

generate:
	make generate-auth-api

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth_v1/auth.proto

#run_client:
#	go run cmd/grpc_client/main.go
#run_server:
#	go run cmd/grpc_server/main.go

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@5.159.103.37:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t c~r.selcloud.ru/k1taro/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAA88zoPyPfjijKbcdHmS0Z8Qt7v2rGIqQG cr.selcloud.ru/k1taro
	docker push cr.selcloud.ru/k1taro/test-server:v0.0.1

