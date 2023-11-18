.PHONY: lint, run, test, coverage, coverage_html
lint:
	gofmt -w .
	golangci-lint run
	go vet ./...

run:
	go run cmd/main.go

test:
	go test -v -timeout 60s -coverprofile=coverage.out -cover ./...
	go tool cover -func coverage.out

coverage:
	go test ./... -coverprofile=coverage.out

coverage_html:
	@$(MAKE) coverage
	go tool cover -html=coverage.out



