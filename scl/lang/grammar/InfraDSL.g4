grammar InfraDSL;

// Parser Rules
program
    : importStatement* statement* EOF
    ;

importStatement
    : 'import' IDENTIFIER ';'
    ;

statement
    : assignment
    | declaration
    | functionDeclaration
    | expressionStatement
    | ifStatement
    ;

assignment
    : IDENTIFIER ':=' expression ';'          // shell := bash;
    | IDENTIFIER '=' expression ';'            // super_user = true;
    ;

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

expressionStatement
    : qualifiedName '(' argumentList? ')' ';'    // Function calls as statements
    ;

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

qualifiedName
    : IDENTIFIER ('.' IDENTIFIER)*
    ;

argumentList
    : expression (',' expression)* ','?    // Allow trailing comma
    ;

ifStatement
    : 'if' expression block ('else' block)?
    ;

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

// Lexer Rules
// Keywords (must come before IDENTIFIER)
IF      : 'if';
ELSE    : 'else';
FN      : 'fn';
DECLARE : 'declare';
IMPORT  : 'import';
BOOL_TYPE    : 'bool';
STRING_TYPE  : 'string';
INT_TYPE     : 'int';
FLOAT_TYPE   : 'float';
LIST_TYPE    : 'list';
BOOLEAN : 'true' | 'false';

// Operators
ASSIGN : ':=';
EQ  : '==';
NEQ : '!=';
LTE : '<=';
GTE : '>=';
LT  : '<';
GT  : '>';
AND : '&&';
OR  : '||';
NOT : '!';
PLUS  : '+';
MINUS : '-';
MULT  : '*';
DIV   : '/';
MOD   : '%';

// Delimiters
LPAREN : '(';
RPAREN : ')';
LBRACE : '{';
RBRACE : '}';
LBRACK : '[';
RBRACK : ']';
SEMICOLON : ';';
COMMA : ',';
COLON : ':';
DOT   : '.';
EQUALS : '=';

// Literals
NUMBER  : [0-9]+ ('.' [0-9]+)?;
STRING  : '"' (~["\\\r\n] | '\\' .)* '"';

// Identifiers (must come after keywords)
IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]*;

// Comments
LINE_COMMENT  : '//' ~[\r\n]* -> skip;
BLOCK_COMMENT : '/*' .*? '*/' -> skip;

// Whitespace
WS : [ \t\r\n]+ -> skip;