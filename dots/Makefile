lint: bin/golangci-lint
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run -v ./...

fmt: bin/golangci-lint
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest fmt -v ./...

bin/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.11.3
	bin/golangci-lint --version
