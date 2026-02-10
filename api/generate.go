package api

//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi.yaml -package accounts -o ./../internal/ports/http/accounts/api.gen.go accounts.yaml
//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi_models.yaml -package accounts -o ./../internal/ports/http/accounts/model.gen.go accounts.yaml

//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi.yaml -package transactions -o ./../internal/ports/http/transactions/api.gen.go transactions.yaml
//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi_models.yaml -package transactions -o ./../internal/ports/http/transactions/model.gen.go transactions.yaml

//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi_models.yaml -package common -o ./../internal/ports/http/common/model.gen.go common.yaml
