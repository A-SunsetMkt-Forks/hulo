package module

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFunctionSymbol(t *testing.T) {
	code := `
pub fn add(a: num, b: num) -> num {
	return $a + $b
}

fn private_func(x: str) -> bool {
	return $x == "test"
}`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试函数符号提取
	functions := module.GetFunctionSymbols()
	assert.Len(t, functions, 2)

	// 检查第一个函数
	addFunc := functions[0]
	assert.Equal(t, "add", addFunc.GetName())
	assert.Equal(t, "num", addFunc.ReturnType)
	assert.True(t, addFunc.IsExported())
	assert.Len(t, addFunc.Parameters, 2)
	assert.Equal(t, "a", addFunc.Parameters[0].Name)
	assert.Equal(t, "num", addFunc.Parameters[0].Type)
	assert.Equal(t, "b", addFunc.Parameters[1].Name)
	assert.Equal(t, "num", addFunc.Parameters[1].Type)

	// 检查第二个函数
	privateFunc := functions[1]
	assert.Equal(t, "private_func", privateFunc.GetName())
	assert.Equal(t, "bool", privateFunc.ReturnType)
	assert.False(t, privateFunc.IsExported())
	assert.Len(t, privateFunc.Parameters, 1)
	assert.Equal(t, "x", privateFunc.Parameters[0].Name)
	assert.Equal(t, "str", privateFunc.Parameters[0].Type)
}

func TestExtractClassSymbol(t *testing.T) {
	code := `
pub class User {
	pub name: str
	age: num
	email: str = "default@example.com"

	pub fn get_name() -> str {
		return $this.name
	}

	fn get_age() -> num {
		return $this.age
	}
}

class InternalClass {
	value: num
}
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 调试：打印符号表
	module.PrintSymbolTable()

	// 测试类符号提取
	userClass := module.LookupClassSymbol("User")
	assert.NotNil(t, userClass)

	assert.Equal(t, "User", userClass.GetName())
	assert.True(t, userClass.IsExported())
	assert.Len(t, userClass.Fields, 3)
	assert.Len(t, userClass.Methods, 2)

	// 检查字段
	nameField := userClass.Fields["name"]
	assert.Equal(t, "name", nameField.Name)
	assert.Equal(t, "str", nameField.Type)
	assert.True(t, nameField.IsPublic)

	ageField := userClass.Fields["age"]
	assert.Equal(t, "age", ageField.Name)
	assert.Equal(t, "num", ageField.Type)
	assert.False(t, ageField.IsPublic)

	emailField := userClass.Fields["email"]
	assert.Equal(t, "email", emailField.Name)
	assert.Equal(t, "str", emailField.Type)
	assert.False(t, emailField.IsPublic)
	assert.NotNil(t, emailField.Default)

	// 检查方法
	getNameMethod := userClass.Methods["get_name"]
	assert.Equal(t, "get_name", getNameMethod.GetName())
	assert.Equal(t, "str", getNameMethod.ReturnType)
	assert.True(t, getNameMethod.IsExported())
	assert.True(t, getNameMethod.IsMethod)
	assert.Equal(t, "User", getNameMethod.Receiver)

	getAgeMethod := userClass.Methods["get_age"]
	assert.Equal(t, "get_age", getAgeMethod.GetName())
	assert.Equal(t, "num", getAgeMethod.ReturnType)
	assert.False(t, getAgeMethod.IsExported())
	assert.True(t, getAgeMethod.IsMethod)

	internalClass := module.LookupClassSymbol("InternalClass")
	assert.NotNil(t, internalClass)
	assert.Equal(t, "InternalClass", internalClass.GetName())
	assert.False(t, internalClass.IsExported())
	assert.Len(t, internalClass.Fields, 1)
	assert.Len(t, internalClass.Methods, 0)
}

func TestExtractVariableSymbol(t *testing.T) {
	code := `
var global_var: num = 42
const PI: num = 3.14159
pub const VERSION: str = "1.0.0"

let count: num = 10
$local_var := 10
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试变量符号提取
	variables := module.GetVariableSymbols()
	assert.Len(t, variables, 3)

	// 检查变量
	globalVar := module.LookupVariableSymbol("global_var")
	pi := module.LookupVariableSymbol("PI")
	version := module.LookupVariableSymbol("VERSION")

	// 验证global_var
	require.NotNil(t, globalVar)
	assert.Equal(t, "global_var", globalVar.GetName())
	assert.Equal(t, "num", globalVar.DataType)
	assert.False(t, globalVar.IsExported())

	// 验证PI
	require.NotNil(t, pi)
	assert.Equal(t, "PI", pi.GetName())
	assert.Equal(t, "num", pi.DataType)
	assert.False(t, pi.IsExported())

	// 验证VERSION
	require.NotNil(t, version)
	assert.Equal(t, "VERSION", version.GetName())
	assert.Equal(t, "str", version.DataType)
	assert.True(t, version.IsExported())

	count := module.LookupVariableSymbol("count")
	assert.Nil(t, count)
}

func TestExtractEnumSymbol(t *testing.T) {
	code := `
pub enum Status {
	Pending
	Approved
	Rejected
}

enum Color {
	Red = "red"
	Green = "green"
	Blue = "blue"
}
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试枚举符号提取
	classes := module.GetClassSymbols() // 枚举被当作类处理
	assert.Len(t, classes, 2)

	// 检查Status枚举
	var statusEnum *ClassSymbol
	var colorEnum *ClassSymbol

	for _, c := range classes {
		switch c.GetName() {
		case "Status":
			statusEnum = c
		case "Color":
			colorEnum = c
		}
	}

	require.NotNil(t, statusEnum)
	assert.Equal(t, "Status", statusEnum.GetName())
	assert.True(t, statusEnum.IsExported())

	require.NotNil(t, colorEnum)
	assert.Equal(t, "Color", colorEnum.GetName())
	assert.False(t, colorEnum.IsExported())
}

func TestExtractTraitSymbol(t *testing.T) {
	code := `
pub trait Printable {
	pub fn print() -> str
}

trait InternalTrait {
	fn internal_method() -> num
}
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试trait符号提取
	classes := module.GetClassSymbols() // trait被当作类处理
	assert.Len(t, classes, 2)

	// 检查Printable trait
	var printableTrait *ClassSymbol
	var internalTrait *ClassSymbol

	for _, c := range classes {
		switch c.GetName() {
		case "Printable":
			printableTrait = c
		case "InternalTrait":
			internalTrait = c
		}
	}

	require.NotNil(t, printableTrait)
	assert.Equal(t, "Printable", printableTrait.GetName())
	assert.True(t, printableTrait.IsExported())

	require.NotNil(t, internalTrait)
	assert.Equal(t, "InternalTrait", internalTrait.GetName())
	assert.False(t, internalTrait.IsExported())
}

func TestExtractModuleSymbol(t *testing.T) {
	code := `
pub mod utils {
	fn helper() -> num {
		return 42
	}
}

mod internal {
	pub fn public_helper() -> str {
		return "helper"
	}
}

pub mod a {
	pub mod b {
		pub fn public_helper() -> str {
			return "helper"
		}
	}
}
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 检查utils模块
	utilsMod := module.LookupNamespaceSymbol("utils")
	internalMod := module.LookupNamespaceSymbol("internal")
	aMod := module.LookupNamespaceSymbol("a")
	bMod := aMod.Symbols["b"].(*NamespaceSymbol)

	require.NotNil(t, utilsMod)
	assert.Equal(t, "utils", utilsMod.GetName())
	assert.True(t, utilsMod.IsExported())

	require.NotNil(t, internalMod)
	assert.Equal(t, "internal", internalMod.GetName())
	assert.False(t, internalMod.IsExported())

	require.NotNil(t, aMod)
	assert.Equal(t, "a", aMod.GetName())
	assert.True(t, aMod.IsExported())

	require.NotNil(t, bMod)
	assert.Equal(t, "b", bMod.GetName())
	assert.True(t, bMod.IsExported())
}

func TestExtractTypeAliasSymbol(t *testing.T) {
	code := `
type UserID = num
type Email = str
pub type Status = str
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试类型别名符号提取
	constants := module.Symbols.GetSymbolsByKind(SymbolConst) // 类型别名被当作常量处理
	assert.Len(t, constants, 3)

	// 检查类型别名
	var userID *ConstantSymbol
	var email *ConstantSymbol
	var status *ConstantSymbol

	for _, c := range constants {
		if cs, ok := AsConstantSymbol(c); ok {
			switch cs.GetName() {
			case "UserID":
				userID = cs
			case "Email":
				email = cs
			case "Status":
				status = cs
			}
		}
	}

	require.NotNil(t, userID)
	assert.Equal(t, "UserID", userID.GetName())
	assert.False(t, userID.IsExported())

	require.NotNil(t, email)
	assert.Equal(t, "Email", email.GetName())
	assert.False(t, email.IsExported())

	require.NotNil(t, status)
	assert.Equal(t, "Status", status.GetName())
	assert.True(t, status.IsExported())
}

func TestSymbolStats(t *testing.T) {
	code := `
var global_var: num = 42
const PI: num = 3.14159

pub fn add(a: num, b: num) -> num {
	return $a + $b
}

fn private_func() -> str {
	return "private"
}

pub class User {
	name: str
	pub fn get_name() -> str {
		return $this.name
	}
}

enum Status {
	Pending
	Approved
}

pub trait Printable {
	fn print() -> str
}

pub mod utils {
	fn helper() -> num {
		return 42
	}
}

type UserID = num
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试符号统计
	stats := module.GetSymbolStats()
	assert.Equal(t, 2, stats["variables"])
	assert.Equal(t, 2, stats["functions"])
	assert.Equal(t, 1, stats["classes"])
	assert.Equal(t, 1, stats["constants"])
	assert.Equal(t, 1, stats["namespaces"])

	// 测试导出符号
	exported := module.GetExportedSymbols()
	assert.Len(t, exported, 4) // add, User, Printable, utils

	// 验证导出的符号
	exportedNames := make(map[string]bool)
	for _, symbol := range exported {
		exportedNames[symbol.GetName()] = true
	}

	assert.True(t, exportedNames["add"])
	assert.True(t, exportedNames["User"])
	assert.True(t, exportedNames["Printable"])
	assert.True(t, exportedNames["utils"])
}

func TestSymbolMangling(t *testing.T) {
	code := `
pub fn add(a: num, b: num) -> num {
	return $a + $b
}

var global_var: num = 42
`

	module := createTestModule(t, code)
	require.NoError(t, module.BuildSymbolTable())

	// 测试符号混淆
	for _, symbol := range module.Symbols.GlobalSymbols {
		assert.True(t, symbol.IsMangled())
		assert.NotEmpty(t, symbol.GetMangledName())
		assert.NotEqual(t, symbol.GetName(), symbol.GetMangledName())
	}

	// 验证混淆后的名称格式
	addFunc := module.GetSymbolByName("add")
	require.NotNil(t, addFunc)
	assert.True(t, addFunc.IsMangled())
	assert.Contains(t, addFunc.GetMangledName(), "_")
}

// 辅助函数：创建测试模块
func createTestModule(t *testing.T, code string) *Module {
	// 解析代码
	file, err := parser.ParseSourceScript(code)
	require.NoError(t, err)
	// ast.Print(file)
	module := &Module{
		ModuleID: 1,
		Name:     "test_module",
		Path:     "/test/module.hl",
		AST:      file,
		Exports:  make(map[string]*ExportInfo),
		Imports:  []*ImportInfo{},
		Symbols:  nil,
		state:    ModuleStateUnresolved,
	}

	return module
}
