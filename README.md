# Service Salon

## Development

### Start
```zshrc
go build
./kedul_server_main
```

### Migrate Up
```zshrc
migrate -database "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable" -path ./migrations up
```

### Migrate Down
```zshrc
migrate -database "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable" -path ./migrations down
```

### Create new migration file
```zshrc
migrate create -ext sql -dir ./migrations -seq create_users_table
```
