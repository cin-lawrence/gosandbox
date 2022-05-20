.PHONY: swagger up down fmt

swagger:
	@$(CURDIR)/scripts/gen-swagger.sh

up:
	@docker-compose up -d

down:
	@docker-compose down --remove-orphans

fmt:
	@go fmt ./...
