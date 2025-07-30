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
    ifStatement
    | loopStatement
    | expressionStatement
    | macroStatement
    | templateStatement
    | variableStatement
;

variableStatement: STATEMENT_LBRACE IDENTIFIER COLON_EQUAL expression STATEMENT_RBRACE;

ifStatement:
    STATEMENT_LBRACE IF expression STATEMENT_RBRACE template elseStatement? STATEMENT_LBRACE (
        END
        | ENDIF
    ) STATEMENT_RBRACE
;

elseStatement: STATEMENT_LBRACE ELSE STATEMENT_RBRACE template;

loopStatement:
    STATEMENT_LBRACE LOOP IDENTIFIER IN IDENTIFIER STATEMENT_RBRACE template STATEMENT_LBRACE (
        END
        | ENDLOOP
    ) STATEMENT_RBRACE
;

expressionStatement: DOUBLE_LBRACE expression DOUBLE_RBRACE;

macroStatement:
    STATEMENT_LBRACE MACRO IDENTIFIER (LPAREN (IDENTIFIER (COMMA IDENTIFIER)*)? RPAREN)? STATEMENT_RBRACE template STATEMENT_LBRACE (
        END
        | ENDMACRO
    ) STATEMENT_RBRACE
;

templateStatement: STATEMENT_LBRACE TEMPLATE IDENTIFIER expression* STATEMENT_RBRACE;

expression: pipelineExpr;

pipelineExpr: functionCall (PIPE functionCall)*;

primaryExpr: NUMBER | STRING | IDENTIFIER | varExpr | LPAREN expression RPAREN;

varExpr: DOLLAR IDENTIFIER;

functionCall:
    IDENTIFIER (primaryExpr)*
    | IDENTIFIER LPAREN primaryExpr? (COMMA primaryExpr)* RPAREN
;
