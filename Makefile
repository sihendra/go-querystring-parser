dep:
	@echo ">> Downloading Dependencies"
	@go mod download

test: dep
	@echo ">> Running Unit Test"
	@go test -count=1 -failfast -cover -covermode=atomic -coverprofile coverage.out ./...
