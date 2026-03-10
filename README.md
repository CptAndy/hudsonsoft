# Universal Point-of-Sale
A universal point-of-sale system that allows for the client to access managerial permissions, business
operations, as well as customer interaction.


## Start

### Prerequisites

- [Go](https://golang.org/doc/install) latest
- [PostgreSQL](https://www.postgresql.org/download/) latest
- [pgAdmin](https://www.pgadmin.org/download/) latest

### Api testing
- [Postman](https://www.postman.com/downloads/)
- [Bruno](https://www.usebruno.com/downloads)


### Set environment variables in .env
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=hudsonsoft
```
### Run database migrations
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/hudsonsoft?sslmode=disable" up
```

### Start the Go server
```bash
go run cmd/api/main.go
```
