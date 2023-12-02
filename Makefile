build:
	docker build -t golang-redis-example-producer -f Dockerfile.producer . ; \
  	docker build -t golang-redis-example-consumer -f Dockerfile.consumer .

run:
	docker-compose up --build
