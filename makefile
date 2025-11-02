# DevOps Projects Makefile
# Author: Swan Htet Aung Phyo
# Purpose: Unified build and deployment automation

# Load environment variables if not in CI
ifneq ($(CI),true)
ifneq ($(GITHUB_ACTIONS),true)
    ifneq ("$(wildcard .env)","")
        include .env
        export $(shell sed 's/=.*//' .env)
    endif
endif
endif

# Project directories
DOCKER_DIR := docker
PROJECTS_DIR := projects
INFRA_DIR := infra
SCRIPTS_DIR := scripts

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help setup verify clean status

# Default target
all: help

## Setup and Installation
setup: ## Set up development environment
	@echo "$(BLUE)Setting up DevOps development environment...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)Creating .env file from template...$(NC)"; \
		cp .env.example .env 2>/dev/null || echo "# Add your environment variables here" > .env; \
	fi
	@echo "$(GREEN)Environment setup complete!$(NC)"
	@echo "$(YELLOW)Please edit .env file with your configuration$(NC)"

verify: ## Verify installation and dependencies
	@echo "$(BLUE)Verifying system dependencies...$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is required but not installed$(NC)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(RED)Docker Compose is required but not installed$(NC)"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "$(YELLOW)Go is recommended for building custom tools$(NC)"; }
	@echo "$(GREEN)System verification complete!$(NC)"

## Docker Environment Management
up-dev: ## Start development environment
	@echo "$(BLUE)Starting development environment...$(NC)"
	@if [ -f $(DOCKER_DIR)/docker-compose.dev.yml ]; then \
		docker-compose -f $(DOCKER_DIR)/docker-compose.dev.yml up -d; \
	else \
		echo "$(YELLOW)Development compose file not found, using default$(NC)"; \
		docker-compose up -d; \
	fi

up-prod: ## Start production environment
	@echo "$(BLUE)Starting production environment...$(NC)"
	@if [ -f $(DOCKER_DIR)/docker-compose.prod.yml ]; then \
		docker-compose -f $(DOCKER_DIR)/docker-compose.prod.yml up -d; \
	else \
		echo "$(RED)Production compose file not found$(NC)"; \
		exit 1; \
	fi

up-nats: ## Start NATS messaging system
	@echo "$(BLUE)Starting NATS messaging system...$(NC)"
	@docker-compose -f docker-compose.nats.yml up -d

## Project-Specific Commands
docker-log-agent: ## Deploy Docker log monitoring agent
	@echo "$(BLUE)Deploying Docker Log Agent...$(NC)"
	@cd $(PROJECTS_DIR)/docker-log-agent && docker-compose up -d
	@echo "$(GREEN)Docker Log Agent deployed successfully$(NC)"

event-driven: ## Start event-driven microservices
	@echo "$(BLUE)Starting event-driven microservices...$(NC)"
	@cd $(PROJECTS_DIR)/event-driven && docker-compose up -d
	@echo "$(GREEN)Event-driven services started$(NC)"

scl-build: ## Build SCL configuration language tool
	@echo "$(BLUE)Building SCL tool...$(NC)"
	@cd $(PROJECTS_DIR)/scl/lang && go build -o scl
	@echo "$(GREEN)SCL tool built successfully$(NC)"

## Infrastructure Management
infra-plan: ## Plan infrastructure changes (Terraform)
	@echo "$(BLUE)Planning infrastructure changes...$(NC)"
	@cd $(INFRA_DIR)/terraform && terraform plan

infra-apply: ## Apply infrastructure changes (Terraform)
	@echo "$(BLUE)Applying infrastructure changes...$(NC)"
	@cd $(INFRA_DIR)/terraform && terraform apply

infra-destroy: ## Destroy infrastructure (Terraform)
	@echo "$(RED)Destroying infrastructure...$(NC)"
	@cd $(INFRA_DIR)/terraform && terraform destroy

pulumi-preview: ## Preview Pulumi infrastructure changes
	@echo "$(BLUE)Previewing Pulumi changes...$(NC)"
	@cd $(INFRA_DIR)/pulumi/event && pulumi preview

pulumi-up: ## Deploy Pulumi infrastructure
	@echo "$(BLUE)Deploying Pulumi infrastructure...$(NC)"
	@cd $(INFRA_DIR)/pulumi/event && pulumi up

ansible-ping: ## Test Ansible connectivity
	@echo "$(BLUE)Testing Ansible connectivity...$(NC)"
	@cd $(INFRA_DIR)/ansible && ansible all -m ping

ansible-deploy: ## Run Ansible deployment playbook
	@echo "$(BLUE)Running Ansible deployment...$(NC)"
	@cd $(INFRA_DIR)/ansible && ansible-playbook -i inventory site.yml

## Development and Testing
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	@for dir in $(PROJECTS_DIR)/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "$(YELLOW)Testing $$dir$(NC)"; \
			cd "$$dir" && go test ./... && cd - >/dev/null; \
		fi; \
	done

build: ## Build all Go projects
	@echo "$(BLUE)Building all Go projects...$(NC)"
	@for dir in $(PROJECTS_DIR)/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "$(YELLOW)Building $$dir$(NC)"; \
			cd "$$dir" && go build ./... && cd - >/dev/null; \
		fi; \
	done

lint: ## Run linting on all projects
	@echo "$(BLUE)Running linters...$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 && \
		for dir in $(PROJECTS_DIR)/*/; do \
			if [ -f "$$dir/go.mod" ]; then \
				echo "$(YELLOW)Linting $$dir$(NC)"; \
				cd "$$dir" && golangci-lint run && cd - >/dev/null; \
			fi; \
		done || echo "$(YELLOW)golangci-lint not installed, skipping Go linting$(NC)"

## Monitoring and Status
status: ## Show status of all services
	@echo "$(BLUE)Checking service status...$(NC)"
	@docker-compose ps 2>/dev/null || echo "$(YELLOW)No default compose services running$(NC)"
	@echo "\n$(BLUE)Docker containers:$(NC)"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

logs: ## Show logs from all services
	@echo "$(BLUE)Showing service logs...$(NC)"
	@docker-compose logs -f --tail=50

logs-agent: ## Show Docker log agent logs
	@echo "$(BLUE)Showing Docker Log Agent logs...$(NC)"
	@cd $(PROJECTS_DIR)/docker-log-agent && docker-compose logs -f

health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "(healthy|unhealthy)" || echo "$(YELLOW)No health checks configured$(NC)"

## Cleanup and Maintenance
clean: ## Clean up containers and volumes
	@echo "$(BLUE)Cleaning up containers and volumes...$(NC)"
	@docker-compose down -v --remove-orphans 2>/dev/null || true
	@cd $(PROJECTS_DIR)/docker-log-agent && docker-compose down -v 2>/dev/null || true
	@cd $(PROJECTS_DIR)/event-driven && docker-compose down -v 2>/dev/null || true
	@echo "$(GREEN)Cleanup complete$(NC)"

clean-all: ## Clean up everything including images
	@echo "$(RED)Cleaning up everything including images...$(NC)"
	@make clean
	@docker system prune -af
	@echo "$(GREEN)Complete cleanup finished$(NC)"

reset: ## Reset entire environment
	@echo "$(RED)Resetting entire environment...$(NC)"
	@make clean-all
	@make setup
	@echo "$(GREEN)Environment reset complete$(NC)"

## Documentation and Help
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@echo "$(YELLOW)Documentation generation not implemented yet$(NC)"

validate: ## Validate all configuration files
	@echo "$(BLUE)Validating configuration files...$(NC)"
	@for file in $(shell find . -name "docker-compose*.yml" -o -name "*.yaml" | grep -v .git); do \
		echo "$(YELLOW)Validating $$file$(NC)"; \
		docker-compose -f "$$file" config >/dev/null 2>&1 || echo "$(RED)Invalid: $$file$(NC)"; \
	done

## Interactive Commands
interactive: ## Interactive compose file selection
	@echo "$(BLUE)Available compose files:$(NC)"
	@find . -name "docker-compose*.yml" -type f | sed 's|^\./||'
	@echo "$(YELLOW)Enter the compose file path:$(NC)"
	@read -p "File: " name; \
	if [ -f "$$name" ]; then \
		docker-compose -f "$$name" up -d; \
	else \
		echo "$(RED)File not found: $$name$(NC)"; \
	fi

## Help and Information
help: ## Show this help message
	@echo "$(GREEN)DevOps Projects - Available Commands$(NC)"
	@echo ""
	@echo "$(BLUE)Setup & Installation:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(setup|verify|install)"
	@echo ""
	@echo "$(BLUE)Environment Management:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(up-|down|start|stop)"
	@echo ""
	@echo "$(BLUE)Project Commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(docker-log|event-driven|scl-|build|test)"
	@echo ""
	@echo "$(BLUE)Infrastructure:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(infra-|pulumi-|ansible-)"
	@echo ""
	@echo "$(BLUE)Monitoring & Maintenance:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(status|logs|health|clean|validate)"
	@echo ""
	@echo "$(BLUE)Other Commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -vE "(setup|verify|install|up-|down|start|stop|docker-log|event-driven|scl-|build|test|infra-|pulumi-|ansible-|status|logs|health|clean|validate)"
	@echo ""
	@echo "$(GREEN)Example Usage:$(NC)"
	@echo "  make setup           # Initial environment setup"
	@echo "  make up-dev          # Start development environment"
	@echo "  make docker-log-agent # Deploy log monitoring"
	@echo "  make status          # Check all services"
	@echo "  make clean           # Clean up resources"