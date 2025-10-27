# SCL - Infrastructure Configuration Language

An **Ansible-like** configuration management tool that can either generate bash scripts or execute commands directly on remote servers via SSH.

## 🚀 Quick Start

### 1. Build SCL
```bash
go build -o scl
```

### 2. Setup Docker Test Environment
```bash
./setup-docker.sh
```

### 3. Test Your Language
```bash
# Compile mode - generates bash script
./scl examples/demo.scl

# Interpret mode - direct execution like Ansible  
sed 's/mode := compile/mode := interpret/g' examples/demo.scl > /tmp/test.scl
./scl /tmp/test.scl --verbose
```

## 📋 Two Execution Modes

### **Compile Mode** (`mode := compile`)
- Generates bash scripts for later execution
- Output saved to `output.sh`
- Run with: `bash output.sh`

### **Interpret Mode** (`mode := interpret`) 
- Executes commands directly via SSH (like Ansible)
- Automatic SSH connectivity checking
- Real-time feedback and error reporting

## 📝 Language Syntax

```scl
import primary;

setting := configuration;
target := ["server.example.com"];
super_user = true;
on_error = rollback;
mode := interpret;  // or compile

declare tools_ready: bool = check(target, ["docker", "nginx"]);
declare app_dir: string = "/opt/myapp";

fn deploy_app(){
    if tools_ready {
        primary.print("Server is ready");
        
        // Copy files to remote server
        copy("config/app.conf", "/etc/myapp/app.conf");
        
        // Create files with permissions
        create("/var/log/myapp", "access.log", "644");
        create("/opt/myapp/scripts", "start.sh", "755");
        
    } else {
        primary.print("Installing tools...");
        install(primary.pkg_snap, "docker");
        install(primary.pkg_snap, "nginx");
    }
}

fn main(){
    deploy_app();
    primary.print("Deployment complete!");
}
```

## 🛠️ Available Functions

| Function | Purpose | Example |
|----------|---------|---------|
| `check(target, ["tools"])` | Verify tools exist on server | `check(target, ["docker", "nginx"])` |
| `copy(source, dest)` | Copy files to remote server | `copy("app.conf", "/etc/app.conf")` |
| `create(dir, file, mode)` | Create files with permissions | `create("/tmp", "script.sh", "755")` |
| `install(pkg_mgr, package)` | Install packages | `install(primary.pkg_snap, "docker")` |
| `test(description)` | Test connectivity | `test("server_check")` |
| `primary.print(msg)` | Display messages | `primary.print("Starting...")` |

## 🧪 Testing

### Prerequisites
- Go 1.19+
- Docker (for SSH testing)

### Run Tests
```bash
# 1. Build SCL
go build -o scl

# 2. Setup Docker SSH environment
./setup-docker.sh

# 3. Test compile mode
./scl examples/demo.scl
cat output.sh  # Check generated bash

# 4. Test interpret mode (change mode in demo.scl first)
./scl examples/demo.scl --verbose

# 5. Verify files on remote server
ssh -p 2222 testuser@localhost "ls -la /tmp/myapp/ /tmp/scripts/"

# 6. Cleanup
docker-compose -f docker-test/docker-compose.yml down
```

## 📁 Project Structure

```
SCL/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
├── internal/
│   ├── nodes/             # Code generation & execution
│   ├── parser/            # Generated ANTLR parser
│   ├── utils/             # Utilities
│   └── validation/        # Field validation
├── grammar/InfraDSL.g4    # Language grammar
├── examples/demo.scl      # Demo SCL file
├── test-files/            # Sample config files
├── docker-test/           # Docker SSH test environment
└── setup-docker.sh       # Docker setup script
```

## 🔧 How It Works

1. **Parse** SCL files using ANTLR grammar
2. **Detect** execution mode (`compile` or `interpret`)
3. **Validate** configuration and SSH connectivity
4. **Execute** either:
   - Generate bash scripts (compile mode)
   - Run commands directly via SSH (interpret mode)

## 🎯 Use Cases

- **Development**: Use interpret mode for quick testing
- **Production**: Use compile mode to generate reviewed scripts
- **CI/CD**: Generate scripts for automated deployments
- **Configuration Management**: Ansible-like server configuration

## 🚨 Error Handling

SCL provides detailed error messages with troubleshooting tips:

```
✗ SSH connectivity check failed: Connection refused
Troubleshooting:
  • Check if SSH server is running on target
  • Verify SSH key authentication is set up
  • Try: ssh-copy-id user@server
  • For Docker testing: Make sure container is running
```

## 📖 Examples

The `examples/demo.scl` file demonstrates:
- SSH connectivity checking
- File copying with automatic directory creation
- File creation with proper permissions
- Conditional logic based on system checks
- Both compile and interpret mode compatibility

Your SCL language is ready for infrastructure configuration management! 🎉

## 🔧 Enhanced SSH DevOps Functionality

SCL now includes comprehensive SSH-based DevOps operations for professional infrastructure management:

### System Operations
- **System Information**: `sysinfo()` - Get comprehensive system details
- **Monitoring**: `monitor()` - Collect real-time system metrics
- **Performance Analysis**: Built-in performance monitoring and analysis

### Package & Service Management
- **Package Management**: `package("install", "nginx")` - Auto-detects package manager (apt/yum/dnf)
- **Service Control**: `service("restart", "nginx")` - Systemd service management
- **Process Management**: Start, stop, restart, enable, disable services

### Container & Application Management
- **Docker Operations**: `docker("ps")`, `docker("start", "container")` - Full Docker support
- **Web Server Management**: Nginx, Apache configuration and control
- **Database Operations**: MySQL, PostgreSQL, MongoDB support

### Security & Compliance
- **Security Audit**: `audit()` - Comprehensive security checks
- **Firewall Management**: `firewall("allow", "80/tcp")` - UFW firewall control
- **Certificate Management**: SSL/TLS certificate operations
- **User Management**: User and group operations with sudo support

### File & Network Operations
- **Enhanced File Operations**: SFTP-based upload/download with proper permissions
- **Network Diagnostics**: `nettest("google.com", "ping")` - Connectivity testing
- **Backup Operations**: Automated backup and restore with compression
- **Log Analysis**: Real-time log monitoring and pattern matching

### System Administration
- **Cron Management**: Schedule and manage cron jobs
- **System Tuning**: Performance optimization and kernel parameter tuning
- **Compliance Checks**: CIS benchmark and PCI DSS compliance validation

## 🐳 Docker Testing Environment

Set up a complete SSH testing environment with Docker:

```bash
# Quick setup
./quick-docker-setup.sh

# Manual setup
docker-compose -f docker-test/docker-compose.yml up -d --build

# Test SSH connection
ssh -p 2222 testuser@localhost  # Password: testpass
```

## 📚 Example Configurations

### Local System Information
```scl
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";

fn main() {
    print("🖥️ Local System Information");
    sysinfo();  // Works locally without SSH
    test();
    create("/tmp", "scl_test.txt", "644");
}
```

### Remote DevOps Management
```scl
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "testuser@localhost:2222";  // Docker test container

fn main() {
    print("🚀 DevOps Infrastructure Management");
    
    // System analysis
    sysinfo();
    monitor();
    
    // Package management
    package("install", "htop");
    package("update", "");
    
    // Service management
    service("status", "ssh");
    service("restart", "nginx");
    
    // Security operations
    firewall("status");
    firewall("allow", "80/tcp");
    audit();
    
    // File operations
    copy("config.txt", "/etc/myapp/config.txt");
    backup("backup", "/etc", "/backup/etc.tar.gz");
    
    // Network diagnostics
    nettest("google.com", "ping");
    nettest("github.com:443", "telnet");
    
    // Log analysis
    logs("/var/log/syslog", "error", 10);
}
```

## 🔑 SSH Authentication

SCL supports multiple SSH authentication methods:

1. **SSH Agent** (recommended)
2. **SSH Keys** with optional passphrase
3. **Password authentication** (for testing)

The system automatically tries SSH agent first, then falls back to key-based authentication.

## 🛠️ Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `sysinfo()` | System information | `sysinfo();` |
| `monitor()` | System metrics | `monitor();` |
| `package(action, name)` | Package management | `package("install", "nginx");` |
| `service(action, name)` | Service control | `service("restart", "nginx");` |
| `docker(action, args...)` | Docker operations | `docker("ps");` |
| `firewall(action, rule)` | Firewall management | `firewall("allow", "80/tcp");` |
| `user(action, username)` | User management | `user("add", "devops");` |
| `cert(action, domain)` | Certificate management | `cert("check_expiry", "example.com");` |
| `cron(action, schedule, cmd)` | Cron management | `cron("add", "0 2 * * *", "/backup.sh");` |
| `nettest(target, type)` | Network diagnostics | `nettest("google.com", "ping");` |
| `logs(file, pattern, lines)` | Log analysis | `logs("/var/log/syslog", "error", 10);` |
| `backup(op, src, dst)` | Backup operations | `backup("backup", "/etc", "/backup.tar.gz");` |
| `audit()` | Security audit | `audit();` |
| `tune(param, value)` | System tuning | `tune("swappiness", "10");` |
| `webserver(type, action)` | Web server control | `webserver("nginx", "reload");` |
| `database(type, action)` | Database operations | `database("mysql", "status");` |

## 🚦 Error Handling

SCL provides comprehensive error handling with detailed troubleshooting information:

- SSH connection failures with specific remediation steps
- Command execution errors with full output
- Network connectivity issues with diagnostic suggestions
- Permission problems with security recommendations

## 📖 Documentation

### Complete Language Documentation
- **[Language Reference](docs/LANGUAGE_REFERENCE.md)** - Complete syntax, tokens, and feature reference
- **[Syntax Guide](docs/SYNTAX_GUIDE.md)** - Comprehensive grammar rules and token definitions
- **[Tutorial](docs/TUTORIAL.md)** - Step-by-step learning guide from basics to advanced
- **[SSH Executor Guide](docs/SSH_EXECUTOR_GUIDE.md)** - SSH functionality and DevOps operations

### Examples and Code Samples
- **[Complete Examples](examples/complete_examples.scl)** - All language features with practical examples
- **[Local Demo](examples/local_demo.scl)** - Local system operations
- **[SSH Demo](examples/ssh_demo_simple.scl)** - Remote SSH operations
- **[DevOps Demo](examples/devops_demo.scl)** - Comprehensive DevOps deployment
- **[Go API Examples](examples/ssh_devops_example.go)** - Direct API usage examples

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.