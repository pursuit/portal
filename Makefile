run:
	go run cmd/api/main.go

pretty:
	gofmt -s -w .

unit-test:
	go test `go list ./... | grep -v cmd | grep -v test-integration | grep -v repo`

test:
	go test ./... -count=1
