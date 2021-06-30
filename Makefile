build:
	docker build . -t pursuit-portal-dock

run:
	docker run --rm --net pursuit_network --name portal -p 5001:5001 pursuit-portal-dock

pretty:
	go fmt `go list ./...`

unit-test:
	go test `go list ./... | grep -v cmd | grep -v test-integration | grep -v repo | grep -v vendor`

test:
	go test `go list ./... | grep -v vendor` -count=1
