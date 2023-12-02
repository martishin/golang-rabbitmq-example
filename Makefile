build:
	golangci-lint run
	docker build -t golang-rabbitmq-example-producer -f Dockerfile.producer .
	docker build -t golang-rabbitmq-example-consumer -f Dockerfile.consumer .

run:
	docker-compose up
