# speedroute
Speedrun route planning algorithm, finding the shortest possible path in a directed graph.

Project is "soon" hosted on [https://speedroute.org](https://speedroute.org).

## Getting started
1. Install go (v1.20 or higher)
2. Install & setup postgresql
3. Install flyway CLI
4. Run `./flyway-migrate.sh`
5. Run `go build` & `go test ./...`
6. Run `golangci-lint -v run`
7. Run `go run main.go`
8. Go to [localhost:8001](https://localhost:8001/).

### Versions used
```
go version go1.14.4 darwin/amd64
flyway-6.5.1
psql (PostgreSQL) 10.12 (Ubuntu 10.12-0ubuntu0.18.04.1)
golangci/tap/golangci-lint 1.27.0
```
