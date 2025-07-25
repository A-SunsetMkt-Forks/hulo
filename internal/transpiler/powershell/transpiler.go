// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"
	"runtime"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/interpreter"
	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	past "github.com/hulo-lang/hulo/syntax/powershell/ast"
	pstok "github.com/hulo-lang/hulo/syntax/powershell/token"
)

var psTokenMap = map[htok.Token]pstok.Token{
	// 比较运算符
	htok.EQ:  pstok.EQ,
	htok.NEQ: pstok.NE,
	htok.LT:  pstok.LT,
	htok.LE:  pstok.LE,
	htok.GT:  pstok.GT,
	htok.GE:  pstok.GE,
	// 逻辑运算符
	htok.AND: pstok.AND,
	htok.OR:  pstok.OR,
	// 算术运算符
	htok.PLUS:     pstok.ADD,
	htok.MINUS:    pstok.SUB,
	htok.ASTERISK: pstok.MUL,
	htok.SLASH:    pstok.DIV,
	// 赋值
	htok.ASSIGN: pstok.ASSIGN,
	htok.DEC:    pstok.DEC,
	htok.INC:    pstok.INC,
}

func psToken(tok htok.Token) pstok.Token {
	if t, ok := psTokenMap[tok]; ok {
		return t
	}
	return 0 // 默认返回 0
}

type PowerShellTranspiler struct {
	opts *config.Huloc

	buffer []past.Stmt

	moduleManager *ModuleManager

	vfs vfs.VFS

	modules        map[string]*Module
	currentModule  *Module
	builtinModules []*Module

	globalSymbols *SymbolTable

	scopeStack container.Stack[ScopeType]

	currentScope *Scope

	// 跟踪类实例变量名到类名的映射
	classInstances map[string]string
}

func Transpile(opts *config.Huloc, vfs vfs.VFS, basePath string) (map[string]string, error) {
	transpiler := NewPowerShellTranspiler(opts, vfs)
	mainFile := opts.Main
	results, err := transpiler.Transpile(mainFile)
	return results, err
}

func NewPowerShellTranspiler(opts *config.Huloc, vfs vfs.VFS) *PowerShellTranspiler {
	env := interpreter.NewEnvironment()
	env.SetWithScope("TARGET", &object.StringValue{Value: "powershell"}, htok.CONST, true)
	env.SetWithScope("OS", &object.StringValue{Value: runtime.GOOS}, htok.CONST, true)
	env.SetWithScope("ARCH", &object.StringValue{Value: runtime.GOARCH}, htok.CONST, true)
	interp := interpreter.NewInterpreter(env)

	return &PowerShellTranspiler{
		opts: opts,
		vfs:  vfs,
		// 初始化模块管理
		moduleManager: &ModuleManager{
			options:  opts,
			vfs:      vfs,
			basePath: "",
			modules:  make(map[string]*Module),
			resolver: &DependencyResolver{
				visited: container.NewMapSet[string](),
				stack:   container.NewMapSet[string](),
				order:   make([]string, 0),
			},
			interp: interp,
		},
		modules: make(map[string]*Module),
		// 初始化符号管理
		globalSymbols: NewSymbolTable("global", opts.EnableMangle),

		// 初始化作用域管理
		scopeStack: container.NewArrayStack[ScopeType](),

		// 初始化类实例映射
		classInstances: make(map[string]string),
	}
}

func (p *PowerShellTranspiler) Transpile(mainFile string) (map[string]string, error) {
	// 1. 解析所有依赖
	modules, err := p.moduleManager.ResolveAllDependencies(mainFile)
	if err != nil {
		return nil, err
	}

	// 2. 按依赖顺序翻译模块
	results := make(map[string]string)
	for _, module := range modules {
		psAst := p.Convert(module.AST)
		if psAst == nil {
			return nil, err
		}
		// 3. 生成 PowerShell 脚本内容
		var script string

		script = past.String(psAst)
		path := module.Path
		if len(path) > 3 && path[len(path)-3:] == ".hl" {
			path = path[:len(path)-3] + ".ps1"
		}
		results[path] = script
	}
	return results, nil
}

// Hulo AST (hast.Node) -> PowerShell AST (past.Node) 分发器和典型节点转换

func (p *PowerShellTranspiler) Convert(node hast.Node) past.Node {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *hast.File:
		return p.ConvertFile(n)
	case *hast.FuncDecl:
		return p.ConvertFuncDecl(n)
	case *hast.AssignStmt:
		return p.ConvertAssignStmt(n)
	case *hast.IfStmt:
		return p.ConvertIfStmt(n)
	case *hast.ExprStmt:
		return p.ConvertExprStmt(n)
	case *hast.Ident:
		return p.ConvertIdent(n)
	case *hast.StringLiteral:
		return p.ConvertStringLiteral(n)
	case *hast.NumericLiteral:
		return p.ConvertNumericLiteral(n)
	case *hast.CallExpr:
		return p.ConvertCallExpr(n)
	case *hast.CmdExpr:
		return p.ConvertCmdExpr(n)
	case *hast.BlockStmt:
		return p.ConvertBlockStmt(n)
	case *hast.BinaryExpr:
		return p.ConvertBinaryExpr(n)
	case *hast.WhileStmt:
		return p.ConvertWhileStmt(n)
	case *hast.ForStmt:
		return p.ConvertForStmt(n)
	case *hast.ForeachStmt:
		return p.ConvertForeachStmt(n)
	case *hast.ReturnStmt:
		return p.ConvertReturnStmt(n)
	case *hast.TrueLiteral:
		return p.ConvertBoolLiteral(n)
	case *hast.FalseLiteral:
		return p.ConvertFalseLiteral(n)
	case *hast.Parameter:
		return p.ConvertParameter(n)
	case *hast.ArrayLiteralExpr:
		return p.ConvertArrayLiteralExpr(n)
	case *hast.ObjectLiteralExpr:
		return p.ConvertObjectLiteralExpr(n)
	case *hast.SelectExpr:
		return p.ConvertSelectExpr(n)
	case *hast.IncDecExpr:
		return p.ConvertIncDecExpr(n)
	case *hast.DoWhileStmt:
		return p.ConvertDoWhileStmt(n)
	case *hast.MatchStmt:
		return p.ConvertSwitchStmt(n)
	// case *hast.CaseClause:
	// 	return p.ConvertCaseClause(n)
	case *hast.BreakStmt:
		return p.ConvertBreakStmt(n)
	case *hast.ContinueStmt:
		return p.ConvertContinueStmt(n)
	case *hast.ThrowStmt:
		return p.ConvertThrowStmt(n)
	case *hast.TryStmt:
		return p.ConvertTryStmt(n)
	case *hast.CatchClause:
		return p.ConvertCatchClause(n)
	case *hast.ClassDecl:
		return p.ConvertClassDecl(n)
	case *hast.EnumDecl:
		return p.ConvertEnumDecl(n)
	case *hast.CommentGroup:
		return p.ConvertCommentGroup(n)
	case *hast.Comment:
		return p.ConvertComment(n)
	case *hast.RefExpr:
		return p.ConvertRefExpr(n)
	case *hast.Import:
		return p.ConvertImport(n)
	default:
		fmt.Printf("Unhandled node type: %T\n", node)
		return nil
	}
}

func (p *PowerShellTranspiler) ConvertImport(node *hast.Import) past.Node {
	// 简单实现：将 import 转为 PowerShell 注释或 source 语句
	// 这里只处理单一 import 路径
	var path string
	if node.ImportSingle != nil {
		path = node.ImportSingle.Path
	} else if node.ImportMulti != nil {
		path = node.ImportMulti.Path
	} else if node.ImportAll != nil {
		path = node.ImportAll.Path
	} else {
		path = "unknown_module"
	}

	path += ".ps1"

	// PowerShell 的 source 语法：. ./module.ps1
	return &past.ExprStmt{X: &past.CmdExpr{
		Cmd:  &past.Ident{Name: "."},
		Args: []past.Expr{&past.StringLit{Val: path}},
	}}
}

func (p *PowerShellTranspiler) ConvertFile(node *hast.File) past.Node {
	docs := make([]*past.CommentGroup, len(node.Docs))
	for i, d := range node.Docs {
		docs[i] = p.Convert(d).(*past.CommentGroup)
	}
	var stmts []past.Stmt
	for _, stmt := range node.Stmts {
		converted := p.Convert(stmt)
		if converted != nil {
			if s, ok := converted.(past.Stmt); ok {
				stmts = append(stmts, s)
			}
		}
	}
	return &past.File{
		Docs:  docs,
		Stmts: stmts,
	}
}

func (p *PowerShellTranspiler) ConvertFuncDecl(node *hast.FuncDecl) past.Node {
	body := p.Convert(node.Body)
	if body == nil {
		return nil
	}
	bodyBlock, ok := body.(*past.BlockStmt)
	if !ok {
		return nil
	}
	return &past.FuncDecl{
		Name: &past.Ident{Name: node.Name.Name},
		Body: bodyBlock,
	}
}

func (p *PowerShellTranspiler) ConvertAssignStmt(node *hast.AssignStmt) past.Node {
	lhs := p.Convert(node.Lhs)
	rhs := p.Convert(node.Rhs)

	if lhs == nil || rhs == nil {
		return nil
	}

	lhsExpr, ok1 := lhs.(past.Expr)
	rhsExpr, ok2 := rhs.(past.Expr)
	if !ok1 || !ok2 {
		return nil
	}

	return &past.AssignStmt{
		Lhs: &past.VarExpr{X: lhsExpr},
		Rhs: rhsExpr,
	}
}

func (p *PowerShellTranspiler) ConvertIfStmt(node *hast.IfStmt) past.Node {
	condNode := p.Convert(node.Cond)
	if condNode == nil {
		return nil
	}
	cond, ok1 := condNode.(past.Expr)
	if !ok1 {
		return nil
	}
	bodyNode := p.Convert(node.Body)
	if bodyNode == nil {
		return nil
	}
	body, ok2 := bodyNode.(*past.BlockStmt)
	if !ok2 {
		return nil
	}
	var elseStmt past.Stmt
	if node.Else != nil {
		converted := p.Convert(node.Else)
		if s, ok := converted.(past.Stmt); ok {
			elseStmt = s
		}
	}
	return &past.IfStmt{
		Cond: cond,
		Body: body,
		Else: elseStmt,
	}
}

func (p *PowerShellTranspiler) ConvertBlockStmt(node *hast.BlockStmt) past.Node {
	var stmts []past.Stmt
	for _, stmt := range node.List {
		converted := p.Convert(stmt)
		if converted == nil {
			continue // 跳过 nil 子节点
		}
		if s, ok := converted.(past.Stmt); ok {
			stmts = append(stmts, s)
		}
	}
	return &past.BlockStmt{List: stmts}
}

func (p *PowerShellTranspiler) ConvertExprStmt(node *hast.ExprStmt) past.Node {
	x := p.Convert(node.X)
	if x == nil {
		return nil // 这里可以保留，表达式语句为 nil 就不输出
	}
	expr, ok := x.(past.Expr)
	if !ok {
		return nil
	}
	return &past.ExprStmt{X: expr}
}

func (p *PowerShellTranspiler) ConvertIdent(node *hast.Ident) past.Node {
	return &past.Ident{Name: node.Name}
}

func (p *PowerShellTranspiler) ConvertStringLiteral(node *hast.StringLiteral) past.Node {
	return &past.StringLit{Val: node.Value}
}

func (p *PowerShellTranspiler) ConvertNumericLiteral(node *hast.NumericLiteral) past.Node {
	return &past.NumericLit{Val: node.Value}
}

func (p *PowerShellTranspiler) ConvertCallExpr(node *hast.CallExpr) past.Node {
	fun := p.Convert(node.Fun)
	if fun == nil {
		return nil
	}
	funExpr, ok := fun.(past.Expr)
	if !ok {
		return nil
	}
	var recv []past.Expr
	for _, arg := range node.Recv {
		converted := p.Convert(arg)
		if converted == nil {
			return nil
		}
		expr, ok := converted.(past.Expr)
		if !ok {
			return nil
		}
		recv = append(recv, expr)
	}
	return &past.CallExpr{Func: funExpr, Recv: recv}
}

func (p *PowerShellTranspiler) ConvertCmdExpr(node *hast.CmdExpr) past.Node {
	fun := p.Convert(node.Cmd)
	if fun == nil {
		return nil
	}
	funExpr, ok := fun.(past.Expr)
	if !ok {
		return nil
	}
	var recv []past.Expr
	for _, arg := range node.Args {
		converted := p.Convert(arg)
		if converted == nil {
			return nil
		}
		expr, ok := converted.(past.Expr)
		if !ok {
			return nil
		}
		recv = append(recv, expr)
	}
	return &past.CallExpr{Func: funExpr, Recv: recv}
}

func (p *PowerShellTranspiler) ConvertBinaryExpr(node *hast.BinaryExpr) past.Node {
	x := p.Convert(node.X)
	y := p.Convert(node.Y)
	if x == nil || y == nil {
		return nil
	}
	xExpr, ok1 := x.(past.Expr)
	yExpr, ok2 := y.(past.Expr)
	if !ok1 || !ok2 {
		return nil
	}
	return &past.BinaryExpr{X: xExpr, Op: psToken(node.Op), Y: yExpr}
}

func (p *PowerShellTranspiler) ConvertWhileStmt(node *hast.WhileStmt) past.Node {
	cond := p.Convert(node.Cond)
	body := p.Convert(node.Body)
	if cond == nil || body == nil {
		return nil
	}
	condExpr, ok1 := cond.(past.Expr)
	bodyBlock, ok2 := body.(*past.BlockStmt)
	if !ok1 || !ok2 {
		return nil
	}
	return &past.WhileStmt{
		Cond: condExpr,
		Body: bodyBlock,
	}
}

func (p *PowerShellTranspiler) ConvertForStmt(node *hast.ForStmt) past.Node {
	init := p.Convert(node.Init)
	cond := p.Convert(node.Cond)
	post := p.Convert(node.Post)
	body := p.Convert(node.Body)
	if init == nil || cond == nil || post == nil || body == nil {
		return nil
	}
	initExpr, ok1 := init.(past.Expr)
	condExpr, ok2 := cond.(past.Expr)
	postExpr, ok3 := post.(past.Expr)
	bodyBlock, ok4 := body.(*past.BlockStmt)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil
	}
	return &past.ForStmt{
		Init: initExpr,
		Cond: condExpr,
		Post: postExpr,
		Body: bodyBlock,
	}
}

func (p *PowerShellTranspiler) ConvertForeachStmt(node *hast.ForeachStmt) past.Node {
	index := p.Convert(node.Index)
	varExpr := p.Convert(node.Var)
	body := p.Convert(node.Body)
	if index == nil || varExpr == nil || body == nil {
		return nil
	}
	indexExpr, ok1 := index.(past.Expr)
	varExprExpr, ok2 := varExpr.(past.Expr)
	bodyBlock, ok3 := body.(*past.BlockStmt)
	if !ok1 || !ok2 || !ok3 {
		return nil
	}
	return &past.ForeachStmt{
		Elm:  indexExpr,
		Elms: varExprExpr,
		Body: bodyBlock,
	}
}

func (p *PowerShellTranspiler) ConvertReturnStmt(node *hast.ReturnStmt) past.Node {
	if node.X == nil {
		return &past.ReturnStmt{}
	}
	x := p.Convert(node.X)
	if x == nil {
		return nil
	}
	xExpr, ok := x.(past.Expr)
	if !ok {
		return nil
	}
	return &past.ReturnStmt{X: xExpr}
}

func (p *PowerShellTranspiler) ConvertBoolLiteral(node *hast.TrueLiteral) past.Node {
	return &past.BoolLit{Val: true}
}
func (p *PowerShellTranspiler) ConvertFalseLiteral(node *hast.FalseLiteral) past.Node {
	return &past.BoolLit{Val: false}
}

func (p *PowerShellTranspiler) ConvertParameter(node *hast.Parameter) past.Node {
	name := p.Convert(node.Name)
	if name == nil {
		return nil
	}
	nameExpr, ok := name.(past.Expr)
	if !ok {
		return nil
	}
	return &past.Parameter{X: nameExpr}
}

func (p *PowerShellTranspiler) ConvertArrayLiteralExpr(node *hast.ArrayLiteralExpr) past.Node {
	var elems []past.Expr
	for _, e := range node.Elems {
		converted := p.Convert(e)
		if converted == nil {
			return nil
		}
		expr, ok := converted.(past.Expr)
		if !ok {
			return nil
		}
		elems = append(elems, expr)
	}
	return &past.ArrayExpr{Elems: elems}
}

func (p *PowerShellTranspiler) ConvertObjectLiteralExpr(node *hast.ObjectLiteralExpr) past.Node {
	var entries []*past.HashEntry
	for _, prop := range node.Props {
		if kv, ok := prop.(*hast.KeyValueExpr); ok {
			key := p.Convert(kv.Key)
			value := p.Convert(kv.Value)
			if key == nil || value == nil {
				return nil
			}

			// 处理键，可能是字符串字面量或标识符
			var keyIdent *past.Ident
			if keyStr, ok := key.(*past.StringLit); ok {
				// 将字符串字面量转换为标识符
				keyIdent = &past.Ident{Name: keyStr.Val}
			} else if keyId, ok := key.(*past.Ident); ok {
				keyIdent = keyId
			} else {
				return nil
			}

			valueExpr, ok2 := value.(past.Expr)
			if !ok2 {
				return nil
			}
			entries = append(entries, &past.HashEntry{Key: keyIdent, Value: valueExpr})
		}
	}
	return &past.HashTable{Entries: entries}
}

func (p *PowerShellTranspiler) ConvertSelectExpr(node *hast.SelectExpr) past.Node {
	x := p.Convert(node.X)
	y := p.Convert(node.Y)
	if x == nil || y == nil {
		return nil
	}
	xExpr, ok1 := x.(past.Expr)
	yExpr, ok2 := y.(past.Expr)
	if !ok1 || !ok2 {
		return nil
	}
	return &past.SelectExpr{X: xExpr, Sel: yExpr}
}

func (p *PowerShellTranspiler) ConvertIncDecExpr(node *hast.IncDecExpr) past.Node {
	x := p.Convert(node.X)
	if x == nil {
		return nil
	}
	xExpr, ok := x.(past.Expr)
	if !ok {
		return nil
	}
	return &past.IncDecExpr{X: xExpr}
}

func (p *PowerShellTranspiler) ConvertDoWhileStmt(node *hast.DoWhileStmt) past.Node {
	body := p.Convert(node.Body)
	cond := p.Convert(node.Cond)
	if body == nil || cond == nil {
		return nil
	}
	bodyBlock, ok1 := body.(*past.BlockStmt)
	condExpr, ok2 := cond.(past.Expr)
	if !ok1 || !ok2 {
		return nil
	}
	return &past.DoWhileStmt{Body: bodyBlock, Cond: condExpr}
}

// 暂不实现 DoUntilStmt

func (p *PowerShellTranspiler) ConvertSwitchStmt(node *hast.MatchStmt) past.Node {
	value := p.Convert(node.Expr)
	if value == nil {
		return nil
	}
	valueExpr, ok := value.(past.Expr)
	if !ok {
		return nil
	}
	var cases []*past.CaseClause
	for _, c := range node.Cases {
		caseNode := p.ConvertCaseClause(c)
		if caseNode == nil {
			return nil
		}
		caseClause, ok := caseNode.(*past.CaseClause)
		if !ok {
			return nil
		}
		cases = append(cases, caseClause)
	}

	var defaultBody *past.BlockStmt
	if node.Default != nil {
		if node.Default.Body != nil {
			converted := p.Convert(node.Default.Body)
			if block, ok := converted.(*past.BlockStmt); ok {
				defaultBody = block
			} else {
				defaultBody = &past.BlockStmt{
					List: []past.Stmt{converted.(past.Stmt)},
				}
			}
		} else {
			defaultBody = &past.BlockStmt{List: []past.Stmt{}}
		}
	}

	return &past.SwitchStmt{Value: valueExpr, Cases: cases, Default: defaultBody}
}

func (p *PowerShellTranspiler) ConvertCaseClause(node *hast.CaseClause) past.Node {
	cond := p.Convert(node.Cond)
	body := p.Convert(node.Body)
	if cond == nil || body == nil {
		return nil
	}
	condExpr, ok1 := cond.(past.Expr)
	bodyBlock, ok2 := body.(*past.BlockStmt)
	if !ok1 || !ok2 {
		return nil
	}
	return &past.CaseClause{Cond: condExpr, Body: bodyBlock}
}

func (p *PowerShellTranspiler) ConvertBreakStmt(node *hast.BreakStmt) past.Node {
	return nil // 或自定义 past.ExprStmt{X: ...}
}

func (p *PowerShellTranspiler) ConvertContinueStmt(node *hast.ContinueStmt) past.Node {
	return nil // 或自定义 past.ExprStmt{X: ...}
}

func (p *PowerShellTranspiler) ConvertThrowStmt(node *hast.ThrowStmt) past.Node {
	return nil // 或自定义 past.ExprStmt{X: ...}
}

func (p *PowerShellTranspiler) ConvertTryStmt(node *hast.TryStmt) past.Node {
	body := p.Convert(node.Body)
	var catches []*past.CatchClause
	for _, c := range node.Catches {
		catches = append(catches, p.ConvertCatchClause(c).(*past.CatchClause))
	}
	return &past.TryStmt{Body: body.(*past.BlockStmt), Catches: catches}
}

func (p *PowerShellTranspiler) ConvertCatchClause(node *hast.CatchClause) past.Node {
	body := p.Convert(node.Body)
	return &past.CatchClause{Body: body.(*past.BlockStmt)}
}

// TrapStmt 暂不实现

func (p *PowerShellTranspiler) ConvertClassDecl(node *hast.ClassDecl) past.Node {
	// // 假设 node.Block 是类体
	// var body past.Node
	// if node.Block != nil {
	// 	body = p.Convert(node.Block)
	// } else {
	// 	body = &past.BlcokStmt{List: nil}
	// }
	// return &past.ClassDecl{Name: &past.Ident{Name: node.Name.Name}, Body: body.(*past.BlcokStmt)}
	return nil
}

func (p *PowerShellTranspiler) ConvertEnumDecl(node *hast.EnumDecl) past.Node {
	// // 假设 node.Block 是枚举体
	// var body past.Node
	// if node.Block != nil {
	// 	body = p.Convert(node.Block)
	// } else {
	// 	body = &past.BlcokStmt{List: nil}
	// }
	// return &past.EnumDecl{Name: &past.Ident{Name: node.Name.Name}, Body: body.(*past.BlcokStmt)}
	return nil
}

func (p *PowerShellTranspiler) ConvertCommentGroup(node *hast.CommentGroup) past.Node {
	var comments []past.Comment
	for _, c := range node.List {
		converted := p.Convert(c)
		if converted != nil {
			if pc, ok := converted.(past.Comment); ok {
				comments = append(comments, pc)
			}
		}
	}
	return &past.CommentGroup{List: comments}
}

func (p *PowerShellTranspiler) ConvertComment(node *hast.Comment) past.Node {
	// 假设 hast.Comment 有 Text 字段
	return &past.SingleLineComment{Text: node.Text}
}

func (p *PowerShellTranspiler) ConvertRefExpr(node *hast.RefExpr) past.Node {
	x := p.Convert(node.X)
	if x == nil {
		return nil
	}
	xExpr, ok := x.(past.Expr)
	if !ok {
		return nil
	}
	return &past.VarExpr{X: xExpr}
}
