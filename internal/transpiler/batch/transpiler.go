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

	fs vfs.VFS

	modules        map[string]*Module
	currentModule  *Module
	builtinModules []*Module

	globalSymbols *SymbolTable

	counter uint64
}

func NewBatchTranspiler(opts *config.Huloc, vfs vfs.VFS) *BatchTranspiler {
	env := interpreter.NewEnvironment()
	env.SetWithScope("TARGET", &object.StringValue{Value: "batch"}, htok.CONST, true)
	env.SetWithScope("OS", &object.StringValue{Value: runtime.GOOS}, htok.CONST, true)
	env.SetWithScope("ARCH", &object.StringValue{Value: runtime.GOARCH}, htok.CONST, true)
	interp := interpreter.NewInterpreter(env)

	return &BatchTranspiler{
		opts: opts,
		fs:   vfs,
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
		return bt.ConvertAssignStmt(node)
	case *hast.RefExpr:
		return bt.ConvertRefExpr(node)
	case *hast.CmdExpr:
		return bt.ConvertCmdExpr(node)
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
		return bt.ConvertWhileStmt(node)
	case *hast.DoWhileStmt:
		return bt.ConvertDoWhileStmt(node)
	case *hast.ForStmt:
		return bt.ConvertForStmt(node)
	case *hast.ForeachStmt:
		return bt.ConvertForeachStmt(node)
	case *hast.ReturnStmt:
		return bt.ConvertReturnStmt(node)
	case *hast.BreakStmt:
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":break"}}}}
	case *hast.ContinueStmt:
		return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":continue"}}}}
	case *hast.ArrayLiteralExpr:
		return bt.ConvertArrayLiteralExpr(node)
	case *hast.ObjectLiteralExpr:
		return &bast.Lit{Val: "[object]"}
	case *hast.UnaryExpr:
		x := bt.Convert(node.X).(bast.Expr)
		return &bast.UnaryExpr{Op: Token(node.Op), X: x}
	case *hast.IncDecExpr:
		return bt.ConvertIncDecExpr(node)
	case *hast.LabelStmt:
		return &bast.LabelStmt{Name: node.Name.Name}
	case *hast.FuncDecl:
		return bt.ConvertFuncDecl(node)
	case *hast.SelectExpr:
		// 只在 CallExpr 里处理，这里降级为字面量
		return &bast.Lit{Val: "[select not supported]"}
	case *hast.Import:
		return bt.ConvertImport(node)
	case *hast.MatchStmt:
		return bt.ConvertMatchStmt(node)
	case *hast.TryStmt, *hast.CatchClause, *hast.ThrowStmt:
		return &bast.Comment{Text: "try/catch/throw not supported in batch, skipped"}
	case *hast.DeclareDecl, *hast.ExternDecl, *hast.Parameter, *hast.ComptimeStmt:
		return &bast.Comment{Text: "declaration/extern/parameter/comptime not supported in batch, skipped"}
	default:
		panic(fmt.Sprintf("unsupported node type: %T", node))
	}
}

func (bt *BatchTranspiler) ConvertForeachStmt(node *hast.ForeachStmt) bast.Node {
	bt.currentModule.Symbols.PushScope(LoopScope, bt.currentModule)
	defer bt.currentModule.Symbols.PopScope()

	x := bt.Convert(node.Index).(bast.Expr)
	elems := bt.Convert(node.Var).(bast.Expr)

	body := bt.Convert(node.Body).(*bast.BlockStmt)
	return &bast.ForStmt{
		X:    x,
		List: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "("}, elems, &bast.Lit{Val: ")"}}},
		Body: body,
	}
}

func (bt *BatchTranspiler) ConvertArrayLiteralExpr(node *hast.ArrayLiteralExpr) bast.Node {
	parts := make([]bast.Expr, len(node.Elems))
	for i, r := range node.Elems {
		parts[i] = bt.Convert(r).(bast.Expr)
	}
	return &bast.Word{Parts: parts}
}

func (bt *BatchTranspiler) ConvertImport(node *hast.Import) bast.Node {
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

func (bt *BatchTranspiler) ConvertRefExpr(node *hast.RefExpr) bast.Node {
	scope := bt.currentModule.Symbols.CurrentScope()
	if scope != nil && scope.Type == AssignScope {
		return &bast.Lit{Val: node.X.(*hast.Ident).Name}
	}
	return &bast.DblQuote{Val: &bast.Lit{Val: node.X.(*hast.Ident).Name}}
}

func (bt *BatchTranspiler) ConvertCmdExpr(node *hast.CmdExpr) bast.Node {
	args := make([]bast.Expr, len(node.Args))
	for i, r := range node.Args {
		args[i] = bt.Convert(r).(bast.Expr)
	}
	return &bast.CmdExpr{Name: bt.Convert(node.Cmd).(bast.Expr), Recv: args}
}

func (bt *BatchTranspiler) ConvertReturnStmt(node *hast.ReturnStmt) bast.Node {
	expr, ok := bt.Convert(node.X).(bast.Expr)
	if ok {
		bt.Emit(&bast.ExprStmt{X: expr})
	}
	return &bast.ExprStmt{X: &bast.Word{Parts: []bast.Expr{&bast.Lit{Val: "goto"}, &bast.Lit{Val: ":eof"}}}}
}

func (bt *BatchTranspiler) ConvertMatchStmt(node *hast.MatchStmt) bast.Node {
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

func (bt *BatchTranspiler) ConvertIncDecExpr(node *hast.IncDecExpr) bast.Node {
	bt.currentModule.Symbols.PushScope(AssignScope, bt.currentModule)
	x := bt.Convert(node.X).(bast.Expr)
	bt.currentModule.Symbols.PopScope()
	if node.Tok.String() == "++" {
		return &bast.AssignStmt{Lhs: x, Rhs: &bast.Word{Parts: []bast.Expr{&bast.DblQuote{Val: x}, &bast.Lit{Val: "+1"}}}}
	} else {
		return &bast.AssignStmt{Lhs: x, Rhs: &bast.Word{Parts: []bast.Expr{&bast.DblQuote{Val: x}, &bast.Lit{Val: "-1"}}}}
	}
}

func (bt *BatchTranspiler) ConvertAssignStmt(node *hast.AssignStmt) bast.Node {
	lhs := bt.Convert(node.Lhs).(bast.Expr)
	rhs := bt.Convert(node.Rhs).(bast.Expr)
	return &bast.AssignStmt{Lhs: lhs, Rhs: rhs}
}

func (bt *BatchTranspiler) ConvertDoWhileStmt(node *hast.DoWhileStmt) bast.Node {
	// emulate do-while: body; for ;; cond; do (body)
	bt.currentModule.Symbols.PushScope(LoopScope, bt.currentModule)
	defer bt.currentModule.Symbols.PopScope()

	loopName := fmt.Sprintf("loop_%d", bt.counter)
	bt.counter++

	body := bt.Convert(node.Body).(bast.Stmt)
	cond := bt.Convert(node.Cond).(bast.Expr)

	ifStmt := &bast.IfStmt{Cond: cond, Body: &bast.BlockStmt{List: []bast.Stmt{&bast.GotoStmt{Label: loopName}}}}
	return &bast.BlockStmt{List: []bast.Stmt{&bast.LabelStmt{Name: loopName}, body, ifStmt}}
}

func (bt *BatchTranspiler) ConvertWhileStmt(node *hast.WhileStmt) bast.Node {
	bt.currentModule.Symbols.PushScope(LoopScope, bt.currentModule)
	defer bt.currentModule.Symbols.PopScope()

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

func (bt *BatchTranspiler) ConvertForStmt(node *hast.ForStmt) bast.Node {
	bt.currentModule.Symbols.PushScope(LoopScope, bt.currentModule)
	defer bt.currentModule.Symbols.PopScope()

	loopName := fmt.Sprintf("loop_%d", bt.counter)
	bt.counter++

	// Convert for loop components
	init := bt.Convert(node.Init).(bast.Stmt)
	cond := bt.Convert(node.Cond).(bast.Expr)
	post := bt.Convert(node.Post).(bast.Stmt)
	body := bt.Convert(node.Body).(bast.Stmt)

	// For loop structure in batch:
	// init
	// :loop_label
	// if not cond goto :end
	// body
	// post
	// goto :loop_label
	// :end

	endLabel := fmt.Sprintf("end_%d", bt.counter-1)

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
