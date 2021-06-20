build:
	docker build . -t pursuit-portal-dock

run:
	docker run --net pursuit_network -p 5001:5001 pursuit-portal-dock

pretty:
	gofmt -s -w .

unit-test:
	go test `go list ./... | grep -v cmd | grep -v test-integration | grep -v repo`

test:
	go test ./... -count=1
