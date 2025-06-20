package interpreter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Interpreter struct {
	debugger *Debugger
}

func (interp *Interpreter) shouldBreak(node ast.Node) bool {
	// pos := node.Pos()
	if bp, ok := interp.debugger.breakpoints["file"][1]; ok {
		fmt.Println(bp, "如果是条件断点，评估条件")
	}
	return false
}

func (interp *Interpreter) Eval(node ast.Node) ast.Node {
	switch node := node.(type) {
	/// Static World

	case *ast.IfStmt:
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

	/// Dynamic World

	case *ast.ComptimeStmt:
		evaluatedObject := interp.executeComptimeBlock(node.X)
		return interp.object2Node(evaluatedObject)
	case *ast.ComptimeExpr:

	}
	return node
}

func (interp *Interpreter) executeComptimeBlock(b ast.Node) object.Value {
	return nil
}

func (interp *Interpreter) object2Node(v object.Value) ast.Node {
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
		return evalIdent(ctx, node)
	case *ast.CallExpr:
		return evalCallExpr(ctx, node)
	case *ast.ExtensionDecl:
		return evalExtensionDecl(ctx, node)
	}
	return nil
}

func evalExtensionDecl(ctx *Context, node *ast.ExtensionDecl) object.Value {
	if node.ExtensionClass != nil {
		// 将类查出来 插入进去
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
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}
	if target := filepath.Join(wd, x.Name()); ctx.os.Exist(target) {
		// ctx.mem.Import(target)
	}

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
	function := Evaluate(ctx, node.Fun)
	args := evalExpressions(ctx, node.Recv)
	// 要根据 args 类型 拿到method
	// TODO function.Type().Call(args...) 自动匹配合适的函数
	function.Type().Method(0).Call(args...)
	// Call 的逻辑中，builtin直接执行，非builtin要执行语法树
	return nil
}

func evalExpressions(ctx *Context, exprs []ast.Expr) []object.Value {
	return nil
}

func evalIdent(ctx *Context, node *ast.Ident) object.Value {
	// math.PI.to_str()
	// PI.to_str()
	// 检查当前 字面量 是什么东西：变量名？包名？函数名？

	// 1. 先找变量 还要支持 import 的时候更新Ctx？
	val, ok := ctx.Get(node.Name)

	if ok {
		return val
	}
	// 2. 从 builtin 包找
	// val, ok = ctx.mem.Get(node.Name)
	if ok {
		return val
	}

	// 3. 还是没有 从 declare 找

	// 4. 从 程序 的 import 找是不是包名

	// 5. 没有找到抛出 unknown 符号 因为可能是没有声明而已？
	return nil
}
