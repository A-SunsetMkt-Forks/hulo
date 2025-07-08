// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"
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
		opts:          opts,
		scopeStack:    container.NewArrayStack[ScopeType](),
		declaredVars:  container.NewMapSet[string](),
		declaredFuncs: container.NewMapSet[string](),
		enableShell:   false,
		declStmts:     container.NewArrayStack[vast.Stmt](),
	}
	return tr.Convert(node), nil
}

type VBScriptTranspiler struct {
	opts *config.VBScriptOptions

	buffer []vast.Stmt

	scopeStack    container.Stack[ScopeType]
	declaredVars  container.Set[string]
	declaredFuncs container.Set[string]

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
		return &vast.Ident{
			Name: node.X.String(),
		}
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
	default:
		fmt.Printf("%T\n", node)
	}
	return nil
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
	fun := v.Convert(node.Fun).(*vast.Ident)

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

	// 转换修饰符
	var mod vtok.Token
	var modPos vtok.Pos
	if node.Pub.IsValid() {
		mod = vtok.PUBLIC
		modPos = vtok.Pos(node.Pub)
	}

	// 转换类体中的语句
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

func (v *VBScriptTranspiler) ConvertFile(node *hast.File) vast.Node {
	docs := make([]*vast.CommentGroup, len(node.Docs))
	for i, d := range node.Docs {
		docs[i] = v.Convert(d).(*vast.CommentGroup)
	}

	stmts := []vast.Stmt{}
	for _, s := range node.Stmts {
		stmt := v.Convert(s)

		stmts = append(stmts, v.Flush()...)

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
