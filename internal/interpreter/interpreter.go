package interpreter

import (
	"fmt"
	"math/big"
	"os"
	"slices"
	"strings"

	"github.com/hulo-lang/hulo/internal/core"
	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
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
	moduleSystem    ModuleSystem

	fs vfs.VFS

	cvt object.Converter
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
	case *ast.File:
		file := &ast.File{
			Docs:  node.Docs,
			Name:  node.Name,
			Decls: node.Decls,
		}
		for _, stmt := range node.Stmts {
			file.Stmts = append(file.Stmts, interp.Eval(stmt).(ast.Stmt))
		}
		return file

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

	case *ast.Import:

	case *ast.ComptimeStmt:
		evaluatedObject := interp.executeComptimeStmt(node.X)
		return interp.object2Node(evaluatedObject)
	case *ast.ComptimeExpr:
		evaluatedObject := interp.executeComptimeExpr(node)
		return interp.object2Node(evaluatedObject)
	}
	return node
}

func (interp *Interpreter) executeComptimeStmt(node ast.Stmt) object.Value {
	switch node := node.(type) {
	case *ast.BlockStmt:
		for _, stmt := range node.List {
			interp.Eval(stmt)
		}
	case *ast.Import:
		return interp.executeImport(node)
	case *ast.ExprStmt:
		return interp.executeComptimeExpr(node.X)
	case *ast.AssignStmt:
		return interp.executeAssignStmt(node)
	case *ast.TypeDecl:
		return interp.executeTypeDecl(node)
	}
	return nil
}

func (interp *Interpreter) executeTypeDecl(node *ast.TypeDecl) object.Value {
	typ, err := interp.cvt.ConvertType(node.Value)
	if err != nil {
		return &object.ErrorValue{Value: err.Error()}
	}

	fmt.Println(typ, node.Name, "TODO: 在类型表注册类型")
	return nil
}

func (interp *Interpreter) executeComptimeExpr(node ast.Expr) object.Value {
	switch node := node.(type) {
	case *ast.Ident:
		return interp.executeIdent(node)
	case *ast.SelectExpr:
		return interp.executeSelectExpr(node)
	case *ast.CallExpr:
		return interp.executeCallExpr(node)
	case *ast.BinaryExpr:
		return interp.executeBinaryExpr(node)
	case *ast.RefExpr:

	default:
		value, err := interp.cvt.ConvertValue(node)
		if value == nil || err != nil {
			return &object.ErrorValue{Value: "failed to convert expression"}
		}
		return value
	}
	return nil
}

func (interp *Interpreter) executeCallExpr(node *ast.CallExpr) object.Value {
	fn := interp.executeComptimeExpr(node.Fun)

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
		values[i] = interp.executeComptimeExpr(expr)
		if values[i] == nil {
			return []object.Value{&object.ErrorValue{Value: "failed to evaluate expression"}}
		}
	}
	return values
}

func (interp *Interpreter) executeBinaryExpr(node *ast.BinaryExpr) object.Value {
	lhs := interp.executeComptimeExpr(node.X)
	rhs := interp.executeComptimeExpr(node.Y)

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
		return &object.NumberValue{Value: lv.Add(lv, rv)}
	case token.MINUS:
		return &object.NumberValue{Value: lv.Sub(lv, rv)}
	case token.ASTERISK:
		return &object.NumberValue{Value: lv.Mul(lv, rv)}
	case token.SLASH:
		return &object.NumberValue{Value: lv.Quo(lv, rv)}
	// case token.MOD:
	// return &object.NumberValue{Value: lv.Mod(lv, rv)}
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
	lv := lhs.Interface().(*object.StringValue)
	rv := rhs.Interface().(*object.StringValue)

	switch op {
	case token.PLUS:
		return &object.StringValue{Value: lv.Value + rv.Value}
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
	lhs := interp.executeComptimeStmt(node.X)

	switch {
	case interp.isPackageName(lhs):
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

func (interp *Interpreter) executeIdent(node *ast.Ident) object.Value {
	v, ok := interp.GetEnvironment().Get(node.Name)
	if !ok {
		return &object.ErrorValue{Value: fmt.Sprintf("identifier %s not found", node.Name)}
	}

	return v
}

func (interp *Interpreter) executeAssignStmt(node *ast.AssignStmt) object.Value {
	// 计算右值
	rhsValue := interp.executeComptimeExpr(node.Rhs)
	// TODO: hulo 支持 let a: null 这种语法，所以可能没值也合法
	if rhsValue == nil {
		return &object.ErrorValue{Value: "failed to evaluate right-hand side expression"}
	}

	// 根据左值类型处理
	switch lhs := node.Lhs.(type) {
	case *ast.Ident:
		// 简单变量声明/赋值
		return interp.handleIdentAssignment(lhs, rhsValue, node.Scope, node.Tok)
	case *ast.RefExpr:
		return interp.handleRefExprAssignment(lhs, rhsValue, node.Scope, node.Tok)
	default:
		return &object.ErrorValue{Value: "unsupported left-hand side expression type"}
	}
}

func (interp *Interpreter) handleRefExprAssignment(ref *ast.RefExpr, value object.Value, scope token.Token, assignTok token.Token) object.Value {
	err := interp.env.Assign(ref.X.String(), value)
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
			if err := interp.env.Declare(name, value, scope); err != nil {
				return &object.ErrorValue{Value: err.Error()}
			}
		} else {
			// 这是变量重新赋值
			if err := interp.env.Assign(name, value); err != nil {
				return &object.ErrorValue{Value: err.Error()}
			}
		}
		return value

	case token.COLON_ASSIGN:
		// := 声明并赋值（类似 Go）
		if err := interp.env.Declare(name, value, token.LET); err != nil {
			return &object.ErrorValue{Value: err.Error()}
		}
		return value

	default:
		return &object.ErrorValue{Value: "unsupported assignment operator"}
	}
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{
		env: NewEnvironment(),
	}

	return interpreter
}

// GetEnvironment 获取环境
func (interp *Interpreter) GetEnvironment() *Environment {
	return interp.env
}

var stdlibs = []string{}

func (interp *Interpreter) resolveModulePath(path string) (string, error) {
	// 1. 相对路径 (./math, ../utils)
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		return interp.resolveRelativePath(path)
	}

	// 2. 标准库路径 (math, io, net)
	if slices.Contains(stdlibs, path) {
		return interp.resolveStdLibPath(path)
	}

	// 3. 第三方库路径 (github.com/hulo-lang/hulo/stdlibs/math)
	if strings.Contains(path, "/") {
		return interp.resolveThirdPartyPath(path)
	}

	return interp.resolveLocalModule(path)
}

func (interp *Interpreter) resolveRelativePath(path string) (string, error) {
	currentDir := interp.currentModule.Path()

	absolutePath := interp.fs.Join(currentDir, path)

	if !interp.fs.Exists(absolutePath) {
		return "", fmt.Errorf("module %s not found", path)
	}

	return absolutePath, nil
}

func (interp *Interpreter) resolveStdLibPath(path string) (string, error) {
	huloPath := os.Getenv("HULO_PATH")

	target := interp.fs.Join(huloPath, path)

	if !interp.fs.Exists(target) {
		return "", fmt.Errorf("module %s not found", path)
	}

	return target, nil
}

func (interp *Interpreter) resolveThirdPartyPath(path string) (string, error) {
	huloPath := os.Getenv("HULO_MODULES")

	target := interp.fs.Join(huloPath, path)

	if !interp.fs.Exists(target) {
		return "", fmt.Errorf("module %s not found", path)
	}

	return target, nil
}

func (interp *Interpreter) resolveLocalModule(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	target := interp.fs.Join(wd, path)

	if !interp.fs.Exists(target) {
		return "", fmt.Errorf("module %s not found", path)
	}

	return target, nil
}

type ModuleSystem interface {
	// 加载模块
	LoadModule(path string) (core.Module, error)

	// 解析符号（跨模块查找）
	ResolveSymbol(name string, currentModule string) (Symbol, error)

	// 获取模块
	GetModule(name string) (core.Module, bool)
}

// Symbol 表示一个符号
type Symbol struct {
	Name   string
	Value  object.Value
	Type   object.Type
	Module string // 所属模块
	Kind   SymbolKind
}

type SymbolKind int

const (
	SymbolValue SymbolKind = iota
	SymbolType
	SymbolFunction
	SymbolModule
)

func (interp *Interpreter) resolveSymbol(name string) (core.Value, error) {
	// 1. 首先在当前模块中查找
	if value, found := interp.currentModule.GetValue(name); found {
		return value, nil
	}

	// 2. 在导入的模块中查找
	for _, module := range interp.importedModules {
		if value, found := module.GetValue(name); found {
			return value, nil
		}
	}

	// 3. 在默认模块中查找（如 std 库）
	if value, found := interp.defaultModule.GetValue(name); found {
		return value, nil
	}

	// 4. 在全局符号表中查找
	if symbol, found := interp.moduleSystem.ResolveSymbol(name, interp.currentModule.Name()); found {
		return symbol.Value, nil
	}

	return nil, fmt.Errorf("symbol %s not found", name)
}

func (interp *Interpreter) getTypeModule(typ object.Type) string {
	// 1. 基本类型属于内置模块
	if interp.isBuiltinType(typ) {
		return "builtin"
	}

	// 2. 检查类型是否来自特定模块
	if moduleType, ok := typ.(*object.ModuleType); ok {
		return moduleType.ModuleName
	}

	// 3. 从类型注册表中查找
	if moduleName, found := interp.typeRegistry.GetTypeModule(typ.Name()); found {
		return moduleName
	}

	// 4. 默认返回当前模块
	return interp.currentModule.Name()
}

func (interp *Interpreter) executeSelectExpr(node *ast.SelectExpr) object.Value {
	// 1. 求值左操作数
	lhs := interp.executeComptimeStmt(node.X)

	// 2. 判断左操作数的类型
	switch {
	case interp.isModuleReference(lhs):
		// 这是一个模块引用 (math.PI)
		return interp.resolveModuleMember(lhs, node.Y)

	case interp.isPackageName(lhs):
		// 这是一个包名 (math.PI)
		return interp.resolvePackageMember(lhs, node.Y)

	default:
		// 这是一个普通的成员访问 (obj.field)
		return interp.resolveObjectMember(lhs, node.Y)
	}
}

func (interp *Interpreter) isModuleReference(v object.Value) bool {
	// 检查是否是模块值
	if moduleValue, ok := v.(*object.ModuleValue); ok {
		return moduleValue != nil
	}
	return false
}

func (interp *Interpreter) resolveModuleMember(moduleValue object.Value, selector ast.Expr) object.Value {
	module := moduleValue.(*object.ModuleValue).Module

	// 获取选择器名称
	var memberName string
	if ident, ok := selector.(*ast.Ident); ok {
		memberName = ident.Name
	} else {
		return &object.ErrorValue{Value: "invalid selector"}
	}

	// 从模块中查找成员
	if value, found := module.GetValue(memberName); found {
		return value
	}

	if typ, found := module.GetType(memberName); found {
		return &object.TypeValue{Type: typ}
	}

	return &object.ErrorValue{Value: fmt.Sprintf("member %s not found in module", memberName)}
}

func (interp *Interpreter) executeImport(node *ast.Import) object.Value {
	// 1. 解析导入语句
	var path string
	var alias string

	switch {
	case node.ImportAll != nil:
		// import * as math from "math"
		path = node.ImportAll.Path
		alias = node.ImportAll.Alias

	case node.ImportSingle != nil:
		// import "math" as math
		path = node.ImportSingle.Path
		alias = node.ImportSingle.Alias

	case node.ImportMulti != nil:
		// import { PI, E } from "math"
		path = node.ImportMulti.Path
		// 处理多个导入项
		for _, field := range node.ImportMulti.List {
			fieldName := field.Field
			fieldAlias := field.Alias
			if fieldAlias == "" {
				fieldAlias = fieldName
			}
			// 注册到当前模块的符号表
			interp.registerImportedSymbol(fieldAlias, path, fieldName)
		}
		return object.NULL
	}

	modulePath, err := interp.resolveModulePath(path)
	if err != nil {
		return &object.ErrorValue{Value: err.Error()}
	}

	// 2. 加载模块
	moduleFile, err := interp.fs.ReadFile(modulePath)
	if err != nil {
		return &object.ErrorValue{Value: err.Error()}
	}

	ast, err := parser.ParseSourceScript(moduleFile)
	if err != nil {
		return &object.ErrorValue{Value: err.Error()}
	}

	// 创建模块
	fmt.Println(ast, alias)
	// 3. 注册模块别名到当前环境
	// if alias != "" {
	// 	interp.env.Set(alias, &object.ModuleValue{Module: module})
	// }

	return object.NULL
}

// ExecuteMain 执行主模块
func (interp *Interpreter) ExecuteMain(mainFile string) error {
	// 1. 解析主模块路径
	modulePath, err := interp.resolveModulePath(mainFile)
	if err != nil {
		return fmt.Errorf("failed to resolve main module path: %w", err)
	}

	// 2. 加载主模块
	mainModule, err := interp.LoadModule(modulePath)
	if err != nil {
		return fmt.Errorf("failed to load main module: %w", err)
	}

	// 3. 设置当前模块
	interp.currentModule = mainModule

	// 4. 执行主模块
	return interp.EvalModule(mainModule)
}

// LoadModule 加载模块（支持缓存）
func (interp *Interpreter) LoadModule(modulePath string) (core.Module, error) {
	// 1. 检查缓存
	if module, exists := interp.modules[modulePath]; exists {
		return module, nil
	}

	// 2. 读取文件
	content, err := interp.fs.ReadFile(modulePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read module file: %w", err)
	}

	// 3. 解析AST
	ast, err := parser.ParseSourceScript(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse module: %w", err)
	}

	// 4. 创建模块环境
	moduleEnv := NewEnvironment()
	module := &ModuleImpl{
		name:  interp.getModuleName(modulePath),
		path:  modulePath,
		env:   moduleEnv,
		ast:   ast,
		state: ModuleLoaded,
	}

	// 5. 缓存模块
	interp.modules[modulePath] = module

	// 6. 处理导入语句（预加载依赖）
	if err := interp.preloadImports(ast, module); err != nil {
		return nil, fmt.Errorf("failed to preload imports: %w", err)
	}

	return module, nil
}

// EvalModule 执行模块
func (interp *Interpreter) EvalModule(module core.Module) error {
	moduleImpl := module.(*ModuleImpl)

	// 1. 设置模块环境为当前环境
	prevEnv := interp.env
	interp.env = moduleImpl.env
	defer func() { interp.env = prevEnv }()

	// 2. 执行模块中的所有语句
	for _, stmt := range moduleImpl.ast.Stmts {
		result := interp.Eval(stmt)
		if result == nil {
			continue
		}

		// 检查是否有错误
		if errorValue, ok := result.(*object.ErrorValue); ok {
			return fmt.Errorf("module execution error: %s", errorValue.Value)
		}
	}

	return nil
}

// preloadImports 预加载导入的模块
func (interp *Interpreter) preloadImports(ast *ast.File, currentModule *ModuleImpl) error {
	for _, stmt := range ast.Stmts {
		if importStmt, ok := stmt.(*ast.Import); ok {
			// 解析导入路径
			var importPath string
			switch {
			case importStmt.ImportAll != nil:
				importPath = importStmt.ImportAll.Path
			case importStmt.ImportSingle != nil:
				importPath = importStmt.ImportSingle.Path
			case importStmt.ImportMulti != nil:
				importPath = importStmt.ImportMulti.Path
			}

			// 解析模块路径
			modulePath, err := interp.resolveModulePath(importPath)
			if err != nil {
				return fmt.Errorf("failed to resolve import path %s: %w", importPath, err)
			}

			// 预加载模块（但不执行）
			_, err = interp.LoadModule(modulePath)
			if err != nil {
				return fmt.Errorf("failed to preload module %s: %w", importPath, err)
			}
		}
	}
	return nil
}

// getModuleName 从路径获取模块名
func (interp *Interpreter) getModuleName(modulePath string) string {
	// 从路径中提取文件名（不含扩展名）
	baseName := interp.fs.Base(modulePath)
	ext := interp.fs.Ext(modulePath)
	if ext != "" {
		baseName = baseName[:len(baseName)-len(ext)]
	}
	return baseName
}
