# Environment Setup Guide

This guide will help you set up your development environment for working with the DevOps projects in this repository.

## Prerequisites

### Required Software

#### Docker & Container Tools
```bash
# macOS (using Homebrew)
brew install docker docker-compose

# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose

# CentOS/RHEL
sudo yum install docker docker-compose
```

#### Programming Languages
```bash
# Go (required for custom tools)
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# Or download from https://golang.org/dl/

# Python (for automation scripts)
# macOS
brew install python3

# Ubuntu/Debian
sudo apt install python3 python3-pip
```

#### Infrastructure Tools
```bash
# Terraform
# macOS
brew install terraform

# Ubuntu/Debian
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform

# Pulumi
curl -fsSL https://get.pulumi.com | sh

# Ansible
pip3 install ansible
```

#### Cloud CLI Tools
```bash
# AWS CLI
# macOS
brew install awscli

# Ubuntu/Debian
sudo apt install awscli

# Or using pip
pip3 install awscli
```

### Optional Tools

#### Development Tools
```bash
# Git (if not already installed)
# macOS
brew install git

# Ubuntu/Debian
sudo apt install git

# Make (for Makefile commands)
# Usually pre-installed on macOS and Linux

# golangci-lint (for Go code linting)
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
```

## Repository Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd DevopsProjects
```

### 2. Environment Configuration
```bash
# Create environment file
make setup

# Or manually create .env file
cp .env.example .env  # if available
# OR
echo "# Add your environment variables here" > .env
```

### 3. Configure Environment Variables
Edit the `.env` file with your specific configuration:

```bash
# AWS Configuration
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_DEFAULT_REGION=us-east-1

# Docker Configuration
DOCKER_HOST=unix:///var/run/docker.sock

# Application Configuration
LOG_LEVEL=info
ENVIRONMENT=development

# Webhook URLs (for Docker Log Agent)
WEBHOOK_URL=https://your-webhook-url.com

# Database Configuration (if applicable)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=devops_db
DB_USER=devops_user
DB_PASSWORD=secure_password
```

### 4. Verify Installation
```bash
# Verify all dependencies
make verify

# Check Docker installation
docker --version
docker-compose --version

# Check Go installation
go version

# Check cloud tools
aws --version
terraform --version
pulumi version
ansible --version
```

## Project-Specific Setup

### Docker Log Agent
```bash
# Navigate to project
cd projects/docker-log-agent

# Build the application
go build -o main .

# Test with Docker Compose
docker-compose up -d
```

### Event-Driven Services
```bash
# Navigate to project
cd projects/event-driven

# Start all services
docker-compose up -d

# Check service status
docker-compose ps
```

### SCL Configuration Language
```bash
# Navigate to project
cd projects/scl/lang

# Build SCL tool
go build -o scl

# Set up Docker test environment
./setup-docker.sh

# Test the tool
./scl examples/demo.scl
```

### Pulumi Infrastructure
```bash
# Navigate to project
cd infra/pulumi/event

# Install dependencies
go mod tidy

# Configure Pulumi
pulumi config set aws:region us-east-1
pulumi config set aws:profile default

# Preview infrastructure
pulumi preview
```

## IDE and Editor Setup

### Visual Studio Code
Recommended extensions:
```json
{
  "recommendations": [
    "ms-vscode.go",
    "hashicorp.terraform",
    "ms-python.python",
    "redhat.ansible",
    "ms-azuretools.vscode-docker",
    "pulumi.pulumi-lsp"
  ]
}
```

### GoLand/IntelliJ IDEA
- Install Go plugin
- Configure Go SDK path
- Set up Docker integration
- Install Terraform plugin

## Network and Security Configuration

### Docker Network Setup
```bash
# Create custom network for projects
docker network create devops-network

# Verify network creation
docker network ls
```

### SSH Key Setup (for remote deployments)
```bash
# Generate SSH key if needed
ssh-keygen -t rsa -b 4096 -C "your-email@example.com"

# Add to SSH agent
ssh-add ~/.ssh/id_rsa

# Copy public key for server access
ssh-copy-id user@your-server.com
```

### Firewall Configuration
```bash
# Allow Docker ports (if using UFW)
sudo ufw allow 2376/tcp  # Docker daemon
sudo ufw allow 8080/tcp  # Application ports
sudo ufw allow 3000/tcp
sudo ufw allow 5432/tcp  # Database ports
```

## Troubleshooting Common Issues

### Docker Permission Issues
```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Or use sudo for docker commands
sudo docker ps
```

### Go Module Issues
```bash
# Clean Go module cache
go clean -modcache

# Reinitialize modules
go mod tidy
go mod download
```

### Port Conflicts
```bash
# Check what's using a port
lsof -i :8080

# Kill process using port
sudo kill -9 $(lsof -t -i:8080)
```

### AWS Credentials
```bash
# Configure AWS CLI
aws configure

# Or set environment variables
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1
```

## Validation Commands

After setup, run these commands to validate your environment:

```bash
# Repository-level validation
make verify
make validate

# Test Docker setup
docker run hello-world

# Test Go compilation
cd projects/docker-log-agent && go build .

# Test infrastructure tools
terraform version
pulumi version
ansible --version

# Test cloud connectivity
aws sts get-caller-identity
```

## Next Steps

1. **Explore Projects**: Start with the [Docker Log Agent](../projects/docker-log-agent/README.md)
2. **Read Documentation**: Review project-specific README files
3. **Run Examples**: Execute the provided examples and demos
4. **Customize Configuration**: Adapt settings for your environment
5. **Deploy Infrastructure**: Try deploying to your cloud environment

## Getting Help

If you encounter issues during setup:

1. Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
2. Review project-specific documentation
3. Verify all prerequisites are installed
4. Check environment variable configuration
5. Ensure proper permissions and network access

---

**Note**: This setup guide is tested on macOS and Linux environments. Windows users should use WSL2 or adapt commands accordingly.