# Service Salon

## Scripts

```zshrc
migrate -database "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable" -path ./migrations up
```

```zshrc
migrate create -ext sql -dir ./migrations -seq create_users_table
```
