// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	btok "github.com/hulo-lang/hulo/syntax/batch/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

func Transpile(opts *config.Huloc, vfs vfs.VFS, basePath string) (map[string]string, error) {
	transpiler := NewBatchTranspiler(opts, vfs)
	mainFile := opts.Main
	results, err := transpiler.Transpile(mainFile)
	return results, err
}

type BatchTranspiler struct {
	opts *config.Huloc

	buffer []bast.Stmt

	moduleManager *ModuleManager

	vfs vfs.VFS

	modules        map[string]*Module
	currentModule  *Module
	builtinModules []*Module

	globalSymbols *SymbolTable

	scopeStack container.Stack[ScopeType]

	currentScope *Scope
}

func NewBatchTranspiler(opts *config.Huloc, vfs vfs.VFS) *BatchTranspiler {
	return &BatchTranspiler{
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
		},
		modules: make(map[string]*Module),
		// 初始化符号管理
		globalSymbols: NewSymbolTable("global", opts.EnableMangle),

		// 初始化作用域管理
		scopeStack: container.NewArrayStack[ScopeType](),
	}
}

func (bt *BatchTranspiler) Transpile(mainFile string) (map[string]string, error) {
	modules, err := bt.moduleManager.ResolveAllDependencies(mainFile)
	if err != nil {
		return nil, err
	}

	// 2. 按依赖顺序翻译模块
	results := make(map[string]string)
	for _, module := range modules {
		bt.currentModule = module
		psAst := bt.Convert(module.AST)
		if psAst == nil {
			return nil, err
		}

		var script string

		script = bast.String(psAst)
		path := module.Path
		if len(path) > 3 && path[len(path)-3:] == ".hl" {
			path = path[:len(path)-3] + ".bat"
		}
		results[path] = script
	}
	return results, nil
}

func (bt *BatchTranspiler) Convert(node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		return bt.ConvertFile(node)
	case *hast.CommentGroup:
		return bt.ConvertCommentGroup(node)
	case *hast.Comment:
		return &bast.Comment{Text: node.Text}
	case *hast.ExprStmt:
		return &bast.ExprStmt{X: bt.Convert(node.X).(bast.Expr)}
	case *hast.AssignStmt:
		lhs := bt.Convert(node.Lhs).(bast.Expr)
		rhs := bt.Convert(node.Rhs).(bast.Expr)
		return &bast.AssignStmt{Lhs: lhs, Rhs: rhs}
	case *hast.RefExpr:
		return &bast.Lit{Val: node.X.(*hast.Ident).Name}
	case *hast.CmdExpr:
		args := make([]bast.Expr, len(node.Args))
		for i, r := range node.Args {
			args[i] = bt.Convert(r).(bast.Expr)
		}
		return &bast.CmdExpr{Name: bt.Convert(node.Cmd).(bast.Expr), Recv: args}
	case *hast.CallExpr:
		return bt.ConvertCallExpr(node)
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
	case *hast.IfStmt:
		return bt.ConvertIfStmt(node)
	case *hast.BinaryExpr:
		return bt.ConvertBinaryExpr(node)
	case *hast.BlockStmt:
		return bt.ConvertBlockStmt(node)
	case *hast.WhileStmt:
		// emulate while with for and if
		// for ;; cond; do (body)
		cond := bt.Convert(node.Cond).(bast.Expr)
		body := bt.Convert(node.Body).(bast.Stmt)
		return &bast.ForStmt{
			X:    &bast.Lit{Val: ""},
			List: &bast.Lit{Val: ""},
			Body: &bast.BlockStmt{List: []bast.Stmt{
				&bast.IfStmt{Cond: cond, Body: body},
			}},
		}
	case *hast.DoWhileStmt:
		// emulate do-while: body; for ;; cond; do (body)
		body := bt.Convert(node.Body).(bast.Stmt)
		cond := bt.Convert(node.Cond).(bast.Expr)
		return &bast.BlockStmt{List: []bast.Stmt{
			body.(bast.Stmt),
			&bast.ForStmt{
				X:    &bast.Lit{Val: ""},
				List: &bast.Lit{Val: ""},
				Body: &bast.BlockStmt{List: []bast.Stmt{
					&bast.IfStmt{Cond: cond, Body: body},
				}},
			},
		}}
	case *hast.ForStmt:
		// emulate for: init; for ;; cond; do (body; post)
		init := bt.Convert(node.Init).(bast.Stmt)
		cond := bt.Convert(node.Cond).(bast.Expr)
		post := bt.Convert(node.Post).(bast.Expr)
		body := bt.Convert(node.Body).(bast.Stmt)
		return &bast.BlockStmt{List: []bast.Stmt{
			init,
			&bast.ForStmt{
				X:    &bast.Lit{Val: ""},
				List: &bast.Lit{Val: ""},
				Body: &bast.BlockStmt{List: []bast.Stmt{
					&bast.IfStmt{Cond: cond, Body: &bast.BlockStmt{List: []bast.Stmt{body, &bast.ExprStmt{X: post}}}},
				}},
			},
		}}
	case *hast.ForeachStmt:
		// batch 不原生支持 foreach，降级为注释
		return &bast.Comment{Text: "foreach not supported in batch, skipped"}
	case *hast.ReturnStmt:
		expr, ok := bt.Convert(node.X).(bast.Expr)
		if ok {
			bt.Emit(&bast.ExprStmt{X: expr})
		}
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":eof"}}}}
	case *hast.BreakStmt:
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":break"}}}}
	case *hast.ContinueStmt:
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":continue"}}}}
	case *hast.ArrayLiteralExpr:
		return &bast.Lit{Val: "[array]"}
	case *hast.ObjectLiteralExpr:
		return &bast.Lit{Val: "[object]"}
	case *hast.UnaryExpr:
		x := bt.Convert(node.X).(bast.Expr)
		return &bast.UnaryExpr{Op: Token(node.Op), X: x}
	case *hast.IncDecExpr:
		x := bt.Convert(node.X).(bast.Expr)
		if node.Tok.String() == "++" {
			return &bast.AssignStmt{Lhs: x, Rhs: &bast.Word{Parts: []bast.Expr{x, &bast.Lit{Val: "+1"}}}}
		} else {
			return &bast.AssignStmt{Lhs: x, Rhs: &bast.Word{Parts: []bast.Expr{x, &bast.Lit{Val: "-1"}}}}
		}
	case *hast.LabelStmt:
		return &bast.LabelStmt{Name: node.Name.Name}
	case *hast.FuncDecl:
		return bt.ConvertFuncDecl(node)
	case *hast.SelectExpr:
		// 只在 CallExpr 里处理，这里降级为字面量
		return &bast.Lit{Val: "[select not supported]"}
	case *hast.Import:
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
	case *hast.MatchStmt:
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
	case *hast.TryStmt, *hast.CatchClause, *hast.ThrowStmt:
		return &bast.Comment{Text: "try/catch/throw not supported in batch, skipped"}
	case *hast.DeclareDecl, *hast.ExternDecl, *hast.Parameter, *hast.ComptimeStmt:
		return &bast.Comment{Text: "declaration/extern/parameter/comptime not supported in batch, skipped"}
	default:
		fmt.Printf("unsupported node type: %T\n", node)
		return &bast.Lit{Val: "[unsupported]"}
	}
}

func (bt *BatchTranspiler) ConvertIfStmt(node *hast.IfStmt) bast.Node {
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

func (bt *BatchTranspiler) ConvertBlockStmt(node *hast.BlockStmt) bast.Node {
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

func (bt *BatchTranspiler) ConvertFuncDecl(node *hast.FuncDecl) bast.Node {
	bt.currentModule.Symbols.PushScope(FunctionScope, bt.currentModule)
	defer bt.currentModule.Symbols.PopScope()
	// batch function: :label ...
	body := bt.Convert(node.Body).(*bast.BlockStmt)
	return &bast.FuncDecl{Name: node.Name.Name, Body: body}
}

func (bt *BatchTranspiler) ConvertBinaryExpr(node *hast.BinaryExpr) bast.Node {
	switch node.Op {
	case htok.PLUS, htok.MINUS, htok.ASTERISK, htok.SLASH, htok.MOD:
		var lhs, rhs bast.Expr
		top := bt.currentModule.Symbols.CurrentScope()

		if top.Type == FunctionScope {
			lhs = &bast.SglQuote{Val: &bast.Lit{Val: "1"}}
			rhs = &bast.SglQuote{Val: &bast.Lit{Val: "2"}}
		} else {
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

func (bt *BatchTranspiler) ConvertCallExpr(node *hast.CallExpr) bast.Node {
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

func (bt *BatchTranspiler) ConvertCommentGroup(node *hast.CommentGroup) bast.Node {
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

func (bt *BatchTranspiler) ConvertFile(node *hast.File) bast.Node {
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

func (bt *BatchTranspiler) Emit(n ...bast.Stmt) {
	bt.buffer = append(bt.buffer, n...)
}

func (bt *BatchTranspiler) Flush() []bast.Stmt {
	stmts := bt.buffer
	bt.buffer = nil
	return stmts
}
