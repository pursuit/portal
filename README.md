# portal

![example workflow](https://github.com/pursuit/portal/actions/workflows/go.yml/badge.svg)

## Development Guide
### Pre-requisite
- [Go 1.16](https://golang.org/doc/install)

### Migration
- [Tool](https://github.com/golang-migrate/migrate)
```
migrate -source file:internal/migration -database postgres://postgres:password@localhost:5432/portal_development?sslmode=disable up
```

### Generate Proto
```
protoc --proto_path=pkg/proto --go_out=internal/proto/out --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_out=internal/proto/out user.proto
```