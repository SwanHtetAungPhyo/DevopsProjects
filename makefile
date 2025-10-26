# Load environment variables if not in CI
ifneq ($(CI),true)
ifneq ($(GITHUB_ACTIONS),true)
    ifneq ("$(wildcard .env)","")
        include .env
        export $(shell sed 's/=.*//' .env)
    endif
endif
endif

DOCKER_DIR=docker

.PHONY: docker-compose-up help

docker-compose-up:
	@echo "Enter the file name:"
	@read -p "File: " name; \
	docker-compose -f $$name up -d


up-dev:
	docker-compose -f $(DOCKER_DIR)/docker-compose.dev.yml up -d

up-prod:
	docker-compose -f $(DOCKER_DIR)/docker-compose.prod.yml up -d

up-nats:
	docker-compose -f docker-compose.nats.yml up -d


dev-up:
	docker-compose -f docker-compose.yml up -d app-one app-two app-three

dev-down:
	docker-compose -f docker-compose.yml down

dev-logs:
	docker-compose -f docker-compose.yml logs -f

# Helper targets
status:
	docker-compose -f docker-compose.yml ps

clean:
	docker-compose -f docker-compose.yml down -v --remove-orphans

help:
	@echo "Available commands:"
	@echo "  make docker-compose-up    - Interactive compose file selection"
	@echo "  make up-dev              - Start development environment"
	@echo "  make up-prod             - Start production environment"
	@echo "  make up-nats             - Start NATS messaging"
	@echo "  make dev-up              - Start all apps"
	@echo "  make dev-down            - Stop all apps"
	@echo "  make status              - Show container status"
	@echo "  make clean               - Clean up containers and volumes"