test:
	@go test ./...

generate:
	@rm -rf ./internal/mocks/*
	@go install go.uber.org/mock/mockgen@latest
	@go generate -x ./internal/...

run-local:
	@docker-compose build
	@docker-compose up -d db
	@docker-compose up flyway
	@docker-compose run --rm bot