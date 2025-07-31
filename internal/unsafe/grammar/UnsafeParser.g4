// $antlr-format alignTrailingComments true, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine true, allowShortBlocksOnASingleLine true, minEmptyLines 0, alignSemicolons ownLine
// $antlr-format alignColons trailing, singleLineOverrulesHangingColon true, alignLexerCommands true, alignLabels true, alignTrailers true

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
parser grammar unsafeParser;

options {
    tokenVocab = unsafeLexer;
}

template: (content | statement)*;

content: TEXT;

statement:
    commentStatement
    | ifStatement
    | loopStatement
    | expressionStatement
    | macroStatement
    | variableStatement
;

commentStatement: COMMENT_START COMMENT_CONTENT COMMENT_END;

variableStatement: STMT_START IDENTIFIER COLON_EQUAL expression STMT_END;

ifStatement:
    STMT_START IF expression STMT_END template elseStatement? STMT_START (END | ENDIF) STMT_END
;

elseStatement: STMT_START ELSE STMT_END template;

loopStatement:
    STMT_START LOOP IDENTIFIER_STMT IN IDENTIFIER_STMT STMT_END template STMT_START (END | ENDLOOP) STMT_END
;

expressionStatement: EXPR_START pipelineExpression EXPR_END;

macroStatement:
    STMT_START MACRO IDENTIFIER_STMT (
        LPAREN_STMT (IDENTIFIER_STMT (COMMA IDENTIFIER_STMT)*)? RPAREN_STMT
    )? STMT_END template STMT_START (END | ENDMACRO) STMT_END
;

pipelineExpression: expression (PIPE expression)*;

expression:
    NUMBER
    | STRING
    | BOOLEAN
    | IDENTIFIER
    | NUMBER_STMT
    | STRING_STMT
    | BOOLEAN_STMT
    | varExpr
    | LPAREN expression RPAREN
    | functionCall
    | logicalOrExpression
    | logicalOrExpressionStmt
;

logicalOrExpression:
    logicalAndExpression (OR logicalAndExpression)*
;

logicalAndExpression:
    equalityExpression (AND equalityExpression)*
;

equalityExpression:
    comparisonExpression ((EQ | NE) comparisonExpression)*
;

comparisonExpression:
    primaryExpression ((GT | LT | GE | LE) primaryExpression)*
;

primaryExpression:
    NUMBER
    | STRING
    | BOOLEAN
    | IDENTIFIER
    | NUMBER_STMT
    | STRING_STMT
    | BOOLEAN_STMT
    | varExpr
    | LPAREN expression RPAREN
    | functionCall
    | NOT primaryExpression
    | NOT_STMT primaryExpressionStmt
;

varExpr: DOLLAR IDENTIFIER;

functionCall: IDENTIFIER (expression)* | IDENTIFIER LPAREN expression? (COMMA expression)* RPAREN;

// Statement mode expressions
logicalOrExpressionStmt:
    logicalAndExpressionStmt (OR_STMT logicalAndExpressionStmt)*
;

logicalAndExpressionStmt:
    equalityExpressionStmt (AND_STMT equalityExpressionStmt)*
;

equalityExpressionStmt:
    comparisonExpressionStmt ((EQ_STMT | NE_STMT) comparisonExpressionStmt)*
;

comparisonExpressionStmt:
    primaryExpressionStmt ((GT_STMT | LT_STMT | GE_STMT | LE_STMT) primaryExpressionStmt)*
;

primaryExpressionStmt:
    NUMBER_STMT
    | STRING_STMT
    | BOOLEAN_STMT
    | IDENTIFIER_STMT
    | LPAREN_STMT expression RPAREN_STMT
    | NOT_STMT primaryExpressionStmt
;
