// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"strconv"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/token"
)

// ===== 基础表达式创建函数 =====

// Word 创建字符串字面量
func Word(v string) *ast.Word {
	return &ast.Word{Val: v, ValPos: token.Pos(1)}
}

// Ident 创建标识符
func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name, NamePos: token.Pos(1)}
}

// Binary 创建二元表达式
func Binary(x ast.Expr, op token.Token, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:     x,
		OpPos: token.Pos(1),
		Op:    op,
		Y:     y,
	}
}

// Unary 创建一元表达式
func Unary(op token.Token, x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		OpPos: token.Pos(1),
		Op:    op,
		X:     x,
	}
}

// ===== 语句创建函数 =====

// ExprStmt 创建表达式语句
func ExprStmt(x ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{X: x}
}

// CmdExpr 创建命令表达式
func CmdExpr(name string, args ...ast.Expr) *ast.CmdExpr {
	return &ast.CmdExpr{
		Name: Ident(name),
		Recv: args,
	}
}

// CmdStmt 创建命令语句
func CmdStmt(name string, args ...ast.Expr) *ast.ExprStmt {
	return ExprStmt(CmdExpr(name, args...))
}

// AssignStmt 创建赋值语句
func AssignStmt(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs:    lhs,
		Assign: token.Pos(1),
		Rhs:    rhs,
	}
}

// LocalAssignStmt 创建 local 赋值语句
func LocalAssignStmt(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Local:  token.Pos(1),
		Lhs:    lhs,
		Assign: token.Pos(1),
		Rhs:    rhs,
	}
}

// ReturnStmt 创建 return 语句
func ReturnStmt(x ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{
		Return: token.Pos(1),
		X:      x,
	}
}

// ===== 控制流语句创建函数 =====

// IfStmt 创建 if 语句
func IfStmt(cond ast.Expr, body, elseBody ast.Stmt) *ast.IfStmt {
	return &ast.IfStmt{
		If:   token.Pos(1),
		Cond: cond,
		Semi: token.Pos(1),
		Then: token.Pos(1),
		Body: body.(*ast.BlockStmt),
		Else: elseBody,
		Fi:   token.Pos(1),
	}
}

// WhileStmt 创建 while 循环
func WhileStmt(cond ast.Expr, body []ast.Stmt) *ast.WhileStmt {
	return &ast.WhileStmt{
		While: token.Pos(1),
		Cond:  cond,
		Semi:  token.Pos(1),
		Do:    token.Pos(1),
		Body: &ast.BlockStmt{
			List: body,
		},
		Done: token.Pos(1),
	}
}

// UntilStmt 创建 until 循环
func UntilStmt(cond ast.Expr, body []ast.Stmt) *ast.UntilStmt {
	return &ast.UntilStmt{
		Until: token.Pos(1),
		Cond:  cond,
		Semi:  token.Pos(1),
		Do:    token.Pos(1),
		Body: &ast.BlockStmt{
			List: body,
		},
		Done: token.Pos(1),
	}
}

// ForInStmt 创建 for-in 循环
func ForInStmt(varName, list ast.Expr, body []ast.Stmt) *ast.ForInStmt {
	return &ast.ForInStmt{
		For:  token.Pos(1),
		Var:  varName,
		In:   token.Pos(1),
		List: list,
		Semi: token.Pos(1),
		Do:   token.Pos(1),
		Body: &ast.BlockStmt{
			List: body,
		},
		Done: token.Pos(1),
	}
}

// ForStmt 创建 C 风格的 for 循环
func ForStmt(init ast.Node, cond ast.Expr, post ast.Node, body []ast.Stmt) *ast.ForStmt {
	return &ast.ForStmt{
		For:    token.Pos(1),
		Lparen: token.Pos(1),
		Init:   init,
		Semi1:  token.Pos(1),
		Cond:   cond,
		Semi2:  token.Pos(1),
		Post:   post,
		Rparen: token.Pos(1),
		Do:     token.Pos(1),
		Body: &ast.BlockStmt{
			List: body,
		},
		Done: token.Pos(1),
	}
}

// CaseStmt 创建 case 语句
func CaseStmt(x ast.Expr, patterns []*ast.CaseClause, elseBody []ast.Stmt) *ast.CaseStmt {
	var elseBlock *ast.BlockStmt
	if len(elseBody) > 0 {
		elseBlock = &ast.BlockStmt{
			List: elseBody,
		}
	}

	return &ast.CaseStmt{
		Case:     token.Pos(1),
		X:        x,
		In:       token.Pos(1),
		Patterns: patterns,
		Else:     elseBlock,
		Esac:     token.Pos(1),
	}
}

// CaseClause 创建 case 子句
func CaseClause(conds []ast.Expr, body []ast.Stmt) *ast.CaseClause {
	return &ast.CaseClause{
		Conds:  conds,
		Rparen: token.Pos(1),
		Body: &ast.BlockStmt{
			List: body,
		},
		Semi: token.Pos(1),
	}
}

// FuncDecl 创建函数声明
func FuncDecl(name string, body []ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Function: token.Pos(1),
		Name:     Ident(name),
		Lparen:   token.Pos(1),
		Rparen:   token.Pos(1),
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

// BlockStmt 创建代码块
func BlockStmt(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: stmts,
	}
}

// ===== 测试表达式创建函数 =====

// TestExpr 创建测试表达式 [ ... ]
func TestExpr(x ast.Expr) *ast.TestExpr {
	return &ast.TestExpr{
		Lbrack: token.Pos(1),
		X:      x,
		Rbrack: token.Pos(1),
	}
}

// ExtendedTestExpr 创建扩展测试表达式 [[ ... ]]
func ExtendedTestExpr(x ast.Expr) *ast.ExtendedTestExpr {
	return &ast.ExtendedTestExpr{
		Lbrack: token.Pos(1),
		X:      x,
		Rbrack: token.Pos(1),
	}
}

// ===== 变量和参数扩展 =====

// VarExp 创建变量扩展表达式 $var
func VarExp(name ast.Expr) *ast.VarExpExpr {
	return &ast.VarExpExpr{
		Dollar: token.Pos(1),
		X:      name,
	}
}

// ParamExp 创建参数扩展表达式 ${var}
func ParamExp(varName ast.Expr, paramExp ast.ParamExp) *ast.ParamExpExpr {
	return &ast.ParamExpExpr{
		Dollar:   token.Pos(1),
		Lbrace:   token.Pos(1),
		Var:      varName,
		ParamExp: paramExp,
		Rbrace:   token.Pos(1),
	}
}

// ===== 比较操作符 =====

// Eq 创建相等比较
func Eq(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.Equal, y)
}

// Ne 创建不等比较
func Ne(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.NotEqual, y)
}

// Lt 创建小于比较
func Lt(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.TsLss, y)
}

// Gt 创建大于比较
func Gt(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.TsGtr, y)
}

// Le 创建小于等于比较
func Le(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.TsLeq, y)
}

// Ge 创建大于等于比较
func Ge(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.TsGeq, y)
}

// ===== 文件测试操作符 =====

// FileExists 创建文件存在测试
func FileExists(file ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-e"), token.TsExists, file)
}

// IsFile 创建普通文件测试
func IsFile(file ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-f"), token.TsRegFile, file)
}

// IsDir 创建目录测试
func IsDir(file ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-d"), token.TsDirect, file)
}

// IsReadable 创建可读测试
func IsReadable(file ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-r"), token.TsRead, file)
}

// IsWritable 创建可写测试
func IsWritable(file ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-w"), token.TsWrite, file)
}

// IsExecutable 创建可执行测试
func IsExecutable(file ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-x"), token.TsExec, file)
}

// ===== 字符串测试操作符 =====

// IsEmpty 创建空字符串测试
func IsEmpty(str ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-z"), token.TsEmpStr, str)
}

// IsNotEmpty 创建非空字符串测试
func IsNotEmpty(str ast.Expr) *ast.BinaryExpr {
	return Binary(Word("-n"), token.TsNempStr, str)
}

// ===== 常用命令 =====

// Echo 创建 echo 命令
func Echo(args ...ast.Expr) *ast.CmdExpr {
	return CmdExpr("echo", args...)
}

// Exit 创建 exit 命令
func Exit(args ...ast.Expr) *ast.CmdExpr {
	if len(args) == 0 {
		return CmdExpr("exit")
	}
	return CmdExpr("exit", args...)
}

// Break 创建 break 命令
func Break(n ...int) *ast.CmdExpr {
	if len(n) == 0 {
		return CmdExpr("break")
	}
	return CmdExpr("break", Word(strconv.Itoa(n[0])))
}

// Continue 创建 continue 命令
func Continue(n ...int) *ast.CmdExpr {
	if len(n) == 0 {
		return CmdExpr("continue")
	}
	return CmdExpr("continue", Word(strconv.Itoa(n[0])))
}

// True 创建 true 命令
func True() *ast.CmdExpr {
	return CmdExpr("true")
}

// False 创建 false 命令
func False() *ast.CmdExpr {
	return CmdExpr("false")
}

// ===== 逻辑操作符 =====

// And 创建逻辑与
func And(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.AndAnd, y)
}

// Or 创建逻辑或
func Or(x, y ast.Expr) *ast.BinaryExpr {
	return Binary(x, token.OrOr, y)
}

// Not 创建逻辑非
func Not(x ast.Expr) *ast.UnaryExpr {
	return Unary(token.ExclMark, x)
}

func Option(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}
