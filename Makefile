pretty:
	gofmt -s -w .

unit-test:
	go test `go list ./... | grep -v cmd | grep -v test-integration`

test:
	go test ./...
