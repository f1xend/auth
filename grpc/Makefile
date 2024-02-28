
LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

up_build:
	-docker compose up --build -d

up:
	-docker compose up -d

down:
	-docker compose down

dock_build_migrator_local:
	docker-compose up --build -d migrator-local

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/golang/protobuf/ptypes

generate:
	make generate-auth-api

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 --proto_path vendor.protogen \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/auth_v1 --grpc-gateway_opt=paths=source_relative \
    --plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
    --validate_out lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
    --plugin=protoc-gen-validate=bin/protoc-gen-validate \
    --openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
    --plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/auth_v1/auth.proto

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

create_t:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} create create_user_table sql

run_client:
	go run cmd/grpc_client/main.go
run_server:
	go run cmd/grpc_server/main.go

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@5.159.103.37:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t c~r.selcloud.ru/k1taro/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAA88zoPyPfjijKbcdHmS0Z8Qt7v2rGIqQG cr.selcloud.ru/k1taro
	docker push cr.selcloud.ru/k1taro/test-server:v0.0.1

ec:
	ECHO ${LOCAL_MIGRATION_DIR}

install_mock:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/f1xend/auth/internal/service/...,github.com/f1xend/auth/internal/api/... -count 5

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone git@github.com:bufbuild/protoc-gen-validate.git vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone git@github.com:googleapis/googleapis.git vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
        	mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
        	git clone git@github.com:grpc-ecosystem/grpc-gateway.git vendor.protogen/openapiv2 &&\
        	mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
        	rm -rf vendor.protogen/openapiv2 ;\
        fi

.PHONY: test-coverage
test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/f1xend/auth/internal/service/...,github.com/f1xend/auth/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

