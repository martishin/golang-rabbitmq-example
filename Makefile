build:
	golangci-lint run
	docker compose build --no-cache

run:
	docker compose up
