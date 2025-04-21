include .env
LOCAL_BIN:=$(PWD)/bin
PATH_PROTO:=$(PWD)/internal/transport/grpc/proto

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/spf13/cobra-cli@v1.3.0
#	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.23.0
#	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0
#	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3
#	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest #v1.36.6
#	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest #v1.71.1

cobra-init:
	$(LOCAL_BIN)/cobra-cli init

clean-cache:
	GOBIN=$(LOCAL_BIN) go clean -cache
	GOBIN=$(LOCAL_BIN) go clean -modcache

generate-proto:
	GOBIN=$(LOCAL_BIN) protoc --proto_path=$(PATH_PROTO) \
	--go_out=$(PATH_PROTO) --go_opt=paths=source_relative \
	--go-grpc_out=$(PATH_PROTO) --go-grpc_opt=paths=source_relative \
	$(PATH_PROTO)/someuser.proto

migration-status:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_MIGRATIONS_DSN) status -v
migration-add:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) create $(name) sql
migration-up:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_MIGRATIONS_DSN) up -v
migration-down:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_MIGRATIONS_DSN) down -v

gen-sql:
	$(LOCAL_BIN)/sqlc generate

swag:
	$(LOCAL_BIN)/swag init -g cmd/up.go

# docker basic
pa:
	docker ps -a
up:
	docker compose up --build -d
down:
	docker compose down
i:
	docker images
b:
	docker build -t someuser .
cleardb:
	rm -r ./mongo_data ./postgres_data
r: down up
cn:
	docker rmi $(docker images -q --filter="dangling=true") # clear images with none status


lu:
	docker logs -f --tail 100 someuser
eu:
	docker exec -it someuser sh
er:
	docker exec -it redis redis-cli