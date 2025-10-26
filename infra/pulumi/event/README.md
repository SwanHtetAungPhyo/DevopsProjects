# Event-Driven AWS Infrastructure

A Pulumi Go program that provisions a complete AWS infrastructure setup including VPC networking and EC2 instances. This project demonstrates Infrastructure as Code best practices using modular components for scalable cloud deployments.

## Architecture Overview

This infrastructure creates:
- **VPC** with custom CIDR block and DNS support
- **Public subnets** across multiple availability zones
- **Internet Gateway** for public internet access
- **Route tables** with proper routing configuration
- **EC2 instance** with security groups and user data
- **Security groups** with HTTP, HTTPS, and SSH access

```mermaid
graph TB
    subgraph "AWS Region (us-east-1)"
        subgraph "VPC (10.0.0.0/16)"
            IGW[Internet Gateway]
            RT[Route Table<br/>0.0.0.0/0 → IGW]
            
            subgraph "AZ: us-east-1a"
                SUBNET1[Public Subnet 1<br/>10.0.1.0/24]
            end
            
            subgraph "AZ: us-east-1b"
                SUBNET2[Public Subnet 2<br/>10.0.2.0/24]
            end
            
            subgraph "Security Group"
                SG[Security Group Rules<br/>• SSH (22)<br/>• HTTP (80)<br/>• HTTPS (443)<br/>• All Outbound]
            end
            
            EC2[EC2 Instance<br/>t2.micro<br/>Amazon Linux 2<br/>Public IP]
        end
    end
    
    INTERNET((Internet))
    
    INTERNET -.-> IGW
    IGW --> RT
    RT --> SUBNET1
    RT --> SUBNET2
    SUBNET1 --> EC2
    SG --> EC2
    
    style VPC fill:#e1f5fe
    style SUBNET1 fill:#f3e5f5
    style SUBNET2 fill:#f3e5f5
    style EC2 fill:#fff3e0
    style SG fill:#e8f5e8
    style IGW fill:#fce4ec
```

## Project Structure

```
.
├── Pulumi.yaml         # Pulumi project configuration
├── go.mod             # Go module dependencies
├── go.sum             # Go module checksums
├── main.go            # Main Pulumi program
├── .env               # Environment variables (AWS credentials)
├── modules/
│   ├── vpc.go         # VPC module with networking components
│   └── ec2.go         # EC2 module with instance and security groups
└── README.md          # This file
```

## Features

### VPC Module
- Creates VPC with configurable CIDR block
- Provisions multiple public subnets across availability zones
- Sets up Internet Gateway for public access
- Configures route tables with internet routing
- Enables DNS hostnames and support

### EC2 Module
- Launches EC2 instances with latest Amazon Linux 2 AMI
- Creates security groups with HTTP (80), HTTPS (443), and SSH (22) access
- Supports custom user data scripts for instance initialization
- Configures public IP assignment
- Allows SSH key pair association

## Prerequisites

- Go 1.23 or later
- Pulumi CLI installed and authenticated
- AWS account with appropriate permissions
- AWS credentials configured (via AWS CLI, environment variables, or IAM roles)

## Configuration

### Environment Variables

Create a `.env` file with your AWS credentials:

```bash
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
```

### Pulumi Configuration

Set the required AWS configuration:

```bash
pulumi config set aws:region us-east-1
pulumi config set aws:profile your-aws-profile
```

## Getting Started

1. Clone and navigate to the project directory

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Configure your AWS credentials in the `.env` file

4. Set Pulumi configuration:
   ```bash
   pulumi config set aws:region us-east-1
   pulumi config set aws:profile default
   ```

5. Preview the infrastructure changes:
   ```bash
   pulumi preview
   ```

6. Deploy the infrastructure:
   ```bash
   pulumi up
   ```

## Outputs

The program exports the following outputs:

- `vpcId`: The ID of the created VPC
- `publicSubnetIds`: Array of public subnet IDs
- `instanceId`: The ID of the EC2 instance
- `instancePublicIp`: Public IP address of the EC2 instance
- `instancePublicDns`: Public DNS name of the EC2 instance

View outputs with:
```bash
pulumi stack output
```

## Customization

### VPC Configuration

Modify the VPC settings in `main.go`:

```go
network, err := modules.NewVPC(ctx, "production-vpc", modules.VPCArgs{
    CidrBlock: "10.0.0.0/16",
    PublicSubnetCidr: []string{
        "10.0.1.0/24",
        "10.0.2.0/24",
    },
    AvailabilityZone: []string{
        "us-east-1a",
        "us-east-1b",
    },
    Tags: map[string]string{
        "Environment": "production",
        "Project":     "your-project",
    },
})
```

### EC2 Configuration

Customize the EC2 instance settings:

```go
instance, err := modules.NewEC2(ctx, "web-server", modules.EC2Args{
    InstanceType: "t3.micro",
    KeyName:      "your-key-pair",
    UserData:     string(userData),
    Tags: map[string]string{
        "Name":        "web-server",
        "Environment": "production",
    },
})
```

## Dependencies

Key dependencies include:
- `github.com/pulumi/pulumi-aws/sdk/v7` - AWS provider for Pulumi
- `github.com/pulumi/pulumi/sdk/v3` - Core Pulumi SDK
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/rs/zerolog` - Structured logging

## Security Considerations

- Store AWS credentials securely using environment variables or IAM roles
- Review security group rules and restrict access as needed
- Use appropriate instance types for your workload
- Consider enabling VPC Flow Logs for network monitoring
- Implement proper tagging strategy for resource management

## Cleanup

To destroy the infrastructure:

```bash
pulumi destroy
```

## Support

For issues and questions:
- Check the Pulumi documentation: https://www.pulumi.com/docs/
- AWS provider reference: https://www.pulumi.com/registry/packages/aws/
- Pulumi Community Slack: https://slack.pulumi.com/