# SCL Syntax Guide

## Complete Token Reference

### Keywords
| Token | Description | Example |
|-------|-------------|---------|
| `if` | Conditional statement | `if (condition) { ... }` |
| `else` | Alternative branch | `if (x) { ... } else { ... }` |
| `fn` | Function declaration | `fn main() { ... }` |
| `declare` | Variable declaration | `declare x: int = 5;` |
| `import` | Module import | `import primary;` |

### Data Type Keywords
| Token | Description | Example |
|-------|-------------|---------|
| `bool` | Boolean type | `declare flag: bool = true;` |
| `string` | String type | `declare name: string = "test";` |
| `int` | Integer type | `declare count: int = 42;` |
| `float` | Float type | `declare rate: float = 3.14;` |
| `list` | Array/List type | `declare items: list = [1, 2, 3];` |

### Boolean Literals
| Token | Description | Example |
|-------|-------------|---------|
| `true` | Boolean true | `ssl_enabled = true;` |
| `false` | Boolean false | `debug_mode = false;` |

### Assignment Operators
| Token | Description | Example |
|-------|-------------|---------|
| `:=` | Variable assignment | `mode := "interpret";` |
| `=` | Value assignment | `super_user = true;` |

### Comparison Operators
| Token | Description | Example |
|-------|-------------|---------|
| `==` | Equal to | `if (env == "prod") { ... }` |
| `!=` | Not equal to | `if (status != "ok") { ... }` |
| `<` | Less than | `if (port < 1024) { ... }` |
| `>` | Greater than | `if (count > 100) { ... }` |
| `<=` | Less than or equal | `if (cpu <= 80) { ... }` |
| `>=` | Greater than or equal | `if (memory >= 4) { ... }` |

### Logical Operators
| Token | Description | Example |
|-------|-------------|---------|
| `&&` | Logical AND | `if (ssl && prod) { ... }` |
| `\|\|` | Logical OR | `if (dev \|\| test) { ... }` |
| `!` | Logical NOT | `if (!disabled) { ... }` |

### Arithmetic Operators
| Token | Description | Example |
|-------|-------------|---------|
| `+` | Addition | `total = base + tax;` |
| `-` | Subtraction | `remaining = total - used;` |
| `*` | Multiplication | `area = width * height;` |
| `/` | Division | `average = sum / count;` |
| `%` | Modulo | `remainder = num % 10;` |

### Delimiters
| Token | Description | Example |
|-------|-------------|---------|
| `(` | Left parenthesis | `fn main() {` |
| `)` | Right parenthesis | `print("hello");` |
| `{` | Left brace | `if (x) { ... }` |
| `}` | Right brace | `fn main() { ... }` |
| `[` | Left bracket | `servers = ["web1", "web2"];` |
| `]` | Right bracket | `ports = [80, 443];` |
| `;` | Semicolon | `mode := "interpret";` |
| `,` | Comma | `copy("src", "dest");` |
| `:` | Colon | `declare x: int = 5;` |
| `.` | Dot | `primary.print("msg");` |

## Grammar Rules (ANTLR4)

### Program Structure
```antlr
program
    : importStatement* statement* EOF
    ;
```

### Import Statements
```antlr
importStatement
    : 'import' IDENTIFIER ';'
    ;
```

### Statements
```antlr
statement
    : assignment
    | declaration
    | functionDeclaration
    | expressionStatement
    | ifStatement
    ;
```

### Variable Assignment
```antlr
assignment
    : IDENTIFIER ':=' expression ';'          // shell := bash;
    | IDENTIFIER '=' expression ';'            // super_user = true;
    ;
```

### Variable Declaration
```antlr
declaration
    : 'declare' IDENTIFIER ':' type '=' expression ';'
    ;

type
    : 'bool'
    | 'string'
    | 'int'
    | 'float'
    | 'list'
    ;
```

### Function Declaration
```antlr
functionDeclaration
    : 'fn' IDENTIFIER '(' parameterList? ')' block
    ;

parameterList
    : parameter (',' parameter)*
    ;

parameter
    : IDENTIFIER ':' type
    ;

block
    : '{' statement* '}'
    ;
```

### Expression Statement
```antlr
expressionStatement
    : qualifiedName '(' argumentList? ')' ';'    // Function calls as statements
    ;
```

### Expressions
```antlr
expression
    : primary                                                           # PrimaryExpr
    | expression '.' IDENTIFIER                                         # MemberAccessExpr
    | expression '(' argumentList? ')'                                  # FunctionCallExpr
    | expression op=('*' | '/' | '%') expression                       # MulDivModExpr
    | expression op=('+' | '-') expression                             # AddSubExpr
    | expression op=('==' | '!=' | '<' | '>' | '<=' | '>=') expression # ComparisonExpr
    | expression op=('&&' | '||') expression                           # LogicalExpr
    | '!' expression                                                    # NotExpr
    ;
```

### Qualified Names and Arguments
```antlr
qualifiedName
    : IDENTIFIER ('.' IDENTIFIER)*
    ;

argumentList
    : expression (',' expression)* ','?    // Allow trailing comma
    ;
```

### If Statements
```antlr
ifStatement
    : 'if' expression block ('else' block)?
    ;
```

### Primary Expressions
```antlr
primary
    : IDENTIFIER
    | STRING
    | NUMBER
    | BOOLEAN
    | array
    | '(' expression ')'
    ;

array
    : '[' (expression (',' expression)*)? ']'
    ;
```

## Lexical Rules

### Literals
```antlr
NUMBER  : [0-9]+ ('.' [0-9]+)?;
STRING  : '"' (~["\\\r\n] | '\\' .)* '"';
BOOLEAN : 'true' | 'false';
```

### Identifiers
```antlr
IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]*;
```

### Comments
```antlr
LINE_COMMENT  : '//' ~[\r\n]* -> skip;
BLOCK_COMMENT : '/*' .*? '*/' -> skip;
```

### Whitespace
```antlr
WS : [ \t\r\n]+ -> skip;
```

## Complete Syntax Examples

### Basic Program Structure
```scl
// Comments are supported
import primary;

// Configuration variables
mode := "interpret";
setting := "configuration";
super_user = true;
on_error = "rollback";
target := "user@server.com";

// Variable declarations with types
declare environment: string = "production";
declare port: int = 8080;
declare ssl_enabled: bool = true;
declare servers: list = ["web1", "web2", "web3"];

// Main function
fn main() {
    print("Starting deployment");
    
    if (ssl_enabled) {
        setup_ssl();
    } else {
        setup_http();
    }
    
    deploy_application();
}

fn setup_ssl() {
    cert("letsencrypt", "example.com");
    firewall("allow", "443/tcp");
}

fn setup_http() {
    firewall("allow", "80/tcp");
}

fn deploy_application() {
    package("install", "nginx");
    service("start", "nginx");
}
```

### Expression Examples
```scl
// Arithmetic expressions
declare total: int = base + tax * rate;
declare average: float = sum / count;

// Comparison expressions
if (port > 1024) {
    print("Non-privileged port");
}

if (environment == "production") {
    backup("backup", "/data", "/backup/data.tar.gz");
}

// Logical expressions
if (ssl_enabled && environment == "production") {
    cert("letsencrypt", domain);
}

if (debug_mode || verbose_logging) {
    logs("/var/log/app.log", "DEBUG", 100);
}

// Complex expressions
if ((cpu_usage > 80.0) && (memory_usage > 90.0)) {
    print("High resource usage detected");
    tune("swappiness", "10");
}
```

### Array and Function Call Examples
```scl
// Array literals
declare web_servers: list = ["web1.example.com", "web2.example.com"];
declare ports: list = [80, 443, 8080];
declare mixed_array: list = ["string", 42, true];

// Function calls with various argument patterns
print("Simple message");
package("install", "nginx");
create("/var/log/app", "app.log", "644");
nettest("google.com", "ping");

// Qualified function calls
primary.print("Using primary module");

// Function calls with trailing commas (allowed)
copy(
    "source.txt",
    "destination.txt",
);
```

### Control Flow Examples
```scl
// Simple if statement
if (tools_ready) {
    print("Tools are available");
}

// If-else statement
if (ssl_enabled) {
    print("SSL is enabled");
    cert("check_expiry", domain);
} else {
    print("SSL is disabled");
    firewall("allow", "80/tcp");
}

// Nested if statements
if (environment == "production") {
    if (ssl_enabled) {
        cert("letsencrypt", domain);
        firewall("allow", "443/tcp");
    } else {
        firewall("allow", "80/tcp");
    }
    
    if (backup_enabled) {
        backup("backup", "/data", "/backup/data.tar.gz");
    }
}

// Complex boolean expressions
if ((environment == "production") && (ssl_enabled || force_ssl)) {
    print("Setting up production SSL");
    cert("letsencrypt", domain);
}
```

### Built-in Function Examples
```scl
// System operations
sysinfo();                              // Get system information
monitor();                              // Monitor system metrics
audit();                                // Security audit

// Package management
package("install", "nginx");            // Install package
package("update", "");                  // Update packages
package("search", "docker");            // Search packages

// Service management
service("start", "nginx");              // Start service
service("stop", "apache2");             // Stop service
service("restart", "mysql");            // Restart service
service("enable", "nginx");             // Enable service
service("status", "docker");            // Check service status

// File operations
create("/var/log/app", "app.log", "644");   // Create file
copy("nginx.conf", "/etc/nginx/nginx.conf"); // Copy file

// Network operations
nettest("google.com", "ping");          // Ping test
nettest("example.com:443", "telnet");   // Port test
nettest("https://api.example.com", "curl"); // HTTP test

// Security operations
firewall("allow", "80/tcp");            // Allow port
firewall("deny", "23/tcp");             // Deny port
user("add", "devops", "sudo");          // Add user

// Container operations
docker("ps");                           // List containers
docker("start", "web-container");       // Start container
docker("logs", "web-container");        // View logs

// Advanced operations
cert("letsencrypt", "example.com");     // SSL certificate
cron("add", "0 2 * * *", "/backup.sh", "root"); // Cron job
backup("backup", "/etc", "/backup/etc.tar.gz"); // Backup
logs("/var/log/syslog", "error", 10);   // Log analysis
tune("swappiness", "10");               // System tuning
```

This comprehensive syntax guide covers all tokens, grammar rules, and syntax patterns available in the SCL language, providing a complete reference for developers using the language.