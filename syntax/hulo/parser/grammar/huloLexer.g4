// $antlr-format alignTrailingComments true, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine true, allowShortBlocksOnASingleLine true, minEmptyLines 0, alignSemicolons ownLine
// $antlr-format alignColons trailing, singleLineOverrulesHangingColon true, alignLexerCommands true, alignLabels true, alignTrailers true

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
lexer grammar huloLexer;

// -----------------------
//
// module

MOD_LIT : 'mod';
USE     : 'use';
IMPORT  : 'import';
FROM    : 'from';

// -----------------------
//
// type system

TYPE   : 'type';
TYPEOF : 'typeof';
AS     : 'as';

// -----------------------
//
// control flow

IF       : 'if';
ELSE     : 'else';
MATCH    : 'match';
DO       : 'do';
LOOP     : 'loop';
IN       : 'in';
OF       : 'of';
RANGE    : 'range';
CONTINUE : 'continue';
BREAK    : 'break';

// -----------------------
//
// modifier

LET      : 'let';
VAR      : 'var';
CONST    : 'const';
STATIC   : 'static';
FINAL    : 'final';
PUB      : 'pub';
REQUIRED : 'required';

// -----------------------
//
// exception handle

TRY     : 'try';
CATCH   : 'catch';
FINALLY : 'finally';
THROW   : 'throw';
THROWS  : 'throws';

// -----------------------
//
// function

FN       : 'fn';
OPERATOR : 'operator';
RETURN   : 'return';

// -----------------------
//
// class and object

ENUM   : 'enum';
CLASS  : 'class';
TRAIT  : 'trait';
IMPL   : 'impl';
FOR    : 'for';
THIS   : 'this';
SUPER  : 'super';
EXTEND : 'extend';

// -----------------------
//
// spec

DECLARE  : 'declare';
DEFER    : 'defer';
COMPTIME : 'comptime';
WHEN     : 'when';
UNSAFE   : 'unsafe';
EXTERN   : 'extern';

DOT        : '.';
DOUBLE_DOT : '..';

COMMA  : ',';
LPAREN : '(';
RPAREN : ')';
LBRACK : '[';
RBRACK : ']';
LBRACE : '{';
RBRACE : '}';
EXP    : '**';
MUL    : '*';
MOD    : '%';
DIV    : '/';
ADD    : '+';
SUB    : '-';

HASH      : '#';
AT        : '@';
QUEST     : '?';
DOLLAR    : '$';
BACKSLASH : '\\';
WILDCARD  : '_';

ASSIGN     : '=';
ADD_ASSIGN : '+=';
SUB_ASSIGN : '-=';
MUL_ASSIGN : '*=';
EXP_ASSIGN : '**=';
DIV_ASSIGN : '/=';
MOD_ASSIGN : '%=';
AND_ASSIGN : '&&=';

LT  : '<';
GT  : '>';
LE  : '<=';
GE  : '>=';
NEQ : '!=';
EQ  : '==';

SHL          : '<<';
SHR          : '>>';
ARROW        : '->';
BACKARROW    : '<-';
DOUBLE_ARROW : '=>';

ELLIPSIS     : '...';
COLON        : ':';
DOUBLE_COLON : '::';
COLON_ASSIGN : ':=';
SEMI         : ';';

AND : '&&';
OR  : '||';
NOT : '!';

INC : '++';
DEC : '--';

BITAND : '&';
BITOR  : '|';
BITXOR : '^';

SINGLE_QUOTE : '\'';
QUOTE        : '"';
TRIPLE_QUOTE : '"""';

NEW    : 'new';
DELETE : 'delete';

NULL : 'null';
NUM  : 'num';
STR  : 'str';
BOOL : 'bool';
ANY  : 'any';

UnsafeLiteral: DOUBLE_LBRACK (BACKSLASH RBRACK | ~[\]])* DOUBLE_RBRACK;

UnsafeBlock: LBRACE (BACKSLASH RBRACE | ~[}])* RBRACE;

DOUBLE_LBRACK : '[[';
DOUBLE_RBRACK : ']]';

fragment TRUE  : 'true';
fragment FALSE : 'false';

NumberLiteral: '-'? [0-9]+ ('.' [0-9]+)?;
StringLiteral:
    SINGLE_QUOTE (ESC | ~['\\\n])* SINGLE_QUOTE // ' '
    | QUOTE (ESC | ~["\\\n])* QUOTE             // " "
    | TRIPLE_QUOTE .*? TRIPLE_QUOTE
    | RawStringLiteral
; // """  """
BoolLiteral: TRUE | FALSE;

RawStringLiteral     : 'r"' (ESC | ~["\\\n])* QUOTE;
FileStringLiteral    : 'f"' (ESC | ~["\\\n])* QUOTE;
CommandStringLiteral : 'c"' (ESC | ~["\\\n])* QUOTE;

ESC: '\\' ['"\\bfnrt];

LineComment  : '//' ~[\n]*;
BlockComment : '/*' .*? '*/';

Identifier : [a-zA-Z_][a-zA-Z0-9_]* ('-' [a-zA-Z0-9_]+)*;
WS         : [ \t\n\r]+ -> skip;
