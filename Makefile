go: 
	go run cmd/*.go
	
build:
	go build -o ./restgopg cmd/*.go

up:
	docker-compose up -d

down:
	docker-compose down

bomb-comments:
	bombardier -c 30 -n 3000 -l http://localhost:8080/api/v1/comment/100