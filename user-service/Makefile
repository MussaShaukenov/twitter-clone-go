cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./cmd/...
	go tool cover -html=coverage.out
	rm coverage.out
goose_up:
	goose -dir internal/migrations postgres "postgresql://postgres:admin@127.0.0.1:5432/user_db?sslmode=disable" up
goose_down:
	goose -dir internal/migrations postgres "postgresql://postgres:admin@127.0.0.1:5432/user_db?sslmode=disable" down
goose_create:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Error: MIGRATION_NAME is required. Usage: make goose_create MIGRATION_NAME=<name>"; \
		exit 1; \
	fi
	goose create -s $(MIGRATION_NAME) sql -dir internal/migrations
