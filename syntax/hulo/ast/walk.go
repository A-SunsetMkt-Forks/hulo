// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Walk traverses an AST in depth-first order;
// if the visitor w returned by v.Visit(node) is not nil,
// Walk visits each of the children of node with the visitor w,
// followed by a call of w.Visit(nil).
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}

	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *File:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.Stmts)
		walkList(v, n.Decls)
	case *FuncDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)

		if n.Name != nil {
			Walk(v, n.Name)
		}

		walkList(v, n.TypeParams)

		walkList(v, n.Recv)

		if n.Type != nil {
			Walk(v, n.Type)
		}

		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ConstructorDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)
		if n.ClsName != nil {
			Walk(v, n.ClsName)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.Recv)
		walkList(v, n.InitFields)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *AssignStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}

		Walk(v, n.Lhs)

		if n.Type != nil {
			Walk(v, n.Type)
		}

		Walk(v, n.Rhs)

	case *ClassDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.TypeParams)
		if n.Parent != nil {
			Walk(v, n.Parent)
		}
		if n.Fields != nil {
			Walk(v, n.Fields)
		}
		walkList(v, n.Ctors)
		walkList(v, n.Methods)

	case *TraitDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Fields != nil {
			Walk(v, n.Fields)
		}

	case *ImplDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Trait != nil {
			Walk(v, n.Trait)
		}
		if n.ImplDeclBody != nil {
			if n.ImplDeclBody.Class != nil {
				Walk(v, n.ImplDeclBody.Class)
			}
			if n.ImplDeclBody.Body != nil {
				Walk(v, n.ImplDeclBody.Body)
			}
		}
		if n.ImplDeclBinding != nil {
			walkList(v, n.ImplDeclBinding.Classes)
		}

	case *EnumDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.TypeParams)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *BasicEnumBody:
		walkList(v, n.Values)

	case *AssociatedEnumBody:
		if n.Fields != nil {
			Walk(v, n.Fields)
		}
		walkList(v, n.Values)
		walkList(v, n.Methods)

	case *ADTEnumBody:
		walkList(v, n.Variants)
		walkList(v, n.Methods)

	case *EnumValue:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}
		walkList(v, n.Data)

	case *EnumVariant:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Fields != nil {
			Walk(v, n.Fields)
		}

	case *TypeParameter:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.Constraints)

	case *TypeReference:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.TypeParams)

	case *UnionType:
		walkList(v, n.Types)

	case *IntersectionType:
		walkList(v, n.Types)

	case *Parameter:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *NamedParameters:
		walkList(v, n.Params)

	case *NullableType:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *TypeDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *ExtensionDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ExtensionEnum:
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ExtensionClass:
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ExtensionTrait:
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ExtensionType:
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ExtensionMod:
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *DeclareDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.X != nil {
			Walk(v, n.X)
		}

	case *UseDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Lhs != nil {
			Walk(v, n.Lhs)
		}
		if n.Rhs != nil {
			Walk(v, n.Rhs)
		}

	case *OperatorDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)
		walkList(v, n.Params)
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ModDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ExternDecl:
		walkList(v, n.List)

	case *GetAccessor:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *SetAccessor:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *FieldList:
		// Field is not a Node type, so we can't walk it directly
		// The FieldList itself is walked, but individual fields are not

	case *Import:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}

	case *AnyLiteral, *StringLiteral, *NumericLiteral,
		*TrueLiteral, *FalseLiteral, *Ident, *Comment:
		// 叶子节点，无需遍历子节点

	case *CommentGroup:
		walkList(v, n.List)

	case *BlockStmt:
		walkList(v, n.List)

	case *CmdStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.Recv)

	case *BreakStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Label != nil {
			Walk(v, n.Label)
		}

	case *ContinueStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}

	case *ComptimeStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
		if n.Else != nil {
			Walk(v, n.Else)
		}

	case *IfStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
		if n.Else != nil {
			Walk(v, n.Else)
		}

	case *LabelStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	case *ForeachStmt:
		if n.Index != nil {
			Walk(v, n.Index)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}
		if n.Var != nil {
			Walk(v, n.Var)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ForInStmt:
		if n.Index != nil {
			Walk(v, n.Index)
		}
		if n.RangeExpr.Start != nil {
			Walk(v, n.RangeExpr.Start)
		}
		if n.RangeExpr.End_ != nil {
			Walk(v, n.RangeExpr.End_)
		}
		if n.RangeExpr.Step != nil {
			Walk(v, n.RangeExpr.Step)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *WhileStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ForStmt:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Post != nil {
			Walk(v, n.Post)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *DoWhileStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *MatchStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Expr != nil {
			Walk(v, n.Expr)
		}
		// CaseClause is not a Node type, so we can't walk it directly
		// The MatchStmt itself is walked, but individual cases are not

	case *ReturnStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.X != nil {
			Walk(v, n.X)
		}

	case *DeferStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.X != nil {
			Walk(v, n.X)
		}

	case *ExprStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)
		if n.X != nil {
			Walk(v, n.X)
		}

	case *UnsafeExpr:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}

	case *TryStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
		walkList(v, n.Catches)
		if n.Finally != nil {
			Walk(v, n.Finally)
		}

	case *CatchClause:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *FinallyStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *ThrowStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		if n.X != nil {
			Walk(v, n.X)
		}

	case *Decorator:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.Recv)

	case *RangeExpr:
		if n.Start != nil {
			Walk(v, n.Start)
		}
		if n.End_ != nil {
			Walk(v, n.End_)
		}
		if n.Step != nil {
			Walk(v, n.Step)
		}

	case *CallExpr:
		if n.Fun != nil {
			Walk(v, n.Fun)
		}
		walkList(v, n.TypeParams)
		walkList(v, n.Recv)

	case *CmdExpr:
		if n.Cmd != nil {
			Walk(v, n.Cmd)
		}
		walkList(v, n.Args)
		walkList(v, n.BuiltinArgs)

	case *CmdSubstExpr:
		if n.Fun != nil {
			Walk(v, n.Fun)
		}
		walkList(v, n.Recv)

	case *NewDelExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *RefExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *IncDecExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *SelectExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Y != nil {
			Walk(v, n.Y)
		}

	case *ModAccessExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Y != nil {
			Walk(v, n.Y)
		}

	case *IndexExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Index != nil {
			Walk(v, n.Index)
		}

	case *SliceExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Low != nil {
			Walk(v, n.Low)
		}
		if n.High != nil {
			Walk(v, n.High)
		}
		if n.Max != nil {
			Walk(v, n.Max)
		}

	case *UnaryExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *BinaryExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Y != nil {
			Walk(v, n.Y)
		}

	case *ComptimeExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *ArrayLiteralExpr:
		walkList(v, n.Elems)

	case *ObjectLiteralExpr:
		walkList(v, n.Props)

	case *NamedObjectLiteralExpr:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkList(v, n.Props)

	case *KeyValueExpr:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkList(v, n.Decs)
		if n.Key != nil {
			Walk(v, n.Key)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *SetLiteralExpr:
		walkList(v, n.Elems)

	case *TupleLiteralExpr:
		walkList(v, n.Elems)

	case *NamedTupleLiteralExpr:
		walkList(v, n.Elems)

	case *TypeLiteral:
		walkList(v, n.Members)

	case *ArrayType:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	case *FunctionType:
		walkList(v, n.Recv)
		if n.RetVal != nil {
			Walk(v, n.RetVal)
		}

	case *TupleType:
		walkList(v, n.Types)

	case *SetType:
		walkList(v, n.Types)

	case *ConditionalType:
		if n.CheckType != nil {
			Walk(v, n.CheckType)
		}
		if n.ExtendsType != nil {
			Walk(v, n.ExtendsType)
		}
		if n.TrueType != nil {
			Walk(v, n.TrueType)
		}
		if n.FalseType != nil {
			Walk(v, n.FalseType)
		}

	case *InferType:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *ConditionalExpr:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.WhenTrue != nil {
			Walk(v, n.WhenTrue)
		}
		if n.WhneFalse != nil {
			Walk(v, n.WhneFalse)
		}

	case *CascadeExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Y != nil {
			Walk(v, n.Y)
		}

	case *TypeofExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *AsExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Y != nil {
			Walk(v, n.Y)
		}

	case *LambdaExpr:
		walkList(v, n.Recv)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *Package:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, file := range n.Files {
			Walk(v, file)
		}
	}

	v.Visit(nil)
}

func walkList[N Node](v Visitor, list []N) {
	for _, node := range list {
		Walk(v, node)
	}
}

type predicate func(Node) bool

func (f predicate) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// WalkIf walks the AST if the predicate returns true.
func WalkIf(node Node, f predicate) {
	Walk(predicate(f), node)
}
