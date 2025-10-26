package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SCL/internal/executor"
)

func main() {
	// Create SSH executor
	sshExec := executor.NewSSHExecutor()
	sshExec.SetVerbose(true)

	// Example 1: Connect with SSH key
	err := sshExec.ConnectWithKey("192.168.1.100", "root", "/home/user/.ssh/id_rsa", "")
	if err != nil {
		log.Printf("Key connection failed, trying password: %v", err)
		
		// Example 2: Connect with password
		err = sshExec.ConnectWithPassword("192.168.1.100", "root", "your_password")
		if err != nil {
			log.Printf("Password connection failed, trying SSH agent: %v", err)
			
			// Example 3: Connect with SSH agent
			err = sshExec.ConnectWithAgent("192.168.1.100", "root")
			if err != nil {
				log.Fatalf("All connection methods failed: %v", err)
			}
		}
	}

	defer sshExec.Close()

	// Example 4: Get system information
	fmt.Println("=== System Information ===")
	sysInfo, err := sshExec.SystemInfo()
	if err != nil {
		log.Printf("Failed to get system info: %v", err)
	} else {
		for key, value := range sysInfo {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Example 5: Package management
	fmt.Println("\n=== Package Management ===")
	err = sshExec.PackageManagement("install", "htop")
	if err != nil {
		log.Printf("Package installation failed: %v", err)
	}

	// Example 6: Service management
	fmt.Println("\n=== Service Management ===")
	err = sshExec.ProcessManagement("status", "nginx")
	if err != nil {
		log.Printf("Service status check failed: %v", err)
	}

	// Example 7: File operations
	fmt.Println("\n=== File Operations ===")
	err = sshExec.FileOperations("upload", "local_file.txt", "/tmp/remote_file.txt")
	if err != nil {
		log.Printf("File upload failed: %v", err)
	}

	err = sshExec.FileOperations("mkdir", "/tmp/test_dir", "", "755")
	if err != nil {
		log.Printf("Directory creation failed: %v", err)
	}

	// Example 8: Network diagnostics
	fmt.Println("\n=== Network Diagnostics ===")
	err = sshExec.NetworkDiagnostics("google.com", "ping")
	if err != nil {
		log.Printf("Network test failed: %v", err)
	}

	// Example 9: Docker operations
	fmt.Println("\n=== Docker Operations ===")
	err = sshExec.DockerOperations("ps")
	if err != nil {
		log.Printf("Docker ps failed: %v", err)
	}

	// Example 10: Log analysis
	fmt.Println("\n=== Log Analysis ===")
	err = sshExec.LogAnalysis("/var/log/syslog", "error", 10)
	if err != nil {
		log.Printf("Log analysis failed: %v", err)
	}

	// Example 11: Monitoring metrics
	fmt.Println("\n=== Monitoring Metrics ===")
	metrics, err := sshExec.MonitoringMetrics()
	if err != nil {
		log.Printf("Failed to get metrics: %v", err)
	} else {
		for key, value := range metrics {
			fmt.Printf("%s: %v\n", key, value)
		}
	}

	// Example 12: Security audit
	fmt.Println("\n=== Security Audit ===")
	issues, err := sshExec.SecurityAudit()
	if err != nil {
		log.Printf("Security audit failed: %v", err)
	} else {
		for _, issue := range issues {
			fmt.Printf("Security check: %s\n", issue)
		}
	}

	// Example 13: Execute command with timeout
	fmt.Println("\n=== Command with Timeout ===")
	output, err := sshExec.RunWithTimeout("sleep 2 && echo 'Command completed'", 5*time.Second)
	if err != nil {
		log.Printf("Timeout command failed: %v", err)
	} else {
		fmt.Printf("Command output: %s\n", string(output))
	}

	// Example 14: Execute local script on remote server
	fmt.Println("\n=== Script Execution ===")
	err = sshExec.ExecuteScript("local_script.sh", "arg1", "arg2")
	if err != nil {
		log.Printf("Script execution failed: %v", err)
	}

	// Example 15: Advanced DevOps operations
	fmt.Println("\n=== Advanced DevOps Operations ===")
	devopsUtils := executor.NewDevOpsUtils(sshExec)

	// Backup operations
	err = devopsUtils.BackupOperations("backup", "/etc", "/tmp/etc_backup.tar.gz")
	if err != nil {
		log.Printf("Backup failed: %v", err)
	}

	// Database operations
	dbParams := map[string]string{
		"user":     "dbuser",
		"password": "dbpass",
		"database": "mydb",
		"output":   "/tmp/db_backup.sql",
	}
	err = devopsUtils.DatabaseOperations("mysql", "backup", dbParams)
	if err != nil {
		log.Printf("Database backup failed: %v", err)
	}

	// Web server operations
	err = devopsUtils.WebServerOperations("nginx", "test")
	if err != nil {
		log.Printf("Nginx test failed: %v", err)
	}

	// Certificate management
	err = devopsUtils.CertificateManagement("check_expiry", "example.com")
	if err != nil {
		log.Printf("Certificate check failed: %v", err)
	}

	// Firewall management
	err = devopsUtils.FirewallManagement("status")
	if err != nil {
		log.Printf("Firewall status failed: %v", err)
	}

	// Cron management
	err = devopsUtils.CronManagement("add", "0 2 * * *", "/usr/bin/backup.sh", "root")
	if err != nil {
		log.Printf("Cron add failed: %v", err)
	}

	// User management
	err = devopsUtils.UserManagement("add", "newuser", "sudo")
	if err != nil {
		log.Printf("User add failed: %v", err)
	}

	// System tuning
	err = devopsUtils.SystemTuning("swappiness", "10")
	if err != nil {
		log.Printf("System tuning failed: %v", err)
	}

	// Compliance check
	complianceResults, err := devopsUtils.ComplianceCheck("cis_benchmark")
	if err != nil {
		log.Printf("Compliance check failed: %v", err)
	} else {
		fmt.Println("Compliance check results:")
		for _, result := range complianceResults {
			fmt.Printf("  %s\n", result)
		}
	}

	// Performance analysis
	perfResults, err := devopsUtils.PerformanceAnalysis(5)
	if err != nil {
		log.Printf("Performance analysis failed: %v", err)
	} else {
		fmt.Println("Performance analysis results:")
		for key, value := range perfResults {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Println("\n=== All operations completed ===")
}