// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"
	"path/filepath"
	"strings"

	"slices"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	btok "github.com/hulo-lang/hulo/syntax/bash/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

type BashTranspiler struct {
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

	// 跟踪类实例变量名到类名的映射
	classInstances map[string]string
}

func NewBashTranspiler(opts *config.Huloc, vfs vfs.VFS) *BashTranspiler {
	return &BashTranspiler{
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

		// 初始化类实例映射
		classInstances: make(map[string]string),
	}
}

func (b *BashTranspiler) Transpile(mainFile string) (map[string]string, error) {
	builtinModules, err := b.moduleManager.ResolveAllDependencies(filepath.Join(b.opts.HuloPath, "core/unsafe/bash/index.hl"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve builtin dependencies: %w", err)
	}

	b.builtinModules = builtinModules

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
		path := strings.Replace(module.Path, ".hl", ".sh", 1)
		results[path] = bast.String(module.Transpiled)
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
			// 处理所有类型的赋值语句
			var varName string
			if refExpr, ok := s.Lhs.(*hast.RefExpr); ok {
				// 从 $p 中提取变量名 p
				if ident, ok := refExpr.X.(*hast.Ident); ok {
					varName = ident.Name
				}
			} else if ident, ok := s.Lhs.(*hast.Ident); ok {
				varName = ident.Name
			}

			if varName != "" {
				switch s.Scope {
				case htok.CONST:
					module.Exports[varName] = &ExportInfo{
						Name:   varName,
						Value:  s,
						Kind:   ExportConstant,
						Public: false,
					}
					module.Symbols.AddConstant(varName, s)
				case htok.VAR:
					module.Exports[varName] = &ExportInfo{
						Name:   varName,
						Value:  s,
						Kind:   ExportVariable,
						Public: false,
					}
					module.Symbols.AddVariable(varName, s)
				default:
					// 普通赋值语句（如 let arr = [1, 2, 3]）
					module.Symbols.AddVariable(varName, s)
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

func Transpile(opts *config.Huloc, vfs vfs.VFS, basePath string) (map[string]string, error) {
	tr := NewBashTranspiler(opts, vfs)
	tr.moduleManager.basePath = basePath
	return tr.Transpile(opts.Main)
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
	case *hast.CmdExpr:
		return b.ConvertCmdExpr(node)
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
	case *hast.ForInStmt:
		return b.ConvertForInStmt(node)
	case *hast.ForeachStmt:
		return b.ConvertForeachStmt(node)
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
	case *hast.IncDecExpr:
		return b.ConvertIncDecExpr(node)
	case *hast.ArrayLiteralExpr:
		return b.ConvertArrayLiteralExpr(node)
	case *hast.ObjectLiteralExpr:
		return b.ConvertObjectLiteralExpr(node)
	case *hast.MatchStmt:
		return b.ConvertMatchStmt(node)
	case *hast.DoWhileStmt:
		return b.ConvertDoWhileStmt(node)
	default:
		fmt.Printf("Unhandled node type: %T\n", node)
		return nil
	}
}

func (b *BashTranspiler) ConvertFile(node *hast.File) bast.Node {
	var stmts []bast.Stmt

	// 添加 shebang
	stmts = append(stmts, &bast.Comment{Text: "!/bin/bash"})

	// 转换文档注释
	if node.Docs != nil {
		for _, doc := range node.Docs {
			converted := b.Convert(doc)
			if converted != nil {
				if stmtNode, ok := converted.(bast.Stmt); ok {
					stmts = append(stmts, stmtNode)
				}
			}
		}
	}

	// 转换所有语句

	for _, stmt := range node.Stmts {
		converted := b.Convert(stmt)
		stmts = append(stmts, b.Flush()...)
		if converted != nil {
			if stmtNode, ok := converted.(bast.Stmt); ok {
				stmts = append(stmts, stmtNode)
			}
		}
	}

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
	// 先转换 RHS
	convertedRhs := b.Convert(node.Rhs)
	if convertedRhs == nil {
		return nil
	}

	// 检查是否是内置命令调用
	if cmdExpr, ok := convertedRhs.(*bast.CmdExpr); ok {
		if cmdExpr.Name.Name == "read" {
			return b.convertReadCommand(node.Lhs, cmdExpr)
		}
		// 检查是否是构造函数调用
		if strings.HasPrefix(cmdExpr.Name.Name, "create_") {
			return b.convertConstructorCall(node.Lhs, cmdExpr)
		}
	}

	rhs, ok := convertedRhs.(bast.Expr)
	if !ok {
		return nil
	}

	// 获取变量名
	var varName string
	if refExpr, ok := node.Lhs.(*hast.RefExpr); ok {
		// 从 $p 中提取变量名 p
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			varName = ident.Name
		}
	} else if ident, ok := node.Lhs.(*hast.Ident); ok {
		varName = ident.Name
	}

	// 如果有变量名，进行符号管理
	if varName != "" && b.currentModule != nil {
		currentScope := b.currentModule.Symbols.CurrentScope()
		if currentScope == nil {
			// 如果没有当前作用域，创建一个全局作用域
			currentScope = b.currentModule.Symbols.PushScope(GlobalScope, b.currentModule)
		}

		// 创建符号
		symbol := &Symbol{
			Name:   varName,
			Type:   SymbolVariable,
			Value:  node,
			Scope:  currentScope,
			Module: b.currentModule,
		}

		// 添加到当前作用域
		b.currentModule.Symbols.AddSymbol(symbol)

		// 分配生成名称
		generatedName := b.currentModule.Symbols.AllocVariableName(varName, currentScope)
		symbol.MangledName = generatedName

		// 如果RHS是对象字面量，建立绑定
		if objLit, ok := node.Rhs.(*hast.ObjectLiteralExpr); ok {
			b.bindObjectLiteral(symbol, objLit)
			// 返回赋值语句，使用混淆后的变量名
			return &bast.AssignStmt{
				Lhs: &bast.Ident{Name: generatedName},
				Rhs: &bast.Ident{Name: symbol.MangledName},
			}
		}

		return &bast.AssignStmt{
			Lhs: &bast.Ident{Name: generatedName},
			Rhs: rhs,
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

	// 将条件转换为算术表达式 (( ... ))
	arithCond := b.convertToArithExpr(cond)

	var elseStmt bast.Stmt
	if node.Else != nil {
		converted := b.Convert(node.Else)
		if stmt, ok := converted.(bast.Stmt); ok {
			elseStmt = stmt
		}
	}

	return &bast.IfStmt{
		If:   btok.Pos(1),
		Cond: arithCond,
		Semi: btok.Pos(1),
		Then: btok.Pos(1),
		Body: body,
		Else: elseStmt,
		Fi:   btok.Pos(1),
	}
}

func (b *BashTranspiler) ConvertBlockStmt(node *hast.BlockStmt) bast.Node {
	// 进入新作用域
	if b.currentModule != nil {
		b.currentModule.Symbols.PushScope(BlockScope, b.currentModule)
		defer b.currentModule.Symbols.PopScope()
	}

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

	// 处理其他算术运算，转换为 $(( ... )) 语法
	if node.Op == htok.PLUS || node.Op == htok.MINUS || node.Op == htok.SLASH {
		var lv, rv string
		switch x := x.(type) {
		case *bast.VarExpExpr:
			lv = x.X.(*bast.Ident).Name
		case *bast.Word:
			lv = x.Val
		}

		switch y := y.(type) {
		case *bast.VarExpExpr:
			rv = y.X.(*bast.Ident).Name
		case *bast.Word:
			rv = y.Val
		}

		// 生成 $((a + b)) 语法
		return &bast.ArithExpr{
			Dollar: btok.Pos(1),
			Lparen: btok.Pos(1),
			X: &bast.BinaryExpr{
				X:  &bast.Ident{Name: lv},
				Op: b.convertOperator(node.Op),
				Y:  &bast.Ident{Name: rv},
			},
			Rparen: btok.Pos(1),
		}
	}

	return &bast.BinaryExpr{
		X:  x,
		Op: b.convertOperator(node.Op),
		Y:  y,
	}
}

func (b *BashTranspiler) ConvertCallExpr(node *hast.CallExpr) bast.Node {
	// 检查是否是类的方法调用：$u.to_str()
	if refExpr, ok := node.Fun.(*hast.RefExpr); ok {
		if selectExpr, ok := refExpr.X.(*hast.SelectExpr); ok {
			if ident, ok := selectExpr.X.(*hast.Ident); ok {
				// 这是一个对象方法调用
				objectVar := ident.Name
				methodName := selectExpr.Y.(*hast.Ident).Name

				// 从映射中获取类名
				className := b.classInstances[objectVar]
				if className == "" {
					// 如果没有找到映射，使用默认处理
					className = strings.Title(objectVar)
				}
				classNameLower := strings.ToLower(className)

				// 生成方法调用：user_to_str "$u"
				bashMethodName := fmt.Sprintf("%s_%s", classNameLower, methodName)

				var args []bast.Expr
				// 第一个参数是对象本身
				args = append(args, &bast.VarExpExpr{X: &bast.Ident{Name: objectVar}})
				// 添加其他参数
				for _, arg := range node.Recv {
					args = append(args, b.Convert(arg).(bast.Expr))
				}

				return &bast.CmdExpr{
					Name: &bast.Ident{Name: bashMethodName},
					Recv: args,
				}
			}
		}
	}

	// 检查是否是构造函数调用：User("John", 20)
	if ident, ok := node.Fun.(*hast.Ident); ok {
		className := ident.Name
		// 检查是否是已知的类
		if b.currentModule != nil && b.currentModule.Symbols.HasClass(className) {
			// 这是构造函数调用
			constructorName := fmt.Sprintf("create_%s", strings.ToLower(className))

			var args []bast.Expr
			for _, arg := range node.Recv {
				args = append(args, b.Convert(arg).(bast.Expr))
			}

			return &bast.CmdExpr{
				Name: &bast.Ident{Name: constructorName},
				Recv: args,
			}
		}
	}

	// 普通函数调用
	fun := b.Convert(node.Fun)
	fnName := ""
	if ident, ok := fun.(*bast.Ident); ok {
		fnName = ident.Name
	} else if varExp, ok := fun.(*bast.VarExpExpr); ok {
		fnName = varExp.X.(*bast.Ident).Name
	}

	var recv []bast.Expr
	for _, arg := range node.Recv {
		recv = append(recv, b.Convert(arg).(bast.Expr))
	}

	return &bast.CmdExpr{
		Name: &bast.Ident{Name: fnName},
		Recv: recv,
	}
}

func (b *BashTranspiler) ConvertCmdExpr(node *hast.CmdExpr) bast.Node {
	// 检查是否是类的方法调用：$u.to_str()
	if refExpr, ok := node.Cmd.(*hast.RefExpr); ok {
		if selectExpr, ok := refExpr.X.(*hast.SelectExpr); ok {
			if ident, ok := selectExpr.X.(*hast.Ident); ok {
				// 这是一个对象方法调用
				objectVar := ident.Name
				methodName := selectExpr.Y.(*hast.Ident).Name

				// 从映射中获取类名
				className := b.classInstances[objectVar]
				if className == "" {
					// 如果没有找到映射，使用默认处理
					className = strings.Title(objectVar)
				}
				classNameLower := strings.ToLower(className)

				// 生成方法调用：user_to_str "$u"
				bashMethodName := fmt.Sprintf("%s_%s", classNameLower, methodName)

				var args []bast.Expr
				// 第一个参数是对象本身
				args = append(args, &bast.VarExpExpr{X: &bast.Ident{Name: objectVar}})
				// 添加其他参数
				for _, arg := range node.Args {
					args = append(args, b.Convert(arg).(bast.Expr))
				}

				return &bast.CmdExpr{
					Name: &bast.Ident{Name: bashMethodName},
					Recv: args,
				}
			}
		}
	}

	// 检查是否是构造函数调用：User("John", 20)
	if ident, ok := node.Cmd.(*hast.Ident); ok {
		className := ident.Name
		// 检查是否是已知的类
		if b.currentModule != nil && b.currentModule.Symbols.HasClass(className) {
			// 这是构造函数调用
			constructorName := fmt.Sprintf("create_%s", strings.ToLower(className))

			var args []bast.Expr
			for _, arg := range node.Args {
				args = append(args, b.Convert(arg).(bast.Expr))
			}

			return &bast.CmdExpr{
				Name: &bast.Ident{Name: constructorName},
				Recv: args,
			}
		}
	}

	// 普通函数调用
	fun := b.Convert(node.Cmd)
	fnName := ""
	if ident, ok := fun.(*bast.Ident); ok {
		fnName = ident.Name
	} else if varExp, ok := fun.(*bast.VarExpExpr); ok {
		fnName = varExp.X.(*bast.Ident).Name
	}

	var recv []bast.Expr
	for _, arg := range node.Args {
		recv = append(recv, b.Convert(arg).(bast.Expr))
	}

	return &bast.CmdExpr{
		Name: &bast.Ident{Name: fnName},
		Recv: recv,
	}
}

func (b *BashTranspiler) ConvertIdent(node *hast.Ident) bast.Node {
	// 查符号表，返回混淆名
	if b.currentModule != nil && b.currentModule.Symbols != nil {
		if symbol := b.currentModule.Symbols.LookupSymbol(node.Name); symbol != nil && symbol.MangledName != "" {
			return &bast.Ident{Name: symbol.MangledName}
		}
	}
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
	className := node.Name.Name

	// 添加到当前模块的符号表
	if b.currentModule != nil {
		b.currentModule.Symbols.AddClass(className, node)
	}

	// 生成构造函数
	constructorName := fmt.Sprintf("create_%s", strings.ToLower(className))

	// 收集字段名
	var fieldNames []string
	for _, field := range node.Fields.List {
		fieldNames = append(fieldNames, field.Name.Name)
	}

	// 生成构造函数体
	var constructorBody []bast.Stmt

	// 添加 local 参数声明
	for i, fieldName := range fieldNames {
		constructorBody = append(constructorBody, &bast.AssignStmt{
			Local: btok.Pos(1),
			Lhs:   &bast.Ident{Name: fieldName},
			Rhs:   &bast.VarExpExpr{X: &bast.Ident{Name: fmt.Sprintf("%d", i+1)}},
		})
	}

	// 创建关联数组声明
	constructorBody = append(constructorBody, &bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{Name: "declare"},
			Recv: []bast.Expr{
				&bast.Word{Val: "-A"},
				&bast.Ident{Name: strings.ToLower(className)},
			},
		},
	})

	// 添加字段赋值
	for _, fieldName := range fieldNames {
		constructorBody = append(constructorBody, &bast.AssignStmt{
			Lhs: &bast.IndexExpr{
				X:      &bast.Ident{Name: strings.ToLower(className)},
				Lbrack: btok.Pos(1),
				Y:      &bast.Word{Val: fmt.Sprintf(`"%s"`, fieldName)},
				Rbrack: btok.Pos(1),
			},
			Rhs: &bast.VarExpExpr{X: &bast.Ident{Name: fieldName}},
		})
	}

	// 返回关联数组
	constructorBody = append(constructorBody, &bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{Name: "echo"},
			Recv: []bast.Expr{
				&bast.Word{Val: fmt.Sprintf(`"$(declare -p %s)"`, strings.ToLower(className))},
			},
		},
	})

	// 生成构造函数
	constructor := &bast.FuncDecl{
		Name: &bast.Ident{Name: constructorName},
		Body: &bast.BlockStmt{List: constructorBody},
	}

	// 将构造函数添加到缓冲区
	b.Emit(constructor)

	// 生成方法
	for _, method := range node.Methods {
		b.generateClassMethod(className, method)
	}

	// 返回一个空的语句（因为构造函数和方法已经通过 Emit 添加）
	return &bast.ExprStmt{
		X: &bast.Word{Val: ""},
	}
}

func (b *BashTranspiler) generateClassMethod(className string, method *hast.FuncDecl) {
	methodName := method.Name.Name
	classNameLower := strings.ToLower(className)

	// 生成方法名：user_to_str, user_greet 等
	bashMethodName := fmt.Sprintf("%s_%s", classNameLower, methodName)

	// 生成方法体
	var methodBody []bast.Stmt

	// 添加 eval 语句来解析传入的对象
	methodBody = append(methodBody, &bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{Name: "eval"},
			Recv: []bast.Expr{
				&bast.Word{Val: fmt.Sprintf(`"declare -A %s=${1}"`, classNameLower)},
			},
		},
	})

	// 处理参数（跳过第一个参数，因为它是对象本身）
	for i, param := range method.Recv {
		if paramNode, ok := param.(*hast.Parameter); ok {
			paramName := paramNode.Name.Name
			// 参数从 $2 开始，因为 $1 是对象
			methodBody = append(methodBody, &bast.AssignStmt{
				Local: btok.Pos(1),
				Lhs:   &bast.Ident{Name: paramName},
				Rhs:   &bast.VarExpExpr{X: &bast.Ident{Name: fmt.Sprintf("%d", i+2)}},
			})
		}
	}

	// 转换方法体
	if method.Body != nil {
		convertedBody := b.Convert(method.Body)
		if block, ok := convertedBody.(*bast.BlockStmt); ok {
			// 替换方法体中的字段访问
			modifiedBody := b.replaceFieldAccesses(block, classNameLower)
			methodBody = append(methodBody, modifiedBody.List...)
		}
	}

	// 生成方法函数
	methodFunc := &bast.FuncDecl{
		Name: &bast.Ident{Name: bashMethodName},
		Body: &bast.BlockStmt{List: methodBody},
	}

	// 将方法添加到缓冲区
	b.Emit(methodFunc)
}

func (b *BashTranspiler) replaceFieldAccesses(block *bast.BlockStmt, classNameLower string) *bast.BlockStmt {
	// 这里需要递归替换所有字段访问
	// 例如：$name 应该替换为 ${user["name"]}
	// 这是一个简化的实现，实际需要更复杂的AST遍历
	return block
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

	// 将条件转换为 Bash 的 test 表达式
	// 例如: $a < 2 转换为 [ "$a" -lt 2 ]
	bashCond := b.convertToBashTest(cond)

	return &bast.WhileStmt{
		While: btok.Pos(1),
		Cond:  bashCond,
		Semi:  btok.Pos(1),
		Do:    btok.Pos(1),
		Body:  body,
		Done:  btok.Pos(1),
	}
}

func (b *BashTranspiler) ConvertForStmt(node *hast.ForStmt) bast.Node {
	// C-style for 循环: loop (init; cond; post) { body }
	// 转换循环体
	body := b.Convert(node.Body).(*bast.BlockStmt)

	// 转换初始化语句
	var init bast.Node
	if node.Init != nil {
		init = b.Convert(node.Init)
	}

	// 转换条件表达式
	var cond bast.Expr
	if node.Cond != nil {
		cond = b.Convert(node.Cond).(bast.Expr)
	}

	// 转换后置表达式
	var post bast.Node
	if node.Post != nil {
		// 特殊处理 IncDecExpr (如 i++)
		if incDec, ok := node.Post.(*hast.IncDecExpr); ok {
			// 将 i++ 转换为 i=i+1
			var tok btok.Token = btok.Plus
			if incDec.Tok == htok.DEC {
				tok = btok.Minus
			}
			post = &bast.AssignStmt{
				Lhs: b.Convert(incDec.X).(bast.Expr),
				Rhs: &bast.BinaryExpr{
					X:  b.Convert(incDec.X).(bast.Expr),
					Op: tok,
					Y:  &bast.Word{Val: "1"},
				},
			}
		} else {
			post = b.Convert(node.Post)
		}
	}

	// ====== 变量名去 $ 处理 ======
	stripVarExp := func(expr bast.Expr) bast.Expr {
		if v, ok := expr.(*bast.VarExpExpr); ok {
			if id, ok := v.X.(*bast.Ident); ok {
				return id
			}
		}
		return expr
	}

	// 处理 init
	if assign, ok := init.(*bast.AssignStmt); ok {
		assign.Lhs = stripVarExp(assign.Lhs)
		init = assign
	}
	// 处理 cond
	if bexpr, ok := cond.(*bast.BinaryExpr); ok {
		bexpr.X = stripVarExp(bexpr.X)
		bexpr.Y = stripVarExp(bexpr.Y)
		cond = bexpr
	}
	// 处理 post
	if assign, ok := post.(*bast.AssignStmt); ok {
		assign.Lhs = stripVarExp(assign.Lhs)
		if bexpr, ok := assign.Rhs.(*bast.BinaryExpr); ok {
			bexpr.X = stripVarExp(bexpr.X)
			bexpr.Y = stripVarExp(bexpr.Y)
			assign.Rhs = bexpr
		}
		post = assign
	}

	// 创建 Bash C-style for 循环
	// for ((i=0; i<10; i++)); do ... done
	return &bast.ForStmt{
		For:    btok.Pos(1),
		Lparen: btok.Pos(1),
		Init:   init,
		Semi1:  btok.Pos(1),
		Cond:   cond,
		Semi2:  btok.Pos(1),
		Post:   post,
		Rparen: btok.Pos(1),
		Do:     btok.Pos(1),
		Body:   body,
		Done:   btok.Pos(1),
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
		return btok.TsLss
	case htok.LE:
		return btok.TsLeq
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

func (b *BashTranspiler) ConvertIncDecExpr(node *hast.IncDecExpr) bast.Node {
	// 将 i++ 转换为 a=$((a + 1))，将 i-- 转换为 a=$((a - 1))
	var op btok.Token
	if node.Tok == htok.INC {
		op = btok.Plus
	} else {
		op = btok.Minus
	}

	// 获取变量名（去掉 $ 前缀）
	var varName string
	if refExpr, ok := node.X.(*hast.RefExpr); ok {
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			varName = ident.Name
		}
	}

	// 创建算术表达式 $((a + 1))
	arithExpr := &bast.ArithExpr{
		Dollar: btok.Pos(1),
		Lparen: btok.Pos(1),
		X: &bast.BinaryExpr{
			X:  &bast.Ident{Name: varName},
			Op: op,
			Y:  &bast.Word{Val: "1"},
		},
		Rparen: btok.Pos(1),
	}

	// 创建赋值语句 a=$((a + 1))
	return &bast.AssignStmt{
		Lhs: &bast.Ident{Name: varName},
		Rhs: arithExpr,
	}
}

// convertToBashTest 将 Hulo 的条件表达式转换为 Bash 的 test 表达式
// 例如: $a < 2 转换为 [ "$a" -lt 2 ]
func (b *BashTranspiler) convertToBashTest(expr bast.Expr) bast.Expr {
	if bexpr, ok := expr.(*bast.BinaryExpr); ok {
		// 获取左右操作数
		left := bexpr.X
		right := bexpr.Y

		// 确保变量引用被正确引用
		if v, ok := left.(*bast.VarExpExpr); ok {
			left = &bast.Word{Val: fmt.Sprintf(`"$%s"`, v.X.(*bast.Ident).Name)}
		}
		if v, ok := right.(*bast.VarExpExpr); ok {
			right = &bast.Word{Val: fmt.Sprintf(`"$%s"`, v.X.(*bast.Ident).Name)}
		}

		// 转换操作符
		var testOp btok.Token
		switch bexpr.Op {
		case btok.TsLeq:
			testOp = btok.TsLeq
		case btok.TsGtr:
			testOp = btok.TsGtr
		case btok.TsGeq:
			testOp = btok.TsGeq
		case btok.TsEql:
			testOp = btok.TsEql
		case btok.TsNeq:
			testOp = btok.TsNeq
		case btok.TsLss:
			testOp = btok.TsLss
		default:
			testOp = btok.TsLss // 默认使用 -lt
		}

		// 创建 test 表达式 [ left -op right ]
		return &bast.TestExpr{
			Lbrack: btok.Pos(1),
			X: &bast.BinaryExpr{
				X:  left,
				Op: testOp,
				Y:  right,
			},
			Rbrack: btok.Pos(1),
		}
	}

	return expr
}

// convertToArithExpr 将 Hulo 的条件表达式转换为 Bash 的算术表达式
// 例如: $a < 2 转换为 (( a < 2 ))
func (b *BashTranspiler) convertToArithExpr(expr bast.Expr) bast.Expr {
	if bexpr, ok := expr.(*bast.BinaryExpr); ok {
		// 获取左右操作数
		left := bexpr.X
		right := bexpr.Y

		// 确保变量引用被正确引用（去掉 $ 前缀）
		// 注意：这里不需要修改，因为VarExpExpr中的X已经是混淆后的名称
		if v, ok := left.(*bast.VarExpExpr); ok {
			left = &bast.Ident{Name: v.X.(*bast.Ident).Name}
		}
		if v, ok := right.(*bast.VarExpExpr); ok {
			right = &bast.Ident{Name: v.X.(*bast.Ident).Name}
		}

		// 转换操作符
		var arithOp btok.Token
		switch bexpr.Op {
		case btok.TsLeq:
			arithOp = btok.LessEqual
		case btok.TsGtr:
			arithOp = btok.RdrOut
		case btok.TsGeq:
			arithOp = btok.GreatEqual
		case btok.TsEql:
			arithOp = btok.Equal
		case btok.TsNeq:
			arithOp = btok.NotEqual
		case btok.TsLss:
			arithOp = btok.RdrIn // 用 -lt 代表 <，在 Bash 算术表达式里就是 <
		default:
			arithOp = btok.RdrIn // 默认用 <
		}

		// 创建算术表达式 (( left op right ))
		return &bast.ArithEvalExpr{
			Lparen: btok.Pos(1),
			X: &bast.BinaryExpr{
				X:  left,
				Op: arithOp,
				Y:  right,
			},
			Rparen: btok.Pos(1),
		}
	}

	return expr
}

// convertReadCommand 将 read("prompt") 转换为 read -p "prompt" variable
func (b *BashTranspiler) convertReadCommand(lhs hast.Expr, callExpr *bast.CmdExpr) bast.Node {
	// 获取变量名
	var varName string
	if refExpr, ok := lhs.(*hast.RefExpr); ok {
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			varName = ident.Name
		}
	} else if ident, ok := lhs.(*hast.Ident); ok {
		varName = ident.Name
	}

	// 获取提示信息
	var prompt string
	if len(callExpr.Recv) > 0 {
		if strLit, ok := callExpr.Recv[0].(*bast.Word); ok {
			prompt = strLit.Val
		}
	}

	// 创建 read -p "prompt" variable 命令
	return &bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{Name: "read"},
			Recv: []bast.Expr{
				&bast.Word{Val: "-p"},
				&bast.Word{Val: prompt},
				&bast.Ident{Name: varName},
			},
		},
	}
}

// convertConstructorCall 将构造函数调用转换为变量赋值
func (b *BashTranspiler) convertConstructorCall(lhs hast.Expr, callExpr *bast.CmdExpr) bast.Node {
	// 获取变量名
	var varName string
	if refExpr, ok := lhs.(*hast.RefExpr); ok {
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			varName = ident.Name
		}
	} else if ident, ok := lhs.(*hast.Ident); ok {
		varName = ident.Name
	}

	// 从构造函数名推断类名
	className := ""
	if strings.HasPrefix(callExpr.Name.Name, "create_") {
		className = strings.Title(strings.TrimPrefix(callExpr.Name.Name, "create_"))
	}

	// 记录变量名到类名的映射
	if className != "" {
		b.classInstances[varName] = className
	}

	// 创建变量赋值：u=$(create_user "John" 20)
	return &bast.AssignStmt{
		Lhs: &bast.Ident{Name: varName},
		Rhs: &bast.CmdSubst{
			Tok: btok.DollParen,
			X:   callExpr,
		},
	}
}

// ConvertArrayLiteralExpr 转换数组字面量表达式
func (b *BashTranspiler) ConvertArrayLiteralExpr(node *hast.ArrayLiteralExpr) bast.Node {
	var elements []bast.Expr
	for _, elem := range node.Elems {
		converted := b.Convert(elem)
		if converted != nil {
			if expr, ok := converted.(bast.Expr); ok {
				elements = append(elements, expr)
			}
		}
	}

	// 在 Bash 中，数组用括号表示：(item1 item2 item3)
	return &bast.ArrExpr{
		Lparen: btok.Pos(1),
		Vars:   elements,
		Rparen: btok.Pos(1),
	}
}

// ConvertObjectLiteralExpr 转换对象字面量表达式
func (b *BashTranspiler) ConvertObjectLiteralExpr(node *hast.ObjectLiteralExpr) bast.Node {
	return &bast.Ident{Name: ""}
	// 生成唯一的变量名
	arrayVarName := fmt.Sprintf("obj_%d", len(b.buffer))

	// 添加 declare -A 声明
	b.Emit(&bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{Name: "declare"},
			Recv: []bast.Expr{
				&bast.Word{Val: "-A"},
				&bast.Ident{Name: arrayVarName},
			},
		},
	})

	// 添加键值对赋值
	for _, prop := range node.Props {
		if kv, ok := prop.(*hast.KeyValueExpr); ok {
			key := b.Convert(kv.Key).(bast.Expr)
			value := b.Convert(kv.Value).(bast.Expr)

			b.Emit(&bast.AssignStmt{
				Lhs: &bast.IndexExpr{
					X:      &bast.Ident{Name: arrayVarName},
					Lbrack: btok.Pos(1),
					Y:      key,
					Rbrack: btok.Pos(1),
				},
				Rhs: value,
			})
		}
	}

	// 返回变量引用
	return &bast.Ident{Name: arrayVarName}
}

// ConvertForInStmt 转换 for-in 循环
func (b *BashTranspiler) ConvertForInStmt(node *hast.ForInStmt) bast.Node {
	// 转换循环变量 - Index 是 *Ident 类型
	loopVarName := node.Index.Name

	// 转换范围表达式
	rangeExpr := b.Convert(&node.RangeExpr).(bast.Expr)

	// 转换循环体
	body := b.Convert(node.Body).(*bast.BlockStmt)

	// 在 Bash 中，for-in 循环格式：for var in list; do ... done
	return &bast.ForInStmt{
		For:  btok.Pos(1),
		Var:  &bast.Ident{Name: loopVarName}, // 使用变量名，不是 $var
		In:   btok.Pos(1),
		List: rangeExpr,
		Semi: btok.Pos(1),
		Do:   btok.Pos(1),
		Body: body,
		Done: btok.Pos(1),
	}
}

// ConvertForeachStmt 转换 foreach 循环
func (b *BashTranspiler) ConvertForeachStmt(node *hast.ForeachStmt) bast.Node {
	if node.Tok == htok.OF {
		return b.ConvertForOfStmt(node)
	}

	// 转换循环变量 - 处理解构赋值和单个变量
	var loopVarName string

	// 检查是否是解构赋值，如 ($key, $value)
	// 暂时简化处理，只处理单个变量
	if refExpr, ok := node.Index.(*hast.RefExpr); ok {
		// 单个变量，如 $item
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			loopVarName = ident.Name
		}
	}
	// 如果没有找到变量名，使用默认值
	if loopVarName == "" {
		loopVarName = "item"
	}

	// 转换要遍历的变量/数组
	var listExpr bast.Expr
	if refExpr, ok := node.Var.(*hast.RefExpr); ok {
		// 如果是变量引用，如 $arr
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			listExpr = &bast.VarExpExpr{X: &bast.Ident{Name: ident.Name}}
		}
	} else if arrayLit, ok := node.Var.(*hast.ArrayLiteralExpr); ok {
		// 如果是数组字面量，如 [0, 1, 2]
		listExpr = b.Convert(arrayLit).(bast.Expr)
	}
	// 如果没有找到列表表达式，使用默认值
	if listExpr == nil {
		listExpr = &bast.Word{Val: "$arr"}
	}

	// 转换循环体
	body := b.Convert(node.Body).(*bast.BlockStmt)

	// 在 Bash 中，foreach 循环格式：for var in list; do ... done
	return &bast.ForInStmt{
		For:  btok.Pos(1),
		Var:  &bast.Ident{Name: loopVarName}, // 使用变量名，不是 $var
		In:   btok.Pos(1),
		List: listExpr,
		Semi: btok.Pos(1),
		Do:   btok.Pos(1),
		Body: body,
		Done: btok.Pos(1),
	}
}

func (b *BashTranspiler) ConvertForOfStmt(node *hast.ForeachStmt) bast.Node {
	// 获取要遍历的关联数组变量名
	var arrayVarName string
	if refExpr, ok := node.Var.(*hast.RefExpr); ok {
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			// 查找符号表中的混淆名称
			if b.currentModule != nil {
				if symbol := b.currentModule.Symbols.LookupSymbol(ident.Name); symbol != nil {
					arrayVarName = symbol.MangledName
				} else {
					arrayVarName = ident.Name
				}
			} else {
				arrayVarName = ident.Name
			}
		}
	}
	if arrayVarName == "" {
		arrayVarName = "config"
	}

	// 如果变量是对象字面量，直接转换它
	if objLit, ok := node.Var.(*hast.ObjectLiteralExpr); ok {
		// 直接转换对象字面量，它会返回变量引用
		arrayVarExpr := b.Convert(objLit).(*bast.Ident)
		arrayVarName = arrayVarExpr.Name
	}

	// 分析解构赋值模式
	var keyVarName, valueVarName string
	var isKeyOnly, isValueOnly bool

	// 检查解构赋值模式
	// 暂时简化处理，只处理单个变量
	if refExpr, ok := node.Index.(*hast.RefExpr); ok {
		// 单个变量，如 $key
		if ident, ok := refExpr.X.(*hast.Ident); ok {
			keyVarName = ident.Name
			isKeyOnly = true
		}
	}

	// 转换循环体
	body := b.Convert(node.Body).(*bast.BlockStmt)

	// 根据解构模式生成不同的 Bash 循环
	if keyVarName != "" && valueVarName != "" {
		// 模式：($key, $value) - 遍历 key 和 value
		return &bast.ForInStmt{
			For:  btok.Pos(1),
			Var:  &bast.Ident{Name: keyVarName},
			In:   btok.Pos(1),
			List: &bast.Word{Val: fmt.Sprintf(`"${!%s[@]}"`, arrayVarName)},
			Semi: btok.Pos(1),
			Do:   btok.Pos(1),
			Body: &bast.BlockStmt{
				List: []bast.Stmt{
					&bast.AssignStmt{
						Lhs: &bast.Ident{Name: valueVarName},
						Rhs: &bast.Word{Val: fmt.Sprintf(`"${%s[$%s]}"`, arrayVarName, keyVarName)},
					},
					body.List[0], // 原始的 echo 语句
				},
			},
			Done: btok.Pos(1),
		}
	} else if isKeyOnly {
		// 模式：($key, _) 或 $key - 只遍历 key
		return &bast.ForInStmt{
			For:  btok.Pos(1),
			Var:  &bast.Ident{Name: keyVarName},
			In:   btok.Pos(1),
			List: &bast.Word{Val: fmt.Sprintf(`"${!%s[@]}"`, arrayVarName)},
			Semi: btok.Pos(1),
			Do:   btok.Pos(1),
			Body: body,
			Done: btok.Pos(1),
		}
	} else if isValueOnly {
		// 模式：(_, $value) - 只遍历 value
		return &bast.ForInStmt{
			For:  btok.Pos(1),
			Var:  &bast.Ident{Name: valueVarName},
			In:   btok.Pos(1),
			List: &bast.Word{Val: fmt.Sprintf(`"${%s[@]}"`, arrayVarName)},
			Semi: btok.Pos(1),
			Do:   btok.Pos(1),
			Body: body,
			Done: btok.Pos(1),
		}
	}

	// 默认情况
	return &bast.ForInStmt{
		For:  btok.Pos(1),
		Var:  &bast.Ident{Name: "key"},
		In:   btok.Pos(1),
		List: &bast.Word{Val: fmt.Sprintf(`"${!%s[@]}"`, arrayVarName)},
		Semi: btok.Pos(1),
		Do:   btok.Pos(1),
		Body: body,
		Done: btok.Pos(1),
	}
}

// bindObjectLiteral 绑定对象字面量到符号
func (b *BashTranspiler) bindObjectLiteral(symbol *Symbol, objLit *hast.ObjectLiteralExpr) {
	// 使用符号的生成名称
	arrayVarName := symbol.MangledName

	// 生成关联数组声明
	b.Emit(&bast.ExprStmt{
		X: &bast.CmdExpr{
			Name: &bast.Ident{Name: "declare"},
			Recv: []bast.Expr{
				&bast.Word{Val: "-A"},
				&bast.Ident{Name: arrayVarName},
			},
		},
	})

	// 添加键值对赋值
	for _, prop := range objLit.Props {
		if kv, ok := prop.(*hast.KeyValueExpr); ok {
			key := b.Convert(kv.Key).(bast.Expr)
			value := b.Convert(kv.Value).(bast.Expr)

			b.Emit(&bast.AssignStmt{
				Lhs: &bast.IndexExpr{
					X:      &bast.Ident{Name: arrayVarName},
					Lbrack: btok.Pos(1),
					Y:      key,
					Rbrack: btok.Pos(1),
				},
				Rhs: value,
			})
		}
	}
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

// ConvertMatchStmt 将 match 语句转换为 case 语句
func (b *BashTranspiler) ConvertMatchStmt(node *hast.MatchStmt) bast.Node {
	// 转换匹配表达式
	matchExpr := b.Convert(node.Expr).(bast.Expr)

	// 收集所有 case 子句
	var patterns []*bast.CaseClause

	// 处理所有 case 子句
	for _, caseClause := range node.Cases {
		// 转换 case 条件
		cond := b.Convert(caseClause.Cond).(bast.Expr)

		// 转换 case 体
		var caseBody *bast.BlockStmt
		if caseClause.Body != nil {
			converted := b.Convert(caseClause.Body)
			if block, ok := converted.(*bast.BlockStmt); ok {
				caseBody = block
			} else {
				caseBody = &bast.BlockStmt{
					List: []bast.Stmt{converted.(bast.Stmt)},
				}
			}
		} else {
			caseBody = &bast.BlockStmt{List: []bast.Stmt{}}
		}

		// 创建 case 子句
		pattern := &bast.CaseClause{
			Conds: []bast.Expr{cond},
			Body:  caseBody,
			Semi:  btok.Pos(1), // ;;
		}

		patterns = append(patterns, pattern)
	}

	// 处理默认分支
	var elseBody *bast.BlockStmt
	if node.Default != nil {
		if node.Default.Body != nil {
			converted := b.Convert(node.Default.Body)
			if block, ok := converted.(*bast.BlockStmt); ok {
				elseBody = block
			} else {
				elseBody = &bast.BlockStmt{
					List: []bast.Stmt{converted.(bast.Stmt)},
				}
			}
		} else {
			elseBody = &bast.BlockStmt{List: []bast.Stmt{}}
		}
	}

	// 创建 case 语句
	return &bast.CaseStmt{
		Case:     btok.Pos(1),
		X:        matchExpr,
		In:       btok.Pos(1),
		Patterns: patterns,
		Else:     elseBody,
		Esac:     btok.Pos(1),
	}
}

// ConvertDoWhileStmt 将 do-while 语句转换为 while 循环
func (b *BashTranspiler) ConvertDoWhileStmt(node *hast.DoWhileStmt) bast.Node {
	// 转换循环体
	body := b.Convert(node.Body).(*bast.BlockStmt)

	// 转换条件表达式
	cond := b.Convert(node.Cond).(bast.Expr)

	// 在 Bash 中，do-while 循环可以通过以下方式实现：
	// 1. 先执行一次循环体
	// 2. 然后使用 while 循环，条件是原来的条件
	// 3. 在 while 循环中再次执行循环体

	// 创建 while 循环
	whileStmt := &bast.WhileStmt{
		While: btok.Pos(1),
		Cond:  cond,
		Semi:  btok.Pos(1),
		Do:    btok.Pos(1),
		Body:  body,
		Done:  btok.Pos(1),
	}

	// 返回 while 语句（因为 do-while 在 Bash 中就是 while 循环）
	return whileStmt
}
