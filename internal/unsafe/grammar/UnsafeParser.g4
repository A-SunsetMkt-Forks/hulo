parser grammar unsafeParser;

options {
	tokenVocab = unsafeLexer;
}

template: (content | statement)*;

content: TEXT;

statement:
	ifStatement
	| loopStatement
	| expressionStatement
	| defineStatement
	| templateStatement
	| variableStatement;

variableStatement: IDENTIFIER COLON_EQUAL expression;

ifStatement:
	DOUBLE_LBRACE IF expression DOUBLE_RBRACE template DOUBLE_LBRACE END DOUBLE_RBRACE;

loopStatement:
	DOUBLE_LBRACE LOOP IDENTIFIER IN IDENTIFIER DOUBLE_RBRACE template DOUBLE_LBRACE END
		DOUBLE_RBRACE;

expressionStatement: DOUBLE_LBRACE expression DOUBLE_RBRACE;

defineStatement:
	DOUBLE_LBRACE DEFINE IDENTIFIER (
		LPAREN (IDENTIFIER (COMMA IDENTIFIER)*)? RPAREN
	)? DOUBLE_RBRACE template DOUBLE_LBRACE END DOUBLE_RBRACE;

templateStatement:
	DOUBLE_LBRACE TEMPLATE IDENTIFIER expression* DOUBLE_RBRACE;

expression: pipelineExpr;

pipelineExpr: primaryExpr (PIPE functionCall)*;

primaryExpr:
	NUMBER
	| STRING
	| IDENTIFIER
	| varExpr
	| LPAREN expression RPAREN;

varExpr: DOLLAR IDENTIFIER;

functionCall:
	IDENTIFIER (primaryExpr)*
	| IDENTIFIER LPAREN primaryExpr? (COMMA primaryExpr)* RPAREN;
