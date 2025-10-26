package executor

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SCL/internal/parser"
	"github.com/SCL/internal/utils"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHExecutor provides enhanced SSH-based remote execution with DevOps functionality
type SSHExecutor struct {
	*parser.BaseInfraDSLListener
	client        *goph.Client
	variables     map[string]string
	imports       map[string]bool
	functions     map[string]parser.IFunctionDeclarationContext
	functionCount int
	commandCount  int
	verbose       bool
	currentFunc   string
	connected     bool
}

// SSHConfig holds SSH connection configuration
type SSHConfig struct {
	Host       string
	Port       string
	User       string
	KeyPath    string
	Password   string
	Passphrase string
	UseAgent   bool
	Timeout    time.Duration
}

func NewSSHExecutor() *SSHExecutor {
	return &SSHExecutor{
		variables: make(map[string]string),
		imports:   make(map[string]bool),
		functions: make(map[string]parser.IFunctionDeclarationContext),
		verbose:   false,
		connected: false,
	}
}

func (e *SSHExecutor) SetVerbose(v bool) {
	e.verbose = v
}

// ConnectWithKey establishes SSH connection using private key
func (e *SSHExecutor) ConnectWithKey(host, user, keyPath, passphrase string) error {
	// Handle host:port format
	actualHost := host
	portStr := "22"
	
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		actualHost = parts[0]
		portStr = parts[1]
	}
	
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port number: %s", portStr)
	}
	
	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Connecting to %s@%s:%d with key %s", user, actualHost, port, keyPath))
	}

	auth, err := goph.Key(keyPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load SSH key: %v", err)
	}

	// Create SSH config with custom port
	config := goph.Config{
		User:     user,
		Addr:     actualHost,
		Port:     uint(port),
		Auth:     auth,
		Timeout:  30 * time.Second,
		Callback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := goph.NewConn(&config)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v", err)
	}

	e.client = client
	e.connected = true
	
	if e.verbose {
		utils.PrintSuccess("✓ SSH connection established")
	}
	
	return nil
}

// ConnectWithPassword establishes SSH connection using password
func (e *SSHExecutor) ConnectWithPassword(host, user, password string) error {
	// Handle host:port format
	actualHost := host
	portStr := "22"
	
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		actualHost = parts[0]
		portStr = parts[1]
	}
	
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port number: %s", portStr)
	}
	
	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Connecting to %s@%s:%d with password", user, actualHost, port))
	}

	// Create SSH config with custom port
	config := goph.Config{
		User:     user,
		Addr:     actualHost,
		Port:     uint(port),
		Auth:     goph.Password(password),
		Timeout:  30 * time.Second,
		Callback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := goph.NewConn(&config)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v", err)
	}

	e.client = client
	e.connected = true
	
	if e.verbose {
		utils.PrintSuccess("✓ SSH connection established")
	}
	
	return nil
}

// ConnectWithAgent establishes SSH connection using SSH agent
func (e *SSHExecutor) ConnectWithAgent(host, user string) error {
	// Handle host:port format
	actualHost := host
	portStr := "22"
	
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		actualHost = parts[0]
		portStr = parts[1]
	}
	
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port number: %s", portStr)
	}
	
	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Connecting to %s@%s:%d using SSH agent", user, actualHost, port))
	}

	auth, err := goph.UseAgent()
	if err != nil {
		return fmt.Errorf("failed to use SSH agent: %v", err)
	}

	// Create SSH config with custom port
	config := goph.Config{
		User:     user,
		Addr:     actualHost,
		Port:     uint(port),
		Auth:     auth,
		Timeout:  30 * time.Second,
		Callback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := goph.NewConn(&config)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v", err)
	}

	e.client = client
	e.connected = true
	
	if e.verbose {
		utils.PrintSuccess("✓ SSH connection established")
	}
	
	return nil
}

// Close closes the SSH connection
func (e *SSHExecutor) Close() error {
	if e.client != nil && e.connected {
		err := e.client.Close()
		e.connected = false
		if e.verbose {
			utils.PrintInfo("SSH connection closed")
		}
		return err
	}
	return nil
}

// SystemInfo retrieves comprehensive system information
func (e *SSHExecutor) SystemInfo() (map[string]string, error) {
	if !e.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	info := make(map[string]string)
	
	commands := map[string]string{
		"hostname":     "hostname",
		"os":          "cat /etc/os-release | grep PRETTY_NAME | cut -d'=' -f2 | tr -d '\"'",
		"kernel":      "uname -r",
		"architecture": "uname -m",
		"uptime":      "uptime -p",
		"load":        "uptime | awk -F'load average:' '{print $2}'",
		"memory":      "free -h | grep Mem | awk '{print \"Used: \" $3 \"/\" $2 \" (\" $3/$2*100 \"%)\"}'",
		"disk":        "df -h / | tail -1 | awk '{print \"Used: \" $3 \"/\" $2 \" (\" $5 \")\"}'",
		"cpu_cores":   "nproc",
		"users":       "who | wc -l",
	}

	for key, cmd := range commands {
		output, err := e.client.Run(cmd)
		if err != nil {
			info[key] = fmt.Sprintf("Error: %v", err)
		} else {
			info[key] = strings.TrimSpace(string(output))
		}
	}

	return info, nil
}

// ProcessManagement provides process control functionality
func (e *SSHExecutor) ProcessManagement(action, service string) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	switch action {
	case "start":
		cmd = fmt.Sprintf("sudo systemctl start %s", service)
	case "stop":
		cmd = fmt.Sprintf("sudo systemctl stop %s", service)
	case "restart":
		cmd = fmt.Sprintf("sudo systemctl restart %s", service)
	case "status":
		cmd = fmt.Sprintf("systemctl status %s", service)
	case "enable":
		cmd = fmt.Sprintf("sudo systemctl enable %s", service)
	case "disable":
		cmd = fmt.Sprintf("sudo systemctl disable %s", service)
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing: %s", cmd))
	}

	output, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Output: %s", string(output)))
	}

	return nil
}

// PackageManagement handles package operations across different distributions
func (e *SSHExecutor) PackageManagement(action, packageName string) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	// Detect package manager
	packageManager, err := e.detectPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %v", err)
	}

	var cmd string
	switch packageManager {
	case "apt":
		switch action {
		case "install":
			cmd = fmt.Sprintf("sudo apt-get update && sudo apt-get install -y %s", packageName)
		case "remove":
			cmd = fmt.Sprintf("sudo apt-get remove -y %s", packageName)
		case "update":
			cmd = "sudo apt-get update && sudo apt-get upgrade -y"
		case "search":
			cmd = fmt.Sprintf("apt-cache search %s", packageName)
		}
	case "yum":
		switch action {
		case "install":
			cmd = fmt.Sprintf("sudo yum install -y %s", packageName)
		case "remove":
			cmd = fmt.Sprintf("sudo yum remove -y %s", packageName)
		case "update":
			cmd = "sudo yum update -y"
		case "search":
			cmd = fmt.Sprintf("yum search %s", packageName)
		}
	case "dnf":
		switch action {
		case "install":
			cmd = fmt.Sprintf("sudo dnf install -y %s", packageName)
		case "remove":
			cmd = fmt.Sprintf("sudo dnf remove -y %s", packageName)
		case "update":
			cmd = "sudo dnf update -y"
		case "search":
			cmd = fmt.Sprintf("dnf search %s", packageName)
		}
	default:
		return fmt.Errorf("unsupported package manager: %s", packageManager)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing: %s", cmd))
	}

	output, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("package operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Package %s %s completed", packageName, action))
	return nil
}

// FileOperations provides comprehensive file management
func (e *SSHExecutor) FileOperations(operation, source, destination string, permissions ...string) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	switch operation {
	case "upload":
		return e.uploadFile(source, destination)
	case "download":
		return e.downloadFile(source, destination)
	case "copy":
		return e.copyFile(source, destination)
	case "move":
		return e.moveFile(source, destination)
	case "delete":
		return e.deleteFile(source)
	case "mkdir":
		return e.createDirectory(source, permissions...)
	case "chmod":
		if len(permissions) > 0 {
			return e.changePermissions(source, permissions[0])
		}
		return fmt.Errorf("chmod requires permissions parameter")
	case "chown":
		if len(permissions) > 0 {
			return e.changeOwnership(source, permissions[0])
		}
		return fmt.Errorf("chown requires owner parameter")
	default:
		return fmt.Errorf("unsupported file operation: %s", operation)
	}
}

// NetworkDiagnostics performs network connectivity tests
func (e *SSHExecutor) NetworkDiagnostics(target string, testType string) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	switch testType {
	case "ping":
		cmd = fmt.Sprintf("ping -c 4 %s", target)
	case "telnet":
		parts := strings.Split(target, ":")
		if len(parts) != 2 {
			return fmt.Errorf("telnet requires host:port format")
		}
		cmd = fmt.Sprintf("timeout 5 telnet %s %s", parts[0], parts[1])
	case "curl":
		cmd = fmt.Sprintf("curl -I -s -o /dev/null -w '%%{http_code}' %s", target)
	case "nslookup":
		cmd = fmt.Sprintf("nslookup %s", target)
	case "traceroute":
		cmd = fmt.Sprintf("traceroute %s", target)
	default:
		return fmt.Errorf("unsupported network test: %s", testType)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Running network test: %s", cmd))
	}

	output, err := e.client.Run(cmd)
	if err != nil {
		utils.PrintInfo(fmt.Sprintf("Network test failed: %v", err))
	}

	utils.PrintInfo(fmt.Sprintf("Network test result:\n%s", string(output)))
	return nil
}

// LogAnalysis provides log file analysis capabilities
func (e *SSHExecutor) LogAnalysis(logFile, pattern string, lines int) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	if pattern != "" {
		cmd = fmt.Sprintf("grep -n '%s' %s | tail -%d", pattern, logFile, lines)
	} else {
		cmd = fmt.Sprintf("tail -%d %s", lines, logFile)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Analyzing log: %s", cmd))
	}

	output, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("log analysis failed: %v", err)
	}

	utils.PrintInfo(fmt.Sprintf("Log analysis result:\n%s", string(output)))
	return nil
}

// DockerOperations provides Docker container management
func (e *SSHExecutor) DockerOperations(action string, args ...string) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	var cmd string
	switch action {
	case "ps":
		cmd = "docker ps -a"
	case "images":
		cmd = "docker images"
	case "start":
		if len(args) > 0 {
			cmd = fmt.Sprintf("docker start %s", args[0])
		}
	case "stop":
		if len(args) > 0 {
			cmd = fmt.Sprintf("docker stop %s", args[0])
		}
	case "restart":
		if len(args) > 0 {
			cmd = fmt.Sprintf("docker restart %s", args[0])
		}
	case "logs":
		if len(args) > 0 {
			cmd = fmt.Sprintf("docker logs %s", args[0])
		}
	case "exec":
		if len(args) > 1 {
			cmd = fmt.Sprintf("docker exec -it %s %s", args[0], strings.Join(args[1:], " "))
		}
	case "pull":
		if len(args) > 0 {
			cmd = fmt.Sprintf("docker pull %s", args[0])
		}
	case "build":
		if len(args) > 1 {
			cmd = fmt.Sprintf("docker build -t %s %s", args[0], args[1])
		}
	default:
		return fmt.Errorf("unsupported docker operation: %s", action)
	}

	if cmd == "" {
		return fmt.Errorf("insufficient arguments for docker %s", action)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing Docker command: %s", cmd))
	}

	output, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("docker operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintInfo(fmt.Sprintf("Docker operation result:\n%s", string(output)))
	return nil
}

// MonitoringMetrics collects system monitoring data
func (e *SSHExecutor) MonitoringMetrics() (map[string]interface{}, error) {
	if !e.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	metrics := make(map[string]interface{})
	
	// CPU usage
	cpuCmd := "top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | cut -d'%' -f1"
	if output, err := e.client.Run(cpuCmd); err == nil {
		metrics["cpu_usage"] = strings.TrimSpace(string(output))
	}

	// Memory usage
	memCmd := "free | grep Mem | awk '{printf \"%.2f\", $3/$2 * 100.0}'"
	if output, err := e.client.Run(memCmd); err == nil {
		metrics["memory_usage"] = strings.TrimSpace(string(output))
	}

	// Disk usage
	diskCmd := "df / | tail -1 | awk '{print $5}' | cut -d'%' -f1"
	if output, err := e.client.Run(diskCmd); err == nil {
		metrics["disk_usage"] = strings.TrimSpace(string(output))
	}

	// Network connections
	netCmd := "netstat -an | wc -l"
	if output, err := e.client.Run(netCmd); err == nil {
		metrics["network_connections"] = strings.TrimSpace(string(output))
	}

	// Load average
	loadCmd := "uptime | awk -F'load average:' '{print $2}'"
	if output, err := e.client.Run(loadCmd); err == nil {
		metrics["load_average"] = strings.TrimSpace(string(output))
	}

	return metrics, nil
}

// SecurityAudit performs basic security checks
func (e *SSHExecutor) SecurityAudit() ([]string, error) {
	if !e.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	var issues []string
	
	checks := map[string]string{
		"SSH root login": "grep '^PermitRootLogin' /etc/ssh/sshd_config || echo 'not configured'",
		"Password authentication": "grep '^PasswordAuthentication' /etc/ssh/sshd_config || echo 'not configured'",
		"Firewall status": "ufw status || iptables -L | head -5",
		"Failed login attempts": "grep 'Failed password' /var/log/auth.log | tail -5 || echo 'no recent failures'",
		"Open ports": "netstat -tuln | grep LISTEN",
		"World writable files": "find /etc -type f -perm -002 2>/dev/null | head -5 || echo 'none found'",
	}

	for check, cmd := range checks {
		output, err := e.client.Run(cmd)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: Error - %v", check, err))
		} else {
			result := strings.TrimSpace(string(output))
			if result != "" {
				issues = append(issues, fmt.Sprintf("%s: %s", check, result))
			}
		}
	}

	return issues, nil
}

// Helper methods

func (e *SSHExecutor) detectPackageManager() (string, error) {
	managers := []string{"apt-get", "yum", "dnf", "pacman", "zypper"}
	
	for _, manager := range managers {
		cmd := fmt.Sprintf("which %s", manager)
		_, err := e.client.Run(cmd)
		if err == nil {
			if manager == "apt-get" {
				return "apt", nil
			}
			return manager, nil
		}
	}
	
	return "", fmt.Errorf("no supported package manager found")
}

func (e *SSHExecutor) uploadFile(localPath, remotePath string) error {
	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Uploading %s to %s", localPath, remotePath))
	}

	err := e.client.Upload(localPath, remotePath)
	if err != nil {
		return fmt.Errorf("upload failed: %v", err)
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Uploaded %s to %s", localPath, remotePath))
	return nil
}

func (e *SSHExecutor) downloadFile(remotePath, localPath string) error {
	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Downloading %s to %s", remotePath, localPath))
	}

	err := e.client.Download(remotePath, localPath)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Downloaded %s to %s", remotePath, localPath))
	return nil
}

func (e *SSHExecutor) copyFile(source, destination string) error {
	cmd := fmt.Sprintf("cp -r %s %s", source, destination)
	_, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("copy failed: %v", err)
	}
	return nil
}

func (e *SSHExecutor) moveFile(source, destination string) error {
	cmd := fmt.Sprintf("mv %s %s", source, destination)
	_, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("move failed: %v", err)
	}
	return nil
}

func (e *SSHExecutor) deleteFile(path string) error {
	cmd := fmt.Sprintf("rm -rf %s", path)
	_, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}

func (e *SSHExecutor) createDirectory(path string, permissions ...string) error {
	perm := "755"
	if len(permissions) > 0 {
		perm = permissions[0]
	}
	
	cmd := fmt.Sprintf("mkdir -p %s && chmod %s %s", path, perm, path)
	_, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("directory creation failed: %v", err)
	}
	return nil
}

func (e *SSHExecutor) changePermissions(path, permissions string) error {
	cmd := fmt.Sprintf("chmod %s %s", permissions, path)
	_, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("chmod failed: %v", err)
	}
	return nil
}

func (e *SSHExecutor) changeOwnership(path, owner string) error {
	cmd := fmt.Sprintf("chown %s %s", owner, path)
	_, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("chown failed: %v", err)
	}
	return nil
}

// RunWithTimeout executes a command with timeout
func (e *SSHExecutor) RunWithTimeout(command string, timeout time.Duration) ([]byte, error) {
	if !e.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return e.client.RunContext(ctx, command)
}

// GetSFTPClient returns an SFTP client for advanced file operations
func (e *SSHExecutor) GetSFTPClient() (*sftp.Client, error) {
	if !e.connected {
		return nil, fmt.Errorf("not connected to remote host")
	}

	return e.client.NewSftp()
}

// ExecuteScript uploads and executes a local script on the remote server
func (e *SSHExecutor) ExecuteScript(localScriptPath string, args ...string) error {
	if !e.connected {
		return fmt.Errorf("not connected to remote host")
	}

	// Upload script to temporary location
	remotePath := fmt.Sprintf("/tmp/%s", filepath.Base(localScriptPath))
	if err := e.uploadFile(localScriptPath, remotePath); err != nil {
		return fmt.Errorf("failed to upload script: %v", err)
	}

	// Make script executable
	if err := e.changePermissions(remotePath, "755"); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	// Execute script
	cmd := remotePath
	if len(args) > 0 {
		cmd += " " + strings.Join(args, " ")
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing script: %s", cmd))
	}

	output, err := e.client.Run(cmd)
	if err != nil {
		return fmt.Errorf("script execution failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintInfo(fmt.Sprintf("Script output:\n%s", string(output)))

	// Clean up
	e.deleteFile(remotePath)

	return nil
}