# SCL Language Reference

## Table of Contents
1. [Overview](#overview)
2. [Language Syntax](#language-syntax)
3. [Tokens and Keywords](#tokens-and-keywords)
4. [Data Types](#data-types)
5. [Variables and Assignments](#variables-and-assignments)
6. [Functions](#functions)
7. [Control Flow](#control-flow)
8. [Built-in Commands](#built-in-commands)
9. [Execution Modes](#execution-modes)
10. [Complete Examples](#complete-examples)

## Overview

SCL (Server Configuration Language) is a declarative infrastructure configuration language that compiles to bash scripts or executes directly via SSH. It provides a clean, readable syntax for DevOps automation and server management.

## Language Syntax

### Basic Structure
```scl
// Comments start with //
/* Multi-line comments
   are also supported */

// Import statements (optional)
import primary;

// Configuration variables
mode := "interpret";
setting := "configuration";
target := "user@server.com";

// Function definitions
fn main() {
    // Function body
    print("Hello, World!");
}
```

## Tokens and Keywords

### Reserved Keywords
```scl
// Control flow
if          // Conditional statement
else        // Alternative branch
fn          // Function declaration

// Variable declaration
declare     // Explicit variable declaration

// Import system
import      // Module import

// Data types
bool        // Boolean type
string      // String type
int         // Integer type
float       // Float type
list        // List/Array type

// Boolean literals
true        // Boolean true
false       // Boolean false
```

### Operators

#### Assignment Operators
```scl
:=          // Variable assignment (shell := bash)
=           // Value assignment (super_user = true)
```

#### Comparison Operators
```scl
==          // Equal to
!=          // Not equal to
<           // Less than
>           // Greater than
<=          // Less than or equal
>=          // Greater than or equal
```

#### Logical Operators
```scl
&&          // Logical AND
||          // Logical OR
!           // Logical NOT
```

#### Arithmetic Operators
```scl
+           // Addition
-           // Subtraction
*           // Multiplication
/           // Division
%           // Modulo
```

### Delimiters
```scl
(           // Left parenthesis
)           // Right parenthesis
{           // Left brace
}           // Right brace
[           // Left bracket
]           // Right bracket
;           // Semicolon (statement terminator)
,           // Comma (separator)
:           // Colon (type annotation)
.           // Dot (member access)
```

### Literals
```scl
// Numbers
42          // Integer
3.14        // Float

// Strings
"hello"     // String literal
"multi word string"

// Arrays
[1, 2, 3]   // Integer array
["a", "b"]  // String array
```

## Data Types

### Primitive Types
```scl
// Boolean
declare is_admin: bool = true;
declare is_enabled: bool = false;

// String
declare hostname: string = "web-server";
declare config_path: string = "/etc/nginx/nginx.conf";

// Integer
declare port: int = 8080;
declare timeout: int = 30;

// Float
declare cpu_threshold: float = 85.5;
declare memory_limit: float = 2.5;

// List/Array
declare servers: list = ["web1", "web2", "web3"];
declare ports: list = [80, 443, 8080];
```

## Variables and Assignments

### Variable Assignment
```scl
// Simple assignment
mode := "interpret";
setting := "configuration";
target := "root@192.168.1.100";

// Value assignment
super_user = true;
on_error = "rollback";
debug_mode = false;

// Array assignment
servers := ["web1.example.com", "web2.example.com"];
ports := [80, 443, 8080];
```

### Variable Declaration
```scl
// Explicit type declaration
declare environment: string = "production";
declare max_connections: int = 1000;
declare ssl_enabled: bool = true;
declare backup_dirs: list = ["/var/www", "/etc/nginx"];

// Type inference
declare auto_string = "inferred as string";
declare auto_number = 42;
declare auto_bool = true;
```

## Functions

### Function Declaration
```scl
// Basic function
fn main() {
    print("Starting deployment...");
    setup_server();
    deploy_application();
    print("Deployment completed!");
}

// Function with parameters (future feature)
fn setup_user(username: string, is_admin: bool) {
    user("add", username);
    if (is_admin) {
        user("sudo", username);
    }
}
```

### Function Calls
```scl
// Built-in function calls
print("Message");
sysinfo();
monitor();

// Function calls with parameters
package("install", "nginx");
service("start", "nginx");
create("/var/log/app", "app.log", "644");

// Qualified function calls
primary.print("Using primary module");
```

## Control Flow

### If-Else Statements
```scl
// Simple if statement
if (tools_ready) {
    print("Tools are available");
    install("snap", "docker");
}

// If-else statement
if (ssl_enabled) {
    print("Setting up SSL certificates");
    cert("generate_self_signed", "example.com");
} else {
    print("SSL disabled, using HTTP only");
    firewall("allow", "80/tcp");
}

// Complex conditions
if (environment == "production") {
    print("Production deployment");
    backup("backup", "/var/www", "/backup/www.tar.gz");
    service("restart", "nginx");
} else {
    print("Development deployment");
    service("reload", "nginx");
}

// Nested if statements
if (super_user) {
    if (ssl_enabled) {
        cert("letsencrypt", "example.com");
        firewall("allow", "443/tcp");
    } else {
        firewall("allow", "80/tcp");
    }
    service("enable", "nginx");
}
```

### Conditional Expressions
```scl
// Boolean expressions
declare tools_ready: bool = check(target, ["docker", "nginx"]);

// Comparison expressions
if (port > 1024) {
    print("Using non-privileged port");
}

if (environment != "development") {
    backup("backup", "/etc", "/backup/etc.tar.gz");
}

// Logical expressions
if (ssl_enabled && environment == "production") {
    cert("letsencrypt", domain);
}

if (debug_mode || verbose_logging) {
    logs("/var/log/app.log", "DEBUG", 50);
}
```

## Built-in Commands

### System Information
```scl
// Get comprehensive system information
sysinfo();

// Monitor system metrics
monitor();
```

### Package Management
```scl
// Install packages (auto-detects package manager)
package("install", "nginx");
package("install", "docker.io");
package("update", "");
package("search", "postgresql");
```

### Service Management
```scl
// Control systemd services
service("start", "nginx");
service("stop", "apache2");
service("restart", "mysql");
service("status", "docker");
service("enable", "nginx");
service("disable", "apache2");
```

### File Operations
```scl
// Create files and directories
create("/var/log/app", "app.log", "644");
create("/etc/app", "config.json", "600");

// Copy files (local to remote)
copy("nginx.conf", "/etc/nginx/nginx.conf");
copy("app.jar", "/opt/app/app.jar");
```

### Network Operations
```scl
// Network diagnostics
nettest("google.com", "ping");
nettest("example.com:443", "telnet");
nettest("https://api.example.com", "curl");
nettest("example.com", "nslookup");
```

### Security Operations
```scl
// Security audit
audit();

// Firewall management
firewall("status");
firewall("enable");
firewall("allow", "22/tcp");
firewall("allow", "80/tcp");
firewall("deny", "23/tcp");
```

### User Management
```scl
// User operations
user("add", "devops", "sudo");
user("passwd", "devops", "secure_password");
user("lock", "olduser");
user("sudo", "devops");
```

### Container Operations
```scl
// Docker operations
docker("ps");
docker("images");
docker("start", "web-container");
docker("stop", "web-container");
docker("logs", "web-container");
```

### Advanced Operations
```scl
// Certificate management
cert("generate_self_signed", "example.com");
cert("check_expiry", "example.com");
cert("letsencrypt", "example.com");

// Cron job management
cron("add", "0 2 * * *", "/usr/bin/backup.sh", "root");
cron("list", "", "", "root");

// Backup operations
backup("backup", "/etc", "/backup/etc.tar.gz");
backup("restore", "/backup/etc.tar.gz", "/restore/path");

// Log analysis
logs("/var/log/syslog", "error", 10);
logs("/var/log/nginx/access.log", "404", 5);

// System tuning
tune("swappiness", "10");
tune("file_limits", "65536");

// Web server operations
webserver("nginx", "test");
webserver("nginx", "reload");
webserver("apache", "restart");

// Database operations
database("mysql", "status");
database("postgresql", "backup", "user=dbuser", "database=mydb", "output=/backup/db.sql");
```

## Execution Modes

### Compile Mode
```scl
// Generates bash script
mode := "compile";
setting := "configuration";
super_user := true;
on_error := "rollback";

fn main() {
    print("This will generate bash code");
    package("install", "nginx");
}
```

### Interpret Mode (SSH Execution)
```scl
// Direct SSH execution
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "user@server.com";

fn main() {
    print("This will execute via SSH");
    sysinfo();
    package("install", "nginx");
}
```

## Complete Examples

### Basic Web Server Setup
```scl
// Basic web server deployment
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "root@web-server.com";

fn main() {
    print("🚀 Setting up web server");
    
    // System check
    sysinfo();
    
    // Install web server
    package("install", "nginx");
    
    // Configure firewall
    firewall("allow", "80/tcp");
    firewall("allow", "443/tcp");
    
    // Start services
    service("enable", "nginx");
    service("start", "nginx");
    
    // Copy configuration
    copy("nginx.conf", "/etc/nginx/nginx.conf");
    
    // Test configuration and reload
    webserver("nginx", "test");
    webserver("nginx", "reload");
    
    print("✅ Web server setup completed");
}
```

### Comprehensive DevOps Deployment
```scl
// Full DevOps deployment with monitoring
import primary;

mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "devops@production-server.com";

declare ssl_enabled: bool = true;
declare environment: string = "production";
declare domain: string = "example.com";

fn main() {
    print("🚀 Starting production deployment");
    
    // Pre-deployment checks
    system_check();
    
    // Security setup
    security_setup();
    
    // Application deployment
    deploy_application();
    
    // Post-deployment verification
    verify_deployment();
    
    print("✅ Production deployment completed");
}

fn system_check() {
    print("📊 System Analysis");
    sysinfo();
    monitor();
    audit();
}

fn security_setup() {
    print("🔒 Security Configuration");
    
    // Firewall setup
    firewall("enable");
    firewall("allow", "22/tcp");
    firewall("allow", "80/tcp");
    
    if (ssl_enabled) {
        firewall("allow", "443/tcp");
        cert("letsencrypt", domain);
    }
    
    // User management
    user("add", "appuser");
    user("passwd", "appuser", "secure_random_password");
}

fn deploy_application() {
    print("📦 Application Deployment");
    
    // Install dependencies
    package("update", "");
    package("install", "nginx");
    package("install", "docker.io");
    
    // Deploy application
    copy("app.jar", "/opt/app/app.jar");
    copy("nginx.conf", "/etc/nginx/nginx.conf");
    
    // Start services
    service("enable", "nginx");
    service("start", "nginx");
    service("enable", "docker");
    service("start", "docker");
    
    // Deploy containers
    docker("pull", "myapp:latest");
    docker("run", "-d", "--name", "myapp", "myapp:latest");
}

fn verify_deployment() {
    print("🔍 Deployment Verification");
    
    // Check services
    service("status", "nginx");
    service("status", "docker");
    
    // Network tests
    nettest("localhost:80", "telnet");
    
    if (ssl_enabled) {
        nettest("localhost:443", "telnet");
    }
    
    // Check logs
    logs("/var/log/nginx/access.log", "", 5);
    
    // Final system check
    monitor();
}
```

### Conditional Deployment Example
```scl
// Environment-specific deployment
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "admin@server.com";

declare environment: string = "staging";
declare enable_monitoring: bool = true;
declare backup_enabled: bool = false;

fn main() {
    print("🚀 Environment-specific deployment");
    
    // Environment-specific configuration
    if (environment == "production") {
        production_setup();
    } else {
        development_setup();
    }
    
    // Common setup
    common_setup();
    
    // Optional features
    if (enable_monitoring) {
        setup_monitoring();
    }
    
    if (backup_enabled) {
        setup_backup();
    }
    
    print("✅ Deployment completed for environment: " + environment);
}

fn production_setup() {
    print("🏭 Production Configuration");
    
    // Production-specific packages
    package("install", "fail2ban");
    package("install", "logrotate");
    
    // Security hardening
    firewall("enable");
    service("enable", "fail2ban");
    
    // SSL certificates
    cert("letsencrypt", "production.example.com");
}

fn development_setup() {
    print("🛠️ Development Configuration");
    
    // Development tools
    package("install", "git");
    package("install", "curl");
    
    // Relaxed firewall for development
    firewall("allow", "3000/tcp");
    firewall("allow", "8080/tcp");
}

fn common_setup() {
    print("⚙️ Common Configuration");
    
    // Basic packages
    package("install", "nginx");
    package("install", "htop");
    
    // Basic services
    service("enable", "nginx");
    service("start", "nginx");
}

fn setup_monitoring() {
    print("📊 Setting up monitoring");
    
    // Install monitoring tools
    package("install", "prometheus-node-exporter");
    service("enable", "prometheus-node-exporter");
    service("start", "prometheus-node-exporter");
    
    // Monitor system
    monitor();
}

fn setup_backup() {
    print("💾 Setting up backup");
    
    // Create backup directory
    create("/backup", "daily", "755");
    
    // Setup backup cron job
    cron("add", "0 2 * * *", "/usr/local/bin/backup.sh", "root");
    
    // Initial backup
    backup("backup", "/etc", "/backup/etc-initial.tar.gz");
}
```

### Error Handling Example
```scl
// Deployment with error handling
mode := "interpret";
setting := "configuration";
super_user := true;
on_error := "rollback";
target := "user@server.com";

declare tools_ready: bool = check(target, ["docker", "nginx"]);

fn main() {
    print("🚀 Deployment with error handling");
    
    // Pre-flight checks
    if (tools_ready) {
        print("✅ Required tools are available");
        deploy();
    } else {
        print("❌ Required tools missing, installing...");
        install_prerequisites();
        deploy();
    }
}

fn install_prerequisites() {
    print("📦 Installing prerequisites");
    
    package("update", "");
    package("install", "docker.io");
    package("install", "nginx");
    
    // Verify installation
    service("status", "docker");
    service("status", "nginx");
}

fn deploy() {
    print("🚀 Starting deployment");
    
    // Deploy with checks
    service("start", "docker");
    
    if (service("status", "docker")) {
        docker("pull", "myapp:latest");
        docker("run", "-d", "myapp:latest");
        print("✅ Application deployed successfully");
    } else {
        print("❌ Docker service failed to start");
    }
}
```

This comprehensive documentation covers all aspects of the SCL language, from basic syntax to advanced DevOps operations. The language provides a clean, declarative way to manage infrastructure with both compile-time bash generation and runtime SSH execution capabilities.