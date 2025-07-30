// $antlr-format alignTrailingComments true, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine true, allowShortBlocksOnASingleLine true, minEmptyLines 0, alignSemicolons ownLine
// $antlr-format alignColons trailing, singleLineOverrulesHangingColon true, alignLexerCommands true, alignLabels true, alignTrailers true

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
lexer grammar unsafeLexer;

COMMENT_LBRACE   : '{#' -> pushMode(COMMENT_MODE);
DOUBLE_LBRACE    : '{{' -> pushMode(EXPRESSION_MODE);
STATEMENT_LBRACE : '{%' -> pushMode(STATEMENT_MODE);

TEXT: ~[{]+;

mode COMMENT_MODE;

COMMENT_RBRACE  : '#}' -> popMode;
COMMENT_CONTENT : .+?  -> skip;

mode EXPRESSION_MODE;

DOUBLE_RBRACE: '}}' -> popMode;

LPAREN : '(';
RPAREN : ')';

LBRACE : '{';
RBRACE : '}';

HASH         : '#';
MOD          : '%';
DOLLAR       : '$';
PIPE         : '|';
COMMA        : ',';
COLON_EQUAL  : ':=';
DOT          : '.';
SINGLE_QUOTE : '\'';
QUOTE        : '"';

STRING : SINGLE_QUOTE (ESC | ~['\\\n])* SINGLE_QUOTE | QUOTE (ESC | ~["\\\n])* QUOTE;
NUMBER : [0-9]+ ('.' [0-9]+)? ([eE] [+\-]? [0-9]+)?;

IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]*;

ESC: '\\' ['"\\bfnrt];

WS: [ \t\r\n]+ -> skip;

mode STATEMENT_MODE;

STATEMENT_RBRACE: '%}' -> popMode;

IF       : 'if';
ELSE     : 'else';
LOOP     : 'loop';
IN       : 'in';
MACRO    : 'macro';
TEMPLATE : 'template';

END      : 'end';
ENDIF    : 'endif';
ENDLOOP  : 'endloop';
ENDMACRO : 'endmacro';

IDENTIFIER_STMT: [a-zA-Z_][a-zA-Z0-9_]*;

WS_STMT: [ \t\r\n]+ -> skip;
