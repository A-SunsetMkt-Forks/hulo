lexer grammar unsafeLexer;

// ================== default mode ==================
DOUBLE_LBRACE: '{{' -> pushMode(TEMPLATE_MODE);
TEXT: ~[{]+;

// ================== template mode ==================
mode TEMPLATE_MODE;

DOUBLE_RBRACE: '}}' -> popMode;

IF: 'if';
LOOP: 'loop';
END: 'end';
IN: 'in';

LPAREN: '(';
RPAREN: ')';

DOLLAR: '$';
PIPE: '|';
COMMA: ',';

IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]*;

STRING: '"' (~["\\] | '\\' .)* '"';

NUMBER: [0-9]+ ('.' [0-9]+)? ([eE] [+\-]? [0-9]+)?;

WS: [ \t\r\n]+ -> skip;