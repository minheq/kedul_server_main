start: ## run the app
	go build && ./kedul_server_main

test: ## run tests
	go test ./handlers -v

db-drop:
	dropdb kedul

db-create:
	createdb kedul

db-migrate-up: ## run database migrations up
	migrate -database "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable" -path ./migrations up

db-migrate-down: ## run database migrations down
	migrate -database "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable" -path ./migrations down

db-create-migration:  ## create database migration file. add name by adding name=<migration_file_name>
	migrate create -ext sql -dir ./migrations -seq $(name)