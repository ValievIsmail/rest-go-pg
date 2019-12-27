go: 
	go run cmd/*.go
	
build:
	go build -o ./rest-go-pg cmd/*.go

up:
	docker-compose up -d

down:
	docker-compose down

docker-build:
	docker build -t rest-go-pg .

docker-run:
	docker run -d -p 8080:8080 rest-go-pg

bomb-comments:
	bombardier -c 30 -n 3000 -l http://localhost:8080/api/v1/comment/100