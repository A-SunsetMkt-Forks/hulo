/*
 Copyright 2025 The Hulo Authors. All rights reserved.
 Use of this source code is governed by a
 MIT-style
 license that can be found in the LICENSE file.
 */
lexer grammar vbsLexer;

CALL: 'Call';
CLASS: 'Class';
END: 'End';
CONST: 'Const';
DIM: 'Dim';
REDIM: 'ReDim';
PRESERVE: 'Preserve';
ON: 'On';
ERROR: 'Error';
RESUME: 'Resume';
NEXT: 'Next';
GOTO: 'GoTo';
PUBLIC: 'Public';
PRIVATE: 'Private';
STOP: 'Stop';
DO: 'Do';
SELECT: 'Select';
CASE: 'Case';
WHILE: 'While';
UNTIL: 'Until';
EXIT: 'Exit';
WEND: 'Wend';
LOOP: 'Loop';
IS: 'Is';
ERASE: 'Erase';
EXECUTE: 'Execute';
EXECUTEGLOBAL: 'ExecuteGlobal';
FOR: 'For';
EACH: 'Each';
IN: 'In';
TO: 'To';
STEP: 'Step';
FUNCTION: 'Function';
WITH: 'WITH';
SUB_LIT: 'SUB';
PROPERTY: 'PROPERTY';
GET: 'Get';
LET: 'Let';
SET: 'Set';
BYVAL: 'ByVal';
BYREF: 'ByRef';
NOTHING: 'Nothing';
OPTION: 'Option';
EXPLICIT: 'Explicit';
IF: 'If';
THEN: 'Then';
ELSEIF: 'ELSEIF';
ELSE: 'Else';

ADD: '+';
SUB: '-';
MUL: '*';
DIV: '/';
IDIV: '\\';
MOD: 'Mod';
EXP: '^';
ASSIGN: '=';

LT: '<';
GT: '>';

COMMA: ',';
LPAREN: '(';
RPAREN: ')';
NUM: '-'? [0-9]+ ('.' [0-9]+)?;
IDENT: [a-zA-Z_][a-zA-Z0-9_]*;

WS: [ \t\n\r]+ -> skip;