package executor

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/SCL/internal/parser"
	"github.com/SCL/internal/utils"
)

// DirectExecutor executes commands directly (interpret mode) - Configuration Management like Ansible
type DirectExecutor struct {
	*parser.BaseInfraDSLListener
	variables     map[string]string
	imports       map[string]bool
	functions     map[string]parser.IFunctionDeclarationContext
	functionCount int
	commandCount  int
	verbose       bool
	currentFunc   string
	sshExecutor   *SSHExecutor
	devopsUtils   *DevOpsUtils
	useSSH        bool
}

func NewDirectExecutor() *DirectExecutor {
	sshExec := NewSSHExecutor()
	return &DirectExecutor{
		variables:   make(map[string]string),
		imports:     make(map[string]bool),
		functions:   make(map[string]parser.IFunctionDeclarationContext),
		verbose:     false,
		sshExecutor: sshExec,
		devopsUtils: NewDevOpsUtils(sshExec),
		useSSH:      false,
	}
}

// EnableSSH enables SSH mode and establishes connection
func (e *DirectExecutor) EnableSSH(host, user, keyPath, passphrase string) error {
	e.sshExecutor.SetVerbose(e.verbose)
	
	var err error
	if keyPath != "" {
		err = e.sshExecutor.ConnectWithKey(host, user, keyPath, passphrase)
	} else {
		err = e.sshExecutor.ConnectWithAgent(host, user)
	}
	
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v", err)
	}
	
	e.useSSH = true
	utils.PrintSuccess("✓ SSH mode enabled")
	return nil
}

// EnableSSHWithPassword enables SSH mode with password authentication
func (e *DirectExecutor) EnableSSHWithPassword(host, user, password string) error {
	e.sshExecutor.SetVerbose(e.verbose)
	
	err := e.sshExecutor.ConnectWithPassword(host, user, password)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v", err)
	}
	
	e.useSSH = true
	utils.PrintSuccess("✓ SSH mode enabled with password")
	return nil
}

// DisableSSH disables SSH mode and closes connection
func (e *DirectExecutor) DisableSSH() error {
	if e.useSSH && e.sshExecutor != nil {
		err := e.sshExecutor.Close()
		e.useSSH = false
		utils.PrintInfo("SSH mode disabled")
		return err
	}
	return nil
}

func (e *DirectExecutor) SetVerbose(v bool) {
	e.verbose = v
}

func (e *DirectExecutor) GetTarget() string {
	return e.variables["target"]
}

func (e *DirectExecutor) CheckSSHConnectivity(target string) error {
	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Testing SSH connection to %s...", target))
	}

	// Check if target includes port (for Docker testing)
	var cmd *exec.Cmd
	if strings.Contains(target, ":") {
		// Handle custom port (e.g., testuser@localhost:2222)
		parts := strings.Split(target, ":")
		host := parts[0]
		port := parts[1]
		cmd = exec.Command("ssh", "-p", port, "-o", "ConnectTimeout=5", "-o", "BatchMode=yes", host, "echo 'SSH connection test successful'")
	} else {
		// Standard SSH (port 22)
		cmd = exec.Command("ssh", "-o", "ConnectTimeout=5", "-o", "BatchMode=yes", target, "echo 'SSH connection test successful'")
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		errorMsg := fmt.Sprintf("SSH connection failed: %v", err)
		if len(output) > 0 {
			errorMsg += fmt.Sprintf("\nOutput: %s", string(output))
		}

		// Add helpful error explanations
		if strings.Contains(err.Error(), "exit status 255") {
			errorMsg += "\nTroubleshooting:"
			errorMsg += "\n  • Check if SSH server is running on target"
			errorMsg += "\n  • Verify SSH key authentication is set up"
			if strings.Contains(target, ":") {
				parts := strings.Split(target, ":")
				errorMsg += "\n  • Try: ssh-copy-id -p " + parts[1] + " " + parts[0]
				errorMsg += "\n  • Or test manually: ssh -p " + parts[1] + " " + parts[0]
			} else {
				errorMsg += "\n  • Try: ssh-copy-id " + target
				errorMsg += "\n  • Or test manually: ssh " + target
			}
			errorMsg += "\n  • For Docker testing: Make sure container is running on port 2222"
		}

		return fmt.Errorf(errorMsg)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("SSH test output: %s", string(output)))
	}

	return nil
}

func (e *DirectExecutor) GetVariables() map[string]string {
	return e.variables
}

func (e *DirectExecutor) GetImports() map[string]bool {
	return e.imports
}

func (e *DirectExecutor) GetFunctionCount() int {
	return e.functionCount
}

func (e *DirectExecutor) GetCommandCount() int {
	return e.commandCount
}

func (e *DirectExecutor) EnterImportStatement(ctx *parser.ImportStatementContext) {
	moduleName := ctx.IDENTIFIER().GetText()
	e.imports[moduleName] = true

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Importing module: %s", moduleName))
	}
}

func (e *DirectExecutor) EnterAssignment(ctx *parser.AssignmentContext) {
	varName := ctx.IDENTIFIER().GetText()
	expr := ctx.Expression().GetText()
	value := e.cleanExpression(expr)
	e.variables[varName] = value

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Assignment: %s = %s", varName, value))
	}
}

func (e *DirectExecutor) EnterDeclaration(ctx *parser.DeclarationContext) {
	varName := ctx.IDENTIFIER().GetText()
	expr := ctx.Expression().GetText()

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Declaring variable: %s", varName))
	}

	// Handle special cases like check() function
	if strings.Contains(expr, "check(") {
		result := e.executeCheckCommand(varName, expr)
		e.variables[varName] = fmt.Sprintf("%t", result)
	} else {
		value := e.cleanExpression(expr)
		e.variables[varName] = value
	}
}

func (e *DirectExecutor) EnterFunctionDeclaration(ctx *parser.FunctionDeclarationContext) {
	funcName := ctx.IDENTIFIER().GetText()
	e.functions[funcName] = ctx
	e.functionCount++

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Registered function: %s()", funcName))
	}
}

// Execute runs the main function
func (e *DirectExecutor) Execute() error {
	mainFunc, exists := e.functions["main"]
	if !exists {
		return fmt.Errorf("main function not found")
	}

	utils.PrintInfo("Executing main function...")
	return e.executeFunction("main", mainFunc)
}

func (e *DirectExecutor) executeFunction(name string, ctx parser.IFunctionDeclarationContext) error {
	e.currentFunc = name

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing function: %s()", name))
	}

	// Execute function body
	block := ctx.Block()
	for _, stmt := range block.AllStatement() {
		if err := e.executeStatement(stmt); err != nil {
			return fmt.Errorf("error in function %s: %v", name, err)
		}
	}

	return nil
}

func (e *DirectExecutor) executeStatement(stmt parser.IStatementContext) error {
	// Use the interface methods to get specific statement types
	if exprStmt := stmt.ExpressionStatement(); exprStmt != nil {
		return e.executeExpressionStatement(exprStmt)
	}

	if ifStmt := stmt.IfStatement(); ifStmt != nil {
		return e.executeIfStatement(ifStmt)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Skipping statement type: %T", stmt))
	}
	return nil
}

func (e *DirectExecutor) executeExpressionStatement(ctx parser.IExpressionStatementContext) error {
	qualifiedName := ctx.QualifiedName().GetText()

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Executing: %s()", qualifiedName))
	}

	return e.handleFunctionCall(qualifiedName, ctx.ArgumentList())
}

func (e *DirectExecutor) executeIfStatement(ctx parser.IIfStatementContext) error {
	condition := ctx.Expression().GetText()
	conditionResult := e.evaluateCondition(condition)

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("If condition '%s' = %t", condition, conditionResult))
	}

	var blockToExecute parser.IBlockContext
	if conditionResult {
		blockToExecute = ctx.AllBlock()[0] // then block
	} else if len(ctx.AllBlock()) > 1 {
		blockToExecute = ctx.AllBlock()[1] // else block
	}

	if blockToExecute != nil {
		for _, stmt := range blockToExecute.AllStatement() {
			if err := e.executeStatement(stmt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *DirectExecutor) handleFunctionCall(qualifiedName string, argList parser.IArgumentListContext) error {
	e.commandCount++

	switch {
	case qualifiedName == "test":
		return e.executeTestCommand()

	case qualifiedName == "print" || qualifiedName == "primary.print":
		return e.executePrintCommand(argList)

	case qualifiedName == "copy":
		return e.executeCopyCommand(argList)

	case qualifiedName == "create":
		return e.executeCreateCommand(argList)

	case qualifiedName == "install":
		return e.executeInstallCommand(argList)

	// Enhanced SSH-based DevOps functions
	case qualifiedName == "sysinfo":
		return e.executeSystemInfoCommand()

	case qualifiedName == "monitor":
		return e.executeMonitorCommand()

	case qualifiedName == "service":
		return e.executeServiceCommand(argList)

	case qualifiedName == "package":
		return e.executePackageCommand(argList)

	case qualifiedName == "docker":
		return e.executeDockerCommand(argList)

	case qualifiedName == "backup":
		return e.executeBackupCommand(argList)

	case qualifiedName == "firewall":
		return e.executeFirewallCommand(argList)

	case qualifiedName == "user":
		return e.executeUserCommand(argList)

	case qualifiedName == "cert":
		return e.executeCertCommand(argList)

	case qualifiedName == "cron":
		return e.executeCronCommand(argList)

	case qualifiedName == "audit":
		return e.executeAuditCommand()

	case qualifiedName == "nettest":
		return e.executeNetworkTestCommand(argList)

	case qualifiedName == "logs":
		return e.executeLogAnalysisCommand(argList)

	case qualifiedName == "tune":
		return e.executeSystemTuningCommand(argList)

	case qualifiedName == "webserver":
		return e.executeWebServerCommand(argList)

	case qualifiedName == "database":
		return e.executeDatabaseCommand(argList)

	default:
		// Check if it's a user-defined function
		if funcCtx, exists := e.functions[qualifiedName]; exists {
			return e.executeFunction(qualifiedName, funcCtx)
		}

		if e.verbose {
			utils.PrintInfo(fmt.Sprintf("Unknown function call: %s()", qualifiedName))
		}
	}

	return nil
}

func (e *DirectExecutor) executeCheckCommand(varName, expr string) bool {
	target := e.variables["target"]

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Checking tools on remote server: %s", target))
	}

	// Always use SSH for remote configuration management
	var cmd *exec.Cmd
	if strings.Contains(target, ":") {
		// Handle custom port
		parts := strings.Split(target, ":")
		host := parts[0]
		port := parts[1]
		cmd = exec.Command("ssh", "-p", port, host, "command -v curl && command -v wget")
	} else {
		cmd = exec.Command("ssh", target, "command -v curl && command -v wget")
	}

	output, err := cmd.CombinedOutput()
	result := err == nil

	if e.verbose {
		if result {
			utils.PrintInfo(fmt.Sprintf("Check result for %s: %t", varName, result))
			if len(output) > 0 {
				utils.PrintInfo(fmt.Sprintf("Command output: %s", string(output)))
			}
		} else {
			utils.PrintInfo(fmt.Sprintf("Check result for %s: %t", varName, result))
			if err != nil {
				utils.PrintInfo(fmt.Sprintf("Check failed: %v", err))
			}
			if len(output) > 0 {
				utils.PrintInfo(fmt.Sprintf("Error output: %s", string(output)))
			}
		}
	}

	return result
}

func (e *DirectExecutor) executeTestCommand() error {
	if e.verbose {
		utils.PrintInfo("Executing test command")
	}

	utils.PrintSuccess("✓ Test access completed")
	return nil
}

func (e *DirectExecutor) executePrintCommand(argList parser.IArgumentListContext) error {
	if argList != nil {
		args := argList.GetText()
		message := strings.Trim(args, "\"")

		utils.PrintSuccess(fmt.Sprintf("📢 %s", message))
	}
	return nil
}

func (e *DirectExecutor) executeCopyCommand(argList parser.IArgumentListContext) error {
	if argList == nil {
		return fmt.Errorf("copy command requires arguments: copy(source, destination)")
	}

	target := e.variables["target"]
	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("copy command requires source and destination: copy(source, destination)")
	}

	source := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	dest := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Copying %s to %s on remote server: %s", source, dest, target))
	}

	// First, create the destination directory on remote server
	destDir := dest
	if !strings.HasSuffix(dest, "/") {
		// If dest is a file path, get the directory part
		lastSlash := strings.LastIndex(dest, "/")
		if lastSlash > 0 {
			destDir = dest[:lastSlash]
		}
	}

	// Create destination directory
	createDirCmd := fmt.Sprintf("mkdir -p %s", destDir)
	var mkdirCmd *exec.Cmd
	if strings.Contains(target, ":") {
		// Handle custom port
		targetParts := strings.Split(target, ":")
		host := targetParts[0]
		port := targetParts[1]
		mkdirCmd = exec.Command("ssh", "-p", port, host, createDirCmd)
	} else {
		mkdirCmd = exec.Command("ssh", target, createDirCmd)
	}

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Creating destination directory: %s", destDir))
	}

	mkdirOutput, mkdirErr := mkdirCmd.CombinedOutput()
	if mkdirErr != nil {
		errorMsg := fmt.Sprintf("failed to create destination directory: %v", mkdirErr)
		if len(mkdirOutput) > 0 {
			errorMsg += fmt.Sprintf("\nCommand output: %s", string(mkdirOutput))
		}
		return fmt.Errorf(errorMsg)
	}

	// Now copy the file using scp
	var scpCmd *exec.Cmd
	if strings.Contains(target, ":") {
		// Handle custom port
		targetParts := strings.Split(target, ":")
		host := targetParts[0]
		port := targetParts[1]
		scpCmd = exec.Command("scp", "-P", port, "-r", source, fmt.Sprintf("%s:%s", host, dest))
	} else {
		scpCmd = exec.Command("scp", "-r", source, fmt.Sprintf("%s:%s", target, dest))
	}

	output, err := scpCmd.CombinedOutput()
	if err != nil {
		errorMsg := fmt.Sprintf("file copy failed: %v", err)
		if len(output) > 0 {
			errorMsg += fmt.Sprintf("\nCommand output: %s", string(output))
		}
		return fmt.Errorf(errorMsg)
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Copied %s to %s:%s", source, target, dest))
	if e.verbose && len(output) > 0 {
		utils.PrintInfo(fmt.Sprintf("Command output: %s", string(output)))
	}
	return nil
}

func (e *DirectExecutor) executeCreateCommand(argList parser.IArgumentListContext) error {
	if argList == nil {
		return fmt.Errorf("create command requires arguments: create(destination, filename, mode)")
	}

	target := e.variables["target"]
	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 3 {
		return fmt.Errorf("create command requires destination, filename, and mode: create(destination, filename, mode)")
	}

	destination := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	filename := strings.Trim(strings.TrimSpace(parts[1]), "\"")
	mode := strings.Trim(strings.TrimSpace(parts[2]), "\"")

	fullPath := fmt.Sprintf("%s/%s", destination, filename)

	if e.verbose {
		utils.PrintInfo(fmt.Sprintf("Creating file %s with mode %s on remote server: %s", fullPath, mode, target))
	}

	createCmd := fmt.Sprintf("mkdir -p %s && touch %s && chmod %s %s", destination, fullPath, mode, fullPath)

	var cmd *exec.Cmd
	if strings.Contains(target, ":") {
		targetParts := strings.Split(target, ":")
		host := targetParts[0]
		port := targetParts[1]
		cmd = exec.Command("ssh", "-p", port, host, createCmd)
	} else {
		cmd = exec.Command("ssh", target, createCmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := fmt.Sprintf("file creation failed: %v", err)
		if len(output) > 0 {
			errorMsg += fmt.Sprintf("\nCommand output: %s", string(output))
		}
		return fmt.Errorf(errorMsg)
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Created file %s with mode %s on %s", fullPath, mode, target))
	if e.verbose && len(output) > 0 {
		utils.PrintInfo(fmt.Sprintf("Command output: %s", string(output)))
	}
	return nil
}

func (e *DirectExecutor) executeInstallCommand(argList parser.IArgumentListContext) error {
	if argList == nil {
		return fmt.Errorf("install command requires arguments")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")
	if len(parts) >= 2 {
		pkg := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		target := e.variables["target"]

		if e.verbose {
			utils.PrintInfo(fmt.Sprintf("Installing %s on remote server: %s", pkg, target))
		}

		// Always use SSH for remote configuration management
		cmdStr := fmt.Sprintf("apt-get update && apt-get install -y %s", pkg)
		var cmd *exec.Cmd
		if strings.Contains(target, ":") {
			// Handle custom port
			parts := strings.Split(target, ":")
			host := parts[0]
			port := parts[1]
			cmd = exec.Command("ssh", "-p", port, host, cmdStr)
		} else {
			cmd = exec.Command("ssh", target, cmdStr)
		}

		if e.verbose {
			utils.PrintInfo(fmt.Sprintf("Executing remote command: %s", cmdStr))
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			errorMsg := fmt.Sprintf("remote install %s failed: %v", pkg, err)
			if len(output) > 0 {
				errorMsg += fmt.Sprintf("\nCommand output: %s", string(output))
			}

			// Add specific error explanations
			if strings.Contains(err.Error(), "exit status 255") {
				errorMsg += "\nReason: SSH connection failed - remote server may be unreachable or SSH key not configured"
			} else if strings.Contains(err.Error(), "connection refused") {
				errorMsg += "\nReason: Connection refused - remote server is not accepting SSH connections"
			} else if strings.Contains(err.Error(), "no route to host") {
				errorMsg += "\nReason: Network unreachable - check network connectivity to remote server"
			} else if strings.Contains(err.Error(), "permission denied") {
				errorMsg += "\nReason: Permission denied - check SSH credentials and user permissions"
			}

			return fmt.Errorf(errorMsg)
		}

		utils.PrintSuccess(fmt.Sprintf("✓ Installed %s on remote server %s", pkg, target))
		if e.verbose && len(output) > 0 {
			utils.PrintInfo(fmt.Sprintf("Command output: %s", string(output)))
		}
	}

	return nil
}

func (e *DirectExecutor) evaluateCondition(condition string) bool {
	// Simple condition evaluation - check if variable exists and is "true" or exit code 0
	if val, exists := e.variables[condition]; exists {
		return val == "true" || val == "0"
	}
	return false
}

func (e *DirectExecutor) cleanExpression(expr string) string {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "[") && strings.HasSuffix(expr, "]") {
		content := strings.Trim(expr, "[]")
		parts := strings.Split(content, ",")
		if len(parts) > 0 {
			return strings.Trim(strings.TrimSpace(parts[0]), "\"")
		}
	}

	return strings.Trim(expr, "\"")
}

// Enhanced SSH-based DevOps command implementations

func (e *DirectExecutor) executeSystemInfoCommand() error {
	if !e.useSSH {
		utils.PrintInfo("System info command requires SSH mode. Use local system info instead.")
		return e.executeLocalSystemInfo()
	}

	if e.verbose {
		utils.PrintInfo("Retrieving remote system information...")
	}

	sysInfo, err := e.sshExecutor.SystemInfo()
	if err != nil {
		return fmt.Errorf("failed to get system info: %v", err)
	}

	utils.PrintSuccess("=== Remote System Information ===")
	for key, value := range sysInfo {
		utils.PrintInfo(fmt.Sprintf("%s: %s", key, value))
	}

	return nil
}

func (e *DirectExecutor) executeLocalSystemInfo() error {
	commands := map[string]string{
		"hostname":     "hostname",
		"os":          "uname -s",
		"kernel":      "uname -r",
		"architecture": "uname -m",
		"uptime":      "uptime",
	}

	utils.PrintSuccess("=== Local System Information ===")
	for key, cmd := range commands {
		output, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			utils.PrintInfo(fmt.Sprintf("%s: Error - %v", key, err))
		} else {
			utils.PrintInfo(fmt.Sprintf("%s: %s", key, strings.TrimSpace(string(output))))
		}
	}

	return nil
}

func (e *DirectExecutor) executeMonitorCommand() error {
	if !e.useSSH {
		return fmt.Errorf("monitor command requires SSH mode")
	}

	if e.verbose {
		utils.PrintInfo("Collecting monitoring metrics...")
	}

	metrics, err := e.sshExecutor.MonitoringMetrics()
	if err != nil {
		return fmt.Errorf("failed to get monitoring metrics: %v", err)
	}

	utils.PrintSuccess("=== Monitoring Metrics ===")
	for key, value := range metrics {
		utils.PrintInfo(fmt.Sprintf("%s: %v", key, value))
	}

	return nil
}

func (e *DirectExecutor) executeServiceCommand(argList parser.IArgumentListContext) error {
	if argList == nil {
		return fmt.Errorf("service command requires arguments: service(action, service_name)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("service command requires action and service name: service(action, service_name)")
	}

	action := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	serviceName := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	if e.useSSH {
		return e.sshExecutor.ProcessManagement(action, serviceName)
	}

	// Local execution fallback
	var cmd *exec.Cmd
	switch action {
	case "start", "stop", "restart", "status", "enable", "disable":
		cmd = exec.Command("systemctl", action, serviceName)
	default:
		return fmt.Errorf("unsupported service action: %s", action)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("service operation failed: %v\nOutput: %s", err, string(output))
	}

	utils.PrintSuccess(fmt.Sprintf("✓ Service %s %s completed", serviceName, action))
	return nil
}

func (e *DirectExecutor) executePackageCommand(argList parser.IArgumentListContext) error {
	if argList == nil {
		return fmt.Errorf("package command requires arguments: package(action, package_name)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("package command requires action and package name: package(action, package_name)")
	}

	action := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	packageName := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	if e.useSSH {
		return e.sshExecutor.PackageManagement(action, packageName)
	}

	// Local execution fallback
	return e.executeInstallCommand(argList)
}

func (e *DirectExecutor) executeDockerCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("docker command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("docker command requires arguments: docker(action, ...args)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 1 {
		return fmt.Errorf("docker command requires at least an action")
	}

	action := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	var dockerArgs []string

	for i := 1; i < len(parts); i++ {
		dockerArgs = append(dockerArgs, strings.Trim(strings.TrimSpace(parts[i]), "\""))
	}

	return e.sshExecutor.DockerOperations(action, dockerArgs...)
}

func (e *DirectExecutor) executeBackupCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("backup command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("backup command requires arguments: backup(operation, source, destination)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 3 {
		return fmt.Errorf("backup command requires operation, source, and destination")
	}

	operation := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	source := strings.Trim(strings.TrimSpace(parts[1]), "\"")
	destination := strings.Trim(strings.TrimSpace(parts[2]), "\"")

	var options []string
	if len(parts) > 3 {
		options = append(options, strings.Trim(strings.TrimSpace(parts[3]), "\""))
	}

	return e.devopsUtils.BackupOperations(operation, source, destination, options...)
}

func (e *DirectExecutor) executeFirewallCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("firewall command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("firewall command requires arguments: firewall(operation, ...params)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 1 {
		return fmt.Errorf("firewall command requires at least an operation")
	}

	operation := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	var params []string

	for i := 1; i < len(parts); i++ {
		params = append(params, strings.Trim(strings.TrimSpace(parts[i]), "\""))
	}

	return e.devopsUtils.FirewallManagement(operation, params...)
}

func (e *DirectExecutor) executeUserCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("user command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("user command requires arguments: user(operation, username, ...params)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("user command requires operation and username")
	}

	operation := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	username := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	var params []string
	for i := 2; i < len(parts); i++ {
		params = append(params, strings.Trim(strings.TrimSpace(parts[i]), "\""))
	}

	return e.devopsUtils.UserManagement(operation, username, params...)
}

func (e *DirectExecutor) executeCertCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("cert command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("cert command requires arguments: cert(operation, domain, ...options)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("cert command requires operation and domain")
	}

	operation := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	domain := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	var options []string
	for i := 2; i < len(parts); i++ {
		options = append(options, strings.Trim(strings.TrimSpace(parts[i]), "\""))
	}

	return e.devopsUtils.CertificateManagement(operation, domain, options...)
}

func (e *DirectExecutor) executeCronCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("cron command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("cron command requires arguments: cron(operation, schedule, command, user)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 1 {
		return fmt.Errorf("cron command requires at least an operation")
	}

	operation := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	
	var schedule, command string
	var user []string

	if len(parts) > 1 {
		schedule = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	}
	if len(parts) > 2 {
		command = strings.Trim(strings.TrimSpace(parts[2]), "\"")
	}
	if len(parts) > 3 {
		user = append(user, strings.Trim(strings.TrimSpace(parts[3]), "\""))
	}

	return e.devopsUtils.CronManagement(operation, schedule, command, user...)
}

func (e *DirectExecutor) executeAuditCommand() error {
	if !e.useSSH {
		return fmt.Errorf("audit command requires SSH mode")
	}

	if e.verbose {
		utils.PrintInfo("Performing security audit...")
	}

	issues, err := e.sshExecutor.SecurityAudit()
	if err != nil {
		return fmt.Errorf("security audit failed: %v", err)
	}

	utils.PrintSuccess("=== Security Audit Results ===")
	for _, issue := range issues {
		utils.PrintInfo(issue)
	}

	return nil
}

func (e *DirectExecutor) executeNetworkTestCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("nettest command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("nettest command requires arguments: nettest(target, test_type)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("nettest command requires target and test type")
	}

	target := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	testType := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	return e.sshExecutor.NetworkDiagnostics(target, testType)
}

func (e *DirectExecutor) executeLogAnalysisCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("logs command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("logs command requires arguments: logs(logfile, pattern, lines)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 1 {
		return fmt.Errorf("logs command requires at least a log file")
	}

	logFile := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	pattern := ""
	lines := 10

	if len(parts) > 1 {
		pattern = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	}
	if len(parts) > 2 {
		if l, err := fmt.Sscanf(strings.TrimSpace(parts[2]), "%d", &lines); err != nil || l != 1 {
			lines = 10
		}
	}

	return e.sshExecutor.LogAnalysis(logFile, pattern, lines)
}

func (e *DirectExecutor) executeSystemTuningCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("tune command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("tune command requires arguments: tune(operation, ...params)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 1 {
		return fmt.Errorf("tune command requires at least an operation")
	}

	operation := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	var params []string

	for i := 1; i < len(parts); i++ {
		params = append(params, strings.Trim(strings.TrimSpace(parts[i]), "\""))
	}

	return e.devopsUtils.SystemTuning(operation, params...)
}

func (e *DirectExecutor) executeWebServerCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("webserver command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("webserver command requires arguments: webserver(server_type, operation)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("webserver command requires server type and operation")
	}

	serverType := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	operation := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	var configPath []string
	if len(parts) > 2 {
		configPath = append(configPath, strings.Trim(strings.TrimSpace(parts[2]), "\""))
	}

	return e.devopsUtils.WebServerOperations(serverType, operation, configPath...)
}

func (e *DirectExecutor) executeDatabaseCommand(argList parser.IArgumentListContext) error {
	if !e.useSSH {
		return fmt.Errorf("database command requires SSH mode")
	}

	if argList == nil {
		return fmt.Errorf("database command requires arguments: database(db_type, operation, params...)")
	}

	args := argList.GetText()
	parts := strings.Split(args, ",")

	if len(parts) < 2 {
		return fmt.Errorf("database command requires database type and operation")
	}

	dbType := strings.Trim(strings.TrimSpace(parts[0]), "\"")
	operation := strings.Trim(strings.TrimSpace(parts[1]), "\"")

	// Parse remaining parameters as key=value pairs
	params := make(map[string]string)
	for i := 2; i < len(parts); i++ {
		param := strings.TrimSpace(parts[i])
		if strings.Contains(param, "=") {
			kv := strings.SplitN(param, "=", 2)
			key := strings.Trim(kv[0], "\"")
			value := strings.Trim(kv[1], "\"")
			params[key] = value
		}
	}

	return e.devopsUtils.DatabaseOperations(dbType, operation, params)
}