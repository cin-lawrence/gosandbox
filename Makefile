.PHONY: swagger build up down fmt test cover

swagger:
	@$(CURDIR)/scripts/gen-swagger.sh

build:
	@docker-compose build

up:
	@docker-compose up -d

down:
	@docker-compose down --remove-orphans

fmt:
	@go fmt ./...

test:
	@go test ./...

cover:
	@go test -cover -coverprofile=c.out ./...
	@go tool cover -html=c.out -o coverage.html
