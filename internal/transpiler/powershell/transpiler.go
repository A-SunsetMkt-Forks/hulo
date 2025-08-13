// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pwsh

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/internal/module"
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	past "github.com/hulo-lang/hulo/syntax/powershell/ast"
)

type Transpiler struct {
	opts *config.PowerShellOptions

	buffer []past.Stmt

	moduleMgr *module.DependecyResolver

	results map[string]*past.File

	callers container.Stack[CallFrame]

	unresolvedSymbols map[string][]linker.UnkownSymbol

	hcrDispatcher *transpiler.HCRDispatcher[past.Node]

	currentModule *module.Module

	// 跟踪类实例变量名到类名的映射
	classInstances map[string]string
}

func Transpile(opts *config.Huloc, fs vfs.VFS, main string) (map[string]string, error) {
	moduleMgr := module.NewDependecyResolver(opts, fs)
	err := module.ResolveAllDependencies(moduleMgr, opts.Main)
	if err != nil {
		return nil, err
	}
	transpiler := NewTranspiler(opts.CompilerOptions.Pwsh, moduleMgr).RegisterDefaultRules()
	err = transpiler.BindRules()
	if err != nil {
		return nil, err
	}

	results := make(map[string]past.Node)
	err = moduleMgr.VisitModules(func(mod *module.Module) error {
		transpiler.currentModule = mod
		batAST := transpiler.Convert(mod.AST)
		results[strings.Replace(mod.Path, ".hl", transpiler.GetTargetExt(), 1)] = batAST
		return nil
	})
	if err != nil {
		return nil, err
	}
	ld := linker.NewLinker(fs)
	ld.Listen(".bat", linker.BeginEnd{Begin: "REM HULO_LINK_BEGIN", End: "REM HULO_LINK_END"})
	ld.Listen(".cmd", linker.BeginEnd{Begin: "REM HULO_LINK_BEGIN", End: "REM HULO_LINK_END"})
	for _, nodes := range transpiler.unresolvedSymbols {
		for _, node := range nodes {
			// 先Load没有在Read, 并且没符号的时候全部导入所有符号，并且还要处理好外部文件和内部文件
			err := ld.Read(node.Source())
			if err != nil {
				return nil, err
			}
			linkable := ld.Load(node.Source())
			symbol := linkable.Lookup(node.Name())
			if symbol == nil {
				return nil, fmt.Errorf("symbol %s not found in %s", node.Name(), node.Source())
			}
			err = node.Link(symbol)
			if err != nil {
				return nil, err
			}
		}
	}
	ret := make(map[string]string)
	for path, node := range results {
		ret[path] = past.String(node)
	}
	return ret, nil
}

func NewTranspiler(opts *config.PowerShellOptions, moduleMgr *module.DependecyResolver) *Transpiler {
	return &Transpiler{
		opts:              opts,
		moduleMgr:         moduleMgr,
		unresolvedSymbols: make(map[string][]linker.UnkownSymbol),
		hcrDispatcher:     transpiler.NewHCRDispatcher[past.Node](),
		callers:           container.NewArrayStack[CallFrame](),
		classInstances:    make(map[string]string),
	}
}

func (p *Transpiler) RegisterDefaultRules() *Transpiler {
	return p
}

func (p *Transpiler) BindRules() error {
	return nil
}

func (p *Transpiler) UnresolvedSymbols() map[string][]linker.UnkownSymbol {
	return p.unresolvedSymbols
}

func (p *Transpiler) GetTargetExt() string {
	return ".ps1"
}

func (p *Transpiler) GetTargetName() string {
	return "powershell"
}

func (p *Transpiler) Convert(node hast.Node) past.Node {
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
	case *hast.UnsafeExpr:
		return p.ConvertUnsafeExpr(n)
	case *hast.ExternDecl:
		return p.ConvertExternDecl(n)
	default:
		fmt.Printf("Unhandled node type: %T\n", node)
		return nil
	}
}

func (p *Transpiler) ConvertExternDecl(node *hast.ExternDecl) past.Node {
	return nil
}

func (p *Transpiler) ConvertUnsafeExpr(node *hast.UnsafeExpr) past.Node {
	return &past.Ident{Name: node.Text}
}

func (p *Transpiler) ConvertImport(node *hast.Import) past.Node {
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

func (p *Transpiler) ConvertFile(node *hast.File) past.Node {
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

func (p *Transpiler) ConvertFuncDecl(node *hast.FuncDecl) past.Node {
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

func (p *Transpiler) ConvertAssignStmt(node *hast.AssignStmt) past.Node {
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

func (p *Transpiler) ConvertIfStmt(node *hast.IfStmt) past.Node {
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

func (p *Transpiler) ConvertBlockStmt(node *hast.BlockStmt) past.Node {
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

func (p *Transpiler) ConvertExprStmt(node *hast.ExprStmt) past.Node {
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

func (p *Transpiler) ConvertIdent(node *hast.Ident) past.Node {
	return &past.Ident{Name: node.Name}
}

func (p *Transpiler) ConvertStringLiteral(node *hast.StringLiteral) past.Node {
	return &past.StringLit{Val: node.Value}
}

func (p *Transpiler) ConvertNumericLiteral(node *hast.NumericLiteral) past.Node {
	return &past.NumericLit{Val: node.Value}
}

func (p *Transpiler) ConvertCallExpr(node *hast.CallExpr) past.Node {
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

func (p *Transpiler) ConvertCmdExpr(node *hast.CmdExpr) past.Node {
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

func (p *Transpiler) ConvertBinaryExpr(node *hast.BinaryExpr) past.Node {
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

func (p *Transpiler) ConvertWhileStmt(node *hast.WhileStmt) past.Node {
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

func (p *Transpiler) ConvertForStmt(node *hast.ForStmt) past.Node {
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

func (p *Transpiler) ConvertForeachStmt(node *hast.ForeachStmt) past.Node {
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

func (p *Transpiler) ConvertReturnStmt(node *hast.ReturnStmt) past.Node {
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

func (p *Transpiler) ConvertBoolLiteral(node *hast.TrueLiteral) past.Node {
	return &past.BoolLit{Val: true}
}
func (p *Transpiler) ConvertFalseLiteral(node *hast.FalseLiteral) past.Node {
	return &past.BoolLit{Val: false}
}

func (p *Transpiler) ConvertParameter(node *hast.Parameter) past.Node {
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

func (p *Transpiler) ConvertArrayLiteralExpr(node *hast.ArrayLiteralExpr) past.Node {
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

func (p *Transpiler) ConvertObjectLiteralExpr(node *hast.ObjectLiteralExpr) past.Node {
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

func (p *Transpiler) ConvertSelectExpr(node *hast.SelectExpr) past.Node {
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

func (p *Transpiler) ConvertIncDecExpr(node *hast.IncDecExpr) past.Node {
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

func (p *Transpiler) ConvertDoWhileStmt(node *hast.DoWhileStmt) past.Node {
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

func (p *Transpiler) ConvertSwitchStmt(node *hast.MatchStmt) past.Node {
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

func (p *Transpiler) ConvertCaseClause(node *hast.CaseClause) past.Node {
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

func (p *Transpiler) ConvertBreakStmt(node *hast.BreakStmt) past.Node {
	return nil // 或自定义 past.ExprStmt{X: ...}
}

func (p *Transpiler) ConvertContinueStmt(node *hast.ContinueStmt) past.Node {
	return nil // 或自定义 past.ExprStmt{X: ...}
}

func (p *Transpiler) ConvertThrowStmt(node *hast.ThrowStmt) past.Node {
	return nil // 或自定义 past.ExprStmt{X: ...}
}

func (p *Transpiler) ConvertTryStmt(node *hast.TryStmt) past.Node {
	body := p.Convert(node.Body)
	var catches []*past.CatchClause
	for _, c := range node.Catches {
		catches = append(catches, p.ConvertCatchClause(c).(*past.CatchClause))
	}
	return &past.TryStmt{Body: body.(*past.BlockStmt), Catches: catches}
}

func (p *Transpiler) ConvertCatchClause(node *hast.CatchClause) past.Node {
	body := p.Convert(node.Body)
	return &past.CatchClause{Body: body.(*past.BlockStmt)}
}

// TrapStmt 暂不实现

func (p *Transpiler) ConvertClassDecl(node *hast.ClassDecl) past.Node {
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

func (p *Transpiler) ConvertEnumDecl(node *hast.EnumDecl) past.Node {
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

func (p *Transpiler) ConvertCommentGroup(node *hast.CommentGroup) past.Node {
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

func (p *Transpiler) ConvertComment(node *hast.Comment) past.Node {
	// 假设 hast.Comment 有 Text 字段
	return &past.SingleLineComment{Text: node.Text}
}

func (p *Transpiler) ConvertRefExpr(node *hast.RefExpr) past.Node {
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
