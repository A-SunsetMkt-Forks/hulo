parser grammar UnsafeParser;

options {
	tokenVocab = unsafeLexer;
}

template: (content | statement)*;

content: TEXT;

statement: ifStatement | loopStatement | expressionStatement;

ifStatement:
	DOUBLE_LBRACE IF expression DOUBLE_RBRACE template DOUBLE_LBRACE END DOUBLE_RBRACE;

loopStatement:
	DOUBLE_LBRACE LOOP IDENTIFIER IN IDENTIFIER DOUBLE_RBRACE template DOUBLE_LBRACE END
		DOUBLE_RBRACE;

expressionStatement: DOUBLE_LBRACE expression DOUBLE_RBRACE;

expression: pipelineExpr;

pipelineExpr: primaryExpr (PIPE functionCall)*;

primaryExpr:
	NUMBER
	| STRING
	| IDENTIFIER
	| LPAREN expression RPAREN;

functionCall:
	IDENTIFIER (primaryExpr)*
	| IDENTIFIER LPAREN primaryExpr? (COMMA primaryExpr)* RPAREN;