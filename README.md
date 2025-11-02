
# DevOps Engineering Portfolio

**Author:** Swan Htet Aung Phyo  
**Focus:** Infrastructure as Code, Container Orchestration, CI/CD, and Cloud-Native Solutions

A comprehensive collection of DevOps projects demonstrating expertise in modern infrastructure management, automation, and cloud technologies. This repository showcases practical implementations of industry-standard DevOps practices and tools.

## Repository Overview

This repository contains production-ready DevOps projects organized into distinct categories:

### Infrastructure as Code
- **Pulumi AWS Infrastructure** - Complete AWS cloud infrastructure provisioning
- **Terraform Configurations** - Multi-cloud infrastructure templates
- **Ansible Playbooks** - Configuration management and automation

### Container & Orchestration
- **Docker Log Agent** - Real-time container monitoring and alerting system
- **Event-Driven Architecture** - Microservices with container orchestration
- **Container Security** - Security scanning and compliance tools

### Custom Tools & Languages
- **SCL (System Configuration Language)** - Custom infrastructure configuration DSL
- **Automation Scripts** - Shell scripts for DevOps workflows

### CI/CD & Automation
- **GitHub Actions Workflows** - Automated testing and deployment pipelines
- **Docker Compose Configurations** - Multi-service application stacks
- **Monitoring & Alerting** - Comprehensive observability solutions

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+
- Python 3.9+
- AWS CLI configured
- Terraform 1.5+
- Ansible 2.14+

### Environment Setup
```bash
# Clone the repository
git clone <repository-url>
cd DevopsProjects

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Install dependencies
make setup

# Verify installation
make verify
```

### Available Commands
```bash
# Development environment
make up-dev              # Start development stack
make up-prod             # Start production stack
make status              # Check container status
make clean               # Clean up resources

# Project-specific commands
make docker-log-agent    # Deploy log monitoring
make event-driven        # Start event-driven services
make infrastructure      # Deploy cloud infrastructure
```

## Project Structure

```
DevopsProjects/
├── projects/                    # Individual DevOps projects
│   ├── docker-log-agent/       # Container monitoring solution
│   ├── event-driven/           # Microservices architecture
│   └── scl/                    # Custom configuration language
├── infra/                      # Infrastructure as Code
│   ├── terraform/              # Terraform configurations
│   ├── pulumi/                 # Pulumi projects
│   ├── ansible/                # Ansible playbooks
│   └── scripts/                # Automation scripts
├── docker/                     # Docker configurations
├── .github/                    # CI/CD workflows
├── scripts/                    # Utility scripts
└── docs/                       # Documentation
```

## Featured Projects

### 1. Docker Log Agent
**Technology Stack:** Go, Docker, Webhooks, YAML  
**Purpose:** Real-time container log monitoring with intelligent alerting

- Monitors Docker container logs for error patterns
- Webhook-based alert forwarding with retry logic
- Container filtering by name or labels
- Minimal resource footprint (<100MB memory)
- Health check endpoints and graceful shutdown

[View Project →](projects/docker-log-agent/)

### 2. Event-Driven Microservices
**Technology Stack:** Go, Docker Compose, NATS, Caddy  
**Purpose:** Scalable microservices architecture with event-driven communication

- Multiple microservices with inter-service communication
- Message queue integration with NATS
- Load balancing and reverse proxy with Caddy
- Container orchestration with Docker Compose

[View Project →](projects/event-driven/)

### 3. SCL Configuration Language
**Technology Stack:** Go, ANTLR, SSH, YAML  
**Purpose:** Custom DSL for infrastructure configuration management

- Ansible-like configuration management
- SSH-based remote execution
- Compile and interpret modes
- Comprehensive DevOps operations support

[View Project →](projects/scl/)

### 4. AWS Infrastructure with Pulumi
**Technology Stack:** Go, Pulumi, AWS, Terraform  
**Purpose:** Production-ready cloud infrastructure provisioning

- VPC with multi-AZ subnets
- EC2 instances with security groups
- Infrastructure as Code best practices
- Modular and reusable components

[View Project →](infra/pulumi/event/)

## Skills Demonstrated

### Infrastructure & Cloud
- Infrastructure as Code (Terraform, Pulumi)
- AWS cloud services and architecture
- Container orchestration and management
- Network security and configuration

### DevOps Practices
- CI/CD pipeline design and implementation
- Configuration management with Ansible
- Monitoring and observability solutions
- Security scanning and compliance

### Programming & Automation
- Go programming for system tools
- Shell scripting and automation
- Custom DSL development
- API design and integration

### Tools & Technologies
- **Containers:** Docker, Docker Compose
- **Cloud:** AWS (VPC, EC2, IAM, CloudWatch)
- **IaC:** Terraform, Pulumi, Ansible
- **Monitoring:** Custom log agents, webhooks
- **Languages:** Go, Python, Shell, YAML
- **CI/CD:** GitHub Actions, automated testing

## Getting Started with Individual Projects

Each project includes comprehensive documentation:

1. **Setup Instructions** - Environment requirements and installation
2. **Configuration Guide** - Customization and deployment options
3. **Usage Examples** - Practical implementation scenarios
4. **Troubleshooting** - Common issues and solutions

## Best Practices Implemented

- **Security First:** Secure defaults, credential management, least privilege
- **Scalability:** Modular design, horizontal scaling, resource optimization
- **Reliability:** Health checks, graceful shutdowns, error handling
- **Maintainability:** Clean code, comprehensive documentation, testing
- **Observability:** Logging, monitoring, alerting, and debugging tools

## Contributing

This repository follows industry-standard DevOps practices:

1. **Code Review:** All changes require review
2. **Testing:** Automated testing for all components
3. **Documentation:** Comprehensive docs for all projects
4. **Security:** Security scanning and vulnerability management

## Contact & Professional Links

**Swan Htet Aung Phyo**
- **Role:** DevOps Engineer
- **Specialization:** Cloud Infrastructure, Container Orchestration, Automation
- **Experience:** Production-scale infrastructure management and deployment

---

*This repository demonstrates practical DevOps engineering skills through real-world projects and industry-standard practices. Each project is designed to showcase different aspects of modern infrastructure management and automation.*