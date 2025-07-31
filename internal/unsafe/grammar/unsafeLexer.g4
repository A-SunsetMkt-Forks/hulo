// $antlr-format alignTrailingComments true, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine true, allowShortBlocksOnASingleLine true, minEmptyLines 0, alignSemicolons ownLine
// $antlr-format alignColons trailing, singleLineOverrulesHangingColon true, alignLexerCommands true, alignLabels true, alignTrailers true

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
lexer grammar unsafeLexer;

COMMENT_START : '{#' -> pushMode(COMMENT_MODE);
EXPR_START    : '{{' -> pushMode(EXPRESSION_MODE);
STMT_START    : '{%' -> pushMode(STATEMENT_MODE);

TEXT: ~[{]+;

fragment TRUE  : 'true';
fragment FALSE : 'false';

fragment Identifier: [a-zA-Z_][a-zA-Z0-9_]*;

fragment SINGLE_QUOTE : '\'';
fragment QUOTE        : '"';
fragment ESC          : '\\' ['"\\bfnrt];

fragment String : SINGLE_QUOTE (ESC | ~['\\\n])* SINGLE_QUOTE | QUOTE (ESC | ~["\\\n])* QUOTE;
fragment Number : '-'? [0-9]+ ('.' [0-9]+)? ([eE] [+\-]? [0-9]+)?;

mode COMMENT_MODE;

COMMENT_END     : '#}' -> popMode;
COMMENT_CONTENT : ~[#}]+;

mode EXPRESSION_MODE;

EXPR_END: '}}' -> popMode;

LPAREN : '(';
RPAREN : ')';

LBRACE : '{';
RBRACE : '}';

HASH        : '#';
MOD         : '%';
DOLLAR      : '$';
PIPE        : '|';
COMMA       : ',';
COLON_EQUAL : ':=';
DOT         : '.';

// Comparison operators
EQ : '==';
NE : '!=';
GT : '>';
LT : '<';
GE : '>=';
LE : '<=';

// Logical operators
AND : '&&';
OR  : '||';
NOT : '!';

STRING  : String;
NUMBER  : Number;
BOOLEAN : TRUE | FALSE;

IDENTIFIER: Identifier;

WS: [ \t\r\n]+ -> skip;

mode STATEMENT_MODE;

STMT_END: '%}' -> popMode;

LPAREN_STMT : '(';
RPAREN_STMT : ')';

IF    : 'if';
ELSE  : 'else';
LOOP  : 'loop';
IN    : 'in';
MACRO : 'macro';

END      : 'end';
ENDIF    : 'endif';
ENDLOOP  : 'endloop';
ENDMACRO : 'endmacro';

// Comparison operators in STATEMENT_MODE
EQ_STMT : '==';
NE_STMT : '!=';
GT_STMT : '>';
LT_STMT : '<';
GE_STMT : '>=';
LE_STMT : '<=';

// Logical operators in STATEMENT_MODE
AND_STMT : '&&';
OR_STMT  : '||';
NOT_STMT : '!';

STRING_STMT  : String;
NUMBER_STMT  : Number;
BOOLEAN_STMT : TRUE | FALSE;

IDENTIFIER_STMT: Identifier;

WS_STMT: [ \t\r\n]+ -> skip;
