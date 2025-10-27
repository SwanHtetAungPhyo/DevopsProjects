# SCL Tutorial: From Basics to DevOps

## Table of Contents
1. [Getting Started](#getting-started)
2. [Basic Syntax](#basic-syntax)
3. [Variables and Data Types](#variables-and-data-types)
4. [Functions](#functions)
5. [Control Flow](#control-flow)
6. [Built-in Commands](#built-in-commands)
7. [Real-World Examples](#real-world-examples)
8. [Best Practices](#best-practices)

## Getting Started

### Installation
```bash
# Clone and build SCL
git clone <repository-url>
cd SCL
go build -o scl .
```

### Your First SCL Program
Create a file called `hello.scl`:

```scl
// hello.scl - Your first SCL program
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";

fn main() {
    print("Hello, SCL World!");
}
```

Run it:
```bash
./scl hello.scl --verbose
```

## Basic Syntax

### Comments
```scl
// Single line comment

/* Multi-line comment
   can span multiple lines */

// Comments can be used anywhere
mode := "interpret"; // End of line comment
```

### Statements and Semicolons
```scl
// Every statement ends with a semicolon
mode := "interpret";
super_user = true;
print("Hello World");
```

### Case Sensitivity
```scl
// SCL is case-sensitive
mode := "interpret";    // Correct
Mode := "interpret";    // Different variable
MODE := "interpret";    // Also different
```

## Variables and Data Types

### Variable Assignment
```scl
// Two types of assignment
mode := "interpret";        // Variable assignment (:=)
super_user = true;          // Value assignment (=)
```

### Basic Data Types
```scl
// Strings
hostname := "web-server";
config_path := "/etc/nginx/nginx.conf";

// Numbers
port := 8080;               // Integer
cpu_threshold := 85.5;      // Float

// Booleans
ssl_enabled := true;
debug_mode := false;

// Arrays
servers := ["web1", "web2", "web3"];
ports := [80, 443, 8080];
```

### Explicit Type Declaration
```scl
// Declare variables with explicit types
declare environment: string = "production";
declare max_connections: int = 1000;
declare ssl_enabled: bool = true;
declare server_list: list = ["web1", "web2"];
```

## Functions

### Function Declaration
```scl
// Basic function
fn greet() {
    print("Hello from function!");
}

// Main function (entry point)
fn main() {
    print("Starting program");
    greet();
    print("Program finished");
}
```

### Function Calls
```scl
fn main() {
    // Call built-in functions
    print("System information:");
    sysinfo();
    
    // Call user-defined functions
    setup_server();
    deploy_app();
}

fn setup_server() {
    print("Setting up server...");
    package("install", "nginx");
}

fn deploy_app() {
    print("Deploying application...");
    service("start", "nginx");
}
```

## Control Flow

### If Statements
```scl
// Simple if
if (ssl_enabled) {
    print("SSL is enabled");
}

// If-else
if (environment == "production") {
    print("Production environment");
} else {
    print("Development environment");
}
```

### Complex Conditions
```scl
// Comparison operators
if (port > 1024) {
    print("Using non-privileged port");
}

if (environment != "development") {
    print("Not in development mode");
}

// Logical operators
if (ssl_enabled && environment == "production") {
    print("Production SSL setup");
}

if (debug_mode || verbose_logging) {
    print("Detailed logging enabled");
}
```

### Nested If Statements
```scl
fn configure_security() {
    if (environment == "production") {
        if (ssl_enabled) {
            print("Setting up production SSL");
            cert("letsencrypt", "example.com");
        } else {
            print("Production without SSL");
            firewall("allow", "80/tcp");
        }
    } else {
        print("Development security setup");
        firewall("allow", "3000/tcp");
    }
}
```

## Built-in Commands

### System Information
```scl
fn system_check() {
    print("=== System Information ===");
    sysinfo();          // Get system details
    
    print("=== System Monitoring ===");
    monitor();          // Get performance metrics
}
```

### Package Management
```scl
fn install_packages() {
    print("Installing required packages...");
    
    // Update package lists
    package("update", "");
    
    // Install packages
    package("install", "nginx");
    package("install", "docker.io");
    package("install", "htop");
    
    print("Package installation completed");
}
```

### Service Management
```scl
fn manage_services() {
    print("Managing services...");
    
    // Start and enable services
    service("start", "nginx");
    service("enable", "nginx");
    
    // Check service status
    service("status", "nginx");
    
    // Restart if needed
    service("restart", "nginx");
}
```

### File Operations
```scl
fn setup_files() {
    print("Setting up files and directories...");
    
    // Create directories and files
    create("/var/log/myapp", "app.log", "644");
    create("/etc/myapp", "config.json", "600");
    
    // Copy configuration files
    copy("nginx.conf", "/etc/nginx/nginx.conf");
    copy("app.jar", "/opt/myapp/app.jar");
}
```

## Real-World Examples

### Example 1: Basic Web Server Setup
```scl
// basic-webserver.scl
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "root@webserver.com";

fn main() {
    print("🚀 Setting up basic web server");
    
    system_check();
    install_webserver();
    configure_firewall();
    start_services();
    
    print("✅ Web server setup completed");
}

fn system_check() {
    print("📊 Checking system...");
    sysinfo();
    monitor();
}

fn install_webserver() {
    print("📦 Installing web server...");
    package("update", "");
    package("install", "nginx");
}

fn configure_firewall() {
    print("🔥 Configuring firewall...");
    firewall("allow", "22/tcp");    // SSH
    firewall("allow", "80/tcp");    // HTTP
    firewall("allow", "443/tcp");   // HTTPS
}

fn start_services() {
    print("⚙️ Starting services...");
    service("enable", "nginx");
    service("start", "nginx");
    service("status", "nginx");
}
```

### Example 2: Conditional Deployment
```scl
// conditional-deployment.scl
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "admin@server.com";

// Configuration variables
declare environment: string = "production";
declare ssl_enabled: bool = true;
declare backup_enabled: bool = true;

fn main() {
    print("🚀 Starting conditional deployment");
    
    // Environment-specific setup
    if (environment == "production") {
        production_setup();
    } else {
        development_setup();
    }
    
    // Common setup for all environments
    common_setup();
    
    // Optional features
    if (ssl_enabled) {
        setup_ssl();
    }
    
    if (backup_enabled) {
        setup_backup();
    }
    
    print("✅ Deployment completed for: " + environment);
}

fn production_setup() {
    print("🏭 Production configuration");
    
    // Production-specific packages
    package("install", "fail2ban");
    package("install", "logrotate");
    
    // Security hardening
    firewall("enable");
    service("enable", "fail2ban");
    
    // Performance tuning
    tune("swappiness", "10");
    tune("file_limits", "65536");
}

fn development_setup() {
    print("🛠️ Development configuration");
    
    // Development tools
    package("install", "git");
    package("install", "curl");
    package("install", "vim");
    
    // Development ports
    firewall("allow", "3000/tcp");
    firewall("allow", "8080/tcp");
}

fn common_setup() {
    print("⚙️ Common setup");
    
    // Basic packages
    package("install", "nginx");
    package("install", "htop");
    
    // Basic services
    service("enable", "nginx");
    service("start", "nginx");
}

fn setup_ssl() {
    print("🔐 Setting up SSL");
    
    if (environment == "production") {
        cert("letsencrypt", "example.com");
    } else {
        cert("generate_self_signed", "dev.example.com");
    }
    
    firewall("allow", "443/tcp");
}

fn setup_backup() {
    print("💾 Setting up backup");
    
    // Create backup directory
    create("/backup", "daily", "755");
    
    // Setup automated backup
    cron("add", "0 2 * * *", "/usr/local/bin/backup.sh", "root");
    
    // Initial backup
    backup("backup", "/etc", "/backup/etc-initial.tar.gz");
}
```

### Example 3: Docker Application Deployment
```scl
// docker-deployment.scl
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "devops@docker-host.com";

declare app_name: string = "myapp";
declare app_version: string = "latest";
declare app_port: int = 3000;

fn main() {
    print("🐳 Docker application deployment");
    
    system_check();
    install_docker();
    deploy_application();
    verify_deployment();
    
    print("✅ Docker deployment completed");
}

fn system_check() {
    print("📊 System check");
    sysinfo();
    monitor();
}

fn install_docker() {
    print("🐳 Installing Docker");
    
    package("update", "");
    package("install", "docker.io");
    
    service("enable", "docker");
    service("start", "docker");
    
    // Add user to docker group
    user("add", "appuser");
    // Note: In real implementation, you'd add user to docker group
}

fn deploy_application() {
    print("🚀 Deploying application");
    
    // Pull application image
    docker("pull", app_name + ":" + app_version);
    
    // Stop existing container if running
    docker("stop", app_name);
    
    // Remove old container
    docker("rm", app_name);
    
    // Run new container
    docker("run", "-d", "--name", app_name, "-p", "80:" + app_port, app_name + ":" + app_version);
    
    // Configure firewall
    firewall("allow", "80/tcp");
}

fn verify_deployment() {
    print("🔍 Verifying deployment");
    
    // Check container status
    docker("ps");
    
    // Check logs
    docker("logs", app_name);
    
    // Test network connectivity
    nettest("localhost:80", "telnet");
    
    print("✅ Deployment verification completed");
}
```

### Example 4: Security Hardening
```scl
// security-hardening.scl
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "root@secure-server.com";

fn main() {
    print("🔒 Security hardening deployment");
    
    system_audit();
    harden_ssh();
    configure_firewall();
    install_security_tools();
    setup_monitoring();
    final_audit();
    
    print("✅ Security hardening completed");
}

fn system_audit() {
    print("🔍 Initial security audit");
    audit();
}

fn harden_ssh() {
    print("🔐 Hardening SSH configuration");
    
    // Note: In real implementation, you'd modify SSH config
    // This is a simplified example
    service("restart", "ssh");
}

fn configure_firewall() {
    print("🔥 Configuring firewall");
    
    // Enable firewall
    firewall("enable");
    
    // Allow only necessary ports
    firewall("allow", "22/tcp");    // SSH
    firewall("allow", "80/tcp");    // HTTP
    firewall("allow", "443/tcp");   // HTTPS
    
    // Deny all other traffic by default
    firewall("deny", "23/tcp");     // Telnet
    firewall("deny", "21/tcp");     // FTP
}

fn install_security_tools() {
    print("🛡️ Installing security tools");
    
    package("install", "fail2ban");
    package("install", "rkhunter");
    package("install", "chkrootkit");
    
    service("enable", "fail2ban");
    service("start", "fail2ban");
}

fn setup_monitoring() {
    print("📊 Setting up security monitoring");
    
    // Setup log monitoring
    logs("/var/log/auth.log", "Failed", 10);
    
    // Setup system monitoring
    monitor();
    
    // Setup automated security scans
    cron("add", "0 3 * * 0", "/usr/bin/rkhunter --check", "root");
}

fn final_audit() {
    print("🔍 Final security audit");
    audit();
    
    print("🔒 Security hardening summary:");
    print("- SSH hardened");
    print("- Firewall configured");
    print("- Security tools installed");
    print("- Monitoring enabled");
}
```

## Best Practices

### 1. Code Organization
```scl
// Good: Organize code with clear sections
// Configuration section
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "user@server.com";

// Variable declarations
declare environment: string = "production";
declare ssl_enabled: bool = true;

// Main function
fn main() {
    // Clear, descriptive function calls
    system_check();
    deploy_application();
    verify_deployment();
}
```

### 2. Meaningful Names
```scl
// Good: Descriptive variable names
declare web_server_port: int = 80;
declare ssl_certificate_path: string = "/etc/ssl/certs/";
declare backup_enabled: bool = true;

// Good: Descriptive function names
fn install_web_server() { ... }
fn configure_ssl_certificates() { ... }
fn setup_database_backup() { ... }
```

### 3. Error Handling
```scl
// Good: Check conditions before proceeding
declare tools_ready: bool = check(target, ["docker", "nginx"]);

fn main() {
    if (tools_ready) {
        deploy_application();
    } else {
        install_prerequisites();
        deploy_application();
    }
}
```

### 4. Comments and Documentation
```scl
// Good: Document complex logic
fn configure_load_balancer() {
    // Configure nginx as a load balancer for multiple backend servers
    // This setup distributes traffic across web1, web2, and web3
    
    copy("nginx-lb.conf", "/etc/nginx/nginx.conf");
    
    // Test configuration before applying
    webserver("nginx", "test");
    webserver("nginx", "reload");
}
```

### 5. Modular Functions
```scl
// Good: Break complex tasks into smaller functions
fn main() {
    system_preparation();
    application_deployment();
    security_configuration();
    monitoring_setup();
}

fn system_preparation() {
    sysinfo();
    package("update", "");
}

fn application_deployment() {
    install_application();
    configure_application();
    start_application();
}

fn security_configuration() {
    configure_firewall();
    setup_ssl();
    harden_system();
}
```

This tutorial provides a comprehensive introduction to SCL, from basic syntax to real-world DevOps scenarios. Each example builds upon previous concepts, making it easy to learn the language progressively.