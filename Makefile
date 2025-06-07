.PHONY: run

api-run:
	go run ./cmd/api api run

migrate-up:
	go run ./cmd/migration migration up

migrate-down:
	go run ./migration/cmd migration down

migrate-create:
	go run ./cmd/migration migration create ${name}

migrate-seed:
	go run ./cmd/migration migration seed

migrate-fresh:
	go run ./cmd/migration migration fresh