# Troubleshooting Guide

This guide helps you resolve common issues encountered when working with the DevOps projects in this repository.

## Quick Diagnostics

### System Health Check
```bash
# Run comprehensive system check
make verify

# Check Docker status
docker --version
docker-compose --version
docker system info

# Check Go installation
go version
go env GOPATH

# Check available resources
df -h
free -h
docker system df
```

## Common Issues and Solutions

### Docker Issues

#### Docker Permission Denied
**Problem**: `permission denied while trying to connect to the Docker daemon socket`

**Solution**:
```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Restart Docker service
sudo systemctl restart docker

# Alternative: Use sudo for docker commands
sudo docker ps
```

#### Docker Compose File Not Found
**Problem**: `docker-compose.yml not found`

**Solution**:
```bash
# Check current directory
pwd
ls -la

# Navigate to correct project directory
cd projects/docker-log-agent
docker-compose up -d

# Or specify file explicitly
docker-compose -f projects/docker-log-agent/docker-compose.yml up -d
```

#### Container Port Conflicts
**Problem**: `port is already allocated`

**Solution**:
```bash
# Find what's using the port
lsof -i :8080
netstat -tulpn | grep :8080

# Kill the process
sudo kill -9 $(lsof -t -i:8080)

# Or change port in docker-compose.yml
ports:
  - "8081:8080"  # Use different host port
```

#### Docker Build Failures
**Problem**: Build context issues or dependency failures

**Solution**:
```bash
# Clean Docker cache
docker system prune -af
docker builder prune -af

# Rebuild without cache
docker-compose build --no-cache

# Check Dockerfile syntax
docker build --dry-run .
```

### Go Development Issues

#### Go Module Issues
**Problem**: `go: module not found` or dependency issues

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Reinitialize modules
go mod tidy
go mod download

# Verify module integrity
go mod verify

# Update dependencies
go get -u ./...
```

#### Go Build Failures
**Problem**: Compilation errors or missing dependencies

**Solution**:
```bash
# Check Go environment
go env

# Verify GOPATH and GOROOT
echo $GOPATH
echo $GOROOT

# Build with verbose output
go build -v ./...

# Check for syntax errors
go vet ./...
gofmt -d .
```

#### Cross-Platform Build Issues
**Problem**: Building for different architectures

**Solution**:
```bash
# Build for Linux (from macOS/Windows)
GOOS=linux GOARCH=amd64 go build -o app-linux ./cmd/app

# Build for multiple platforms
make build-all

# Use Docker for consistent builds
docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.21 go build -v
```

### Infrastructure Issues

#### AWS Credentials
**Problem**: AWS authentication failures

**Solution**:
```bash
# Configure AWS CLI
aws configure

# Check current credentials
aws sts get-caller-identity

# Set environment variables
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1

# Use AWS profiles
aws configure --profile devops
export AWS_PROFILE=devops
```

#### Terraform Issues
**Problem**: Terraform state or provider issues

**Solution**:
```bash
# Initialize Terraform
cd infra/terraform
terraform init

# Refresh state
terraform refresh

# Fix state issues
terraform state list
terraform state show aws_instance.example

# Force unlock (if locked)
terraform force-unlock LOCK_ID

# Validate configuration
terraform validate
terraform plan
```

#### Pulumi Issues
**Problem**: Pulumi stack or authentication issues

**Solution**:
```bash
# Login to Pulumi
pulumi login

# Select correct stack
pulumi stack select dev

# Refresh stack state
pulumi refresh

# Check configuration
pulumi config
pulumi config set aws:region us-east-1

# Preview changes
pulumi preview --diff
```

### Network and Connectivity Issues

#### SSH Connection Failures
**Problem**: Cannot connect to remote servers

**Solution**:
```bash
# Test SSH connectivity
ssh -v user@server.com

# Check SSH key
ssh-add -l
ssh-keygen -t rsa -b 4096

# Copy SSH key to server
ssh-copy-id user@server.com

# Use specific key
ssh -i ~/.ssh/specific_key user@server.com
```

#### DNS Resolution Issues
**Problem**: Cannot resolve hostnames

**Solution**:
```bash
# Test DNS resolution
nslookup google.com
dig google.com

# Check /etc/hosts
cat /etc/hosts

# Flush DNS cache (macOS)
sudo dscacheutil -flushcache

# Flush DNS cache (Linux)
sudo systemctl restart systemd-resolved
```

#### Firewall Issues
**Problem**: Blocked ports or connections

**Solution**:
```bash
# Check firewall status (Ubuntu)
sudo ufw status

# Allow specific ports
sudo ufw allow 8080/tcp
sudo ufw allow ssh

# Check iptables rules
sudo iptables -L

# Disable firewall temporarily (for testing)
sudo ufw disable
```

### Application-Specific Issues

#### Docker Log Agent Issues

**Problem**: Agent not detecting container logs

**Solution**:
```bash
# Check Docker socket permissions
ls -la /var/run/docker.sock

# Verify container filters
docker ps --format "table {{.Names}}\t{{.Labels}}"

# Check agent logs
docker logs docker-log-agent

# Test webhook connectivity
curl -X POST -H "Content-Type: application/json" \
  -d '{"test": "message"}' \
  $WEBHOOK_URL
```

**Problem**: High memory usage

**Solution**:
```bash
# Check memory usage
docker stats docker-log-agent

# Reduce buffer size in config
ALERT_AGENT_BUFFER_SIZE=500

# Limit container memory
docker run --memory=100m docker-log-agent
```

#### Event-Driven Services Issues

**Problem**: Services cannot communicate

**Solution**:
```bash
# Check service status
docker-compose ps

# Verify network connectivity
docker network ls
docker network inspect event-driven_default

# Test NATS connectivity
docker exec -it nats nats-cli pub test "hello"
docker exec -it nats nats-cli sub test

# Check service logs
docker-compose logs app-one
docker-compose logs nats
```

#### SCL Language Issues

**Problem**: SCL compilation errors

**Solution**:
```bash
# Check ANTLR grammar
cd projects/scl/lang
antlr4 -Dlanguage=Go InfraDSL.g4

# Rebuild parser
go generate ./...

# Test with simple script
echo 'fn main() { print("test"); }' > test.scl
./scl test.scl

# Enable verbose output
./scl --verbose examples/demo.scl
```

### Performance Issues

#### Slow Docker Builds
**Problem**: Docker builds taking too long

**Solution**:
```bash
# Use BuildKit
export DOCKER_BUILDKIT=1

# Optimize Dockerfile
# - Use multi-stage builds
# - Minimize layers
# - Use .dockerignore

# Use build cache
docker build --cache-from myapp:latest .

# Parallel builds
docker-compose build --parallel
```

#### High Resource Usage
**Problem**: Containers using too much CPU/memory

**Solution**:
```bash
# Monitor resource usage
docker stats

# Set resource limits
docker run --cpus="0.5" --memory="512m" myapp

# In docker-compose.yml:
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

### CI/CD Issues

#### GitHub Actions Failures
**Problem**: CI/CD pipeline failures

**Solution**:
```bash
# Check workflow syntax
# Use GitHub Actions extension in VS Code

# Validate YAML
yamllint .github/workflows/ci.yml

# Test locally with act
act -j test

# Check secrets and environment variables
# Ensure all required secrets are set in GitHub
```

#### Docker Registry Issues
**Problem**: Cannot push/pull images

**Solution**:
```bash
# Login to registry
docker login

# Check image tags
docker images

# Retry with explicit registry
docker tag myapp:latest docker.io/username/myapp:latest
docker push docker.io/username/myapp:latest
```

## Environment-Specific Issues

### macOS Issues

#### Docker Desktop Problems
```bash
# Restart Docker Desktop
# Use Docker Desktop GUI or:
killall Docker && open /Applications/Docker.app

# Reset Docker Desktop
# Docker Desktop > Troubleshoot > Reset to factory defaults

# Check Docker Desktop resources
# Docker Desktop > Preferences > Resources
```

#### Homebrew Issues
```bash
# Update Homebrew
brew update
brew upgrade

# Fix permissions
sudo chown -R $(whoami) $(brew --prefix)/*

# Reinstall problematic packages
brew uninstall docker
brew install docker
```

### Linux Issues

#### SystemD Service Issues
```bash
# Check service status
sudo systemctl status docker

# Restart services
sudo systemctl restart docker
sudo systemctl enable docker

# Check logs
sudo journalctl -u docker.service
```

#### Package Manager Issues
```bash
# Update package lists
sudo apt update

# Fix broken packages
sudo apt --fix-broken install

# Clean package cache
sudo apt clean
sudo apt autoremove
```

## Debugging Techniques

### Enable Debug Logging
```bash
# Docker Log Agent
export LOG_LEVEL=debug

# Docker daemon
sudo dockerd --debug

# Go applications
export DEBUG=true
export VERBOSE=true
```

### Network Debugging
```bash
# Test connectivity
ping google.com
telnet github.com 443
curl -v https://api.github.com

# Check routing
traceroute google.com
ip route show

# Monitor network traffic
sudo tcpdump -i any port 80
```

### Container Debugging
```bash
# Inspect container
docker inspect container_name

# Execute shell in container
docker exec -it container_name /bin/sh

# Check container logs
docker logs --follow container_name

# Monitor container events
docker events --filter container=container_name
```

## Getting Help

### Log Collection
When reporting issues, include:

```bash
# System information
uname -a
docker version
docker-compose version
go version

# Error logs
docker logs container_name > error.log
journalctl -u docker.service > docker.log

# Configuration
cat .env
docker-compose config
```

### Useful Commands for Diagnostics
```bash
# Complete system check
make verify 2>&1 | tee system-check.log

# Docker system information
docker system info > docker-info.txt
docker system df > docker-usage.txt

# Network information
ip addr show > network-config.txt
netstat -tulpn > network-ports.txt
```

### Community Resources
- **Docker Documentation**: https://docs.docker.com/
- **Go Documentation**: https://golang.org/doc/
- **AWS Documentation**: https://docs.aws.amazon.com/
- **Terraform Documentation**: https://www.terraform.io/docs/
- **Pulumi Documentation**: https://www.pulumi.com/docs/

### Emergency Recovery
```bash
# Stop all containers
docker stop $(docker ps -aq)

# Remove all containers
docker rm $(docker ps -aq)

# Clean everything
make clean-all

# Reset environment
make reset
```

---

**Note**: Always backup important data before performing destructive operations. Test solutions in a development environment first.