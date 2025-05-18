.PHONY: run

api-run:
	go run . api run

migrate-up:
	go run . migration up

migrate-down:
	go run . migration down

migrate-create:
	go run . migration create ${name}

migrate-seed:
	go run . migration seed