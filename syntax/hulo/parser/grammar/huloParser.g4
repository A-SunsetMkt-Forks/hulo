// $antlr-format alignTrailingComments true, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine true, allowShortBlocksOnASingleLine true, minEmptyLines 0, alignSemicolons ownLine
// $antlr-format alignColons trailing, singleLineOverrulesHangingColon true, alignLexerCommands true, alignLabels true, alignTrailers true

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
parser grammar huloParser;

options {
    tokenVocab = huloLexer;
    // superClass = huloParserBase;
}

file: (statement SEMI*)*;

block: LBRACE (statement SEMI*)* RBRACE;

comment: LineComment | BlockComment;

statement:
    comment
    | importDeclaration
    | moduleDeclaration
    | functionDeclaration
    | classDeclaration
    | enumDeclaration
    | traitDeclaration
    | implDeclaration
    | typeDeclaration
    | extendDeclaration
    | externDeclaration
    | macroStatement
    | assignStatement
    | lambdaAssignStatement
    | tryStatement
    | throwStatement
    | returnStatement
    | breakStatement
    | continueStatement
    | ifStatement
    | matchStatement
    | loopStatement
    | returnStatement
    | deferStatement
    | declareStatement
    | channelInputStatement
    | expressionStatement
;

// -----------------------
//
// expression

// assignModifier ident = expr
assignStatement:
    assignModifier = (LET | CONST | VAR) (Identifier | variableNames) (COLON type)? (
        ASSIGN (expression | matchStatement)
    )?
    | (variableExpression | variableNullableExpressions) (
        ASSIGN
        | ADD_ASSIGN
        | SUB_ASSIGN
        | MUL_ASSIGN
        | DIV_ASSIGN
        | MOD_ASSIGN
        | AND_ASSIGN
        | EXP_ASSIGN
    ) (expression | matchStatement)
;

lambdaAssignStatement: variableExpression COLON_ASSIGN (expression | matchStatement);

variableNullableExpressions: variableNullableExpression (COMMA variableNullableExpression)*;

variableNullableExpression: variableExpression | WILDCARD;

variableNames: LBRACE variableNameList RBRACE;

variableNameList: variableName (COMMA variableName)*;

// $i++ * ++$j - $k << 10 > 10 | true ^ 1 ? "hello" : null
conditionalExpression:
    conditionalBoolExpression (QUEST conditionalExpression COLON conditionalExpression)?
;

// $i++ * ++$j - $k << 10 > 10 | true ^ 1
conditionalBoolExpression:
    logicalExpression (conditionalOp = (BITAND | BITOR | BITXOR | AND | OR) logicalExpression)*
;

/* logicalExpression:
 * $i++ * ++$j - $k << 10 > 10
 * typeof ($i++ * ++$j - $k << 10 > 10) == "bool"
 */
logicalExpression:
    shiftExpression (logicalOp = (GT | LT | GE | LE | EQ | NEQ | MOD) shiftExpression)*
    | typeofExpression (EQ | NEQ) StringLiteral
;

// $i++ * ++$j - $k << 10
shiftExpression: addSubExpression (shiftOp = (SHL | SHR) addSubExpression)*;

// $i++ * ++$j - $k
addSubExpression: mulDivExpression (addSubOp = (ADD | SUB) mulDivExpression)*;

// $i++ * ++$j
mulDivExpression: incDecExpression (mulDivOp = (MUL | DIV) incDecExpression)*;

incDecExpression: preIncDecExpression | postIncDecExpression;

// ++$i
preIncDecExpression: (INC | DEC)? factor;

// $i++
postIncDecExpression: factor (INC | DEC)?;

/*
 * 3.14 -> literal
 * time.now() -> callExpression
 * (3.14 + (10 * $a - 3) / 5 % ++$i) >= 0 ? "hello" : true
 * [1, "abc"] >> env
 */
factor:
    LPAREN factor RPAREN
    | SUB factor
    | literal
    | listExpression
    | tripleExpression
    | mapExpression
    | methodExpression
    | variableExpression
    | callExpression
    | fileExpression
;

tripleExpression: LPAREN expression (COMMA expression)* RPAREN;

listExpression: LBRACK expression (COMMA expression)* RBRACK;

newDelExpression: (NEW | DELETE) callExpression;

mapExpression : LBRACE pair (COMMA pair)* RBRACE;
pair          : (NumberLiteral | BoolLiteral | StringLiteral | Identifier) COLON expression;

variableExpression: DOLLAR memberAccess;

methodExpression: DOLLAR memberAccess LPAREN receiverArgumentList? RPAREN;

expression:
    lambdaExpression
    | conditionalExpression
    | newDelExpression
    | classInitializeExpression
    | typeofExpression
    | channelOutputExpression
    | commandExpression
    | unsafeExpression
    | comptimeExpression
;

// expr (, expr)*
expressionList: expression (COMMA expression)*;

/* memberAccess:
 * std::time.now -> IDENT (DOUBLE_COLON IDENT (DOT IDENT LPAREN RPAREN))
 * time.now -> IDENT (DOT IDENT LPAREN expression RPAREN)
 * time.utc::now -> IDENT (DOT IDENT (DOUBLE_COLON IDENT LPAREN RPAREN (DOT IDENT LPAREN RPAREN (DOT IDENT LPAREN RPAREN))))
 * new builtin.std::chan<str> -> NEW memberAccess(IDENT DOT IDENT DOUBLE_COLON IDENT genericType LPAREN RPAREN)
 * pkg.scope::arr[1][8] -> IDENT (DOT IDENT (DOUBLE_COLON IDENT (LBRACK expression RBRACK (LBRACK expression RBRACK (DOT IDENT LPAREN RPAREN)))))
 * this.ident -> THIS (DOT IDENT)
 * super.classA.call -> SUPER (DOT IDENT (DOT IDENT LPAREN RPAREN))
 * super -> SUPER LPAREN RPAREN
 * 10.to_str -> NumberLiteral (DOT IDENT LPAREN RPAREN)
 * "abc".split -> StringLiteral (DOT IDENT LPAREN RPAREN)
 * str.split -> STR (DOT IDENT LPAREN funcArgumentList? RPAREN)
 * echo!
 * time.now!
 */
memberAccess:
    (
        Identifier genericArguments
        | Identifier memberAccessPoint?
        | (STR | NUM | BOOL) memberAccessPoint
        | literal memberAccessPoint
        | THIS memberAccessPoint?
        | SUPER memberAccessPoint?
    ) NOT?
;

fileExpression: FileStringLiteral callExpressionLinkedList?;

/* callExpression:
 * memberAccess(...).call()..call2()..call3()
 * println!()
 */
callExpression: memberAccess NOT? LPAREN receiverArgumentList? RPAREN callExpressionLinkedList?;

callExpressionLinkedList:
    DOT Identifier LPAREN receiverArgumentList? RPAREN callExpressionLinkedList?
    | DOUBLE_DOT Identifier LPAREN receiverArgumentList? RPAREN callExpressionLinkedList?
;

memberAccessPoint:
    DOT Identifier memberAccessPoint?
    | DOT Identifier genericArguments memberAccessPoint?
    | DOUBLE_COLON Identifier memberAccessPoint?
    | DOUBLE_COLON variableExpression memberAccessPoint?
    | LBRACK expression RBRACK memberAccessPoint?
;

literal: NumberLiteral | BoolLiteral | StringLiteral | NULL;

/* funcArgumentList:
 * 10, zero(), name: "hulo" -> expressionList COMMA namedArgumentList
 * name: "hulo", version: "^0.0.1" -> namedArgumentList
*/
receiverArgumentList: expressionList (COMMA namedArgumentList)? | namedArgumentList;

// ident: expr (, ident: expr)*
namedArgumentList: namedArgument (COMMA namedArgument)*;

// ident: expr
namedArgument: Identifier COLON expression;

// _ or ident: type
variableName: WILDCARD | Identifier (COLON type)?;

rangeExpression: NumberLiteral DOUBLE_DOT NumberLiteral;

expressionStatement: expression;

// -----------------------
//
// command

option: (shortOption | longOption);

// -a -abc
shortOption: SUB Identifier;

// --abc --a-b
longOption: DEC Identifier;

/*
 * std::grep -o "abc" "a.txt"
 * my_cmd.grep -o "abc" "a.txt"
 * my_cmd.std::grep -o "abc" "a.txt" "b.txt" -- -os platform::win32
 */
commandExpression:
    (memberAccess | CommandStringLiteral) (option conditionalExpression?)* (
        conditionalExpression
        | memberAccess
    )* (BITAND DEC (option conditionalExpression?)*)? (commandJoin | commandStream)?
;

commandJoin: (AND | OR) commandExpression;

commandStream: (BITOR | LT | GT | SHL | SHR) commandExpression;

commandMemberAccess: Identifier memberAccessPoint?;

commandAccessPoint:
    DOT Identifier commandAccessPoint?
    | DOUBLE_COLON Identifier commandAccessPoint?
;

// -----------------------
//
// function

/* receiverParameters:
 * (x: str, y: bool = true)
 * (x: str, y: bool = true, {kind: str = "obj", required z: num})
 * ({kind: str = "obj", required z: num})
 */
receiverParameters:
    LPAREN (
        receiverParameterList?
        | receiverParameterList (COMMA namedParameters)?
        | namedParameters
    ) RPAREN
;

/* receiverParameterList:
 * s: str, i: num = 10 -> receiverParameter (COMMA receiverParameter)
 */
receiverParameterList: receiverParameter (COMMA receiverParameter)*;

/* receiverParameter:
 * s: str -> IDENT (COLON type)
 * s: str = "hello world" -> IDENT (COLON type) (ASSIGN expression)
 * title?: str -> IDENT QUEST (COLON type)
 * title?: str = "default" -> IDENT QUEST (COLON type) (ASSIGN expression)
 */
receiverParameter: Identifier QUEST? (COLON type)? (ASSIGN expression)?;

namedParameters: LBRACE namedParameterList? RBRACE;

namedParameterList: namedParameter (COMMA namedParameter)*;

/* namedParameter:
 * name: str
 * required name: str = "hulo"
 * u: user = {name: "hulo", age: 10 + 5}
 */
namedParameter: REQUIRED? receiverParameter;

returnStatement: RETURN expressionList?;

functionDeclaration: standardFunctionDeclaration | lambdaFunctionDeclaration | functionSignature;

operatorIdentifier:
    OPERATOR (
        ADD
        | SUB
        | MUL
        | DIV
        | MOD
        | LT
        | GT
        | SHL
        | SHR
        | HASH
        | NEW
        | DELETE
        | LBRACK RBRACK
    )
;

standardFunctionDeclaration:
    functionModifier* FN (Identifier | operatorIdentifier) genericParameters? receiverParameters THROWS? (
        ARROW functionReturnValue
    )? block
;

functionReturnValue: type | LPAREN typeList RPAREN;

/* lambdaFunctionDeclaration:
 * fn f() => println("hello world")
 * fn f<T>(x: T) => println("type is ${typeof $x}")
 */
lambdaFunctionDeclaration:
    functionModifier* FN (Identifier | operatorIdentifier) genericParameters? lambdaExpression
;

/* lambdaExpression:
 * (x) => $x * 2
 * (x) => { $x, true }
 * (x) => { if $x > 10 { return 10 } return -1 }
 */
lambdaExpression: receiverParameters DOUBLE_ARROW lambdaBody;

lambdaBody: expression | LBRACE expressionList RBRACE | block;

/* functionSignature:
 * fn InputBox(prompt: str, title?: str, default?: str, xpos?: num, ypos?: num, helpfile?: str, context?: num) -> any;
 */
functionSignature:
    functionModifier* FN (Identifier | operatorIdentifier) genericParameters? receiverParameters THROWS? (
        ARROW functionReturnValue
    )?
;

functionModifier: PUB | COMPTIME;

macroStatement: AT memberAccess (LPAREN receiverArgumentList? RPAREN)?;

// -----------------------
//
// class and object

// ? class -------------
classDeclaration:
    macroStatement* classModifier* CLASS Identifier genericParameters? (COLON classSuper)? classBody
;

classModifier: PUB;

classSuper: memberAccess (COMMA memberAccess)*;

classBody: LBRACE ((comment | classMember | classMethod | classBuiltinMethod) SEMI*)* RBRACE;

classMember: macroStatement* classMemberModifier* Identifier COLON type (ASSIGN expression)?;

classMemberModifier: PUB | STATIC | FINAL | CONST;

classMethod:
    macroStatement* classMethodModifier* (standardFunctionDeclaration | lambdaFunctionDeclaration)
;

classMethodModifier: PUB | STATIC;

/*
 * const user({ required this.name, this.age = 10 })
 * const user(age: num, { required this.name }) throws {}
 * const user(name: str, age: num) => {}
 */
classBuiltinMethod:
    macroStatement* classBuiltinMethodModifier? Identifier classBuiltinParameters (
        THROWS? block
        | DOUBLE_ARROW block
    )?
;

classBuiltinMethodModifier: CONST;

/* classBuiltinParameters:
 * ()
 * (kind: str)
 * ({ required kind: str, age: num })
 * (kind: str, { required age: num })
 * (kind: str, { required age: num, this.name = "hello" })
 */
classBuiltinParameters:
    LPAREN (
        receiverParameterList?
        | receiverParameterList (COMMA classNamedParameters)?
        | classNamedParameters
    ) RPAREN
;

// { required kind: str, this.name = "hello", required this.age }
classNamedParameters: LBRACE classNamedParameterList? RBRACE;

classNamedParameterList: classNamedParameter (COMMA classNamedParameter)*;

/* classNamedParameter:
 * required this.name
 * required super.name
 * required super.classA.name = 10
 * kind: str = "hello"
 */
classNamedParameter:
    REQUIRED? receiverParameter
    | REQUIRED? (THIS | SUPER) classNamedParameterAccessPoint (ASSIGN expression)?
;

classNamedParameterAccessPoint: DOT Identifier classNamedParameterAccessPoint?;

// user{name: "hello", age: 0}
classInitializeExpression: Identifier LBRACE namedArgumentList RBRACE;

// ? enum -------------
enumDeclaration:
    macroStatement* enumModifier? ENUM Identifier genericParameters? (
        enumBodySimple
        | enumBodyAssociated
        | enumBodyADT
    )
;

enumModifier: PUB | FINAL;

// Basic enum: enum Status { Pending, Approved, Rejected }
enumBodySimple: LBRACE enumValue (COMMA comment* enumValue)* RBRACE;

// Associated value enum: enum Protocol { port: num; tcp(6), udp(17) }
// Complex associated enum: enum Protocol { final port: num; const Protocol(...); tcp(6), udp(17); fn get_port() -> num { ... } }
enumBodyAssociated:
    LBRACE enumAssociatedFields? enumAssociatedConstructor? enumAssociatedValues SEMI? enumAssociatedMethods? RBRACE
;

// ADT enum: enum NetworkPacket { TCP { src_port: num, dst_port: num }, UDP { port: num } }
enumBodyADT: LBRACE enumVariant (COMMA comment* enumVariant)* enumMethods? RBRACE;

// Enum value with optional assignment: Pending or Pending = 0
enumValue: Identifier (ASSIGN expression)? comment?;

// Associated fields declaration: port: num or final port: num
enumAssociatedFields: enumAssociatedField (COMMA enumAssociatedField)* SEMI;

enumAssociatedField: enumFieldModifier* Identifier COLON type;

enumFieldModifier: FINAL | CONST;

enumField: Identifier COLON type;

// Associated values: tcp(6), udp(17)
enumAssociatedValues: enumAssociatedValue (COMMA comment* enumAssociatedValue)*;

enumAssociatedValue: Identifier LPAREN expressionList? RPAREN;

// Associated enum constructor: const Protocol($this.port = -1) or const Protocol.One(): $this.port = 1 {}
enumAssociatedConstructor: enumConstructor+;

// Enum constructor: const Protocol(...) or const Protocol.One(...)
// Examples:
// const Protocol($this.port = -1);
// const Protocol.One(): $this.port = 1 {}
// const Protocol.Port(v: num): $this.port = v {}
// const Protocol(zero: bool, ...v: num) { ... }
enumConstructor:
    macroStatement* enumBuiltinMethodModifier? enumConstructorName enumConstructorParameters (
        enumConstructorInit? THROWS? block
        | enumConstructorInit? DOUBLE_ARROW block
        | enumConstructorInit? SEMI
    )?
;

// Protocol or Protocol.One
enumConstructorName: Identifier (DOT Identifier)?;

// Enum constructor parameters: ($this.port = -1) or (v: num) or (zero: bool, ...v: num)
enumConstructorParameters:
    LPAREN (
        enumConstructorDirectInit
        | receiverParameterList (COMMA enumConstructorDirectInit)*
        | receiverParameterList
    )? RPAREN
;

// $this.port = -1 (直接赋值，不需要冒号)
enumConstructorDirectInit: variableExpression ASSIGN expression;

// : $this.port = 1 (冒号后直接赋值)
enumConstructorInit: COLON enumConstructorDirectInit;

// Associated enum methods
enumAssociatedMethods: enumMethod+;

// ADT variant: TCP { src_port: num, dst_port: num }
enumVariant: Identifier LBRACE enumField (COMMA enumField)* RBRACE;

// Enum methods (for ADT and associated enums)
enumMethods: enumMethod+;

enumMethod:
    macroStatement* enumMethodModifier* (standardFunctionDeclaration | lambdaFunctionDeclaration)
;

enumMethodModifier: PUB | STATIC;

// Legacy support
enumMember: macroStatement* enumMemberModifier* Identifier COLON type;

enumMemberModifier: PUB | FINAL;

enumBuiltinMethodModifier: CONST;

enumInitialize: (enumInitializeMember COMMA)* enumInitializeMember SEMI;

enumInitializeMember: Identifier LPAREN expressionList? RPAREN;

// ? trait -------------

traitDeclaration: macroStatement* traitModifier? TRAIT Identifier genericParameters? traitBody;

traitModifier: PUB;

traitBody: LBRACE ((comment | traitMember) SEMI*)* RBRACE;

traitMember: traitMemberModifier* (Identifier | operatorIdentifier) COLON type;

traitMemberModifier: PUB | FINAL | CONST | STATIC;

// ? impl -------------

implDeclaration: IMPL memberAccess FOR (implDeclarationBody | implDeclarationBinding);

implDeclarationBinding: memberAccess (COMMA memberAccess)*;

implDeclarationBody: memberAccess classBody;

// ? extend -------------

extendDeclaration: EXTEND (extendEnum | extendClass | extendTrait | extendType | extendMod);

extendEnum: ENUM Identifier (enumBodySimple | enumBodyAssociated | enumBodyADT);

extendClass: CLASS Identifier classBody;

extendTrait: TRAIT Identifier traitBody;

extendType: TYPE Identifier LBRACE type? RBRACE;

extendMod: MOD_LIT Identifier LBRACE (moduleStatement)* RBRACE;

// -----------------------
//
// module

importDeclaration: IMPORT (importAll | importSingle | importMulti);

/* importSingle:
 * "time"
 * "time" as t
 */
importSingle: StringLiteral asIdentifier?;

/* importSingle:
 * * from "time"
 * * as t from "time"
 */
importAll: MUL asIdentifier? FROM StringLiteral;

// { now as utc_now, time } from "time"
importMulti:
    LBRACE (identifierAsIdentifier (COMMA identifierAsIdentifier)*)? RBRACE FROM StringLiteral
;

asIdentifier: AS Identifier;

/* identifierAsIdentifier:
 * math -> Identifier
 * time as t -> Identifier asIdentifier
 */
identifierAsIdentifier: Identifier asIdentifier?;

moduleDeclaration: PUB? MOD_LIT Identifier LBRACE (moduleStatement)* RBRACE;

moduleStatement:
    useDeclaration
    | moduleDeclaration
    | classDeclaration
    | enumDeclaration
    | traitDeclaration
    | functionDeclaration
    | extendDeclaration
    | assignStatement
;

/* useDeclaration:
 * use color as c
 * use time.std::now
 * use * as t from "time"
 * use {now as n, date} from "time"
 */
useDeclaration: USE (useSingle | useMulti | useAll);

useSingle: Identifier asIdentifier?;

useMulti: LBRACE identifierAsIdentifier (COMMA identifierAsIdentifier)* RBRACE FROM StringLiteral;

useAll: MUL asIdentifier? FROM StringLiteral;

// -----------------------
//
// type system

typeDeclaration: TYPE Identifier genericParameters? ASSIGN type;

// <type (, type)*>
genericArguments: LT typeList? GT;

// <T: type, ...>
genericParameters: LT genericParameterList GT;

// T: type, ...
genericParameterList: genericParameter (COMMA genericParameter)*;

// T: type
genericParameter: Identifier (COLON type)?;

compositeType: (BITOR type | BITAND type) compositeType?;

/* type:
 * str | num? | bool
 * 'udp' | 'tcp'
 * str & num
 * error::runtime | chan<str> | std::time.date
 * user[5][3]
 * ...user
 * { name: str, age: num }
 * [str, num, bool]
 */
type: (
        (STR | NUM | BOOL | ANY | Identifier) typeAccessPoint?
        | memberAccess
        | StringLiteral
        | functionType
        | objectType
        | tupleType
    ) QUEST? compositeType?
;

//// TODO

typeLiteral: (STR | NUM | BOOL | ANY | Identifier);

// str?
nullableType: type QUEST;

// str | num
unionType: type (BITOR type)*;

// str & num
intersectionType: type (BITAND type)*;

// str[5]
arrayType: type LBRACK NumberLiteral RBRACK;

///// ENDING - TODO

typeAccessPoint: LBRACK NumberLiteral RBRACK typeAccessPoint?;

typeList: type (COMMA type)*;

typeofExpression: TYPEOF expression;

asExpression: variableExpression AS type;

// { name: str, age: num }
objectType: LBRACE objectTypeMember (COMMA objectTypeMember)* RBRACE;

// name: str
objectTypeMember: Identifier COLON type;

// [str, num, bool]
tupleType: LBRACK typeList RBRACK;

/* functionType:
 * (x: str)
 * (str) -> num
 * (x: str, y: num) throws -> (num, user)
 */
functionType: receiverParameters THROWS? (ARROW functionReturnValue)?;

// -----------------------
//
// exception

tryStatement: TRY block catchClause* finallyClause?;

/* catchClause:
 * catch (e: error::runtime) {}
 * catch error::runtime {}
 * catch {}
 */
catchClause: CATCH (catchClauseReceiver | memberAccess)? block;

/* catchClauseReceiver:
 * (e)
 * (e: error)
 * (e: error::runtime)
 */
catchClauseReceiver: LPAREN Identifier (COLON type)? RPAREN;

finallyClause: FINALLY block;

throwStatement: THROW (newDelExpression | callExpression);

// -----------------------
//
// control flow

breakStatement: BREAK Identifier?;

continueStatement: CONTINUE;

// if statement -----------------------

ifStatement: IF conditionalExpression block (ELSE (ifStatement | block))?;

// match statement -----------------------

matchStatement:
    MATCH expression LBRACE (matchCaseClause COMMA comment*)* (matchDefaultClause COMMA? comment*)? RBRACE
;

matchCaseClause: (type | matchEnum | memberAccess | matchTriple | expression | rangeExpression) DOUBLE_ARROW matchCaseBody
;

matchEnum: memberAccess LPAREN matchEnumMember (COMMA matchEnumMember)* RPAREN (IF expression)?;

matchEnumMember: Identifier | literal | WILDCARD;

matchTriple: LPAREN matchTripleValue (COMMA matchTripleValue)* RPAREN (IF expression)?;

matchTripleValue: Identifier | WILDCARD;

matchCaseBody: expression | returnStatement | block;

matchDefaultClause: WILDCARD DOUBLE_ARROW matchCaseBody;

// loop statement -----------------------

loopStatement:
    loopLabel? (
        whileStatement
        | doWhileStatement
        | rangeStatement
        | forStatement
        | foreachStatement
    )
;

loopLabel: Identifier COLON;

// loop ($item, $index) in $arr {
//     echo $item
// }
// loop $item in $arr {
//     echo $item
// }
foreachStatement: LOOP foreachClause (IN | OF) expression block;

foreachClause: foreachVariableName | LPAREN foreachVariableName COMMA foreachVariableName RPAREN;

foreachVariableName: variableName | variableExpression;

forStatement: LOOP LPAREN? forClause RPAREN? block;

forClause: statement? SEMI expression SEMI expression?;

rangeStatement: LOOP Identifier IN rangeClause block;

rangeClause: RANGE LPAREN expression COMMA expression (COMMA expression)? RPAREN;

doWhileStatement: DO block LOOP LPAREN expression RPAREN;

whileStatement: LOOP expression? block;

// -----------------------
//
// spec

deferStatement: DEFER callExpression;

/* declareStatement:
 declare {
    class user {
        name: str
        age: num
    }
 }

 declare class caller {
    file: str
    func: str
    line: num
 }
 */
declareStatement:
    DECLARE (
        block
        | classDeclaration
        | enumDeclaration
        | traitDeclaration
        | functionDeclaration
        | moduleDeclaration
        | typeDeclaration
    )
;

channelInputStatement:
    variableExpression BACKARROW channelPayload = (StringLiteral | NumberLiteral) BITAND?
;
channelOutputExpression: BACKARROW variableExpression;

unsafeExpression: UnsafeLiteral;

comptimeExpression: COMPTIME block;

// -----------------------
//
// extern declaration

externDeclaration: EXTERN externList;

externList: externItem (COMMA externItem)*;

externItem: Identifier (COLON type)?;
