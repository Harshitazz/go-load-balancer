run:
	go run main.go

build:
	go build -o load-balancer main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

clean:
	rm -f load-balancer
