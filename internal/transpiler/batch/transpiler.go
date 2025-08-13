// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package batch

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/internal/module"
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/internal/unsafe"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	"github.com/hulo-lang/hulo/syntax/batch/astutil"
	btok "github.com/hulo-lang/hulo/syntax/batch/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

func Transpile(opts *config.Huloc, fs vfs.VFS, main string) (map[string]string, error) {
	moduleMgr := module.NewDependecyResolver(opts, fs)
	err := module.ResolveAllDependencies(moduleMgr, opts.Main)
	if err != nil {
		return nil, err
	}
	transpiler := NewTranspiler(opts.CompilerOptions.Batch, moduleMgr).RegisterDefaultRules()
	err = transpiler.BindRules()
	if err != nil {
		return nil, err
	}

	results := make(map[string]bast.Node)
	err = moduleMgr.VisitModules(func(mod *module.Module) error {
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
		ret[path] = bast.String(node)
	}
	return ret, nil
}

type UnresolvedSymbol struct {
	RefCount int // 引用次数，如果为0，则可以删除
	AST      *hast.UnresolvedSymbol
	node     *bast.Lit
}

func (sym *UnresolvedSymbol) Name() string {
	return sym.AST.Symbol
}

func (sym *UnresolvedSymbol) Source() string {
	return sym.AST.Path
}

func (sym *UnresolvedSymbol) Link(symbol *linker.LinkableSymbol) error {
	if symbol == nil {
		return fmt.Errorf("symbol is nil")
	}
	sym.node.Val = symbol.Text()
	return nil
}

type Transpiler struct {
	opts *config.BatchOptions

	buffer []bast.Stmt

	moduleMgr *module.DependecyResolver

	results map[string]*bast.File

	callers container.Stack[CallFrame]

	unresolvedSymbols map[string][]linker.UnkownSymbol

	hcrDispatcher *transpiler.HCRDispatcher[bast.Node]

	counter uint64

	unsafeEngine *unsafe.TemplateEngine
}

func NewTranspiler(opts *config.BatchOptions, moduleMgr *module.DependecyResolver) *Transpiler {
	return &Transpiler{
		counter:           0,
		opts:              opts,
		moduleMgr:         moduleMgr,
		unresolvedSymbols: make(map[string][]linker.UnkownSymbol),
		hcrDispatcher:     transpiler.NewHCRDispatcher[bast.Node](),
		callers:           container.NewArrayStack[CallFrame](),
		unsafeEngine:      unsafe.NewTemplateEngine(),
	}
}

func (t *Transpiler) HCRDispatcher() *transpiler.HCRDispatcher[bast.Node] {
	return t.hcrDispatcher
}

func (t *Transpiler) RegisterDefaultRules() *Transpiler {
	t.hcrDispatcher.Register(RuleCommentSyntax, &REMCommentHandler{}, &DoubleColonCommentHandler{})
	t.hcrDispatcher.Register(RuleBoolFormat, &BoolAsNumberHandler{}, &BoolAsStringHandler{}, &BoolAsCmdHandler{})
	return t
}

func (t *Transpiler) BindRules() error {
	opts := t.opts
	err := t.hcrDispatcher.Bind(RuleCommentSyntax, opts.CommentSyntax)
	if err != nil {
		return err
	}
	err = t.hcrDispatcher.Bind(RuleBoolFormat, opts.BoolFormat)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transpiler) UnresolvedSymbols() map[string][]linker.UnkownSymbol {
	return t.unresolvedSymbols
}

func (t *Transpiler) GetTargetExt() string {
	return t.opts.ExtFileName
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
	case *hast.TrueLiteral, *hast.FalseLiteral:
		return t.invokeHCR(RuleBoolFormat, node)
	case *hast.BreakStmt:
		return t.ConvertBreakStmt(node)
	case *hast.ContinueStmt:
		return t.ConvertContinueStmt(node)
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
	case *hast.CmdExpr:
		return t.ConvertCmdExpr(node)
	case *hast.ReturnStmt:
		return t.ConvertReturnStmt(node)
	case *hast.MatchStmt:
		return t.ConvertMatchStmt(node)
	case *hast.IncDecExpr:
		return t.ConvertIncDecExpr(node)
	case *hast.AssignStmt:
		return t.ConvertAssignStmt(node)
	case *hast.DoWhileStmt:
		return t.ConvertDoWhileStmt(node)
	case *hast.WhileStmt:
		return t.ConvertWhileStmt(node)
	case *hast.IfStmt:
		return t.ConvertIfStmt(node)
	case *hast.BlockStmt:
		return t.ConvertBlockStmt(node)
	case *hast.FuncDecl:
		return t.ConvertFuncDecl(node)
	case *hast.BinaryExpr:
		return t.ConvertBinaryExpr(node)
	case *hast.CallExpr:
		return t.ConvertCallExpr(node)
	case *hast.CommentGroup:
		return t.ConvertCommentGroup(node)
	case *hast.Import:
		return t.ConvertImport(node)
	case *hast.RefExpr:
		return t.ConvertRefExpr(node)
	case *hast.ArrayLiteralExpr:
		return t.ConvertArrayLiteralExpr(node)
	case *hast.ForeachStmt:
		return t.ConvertForeachStmt(node)
	case *hast.UnresolvedSymbol:
		return t.ConvertUnresolvedSymbol(node)
	case *hast.UnsafeExpr:
		return t.ConvertUnsafeExpr(node)
	default:
		panic(fmt.Sprintf("unsupported node type: %T", node))
	}
}

func (t *Transpiler) ConvertUnsafeExpr(node *hast.UnsafeExpr) bast.Node {
	// TODO 这里需要模板引擎执行，但是模板引擎缺少编译期变量，没法正确解析，需要解释器导出他的环境表
	output, err := t.unsafeEngine.Execute(node.Text)
	if err != nil {
		return nil
	}
	return astutil.Lit(output)
}

func (t *Transpiler) ConvertContinueStmt(node *hast.ContinueStmt) bast.Node {
	caller, ok := t.callers.Peek()
	if ok && caller.Kind() == CallerLoop {
		return astutil.Goto(caller.(*LoopFrame).StartLabel)
	}
	return nil
}

func (t *Transpiler) ConvertBreakStmt(node *hast.BreakStmt) bast.Node {
	caller, ok := t.callers.Peek()
	if ok && caller.Kind() == CallerLoop {
		return astutil.Goto(caller.(*LoopFrame).EndLabel)
	}
	return nil
}

func (t *Transpiler) ConvertUnresolvedSymbol(node *hast.UnresolvedSymbol) bast.Node {
	sym := astutil.Lit("[PLACEHOLDER]")
	t.unresolvedSymbols[node.Path] = append(t.unresolvedSymbols[node.Path],
		&UnresolvedSymbol{AST: node, node: sym})
	return astutil.ExprStmt(sym)
}

func (t *Transpiler) ConvertForStmt(node *hast.ForStmt) bast.Node {
	loopName := fmt.Sprintf("loop_%d", t.counter)
	endLabel := fmt.Sprintf("end_%d", t.counter)
	t.counter++

	t.callers.Push(&LoopFrame{
		StartLabel: loopName,
		EndLabel:   endLabel,
	})
	defer t.callers.Pop()

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

	// Create the loop structure
	loopStmts := []bast.Stmt{
		init,                    // initialization
		astutil.Label(loopName), // loop start label
		astutil.If(astutil.Unary(btok.NOT, cond), astutil.Block(astutil.Goto(endLabel)), nil),
		body,                    // loop body
		post,                    // post iteration
		astutil.Goto(loopName),  // jump back to loop start
		astutil.Label(endLabel), // loop end label
	}

	return astutil.Block(loopStmts...)
}

func (bt *Transpiler) ConvertForeachStmt(node *hast.ForeachStmt) bast.Node {
	bt.callers.Push(&LoopFrame{})
	defer bt.callers.Pop()

	x := bt.Convert(node.Index).(bast.Expr)
	elems := bt.Convert(node.Var).(bast.Expr)

	body := bt.Convert(node.Body).(*bast.BlockStmt)
	return astutil.For(x, astutil.Word(astutil.Lit("("), elems, astutil.Lit(")")), body)
}

func (bt *Transpiler) ConvertArrayLiteralExpr(node *hast.ArrayLiteralExpr) bast.Node {
	parts := make([]bast.Expr, len(node.Elems))
	for i, r := range node.Elems {
		parts[i] = bt.Convert(r).(bast.Expr)
	}
	return astutil.Word(parts...)
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
	scope, ok := bt.callers.Peek()
	if ok && scope.Kind() == CallerAssign {
		return astutil.Lit(node.X.(*hast.Ident).Name)
	}
	return astutil.DblQuote(astutil.Lit(node.X.(*hast.Ident).Name))
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
		bt.Emit(astutil.ExprStmt(expr))
	}
	return astutil.ExprStmt(astutil.Word(astutil.Lit("goto"), astutil.Lit(":eof")))
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
		ifStmt := astutil.If(astutil.Binary(astutil.DblQuote(astutil.Lit(name)), btok.EQU, cond), body, nil)
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
	bt.callers.Push(&AssignFrame{})
	x := bt.Convert(node.X).(bast.Expr)
	bt.callers.Pop()
	if node.Tok == htok.INC {
		return astutil.Assign(x, astutil.Word(astutil.DblQuote(x), astutil.Lit("+1")))
	}
	return astutil.Assign(x, astutil.Word(astutil.DblQuote(x), astutil.Lit("-1")))
}

func (bt *Transpiler) ConvertAssignStmt(node *hast.AssignStmt) bast.Node {
	lhs := bt.Convert(node.Lhs).(bast.Expr)
	rhs := bt.Convert(node.Rhs).(bast.Expr)
	return astutil.Assign(lhs, rhs)
}

func (bt *Transpiler) ConvertDoWhileStmt(node *hast.DoWhileStmt) bast.Node {
	loopName := fmt.Sprintf("loop_%d", bt.counter)
	bt.counter++
	// emulate do-while: body; for ;; cond; do (body)
	bt.callers.Push(&LoopFrame{
		StartLabel: loopName,
	})
	defer bt.callers.Pop()

	body := bt.Convert(node.Body).(bast.Stmt)
	cond := bt.Convert(node.Cond).(bast.Expr)

	ifStmt := astutil.If(astutil.Unary(btok.NOT, cond), astutil.Block(astutil.Goto(loopName)), nil)
	return astutil.Block(astutil.Label(loopName), body, ifStmt)
}

func (bt *Transpiler) ConvertWhileStmt(node *hast.WhileStmt) bast.Node {
	loopName := fmt.Sprintf("loop_%d", bt.counter)
	endName := fmt.Sprintf("end_%d", bt.counter)
	bt.counter++

	bt.callers.Push(&LoopFrame{
		StartLabel: loopName,
		EndLabel:   endName,
	})
	defer bt.callers.Pop()

	stmts := []bast.Stmt{astutil.Goto(loopName)}
	if node.Cond != nil {
		cond := bt.Convert(node.Cond).(bast.Expr)
		ifStmt := astutil.If(astutil.Unary(btok.NOT, cond), astutil.Block(astutil.Goto(endName)), nil)
		stmts = append(stmts, ifStmt)
	}

	body := bt.Convert(node.Body).(bast.Stmt)
	stmts = append(stmts, body, astutil.Goto(endName))

	return astutil.Block(stmts...)
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
	return astutil.If(cond, body, elseStmt)
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
	return astutil.Block(stmts...)
}

func (bt *Transpiler) ConvertFuncDecl(node *hast.FuncDecl) bast.Node {
	bt.callers.Push(&FunctionFrame{})
	defer bt.callers.Pop()
	// batch function: :label ...
	body := bt.Convert(node.Body).(*bast.BlockStmt)
	return astutil.FuncDecl(node.Name.Name, body)
}

func (bt *Transpiler) ConvertBinaryExpr(node *hast.BinaryExpr) bast.Node {
	switch node.Op {
	case htok.PLUS, htok.MINUS, htok.ASTERISK, htok.SLASH, htok.MOD:
		var lhs, rhs bast.Expr
		top, _ := bt.callers.Peek()

		switch top.Kind() {
		case CallerFunction:
			lhs = astutil.SglQuote(astutil.Lit("1"))
			rhs = astutil.SglQuote(astutil.Lit("2"))
		case CallerLoop:
			bt.Emit(&bast.ExprStmt{})
			return nil
		default:
			lhs = bt.Convert(node.X).(bast.Expr)
			rhs = bt.Convert(node.Y).(bast.Expr)
		}

		bt.Emit(&bast.AssignStmt{
			Opt: astutil.Lit("/a"),
			Lhs: astutil.Lit("result"),
			Rhs: astutil.Binary(lhs, Token(node.Op), rhs),
		})
		return nil
	default:
		return astutil.Binary(bt.Convert(node.X).(bast.Expr), Token(node.Op), bt.Convert(node.Y).(bast.Expr))
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
			bt.Emit(astutil.ExprStmt(astutil.Word(append(parts, args...)...)))
			return astutil.DblQuote(astutil.Lit("result"))
		}
	}
	// 其它情况按原有逻辑
	recv := make([]bast.Expr, len(node.Recv))
	for i, r := range node.Recv {
		recv[i] = bt.Convert(r).(bast.Expr)
	}
	return &bast.CallExpr{Fun: bt.Convert(node.Fun).(bast.Expr), Recv: recv}
}

func (bt *Transpiler) invokeHCR(name transpiler.RuleID, node hast.Node) bast.Node {
	rule, err := bt.hcrDispatcher.Get(name)
	if err != nil {
		return nil
	}
	converted, err := rule.Apply(bt, node)
	if err != nil {
		return nil
	}
	return converted
}

func (bt *Transpiler) ConvertCommentGroup(node *hast.CommentGroup) bast.Node {
	return bt.invokeHCR(RuleCommentSyntax, node)
}

func (bt *Transpiler) ConvertFile(node *hast.File) bast.Node {
	docs := make([]*bast.CommentGroup, len(node.Docs))
	for i, d := range node.Docs {
		docs[i] = bt.Convert(d).(*bast.CommentGroup)
	}
	stmts := []bast.Stmt{astutil.ExprStmt(astutil.Word(astutil.Lit("@echo"), astutil.Lit("off")))}
	for _, s := range node.Stmts {
		stmt := bt.Convert(s)

		stmts = append(stmts, bt.Flush()...)
		if s == nil {
			continue
		}
		stmts = append(stmts, stmt.(bast.Stmt))
	}
	return astutil.File(docs, stmts...)
}

func (bt *Transpiler) Emit(n ...bast.Stmt) {
	bt.buffer = append(bt.buffer, n...)
}

func (bt *Transpiler) Flush() []bast.Stmt {
	stmts := bt.buffer
	bt.buffer = nil
	return stmts
}
