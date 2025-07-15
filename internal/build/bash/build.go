// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"
	"strings"

	"slices"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	btok "github.com/hulo-lang/hulo/syntax/bash/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

func Translate(opts *config.BashOptions, node hast.Node) (bast.Node, error) {
	bnode := translate2Bash(opts, node)
	return bnode, nil
}

func translate2Bash(opts *config.BashOptions, node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		stmts := []bast.Stmt{&bast.Comment{Text: "!/bin/bash"}}
		for _, s := range node.Stmts {
			stmts = append(stmts, translate2Bash(opts, s).(bast.Stmt))
		}

		return &bast.File{
			Stmts: stmts,
		}
	case *hast.IfStmt:
		// opts.Boolean = "number & >= 1.2.3"
		parse, err := hcrDispatcher.Get(opts.BooleanFormat)
		if err != nil {
			return nil
		}
		cond, err := parse.Apply(&hast.IfStmt{}, node.Cond)
		if err != nil {
			return nil
		}
		return &bast.IfStmt{
			Cond: cond.(bast.Expr),
			Body: translate2Bash(opts, node.Body).(*bast.BlockStmt),
		}
	case *hast.BlockStmt:
		return &bast.BlockStmt{}
	case *hast.BinaryExpr:
		return &bast.BinaryExpr{
			X:  translate2Bash(opts, node.X).(bast.Expr),
			Op: Token(node.Op),
			Y:  translate2Bash(opts, node.Y).(bast.Expr),
		}
	case *hast.RefExpr:
		return &bast.VarExpExpr{
			X: translate2Bash(opts, node.X).(*bast.Ident),
		}
	case *hast.Ident:
		return &bast.Ident{
			Name: node.Name,
		}
	case *hast.StringLiteral:
		return &bast.Word{
			Val: fmt.Sprintf(`"%s"`, node.Value),
		}
	case *hast.NumericLiteral:
		return &bast.Word{
			Val: node.Value,
		}
	case *hast.TrueLiteral:
		return &bast.Word{
			Val: "true",
		}
	case *hast.FalseLiteral:
		return &bast.Word{
			Val: "false",
		}
	case *hast.ExprStmt:
		return &bast.ExprStmt{
			X: translate2Bash(opts, node.X).(bast.Expr),
		}
	case *hast.CallExpr:
		fun := translate2Bash(opts, node.Fun).(*bast.Ident)

		recv := []bast.Expr{}
		for _, r := range node.Recv {
			recv = append(recv, translate2Bash(opts, r).(bast.Expr))
		}

		return &bast.CmdExpr{
			Name: fun,
			Recv: recv,
		}
	case *hast.Comment:
		return &bast.Comment{
			Text: node.Text,
		}
	default:
		fmt.Printf("%T\n", node)
	}

	return nil
}

type BashTranspiler struct {
	opts *config.BashOptions

	buffer []bast.Stmt

	moduleManager *ModuleManager

	vfs vfs.VFS

	modules       map[string]*Module
	currentModule *Module

	globalSymbols *SymbolTable

	scopeStack container.Stack[ScopeType]
}

func NewBashTranspiler(opts *config.BashOptions, vfs vfs.VFS) *BashTranspiler {
	return &BashTranspiler{
		opts: opts,
		vfs:  vfs,
		// 初始化模块管理
		moduleManager: &ModuleManager{
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
		globalSymbols: NewSymbolTable(),

		// 初始化作用域管理
		scopeStack: container.NewArrayStack[ScopeType](),
	}
}

func (b *BashTranspiler) Transpile(mainFile string) (map[string]string, error) {
	// 1. 解析所有依赖
	modules, err := b.moduleManager.ResolveAllDependencies(mainFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// 2. 按依赖顺序翻译模块
	results := make(map[string]string)
	for _, module := range modules {
		if err := b.translateModule(module); err != nil {
			return nil, fmt.Errorf("failed to translate module %s: %w", module.Path, err)
		}
		results[module.Path] = bast.String(module.Transpiled)
	}

	return results, nil
}

func (b *BashTranspiler) translateModule(module *Module) error {
	// 设置当前模块
	b.currentModule = module

	// 确保模块有符号表
	if module.Symbols == nil {
		module.Symbols = &SymbolTable{}
	}

	// 分析模块的导出符号
	b.analyzeModuleExports(module)

	// 翻译模块内容
	module.Transpiled = b.Convert(module.AST).(*bast.File)

	// 缓存翻译结果
	// b.translatedModules[module.Path] = module.Transpiled

	return nil
}

func (b *BashTranspiler) analyzeModuleExports(module *Module) {
	for _, stmt := range module.AST.Stmts {
		switch s := stmt.(type) {
		case *hast.FuncDecl:
			// 检查是否有 pub 修饰符
			for _, modifier := range s.Modifiers {
				if _, ok := modifier.(*hast.PubModifier); ok {
					module.Exports[s.Name.Name] = &ExportInfo{
						Name:   s.Name.Name,
						Value:  s,
						Kind:   ExportFunction,
						Public: true,
					}
					module.Symbols.AddFunction(s.Name.Name, s)
					break
				}
			}

			if len(s.Modifiers) == 0 {
				module.Exports[s.Name.Name] = &ExportInfo{
					Name:   s.Name.Name,
					Value:  s,
					Kind:   ExportFunction,
					Public: false,
				}
				module.Symbols.AddFunction(s.Name.Name, s)
			}

		case *hast.ClassDecl:
			// 检查是否有 pub 修饰符
			if s.Pub.IsValid() {
				module.Exports[s.Name.Name] = &ExportInfo{
					Name:   s.Name.Name,
					Value:  s,
					Kind:   ExportClass,
					Public: true,
				}
				module.Symbols.AddClass(s.Name.Name, s)
			}

		case *hast.AssignStmt:
			// pub var a = 10
			// TODO: 要实现 pub 修饰符
			switch s.Scope {
			case htok.CONST:
				if ident, ok := s.Lhs.(*hast.Ident); ok {
					module.Exports[ident.Name] = &ExportInfo{
						Name:   ident.Name,
						Value:  s,
						Kind:   ExportConstant,
						Public: false,
					}
					module.Symbols.AddConstant(ident.Name, s)
				}
			case htok.VAR:
				if ident, ok := s.Lhs.(*hast.Ident); ok {
					module.Exports[ident.Name] = &ExportInfo{
						Name:   ident.Name,
						Value:  s,
						Kind:   ExportVariable,
						Public: false,
					}
					module.Symbols.AddVariable(ident.Name, s)
				}
			}

		case *hast.EnumDecl:
			module.Exports[s.Name.Name] = &ExportInfo{
				Name:   s.Name.Name,
				Value:  s,
				Kind:   ExportConstant,
				Public: false,
			}
			module.Symbols.AddConstant(s.Name.Name, s)
		}
	}
}

func Transpile(opts *config.BashOptions, node hast.Node, vfs vfs.VFS, basePath string, huloPath string, mainFile string, popts ...parser.ParserOptions) (map[string]string, error) {
	tr := NewBashTranspiler(opts, vfs)
	tr.moduleManager.basePath = basePath
	return tr.Transpile(mainFile)
}

func (b *BashTranspiler) Convert(node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		return b.ConvertFile(node)
	case *hast.CommentGroup:
		return b.ConvertCommentGroup(node)
	case *hast.ExprStmt:
		return b.ConvertExprStmt(node)
	case *hast.AssignStmt:
		return b.ConvertAssignStmt(node)
	case *hast.IfStmt:
		return b.ConvertIfStmt(node)
	case *hast.BlockStmt:
		return b.ConvertBlockStmt(node)
	case *hast.BinaryExpr:
		return b.ConvertBinaryExpr(node)
	case *hast.CallExpr:
		return b.ConvertCallExpr(node)
	case *hast.Ident:
		return b.ConvertIdent(node)
	case *hast.StringLiteral:
		return b.ConvertStringLiteral(node)
	case *hast.NumericLiteral:
		return b.ConvertNumericLiteral(node)
	case *hast.TrueLiteral:
		return b.ConvertTrueLiteral(node)
	case *hast.FalseLiteral:
		return b.ConvertFalseLiteral(node)
	case *hast.FuncDecl:
		return b.ConvertFuncDecl(node)
	case *hast.ClassDecl:
		return b.ConvertClassDecl(node)
	case *hast.Import:
		return b.ConvertImport(node)
	case *hast.WhileStmt:
		return b.ConvertWhileStmt(node)
	case *hast.ForStmt:
		return b.ConvertForStmt(node)
	case *hast.ReturnStmt:
		return b.ConvertReturnStmt(node)
	case *hast.Comment:
		return b.ConvertComment(node)
	case *hast.SelectExpr:
		return b.ConvertSelectExpr(node)
	case *hast.RefExpr:
		return b.ConvertRefExpr(node)
	case *hast.Parameter:
		return b.ConvertParameter(node)
	default:
		fmt.Printf("Unhandled node type: %T\n", node)
		return nil
	}
}

func (b *BashTranspiler) ConvertFile(node *hast.File) bast.Node {
	var stmts []bast.Stmt

	// 添加 shebang
	stmts = append(stmts, &bast.Comment{
		Text: "#!/bin/bash",
	})

	// 转换所有语句
	for _, stmt := range node.Stmts {
		converted := b.Convert(stmt)
		if converted != nil {
			if stmtNode, ok := converted.(bast.Stmt); ok {
				stmts = append(stmts, stmtNode)
			}
		}
	}

	// 添加缓冲区中的语句
	stmts = append(stmts, b.Flush()...)

	return &bast.File{
		Stmts: stmts,
	}
}

func (b *BashTranspiler) ConvertCommentGroup(node *hast.CommentGroup) bast.Node {
	var comments []*bast.Comment
	for _, comment := range node.List {
		comments = append(comments, &bast.Comment{
			Text: comment.Text,
		})
	}
	return &bast.CommentGroup{
		List: comments,
	}
}

func (b *BashTranspiler) ConvertComment(node *hast.Comment) bast.Node {
	return &bast.Comment{
		Text: node.Text,
	}
}

func (b *BashTranspiler) ConvertExprStmt(node *hast.ExprStmt) bast.Node {
	expr := b.Convert(node.X)
	if expr == nil {
		return nil
	}

	if stmt, ok := expr.(bast.Stmt); ok {
		return stmt
	}

	return &bast.ExprStmt{
		X: expr.(bast.Expr),
	}
}

func (b *BashTranspiler) ConvertAssignStmt(node *hast.AssignStmt) bast.Node {
	rhs := b.Convert(node.Rhs).(bast.Expr)

	// 处理 $p := "hello" 的情况
	// 在 Hulo 中，$p := "hello" 等价于 let p = "hello"
	if refExpr, ok := node.Lhs.(*hast.RefExpr); ok {
		// 从 $p 中提取变量名 p
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			return &bast.AssignStmt{
				Lhs: &bast.Ident{Name: ident.Name},
				Rhs: rhs,
			}
		}
	}

	// 处理普通赋值 p = "hello"
	lhs := b.Convert(node.Lhs).(bast.Expr)
	return &bast.AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}

func (b *BashTranspiler) ConvertIfStmt(node *hast.IfStmt) bast.Node {
	cond := b.Convert(node.Cond).(bast.Expr)
	body := b.Convert(node.Body).(*bast.BlockStmt)

	var elseStmt bast.Stmt
	if node.Else != nil {
		converted := b.Convert(node.Else)
		if stmt, ok := converted.(bast.Stmt); ok {
			elseStmt = stmt
		}
	}

	return &bast.IfStmt{
		Cond: cond,
		Body: body,
		Else: elseStmt,
	}
}

func (b *BashTranspiler) ConvertBlockStmt(node *hast.BlockStmt) bast.Node {
	var stmts []bast.Stmt
	for _, stmt := range node.List {
		converted := b.Convert(stmt)
		if converted != nil {
			if stmtNode, ok := converted.(bast.Stmt); ok {
				stmts = append(stmts, stmtNode)
			}
		}
	}
	return &bast.BlockStmt{
		List: stmts,
	}
}

func (b *BashTranspiler) ConvertBinaryExpr(node *hast.BinaryExpr) bast.Node {
	x := b.Convert(node.X).(bast.Expr)
	y := b.Convert(node.Y).(bast.Expr)

	// 特殊处理乘法运算，转换为 echo "..." | bc

	// $x * $x
	// $x * 10
	// 10 * 3.14
	// $x * $y
	if node.Op == htok.ASTERISK {
		var lv, rv string
		switch x := x.(type) {
		case *bast.VarExpExpr:
			lv = "$" + x.X.(*bast.Ident).Name
		case *bast.Word:
			lv = x.Val
		}

		switch y := y.(type) {
		case *bast.VarExpExpr:
			rv = "$" + y.X.(*bast.Ident).Name
		case *bast.Word:
			rv = y.Val
		}

		var expr string
		// 生成 echo "x * y" | bc 语句
		if lv == rv {
			expr = fmt.Sprintf("scale=2; sqrt(%s)", lv)
		} else {
			expr = fmt.Sprintf("%s * %s", lv, rv)
		}
		return &bast.CmdSubst{
			Tok: btok.DollParen,
			X: &bast.PipelineExpr{
				CtrOp: btok.Or,
				Cmds: []bast.Expr{
					&bast.CmdExpr{
						Name: &bast.Ident{Name: "echo"},
						Recv: []bast.Expr{
							&bast.Word{Val: fmt.Sprintf(`"%s"`, expr)},
						},
					},
					&bast.CmdExpr{
						Name: &bast.Ident{Name: "bc"},
						Recv: []bast.Expr{},
					},
				},
			},
		}
	}

	return &bast.BinaryExpr{
		X:  x,
		Op: b.convertOperator(node.Op),
		Y:  y,
	}
}

func (b *BashTranspiler) ConvertCallExpr(node *hast.CallExpr) bast.Node {
	fun := b.Convert(node.Fun)
	fnName := ""
	if ident, ok := fun.(*bast.Ident); ok {
		fnName = ident.Name
	} else if varExp, ok := fun.(*bast.VarExpExpr); ok {
		fnName = varExp.X.(*bast.Ident).Name
	}
	// TODO: RefExpr (SelectExpr) 这里错了 Hulo 分析器解析有问题
	// 应该是 SelectExpr 的 X 是 RefExpr
	if refExpr, ok := node.Fun.(*hast.RefExpr); ok {
		fmt.Printf("refExpr: %T\n", refExpr.X)
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			fmt.Printf("ident: %s %T\n", ident.Name, ident)
		}
	}
	fmt.Printf("fnName: %s %T %T\n", fnName, fun, node.Fun)

	var recv []bast.Expr
	for _, arg := range node.Recv {
		recv = append(recv, b.Convert(arg).(bast.Expr))
	}

	return &bast.CmdExpr{
		Name: &bast.Ident{Name: fnName},
		Recv: recv,
	}
}

func (b *BashTranspiler) ConvertIdent(node *hast.Ident) bast.Node {
	return &bast.Ident{
		Name: node.Name,
	}
}

func (b *BashTranspiler) ConvertStringLiteral(node *hast.StringLiteral) bast.Node {
	return &bast.Word{
		Val: fmt.Sprintf(`"%s"`, node.Value),
	}
}

func (b *BashTranspiler) ConvertNumericLiteral(node *hast.NumericLiteral) bast.Node {
	return &bast.Word{
		Val: node.Value,
	}
}

func (b *BashTranspiler) ConvertTrueLiteral(node *hast.TrueLiteral) bast.Node {
	return &bast.Word{
		Val: "true",
	}
}

func (b *BashTranspiler) ConvertFalseLiteral(node *hast.FalseLiteral) bast.Node {
	return &bast.Word{
		Val: "false",
	}
}

func (b *BashTranspiler) ConvertFuncDecl(node *hast.FuncDecl) bast.Node {
	name := b.Convert(node.Name).(*bast.Ident)

	// 处理参数，生成 local 声明
	var paramStmts []bast.Stmt
	var paramNames []string

	for i, param := range node.Recv {
		if paramNode, ok := param.(*hast.Parameter); ok {
			paramName := paramNode.Name.Name
			paramNames = append(paramNames, paramName)

			// 生成 local x=$1 语句
			paramStmts = append(paramStmts, &bast.AssignStmt{
				Local: btok.Pos(1), // 标记为 local
				Lhs:   &bast.Ident{Name: paramName},
				Rhs:   &bast.VarExpExpr{X: &bast.Ident{Name: fmt.Sprintf("%d", i+1)}},
			})
		}
	}

	// 转换函数体
	originalBody := b.Convert(node.Body).(*bast.BlockStmt)

	// 将参数声明添加到函数体开头
	var newBodyStmts []bast.Stmt
	newBodyStmts = append(newBodyStmts, paramStmts...)
	newBodyStmts = append(newBodyStmts, originalBody.List...)

	newBody := &bast.BlockStmt{
		List: newBodyStmts,
	}

	// 添加到当前模块的符号表
	if b.currentModule != nil {
		b.currentModule.Symbols.AddFunction(name.Name, node)
	}

	return &bast.FuncDecl{
		Name: name,
		Body: newBody,
	}
}

func (b *BashTranspiler) ConvertClassDecl(node *hast.ClassDecl) bast.Node {
	name := b.Convert(node.Name).(*bast.Ident)

	// 添加到当前模块的符号表
	if b.currentModule != nil {
		b.currentModule.Symbols.AddClass(name.Name, node)
	}

	// 对于Bash，类通常转换为函数或结构
	// 这里简化为函数声明
	return &bast.FuncDecl{
		Name: name,
		Body: &bast.BlockStmt{},
	}
}

func (b *BashTranspiler) ConvertImport(node *hast.Import) bast.Node {
	// 对于Bash，导入通常转换为source命令
	var path string

	// 处理不同的导入类型
	switch {
	case node.ImportSingle != nil:
		path = node.ImportSingle.Path
	case node.ImportMulti != nil:
		path = node.ImportMulti.Path
	case node.ImportAll != nil:
		path = node.ImportAll.Path
	default:
		// 如果没有具体的导入信息，使用默认路径
		path = "unknown_module"
	}

	// 确保路径有正确的扩展名
	if !strings.HasSuffix(path, ".sh") && !strings.HasSuffix(path, ".bash") {
		path = path + ".sh"
	}

	return &bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{
				Name: "source",
			},
			Recv: []bast.Expr{
				&bast.Word{Val: fmt.Sprintf(`"%s"`, path)},
			},
		},
	}
}

func (b *BashTranspiler) ConvertWhileStmt(node *hast.WhileStmt) bast.Node {
	cond := b.Convert(node.Cond).(bast.Expr)
	body := b.Convert(node.Body).(*bast.BlockStmt)

	return &bast.WhileStmt{
		Cond: cond,
		Body: body,
	}
}

func (b *BashTranspiler) ConvertForStmt(node *hast.ForStmt) bast.Node {
	init := b.Convert(node.Init).(bast.Stmt)
	cond := b.Convert(node.Cond).(bast.Expr)
	post := b.Convert(node.Post).(bast.Stmt)
	body := b.Convert(node.Body).(*bast.BlockStmt)

	return &bast.ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: body,
	}
}

func (b *BashTranspiler) ConvertReturnStmt(node *hast.ReturnStmt) bast.Node {
	var x bast.Expr
	if node.X != nil {
		x = b.Convert(node.X).(bast.Expr)
	} else {
		// 如果没有返回值，使用默认值 0
		x = &bast.Word{Val: "0"}
	}

	return &bast.ReturnStmt{
		X: x,
	}
}

func (b *BashTranspiler) convertOperator(op htok.Token) btok.Token {
	switch op {
	case htok.PLUS:
		return btok.Plus
	case htok.MINUS:
		return btok.Minus
	case htok.ASTERISK:
		return btok.Star
	case htok.SLASH:
		return btok.Slash
	case htok.EQ:
		return btok.Equal
	case htok.NEQ:
		return btok.NotEqual
	case htok.LT:
		return btok.LessEqual
	case htok.LE:
		return btok.LessEqual
	case htok.GT:
		return btok.TsGtr
	case htok.GE:
		return btok.GreatEqual
	default:
		return btok.Illegal
	}
}

func (b *BashTranspiler) needsDeclaration(expr bast.Expr) bool {
	if ident, ok := expr.(*bast.Ident); ok {
		// 检查当前模块的符号表中是否已声明
		if b.currentModule != nil && b.currentModule.Symbols != nil {
			return !b.currentModule.Symbols.HasVariable(ident.Name)
		}
		return true // 如果没有当前模块，默认需要声明
	}
	return false
}

func (b *BashTranspiler) Emit(stmt bast.Stmt) {
	b.buffer = append(b.buffer, stmt)
}

func (v *BashTranspiler) Flush() []bast.Stmt {
	stmts := v.buffer
	v.buffer = nil
	return stmts
}

// ConvertSelectExpr 处理选择表达式，区分模块访问和对象方法调用
func (b *BashTranspiler) ConvertSelectExpr(node *hast.SelectExpr) bast.Node {
	x := b.Convert(node.X).(bast.Expr)
	y := b.Convert(node.Y).(bast.Expr)

	// 检查是否是模块访问
	if b.isModuleAccess(node.X) {
		// 模块访问：math.PI -> 直接使用符号名
		// 模块访问：math.$PI -> 直接使用符号名（去掉$前缀）
		if ident, ok := y.(*bast.Ident); ok {
			return ident
		}
		// 处理 math.$PI 的情况，其中 Y 是 VarExpExpr
		if varExp, ok := y.(*bast.VarExpExpr); ok {
			if ident, ok := varExp.X.(*bast.Ident); ok {
				return ident
			}
		}
	}

	// 对象方法调用：$p.greet() -> 对于 Bash，我们简化为直接调用
	// 因为 Bash 不支持对象方法调用，所以这里需要特殊处理
	_ = x // 暂时忽略左操作数，因为 Bash 不支持对象方法调用
	if ident, ok := y.(*bast.Ident); ok {
		// 如果是方法调用，可能需要特殊处理
		// 这里暂时返回标识符，实际使用时需要根据上下文决定
		return ident
	}
	fmt.Printf("%T %T\n", x, y)
	return y
}

// convertBuiltinMethod 转换内置方法调用
func (b *BashTranspiler) convertBuiltinMethod(obj *hast.RefExpr, methodName string, _ []hast.Expr) bast.Node {
	varName := obj.X.(*hast.Ident).Name

	switch methodName {
	case "length":
		// $p.length() -> length=${#p}
		tempVar := fmt.Sprintf("length_%s", varName)

		// 生成 length_var=${#var} 语句
		assignStmt := &bast.AssignStmt{
			Lhs: &bast.Ident{Name: tempVar},
			Rhs: &bast.ParamExpExpr{
				Dollar: btok.Pos(1),
				Lbrace: btok.Pos(1),
				Var:    &bast.Ident{Name: varName},
				ParamExp: &bast.LengthExp{
					Hash: btok.Pos(1),
				},
				Rbrace: btok.Pos(1),
			},
		}

		// 将赋值语句添加到缓冲区
		b.Emit(assignStmt)

		// 返回临时变量引用
		return &bast.VarExpExpr{
			X: &bast.Ident{Name: tempVar},
		}

	default:
		// 未知方法，返回错误或默认值
		return &bast.Word{Val: "0"}
	}
}

// ConvertRefExpr 处理引用表达式 $x
func (b *BashTranspiler) ConvertRefExpr(node *hast.RefExpr) bast.Node {
	x := b.Convert(node.X).(bast.Expr)

	// 对于 Bash，$x 转换为变量引用
	return &bast.VarExpExpr{
		X: x,
	}
}

// isModuleAccess 判断是否是模块访问
func (b *BashTranspiler) isModuleAccess(expr hast.Expr) bool {
	// 检查是否是标识符（模块名）
	if ident, ok := expr.(*hast.Ident); ok {
		// 检查是否是导入的模块
		if b.currentModule != nil {
			for _, importInfo := range b.currentModule.Imports {
				// 检查是否是导入的模块
				if importInfo.Alias == ident.Name ||
					(importInfo.Kind == ImportSingle && b.getModuleName(importInfo.ModulePath) == ident.Name) {
					return true
				}
			}
		}

		// 检查是否是全局模块（如 std 库）
		if b.isGlobalModule(ident.Name) {
			return true
		}
	}

	return false
}

// isGlobalModule 判断是否是全局模块
func (b *BashTranspiler) isGlobalModule(name string) bool {
	// 常见的全局模块
	globalModules := []string{"std", "math", "io", "net", "time", "fs"}
	return slices.Contains(globalModules, name)
}

// ConvertParameter 处理函数参数
func (b *BashTranspiler) ConvertParameter(node *hast.Parameter) bast.Node {
	// 对于 Bash，我们只需要参数名
	return b.Convert(node.Name).(bast.Expr)
}

// getModuleName 从路径获取模块名
func (b *BashTranspiler) getModuleName(path string) string {
	// 简单的实现，从路径中提取文件名（不含扩展名）
	// 这里可以复用 module.go 中的逻辑
	baseName := path
	if idx := len(baseName) - 1; idx >= 0 && baseName[idx] == '/' {
		baseName = baseName[:idx]
	}
	if idx := len(baseName) - 1; idx >= 0 && baseName[idx] == '\\' {
		baseName = baseName[:idx]
	}
	if idx := len(baseName) - 1; idx >= 0 && baseName[idx] == '.' {
		baseName = baseName[:idx]
	}
	return baseName
}
