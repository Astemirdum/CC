
.PHONY: example
example:
	go run cmd/main.go example

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run --fix

.PHONY: test
test:
	go test -v -race -timeout 90s -count=1 -shuffle=on  -coverprofile cover.out ./...
	@go tool cover -func cover.out | grep total | awk '{print $3}'
	go tool cover -html="cover.out" -o coverage.html
