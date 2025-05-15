/*
 Copyright 2025 The Hulo Authors. All rights reserved.
 Use of this source code is governed by a
 MIT-style
 license that can be found in the LICENSE file.
 */
parser grammar vbsParser;

options {
	tokenVocab = vbsLexer;
}

file: stmt*;

block: stmt*;

stmt:
	dimDecl
	| subDecl
	| dimDecl
	| funcDecl
	| propertyDecl
	| constStmt
	| withStmt
	| forStmt
	| callStmt
	| exitStmt
	| funcDecl
	| assignStmt
	| optionStmt
	| whileWendStmt
	| forEachStmt
	| selectStmt
	| doWhileStmt
	| ifStmt;

classDecl: CLASS IDENT classBody END CLASS;

classMemeber: (PUBLIC | PRIVATE) expr (COMMA expr)*;

classBody:
	classMemeber
	| dimDecl
	| funcDecl
	| propertyDecl
	| subDecl;

propertyDecl:
	export? PROPERTY (GET | SET | LET) IDENT block END PROPERTY;

dimDecl: DIM IDENT subscripts? (COMMA IDENT subscripts?)*;

reDimDecl: REDIM PRESERVE? expr;

subDecl: export? SUB IDENT block END SUB;

constStmt: export? CONST IDENT ASSIGN;

export: PUBLIC | PRIVATE;

subscripts: LPAREN RPAREN;

funcDecl: export? FUNCTION IDENT END FUNCTION;

onErrorStmt: ON ERROR (onErrorResume | onErrorGoto);

onErrorResume: RESUME NEXT;

onErrorGoto: GOTO NUM;

callStmt: CALL IDENT;

forStmt: FOR expr TO expr (STEP expr) block NEXT;

exitStmt: EXIT (DO | FOR | FUNCTION | PROPERTY | SUB_LIT)?;

assignStmt: export? (CONST | SET) IDENT ASSIGN expr;

withStmt: WITH expr block END WITH;

optionStmt: OPTION EXPLICIT;

ifStmt:
	IF expr THEN block (ELSEIF expr THEN block)* (ELSE block)? END IF;

forEachStmt: FOR EACH expr IN expr block NEXT expr?;

doWhileStmt: doWhilePreStmt | doWhilePostStmt;

doWhilePreStmt: DO (WHILE | UNTIL) expr block LOOP;

doWhilePostStmt: DO block LOOP (WHILE | UNTIL) expr;

whileWendStmt: WHILE expr block WEND;

selectStmt:
	SELECT CASE expr (CASE expr block)* (CASE ELSE block)? END SELECT;

expr: IDENT | NUM | NOTHING;