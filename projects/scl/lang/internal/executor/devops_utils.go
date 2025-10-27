package executor

import (
	"fmt"
	"strings"
	"time"

	"github.com/SCL/internal/utils"
)

// DevOpsUtils provides additional DevOps-specific functionality
type DevOpsUtils struct {
	executor *SSHExecutor
}

func NewDevOpsUtils(executor *SSHExecutor) *DevOpsUtils {
	return &DevOpsUtils{
		executor: executor,
	}
}

// BackupOperations handles backup and restore operations
func (d *DevOpsUtils) BackupOperations(operation, source, destination string, options ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	switch operation {
	case "backup":
		// Create compressed backup
		timestamp := time.Now().Format("20060102_150405")
		backupName := fmt.Sprintf("backup_%s.tar.gz", timestamp)
		if destination == "" {
			destination = fmt.Sprintf("/tmp/%s", backupName)
		}
		cmd = fmt.Sprintf("tar -czf %s %s", destination, source)
		
	case "restore":
		cmd = fmt.Sprintf("tar -xzf %s -C %s", source, destination)
		
	case "sync":
		// Rsync-like operation
		excludePattern := ""
		if len(options) > 0 {
			excludePattern = fmt.Sprintf("--exclude='%s'", options[0])
		}
		cmd = fmt.Sprintf("rsync -av %s %s %s", excludePattern, source, destination)
		
	default:
		return fmt.Errorf("unsupported backup operation: %s", operation)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing backup operation: %s", cmd))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("backup operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Backup operation '%s' completed", operation))
	return nil
}

// DatabaseOperations provides database management functionality
func (d *DevOpsUtils) DatabaseOperations(dbType, operation string, params map[string]string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	
	switch dbType {
	case "mysql":
		switch operation {
		case "backup":
			cmd = fmt.Sprintf("mysqldump -u%s -p%s %s > %s", 
				params["user"], params["password"], params["database"], params["output"])
		case "restore":
			cmd = fmt.Sprintf("mysql -u%s -p%s %s < %s", 
				params["user"], params["password"], params["database"], params["input"])
		case "status":
			cmd = "systemctl status mysql"
		}
		
	case "postgresql":
		switch operation {
		case "backup":
			cmd = fmt.Sprintf("pg_dump -U %s -h localhost %s > %s", 
				params["user"], params["database"], params["output"])
		case "restore":
			cmd = fmt.Sprintf("psql -U %s -h localhost %s < %s", 
				params["user"], params["database"], params["input"])
		case "status":
			cmd = "systemctl status postgresql"
		}
		
	case "mongodb":
		switch operation {
		case "backup":
			cmd = fmt.Sprintf("mongodump --db %s --out %s", params["database"], params["output"])
		case "restore":
			cmd = fmt.Sprintf("mongorestore --db %s %s", params["database"], params["input"])
		case "status":
			cmd = "systemctl status mongod"
		}
		
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing database operation: %s", operation))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("database operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Database %s operation completed", operation))
	return nil
}

// WebServerOperations manages web servers (Nginx, Apache)
func (d *DevOpsUtils) WebServerOperations(serverType, operation string, configPath ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	
	switch serverType {
	case "nginx":
		switch operation {
		case "test":
			cmd = "nginx -t"
		case "reload":
			cmd = "nginx -s reload"
		case "status":
			cmd = "systemctl status nginx"
		case "start":
			cmd = "systemctl start nginx"
		case "stop":
			cmd = "systemctl stop nginx"
		case "restart":
			cmd = "systemctl restart nginx"
		case "logs":
			cmd = "tail -f /var/log/nginx/error.log"
		}
		
	case "apache":
		switch operation {
		case "test":
			cmd = "apache2ctl configtest"
		case "reload":
			cmd = "systemctl reload apache2"
		case "status":
			cmd = "systemctl status apache2"
		case "start":
			cmd = "systemctl start apache2"
		case "stop":
			cmd = "systemctl stop apache2"
		case "restart":
			cmd = "systemctl restart apache2"
		case "logs":
			cmd = "tail -f /var/log/apache2/error.log"
		}
		
	default:
		return fmt.Errorf("unsupported web server: %s", serverType)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing web server operation: %s %s", serverType, operation))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("web server operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintInfo(fmt.Sprintf("Web server operation result:\n%s", string(output)))
	return nil
}

// CertificateManagement handles SSL/TLS certificates
func (d *DevOpsUtils) CertificateManagement(operation, domain string, options ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	
	switch operation {
	case "generate_self_signed":
		keyPath := fmt.Sprintf("/etc/ssl/private/%s.key", domain)
		certPath := fmt.Sprintf("/etc/ssl/certs/%s.crt", domain)
		cmd = fmt.Sprintf("openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout %s -out %s -subj '/CN=%s'", 
			keyPath, certPath, domain)
			
	case "check_expiry":
		if len(options) > 0 {
			cmd = fmt.Sprintf("openssl x509 -in %s -noout -dates", options[0])
		} else {
			cmd = fmt.Sprintf("echo | openssl s_client -servername %s -connect %s:443 2>/dev/null | openssl x509 -noout -dates", 
				domain, domain)
		}
		
	case "letsencrypt":
		cmd = fmt.Sprintf("certbot --nginx -d %s --non-interactive --agree-tos --email admin@%s", domain, domain)
		
	case "renew":
		cmd = "certbot renew --dry-run"
		
	default:
		return fmt.Errorf("unsupported certificate operation: %s", operation)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing certificate operation: %s", operation))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("certificate operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Certificate operation '%s' completed", operation))
	return nil
}

// FirewallManagement handles firewall configuration
func (d *DevOpsUtils) FirewallManagement(operation string, params ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	
	switch operation {
	case "status":
		cmd = "ufw status verbose"
		
	case "enable":
		cmd = "ufw --force enable"
		
	case "disable":
		cmd = "ufw disable"
		
	case "allow":
		if len(params) > 0 {
			cmd = fmt.Sprintf("ufw allow %s", params[0])
		}
		
	case "deny":
		if len(params) > 0 {
			cmd = fmt.Sprintf("ufw deny %s", params[0])
		}
		
	case "delete":
		if len(params) > 0 {
			cmd = fmt.Sprintf("ufw delete %s", params[0])
		}
		
	case "reset":
		cmd = "ufw --force reset"
		
	case "list_rules":
		cmd = "ufw status numbered"
		
	default:
		return fmt.Errorf("unsupported firewall operation: %s", operation)
	}

	if cmd == "" {
		return fmt.Errorf("insufficient parameters for firewall operation: %s", operation)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing firewall operation: %s", cmd))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("firewall operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintInfo(fmt.Sprintf("Firewall operation result:\n%s", string(output)))
	return nil
}

// CronManagement handles cron job operations
func (d *DevOpsUtils) CronManagement(operation, schedule, command string, user ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	targetUser := "root"
	if len(user) > 0 {
		targetUser = user[0]
	}

	var cmd string
	
	switch operation {
	case "add":
		cronEntry := fmt.Sprintf("%s %s", schedule, command)
		cmd = fmt.Sprintf("(crontab -u %s -l 2>/dev/null; echo '%s') | crontab -u %s -", 
			targetUser, cronEntry, targetUser)
			
	case "list":
		cmd = fmt.Sprintf("crontab -u %s -l", targetUser)
		
	case "remove":
		cmd = fmt.Sprintf("crontab -u %s -l | grep -v '%s' | crontab -u %s -", 
			targetUser, command, targetUser)
			
	case "clear":
		cmd = fmt.Sprintf("crontab -u %s -r", targetUser)
		
	default:
		return fmt.Errorf("unsupported cron operation: %s", operation)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing cron operation: %s", operation))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("cron operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Cron operation '%s' completed", operation))
	if len(output) > 0 {
		utils.PrintInfo(fmt.Sprintf("Output:\n%s", string(output)))
	}
	return nil
}

// UserManagement handles user and group operations
func (d *DevOpsUtils) UserManagement(operation, username string, params ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	
	switch operation {
	case "add":
		cmd = fmt.Sprintf("useradd -m -s /bin/bash %s", username)
		if len(params) > 0 {
			cmd += fmt.Sprintf(" -G %s", params[0]) // groups
		}
		
	case "delete":
		cmd = fmt.Sprintf("userdel -r %s", username)
		
	case "passwd":
		if len(params) > 0 {
			cmd = fmt.Sprintf("echo '%s:%s' | chpasswd", username, params[0])
		}
		
	case "lock":
		cmd = fmt.Sprintf("usermod -L %s", username)
		
	case "unlock":
		cmd = fmt.Sprintf("usermod -U %s", username)
		
	case "info":
		cmd = fmt.Sprintf("id %s && finger %s", username, username)
		
	case "list":
		cmd = "cut -d: -f1 /etc/passwd | sort"
		
	case "sudo":
		cmd = fmt.Sprintf("usermod -aG sudo %s", username)
		
	default:
		return fmt.Errorf("unsupported user operation: %s", operation)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing user operation: %s", operation))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("user operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ User operation '%s' completed", operation))
	if len(output) > 0 {
		utils.PrintInfo(fmt.Sprintf("Output:\n%s", string(output)))
	}
	return nil
}

// SystemTuning provides system performance tuning
func (d *DevOpsUtils) SystemTuning(operation string, params ...string) error {
	if !d.executor.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	
	switch operation {
	case "swappiness":
		if len(params) > 0 {
			cmd = fmt.Sprintf("echo 'vm.swappiness = %s' >> /etc/sysctl.conf && sysctl -p", params[0])
		}
		
	case "file_limits":
		if len(params) > 0 {
			cmd = fmt.Sprintf("echo '* soft nofile %s' >> /etc/security/limits.conf && echo '* hard nofile %s' >> /etc/security/limits.conf", 
				params[0], params[0])
		}
		
	case "tcp_tuning":
		cmd = `echo 'net.core.rmem_max = 16777216' >> /etc/sysctl.conf &&
		       echo 'net.core.wmem_max = 16777216' >> /etc/sysctl.conf &&
		       echo 'net.ipv4.tcp_rmem = 4096 87380 16777216' >> /etc/sysctl.conf &&
		       echo 'net.ipv4.tcp_wmem = 4096 65536 16777216' >> /etc/sysctl.conf &&
		       sysctl -p`
		       
	case "kernel_params":
		cmd = "sysctl -a | grep -E '(vm|net|kernel)' | head -20"
		
	default:
		return fmt.Errorf("unsupported tuning operation: %s", operation)
	}

	if d.executor.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing system tuning: %s", operation))
	}

	output, err := d.executor.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("system tuning failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ System tuning '%s' completed", operation))
	return nil
}

// ComplianceCheck performs security and compliance checks
func (d *DevOpsUtils) ComplianceCheck(checkType string) ([]string, error) {
	if !d.executor.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	var results []string
	
	switch checkType {
	case "cis_benchmark":
		checks := map[string]string{
			"Password policy": "grep -E '^(PASS_MAX_DAYS|PASS_MIN_DAYS|PASS_WARN_AGE)' /etc/login.defs",
			"SSH configuration": "grep -E '^(Protocol|PermitRootLogin|PasswordAuthentication)' /etc/ssh/sshd_config",
			"File permissions": "find /etc -type f -perm -002 | head -5",
			"World writable directories": "find / -type d -perm -002 2>/dev/null | head -5",
		}
		
		for check, cmd := range checks {
			output, err := d.executor.client.Run(cmd)
			if err != nil {
				results = append(results, fmt.Sprintf("%s: ERROR - %v", check, err))
			} else {
				results = append(results, fmt.Sprintf("%s: %s", check, strings.TrimSpace(string(output))))
			}
		}
		
	case "pci_dss":
		checks := map[string]string{
			"Firewall status": "ufw status",
			"Antivirus status": "systemctl status clamav-daemon 2>/dev/null || echo 'not installed'",
			"Log monitoring": "systemctl status rsyslog",
			"Access controls": "cat /etc/passwd | grep -E '^(root|admin):'",
		}
		
		for check, cmd := range checks {
			output, err := d.executor.client.Run(cmd)
			if err != nil {
				results = append(results, fmt.Sprintf("%s: ERROR - %v", check, err))
			} else {
				results = append(results, fmt.Sprintf("%s: %s", check, strings.TrimSpace(string(output))))
			}
		}
		
	default:
		return nil, fmt.Errorf("unsupported compliance check: %s", checkType)
	}

	return results, nil
}

// PerformanceAnalysis analyzes system performance
func (d *DevOpsUtils) PerformanceAnalysis(duration int) (map[string]string, error) {
	if !d.executor.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	analysis := make(map[string]string)
	
	// CPU analysis
	cpuCmd := fmt.Sprintf("sar -u %d 1 | tail -1", duration)
	if output, err := d.executor.client.Run(cpuCmd); err == nil {
		analysis["cpu_analysis"] = strings.TrimSpace(string(output))
	}

	// Memory analysis
	memCmd := "free -h && echo '---' && cat /proc/meminfo | grep -E '(MemTotal|MemFree|Buffers|Cached)'"
	if output, err := d.executor.client.Run(memCmd); err == nil {
		analysis["memory_analysis"] = strings.TrimSpace(string(output))
	}

	// Disk I/O analysis
	ioCmd := fmt.Sprintf("iostat -x %d 1 | tail -10", duration)
	if output, err := d.executor.client.Run(ioCmd); err == nil {
		analysis["io_analysis"] = strings.TrimSpace(string(output))
	}

	// Network analysis
	netCmd := "ss -tuln | wc -l && echo 'Active connections:' && ss -tuln | head -10"
	if output, err := d.executor.client.Run(netCmd); err == nil {
		analysis["network_analysis"] = strings.TrimSpace(string(output))
	}

	// Process analysis
	procCmd := "ps aux --sort=-%cpu | head -10"
	if output, err := d.executor.client.Run(procCmd); err == nil {
		analysis["process_analysis"] = strings.TrimSpace(string(output))
	}

	return analysis, nil
}