POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=denetdb
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
URL="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DB)?sslmode=disable"
.PHONY: migrate migrate-u migrate-d migrate-f1 migrate-f2 u ud d dv run run-prod
migrate:
	go run cmd/migrate/migrate.go --url=$(URL)
migrate-u:
	go run cmd/migrate/migrate.go --url=$(URL) --steps=1
migrate-d:
	go run cmd/migrate/migrate.go --url=$(URL) --steps=-1
migrate-f1:
	go run cmd/migrate/migrate.go --url=$(URL) --force=true --steps=1
migrate-f2:
	go run cmd/migrate/migrate.go --url=$(URL) --force=true --steps=2

u:
	docker compose up --build --remove-orphans
ud:
	docker compose up --build --remove-orphans -d
d:
	docker compose down 
dv:
	docker compose down 
run:
	go run cmd/denet/main.go
run-prod:
	go run cmd/denet/main.go --env=prod