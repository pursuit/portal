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
