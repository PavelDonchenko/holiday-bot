.PHONY: lint, run, test, coverage, coverage_html
lint:
	gofmt -w .
	goimports -local git.foxminded.ua/foxstudent106361/holiday-bot -w .
	golangci-lint run
	go vet ./...

run_worker:
	go run cmd/main.go worker

run_bot:
	go run cmd/main.go bot

compose_up:
	docker-compose -f docker-compose.yml up --build --remove-orphans

test:
	go test -v -timeout 60s -coverprofile=coverage.out -cover ./...
	go tool cover -func coverage.out

coverage:
	go test ./... -coverprofile=coverage.out

coverage_html:
	@$(MAKE) coverage
	go tool cover -html=coverage.out



