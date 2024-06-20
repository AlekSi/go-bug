all:
	go clean -cache
	go run getwd.go
	go test -v getwd_test.go
	go build -o getwd_vet getwd_vet.go
	go vet -vettool=getwd_vet
