// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hulo-lang/hulo/internal/build"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	vtok "github.com/hulo-lang/hulo/syntax/vbs/token"
)

func TranspileToVBScript(opts *config.VBScriptOptions, node hast.Node) (vast.Node, error) {
	tr := &VBScriptTranspiler{
		opts:            opts,
		scopeStack:      container.NewArrayStack[ScopeType](),
		declaredVars:    container.NewMapSet[string](),
		declaredFuncs:   container.NewMapSet[string](),
		enableShell:     false,
		declStmts:       container.NewArrayStack[vast.Stmt](),
		declaredClasses: container.NewMapSet[string](),
		classSymbols:    NewClassSymbolTable(),
		enumSymbols:     NewEnumSymbolTable(),
		declaredConsts:  container.NewMapSet[string](),
	}
	return tr.Convert(node), nil
}

type VBScriptTranspiler struct {
	opts *config.VBScriptOptions

	buffer []vast.Stmt

	scopeStack      container.Stack[ScopeType]
	declaredVars    container.Set[string]
	declaredFuncs   container.Set[string]
	declaredClasses container.Set[string]
	classSymbols    *ClassSymbolTable
	enumSymbols     *EnumSymbolTable
	declaredConsts  container.Set[string] // 跟踪已声明的常量

	declStmts container.Stack[vast.Stmt]

	enableShell bool

	alloc *build.Allocator
}

type ScopeType int

const (
	GlobalScope ScopeType = iota
	ForScope
	MatchScope
	IfScope
)

// ClassSymbol 表示类的符号表信息
type ClassSymbol struct {
	Name   string   // 类名
	Fields []string // 字段名列表
	Public bool     // 是否为公共类
}

// ClassSymbolTable 类符号表
type ClassSymbolTable struct {
	classes map[string]*ClassSymbol
}

func NewClassSymbolTable() *ClassSymbolTable {
	return &ClassSymbolTable{
		classes: make(map[string]*ClassSymbol),
	}
}

func (cst *ClassSymbolTable) AddClass(name string, fields []string, public bool) {
	cst.classes[name] = &ClassSymbol{
		Name:   name,
		Fields: fields,
		Public: public,
	}
}

func (cst *ClassSymbolTable) GetClass(name string) (*ClassSymbol, bool) {
	class, exists := cst.classes[name]
	return class, exists
}

func (cst *ClassSymbolTable) HasClass(name string) bool {
	_, exists := cst.classes[name]
	return exists
}

// EnumValue 表示枚举值的信息
type EnumValue struct {
	Name  string // 枚举值名
	Value string // 实际值
	Type  string // 值类型：number, string, boolean
}

// EnumSymbol 表示枚举的符号表信息
type EnumSymbol struct {
	Name   string       // 枚举名
	Values []*EnumValue // 枚举值列表
}

// EnumSymbolTable 枚举符号表
type EnumSymbolTable struct {
	enums map[string]*EnumSymbol
}

func NewEnumSymbolTable() *EnumSymbolTable {
	return &EnumSymbolTable{
		enums: make(map[string]*EnumSymbol),
	}
}

func (est *EnumSymbolTable) AddEnum(name string, values []*EnumValue) {
	est.enums[name] = &EnumSymbol{
		Name:   name,
		Values: values,
	}
}

func (est *EnumSymbolTable) GetEnum(name string) (*EnumSymbol, bool) {
	enum, exists := est.enums[name]
	return enum, exists
}

func (est *EnumSymbolTable) HasEnum(name string) bool {
	_, exists := est.enums[name]
	return exists
}

func (est *EnumSymbolTable) HasEnumValue(enumName, valueName string) bool {
	if enum, exists := est.enums[enumName]; exists {
		for _, value := range enum.Values {
			if value.Name == valueName {
				return true
			}
		}
	}
	return false
}

func (v *VBScriptTranspiler) pushScope(scope ScopeType) {
	v.scopeStack.Push(scope)
}

func (v *VBScriptTranspiler) popScope() {
	v.scopeStack.Pop()
}

func (v *VBScriptTranspiler) currentScope() (ScopeType, bool) {
	return v.scopeStack.Peek()
}

func (v *VBScriptTranspiler) Convert(node hast.Node) vast.Node {
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
		return &vast.ExprStmt{
			X: v.Convert(node.X).(vast.Expr),
		}
	case *hast.UnsafeStmt:
		return v.ConvertUnsafeStmt(node)
	case *hast.ExternDecl:
		return v.ConvertExternDecl(node)
	case *hast.CallExpr:
		return v.ConvertCallExpr(node)
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
			v.declaredVars.Add(expr.Name)
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
	if v.enumSymbols.HasEnum(moduleName) && v.enumSymbols.HasEnumValue(moduleName, memberName) {
		// 检查是否是纯数字枚举，如果是则直接返回数值
		if v.isNumericEnum(moduleName) {
			// 对于数字枚举，直接返回数值
			return v.getEnumNumericValue(moduleName, memberName)
		} else {
			// 对于字符串/混合枚举，生成常量声明并返回常量名
			constName := fmt.Sprintf("%s_%s", moduleName, memberName)

			// 检查是否已经声明过这个常量
			if !v.declaredConsts.Contains(constName) {
				constValue := v.getEnumStringValue(moduleName, memberName)

				// 生成常量声明
				constStmt := &vast.ConstStmt{
					Lhs: &vast.Ident{Name: constName},
					Rhs: constValue,
				}
				v.Emit(constStmt)

				// 标记为已声明
				v.declaredConsts.Add(constName)
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
	if enum, exists := v.enumSymbols.GetEnum(enumName); exists {
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
	if enum, exists := v.enumSymbols.GetEnum(enumName); exists {
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
	if enum, exists := v.enumSymbols.GetEnum(enumName); exists {
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

		// TODO 去除后缀

		// 解析 import 指代的文件路径

		return &vast.ExprStmt{
			X: &vast.CmdExpr{
				Cmd: &vast.Ident{Name: "Import"},
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
		v.declaredFuncs.Add(fnName.Name)
		return nil
	case *hast.BlockStmt:
		for _, stmt := range n.List {
			switch stmt := stmt.(type) {
			case *hast.AssignStmt:
				lhs := v.Convert(stmt.Lhs).(vast.Expr)
				v.declaredVars.Add(lhs.String())
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

	stmt, ok := v.declStmts.Peek()
	if !ok {
		panic("return statement outside of function")
	}

	var fn *vast.FuncDecl
	if fn, ok = stmt.(*vast.FuncDecl); !ok {
		panic("return statement outside of function")
	}

	return &vast.AssignStmt{
		Lhs: &vast.Ident{Name: fn.Name.Name},
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
	v.enumSymbols.AddEnum(enumName, enumValues)

	// 不生成任何语句，只收集符号表信息
	return nil
}

type StringPart struct {
	Text       string
	IsVariable bool
	Expr       vast.Expr
}

func ParseStringInterpolation(s string) []StringPart {
	var parts []StringPart
	var current strings.Builder
	var i int

	for i < len(s) {
		if s[i] == '$' {
			// 保存当前文本部分
			if current.Len() > 0 {
				parts = append(parts, StringPart{
					Text:       current.String(),
					IsVariable: false,
				})
				current.Reset()
			}

			i++ // 跳过 $

			if i < len(s) && s[i] == '{' {
				// ${name} 形式
				i++ // 跳过 {
				var varName strings.Builder
				for i < len(s) && s[i] != '}' {
					varName.WriteByte(s[i])
					i++
				}
				if i < len(s) {
					i++ // 跳过 }
				}

				parts = append(parts, StringPart{
					Text:       "",
					IsVariable: true,
					Expr: &vast.Ident{
						Name: varName.String(),
					},
				})
			} else {
				// $name 形式
				var varName strings.Builder
				for i < len(s) && (s[i] == '_' || (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= '0' && s[i] <= '9')) {
					varName.WriteByte(s[i])
					i++
				}

				parts = append(parts, StringPart{
					Text:       "",
					IsVariable: true,
					Expr: &vast.Ident{
						Name: varName.String(),
					},
				})
			}
		} else {
			current.WriteByte(s[i])
			i++
		}
	}

	// 添加最后的文本部分
	if current.Len() > 0 {
		parts = append(parts, StringPart{
			Text:       current.String(),
			IsVariable: false,
		})
	}

	return parts
}

// escapeVBScriptString 转义VBScript字符串中的双引号
func (v *VBScriptTranspiler) escapeVBScriptString(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}

func (v *VBScriptTranspiler) ConvertCallExpr(node *hast.CallExpr) vast.Node {
	convertedFun := v.Convert(node.Fun)

	// Handle different types of function expressions
	switch fun := convertedFun.(type) {
	case *vast.Ident:
		// Simple function name
		if v.declaredClasses.Contains(fun.Name) {
			return &vast.CmdExpr{
				Cmd: &vast.Ident{Name: "New"},
				Recv: []vast.Expr{
					&vast.Ident{Name: fun.Name},
				},
			}
		}

		if !v.declaredFuncs.Contains(fun.Name) {
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
	case *vast.SelectorExpr:
		// Handle $p.greet() case - the RefExpr was converted to SelectorExpr
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

func (v *VBScriptTranspiler) ConvertFuncDecl(node *hast.FuncDecl) vast.Node {
	var ret *vast.FuncDecl = &vast.FuncDecl{}
	v.declStmts.Push(ret)
	defer v.declStmts.Pop()

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

	if v.declaredFuncs.Contains(ret.Name.Name) {
		panic(fmt.Sprintf("function %s already declared", ret.Name.Name))
	}
	v.declaredFuncs.Add(ret.Name.Name)

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
		stmts = append(stmts, stmt.(vast.Stmt))
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
	v.declaredClasses.Add(className.Name)

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
	v.classSymbols.AddClass(className.Name, fieldNames, isPublic)

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
	v.pushScope(ForScope)
	defer v.popScope()

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
			if v.declaredClasses.Contains(ident.Name) {
				// 这是构造函数调用赋值，需要特殊处理
				return v.convertConstructorAssignment(lhs, ident.Name, callExpr.Recv)
			}
		}
	}

	if v.needsDimDeclaration(lhs) {
		v.Emit(&vast.DimDecl{List: []vast.Expr{lhs}})
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
	if scope, ok := v.currentScope(); ok {
		// 在循环作用域中不需要Dim声明
		if scope == ForScope {
			return false
		}
	}

	// 检查变量是否已经声明过
	if ident, ok := expr.(*vast.Ident); ok {
		if v.declaredVars.Contains(ident.Name) {
			return false
		}
		// 标记为已声明
		v.declaredVars.Add(ident.Name)
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
	if class, exists := v.classSymbols.GetClass(className); exists {
		return class.Fields
	}
	return []string{}
}

// ConvertMatchStmt 将 Hulo 的 match 语句转换为 VBScript 的 if-elseif-else 语句
func (v *VBScriptTranspiler) ConvertMatchStmt(node *hast.MatchStmt) vast.Node {
	v.pushScope(MatchScope)
	defer v.popScope()

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

// convertMatchCondition 转换 match 语句的条件表达式
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
