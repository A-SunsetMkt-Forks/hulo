// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	vtok "github.com/hulo-lang/hulo/syntax/vbs/token"
)

func Transpile(opts *config.Huloc, vfs vfs.VFS, basePath string) (map[string]string, error) {
	tr := NewVBScriptTranspiler(opts, vfs)
	tr.moduleManager.basePath = basePath
	return tr.Transpile(opts.Main)
}

type VBScriptTranspiler struct {
	opts *config.Huloc

	buffer []vast.Stmt

	moduleManager *ModuleManager

	vfs vfs.VFS

	modules       map[string]*Module
	currentModule *Module

	builtinModules []*Module

	globalSymbols *SymbolTable

	scopeStack container.Stack[ScopeType]

	currentScope *Scope

	enableShell bool

	undefinedSymbols container.Set[string] // TODO: 需要移动到 Module 里面
}

func NewVBScriptTranspiler(opts *config.Huloc, vfs vfs.VFS) *VBScriptTranspiler {
	return &VBScriptTranspiler{
		opts:             opts,
		vfs:              vfs,
		undefinedSymbols: container.NewMapSet[string](),
		// 初始化模块管理
		moduleManager: &ModuleManager{
			isFirst:  true,
			options:  opts,
			vfs:      vfs,
			basePath: "",
			modules:  make(map[string]*Module),
			resolver: &DependencyResolver{
				visited: container.NewMapSet[string](),
				stack:   container.NewMapSet[string](),
				order:   make([]string, 0),
			},
			externalVBS: make(map[string]string),
		},
		modules: make(map[string]*Module),

		// 初始化符号管理
		globalSymbols: NewSymbolTable("global", opts.EnableMangle),

		// 初始化作用域管理
		scopeStack: container.NewArrayStack[ScopeType](),
	}
}

func (b *VBScriptTranspiler) Transpile(mainFile string) (map[string]string, error) {
	// 解析内置模块
	builtinModules, err := b.moduleManager.ResolveAllDependencies(filepath.Join(b.opts.HuloPath, "core/unsafe/vbs/index.hl"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve builtin dependencies: %w", err)
	}

	b.builtinModules = builtinModules

	// 1. 解析所有依赖
	modules, err := b.moduleManager.ResolveAllDependencies(mainFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// 2. 将解析的模块同步到 transpiler 的 modules 中
	for _, module := range modules {
		b.modules[module.Path] = module
		log.Debugf("Loaded module: %s (path: %s)", module.Name, module.Path)
	}

	// 3. 分析所有模块的导出符号
	for _, module := range modules {
		if module.Symbols == nil {
			module.Symbols = NewSymbolTable(module.Name, false)
		}
		b.analyzeModuleExports(module)
	}

	// 4. 按依赖顺序翻译模块
	results := make(map[string]string)
	for _, module := range modules {
		if err := b.translateModule(module); err != nil {
			return nil, fmt.Errorf("failed to translate module %s: %w", module.Path, err)
		}
		path := strings.Replace(module.Path, ".hl", ".vbs", 1)
		results[path] = vast.String(module.Transpiled)
	}

	return results, nil
}

func (b *VBScriptTranspiler) translateModule(module *Module) error {
	// 设置当前模块
	b.currentModule = module

	// 翻译模块内容
	module.Transpiled = b.Convert(module.AST).(*vast.File)

	return nil
}

func (b *VBScriptTranspiler) analyzeModuleExports(module *Module) {
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

		case *hast.DeclareDecl:
			// 处理 declare fn
			if funcDecl, ok := s.X.(*hast.FuncDecl); ok {
				module.Exports[funcDecl.Name.Name] = &ExportInfo{
					Name:   funcDecl.Name.Name,
					Value:  funcDecl,
					Kind:   ExportFunction,
					Public: true, // declare 默认视为可用
				}
				module.Symbols.AddFunction(funcDecl.Name.Name, funcDecl)
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
				fields := make([]string, 0)
				for _, field := range s.Fields.List {
					fields = append(fields, field.Name.Name)
				}
				module.Symbols.AddClass(s.Name.Name, &ClassSymbol{
					Name:   s.Name.Name,
					Fields: fields,
					Public: s.Pub.IsValid(),
				})
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

func (v *VBScriptTranspiler) Convert(node hast.Node) vast.Node {
	switch node := node.(type) {
	case *hast.File:
		return v.ConvertFile(node)
	case *hast.CommentGroup:
		docs := make([]*vast.Comment, len(node.List))
		var tok vtok.Token
		if v.opts.CompilerOptions.VBScript.CommentSyntax == "'" {
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
	case *hast.UnsafeStmt:
		return v.ConvertUnsafeStmt(node)
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
	case *hast.DeclareDecl:
		return v.ConvertDeclareDecl(node)
	case *hast.Import:
		return v.ConvertImport(node)
	case *hast.EnumDecl:
		return v.ConvertEnumDecl(node)
	case *hast.ModAccessExpr:
		return v.ConvertModAccessExpr(node)
	default:
		fmt.Printf("%T\n", node)
	}
	return nil
}

func (v *VBScriptTranspiler) ConvertUnsafeStmt(node *hast.UnsafeStmt) vast.Node {
	return &vast.ExprStmt{
		X: &vast.Ident{
			Name: node.Text,
		},
	}
}

func (v *VBScriptTranspiler) ConvertExternDecl(node *hast.ExternDecl) vast.Node {
	for _, item := range node.List {
		expr := v.Convert(item).(vast.Expr)
		if expr, ok := expr.(*vast.Ident); ok {
			v.currentModule.Symbols.AddVariable(expr.Name, item)
		}
	}
	return nil
}

func (v *VBScriptTranspiler) ConvertSelectExpr(node *hast.SelectExpr) vast.Node {
	return &vast.SelectorExpr{
		X:   v.Convert(node.X).(vast.Expr),
		Sel: v.Convert(node.Y).(vast.Expr),
	}
}

// ConvertModAccessExpr 将模块访问表达式转换为相应的常量引用
// 例如：Status::Pending 转换为 0（数字枚举）或生成常量声明（字符串枚举）
func (v *VBScriptTranspiler) ConvertModAccessExpr(node *hast.ModAccessExpr) vast.Node {
	// 获取模块名（枚举名）
	moduleName := v.Convert(node.X).(*vast.Ident).Name

	// 获取成员名（枚举值）
	memberName := v.Convert(node.Y).(*vast.Ident).Name

	// 检查是否是枚举访问
	if v.currentModule.Symbols.HasEnum(moduleName) && v.currentModule.Symbols.HasEnumValue(moduleName, memberName) {
		// 检查是否是纯数字枚举，如果是则直接返回数值
		if v.isNumericEnum(moduleName) {
			// 对于数字枚举，直接返回数值
			return v.getEnumNumericValue(moduleName, memberName)
		} else {
			// 对于字符串/混合枚举，生成常量声明并返回常量名
			constName := fmt.Sprintf("%s_%s", moduleName, memberName)

			// 检查是否已经声明过这个常量
			if !v.currentModule.Symbols.HasConstant(constName) {
				constValue := v.getEnumStringValue(moduleName, memberName)

				// 生成常量声明
				constStmt := &vast.ConstStmt{
					Lhs: &vast.Ident{Name: constName},
					Rhs: constValue,
				}
				v.Emit(constStmt)

				// 标记为已声明
				v.currentModule.Symbols.AddConstant(constName, node)
			}

			// 返回常量名
			return &vast.Ident{Name: constName}
		}
	}

	// 不是枚举访问，保持原样（可能需要其他处理）
	return &vast.SelectorExpr{
		X:   v.Convert(node.X).(vast.Expr),
		Sel: v.Convert(node.Y).(vast.Expr),
	}
}

// isNumericEnum 检查是否是纯数字枚举
func (v *VBScriptTranspiler) isNumericEnum(enumName string) bool {
	if enum, exists := v.currentModule.Symbols.GetEnum(enumName); exists {
		for _, value := range enum.Values {
			if value.Type != "number" {
				return false
			}
		}
		return true
	}
	return false
}

// getEnumNumericValue 获取枚举的数值
func (v *VBScriptTranspiler) getEnumNumericValue(enumName, valueName string) vast.Expr {
	if enum, exists := v.currentModule.Symbols.GetEnum(enumName); exists {
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
func (v *VBScriptTranspiler) getEnumStringValue(enumName, valueName string) vast.Expr {
	if enum, exists := v.currentModule.Symbols.GetEnum(enumName); exists {
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

func (v *VBScriptTranspiler) ConvertImport(node *hast.Import) vast.Node {
	if node.ImportSingle != nil {
		importPath := node.ImportSingle.Path
		v.undefinedSymbols.Add("Import")
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

func (v *VBScriptTranspiler) ConvertDeclareDecl(node *hast.DeclareDecl) vast.Node {
	switch n := node.X.(type) {
	case *hast.FuncDecl:
		fnName := v.Convert(n.Name).(*vast.Ident)
		v.currentModule.Symbols.AddFunction(fnName.Name, n)
		return nil
	case *hast.BlockStmt:
		for _, stmt := range n.List {
			switch stmt := stmt.(type) {
			case *hast.AssignStmt:
				lhs := v.Convert(stmt.Lhs).(vast.Expr)
				v.currentModule.Symbols.AddVariable(lhs.String(), stmt)
			default:
			}
		}
		return nil
	default:
		panic(fmt.Sprintf("unsupported declare decl: %T", n))
	}
}

func (v *VBScriptTranspiler) ConvertIfStmt(node *hast.IfStmt) vast.Node {
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

func (v *VBScriptTranspiler) ConvertReturnStmt(node *hast.ReturnStmt) vast.Node {
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

	scope := v.currentModule.Symbols.CurrentScope()
	if scope == nil || scope.Type != FunctionScope {
		panic("return statement outside of function")
	}

	return &vast.AssignStmt{
		Lhs: &vast.Ident{Name: scope.ID},
		Rhs: x,
	}
}

func (v *VBScriptTranspiler) ConvertStringLiteral(node *hast.StringLiteral) vast.Node {
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
func (v *VBScriptTranspiler) ConvertEnumDecl(node *hast.EnumDecl) vast.Node {
	enumName := node.Name.Name
	var enumValues []*EnumValue // 收集枚举值用于符号表

	// 根据枚举体类型处理
	switch body := node.Body.(type) {
	case *hast.BasicEnumBody:
		// 简单枚举：enum Status { Pending, Approved, Rejected }
		// 检查是否有显式值，如果有，按递增逻辑；如果没有，按索引逻辑
		hasExplicitValues := false
		for _, value := range body.Values {
			if value.Value != nil {
				hasExplicitValues = true
				break
			}
		}

		if hasExplicitValues {
			// 有关联值的枚举，使用递增逻辑
			lastValue := 0
			for _, value := range body.Values {
				var constValue vast.Expr
				var enumValue *EnumValue

				if value.Value != nil {
					// 如果有显式值，使用它
					constValue = v.Convert(value.Value).(vast.Expr)

					// 确定值的类型
					var valueType string
					if lit, ok := constValue.(*vast.BasicLit); ok {
						switch lit.Kind {
						case vtok.INTEGER:
							valueType = "number"
							// 更新 lastValue（如果是数字的话）
							if val, err := strconv.Atoi(lit.Value); err == nil {
								lastValue = val
							}
						case vtok.STRING:
							valueType = "string"
						case vtok.BOOLEAN:
							valueType = "boolean"
						default:
							valueType = "number"
						}
					} else {
						valueType = "number"
					}

					// 创建枚举值对象
					var rawValue string
					if lit, ok := constValue.(*vast.BasicLit); ok {
						rawValue = lit.Value
					} else {
						rawValue = constValue.String()
					}

					enumValue = &EnumValue{
						Name:  value.Name.Name,
						Value: rawValue,
						Type:  valueType,
					}
				} else {
					// 没有显式值，使用上一个值 + 1
					lastValue++
					constValue = &vast.BasicLit{
						Kind:  vtok.INTEGER,
						Value: fmt.Sprintf("%d", lastValue),
					}

					// 创建枚举值对象
					enumValue = &EnumValue{
						Name:  value.Name.Name,
						Value: fmt.Sprintf("%d", lastValue),
						Type:  "number",
					}
				}

				enumValues = append(enumValues, enumValue) // 收集枚举值
			}
		} else {
			// 简单枚举，使用索引逻辑
			for i, value := range body.Values {
				// 创建枚举值对象
				enumValue := &EnumValue{
					Name:  value.Name.Name,
					Value: fmt.Sprintf("%d", i),
					Type:  "number",
				}
				enumValues = append(enumValues, enumValue) // 收集枚举值
			}
		}

	case *hast.AssociatedEnumBody:
		// 关联枚举：enum HttpCode { OK = 200, NotFound = 404, ... }
		lastValue := 0 // 跟踪上一个值，用于自动递增
		for _, value := range body.Values {
			var constValue vast.Expr
			var enumValue *EnumValue

			if value.Data != nil && len(value.Data) > 0 {
				// 如果有数据，使用第一个值
				constValue = v.Convert(value.Data[0]).(vast.Expr)

				// 确定值的类型
				var valueType string
				if lit, ok := constValue.(*vast.BasicLit); ok {
					switch lit.Kind {
					case vtok.INTEGER:
						valueType = "number"
						// 更新 lastValue（如果是数字的话）
						if val, err := strconv.Atoi(lit.Value); err == nil {
							lastValue = val
						}
					case vtok.STRING:
						valueType = "string"
					case vtok.BOOLEAN:
						valueType = "boolean"
					default:
						valueType = "number"
					}
				} else {
					valueType = "number"
				}

				// 创建枚举值对象
				var rawValue string
				if lit, ok := constValue.(*vast.BasicLit); ok {
					rawValue = lit.Value
				} else {
					rawValue = constValue.String()
				}

				enumValue = &EnumValue{
					Name:  value.Name.Name,
					Value: rawValue,
					Type:  valueType,
				}
			} else {
				// 否则使用上一个值 + 1
				lastValue++
				constValue = &vast.BasicLit{
					Kind:  vtok.INTEGER,
					Value: fmt.Sprintf("%d", lastValue),
				}

				// 创建枚举值对象
				enumValue = &EnumValue{
					Name:  value.Name.Name,
					Value: fmt.Sprintf("%d", lastValue),
					Type:  "number",
				}
			}

			enumValues = append(enumValues, enumValue) // 收集枚举值
		}

	case *hast.ADTEnumBody:
		// ADT 枚举：enum Result { Success(num), Error(str) }
		for i, variant := range body.Variants {
			// 创建枚举值对象
			enumValue := &EnumValue{
				Name:  variant.Name.Name,
				Value: fmt.Sprintf("%d", i),
				Type:  "number",
			}
			enumValues = append(enumValues, enumValue) // 收集枚举值
		}
	}

	// 将枚举信息添加到符号表
	v.currentModule.Symbols.AddEnum(enumName, enumValues)

	// 不生成任何语句，只收集符号表信息
	return nil
}

func (v *VBScriptTranspiler) escapeVBScriptString(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}

func (v *VBScriptTranspiler) canResolveFunction(fun *vast.Ident) bool {
	log.Debugf("[canResolveFunction] 查找函数: %s", fun.Name)
	log.Debugf("[canResolveFunction] 当前模块: %s", v.currentModule.Path)
	log.Debugf("[canResolveFunction] 当前模块符号表函数: %v", getMapKeys(v.currentModule.Symbols.functions))
	log.Debugf("[canResolveFunction] 当前模块依赖: %v", v.currentModule.Dependencies)

	// 1. 查当前模块
	if v.currentModule.Symbols.HasFunction(fun.Name) {
		log.Debugf("[canResolveFunction] 在当前模块符号表找到: %s", fun.Name)
		return true
	}

	// 2. 递归查所有依赖
	visited := make(map[string]bool)
	found := v.findFunctionInDependencies(v.currentModule, fun.Name, visited)
	log.Debugf("[canResolveFunction] 递归依赖查找结果: %v", found)

	// 3. 查 moduleManager.externalVBS
	if _, ok := v.moduleManager.externalVBS[fun.Name]; ok {
		log.Debugf("[canResolveFunction] 在 moduleManager.externalVBS 找到: %s", fun.Name)
		return true
	}

	return found
}

// findFunctionInDependencies 递归查找模块及其所有依赖中是否存在指定函数
func (v *VBScriptTranspiler) findFunctionInDependencies(module *Module, funcName string, visited map[string]bool) bool {
	// 防止循环依赖
	if visited[module.Path] {
		return false
	}
	visited[module.Path] = true

	// 查当前模块
	if module.Symbols.HasFunction(funcName) {
		return true
	}

	// 递归查所有依赖
	for _, dep := range module.Dependencies {
		if depModule, exists := v.modules[dep]; exists {
			if v.findFunctionInDependencies(depModule, funcName, visited) {
				return true
			}
		}
	}

	return false
}

// 辅助函数：打印 map 的所有 key
func getMapKeys(m map[string]hast.Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// findFunctionInModuleDependencies 递归查找指定模块及其依赖中是否存在指定函数
func (v *VBScriptTranspiler) findFunctionInModuleDependencies(module *Module, targetModuleName, funcName string, visited map[string]bool) bool {
	// 防止循环依赖
	if visited[module.Path] {
		return false
	}
	visited[module.Path] = true

	// 检查当前模块是否是目标模块
	if strings.HasSuffix(module.Path, targetModuleName) || module.Path == targetModuleName {
		if module.Symbols.HasFunction(funcName) {
			return true
		}
	}

	// 检查内置模块
	for _, builtinModule := range v.builtinModules {
		if builtinModule.Symbols.HasFunction(funcName) {
			return true
		}
	}

	// 递归检查所有依赖
	for _, dep := range module.Dependencies {
		if depModule, exists := v.modules[dep]; exists {
			if v.findFunctionInModuleDependencies(depModule, targetModuleName, funcName, visited) {
				return true
			}
		}
	}

	return false
}

func (v *VBScriptTranspiler) ConvertCallExpr(node *hast.CallExpr) vast.Node {
	convertedFun := v.Convert(node.Fun)

	// Handle different types of function expressions
	switch fun := convertedFun.(type) {
	case *vast.Ident:
		// Simple function name
		if v.currentModule.Symbols.HasClass(fun.Name) {
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
				recv = append(recv, v.Convert(r).(vast.Expr).String())
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
		if moduleIdent, ok := fun.X.(*vast.Ident); ok {
			if funcIdent, ok := fun.Sel.(*vast.Ident); ok {
				// Check if this is a module access (e.g., utils.Calculate)
				// vs object method call (e.g., $p2.greet())
				visited := make(map[string]bool)
				if v.findFunctionInModuleDependencies(v.currentModule, moduleIdent.Name, funcIdent.Name, visited) {
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

func (v *VBScriptTranspiler) ConvertCmdExpr(node *hast.CmdExpr) vast.Node {
	convertedFun := v.Convert(node.Cmd)

	// Handle different types of function expressions
	switch fun := convertedFun.(type) {
	case *vast.Ident:
		// Simple function name
		if v.currentModule.Symbols.HasClass(fun.Name) {
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
				arg := v.Convert(r).(vast.Expr).String()
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
		if moduleIdent, ok := fun.X.(*vast.Ident); ok {
			if funcIdent, ok := fun.Sel.(*vast.Ident); ok {
				// Check if this is a module access (e.g., utils.Calculate)
				// vs object method call (e.g., $p2.greet())
				visited := make(map[string]bool)
				if v.findFunctionInModuleDependencies(v.currentModule, moduleIdent.Name, funcIdent.Name, visited) {
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

func (v *VBScriptTranspiler) ConvertFuncDecl(node *hast.FuncDecl) vast.Node {
	var ret *vast.FuncDecl = &vast.FuncDecl{}
	v.currentModule.Symbols.PushScope(FunctionScope, v.currentModule, node.Name.Name)
	defer v.currentModule.Symbols.PopScope()

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
	v.currentModule.Symbols.AddFunction(ret.Name.Name, node)

	return ret
}

func (v *VBScriptTranspiler) ConvertParameter(node *hast.Parameter) vast.Node {
	return &vast.Field{
		Name: v.Convert(node.Name).(*vast.Ident),
		Tok:  vtok.BYVAL,
	}
}

func (v *VBScriptTranspiler) ConvertBlock(node *hast.BlockStmt) vast.Node {
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

func (v *VBScriptTranspiler) ConvertIncDecExpr(node *hast.IncDecExpr) vast.Node {
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

func (v *VBScriptTranspiler) ConvertClassDecl(node *hast.ClassDecl) vast.Node {
	// 转换类名
	className := v.Convert(node.Name).(*vast.Ident)

	// 转换修饰符
	var mod vtok.Token
	var modPos vtok.Pos
	var isPublic bool
	if node.Pub.IsValid() {
		mod = vtok.PUBLIC
		modPos = vtok.Pos(node.Pub)
		isPublic = true
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

	// 将类信息添加到符号表
	v.currentModule.Symbols.AddClass(className.Name, &ClassSymbol{
		Name:   className.Name,
		Fields: fieldNames,
		Public: isPublic,
	})

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

func (v *VBScriptTranspiler) ConvertFile(node *hast.File) vast.Node {
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

	if v.currentModule.Path == v.opts.Main {
		for _, sym := range v.undefinedSymbols.Items() {
			stmts = append([]vast.Stmt{
				&vast.ExprStmt{
					X: &vast.BasicLit{
						Value: v.moduleManager.externalVBS[sym],
					},
				},
			}, stmts...)
			v.undefinedSymbols.Remove(sym)
		}
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

func (v *VBScriptTranspiler) ConvertForStmt(node *hast.ForStmt) vast.Node {
	v.currentModule.Symbols.PushScope(ForScope, v.currentModule)
	defer v.currentModule.Symbols.PopScope()

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
func (v *VBScriptTranspiler) isObjectExpression(expr vast.Expr) bool {
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

func (v *VBScriptTranspiler) ConvertAssignStmt(node *hast.AssignStmt) vast.Node {
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
			if v.currentModule.Symbols.HasClass(ident.Name) {
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

func (v *VBScriptTranspiler) Emit(n ...vast.Stmt) {
	v.buffer = append(v.buffer, n...)
}

func (v *VBScriptTranspiler) Flush() []vast.Stmt {
	stmts := v.buffer
	v.buffer = nil
	return stmts
}

func (v *VBScriptTranspiler) needsDimDeclaration(expr vast.Expr) bool {
	// 检查当前作用域
	if scope := v.currentModule.Symbols.CurrentScope(); scope != nil {
		// 在循环作用域中不需要Dim声明
		if scope.Type == ForScope {
			return false
		}
	}

	// 检查变量是否已经声明过
	if ident, ok := expr.(*vast.Ident); ok {
		if v.currentModule.Symbols.HasVariable(ident.Name) {
			return false
		}
		// 标记为已声明
		v.currentModule.Symbols.AddVariable(ident.Name, nil)
		return true
	}

	return false
}

// convertConstructorAssignment 处理构造函数调用赋值
// 例如：let p2 = Person("Jerry", 20) 转换为：
// Dim p2
// Set p2 = New Person
// p2.name = "Jerry"
// p2.age = 20
func (v *VBScriptTranspiler) convertConstructorAssignment(lhs vast.Expr, className string, args []hast.Expr) vast.Node {
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
func (v *VBScriptTranspiler) getClassFieldNames(className string) []string {
	if class, exists := v.currentModule.Symbols.GetClass(className); exists {
		return class.Fields
	}
	return []string{}
}

// ConvertMatchStmt 将 Hulo 的 match 语句转换为 VBScript 的 if-elseif-else 语句
func (v *VBScriptTranspiler) ConvertMatchStmt(node *hast.MatchStmt) vast.Node {
	v.currentModule.Symbols.PushScope(BlockScope, v.currentModule)
	defer v.currentModule.Symbols.PopScope()

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

func (v *VBScriptTranspiler) convertMatchCondition(matchExpr vast.Expr, cond hast.Expr) vast.Expr {
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
func (v *VBScriptTranspiler) ConvertArrayLiteral(node *hast.ArrayLiteralExpr) vast.Node {
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
func (v *VBScriptTranspiler) ConvertObjectLiteral(node *hast.ObjectLiteralExpr) vast.Node {
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
func (v *VBScriptTranspiler) ConvertForeachStmt(node *hast.ForeachStmt) vast.Node {
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
