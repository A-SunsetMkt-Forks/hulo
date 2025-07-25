// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"fmt"
	"math/big"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/core"
	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Interpreter struct {
	debugger *Debugger

	// 当前作用域
	env *Environment

	// 模块映射
	modules map[string]*Environment

	currentModule core.Module

	importedModules map[string]core.Module
	defaultModule   core.Module

	fs vfs.VFS

	cvt object.Converter
}

func NewInterpreter(env *Environment) *Interpreter {
	return &Interpreter{
		env: env,
		cvt: &object.ASTConverter{},
	}
}

func (interp *Interpreter) shouldBreak(node ast.Node) bool {
	// pos := node.Pos()
	if bp, ok := interp.debugger.breakpoints["file"][1]; ok {
		fmt.Println(bp, "如果是条件断点，评估条件")
	}
	return false
}

func (interp *Interpreter) Eval(node ast.Node) ast.Node {
	log.Debugf("enter %T", node)
	defer log.Debugf("exit %T", node)
	switch node := node.(type) {
	/// Static World
	case *ast.File:
		return interp.rebuildFile(node)
	case *ast.IfStmt:
		return interp.rebuildIfStmt(node)
	case *ast.AssignStmt:
		return interp.rebuildAssignStmt(node)
	case *ast.ExprStmt:
		return interp.rebuildExprStmt(node)
	case *ast.ReturnStmt:
		return interp.rebuildReturnStmt(node)
	case *ast.ForStmt:
		return interp.rebuildForStmt(node)
	case *ast.ForeachStmt:
		return interp.rebuildForeachStmt(node)
	case *ast.WhileStmt:
		return interp.rebuildWhileStmt(node)
	case *ast.DoWhileStmt:
		return interp.rebuildDoWhileStmt(node)
	case *ast.MatchStmt:
		return interp.rebuildMatchStmt(node)
	case *ast.FuncDecl:
		return interp.rebuildFuncDecl(node)
	case *ast.CommentGroup:
		return interp.rebuildCommentGroup(node)
	case *ast.UnsafeStmt:
		return interp.rebuildUnsafeStmt(node)
	case *ast.ExternDecl:
		return interp.rebuildExternDecl(node)
	case *ast.CallExpr:
		return interp.rebuildCallExpr(node)
	case *ast.CmdExpr:
		return interp.rebuildCmdExpr(node)
	case *ast.Ident:
		return interp.rebuildIdent(node)
	case *ast.NumericLiteral:
		return interp.rebuildNumericLiteral(node)
	case *ast.StringLiteral:
		return interp.rebuildStringLiteral(node)
	case *ast.TrueLiteral:
		return interp.rebuildTrueLiteral(node)
	case *ast.FalseLiteral:
		return interp.rebuildFalseLiteral(node)
	case *ast.NullLiteral:
		return interp.rebuildNullLiteral(node)
	case *ast.ArrayLiteralExpr:
		return interp.rebuildArrayLiteralExpr(node)
	case *ast.ObjectLiteralExpr:
		return interp.rebuildObjectLiteralExpr(node)
	case *ast.RefExpr:
		return interp.rebuildRefExpr(node)
	case *ast.IncDecExpr:
		return interp.rebuildIncDecExpr(node)
	case *ast.BinaryExpr:
		return interp.rebuildBinaryExpr(node)
	case *ast.SelectExpr:
		return interp.rebuildSelectExpr(node)
	case *ast.TypeDecl:
		return interp.rebuildTypeDecl(node)
	case *ast.BlockStmt:
		return interp.rebuildBlockStmt(node)
	case *ast.Parameter:
		return interp.rebuildParameter(node)
	case *ast.ClassDecl:
		return interp.rebuildClassDecl(node)
	case *ast.DeclareDecl:
		return interp.rebuildDeclareDecl(node)
	case *ast.EnumDecl:
		return interp.rebuildEnumDecl(node)
	case *ast.ModAccessExpr:
		return interp.rebuildModAccessExpr(node)
	/// Dynamic World

	case *ast.Import:

	case *ast.ComptimeStmt:
		if node.Cond != nil {
			evaluatedObject := interp.evalComptimeWhenStmt(node)
			return interp.object2Node(evaluatedObject)
		}
		evaluatedObject := interp.evalComptimeStmt(node.Body)
		return interp.object2Node(evaluatedObject)

	case *ast.ComptimeExpr:
		evaluatedObject := interp.evalComptimeExpr(node)
		return interp.object2Node(evaluatedObject)
	}
	return node
}

func (interp *Interpreter) rebuildForeachStmt(node *ast.ForeachStmt) ast.Node {
	return node
}

func (interp *Interpreter) rebuildTypeDecl(node *ast.TypeDecl) ast.Node {
	return node
}

func (interp *Interpreter) rebuildEnumDecl(node *ast.EnumDecl) ast.Node {
	var newBody ast.EnumBody
	switch body := node.Body.(type) {
	case *ast.BasicEnumBody:
		newValues := make([]*ast.EnumValue, len(body.Values))
		for i, v := range body.Values {
			newValues[i] = &ast.EnumValue{
				Docs:  v.Docs,
				Name:  v.Name,
				Value: interp.Eval(v.Value).(ast.Expr),
			}
		}
		newBody = &ast.BasicEnumBody{
			Lbrace: body.Lbrace,
			Values: newValues,
			Rbrace: body.Rbrace,
		}
	case *ast.ADTEnumBody:
		newValues := make([]*ast.EnumVariant, len(body.Variants))
		for i, v := range body.Variants {
			newValues[i] = &ast.EnumVariant{
				Docs:   v.Docs,
				Name:   v.Name,
				Fields: v.Fields,
			}
		}
		newMethods := make([]ast.Stmt, len(body.Methods))
		for i, m := range body.Methods {
			newMethods[i] = interp.Eval(m).(ast.Stmt)
		}
		newBody = &ast.ADTEnumBody{
			Lbrace:   body.Lbrace,
			Variants: newValues,
			Methods:  newMethods,
			Rbrace:   body.Rbrace,
		}
	case *ast.AssociatedEnumBody:
		newFields := make([]*ast.Field, len(body.Fields.List))
		for i, f := range body.Fields.List {
			newFields[i] = &ast.Field{
				Docs:      f.Docs,
				Modifiers: f.Modifiers,
				Name:      f.Name,
				Type:      interp.Eval(f.Type).(ast.Expr),
				Value:     interp.Eval(f.Value).(ast.Expr),
			}
		}
		newValues := make([]*ast.EnumValue, len(body.Values))
		for i, v := range body.Values {
			newValues[i] = &ast.EnumValue{
				Docs:  v.Docs,
				Name:  v.Name,
				Value: interp.Eval(v.Value).(ast.Expr),
			}
		}
		newMethods := make([]ast.Stmt, len(body.Methods))
		for i, m := range body.Methods {
			newMethods[i] = interp.Eval(m).(ast.Stmt)
		}
		newBody = &ast.AssociatedEnumBody{
			Lbrace:  body.Lbrace,
			Fields:  &ast.FieldList{List: newFields},
			Values:  newValues,
			Methods: newMethods,
			Rbrace:  body.Rbrace,
		}
	}

	return &ast.EnumDecl{
		Docs:       node.Docs,
		Decs:       node.Decs,
		Modifiers:  node.Modifiers,
		Name:       node.Name,
		TypeParams: node.TypeParams,
		Body:       newBody,
	}
}

func (interp *Interpreter) rebuildModAccessExpr(node *ast.ModAccessExpr) ast.Node {
	return node
}

func (interp *Interpreter) rebuildParameter(node *ast.Parameter) ast.Node {
	return node
}

func (interp *Interpreter) rebuildClassDecl(node *ast.ClassDecl) ast.Node {
	newFields := make([]*ast.Field, len(node.Fields.List))
	for i, f := range node.Fields.List {
		newFields[i] = &ast.Field{
			Docs:      f.Docs,
			Modifiers: f.Modifiers,
			Name:      f.Name,
			Type:      interp.Eval(f.Type).(ast.Expr),
			Value:     interp.Eval(f.Value).(ast.Expr),
		}
	}

	newCtos := make([]*ast.ConstructorDecl, len(node.Ctors))
	for i, c := range node.Ctors {
		newInitFields := make([]ast.Expr, len(c.InitFields))
		for i, f := range c.InitFields {
			newInitFields[i] = interp.Eval(f).(ast.Expr)
		}
		newCtorBody := interp.Eval(c.Body).(*ast.BlockStmt)
		newCtos[i] = &ast.ConstructorDecl{
			Docs:       c.Docs,
			Decs:       c.Decs,
			Modifiers:  c.Modifiers,
			ClsName:    node.Name,
			Name:       c.Name,
			Body:       newCtorBody,
			InitFields: newInitFields,
		}
	}

	newMethods := make([]*ast.FuncDecl, len(node.Methods))
	for i, m := range node.Methods {
		newMethods[i] = interp.Eval(m).(*ast.FuncDecl)
	}

	return &ast.ClassDecl{
		Docs:      node.Docs,
		Decs:      node.Decs,
		Modifiers: node.Modifiers,
		Name:      node.Name,
		Fields:    &ast.FieldList{List: newFields},
		Ctors:     newCtos,
		Methods:   newMethods,
	}
}

func (interp *Interpreter) rebuildDeclareDecl(node *ast.DeclareDecl) ast.Node {
	return node
}

func (interp *Interpreter) rebuildBlockStmt(node *ast.BlockStmt) ast.Node {
	newList := make([]ast.Stmt, len(node.List))
	for i, s := range node.List {
		newList[i] = interp.Eval(s).(ast.Stmt)
	}
	return &ast.BlockStmt{List: newList}
}

func (interp *Interpreter) rebuildSelectExpr(node *ast.SelectExpr) ast.Node {
	newX := interp.Eval(node.X)
	newY := interp.Eval(node.Y)
	return &ast.SelectExpr{X: newX.(ast.Expr), Y: newY.(ast.Expr)}
}

func (interp *Interpreter) rebuildBinaryExpr(node *ast.BinaryExpr) ast.Node {
	newX := interp.Eval(node.X)
	newY := interp.Eval(node.Y)
	return &ast.BinaryExpr{X: newX.(ast.Expr), Op: node.Op, Y: newY.(ast.Expr)}
}

func (interp *Interpreter) rebuildRefExpr(node *ast.RefExpr) ast.Node {
	return node
}

func (interp *Interpreter) rebuildIncDecExpr(node *ast.IncDecExpr) ast.Node {
	return node
}

func (interp *Interpreter) rebuildObjectLiteralExpr(node *ast.ObjectLiteralExpr) ast.Node {
	newFields := make([]ast.Expr, len(node.Props))
	for i, f := range node.Props {
		kv := interp.Eval(f).(*ast.KeyValueExpr)
		newFields[i] = &ast.KeyValueExpr{
			Key:   interp.Eval(kv.Key).(ast.Expr),
			Value: interp.Eval(kv.Value).(ast.Expr),
		}
	}
	return &ast.ObjectLiteralExpr{Props: newFields}
}

func (interp *Interpreter) rebuildNullLiteral(node *ast.NullLiteral) ast.Node {
	return node
}

func (interp *Interpreter) rebuildArrayLiteralExpr(node *ast.ArrayLiteralExpr) ast.Node {
	return node
}

func (interp *Interpreter) rebuildTrueLiteral(node *ast.TrueLiteral) ast.Node {
	return node
}

func (interp *Interpreter) rebuildFalseLiteral(node *ast.FalseLiteral) ast.Node {
	return node
}

func (interp *Interpreter) rebuildNumericLiteral(node *ast.NumericLiteral) ast.Node {
	return node
}

func (interp *Interpreter) rebuildStringLiteral(node *ast.StringLiteral) ast.Node {
	return node
}

func (interp *Interpreter) rebuildIdent(node *ast.Ident) ast.Node {
	return node
}

func (interp *Interpreter) rebuildCmdExpr(node *ast.CmdExpr) ast.Node {
	newCmd := interp.Eval(node.Cmd)
	newArgs := make([]ast.Expr, len(node.Args))
	for i, a := range node.Args {
		newArgs[i] = interp.Eval(a).(ast.Expr)
	}
	return &ast.CmdExpr{Cmd: newCmd.(ast.Expr), Args: newArgs}
}

func (interp *Interpreter) rebuildCallExpr(node *ast.CallExpr) ast.Node {
	newFun := interp.Eval(node.Fun)
	newRecv := make([]ast.Expr, len(node.Recv))
	for i, r := range node.Recv {
		newRecv[i] = interp.Eval(r).(ast.Expr)
	}
	return &ast.CallExpr{Fun: newFun.(ast.Expr), Recv: newRecv}
}

func (interp *Interpreter) rebuildUnsafeStmt(node *ast.UnsafeStmt) ast.Node {
	return node
}

func (interp *Interpreter) rebuildExternDecl(node *ast.ExternDecl) ast.Node {
	return node
}

func (interp *Interpreter) rebuildCommentGroup(node *ast.CommentGroup) ast.Node {
	return node
}

func (interp *Interpreter) rebuildForStmt(node *ast.ForStmt) ast.Node {
	newInit := interp.Eval(node.Init)
	newCond := interp.Eval(node.Cond)
	newPost := interp.Eval(node.Post)
	newBody := interp.Eval(node.Body)

	return &ast.ForStmt{
		Init: newInit.(ast.Stmt),
		Cond: newCond.(ast.Expr),
		Post: newPost.(ast.Expr),
		Body: newBody.(*ast.BlockStmt),
	}
}

func (interp *Interpreter) rebuildWhileStmt(node *ast.WhileStmt) ast.Node {
	newCond := interp.Eval(node.Cond)
	newBody := interp.Eval(node.Body)
	return &ast.WhileStmt{
		Cond: newCond.(ast.Expr),
		Body: newBody.(*ast.BlockStmt),
	}
}

func (interp *Interpreter) rebuildDoWhileStmt(node *ast.DoWhileStmt) ast.Node {
	newBody := interp.Eval(node.Body)
	newCond := interp.Eval(node.Cond)
	return &ast.DoWhileStmt{
		Body: newBody.(*ast.BlockStmt),
		Cond: newCond.(ast.Expr),
	}
}

func (interp *Interpreter) rebuildMatchStmt(node *ast.MatchStmt) ast.Node {
	newExpr := interp.Eval(node.Expr)
	newCases := make([]*ast.CaseClause, len(node.Cases))
	for i, c := range node.Cases {
		newCases[i] = &ast.CaseClause{
			Cond: interp.Eval(c.Cond).(ast.Expr),
			Body: interp.Eval(c.Body).(*ast.BlockStmt),
		}
	}
	var newDefault *ast.CaseClause
	if node.Default != nil {
		newDefault = &ast.CaseClause{
			Cond: nil,
			Body: interp.Eval(node.Default.Body).(*ast.BlockStmt),
		}
	}
	return &ast.MatchStmt{
		Expr:    newExpr.(ast.Expr),
		Cases:   newCases,
		Default: newDefault,
	}
}

func (interp *Interpreter) rebuildFuncDecl(node *ast.FuncDecl) ast.Node {
	newRecv := make([]ast.Expr, len(node.Recv))
	for i, r := range node.Recv {
		newRecv[i] = interp.Eval(r).(ast.Expr)
	}
	newBody := interp.Eval(node.Body).(*ast.BlockStmt)
	return &ast.FuncDecl{
		Docs:      node.Docs,
		Modifiers: node.Modifiers,
		Recv:      newRecv,
		Name:      node.Name,
		Type:      node.Type,
		Body:      newBody,
	}
}

func (interp *Interpreter) rebuildReturnStmt(node *ast.ReturnStmt) ast.Node {
	newX := interp.Eval(node.X)
	return &ast.ReturnStmt{
		X: newX.(ast.Expr),
	}
}

func (interp *Interpreter) rebuildExprStmt(node *ast.ExprStmt) ast.Node {
	newX := interp.Eval(node.X)
	if v, ok := newX.(ast.Stmt); ok {
		return v
	}
	return &ast.ExprStmt{
		X: newX.(ast.Expr),
	}
}

func (interp *Interpreter) rebuildAssignStmt(node *ast.AssignStmt) ast.Node {
	newLhs := interp.Eval(node.Lhs)
	newRhs := interp.Eval(node.Rhs)

	// comptime {1+2} 的时候返回的是 ExprStmt
	if es, ok := newRhs.(*ast.ExprStmt); ok {
		newRhs = es.X
	}
	return &ast.AssignStmt{
		Scope: node.Scope,
		Lhs:   newLhs.(ast.Expr),
		Tok:   node.Tok,
		Rhs:   newRhs.(ast.Expr),
	}
}

func (interp *Interpreter) rebuildIfStmt(node *ast.IfStmt) ast.Node {
	newCond := interp.Eval(node.Cond)
	newBody := interp.Eval(node.Body)
	var newElse ast.Stmt
	if node.Else != nil {
		newElse = interp.Eval(node.Else).(ast.Stmt)
	}

	return &ast.IfStmt{
		If:   node.If,
		Cond: newCond.(ast.Expr),
		Body: newBody.(*ast.BlockStmt),
		Else: newElse,
	}
}

func (interp *Interpreter) rebuildFile(node *ast.File) ast.Node {
	file := &ast.File{
		Docs:  node.Docs,
		Name:  node.Name,
		Decls: node.Decls,
	}
	for _, stmt := range node.Stmts {
		file.Stmts = append(file.Stmts, interp.Eval(stmt).(ast.Stmt))
	}
	return file
}

func (interp *Interpreter) evalComptimeWhenStmt(node *ast.ComptimeStmt) object.Value {
	evaluatedObject := interp.evalComptimeExpr(node.Cond)
	if evaluatedObject == object.TRUE {
		return interp.evalComptimeStmt(node.Body)
	}
	if node.Else != nil {
		if elseStmt, ok := node.Else.(*ast.ComptimeStmt); ok {
			if elseStmt.Cond != nil {
				return interp.evalComptimeWhenStmt(elseStmt)
			}
			return interp.evalComptimeStmt(elseStmt.Body)
		}
	}
	return &object.ErrorValue{Value: "comptime when statement is not true"}
}

type WarpValue struct {
	AST ast.Node
}

func (w *WarpValue) Type() object.Type {
	return nil
}

func (w *WarpValue) Text() string {
	return "ast.Node"
}

func (w *WarpValue) Interface() any {
	return w.AST
}

func (interp *Interpreter) evalComptimeStmt(node ast.Stmt) object.Value {
	log.Debugf("enter %T", node)
	defer log.Debugf("exit %T", node)
	switch node := node.(type) {
	case *ast.BlockStmt:
		// 默认最后一行为给编译器的值
		if len(node.List) == 1 {
			return &WarpValue{AST: node.List[0]}
		}
		for _, stmt := range node.List {
			obj := interp.evalComptimeStmt(stmt)
			if returnValue, ok := obj.(*ReturnValue); ok {
				log.Debugf("return value: %T", returnValue.Value)
				return returnValue
			}
		}

	case *ast.Import:
		// return interp.executeImport(node)
	case *ast.ExprStmt:
		return interp.evalComptimeExpr(node.X)
	case *ast.AssignStmt:
		return interp.evalAssignStmt(node)
	case *ast.TypeDecl:
		return interp.evalTypeDecl(node)
	case *ast.ReturnStmt:
		return interp.evalReturnStmt(node)
	case *ast.IfStmt:
		return interp.evalIfStmt(node)
	case *ast.ForStmt:
		return interp.evalForStmt(node)
	case *ast.WhileStmt:
		return interp.evalWhileStmt(node)
	case *ast.DoWhileStmt:
		return interp.evalDoWhileStmt(node)
	case *ast.MatchStmt:
		return interp.evalMatchStmt(node)
	case *ast.FuncDecl:
		return interp.evalFuncDecl(node)
	default:
		panic("unknown comptime statement:" + fmt.Sprintf("%T", node))
	}
	return nil
}

func (interp *Interpreter) evalFuncDecl(node *ast.FuncDecl) object.Value {
	name := node.Name.Name

	builder := object.NewFunctionBuilder(name)

	// 参数
	for _, param := range node.Recv {
		var paramName string
		var paramType object.Type
		if p, ok := param.(*ast.Parameter); ok {
			paramName = p.Name.Name
			var err error
			paramType, err = interp.cvt.ConvertType(p.Type)
			if err != nil {
				paramType = object.GetAnyType()
			}
		} else if id, ok := param.(*ast.Ident); ok {
			paramName = id.Name
			paramType = object.GetAnyType()
		} else {
			paramName = "param"
			paramType = object.GetAnyType()
		}
		builder = builder.WithParameter(paramName, paramType)
	}

	// 返回类型
	retType := object.GetAnyType()
	if node.Type != nil {
		if ft, ok := node.Type.(*ast.FunctionType); ok {
			if ft.RetVal != nil {
				if t, err := interp.cvt.ConvertType(ft.RetVal); err == nil {
					retType = t
				}
			}
		}
	}
	builder = builder.WithReturnType(retType)
	builder = builder.WithBody(node)

	newFnType := builder.Build()

	fnValue, exists := interp.env.GetValue(name)
	if exists {
		if v, ok := fnValue.(*object.FunctionValue); ok {
			fnType := v.Type().(*object.FunctionType)
			sigs := newFnType.Signatures()
			for _, sig := range sigs {
				if sig != nil {
					fnType.AppendSignatures([]*object.FunctionSignature{sig})
				}
			}
			return nil
		}
	}
	interp.env.Set(name, object.NewFunctionValue(newFnType))
	return nil
}

func (interp *Interpreter) evalMatchStmt(node *ast.MatchStmt) object.Value {
	prevEnv := interp.env
	interp.env = interp.env.Fork()
	defer func() { interp.env = prevEnv }()

	lv := interp.evalComptimeExpr(node.Expr)
	for _, clause := range node.Cases {
		if clause.Cond != nil {
			rv := interp.evalComptimeExpr(clause.Cond)
			if rv.Text() == lv.Text() {
				return interp.evalComptimeStmt(clause.Body)
			}
		}
	}

	if node.Default != nil {
		return interp.evalComptimeStmt(node.Default.Body)
	}
	return nil
}

func (interp *Interpreter) evalDoWhileStmt(node *ast.DoWhileStmt) object.Value {
	prevEnv := interp.env
	interp.env = interp.env.Fork()
	defer func() { interp.env = prevEnv }()

	for {
		interp.evalComptimeStmt(node.Body)
		cond := interp.evalComptimeExpr(node.Cond)
		if cond == object.FALSE {
			break
		}
	}
	return nil
}

func (interp *Interpreter) evalWhileStmt(node *ast.WhileStmt) object.Value {
	prevEnv := interp.env
	interp.env = interp.env.Fork()
	defer func() { interp.env = prevEnv }()

	for {
		cond := interp.evalComptimeExpr(node.Cond)
		if cond == object.FALSE {
			break
		}
		interp.evalComptimeStmt(node.Body)
	}
	return nil
}

func (interp *Interpreter) evalForStmt(node *ast.ForStmt) object.Value {
	prevEnv := interp.env
	interp.env = interp.env.Fork()
	defer func() { interp.env = prevEnv }()

	if node.Init != nil {
		interp.evalComptimeStmt(node.Init)
	}
	for {
		cond := true
		if node.Cond != nil {
			condValue := interp.evalComptimeExpr(node.Cond)
			if boolVal, ok := condValue.(interface{ Interface() any }); ok {
				if v, ok := boolVal.Interface().(bool); ok {
					cond = v
				} else {
					cond = false
				}
			} else {
				cond = false
			}
		}
		if !cond {
			break
		}
		result := interp.evalComptimeStmt(node.Body)
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue
		}
		if node.Post != nil {
			interp.evalComptimeExpr(node.Post)
		}
	}
	return nil
}

func (interp *Interpreter) evalIfStmt(node *ast.IfStmt) object.Value {
	cond := interp.evalComptimeExpr(node.Cond)
	if cond == object.TRUE {
		return interp.evalComptimeStmt(node.Body)
	}
	for node.Else != nil {
		switch el := node.Else.(type) {
		case *ast.IfStmt:
			cond = interp.evalComptimeExpr(el.Cond)
			if cond == object.TRUE {
				return interp.evalComptimeStmt(el.Body)
			}
			node.Else = el.Else
		case *ast.BlockStmt:
			return interp.evalComptimeStmt(el)
		}
	}
	return nil
}

type ReturnValue struct {
	Value object.Value
}

func (r *ReturnValue) Type() object.Type {
	return r.Value.Type()
}

func (r *ReturnValue) Text() string {
	return r.Value.Text()
}

func (r *ReturnValue) Interface() any {
	return r.Value.Interface()
}

func (interp *Interpreter) evalReturnStmt(node *ast.ReturnStmt) object.Value {
	if node.X == nil {
		return &ReturnValue{Value: object.NULL}
	}
	return &ReturnValue{Value: interp.evalComptimeExpr(node.X)}
}

func (interp *Interpreter) evalTypeDecl(node *ast.TypeDecl) object.Value {
	typ, err := interp.cvt.ConvertType(node.Value)
	if err != nil {
		return &object.ErrorValue{Value: err.Error()}
	}

	fmt.Println(typ, node.Name, "TODO: 在类型表注册类型")
	return nil
}

func (interp *Interpreter) evalComptimeExpr(node ast.Expr) object.Value {
	switch node := node.(type) {
	case *ast.Ident:
		return interp.evalIdent(node)
	case *ast.SelectExpr:
		return interp.executeSelectExpr(node)
	case *ast.CallExpr:
		return interp.executeCallExpr(node)
	case *ast.CmdExpr:
		return interp.evalCmdExpr(node)
	case *ast.BinaryExpr:
		return interp.evalBinaryExpr(node)
	case *ast.RefExpr:
		return interp.evalIdent(node.X.(*ast.Ident))
	case *ast.IncDecExpr:
		return interp.evalIncDecExpr(node)
	default:
		value, err := interp.cvt.ConvertValue(node)
		if value == nil || err != nil {
			return &object.ErrorValue{Value: "failed to convert expression"}
		}
		return value
	}
}

func (interp *Interpreter) evalIncDecExpr(node *ast.IncDecExpr) object.Value {
	var name string
	if ident, ok := node.X.(*ast.Ident); ok {
		name = ident.Name
	} else if ref, ok := node.X.(*ast.RefExpr); ok {
		name = ref.X.(*ast.Ident).Name
	}

	value, ok := interp.env.GetValue(name)
	if !ok {
		return &object.ErrorValue{Value: "variable not found"}
	}
	if num, ok := value.(*object.NumberValue); ok {
		if node.Tok == token.INC {
			num.Value.Add(num.Value, big.NewFloat(1))
		} else {
			num.Value.Sub(num.Value, big.NewFloat(1))
		}
	}
	interp.env.Assign(name, value)
	return value
}

func (interp *Interpreter) evalCmdExpr(node *ast.CmdExpr) object.Value {
	fn := interp.evalComptimeExpr(node.Cmd)

	switch fn.Type().Kind() {
	case object.O_FUNC:
		if funcType, ok := fn.Type().(*object.FunctionType); ok {
			evaluator := NewFunctionEvaluator(interp)
			args := interp.executeExprList(node.Args)
			// TODO: 要判断错误
			result, err := funcType.Call(args, nil, evaluator)
			if err != nil {
				return &object.ErrorValue{Value: err.Error()}
			}
			return result
		}
	default:
		return &object.ErrorValue{Value: "unsupported function type"}
	}
	return nil
}

func (interp *Interpreter) executeCallExpr(node *ast.CallExpr) object.Value {
	fn := interp.evalComptimeExpr(node.Fun)

	switch fn.Type().Kind() {
	case object.O_FUNC:
		if funcType, ok := fn.Type().(*object.FunctionType); ok {
			evaluator := NewFunctionEvaluator(interp)
			args := interp.executeExprList(node.Recv)
			// TODO: 要判断错误
			result, err := funcType.Call(args, nil, evaluator)
			if err != nil {
				return &object.ErrorValue{Value: err.Error()}
			}
			return result
		}
	default:
		return &object.ErrorValue{Value: "unsupported function type"}
	}
	return nil
}

func (interp *Interpreter) executeExprList(exprs []ast.Expr) []object.Value {
	values := make([]object.Value, len(exprs))
	for i, expr := range exprs {
		values[i] = interp.evalComptimeExpr(expr)
		if values[i] == nil {
			return []object.Value{&object.ErrorValue{Value: "failed to evaluate expression"}}
		}
	}
	return values
}

func (interp *Interpreter) evalBinaryExpr(node *ast.BinaryExpr) object.Value {
	lhs := interp.evalComptimeExpr(node.X)
	rhs := interp.evalComptimeExpr(node.Y)

	switch {
	case lhs.Type().Kind() == object.O_NUM && rhs.Type().Kind() == object.O_NUM:
		return interp.executeBinaryExprNumber(lhs, rhs, node.Op)
	case lhs.Type().Kind() == object.O_STR && rhs.Type().Kind() == object.O_STR:
		return interp.executeBinaryExprString(lhs, rhs, node.Op)
	case lhs.Type().Kind() == object.O_BOOL && rhs.Type().Kind() == object.O_BOOL:
		return interp.executeBinaryExprBool(lhs, rhs, node.Op)
	case lhs.Type().Kind() != rhs.Type().Kind():
		return &object.ErrorValue{Value: "type mismatch"}
	}
	return &object.ErrorValue{Value: "unknown binary expression"}
}

func (interp *Interpreter) executeBinaryExprBool(lhs, rhs object.Value, op token.Token) object.Value {
	lv := lhs.Interface().(bool)
	rv := rhs.Interface().(bool)

	switch op {
	case token.AND:
		return interp.nativeBoolObject(lv && rv)
	case token.OR:
		return interp.nativeBoolObject(lv || rv)
	}
	return &object.ErrorValue{Value: "unknown binary expression"}
}

func (interp *Interpreter) executeBinaryExprNumber(lhs, rhs object.Value, op token.Token) object.Value {
	lv := lhs.Interface().(*big.Float)
	rv := rhs.Interface().(*big.Float)

	switch op {
	case token.PLUS:
		return &object.NumberValue{Value: new(big.Float).Add(lv, rv)}
	case token.MINUS:
		return &object.NumberValue{Value: new(big.Float).Sub(lv, rv)}
	case token.ASTERISK:
		return &object.NumberValue{Value: new(big.Float).Mul(lv, rv)}
	case token.SLASH:
		return &object.NumberValue{Value: new(big.Float).Quo(lv, rv)}
	// case token.MOD:
	// 	return &object.NumberValue{Value: new(big.Float).Mod(lv, rv)}
	case token.LT:
		return interp.nativeBoolObject(lv.Cmp(rv) < 0)
	case token.GT:
		return interp.nativeBoolObject(lv.Cmp(rv) > 0)
	case token.EQ:
		return interp.nativeBoolObject(lv.Cmp(rv) == 0)
	case token.NEQ:
		return interp.nativeBoolObject(lv.Cmp(rv) != 0)
	case token.LE:
		return interp.nativeBoolObject(lv.Cmp(rv) <= 0)
	case token.GE:
		return interp.nativeBoolObject(lv.Cmp(rv) >= 0)
	case token.AND:
	}
	return &object.ErrorValue{Value: "unknown binary expression"}
}

func (interp *Interpreter) executeBinaryExprString(lhs, rhs object.Value, op token.Token) object.Value {
	lv := lhs.Interface().(string)
	rv := rhs.Interface().(string)

	switch op {
	case token.PLUS:
		return &object.StringValue{Value: lv + rv}
	case token.EQ:
		return interp.nativeBoolObject(lv == rv)
	case token.NEQ:
		return interp.nativeBoolObject(lv != rv)
	case token.LT:
		return interp.nativeBoolObject(lv < rv)
	case token.GT:
		return interp.nativeBoolObject(lv > rv)
	case token.LE:
		return interp.nativeBoolObject(lv <= rv)
	case token.GE:
		return interp.nativeBoolObject(lv >= rv)
	}
	return &object.ErrorValue{Value: "unknown binary expression"}
}

func (interp *Interpreter) nativeBoolObject(v bool) object.Value {
	if v {
		return object.TRUE
	}
	return object.FALSE
}

func (interp *Interpreter) executeBasicLit(node *ast.BasicLit) object.Value {
	switch node.Kind {
	case token.NUM:
		o := &object.NumberValue{}
		o.Value.SetString(node.Value)
		return o
	case token.STR:
		return &object.StringValue{Value: node.Value}
	case token.TRUE:
		return object.TRUE
	case token.FALSE:
		return object.FALSE
	case token.NULL:
		return object.NULL
	}
	return &object.ErrorValue{Value: "unknown basic literal"}
}

func (interp *Interpreter) executeSelectExpr(node *ast.SelectExpr) object.Value {
	lhs := interp.evalComptimeExpr(node.X)

	switch lhs.Type().Kind() {
	case object.O_STR:
		return interp.executePackageSelector(lhs, node.Y)
	case object.O_LITERAL:
		return interp.executePackageSelector(lhs, node.Y)
	}
	return nil
}

func (interp *Interpreter) isPackageName(v object.Value) bool {
	return v.Type().Kind() == object.O_STR || v.Type().Kind() == object.O_LITERAL
}

func (interp *Interpreter) executePackageSelector(pkg object.Value, selector ast.Expr) object.Value {
	// interp.env.Get(pkgName)
	return nil
}

func (interp *Interpreter) object2Node(v object.Value) ast.Node {
	if warp, ok := v.(*WarpValue); ok {
		return warp.AST
	}
	if returnValue, ok := v.(*ReturnValue); ok {
		switch returnValue.Value.Type().Kind() {
		case object.O_NUM:
			return &ast.NumericLiteral{
				Value: returnValue.Value.Text(),
			}
		case object.O_STR:
			return &ast.StringLiteral{
				Value: returnValue.Value.Interface().(string),
			}
		case object.O_BOOL:
			if returnValue.Value.Interface().(bool) {
				return &ast.TrueLiteral{}
			}
			return &ast.FalseLiteral{}
		case object.O_NULL:
			return &ast.NullLiteral{}
		}
	}
	return nil
}

func Evaluate(ctx *Context, node ast.Node) object.Value {
	switch node := node.(type) {
	case *ast.File:
		return evalFile(ctx, node)
	case *ast.FuncDecl:
		return evalFuncDecl(ctx, node)
	case *ast.ModDecl:
	case *ast.SelectExpr:
		return evalSelectExpr(ctx, node)
	case *ast.Ident:
		// return evalIdent(ctx, node)
	case *ast.CallExpr:
		return evalCallExpr(ctx, node)
	case *ast.ExtensionDecl:
		return evalExtensionDecl(ctx, node)
	}
	return nil
}

func evalExtensionDecl(ctx *Context, node *ast.ExtensionDecl) object.Value {
	switch body := node.Body.(type) {
	case *ast.ExtensionClass:
		// 将类查出来 插入进去
		_ = body // TODO: use body.Body to process extension
		ctx.Get(node.Name.Name)
	}
	return nil
}

// x.y
func evalSelectExpr(ctx *Context, node *ast.SelectExpr) object.Value {
	// math.PI.to_str()
	// PI.to_str()
	// 检查 X 是什么东西
	x := Evaluate(ctx, node.X)
	x.Type() // 可以知道是什么鬼东西了这下
	// 1. 是包名 走 vfs 读取代码
	// 要判断包名是不是当前路径下面的文件

	// 可能是一个别名
	// import * as math from "net"
	// import "math"
	// import * from "math"

	// 要在这里遍历ast看看真正的文件名

	// 是直接引入
	// wd, err := os.Getwd()
	// if err != nil {
	// 	return nil
	// }
	// if target := filepath.Join(wd, x.Name()); ctx.os.Exist(target) {
	// 	// ctx.mem.Import(target)
	// }

	// 不是判断 是第三方还是标准库 要拿到 HULOPATH 这个变量 指出第三方库和标准库存储的父路径

	// ctx.mem.Import(filepath.Join("HULOPATH", x.Name()))

	// 引入后开始访问 y 表达式
	{
		// 拿库的上下文去访问？ 这样就能找到 Y 的定义
		// y := Evaluate(ctx, node.Y) // 如果 y 是函数要把自己传入进去吧？？？？

		// ctx.mem.Get(y.Name())
	}

	// 2. 是变量名
	// PI
	x.Type()

	{
		// y := Evaluate(ctx, node.Y) // 在去访问 to_Str()
		// ctx.mem.Get(y.Name())
	}

	// 3. 是类名

	// 4.
	return nil
}

func evalFile(ctx *Context, node *ast.File) object.Value {
	for _, decl := range node.Decls {
		Evaluate(ctx, decl)
	}
	return nil
}

func evalFuncDecl(ctx *Context, node *ast.FuncDecl) object.Value {
	return nil
}

func evalCallExpr(ctx *Context, node *ast.CallExpr) object.Value {
	// 需要先找到函数的定义
	// 拿到 to_str 的定义
	// function := Evaluate(ctx, node.Fun)
	// args := evalExpressions(ctx, node.Recv)
	// 要根据 args 类型 拿到method
	// TODO function.Type().Call(args...) 自动匹配合适的函数
	// function.Type().Method(0).Call(args...)
	// Call 的逻辑中，builtin直接执行，非builtin要执行语法树
	return nil
}

func evalExpressions(ctx *Context, exprs []ast.Expr) []object.Value {
	return nil
}

func (interp *Interpreter) evalIdent(node *ast.Ident) object.Value {
	if v, ok := interp.env.GetValue(node.Name); ok {
		return v
	}

	if builtin, ok := builtin[node.Name]; ok {
		return builtin
	}

	return &object.ErrorValue{Value: fmt.Sprintf("identifier %s not found", node.Name)}
}

func (interp *Interpreter) evalAssignStmt(node *ast.AssignStmt) object.Value {
	// 计算右值
	rhsValue := interp.evalComptimeExpr(node.Rhs)
	// TODO: hulo 支持 let a: null 这种语法，所以可能没值也合法
	if rhsValue == nil {
		return &object.ErrorValue{Value: "failed to evaluate right-hand side expression"}
	}

	// 根据左值类型处理
	switch lhs := node.Lhs.(type) {
	case *ast.Ident:
		return interp.handleIdentAssignment(lhs, rhsValue, node.Scope, node.Tok)
	case *ast.RefExpr:
		return interp.handleRefExprAssignment(lhs, rhsValue, node.Scope, node.Tok)
	default:
		return &object.ErrorValue{Value: "unsupported left-hand side expression type"}
	}
}

func (interp *Interpreter) handleRefExprAssignment(ref *ast.RefExpr, value object.Value, scope token.Token, assignTok token.Token) object.Value {
	name := ref.X.(*ast.Ident).Name

	var err error
	switch assignTok {
	case token.COLON_ASSIGN:
		err = interp.env.Declare(ref.X.(*ast.Ident).Name, cloneValue(value), token.LET)
	case token.ASSIGN:
		err = interp.env.Assign(name, value)
	case token.PLUS_ASSIGN:
		oldValue, ok := interp.env.GetValue(name)
		if !ok {
			return &object.ErrorValue{Value: "variable not found"}
		}
		if num, ok := oldValue.(*object.NumberValue); ok {
			num.Value.Add(num.Value, value.(*object.NumberValue).Value)
		}
	case token.MINUS_ASSIGN:
		oldValue, ok := interp.env.GetValue(name)
		if !ok {
			return &object.ErrorValue{Value: "variable not found"}
		}
		if num, ok := oldValue.(*object.NumberValue); ok {
			num.Value.Sub(num.Value, value.(*object.NumberValue).Value)
		}
	case token.ASTERISK_ASSIGN:
		oldValue, ok := interp.env.GetValue(name)
		if !ok {
			return &object.ErrorValue{Value: "variable not found"}
		}
		if num, ok := oldValue.(*object.NumberValue); ok {
			num.Value.Mul(num.Value, value.(*object.NumberValue).Value)
		}
	case token.SLASH_ASSIGN:
		oldValue, ok := interp.env.GetValue(name)
		if !ok {
			return &object.ErrorValue{Value: "variable not found"}
		}
		if num, ok := oldValue.(*object.NumberValue); ok {
			num.Value.Quo(num.Value, value.(*object.NumberValue).Value)
		}
	}

	if err != nil {
		return &object.ErrorValue{Value: err.Error()}
	}
	return nil
}

func (interp *Interpreter) handleIdentAssignment(ident *ast.Ident, value object.Value, scope token.Token, assignTok token.Token) object.Value {
	name := ident.Name

	switch assignTok {
	case token.ASSIGN:
		// 简单赋值
		if scope != token.ILLEGAL {
			// 这是变量声明
			if err := interp.env.Declare(name, cloneValue(value), scope); err != nil {
				return &object.ErrorValue{Value: err.Error()}
			}
		} else {
			// 这是变量重新赋值
			if err := interp.env.Assign(name, cloneValue(value)); err != nil {
				return &object.ErrorValue{Value: err.Error()}
			}
		}
		return value

	case token.COLON_ASSIGN:
		// := 声明并赋值（类似 Go）
		if err := interp.env.Declare(name, cloneValue(value), token.LET); err != nil {
			return &object.ErrorValue{Value: err.Error()}
		}
		return value

	default:
		return &object.ErrorValue{Value: "unsupported assignment operator"}
	}
}

// cloneValue 工具函数
func cloneValue(value object.Value) object.Value {
	switch v := value.(type) {
	case *object.NumberValue:
		return v.Clone()
	case *object.StringValue:
		return v.Clone()
	case *object.BoolValue:
		return v.Clone()
	default:
		return value // 其他类型暂时直接返回
	}
}

// func NewInterpreter() *Interpreter {
// 	interpreter := &Interpreter{
// 		env: NewEnvironment(),
// 	}

// 	return interpreter
// }

// // GetEnvironment 获取环境
// func (interp *Interpreter) GetEnvironment() *Environment {
// 	return interp.env
// }

// var stdlibs = []string{}

// func (interp *Interpreter) resolveModulePath(path string) (string, error) {
// 	// 1. 相对路径 (./math, ../utils)
// 	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
// 		return interp.resolveRelativePath(path)
// 	}

// 	// 2. 标准库路径 (math, io, net)
// 	if slices.Contains(stdlibs, path) {
// 		return interp.resolveStdLibPath(path)
// 	}

// 	// 3. 第三方库路径 (github.com/hulo-lang/hulo/stdlibs/math)
// 	if strings.Contains(path, "/") {
// 		return interp.resolveThirdPartyPath(path)
// 	}

// 	return interp.resolveLocalModule(path)
// }

// func (interp *Interpreter) resolveRelativePath(path string) (string, error) {
// 	currentDir := interp.currentModule.Path()

// 	absolutePath := interp.fs.Join(currentDir, path)

// 	if !interp.fs.Exists(absolutePath) {
// 		return "", fmt.Errorf("module %s not found", path)
// 	}

// 	return absolutePath, nil
// }

// func (interp *Interpreter) resolveStdLibPath(path string) (string, error) {
// 	huloPath := os.Getenv("HULO_PATH")

// 	target := interp.fs.Join(huloPath, path)

// 	if !interp.fs.Exists(target) {
// 		return "", fmt.Errorf("module %s not found", path)
// 	}

// 	return target, nil
// }

// func (interp *Interpreter) resolveThirdPartyPath(path string) (string, error) {
// 	huloPath := os.Getenv("HULO_MODULES")

// 	target := interp.fs.Join(huloPath, path)

// 	if !interp.fs.Exists(target) {
// 		return "", fmt.Errorf("module %s not found", path)
// 	}

// 	return target, nil
// }

// func (interp *Interpreter) resolveLocalModule(path string) (string, error) {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return "", err
// 	}

// 	target := interp.fs.Join(wd, path)

// 	if !interp.fs.Exists(target) {
// 		return "", fmt.Errorf("module %s not found", path)
// 	}

// 	return target, nil
// }

// type ModuleSystem interface {
// 	// 加载模块
// 	LoadModule(path string) (core.Module, error)

// 	// 解析符号（跨模块查找）
// 	ResolveSymbol(name string, currentModule string) (Symbol, error)

// 	// 获取模块
// 	GetModule(name string) (core.Module, bool)
// }

// // Symbol 表示一个符号
// type Symbol struct {
// 	Name   string
// 	Value  object.Value
// 	Type   object.Type
// 	Module string // 所属模块
// 	Kind   SymbolKind
// }

// type SymbolKind int

// const (
// 	SymbolValue SymbolKind = iota
// 	SymbolType
// 	SymbolFunction
// 	SymbolModule
// )

// func (interp *Interpreter) resolveSymbol(name string) (core.Value, error) {
// 	// 1. 首先在当前模块中查找
// 	if value, found := interp.currentModule.GetValue(name); found {
// 		return value, nil
// 	}

// 	// 2. 在导入的模块中查找
// 	for _, module := range interp.importedModules {
// 		if value, found := module.GetValue(name); found {
// 			return value, nil
// 		}
// 	}

// 	// 3. 在默认模块中查找（如 std 库）
// 	if value, found := interp.defaultModule.GetValue(name); found {
// 		return value, nil
// 	}

// 	// 4. 在全局符号表中查找
// 	if symbol, found := interp.moduleSystem.ResolveSymbol(name, interp.currentModule.Name()); found {
// 		return symbol.Value, nil
// 	}

// 	return nil, fmt.Errorf("symbol %s not found", name)
// }

// func (interp *Interpreter) getTypeModule(typ object.Type) string {
// 	// 1. 基本类型属于内置模块
// 	if interp.isBuiltinType(typ) {
// 		return "builtin"
// 	}

// 	// 2. 检查类型是否来自特定模块
// 	if moduleType, ok := typ.(*object.ModuleType); ok {
// 		return moduleType.ModuleName
// 	}

// 	// 3. 从类型注册表中查找
// 	if moduleName, found := interp.typeRegistry.GetTypeModule(typ.Name()); found {
// 		return moduleName
// 	}

// 	// 4. 默认返回当前模块
// 	return interp.currentModule.Name()
// }

// func (interp *Interpreter) executeSelectExpr(node *ast.SelectExpr) object.Value {
// 	// 1. 求值左操作数
// 	lhs := interp.executeComptimeStmt(node.X)

// 	// 2. 判断左操作数的类型
// 	switch {
// 	case interp.isModuleReference(lhs):
// 		// 这是一个模块引用 (math.PI)
// 		return interp.resolveModuleMember(lhs, node.Y)

// 	case interp.isPackageName(lhs):
// 		// 这是一个包名 (math.PI)
// 		return interp.resolvePackageMember(lhs, node.Y)

// 	default:
// 		// 这是一个普通的成员访问 (obj.field)
// 		return interp.resolveObjectMember(lhs, node.Y)
// 	}
// }

// func (interp *Interpreter) isModuleReference(v object.Value) bool {
// 	// 检查是否是模块值
// 	if moduleValue, ok := v.(*object.ModuleValue); ok {
// 		return moduleValue != nil
// 	}
// 	return false
// }

// func (interp *Interpreter) resolveModuleMember(moduleValue object.Value, selector ast.Expr) object.Value {
// 	module := moduleValue.(*object.ModuleValue).Module

// 	// 获取选择器名称
// 	var memberName string
// 	if ident, ok := selector.(*ast.Ident); ok {
// 		memberName = ident.Name
// 	} else {
// 		return &object.ErrorValue{Value: "invalid selector"}
// 	}

// 	// 从模块中查找成员
// 	if value, found := module.GetValue(memberName); found {
// 		return value
// 	}

// 	if typ, found := module.GetType(memberName); found {
// 		return &object.TypeValue{Type: typ}
// 	}

// 	return &object.ErrorValue{Value: fmt.Sprintf("member %s not found in module", memberName)}
// }

// func (interp *Interpreter) executeImport(node *ast.Import) object.Value {
// 	// 1. 解析导入语句
// 	var path string
// 	var alias string

// 	switch {
// 	case node.ImportAll != nil:
// 		// import * as math from "math"
// 		path = node.ImportAll.Path
// 		alias = node.ImportAll.Alias

// 	case node.ImportSingle != nil:
// 		// import "math" as math
// 		path = node.ImportSingle.Path
// 		alias = node.ImportSingle.Alias

// 	case node.ImportMulti != nil:
// 		// import { PI, E } from "math"
// 		path = node.ImportMulti.Path
// 		// 处理多个导入项
// 		for _, field := range node.ImportMulti.List {
// 			fieldName := field.Field
// 			fieldAlias := field.Alias
// 			if fieldAlias == "" {
// 				fieldAlias = fieldName
// 			}
// 			// 注册到当前模块的符号表
// 			interp.registerImportedSymbol(fieldAlias, path, fieldName)
// 		}
// 		return object.NULL
// 	}

// 	modulePath, err := interp.resolveModulePath(path)
// 	if err != nil {
// 		return &object.ErrorValue{Value: err.Error()}
// 	}

// 	// 2. 加载模块
// 	moduleFile, err := interp.fs.ReadFile(modulePath)
// 	if err != nil {
// 		return &object.ErrorValue{Value: err.Error()}
// 	}

// 	ast, err := parser.ParseSourceScript(moduleFile)
// 	if err != nil {
// 		return &object.ErrorValue{Value: err.Error()}
// 	}

// 	// 创建模块
// 	fmt.Println(ast, alias)
// 	// 3. 注册模块别名到当前环境
// 	// if alias != "" {
// 	// 	interp.env.Set(alias, &object.ModuleValue{Module: module})
// 	// }

// 	return object.NULL
// }

// // ExecuteMain 执行主模块
// func (interp *Interpreter) ExecuteMain(mainFile string) error {
// 	// 1. 解析主模块路径
// 	modulePath, err := interp.resolveModulePath(mainFile)
// 	if err != nil {
// 		return fmt.Errorf("failed to resolve main module path: %w", err)
// 	}

// 	// 2. 加载主模块
// 	mainModule, err := interp.LoadModule(modulePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to load main module: %w", err)
// 	}

// 	// 3. 设置当前模块
// 	interp.currentModule = mainModule

// 	// 4. 执行主模块
// 	return interp.EvalModule(mainModule)
// }

// // LoadModule 加载模块（支持缓存）
// func (interp *Interpreter) LoadModule(modulePath string) (core.Module, error) {
// 	// 1. 检查缓存
// 	if module, exists := interp.modules[modulePath]; exists {
// 		return module, nil
// 	}

// 	// 2. 读取文件
// 	content, err := interp.fs.ReadFile(modulePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read module file: %w", err)
// 	}

// 	// 3. 解析AST
// 	ast, err := parser.ParseSourceScript(content)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse module: %w", err)
// 	}

// 	// 4. 创建模块环境
// 	moduleEnv := NewEnvironment()
// 	module := &ModuleImpl{
// 		name:  interp.getModuleName(modulePath),
// 		path:  modulePath,
// 		env:   moduleEnv,
// 		ast:   ast,
// 		state: ModuleLoaded,
// 	}

// 	// 5. 缓存模块
// 	interp.modules[modulePath] = module

// 	// 6. 处理导入语句（预加载依赖）
// 	if err := interp.preloadImports(ast, module); err != nil {
// 		return nil, fmt.Errorf("failed to preload imports: %w", err)
// 	}

// 	return module, nil
// }

// // EvalModule 执行模块
// func (interp *Interpreter) EvalModule(module core.Module) error {
// 	moduleImpl := module.(*ModuleImpl)

// 	// 1. 设置模块环境为当前环境
// 	prevEnv := interp.env
// 	interp.env = moduleImpl.env
// 	defer func() { interp.env = prevEnv }()

// 	// 2. 执行模块中的所有语句
// 	for _, stmt := range moduleImpl.ast.Stmts {
// 		result := interp.Eval(stmt)
// 		if result == nil {
// 			continue
// 		}

// 		// 检查是否有错误
// 		if errorValue, ok := result.(*object.ErrorValue); ok {
// 			return fmt.Errorf("module execution error: %s", errorValue.Value)
// 		}
// 	}

// 	return nil
// }

// // preloadImports 预加载导入的模块
// func (interp *Interpreter) preloadImports(ast *ast.File, currentModule *ModuleImpl) error {
// 	for _, stmt := range ast.Stmts {
// 		if importStmt, ok := stmt.(*ast.Import); ok {
// 			// 解析导入路径
// 			var importPath string
// 			switch {
// 			case importStmt.ImportAll != nil:
// 				importPath = importStmt.ImportAll.Path
// 			case importStmt.ImportSingle != nil:
// 				importPath = importStmt.ImportSingle.Path
// 			case importStmt.ImportMulti != nil:
// 				importPath = importStmt.ImportMulti.Path
// 			}

// 			// 解析模块路径
// 			modulePath, err := interp.resolveModulePath(importPath)
// 			if err != nil {
// 				return fmt.Errorf("failed to resolve import path %s: %w", importPath, err)
// 			}

// 			// 预加载模块（但不执行）
// 			_, err = interp.LoadModule(modulePath)
// 			if err != nil {
// 				return fmt.Errorf("failed to preload module %s: %w", importPath, err)
// 			}
// 		}
// 	}
// 	return nil
// }

// // getModuleName 从路径获取模块名
// func (interp *Interpreter) getModuleName(modulePath string) string {
// 	// 从路径中提取文件名（不含扩展名）
// 	baseName := interp.fs.Base(modulePath)
// 	ext := interp.fs.Ext(modulePath)
// 	if ext != "" {
// 		baseName = baseName[:len(baseName)-len(ext)]
// 	}
// 	return baseName
// }
