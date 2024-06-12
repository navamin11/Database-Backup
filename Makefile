override filename = dbBackup
version = v1.0.0

.PHONY: clean-docker-images

run:
	@echo "Run Go Application"
	go run cmd/main.go

tidy:
	@echo "Install Packages"
	go mod tidy

rebuild-app:
	docker image prune -f && docker compose up --build myserver

down:
	docker compose down && docker image prune -f
		
ps:
	docker ps -a
	docker compose ps -a

up:
	docker compose up

rebuild-all:
	docker compose down && docker image prune -f && docker compose build --no-cache && docker compose up

clean:
	@echo "Removing all Docker images"
	@docker rmi $$(docker images -q) || true

# dump:
# 	docker exec postgres psql -U postgres pagila < pagila-data.sql

build:
	@echo "Build Application"
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/${filename}-${version}-linux cmd/main.go
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/${filename}-${version}-windows.exe cmd/main.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/${filename}-${version}-mac cmd/main.go
	chmod 755 ./bin/${filename}-${version}*