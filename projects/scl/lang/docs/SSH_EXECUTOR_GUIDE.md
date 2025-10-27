# Enhanced SSH Executor with DevOps Functionality

This enhanced SSH executor provides comprehensive remote server management capabilities using the `goph` SSH library, designed specifically for DevOps engineers and system administrators.

## Features

### Core SSH Functionality
- **Multiple Authentication Methods**: SSH keys, passwords, and SSH agent
- **Secure Connections**: Built on the robust `goph` library
- **File Transfer**: Upload/download files with SFTP support
- **Command Execution**: Run commands with timeout support
- **Script Execution**: Upload and execute local scripts remotely

### DevOps Operations

#### 1. System Information & Monitoring
```go
// Get comprehensive system information
sysInfo, err := sshExec.SystemInfo()

// Collect monitoring metrics
metrics, err := sshExec.MonitoringMetrics()

// Performance analysis
perfResults, err := devopsUtils.PerformanceAnalysis(5)
```

#### 2. Package Management
```go
// Install packages (auto-detects package manager: apt, yum, dnf)
err := sshExec.PackageManagement("install", "nginx")
err := sshExec.PackageManagement("update", "")
err := sshExec.PackageManagement("search", "docker")
```

#### 3. Service Management
```go
// Control systemd services
err := sshExec.ProcessManagement("start", "nginx")
err := sshExec.ProcessManagement("stop", "apache2")
err := sshExec.ProcessManagement("restart", "mysql")
err := sshExec.ProcessManagement("status", "docker")
```

#### 4. File Operations
```go
// Comprehensive file management
err := sshExec.FileOperations("upload", "local.txt", "/remote/path.txt")
err := sshExec.FileOperations("download", "/remote/file", "./local/file")
err := sshExec.FileOperations("mkdir", "/new/directory", "", "755")
err := sshExec.FileOperations("chmod", "/path/to/file", "", "644")
err := sshExec.FileOperations("chown", "/path/to/file", "", "user:group")
```

#### 5. Network Diagnostics
```go
// Network connectivity testing
err := sshExec.NetworkDiagnostics("google.com", "ping")
err := sshExec.NetworkDiagnostics("example.com:80", "telnet")
err := sshExec.NetworkDiagnostics("https://api.example.com", "curl")
err := sshExec.NetworkDiagnostics("example.com", "nslookup")
```

#### 6. Docker Operations
```go
// Docker container management
err := sshExec.DockerOperations("ps")
err := sshExec.DockerOperations("start", "container_name")
err := sshExec.DockerOperations("logs", "container_name")
err := sshExec.DockerOperations("exec", "container_name", "bash")
```

#### 7. Database Operations
```go
// Database backup and restore
dbParams := map[string]string{
    "user":     "dbuser",
    "password": "dbpass",
    "database": "mydb",
    "output":   "/backup/db.sql",
}
err := devopsUtils.DatabaseOperations("mysql", "backup", dbParams)
err := devopsUtils.DatabaseOperations("postgresql", "restore", dbParams)
```

#### 8. Web Server Management
```go
// Nginx/Apache operations
err := devopsUtils.WebServerOperations("nginx", "test")
err := devopsUtils.WebServerOperations("nginx", "reload")
err := devopsUtils.WebServerOperations("apache", "restart")
```

#### 9. SSL/TLS Certificate Management
```go
// Certificate operations
err := devopsUtils.CertificateManagement("generate_self_signed", "example.com")
err := devopsUtils.CertificateManagement("check_expiry", "example.com")
err := devopsUtils.CertificateManagement("letsencrypt", "example.com")
```

#### 10. Firewall Management
```go
// UFW firewall operations
err := devopsUtils.FirewallManagement("status")
err := devopsUtils.FirewallManagement("allow", "22/tcp")
err := devopsUtils.FirewallManagement("deny", "80/tcp")
err := devopsUtils.FirewallManagement("enable")
```

#### 11. Cron Job Management
```go
// Cron operations
err := devopsUtils.CronManagement("add", "0 2 * * *", "/usr/bin/backup.sh", "root")
err := devopsUtils.CronManagement("list", "", "", "root")
err := devopsUtils.CronManagement("remove", "", "/usr/bin/backup.sh", "root")
```

#### 12. User Management
```go
// User and group operations
err := devopsUtils.UserManagement("add", "newuser", "sudo")
err := devopsUtils.UserManagement("passwd", "username", "newpassword")
err := devopsUtils.UserManagement("sudo", "username")
err := devopsUtils.UserManagement("lock", "username")
```

#### 13. Backup Operations
```go
// Backup and restore
err := devopsUtils.BackupOperations("backup", "/etc", "/backup/etc.tar.gz")
err := devopsUtils.BackupOperations("restore", "/backup/etc.tar.gz", "/restore/path")
err := devopsUtils.BackupOperations("sync", "/source", "/destination", "*.log")
```

#### 14. System Tuning
```go
// Performance tuning
err := devopsUtils.SystemTuning("swappiness", "10")
err := devopsUtils.SystemTuning("file_limits", "65536")
err := devopsUtils.SystemTuning("tcp_tuning")
```

#### 15. Security & Compliance
```go
// Security audit
issues, err := sshExec.SecurityAudit()

// Compliance checks
complianceResults, err := devopsUtils.ComplianceCheck("cis_benchmark")
complianceResults, err := devopsUtils.ComplianceCheck("pci_dss")
```

## Connection Examples

### SSH Key Authentication
```go
sshExec := executor.NewSSHExecutor()
err := sshExec.ConnectWithKey("192.168.1.100", "root", "/home/user/.ssh/id_rsa", "")
defer sshExec.Close()
```

### Password Authentication
```go
err := sshExec.ConnectWithPassword("192.168.1.100", "root", "password")
```

### SSH Agent Authentication
```go
err := sshExec.ConnectWithAgent("192.168.1.100", "root")
```

### Protected Private Key
```go
err := sshExec.ConnectWithKey("192.168.1.100", "root", "/home/user/.ssh/id_rsa", "passphrase")
```

## Advanced Features

### Command Timeout
```go
output, err := sshExec.RunWithTimeout("long_running_command", 30*time.Second)
```

### SFTP Operations
```go
sftp, err := sshExec.GetSFTPClient()
if err == nil {
    file, err := sftp.Create("/remote/path/file.txt")
    file.Write([]byte("Hello World"))
    file.Close()
}
```

### Script Execution
```go
err := sshExec.ExecuteScript("local_script.sh", "arg1", "arg2")
```

## Error Handling

The executor provides detailed error messages with context:

```go
if err != nil {
    log.Printf("Operation failed: %v", err)
    // Error messages include command output and troubleshooting hints
}
```

## Best Practices

1. **Always close connections**: Use `defer sshExec.Close()`
2. **Enable verbose mode for debugging**: `sshExec.SetVerbose(true)`
3. **Handle errors appropriately**: Check return values and log errors
4. **Use timeouts for long operations**: Prevent hanging connections
5. **Validate inputs**: Ensure parameters are properly formatted
6. **Test connectivity first**: Use `SystemInfo()` to verify connection

## Security Considerations

- Use SSH keys instead of passwords when possible
- Implement proper key management and rotation
- Use SSH agent for secure key handling
- Validate all user inputs to prevent command injection
- Use least privilege principle for user accounts
- Enable firewall and proper access controls
- Regular security audits and compliance checks

## Dependencies

```go
go get github.com/melbahja/goph
```

The executor automatically handles all SSH and SFTP operations through the goph library, providing a robust and secure foundation for remote server management.

## Integration with SCL

This SSH executor can be integrated into the SCL (Server Configuration Language) system to provide remote execution capabilities for infrastructure automation and configuration management tasks.