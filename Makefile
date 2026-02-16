
mod:
	go get -tool "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	sudo rm -r ./docker/bank-db-data
	go mod tidy
	cd tools && go mod tidy

generate:
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/oapi.yaml -package accounts -o ./internal/ports/http/accounts/api.gen.go api/accounts.yaml
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/oapi_models.yaml -package accounts -o ./internal/ports/http/accounts/model.gen.go api/accounts.yaml
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/oapi.yaml -package transactions -o ./internal/ports/http/transactions/api.gen.go api/transactions.yaml
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/oapi_models.yaml -package transactions -o ./internal/ports/http/transactions/model.gen.go api/transactions.yaml
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/oapi_models.yaml -package common -o ./internal/ports/http/common/model.gen.go api/common.yaml

migrate_up:
	go build -o ./bin/migrate ./cmd/migrations
	./bin/migrate up -c ./config/config_local.yaml -m ./internal/ports/db/migrations

migrate_down:
	go build -o ./bin/migrate ./cmd/migrations
	./bin/migrate down -c ./config/config_local.yaml -m ./internal/ports/db/migrations

build_app:
	go build -o ./bin/bankapp ./cmd/bankapp