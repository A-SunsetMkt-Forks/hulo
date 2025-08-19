// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vbs

import (
	"fmt"
	"strings"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/internal/module"
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	vtok "github.com/hulo-lang/hulo/syntax/vbs/token"
)

func Transpile(opts *config.Huloc, fs vfs.VFS, basePath string) (map[string]string, error) {
	moduleMgr := module.NewDependecyResolver(opts, fs)
	err := module.ResolveAllDependencies(moduleMgr, opts.Main)
	if err != nil {
		return nil, err
	}
	transpiler := NewTranspiler(opts.CompilerOptions.VBScript, moduleMgr).RegisterDefaultRules()
	err = transpiler.BindRules()
	if err != nil {
		return nil, err
	}

	results := make(map[string]vast.Node)
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
	ld.Listen(".sh", linker.BeginEnd{Begin: "# HULO_LINK_BEGIN", End: "# HULO_LINK_END"})

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
		ret[path] = vast.String(node)
	}
	return ret, nil
}

type Transpiler struct {
	opts *config.VBScriptOptions

	buffer []vast.Stmt

	moduleMgr *module.DependecyResolver

	currentModule *module.Module

	enableShell bool

	results map[string]*vast.File

	callers container.Stack[CallFrame]

	unresolvedSymbols map[string][]linker.UnkownSymbol

	hcrDispatcher *transpiler.HCRDispatcher[vast.Node]
}

func NewTranspiler(opts *config.VBScriptOptions, moduleMgr *module.DependecyResolver) *Transpiler {
	return &Transpiler{
		opts:              opts,
		moduleMgr:         moduleMgr,
		results:           make(map[string]*vast.File),
		unresolvedSymbols: make(map[string][]linker.UnkownSymbol),
		hcrDispatcher:     transpiler.NewHCRDispatcher[vast.Node](),
		callers:           container.NewArrayStack[CallFrame](),
	}
}

type VBScriptSymbol struct {
	AST  *hast.UnresolvedSymbol
	node *vast.BasicLit
}

func (s *VBScriptSymbol) Source() string {
	return s.AST.Path
}

func (s *VBScriptSymbol) Name() string {
	return s.AST.Symbol
}

func (s *VBScriptSymbol) Link(symbol *linker.LinkableSymbol) error {
	if symbol == nil {
		return fmt.Errorf("symbol is nil")
	}
	s.node.Value = symbol.Text()
	return nil
}

func (t *Transpiler) RegisterDefaultRules() *Transpiler {
	return t
}

func (t *Transpiler) BindRules() error {
	return nil
}

func (t *Transpiler) UnresolvedSymbols() map[string][]linker.UnkownSymbol {
	return t.unresolvedSymbols
}

func (t *Transpiler) GetTargetExt() string {
	return ".sh"
}

func (t *Transpiler) GetTargetName() string {
	return "bash"
}

func (v *Transpiler) Convert(node hast.Node) vast.Node {
	switch node := node.(type) {
	case *hast.File:
		return v.ConvertFile(node)
	case *hast.CommentGroup:
		docs := make([]*vast.Comment, len(node.List))
		var tok vtok.Token
		if v.opts.CommentSyntax == "'" {
			tok = vtok.SGL_QUOTE
		} else {
			tok = vtok.REM
		}

		for i, d := range node.List {
			docs[i] = &vast.Comment{
				Tok:  tok,
				Text: d.Text,
			}
		}
		return &vast.CommentGroup{List: docs}
	case *hast.ExprStmt:
		expr := v.Convert(node.X)
		if s, ok := expr.(vast.Stmt); ok {
			return s
		}
		return &vast.ExprStmt{
			X: expr.(vast.Expr),
		}
	case *hast.UnsafeExpr:
		return v.ConvertUnsafeExpr(node)
	case *hast.ExternDecl:
		return v.ConvertExternDecl(node)
	case *hast.CallExpr:
		return v.ConvertCallExpr(node)
	case *hast.CmdExpr:
		return v.ConvertCmdExpr(node)
	case *hast.Ident:
		return &vast.Ident{Name: node.Name}
	case *hast.NumericLiteral:
		return &vast.BasicLit{
			Kind:  vtok.INTEGER,
			Value: node.Value,
		}
	case *hast.StringLiteral:
		return v.ConvertStringLiteral(node)
	case *hast.TrueLiteral:
		return &vast.BasicLit{
			Kind:  vtok.BOOLEAN,
			Value: "true",
		}
	case *hast.FalseLiteral:
		return &vast.BasicLit{
			Kind:  vtok.BOOLEAN,
			Value: "false",
		}
	case *hast.ArrayLiteralExpr:
		return v.ConvertArrayLiteral(node)
	case *hast.ObjectLiteralExpr:
		return v.ConvertObjectLiteral(node)
	case *hast.ForeachStmt:
		return v.ConvertForeachStmt(node)
	// case *hast.BasicLit:
	// 	return &vast.BasicLit{
	// 		Kind:  Token(node.Kind),
	// 		Value: node.Value,
	// 	}
	case *hast.AssignStmt:
		return v.ConvertAssignStmt(node)
	case *hast.IfStmt:
		return v.ConvertIfStmt(node)
	case *hast.BlockStmt:
		return v.ConvertBlock(node)
	case *hast.BinaryExpr:
		x := v.Convert(node.X).(vast.Expr)
		y := v.Convert(node.Y).(vast.Expr)
		return &vast.BinaryExpr{
			X:  x,
			Op: Token(node.Op),
			Y:  y,
		}
	case *hast.RefExpr:
		return v.Convert(node.X)
	case *hast.SelectExpr:
		return v.ConvertSelectExpr(node)
	case *hast.WhileStmt:
		body := v.Convert(node.Body).(*vast.BlockStmt)
		var (
			cond vast.Expr
			tok  vtok.Token
		)
		if node.Cond != nil {
			tok = vtok.WHILE
			cond = v.Convert(node.Cond).(vast.Expr)
		}
		return &vast.DoLoopStmt{
			Pre:  true,
			Tok:  tok,
			Cond: cond,
			Body: body,
		}
	case *hast.DoWhileStmt:
		body := v.Convert(node.Body).(*vast.BlockStmt)
		cond := v.Convert(node.Cond).(vast.Expr)
		return &vast.DoLoopStmt{
			Body: body,
			Tok:  vtok.WHILE,
			Cond: cond,
		}
	case *hast.ForStmt:
		return v.ConvertForStmt(node)
	// case *hast.RangeStmt:
	case *hast.IncDecExpr:
		return v.ConvertIncDecExpr(node)
	case *hast.FuncDecl:
		return v.ConvertFuncDecl(node)
	case *hast.Parameter:
		return v.ConvertParameter(node)
	case *hast.ReturnStmt:
		return v.ConvertReturnStmt(node)
	case *hast.ClassDecl:
		return v.ConvertClassDecl(node)
	case *hast.MatchStmt:
		return v.ConvertMatchStmt(node)
	case *hast.Import:
		return v.ConvertImport(node)
	case *hast.EnumDecl:
		return v.ConvertEnumDecl(node)
	default:
		fmt.Printf("%T\n", node)
	}
	return nil
}

func (v *Transpiler) ConvertUnsafeExpr(node *hast.UnsafeExpr) vast.Node {
	return &vast.ExprStmt{
		X: &vast.Ident{
			Name: node.Text,
		},
	}
}

func (v *Transpiler) ConvertExternDecl(node *hast.ExternDecl) vast.Node {
	for _, item := range node.List {
		expr := v.Convert(item).(vast.Expr)
		if expr, ok := expr.(*vast.Ident); ok {
			v.currentModule.Symbols.GlobalSymbols[expr.Name] = &module.BaseSymbol{}
		}
	}
	return nil
}

func (v *Transpiler) ConvertSelectExpr(node *hast.SelectExpr) vast.Node {
	return &vast.SelectorExpr{
		X:   v.Convert(node.X).(vast.Expr),
		Sel: v.Convert(node.Y).(vast.Expr),
	}
}

// isNumericEnum 检查是否是纯数字枚举
func (v *Transpiler) isNumericEnum(enumName string) bool {
	if enum := v.currentModule.LookupEnumSymbol(enumName); enum != nil {
		for _, value := range enum.Values {
			if value.Type != "num" {
				return false
			}
		}
		return true
	}
	return false
}

// getEnumNumericValue 获取枚举的数值
func (v *Transpiler) getEnumNumericValue(enumName, valueName string) vast.Expr {
	if enum := v.currentModule.LookupEnumSymbol(enumName); enum != nil {
		for _, value := range enum.Values {
			if value.Name == valueName {
				return &vast.BasicLit{
					Kind:  vtok.INTEGER,
					Value: value.Value,
				}
			}
		}
	}
	// 如果找不到，返回 0
	return &vast.BasicLit{
		Kind:  vtok.INTEGER,
		Value: "0",
	}
}

// getEnumStringValue 获取枚举的字符串值
func (v *Transpiler) getEnumStringValue(enumName, valueName string) vast.Expr {
	if enum := v.currentModule.LookupEnumSymbol(enumName); enum != nil {
		for _, value := range enum.Values {
			if value.Name == valueName {
				// 根据类型返回相应的字面量
				switch value.Type {
				case "string":
					return &vast.BasicLit{
						Kind:  vtok.STRING,
						Value: value.Value,
					}
				case "boolean":
					return &vast.BasicLit{
						Kind:  vtok.BOOLEAN,
						Value: value.Value,
					}
				default:
					return &vast.BasicLit{
						Kind:  vtok.INTEGER,
						Value: value.Value,
					}
				}
			}
		}
	}
	// 如果找不到，返回空字符串
	return &vast.BasicLit{
		Kind:  vtok.STRING,
		Value: "",
	}
}

func (v *Transpiler) ConvertImport(node *hast.Import) vast.Node {
	if v.unresolvedSymbols["Import"] == nil {
		node := &vast.BasicLit{
			Kind:  vtok.INTEGER,
			Value: "[PLACEHOLDER]",
		}
		v.Emit(&vast.ExprStmt{X: node})
		// Import 作为 Builtin 直接在未知符号声明一个即可，不然没法做这个需求
		v.unresolvedSymbols["Import"] = []linker.UnkownSymbol{
			&VBScriptSymbol{
				AST: &hast.UnresolvedSymbol{
					Symbol: "Import",
					Path:   "$HULO_PATH/core/unsafe/vbs/import.vbs",
				},
				node: node,
			},
		}
	}

	if node.ImportSingle != nil {
		importPath := node.ImportSingle.Path

		// 生成运行时导入代码
		// 使用Import函数在运行时加载模块
		return &vast.ExprStmt{
			X: &vast.CallExpr{
				Func: &vast.Ident{Name: "Import"},
				Recv: []vast.Expr{
					&vast.BasicLit{Kind: vtok.STRING, Value: fmt.Sprintf("%s.vbs", importPath)},
				},
			},
		}
	}
	return nil
}

func (v *Transpiler) ConvertIfStmt(node *hast.IfStmt) vast.Node {
	cond := v.Convert(node.Cond).(vast.Expr)
	body := v.Convert(node.Body).(*vast.BlockStmt)
	var elseStmt vast.Stmt
	if node.Else != nil {
		converted := v.Convert(node.Else)
		switch converted := converted.(type) {
		case *vast.BlockStmt:
			elseStmt = converted
		case *vast.IfStmt:
			elseStmt = converted
		default:
			// 如果不是 BlockStmt 或 IfStmt，包装成 BlockStmt
			elseStmt = &vast.BlockStmt{
				List: []vast.Stmt{converted.(vast.Stmt)},
			}
		}
	}
	return &vast.IfStmt{
		Cond: cond,
		Body: body,
		Else: elseStmt,
	}
}

func (v *Transpiler) ConvertReturnStmt(node *hast.ReturnStmt) vast.Node {
	var x vast.Expr
	if node.X != nil {
		x = v.Convert(node.X).(vast.Expr)
	} else {
		// 如果没有返回值，使用空字符串
		x = &vast.BasicLit{
			Kind:  vtok.STRING,
			Value: "",
		}
	}

	caller, ok := v.callers.Peek()
	if !ok {
		panic("return statement outside of function")
	}

	return &vast.AssignStmt{
		Lhs: &vast.Ident{Name: caller.(*FunctionFrame).Name},
		Rhs: x,
	}
}

func (v *Transpiler) ConvertStringLiteral(node *hast.StringLiteral) vast.Node {
	// 检查是否包含字符串插值
	if !strings.Contains(node.Value, "$") {
		return &vast.BasicLit{
			Kind:  vtok.STRING,
			Value: node.Value,
		}
	}

	// 解析字符串插值
	parts := ParseStringInterpolation(node.Value)

	// 如果只有一个部分，直接返回
	if len(parts) == 1 {
		if parts[0].IsVariable {
			return parts[0].Expr
		}
		return &vast.BasicLit{
			Kind:  vtok.STRING,
			Value: parts[0].Text,
		}
	}

	// 构建字符串连接表达式
	var result vast.Expr
	for i, part := range parts {
		var expr vast.Expr
		if part.IsVariable {
			expr = part.Expr
		} else {
			expr = &vast.BasicLit{
				Kind:  vtok.STRING,
				Value: v.escapeVBScriptString(part.Text),
			}
		}

		if i == 0 {
			result = expr
		} else {
			result = &vast.BinaryExpr{
				X:  result,
				Op: vtok.CONCAT,
				Y:  expr,
			}
		}
	}

	return result
}

// ConvertEnumDecl 将 Hulo 的枚举声明转换为 VBScript 的常量声明
func (v *Transpiler) ConvertEnumDecl(node *hast.EnumDecl) vast.Node {
	return nil
}

func (v *Transpiler) escapeVBScriptString(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}

func (v *Transpiler) canResolveFunction(fun *vast.Ident) bool {
	return v.currentModule.LookupFunctionSymbol(fun.Name) != nil
}

func (v *Transpiler) ConvertCallExpr(node *hast.CallExpr) vast.Node {
	convertedFun := v.Convert(node.Fun)

	// Handle different types of function expressions
	switch fun := convertedFun.(type) {
	case *vast.Ident:
		// Simple function name
		if v.currentModule.LookupClassSymbol(fun.Name) != nil {
			return &vast.CmdExpr{
				Cmd: &vast.Ident{Name: "New"},
				Recv: []vast.Expr{
					&vast.Ident{Name: fun.Name},
				},
			}
		}

		if !v.canResolveFunction(fun) {
			log.Debugf("has no function %s, enable shell", fun.Name)
			v.enableShell = true

			recv := []string{}
			for _, r := range node.Recv {
				recv = append(recv, vast.String(v.Convert(r).(vast.Expr)))
			}

			return &vast.SelectorExpr{
				X: &vast.Ident{Name: "shell"},
				Sel: &vast.CallExpr{
					Func: &vast.Ident{Name: "Exec"},
					Recv: []vast.Expr{
						&vast.BasicLit{
							Kind:  vtok.STRING,
							Value: fmt.Sprintf("cmd.exe /c %s %s", fun.Name, strings.Join(recv, " ")),
						},
					},
				},
			}
		} else {
			recv := []vast.Expr{}
			for _, r := range node.Recv {
				recv = append(recv, v.Convert(r).(vast.Expr))
			}
			return &vast.CallExpr{
				Func: fun,
				Recv: recv,
			}
		}
	// TODO: 这里有很大的问题 callExpr 居然套娃了 select, 说明语法树解析出问题了，现在暂时将错就错
	case *vast.SelectorExpr:
		// Handle module access like utils.Calculate() or $p.greet()
		// Check if this is a module access (e.g., utils.Calculate)
		if _, ok := fun.X.(*vast.Ident); ok {
			if funcIdent, ok := fun.Sel.(*vast.Ident); ok {
				// Check if this is a module access (e.g., utils.Calculate)
				// vs object method call (e.g., $p2.greet())
				if v.currentModule.LookupFunctionSymbol(funcIdent.Name) != nil {
					recv := []vast.Expr{}
					for _, r := range node.Recv {
						recv = append(recv, v.Convert(r).(vast.Expr))
					}
					return &vast.CallExpr{
						Func: funcIdent, // Just use the function name
						Recv: recv,
					}
				}
				// Otherwise, it's an object method call, keep the selector expression
			}
		}

		// Handle other selector expressions (e.g., $p.greet())
		recv := []vast.Expr{}
		for _, r := range node.Recv {
			recv = append(recv, v.Convert(r).(vast.Expr))
		}
		return &vast.CallExpr{
			Func: fun,
			Recv: recv,
		}
	default:
		// For other cases, just convert the function and arguments
		recv := []vast.Expr{}
		for _, r := range node.Recv {
			recv = append(recv, v.Convert(r).(vast.Expr))
		}
		return &vast.CallExpr{
			Func: convertedFun.(vast.Expr),
			Recv: recv,
		}
	}
}

func (v *Transpiler) ConvertCmdExpr(node *hast.CmdExpr) vast.Node {
	convertedFun := v.Convert(node.Cmd)

	// Handle different types of function expressions
	switch fun := convertedFun.(type) {
	case *vast.Ident:
		// Simple function name
		if v.currentModule.LookupClassSymbol(fun.Name) != nil {
			return &vast.CmdExpr{
				Cmd: &vast.Ident{Name: "New"},
				Recv: []vast.Expr{
					&vast.Ident{Name: fun.Name},
				},
			}
		}

		if !v.canResolveFunction(fun) {
			log.Debugf("has no function %s, enable shell", fun.Name)
			v.enableShell = true

			recv := []string{}
			for _, r := range node.Args {
				arg := vast.String(v.Convert(r).(vast.Expr))
				if len(arg) >= 2 && arg[0] == '"' && arg[len(arg)-1] == '"' {
					arg = fmt.Sprintf("\"\"%s\"\"", arg[1:len(arg)-1])
				}
				recv = append(recv, arg)
			}

			return &vast.SelectorExpr{
				X: &vast.Ident{Name: "shell"},
				Sel: &vast.CallExpr{
					Func: &vast.Ident{Name: "Exec"},
					Recv: []vast.Expr{
						&vast.BasicLit{
							Kind:  vtok.STRING,
							Value: fmt.Sprintf("cmd.exe /c %s %s", fun.Name, strings.Join(recv, " ")),
						},
					},
				},
			}
		} else {
			recv := []vast.Expr{}
			for _, r := range node.Args {
				recv = append(recv, v.Convert(r).(vast.Expr))
			}
			return &vast.CallExpr{
				Func: fun,
				Recv: recv,
			}
		}
	// TODO: 这里有很大的问题 callExpr 居然套娃了 select, 说明语法树解析出问题了，现在暂时将错就错
	case *vast.SelectorExpr:
		// Handle module access like utils.Calculate() or $p.greet()
		// Check if this is a module access (e.g., utils.Calculate)
		if _, ok := fun.X.(*vast.Ident); ok {
			if funcIdent, ok := fun.Sel.(*vast.Ident); ok {
				// Check if this is a module access (e.g., utils.Calculate)
				// vs object method call (e.g., $p2.greet())
				if v.currentModule.LookupFunctionSymbol(funcIdent.Name) != nil {
					recv := []vast.Expr{}
					for _, r := range node.Args {
						recv = append(recv, v.Convert(r).(vast.Expr))
					}
					return &vast.CallExpr{
						Func: funcIdent, // Just use the function name
						Recv: recv,
					}
				}
				// Otherwise, it's an object method call, keep the selector expression
			}
		}

		// Handle other selector expressions (e.g., $p.greet())
		recv := []vast.Expr{}
		for _, r := range node.Args {
			recv = append(recv, v.Convert(r).(vast.Expr))
		}
		return &vast.CallExpr{
			Func: fun,
			Recv: recv,
		}
	default:
		// For other cases, just convert the function and arguments
		recv := []vast.Expr{}
		for _, r := range node.Args {
			recv = append(recv, v.Convert(r).(vast.Expr))
		}
		return &vast.CallExpr{
			Func: convertedFun.(vast.Expr),
			Recv: recv,
		}
	}
}

func (v *Transpiler) ConvertFuncDecl(node *hast.FuncDecl) vast.Node {
	var ret *vast.FuncDecl = &vast.FuncDecl{}
	v.callers.Push(&FunctionFrame{Name: node.Name.Name})
	defer v.callers.Pop()

	// 处理修饰符
	var mod vtok.Token
	var modPos vtok.Pos
	for _, modifier := range node.Modifiers {
		if pubMod, ok := modifier.(*hast.PubModifier); ok {
			mod = vtok.PUBLIC
			modPos = vtok.Pos(pubMod.Pub)
			break
		}
	}

	ret.Mod = mod
	ret.ModPos = modPos
	ret.Function = vtok.Pos(node.Fn)
	ret.Name = v.Convert(node.Name).(*vast.Ident)
	ret.Recv = []*vast.Field{}
	for _, r := range node.Recv {
		ret.Recv = append(ret.Recv, v.Convert(r).(*vast.Field))
	}
	ret.Body = v.Convert(node.Body).(*vast.BlockStmt)
	ret.EndFunc = vtok.Pos(node.Body.Rbrace)

	// if v.currentModule.Symbols.HasFunction(ret.Name.Name) {
	// 	panic(fmt.Sprintf("function %s already declared", ret.Name.Name))
	// }

	return ret
}

func (v *Transpiler) ConvertParameter(node *hast.Parameter) vast.Node {
	return &vast.Field{
		Name: v.Convert(node.Name).(*vast.Ident),
		Tok:  vtok.BYVAL,
	}
}

func (v *Transpiler) ConvertBlock(node *hast.BlockStmt) vast.Node {
	stmts := []vast.Stmt{}
	for _, s := range node.List {
		stmt := v.Convert(s)
		stmts = append(stmts, v.Flush()...)
		if stmt == nil {
			continue
		}
		if s, ok := stmt.(vast.Stmt); ok {
			stmts = append(stmts, s)
		} else {
			// 不是vast.Stmt类型，忽略或报错
			// 可以选择panic或者continue，这里选择continue
			continue
		}
	}
	return &vast.BlockStmt{List: stmts}
}

func (v *Transpiler) ConvertIncDecExpr(node *hast.IncDecExpr) vast.Node {
	counter := v.Convert(node.X).(*vast.Ident)
	var op vtok.Token
	if node.Tok == htok.INC {
		op = vtok.ADD
	} else {
		op = vtok.SUB
	}

	return &vast.AssignStmt{
		Lhs: counter,
		Rhs: &vast.BinaryExpr{
			X:  counter,
			Op: op,
			Y: &vast.BasicLit{
				Kind:  vtok.INTEGER,
				Value: "1",
			},
		},
	}
}

func (v *Transpiler) ConvertClassDecl(node *hast.ClassDecl) vast.Node {
	// 转换类名
	className := v.Convert(node.Name).(*vast.Ident)

	// 转换修饰符
	var mod vtok.Token
	var modPos vtok.Pos

	if node.Pub.IsValid() {
		mod = vtok.PUBLIC
		modPos = vtok.Pos(node.Pub)
	}

	// 收集字段信息
	var fieldNames []string
	var stmts []vast.Stmt

	// 处理字段
	if node.Fields != nil {
		for _, field := range node.Fields.List {
			// 检查字段是否有 pub 修饰符
			var isPublic bool
			for _, modifier := range field.Modifiers {
				if _, ok := modifier.(*hast.PubModifier); ok {
					isPublic = true
					break
				}
			}

			// 将字段转换为相应的声明
			fieldName := v.Convert(field.Name).(*vast.Ident)
			fieldNames = append(fieldNames, fieldName.Name)

			if isPublic {
				// 公共字段使用 PublicStmt
				publicStmt := &vast.PublicStmt{
					Public: vtok.Pos(field.Name.Pos()),
					List:   []vast.Expr{fieldName},
				}
				stmts = append(stmts, publicStmt)
			} else {
				// 私有字段使用 PrivateStmt
				privateStmt := &vast.PrivateStmt{
					Private: vtok.Pos(field.Name.Pos()),
					List:    []vast.Expr{fieldName},
				}
				stmts = append(stmts, privateStmt)
			}

			// 如果有默认值，添加赋值语句
			if field.Value != nil {
				assignStmt := &vast.AssignStmt{
					Lhs: fieldName,
					Rhs: v.Convert(field.Value).(vast.Expr),
				}
				stmts = append(stmts, assignStmt)
			}
		}
	}

	// 处理方法
	for _, method := range node.Methods {
		funcDecl := v.Convert(method).(*vast.FuncDecl)
		stmts = append(stmts, funcDecl)
	}

	// 处理构造函数
	for _, ctor := range node.Ctors {
		// 构造函数转换为函数声明
		ctorName := v.Convert(ctor.Name).(*vast.Ident)
		ctorBody := v.Convert(ctor.Body).(*vast.BlockStmt)

		// 转换参数
		var recv []*vast.Field
		for _, param := range ctor.Recv {
			if param, ok := param.(*hast.Parameter); ok {
				field := &vast.Field{
					TokPos: vtok.Pos(param.Name.Pos()),
					Tok:    vtok.BYVAL,
					Name:   v.Convert(param.Name).(*vast.Ident),
				}
				recv = append(recv, field)
			}
		}

		funcDecl := &vast.FuncDecl{
			Mod:      mod,
			ModPos:   modPos,
			Function: vtok.Pos(ctor.Name.Pos()),
			Name:     ctorName,
			Recv:     recv,
			Body:     ctorBody,
			EndFunc:  vtok.Pos(ctor.Body.Rbrace),
		}
		stmts = append(stmts, funcDecl)
	}

	return &vast.ClassDecl{
		Mod:      mod,
		ModPos:   modPos,
		Class:    vtok.Pos(node.Class),
		Name:     className,
		Stmts:    stmts,
		EndClass: vtok.Pos(node.Rbrace),
	}
}

func (v *Transpiler) ConvertFile(node *hast.File) vast.Node {
	docs := make([]*vast.CommentGroup, len(node.Docs))
	for i, d := range node.Docs {
		docs[i] = v.Convert(d).(*vast.CommentGroup)
	}

	stmts := []vast.Stmt{}
	for _, s := range node.Stmts {
		stmt := v.Convert(s)

		stmts = append(stmts, v.Flush()...)

		if stmt == nil {
			continue
		}
		stmts = append(stmts, stmt.(vast.Stmt))
	}

	if v.enableShell {
		stmts = append([]vast.Stmt{
			&vast.SetStmt{
				Lhs: &vast.Ident{Name: "shell"},
				Rhs: &vast.CallExpr{
					Func: &vast.Ident{Name: "CreateObject"},
					Recv: []vast.Expr{
						&vast.BasicLit{
							Kind:  vtok.STRING,
							Value: "WScript.Shell",
						},
					},
				},
			},
		}, stmts...)
	}

	return &vast.File{
		Doc:   docs,
		Stmts: stmts,
	}
}

func (v *Transpiler) ConvertForStmt(node *hast.ForStmt) vast.Node {
	v.callers.Push(&LoopFrame{})
	defer v.callers.Pop()

	init := v.Convert(node.Init).(vast.Stmt)
	cond := v.Convert(node.Cond).(vast.Expr)
	post := v.Convert(node.Post).(vast.Stmt)

	stmts := []vast.Stmt{}
	body := v.Convert(node.Body).(*vast.BlockStmt)

	stmts = append(stmts, v.Flush()...)
	stmts = append(stmts, body.List...)
	stmts = append(stmts, post)

	v.Emit(init)

	return &vast.DoLoopStmt{
		Pre:  true,
		Tok:  vtok.WHILE,
		Cond: cond,
		Body: &vast.BlockStmt{List: stmts},
	}
}

// isObjectExpression 判断表达式是否是对象类型
func (v *Transpiler) isObjectExpression(expr vast.Expr) bool {
	switch e := expr.(type) {
	case *vast.CallExpr:
		// CreateObject 调用返回对象
		if ident, ok := e.Func.(*vast.Ident); ok && ident.Name == "CreateObject" {
			return true
		}
		// Array 调用返回数组，不需要 Set
		if ident, ok := e.Func.(*vast.Ident); ok && ident.Name == "Array" {
			return false
		}
	case *vast.NewExpr:
		// New 表达式创建对象
		return true
	case *vast.Ident:
		// 检查是否是已知的对象变量（如Dictionary）
		if e.Name == "dict_0" || strings.HasPrefix(e.Name, "dict_") {
			return true
		}
	}
	return false
}

func (v *Transpiler) ConvertAssignStmt(node *hast.AssignStmt) vast.Node {
	lhs := v.Convert(node.Lhs).(vast.Expr)
	rhs := v.Convert(node.Rhs).(vast.Expr)
	if node.Scope == htok.CONST {
		return &vast.ConstStmt{
			Lhs: lhs.(*vast.Ident),
			Rhs: rhs.(*vast.BasicLit),
		}
	}

	// 检查是否是构造函数调用赋值
	if callExpr, ok := node.Rhs.(*hast.CallExpr); ok {
		if ident, ok := callExpr.Fun.(*hast.Ident); ok {
			if v.currentModule.LookupClassSymbol(ident.Name) != nil {
				// 这是构造函数调用赋值，需要特殊处理
				return v.convertConstructorAssignment(lhs, ident.Name, callExpr.Recv)
			}
		}
	}

	if v.needsDimDeclaration(lhs) {
		v.Emit(&vast.DimDecl{List: []vast.Expr{lhs}})
	}

	// 检查右侧表达式是否是对象类型
	if v.isObjectExpression(rhs) {
		// 对象赋值需要使用 Set
		return &vast.SetStmt{
			Lhs: lhs,
			Rhs: rhs,
		}
	}

	return &vast.AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}

func (v *Transpiler) Emit(n ...vast.Stmt) {
	v.buffer = append(v.buffer, n...)
}

func (v *Transpiler) Flush() []vast.Stmt {
	stmts := v.buffer
	v.buffer = nil
	return stmts
}

func (v *Transpiler) needsDimDeclaration(expr vast.Expr) bool {
	// 检查当前作用域
	caller, ok := v.callers.Peek()
	if ok && caller.Kind() == CallerLoop {
		// 在循环作用域中不需要Dim声明
		return false
	}

	// 检查变量是否已经声明过
	if ident, ok := expr.(*vast.Ident); ok {
		return v.currentModule.Symbols.HasSymbol(ident.Name)
	}

	return false
}

// convertConstructorAssignment 处理构造函数调用赋值
// 例如：let p2 = Person("Jerry", 20) 转换为：
// Dim p2
// Set p2 = New Person
// p2.name = "Jerry"
// p2.age = 20
func (v *Transpiler) convertConstructorAssignment(lhs vast.Expr, className string, args []hast.Expr) vast.Node {
	// 检查是否需要 Dim 声明
	if v.needsDimDeclaration(lhs) {
		v.Emit(&vast.DimDecl{List: []vast.Expr{lhs}})
	}

	// 创建对象实例
	createStmt := &vast.SetStmt{
		Lhs: lhs,
		Rhs: &vast.NewExpr{
			X: &vast.Ident{Name: className},
		},
	}

	// 将创建语句添加到缓冲区
	v.Emit(createStmt)

	// 根据类的字段顺序，为字段赋值
	// 这里我们需要知道类的字段信息，暂时使用简单的映射
	fieldNames := v.getClassFieldNames(className)

	var stmts []vast.Stmt
	for i, arg := range args {
		if i < len(fieldNames) {
			fieldName := fieldNames[i]
			assignStmt := &vast.AssignStmt{
				Lhs: &vast.SelectorExpr{
					X:   lhs,
					Sel: &vast.Ident{Name: fieldName},
				},
				Rhs: v.Convert(arg).(vast.Expr),
			}
			stmts = append(stmts, assignStmt)
		}
	}

	// 将字段赋值语句添加到缓冲区
	for _, stmt := range stmts {
		v.Emit(stmt)
	}

	// 返回 nil，因为实际的语句已经通过 Emit 添加了
	return nil
}

// getClassFieldNames 获取类的字段名列表
func (v *Transpiler) getClassFieldNames(className string) []string {
	sym := v.currentModule.LookupClassSymbol(className)
	if sym == nil {
		return []string{}
	}
	fields := sym.Fields()
	fieldNames := make([]string, len(fields))
	i := 0
	for _, field := range fields {
		fieldNames[i] = field.Name()
		i++
	}
	return fieldNames
}

// ConvertMatchStmt 将 Hulo 的 match 语句转换为 VBScript 的 if-elseif-else 语句
func (v *Transpiler) ConvertMatchStmt(node *hast.MatchStmt) vast.Node {
	v.callers.Push(&BlockFrame{})
	defer v.callers.Pop()

	// 转换匹配表达式
	matchExpr := v.Convert(node.Expr).(vast.Expr)

	// 构建 if-elseif-else 链
	var currentIf *vast.IfStmt
	var firstIf *vast.IfStmt

	// 处理所有 case 子句
	for _, caseClause := range node.Cases {
		// 转换条件表达式
		cond := v.convertMatchCondition(matchExpr, caseClause.Cond)

		// 转换 case 体
		var caseBody vast.Stmt
		if caseClause.Body != nil {
			converted := v.Convert(caseClause.Body)
			if block, ok := converted.(*vast.BlockStmt); ok {
				caseBody = block
			} else {
				caseBody = &vast.BlockStmt{
					List: []vast.Stmt{converted.(vast.Stmt)},
				}
			}
		} else {
			caseBody = &vast.BlockStmt{List: []vast.Stmt{}}
		}

		// 创建 if 语句
		ifStmt := &vast.IfStmt{
			Cond: cond,
			Body: caseBody.(*vast.BlockStmt),
		}

		if currentIf == nil {
			// 第一个 if 语句
			firstIf = ifStmt
			currentIf = ifStmt
		} else {
			// 将当前 if 语句作为前一个的 else 分支
			currentIf.Else = ifStmt
			currentIf = ifStmt
		}
	}

	// 处理默认分支
	if node.Default != nil {
		var defaultBody vast.Stmt
		if node.Default.Body != nil {
			converted := v.Convert(node.Default.Body)
			if block, ok := converted.(*vast.BlockStmt); ok {
				defaultBody = block
			} else {
				defaultBody = &vast.BlockStmt{
					List: []vast.Stmt{converted.(vast.Stmt)},
				}
			}
		} else {
			defaultBody = &vast.BlockStmt{List: []vast.Stmt{}}
		}

		if currentIf != nil {
			currentIf.Else = defaultBody
		} else {
			// 如果没有 case 子句，只有默认分支
			firstIf = &vast.IfStmt{
				If: vtok.Pos(1), // 使用动态位置
				Cond: &vast.BasicLit{
					Kind:  vtok.BOOLEAN,
					Value: "True",
				},
				Then:  vtok.Pos(2), // 使用动态位置
				Body:  defaultBody.(*vast.BlockStmt),
				EndIf: vtok.Pos(3), // 使用动态位置
			}
		}
	}

	return firstIf
}

func (v *Transpiler) convertMatchCondition(matchExpr vast.Expr, cond hast.Expr) vast.Expr {
	switch c := cond.(type) {
	case *hast.StringLiteral:
		// 字符串字面量匹配
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y: &vast.BasicLit{
				Kind:  vtok.STRING,
				Value: c.Value,
			},
		}
	case *hast.NumericLiteral:
		// 数字字面量匹配
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y: &vast.BasicLit{
				Kind:  vtok.INTEGER,
				Value: c.Value,
			},
		}
	case *hast.TrueLiteral:
		// true 字面量匹配
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y: &vast.BasicLit{
				Kind:  vtok.BOOLEAN,
				Value: "True",
			},
		}
	case *hast.FalseLiteral:
		// false 字面量匹配
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y: &vast.BasicLit{
				Kind:  vtok.BOOLEAN,
				Value: "False",
			},
		}
	case *hast.Ident:
		// 标识符匹配（可能是枚举值或其他变量）
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y:  &vast.Ident{Name: c.Name},
		}
	case *hast.BasicLit:
		// 基本字面量匹配
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y: &vast.BasicLit{
				Kind:  vtok.Token(c.Kind),
				Value: c.Value,
			},
		}
	default:
		// 对于其他类型的表达式，直接使用相等比较
		return &vast.BinaryExpr{
			X:  matchExpr,
			Op: vtok.EQ,
			Y:  v.Convert(c).(vast.Expr),
		}
	}
}

// ConvertArrayLiteral 将Hulo数组字面量转换为VBScript的Array函数调用
func (v *Transpiler) ConvertArrayLiteral(node *hast.ArrayLiteralExpr) vast.Node {
	var args []vast.Expr

	// 转换所有数组元素
	for _, elem := range node.Elems {
		converted := v.Convert(elem)
		if converted != nil {
			args = append(args, converted.(vast.Expr))
		}
	}

	// 创建Array函数调用
	return &vast.CallExpr{
		Func:   &vast.Ident{Name: "Array"},
		Lparen: vtok.Pos(node.Lbrack),
		Recv:   args,
		Rparen: vtok.Pos(node.Rbrack),
	}
}

// ConvertObjectLiteral 将Hulo对象字面量转换为VBScript的Dictionary对象
func (v *Transpiler) ConvertObjectLiteral(node *hast.ObjectLiteralExpr) vast.Node {
	// 生成唯一的变量名
	dictVarName := fmt.Sprintf("dict_%d", len(v.buffer))

	// 生成创建Dictionary的代码
	createDict := &vast.SetStmt{
		Lhs: &vast.Ident{Name: dictVarName},
		Rhs: &vast.CallExpr{
			Func: &vast.Ident{Name: "CreateObject"},
			Recv: []vast.Expr{
				&vast.BasicLit{
					Kind:  vtok.STRING,
					Value: "Scripting.Dictionary",
				},
			},
		},
	}

	// 将创建语句添加到缓冲区
	v.Emit(createDict)

	// 处理所有属性
	for _, prop := range node.Props {
		if keyValue, ok := prop.(*hast.KeyValueExpr); ok {
			// 转换键和值
			key := v.Convert(keyValue.Key).(vast.Expr)
			value := v.Convert(keyValue.Value).(vast.Expr)

			// 生成Add方法调用
			addCall := &vast.CmdExpr{
				Cmd: &vast.SelectorExpr{
					X:   &vast.Ident{Name: dictVarName},
					Sel: &vast.Ident{Name: "Add"},
				},
				Recv: []vast.Expr{key, value},
			}

			// 将Add调用添加到缓冲区
			v.Emit(&vast.ExprStmt{X: addCall})
		}
	}

	// 返回Dictionary变量引用
	return &vast.Ident{Name: dictVarName}
}

const WILDCARD = "_"

// ConvertForeachStmt 将Hulo的foreach语句转换为VBScript的For Each语句
func (v *Transpiler) ConvertForeachStmt(node *hast.ForeachStmt) vast.Node {
	// 转换集合表达式
	collectionExpr := v.Convert(node.Var).(vast.Expr)

	// 如果是对象遍历（使用of关键字），需要特殊处理
	if node.Tok == htok.OF {
		// 对于对象遍历，我们需要遍历Dictionary的Keys
		// 修改集合表达式为 collectionExpr.Keys
		collectionExpr = &vast.SelectorExpr{
			X:   collectionExpr,
			Sel: &vast.Ident{Name: "Keys"},
		}

		// 转换循环变量
		var elemExpr vast.Expr
		if node.Index != nil {
			// 检查是否是下划线占位符
			if ident, ok := node.Index.(*hast.Ident); ok && ident.Name == WILDCARD {
				// 使用临时变量名
				elemExpr = &vast.Ident{Name: "tempKey"}
			} else {
				elemExpr = v.Convert(node.Index).(vast.Expr)
			}
		} else {
			// 如果没有指定变量名，使用默认的key
			elemExpr = &vast.Ident{Name: "key"}
		}

		// 转换循环体
		body := v.Convert(node.Body).(*vast.BlockStmt)

		// 如果有两个变量（key和value），需要在循环体中添加value变量的赋值
		if node.Value != nil {
			// 检查value是否是下划线占位符
			if valueIdent, ok := node.Value.(*hast.Ident); ok && valueIdent.Name == WILDCARD {
				// 如果value是下划线，不需要赋值语句，直接使用循环体
				// 但需要确保key变量名是有效的
				if ident, ok := node.Index.(*hast.Ident); ok && ident.Name == WILDCARD {
					// 如果key也是下划线，使用临时变量名
					elemExpr = &vast.Ident{Name: "tempKey"}
				}
			} else {
				// 创建新的循环体，在开头添加value赋值语句
				var newBodyStmts []vast.Stmt

				// 添加 value = collectionExpr.Item(key) 语句
				valueAssignStmt := &vast.AssignStmt{
					Lhs: v.Convert(node.Value).(vast.Expr),
					Rhs: &vast.CallExpr{
						Func: &vast.SelectorExpr{
							X:   v.Convert(node.Var).(vast.Expr),
							Sel: &vast.Ident{Name: "Item"},
						},
						Recv: []vast.Expr{elemExpr},
					},
				}
				newBodyStmts = append(newBodyStmts, valueAssignStmt)

				// 添加原有的循环体语句
				newBodyStmts = append(newBodyStmts, body.List...)

				body = &vast.BlockStmt{List: newBodyStmts}
			}
		}

		// 创建For Each语句
		forEachStmt := &vast.ForEachStmt{
			For:   vtok.Pos(node.Loop),
			Each:  vtok.Pos(node.Loop), // 使用相同位置
			Elem:  elemExpr,
			In:    vtok.Pos(node.Tok),
			Group: collectionExpr,
			Body:  body,
			Next:  vtok.Pos(node.Body.Rbrace), // 使用循环体结束位置
		}

		return forEachStmt
	} else {
		// 对于数组遍历（使用in关键字），直接处理
		var elemExpr vast.Expr
		if node.Index != nil {
			elemExpr = v.Convert(node.Index).(vast.Expr)
		} else {
			// 如果没有指定变量名，使用默认的item
			elemExpr = &vast.Ident{Name: "item"}
		}

		// 转换循环体
		body := v.Convert(node.Body).(*vast.BlockStmt)

		// 创建For Each语句
		forEachStmt := &vast.ForEachStmt{
			For:   vtok.Pos(node.Loop),
			Each:  vtok.Pos(node.Loop), // 使用相同位置
			Elem:  elemExpr,
			In:    vtok.Pos(node.Tok),
			Group: collectionExpr,
			Body:  body,
			Next:  vtok.Pos(node.Body.Rbrace), // 使用循环体结束位置
		}

		return forEachStmt
	}
}
