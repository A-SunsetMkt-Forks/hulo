// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/module"
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	btok "github.com/hulo-lang/hulo/syntax/batch/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

func Transpile(opts *config.Huloc, fs vfs.VFS, main string) (map[string]string, error) {
	moduleMgr := module.NewDependecyResolver()
	err := module.ResolveAllDependencies(moduleMgr, opts.Main)
	if err != nil {
		return nil, err
	}
	transpiler := &Transpiler{
		opts:      opts,
		moduleMgr: moduleMgr,
	}
	results := make(map[string]string)
	err = moduleMgr.VisitModules(func(mod *module.Module) error {
		bashAST := transpiler.Convert(mod.AST)
		results[mod.Path] = bast.String(bashAST)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

type UnresolvedSymbol struct {
	File     string
	RefCount int // 引用次数，如果为0，则可以删除
}

type Transpiler struct {
	opts *config.Huloc

	buffer []bast.Stmt

	moduleMgr *module.DependecyResolver

	results map[string]*bast.File

	// ! 一定要重建作用域，在元编程的时候会破坏作用域，如果被元编程破坏的话 还要重新修复语法树
	// 如果假设语法树修复了，能不能从 Node -> Scope 的映射获取呢？
	// 这里为了简便设计 先手动管理作用域
	scopes *ScopeStack

	unresolvedSymbols map[string][]string

	// HCRMgr
	strategies map[string]transpiler.Strategy[any]

	counter uint64
}

func (t *Transpiler) UnresolvedSymbols() []any {
	return nil
}

func (t *Transpiler) GetTargetType() bast.Node {
	return &bast.File{}
}

func (t *Transpiler) GetTargetExt() string {
	return ".bat"
}

func (t *Transpiler) GetTargetName() string {
	return "batch"
}

func (t *Transpiler) Convert(node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		return t.ConvertFile(node)
	case *hast.ForStmt:
		return t.ConvertForStmt(node)
	case *hast.Comment:
		return &bast.Comment{Text: node.Text}
	case *hast.ExprStmt:
		return &bast.ExprStmt{X: t.Convert(node.X).(bast.Expr)}
	case *hast.Ident:
		return &bast.Lit{Val: node.Name}
	case *hast.StringLiteral:
		return &bast.Lit{Val: node.Value}
	case *hast.NumericLiteral:
		return &bast.Lit{Val: node.Value}
	case *hast.TrueLiteral:
		return &bast.Lit{Val: "true"}
	case *hast.FalseLiteral:
		return &bast.Lit{Val: "false"}
	case *hast.BreakStmt:
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":break"}}}}
	case *hast.ContinueStmt:
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":continue"}}}}
	case *hast.ObjectLiteralExpr:
		return &bast.Lit{Val: "[object]"}
	case *hast.UnaryExpr:
		x := t.Convert(node.X).(bast.Expr)
		return &bast.UnaryExpr{Op: Token(node.Op), X: x}
	case *hast.LabelStmt:
		return &bast.LabelStmt{Name: node.Name.Name}
	case *hast.SelectExpr:
		// 只在 CallExpr 里处理，这里降级为字面量
		return &bast.Lit{Val: "[select not supported]"}
	case *hast.TryStmt, *hast.CatchClause, *hast.ThrowStmt:
		return &bast.Comment{Text: "try/catch/throw not supported in batch, skipped"}
	case *hast.DeclareDecl, *hast.ExternDecl, *hast.Parameter, *hast.ComptimeStmt:
		return &bast.Comment{Text: "declaration/extern/parameter/comptime not supported in batch, skipped"}
	default:
		panic(fmt.Sprintf("unsupported node type: %T", node))
	}
}

func (t *Transpiler) ConvertForStmt(node *hast.ForStmt) bast.Node {
	t.scopes.Push(&Scope{Type: LoopScope})
	defer t.scopes.Pop()

	loopName := fmt.Sprintf("loop_%d", t.counter)
	t.counter++

	// Convert for loop components
	init := t.Convert(node.Init).(bast.Stmt)
	cond := t.Convert(node.Cond).(bast.Expr)
	post := t.Convert(node.Post).(bast.Stmt)
	body := t.Convert(node.Body).(bast.Stmt)

	// For loop structure in batch:
	// init
	// :loop_label
	// if not cond goto :end
	// body
	// post
	// goto :loop_label
	// :end

	endLabel := fmt.Sprintf("end_%d", t.counter-1)

	// Create the loop structure
	loopStmts := []bast.Stmt{
		init,                            // initialization
		&bast.LabelStmt{Name: loopName}, // loop start label
		&bast.IfStmt{ // condition check
			Cond: &bast.UnaryExpr{Op: btok.NOT, X: cond},
			Body: &bast.BlockStmt{List: []bast.Stmt{
				&bast.GotoStmt{Label: endLabel},
			}},
		},
		body,                            // loop body
		post,                            // post iteration
		&bast.GotoStmt{Label: loopName}, // jump back to loop start
		&bast.LabelStmt{Name: endLabel}, // loop end label
	}

	return &bast.BlockStmt{List: loopStmts}
}

func (bt *Transpiler) ConvertForeachStmt(node *hast.ForeachStmt) bast.Node {
	bt.scopes.Push(&Scope{Type: LoopScope})
	defer bt.scopes.Pop()

	x := bt.Convert(node.Index).(bast.Expr)
	elems := bt.Convert(node.Var).(bast.Expr)

	body := bt.Convert(node.Body).(*bast.BlockStmt)
	return &bast.ForStmt{
		X:    x,
		List: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "("}, elems, &bast.Lit{Val: ")"}}},
		Body: body,
	}
}

func (bt *Transpiler) ConvertArrayLiteralExpr(node *hast.ArrayLiteralExpr) bast.Node {
	parts := make([]bast.Expr, len(node.Elems))
	for i, r := range node.Elems {
		parts[i] = bt.Convert(r).(bast.Expr)
	}
	return &bast.Word{Parts: parts}
}

func (bt *Transpiler) ConvertImport(node *hast.Import) bast.Node {
	if node.ImportSingle != nil {
		return &bast.CallStmt{
			Name:   node.ImportSingle.Path + ".bat",
			IsFile: true,
		}
	}
	if node.ImportAll != nil {
		return &bast.CallStmt{
			Name:   node.ImportAll.Path + ".bat",
			IsFile: true,
		}
	}
	return nil
}

func (bt *Transpiler) ConvertRefExpr(node *hast.RefExpr) bast.Node {
	scope := bt.scopes.Current()
	if scope != nil && scope.Type == AssignScope {
		return &bast.Lit{Val: node.X.(*hast.Ident).Name}
	}
	return &bast.DblQuote{Val: &bast.Lit{Val: node.X.(*hast.Ident).Name}}
}

func (bt *Transpiler) ConvertCmdExpr(node *hast.CmdExpr) bast.Node {
	args := make([]bast.Expr, len(node.Args))
	for i, r := range node.Args {
		args[i] = bt.Convert(r).(bast.Expr)
	}
	return &bast.CmdExpr{Name: bt.Convert(node.Cmd).(bast.Expr), Recv: args}
}

func (bt *Transpiler) ConvertReturnStmt(node *hast.ReturnStmt) bast.Node {
	expr, ok := bt.Convert(node.X).(bast.Expr)
	if ok {
		bt.Emit(&bast.ExprStmt{X: expr})
	}
	return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":eof"}}}}
}

func (bt *Transpiler) ConvertMatchStmt(node *hast.MatchStmt) bast.Node {
	// 用 if-else 链模拟 match
	var firstIf *bast.IfStmt
	var currentIf *bast.IfStmt
	var name string
	if id, ok := node.Expr.(*hast.Ident); ok {
		name = id.Name
	} else if ref, ok := node.Expr.(*hast.RefExpr); ok {
		name = ref.X.(*hast.Ident).Name
	}
	for _, cc := range node.Cases {
		cond := bt.Convert(cc.Cond).(bast.Expr)
		body := bt.Convert(cc.Body).(bast.Stmt)
		ifStmt := &bast.IfStmt{Cond: &bast.BinaryExpr{X: &bast.DblQuote{Val: &bast.Lit{Val: name}}, Op: btok.EQU, Y: cond}, Body: body}
		if firstIf == nil {
			firstIf = ifStmt
			currentIf = ifStmt
		} else {
			currentIf.Else = ifStmt
			currentIf = ifStmt
		}
	}
	if node.Default != nil {
		defaultBody := bt.Convert(node.Default.Body).(bast.Stmt)
		currentIf.Else = defaultBody
	}
	if firstIf == nil {
		return &bast.Lit{Val: "[empty match]"}
	}
	return firstIf
}

func (bt *Transpiler) ConvertIncDecExpr(node *hast.IncDecExpr) bast.Node {
	bt.scopes.Push(&Scope{Type: AssignScope})
	x := bt.Convert(node.X).(bast.Expr)
	bt.scopes.Pop()
	if node.Tok == htok.INC {
		return &bast.AssignStmt{Lhs: x, Rhs: &bast.Word{Parts: []bast.Expr{&bast.DblQuote{Val: x}, &bast.Lit{Val: "+1"}}}}
	} else {
		return &bast.AssignStmt{Lhs: x, Rhs: &bast.Word{Parts: []bast.Expr{&bast.DblQuote{Val: x}, &bast.Lit{Val: "-1"}}}}
	}
}

func (bt *Transpiler) ConvertAssignStmt(node *hast.AssignStmt) bast.Node {
	lhs := bt.Convert(node.Lhs).(bast.Expr)
	rhs := bt.Convert(node.Rhs).(bast.Expr)
	return &bast.AssignStmt{Lhs: lhs, Rhs: rhs}
}

func (bt *Transpiler) ConvertDoWhileStmt(node *hast.DoWhileStmt) bast.Node {
	// emulate do-while: body; for ;; cond; do (body)
	bt.scopes.Push(&Scope{Type: LoopScope})
	defer bt.scopes.Pop()

	loopName := fmt.Sprintf("loop_%d", bt.counter)
	bt.counter++

	body := bt.Convert(node.Body).(bast.Stmt)
	cond := bt.Convert(node.Cond).(bast.Expr)

	ifStmt := &bast.IfStmt{Cond: cond, Body: &bast.BlockStmt{List: []bast.Stmt{&bast.GotoStmt{Label: loopName}}}}
	return &bast.BlockStmt{List: []bast.Stmt{&bast.LabelStmt{Name: loopName}, body, ifStmt}}
}

func (bt *Transpiler) ConvertWhileStmt(node *hast.WhileStmt) bast.Node {
	bt.scopes.Push(&Scope{Type: LoopScope})
	defer bt.scopes.Pop()

	loopName := fmt.Sprintf("loop_%d", bt.counter)
	endName := fmt.Sprintf("end_%d", bt.counter)
	bt.counter++

	stmts := []bast.Stmt{&bast.GotoStmt{Label: loopName}}
	if node.Cond != nil {
		cond := bt.Convert(node.Cond).(bast.Expr)
		ifStmt := &bast.IfStmt{Cond: &bast.UnaryExpr{Op: btok.NOT, X: cond}, Body: &bast.BlockStmt{List: []bast.Stmt{&bast.GotoStmt{Label: endName}}}}
		stmts = append(stmts, ifStmt)
	}

	body := bt.Convert(node.Body).(bast.Stmt)
	stmts = append(stmts, body, &bast.GotoStmt{Label: endName})

	return &bast.BlockStmt{List: stmts}
}

func (bt *Transpiler) ConvertIfStmt(node *hast.IfStmt) bast.Node {
	condNode := bt.Convert(node.Cond)
	if condNode == nil {
		return nil
	}
	cond, ok1 := condNode.(bast.Expr)
	if !ok1 {
		return nil
	}
	bodyNode := bt.Convert(node.Body)
	if bodyNode == nil {
		return nil
	}
	body, ok2 := bodyNode.(*bast.BlockStmt)
	if !ok2 {
		return nil
	}
	var elseStmt bast.Stmt
	if node.Else != nil {
		converted := bt.Convert(node.Else)
		if s, ok := converted.(bast.Stmt); ok {
			elseStmt = s
		}
	}
	return &bast.IfStmt{
		Cond: cond,
		Body: body,
		Else: elseStmt,
	}
}

func (bt *Transpiler) ConvertBlockStmt(node *hast.BlockStmt) bast.Node {
	stmts := make([]bast.Stmt, 0, len(node.List))
	for _, s := range node.List {
		stmt := bt.Convert(s)
		stmts = append(stmts, bt.Flush()...)
		if stmt == nil {
			continue
		}
		if s, ok := stmt.(bast.Stmt); ok {
			stmts = append(stmts, s)
		} else {
			// 不是vast.Stmt类型，忽略或报错
			// 可以选择panic或者continue，这里选择continue
			continue
		}
	}
	return &bast.BlockStmt{List: stmts}
}

func (bt *Transpiler) ConvertFuncDecl(node *hast.FuncDecl) bast.Node {
	bt.scopes.Push(&Scope{Type: FunctionScope})
	defer bt.scopes.Pop()
	// batch function: :label ...
	body := bt.Convert(node.Body).(*bast.BlockStmt)
	return &bast.FuncDecl{Name: node.Name.Name, Body: body}
}

func (bt *Transpiler) ConvertBinaryExpr(node *hast.BinaryExpr) bast.Node {
	switch node.Op {
	case htok.PLUS, htok.MINUS, htok.ASTERISK, htok.SLASH, htok.MOD:
		var lhs, rhs bast.Expr
		top := bt.scopes.Current()

		switch top.Type {
		case FunctionScope:
			lhs = &bast.SglQuote{Val: &bast.Lit{Val: "1"}}
			rhs = &bast.SglQuote{Val: &bast.Lit{Val: "2"}}
		case LoopScope:
			bt.Emit(&bast.ExprStmt{})
			return nil
		default:
			lhs = bt.Convert(node.X).(bast.Expr)
			rhs = bt.Convert(node.Y).(bast.Expr)
		}

		bt.Emit(&bast.AssignStmt{
			Opt: &bast.Lit{Val: "/a"},
			Lhs: &bast.Lit{Val: "result"},
			Rhs: &bast.BinaryExpr{X: lhs, Op: Token(node.Op), Y: rhs},
		})
		return nil
	default:
		return &bast.BinaryExpr{X: bt.Convert(node.X).(bast.Expr), Op: Token(node.Op), Y: bt.Convert(node.Y).(bast.Expr)}
	}
}

func (bt *Transpiler) ConvertCallExpr(node *hast.CallExpr) bast.Node {
	if sel, ok := node.Fun.(*hast.SelectExpr); ok {
		_, ok1 := bt.Convert(sel.X).(*bast.Lit)
		funLit, ok2 := bt.Convert(sel.Y).(*bast.Lit)
		if ok1 && ok2 {
			args := make([]bast.Expr, 0, len(node.Recv))
			for _, r := range node.Recv {
				args = append(args, bt.Convert(r).(bast.Expr))
			}
			parts := []bast.Expr{
				&bast.Lit{Val: "call"},
				// &bast.Lit{Val: modLit.Val + ".bat"},
				&bast.Lit{Val: funLit.Val},
			}
			bt.Emit(&bast.ExprStmt{X: &bast.Word{Parts: append(parts, args...)}})
			return &bast.DblQuote{Val: &bast.Lit{Val: "result"}}
		}
	}
	// 其它情况按原有逻辑
	recv := make([]bast.Expr, len(node.Recv))
	for i, r := range node.Recv {
		recv[i] = bt.Convert(r).(bast.Expr)
	}
	return &bast.CallExpr{Fun: bt.Convert(node.Fun).(bast.Expr), Recv: recv}
}

func (bt *Transpiler) ConvertCommentGroup(node *hast.CommentGroup) bast.Node {
	var tok btok.Token
	if bt.opts.CompilerOptions.Batch != nil && bt.opts.CompilerOptions.Batch.CommentSyntax == "::" {
		tok = btok.DOUBLE_COLON
	} else {
		tok = btok.REM
	}
	docs := make([]*bast.Comment, len(node.List))
	for i, d := range node.List {
		docs[i] = &bast.Comment{Tok: tok, Text: d.Text}
	}
	return &bast.CommentGroup{Comments: docs}
}

func (bt *Transpiler) ConvertFile(node *hast.File) bast.Node {
	docs := make([]*bast.CommentGroup, len(node.Docs))
	for i, d := range node.Docs {
		docs[i] = bt.Convert(d).(*bast.CommentGroup)
	}
	stmts := []bast.Stmt{&bast.ExprStmt{
		X: &bast.Word{
			Parts: []bast.Expr{
				&bast.Lit{Val: "@echo"},
				&bast.Lit{Val: "off"},
			},
		},
	}}
	for _, s := range node.Stmts {
		stmt := bt.Convert(s)

		stmts = append(stmts, bt.Flush()...)
		if s == nil {
			continue
		}
		stmts = append(stmts, stmt.(bast.Stmt))
	}
	return &bast.File{Docs: docs, Stmts: stmts}
}

func (bt *Transpiler) Emit(n ...bast.Stmt) {
	bt.buffer = append(bt.buffer, n...)
}

func (bt *Transpiler) Flush() []bast.Stmt {
	stmts := bt.buffer
	bt.buffer = nil
	return stmts
}
