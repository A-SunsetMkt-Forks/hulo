// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package unsafe

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/internal/unsafe/generated"
)

// Context represents the execution context with variables and functions
type Context struct {
	Variables map[string]any
	Functions map[string]any
	Macros    map[string]*Macro
	Parent    *Context
}

// Macro represents a macro definition
type Macro struct {
	Name       string
	Parameters []string
	Template   *generated.TemplateContext
}

// NewContext creates a new execution context
func NewContext() *Context {
	return &Context{
		Variables: make(map[string]any),
		Functions: make(map[string]any),
		Macros:    make(map[string]*Macro),
	}
}

// SetVariable sets a variable in the current context
func (c *Context) SetVariable(name string, value any) {
	c.Variables[name] = value
}

// GetVariable gets a variable from current or parent context
func (c *Context) GetVariable(name string) (any, bool) {
	if value, exists := c.Variables[name]; exists {
		return value, true
	}
	if c.Parent != nil {
		return c.Parent.GetVariable(name)
	}
	return nil, false
}

// SetFunction sets a function in the current context
func (c *Context) SetFunction(name string, fn any) {
	c.Functions[name] = fn
}

// GetFunction gets a function from current or parent context
func (c *Context) GetFunction(name string) (any, bool) {
	if fn, exists := c.Functions[name]; exists {
		return fn, true
	}
	if c.Parent != nil {
		return c.Parent.GetFunction(name)
	}
	return nil, false
}

// SetMacro sets a macro in the current context
func (c *Context) SetMacro(name string, macro *Macro) {
	c.Macros[name] = macro
}

// GetMacro gets a macro from current or parent context
func (c *Context) GetMacro(name string) (*Macro, bool) {
	if macro, exists := c.Macros[name]; exists {
		return macro, true
	}
	if c.Parent != nil {
		return c.Parent.GetMacro(name)
	}
	return nil, false
}

// TemplateEngine represents the template engine
type TemplateEngine struct {
	context *Context
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		context: NewContext(),
	}

	// Register built-in functions
	engine.RegisterBuiltins()

	return engine
}

// RegisterBuiltins registers built-in functions
func (engine *TemplateEngine) RegisterBuiltins() {
	engine.context.SetFunction("sum", func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		var sum float64
		for _, arg := range args {
			switch v := arg.(type) {
			case int:
				sum += float64(v)
			case float64:
				sum += v
			case string:
				if num, err := strconv.ParseFloat(v, 64); err == nil {
					sum += num
				}
			}
		}
		return sum
	})

	engine.context.SetFunction("to_str", func(args ...any) string {
		if len(args) == 0 {
			return ""
		}
		return fmt.Sprintf("%v", args[0])
	})

	engine.context.SetFunction("len", func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		v := reflect.ValueOf(args[0])
		switch v.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			return v.Len()
		default:
			return 0
		}
	})
}

// SetVariable sets a variable in the engine context
func (engine *TemplateEngine) SetVariable(name string, value any) {
	engine.context.SetVariable(name, value)
}

// SetFunction sets a function in the engine context
func (engine *TemplateEngine) SetFunction(name string, fn any) {
	engine.context.SetFunction(name, fn)
}

// Execute executes a template with the given context
func (engine *TemplateEngine) Execute(template string) (string, error) {
	stream := antlr.NewInputStream(template)
	lex := generated.NewunsafeLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	parser := generated.NewunsafeParser(tokens)

	visitor := &TemplateVisitor{
		engine:  engine,
		context: engine.context,
	}

	result, err := visitor.VisitTemplate(parser.Template().(*generated.TemplateContext))
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

// TemplateVisitor implements the visitor pattern for template execution
type TemplateVisitor struct {
	engine  *TemplateEngine
	context *Context
}

// VisitTemplate visits the template node
func (v *TemplateVisitor) VisitTemplate(ctx *generated.TemplateContext) (any, error) {
	var result strings.Builder

	for _, child := range ctx.GetChildren() {
		switch node := child.(type) {
		case *generated.ContentContext:
			if content := node.GetText(); content != "" {
				result.WriteString(content)
			}
		case *generated.StatementContext:
			if stmt, err := v.VisitStatement(node); err != nil {
				return nil, err
			} else if stmt != nil {
				// 直接写入字符串，避免额外的格式化
				if str, ok := stmt.(string); ok {
					result.WriteString(str)
				} else {
					result.WriteString(fmt.Sprintf("%v", stmt))
				}
			}
		}
	}

	return result.String(), nil
}

// VisitStatement visits statement nodes
func (v *TemplateVisitor) VisitStatement(ctx *generated.StatementContext) (any, error) {
	if ctx.CommentStatement() != nil {
		return v.VisitCommentStatement(ctx.CommentStatement().(*generated.CommentStatementContext))
	}
	if ctx.IfStatement() != nil {
		return v.VisitIfStatement(ctx.IfStatement().(*generated.IfStatementContext))
	}
	if ctx.LoopStatement() != nil {
		return v.VisitLoopStatement(ctx.LoopStatement().(*generated.LoopStatementContext))
	}
	if ctx.ExpressionStatement() != nil {
		return v.VisitExpressionStatement(ctx.ExpressionStatement().(*generated.ExpressionStatementContext))
	}
	if ctx.MacroStatement() != nil {
		return v.VisitMacroStatement(ctx.MacroStatement().(*generated.MacroStatementContext))
	}
	if ctx.VariableStatement() != nil {
		return v.VisitVariableStatement(ctx.VariableStatement().(*generated.VariableStatementContext))
	}
	return "", nil
}

// VisitCommentStatement visits comment statements (ignored)
func (v *TemplateVisitor) VisitCommentStatement(ctx *generated.CommentStatementContext) (any, error) {
	return "", nil // Comments are ignored
}

// VisitVariableStatement visits variable assignment statements
func (v *TemplateVisitor) VisitVariableStatement(ctx *generated.VariableStatementContext) (any, error) {
	name := ctx.IDENTIFIER().GetText()
	expr, err := v.VisitExpression(ctx.Expression().(*generated.ExpressionContext))
	if err != nil {
		return nil, err
	}

	v.context.SetVariable(name, expr)
	return "", nil
}

// VisitIfStatement visits if statements
func (v *TemplateVisitor) VisitIfStatement(ctx *generated.IfStatementContext) (any, error) {
	condition, err := v.VisitExpression(ctx.Expression().(*generated.ExpressionContext))
	if err != nil {
		return nil, err
	}

	// Evaluate condition
	if isTrue(condition) {
		// Execute the if block template
		return v.VisitTemplate(ctx.Template().(*generated.TemplateContext))
	}

	// Check for else statement
	if ctx.ElseStatement() != nil {
		return v.VisitElseStatement(ctx.ElseStatement().(*generated.ElseStatementContext))
	}

	return "", nil
}

// VisitElseStatement visits else statements
func (v *TemplateVisitor) VisitElseStatement(ctx *generated.ElseStatementContext) (any, error) {
	return v.VisitTemplate(ctx.Template().(*generated.TemplateContext))
}

// VisitLoopStatement visits loop statements
func (v *TemplateVisitor) VisitLoopStatement(ctx *generated.LoopStatementContext) (any, error) {
	itemName := ctx.IDENTIFIER_STMT(0).GetText()
	collectionName := ctx.IDENTIFIER_STMT(1).GetText()

	// Get collection from context
	collection, exists := v.context.GetVariable(collectionName)
	if !exists {
		return "", fmt.Errorf("variable '%s' not found", collectionName)
	}

	// Create child context for loop
	childContext := &Context{
		Variables: make(map[string]any),
		Functions: v.context.Functions,
		Parent:    v.context,
	}

	// Save original visitor context
	originalContext := v.context
	v.context = childContext
	defer func() { v.context = originalContext }()

	var result strings.Builder

	// Iterate over collection
	collectionValue := reflect.ValueOf(collection)
	switch collectionValue.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < collectionValue.Len(); i++ {
			childContext.SetVariable(itemName, collectionValue.Index(i).Interface())
			loopResult, err := v.VisitTemplate(ctx.Template().(*generated.TemplateContext))
			if err != nil {
				return nil, err
			}
			result.WriteString(fmt.Sprintf("%v", loopResult))
		}
	case reflect.Map:
		for _, key := range collectionValue.MapKeys() {
			childContext.SetVariable(itemName, collectionValue.MapIndex(key).Interface())
			loopResult, err := v.VisitTemplate(ctx.Template().(*generated.TemplateContext))
			if err != nil {
				return nil, err
			}
			result.WriteString(fmt.Sprintf("%v", loopResult))
		}
	default:
		return "", fmt.Errorf("cannot iterate over type %T", collection)
	}

	return result.String(), nil
}

// VisitExpressionStatement visits expression statements
func (v *TemplateVisitor) VisitExpressionStatement(ctx *generated.ExpressionStatementContext) (any, error) {
	return v.VisitPipelineExpression(ctx.PipelineExpression().(*generated.PipelineExpressionContext))
}

// VisitPipelineExpression visits pipeline expressions
func (v *TemplateVisitor) VisitPipelineExpression(ctx *generated.PipelineExpressionContext) (any, error) {
	result, err := v.VisitExpression(ctx.Expression(0).(*generated.ExpressionContext))
	if err != nil {
		return nil, err
	}

	// Process pipeline
	for i := 1; i < len(ctx.AllExpression()); i++ {
		pipeExpr := ctx.Expression(i).(*generated.ExpressionContext)

		// Check if this is a function call in the pipeline
		if pipeExpr.FunctionCall() != nil {
			// Get the function name
			fnName := pipeExpr.FunctionCall().(*generated.FunctionCallContext).IDENTIFIER().GetText()

			// Get the function
			fn, exists := v.context.GetFunction(fnName)
			if !exists {
				return nil, fmt.Errorf("function '%s' not found", fnName)
			}

			// Collect function arguments
			var args []any
			args = append(args, result) // First argument is the previous result

			// Add any additional arguments from the function call
			for _, expr := range pipeExpr.FunctionCall().(*generated.FunctionCallContext).AllExpression() {
				if arg, err := v.VisitExpression(expr.(*generated.ExpressionContext)); err != nil {
					return nil, err
				} else {
					args = append(args, arg)
				}
			}

			// Call the function
			fnValue := reflect.ValueOf(fn)
			if fnValue.Type().Kind() != reflect.Func {
				return nil, fmt.Errorf("'%s' is not a function", fnName)
			}

			// Convert args to reflect values
			argValues := make([]reflect.Value, len(args))
			for i, arg := range args {
				argValues[i] = reflect.ValueOf(arg)
			}

			results := fnValue.Call(argValues)
			if len(results) == 0 {
				return nil, nil
			}
			result = results[0].Interface()
		} else {
			// Simple expression in pipeline (like a variable or literal)
			if pipeResult, err := v.VisitExpression(pipeExpr); err != nil {
				return nil, err
			} else {
				// Apply as a single-argument function
				if fn, ok := v.context.GetFunction(fmt.Sprintf("%v", pipeResult)); ok {
					if reflect.TypeOf(fn).Kind() == reflect.Func {
						fnValue := reflect.ValueOf(fn)
						if fnValue.Type().NumIn() == 1 {
							result = fnValue.Call([]reflect.Value{reflect.ValueOf(result)})[0].Interface()
						}
					}
				}
			}
		}
	}

	return result, nil
}

// VisitExpression visits expression nodes
func (v *TemplateVisitor) VisitExpression(ctx *generated.ExpressionContext) (any, error) {
	// Handle logical expressions first
	if ctx.LogicalOrExpression() != nil {
		return v.VisitLogicalOrExpression(ctx.LogicalOrExpression().(*generated.LogicalOrExpressionContext))
	}

	// Handle statement mode logical expressions
	if ctx.LogicalOrExpressionStmt() != nil {
		return v.VisitLogicalOrExpressionStmt(ctx.LogicalOrExpressionStmt().(*generated.LogicalOrExpressionStmtContext))
	}

	// Handle basic expressions
	if ctx.NUMBER() != nil {
		text := ctx.NUMBER().GetText()
		if strings.Contains(text, ".") {
			if num, err := strconv.ParseFloat(text, 64); err == nil {
				return num, nil
			}
		}
		if num, err := strconv.Atoi(text); err == nil {
			return num, nil
		}
		return text, nil
	}

	if ctx.STRING() != nil {
		text := ctx.STRING().GetText()
		// Remove quotes
		return strings.Trim(text, `"'`), nil
	}

	if ctx.BOOLEAN() != nil {
		text := ctx.BOOLEAN().GetText()
		return text == "true", nil
	}

	if ctx.BOOLEAN_STMT() != nil {
		text := ctx.BOOLEAN_STMT().GetText()
		return text == "true", nil
	}

	if ctx.IDENTIFIER() != nil {
		name := ctx.IDENTIFIER().GetText()
		if value, exists := v.context.GetVariable(name); exists {
			return value, nil
		}
		return name, nil // Return as string if not found
	}

	if ctx.VarExpr() != nil {
		return v.VisitVarExpr(ctx.VarExpr().(*generated.VarExprContext))
	}

	if ctx.FunctionCall() != nil {
		return v.VisitFunctionCall(ctx.FunctionCall().(*generated.FunctionCallContext))
	}

	return "", nil
}

// VisitVarExpr visits variable expressions
func (v *TemplateVisitor) VisitVarExpr(ctx *generated.VarExprContext) (any, error) {
	name := ctx.IDENTIFIER().GetText()
	if value, exists := v.context.GetVariable(name); exists {
		return value, nil
	}
	return nil, fmt.Errorf("variable '$%s' not found", name)
}

// VisitFunctionCall visits function calls
func (v *TemplateVisitor) VisitFunctionCall(ctx *generated.FunctionCallContext) (any, error) {
	name := ctx.IDENTIFIER().GetText()

	// First check if it's a macro
	if macro, exists := v.context.GetMacro(name); exists {
		return v.callMacro(macro, ctx.AllExpression())
	}

	// Then check if it's a function
	fn, exists := v.context.GetFunction(name)
	if !exists {
		return nil, fmt.Errorf("function '%s' not found", name)
	}

	// Collect arguments
	var args []any
	for _, expr := range ctx.AllExpression() {
		if arg, err := v.VisitExpression(expr.(*generated.ExpressionContext)); err != nil {
			return nil, err
		} else {
			args = append(args, arg)
		}
	}

	// Call function
	fnValue := reflect.ValueOf(fn)
	if fnValue.Type().Kind() != reflect.Func {
		return nil, fmt.Errorf("'%s' is not a function", name)
	}

	// Convert args to reflect values
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}

	results := fnValue.Call(argValues)
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].Interface(), nil
}

// VisitLogicalOrExpression visits logical OR expressions
func (v *TemplateVisitor) VisitLogicalOrExpression(ctx *generated.LogicalOrExpressionContext) (any, error) {
	if len(ctx.AllLogicalAndExpression()) == 1 {
		return v.VisitLogicalAndExpression(ctx.LogicalAndExpression(0).(*generated.LogicalAndExpressionContext))
	}

	// Handle multiple OR expressions
	var result any
	for i, expr := range ctx.AllLogicalAndExpression() {
		value, err := v.VisitLogicalAndExpression(expr.(*generated.LogicalAndExpressionContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// OR logic: if any operand is true, result is true
			if isTrue(result) || isTrue(value) {
				result = true
			} else {
				result = false
			}
		}
	}

	return result, nil
}

// VisitLogicalAndExpression visits logical AND expressions
func (v *TemplateVisitor) VisitLogicalAndExpression(ctx *generated.LogicalAndExpressionContext) (any, error) {
	if len(ctx.AllEqualityExpression()) == 1 {
		return v.VisitEqualityExpression(ctx.EqualityExpression(0).(*generated.EqualityExpressionContext))
	}

	// Handle multiple AND expressions
	var result any
	for i, expr := range ctx.AllEqualityExpression() {
		value, err := v.VisitEqualityExpression(expr.(*generated.EqualityExpressionContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// AND logic: if any operand is false, result is false
			if !isTrue(result) || !isTrue(value) {
				result = false
			} else {
				result = true
			}
		}
	}

	return result, nil
}

// VisitEqualityExpression visits equality expressions
func (v *TemplateVisitor) VisitEqualityExpression(ctx *generated.EqualityExpressionContext) (any, error) {
	if len(ctx.AllComparisonExpression()) == 1 {
		return v.VisitComparisonExpression(ctx.ComparisonExpression(0).(*generated.ComparisonExpressionContext))
	}

	// Handle multiple equality expressions
	var result any
	for i, expr := range ctx.AllComparisonExpression() {
		value, err := v.VisitComparisonExpression(expr.(*generated.ComparisonExpressionContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// Get the operator
			var operator string
			if i-1 < len(ctx.AllEQ()) {
				operator = "=="
			} else if i-1 < len(ctx.AllEQ())+len(ctx.AllNE()) {
				operator = "!="
			}

			switch operator {
			case "==":
				result = result == value
			case "!=":
				result = result != value
			}
		}
	}

	return result, nil
}

// VisitComparisonExpression visits comparison expressions
func (v *TemplateVisitor) VisitComparisonExpression(ctx *generated.ComparisonExpressionContext) (any, error) {
	if len(ctx.AllPrimaryExpression()) == 1 {
		return v.VisitPrimaryExpression(ctx.PrimaryExpression(0).(*generated.PrimaryExpressionContext))
	}

	// Handle multiple comparison expressions
	var result any
	for i, expr := range ctx.AllPrimaryExpression() {
		value, err := v.VisitPrimaryExpression(expr.(*generated.PrimaryExpressionContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// Get the operator
			var operator string
			if i-1 < len(ctx.AllGT()) {
				operator = ">"
			} else if i-1 < len(ctx.AllGT())+len(ctx.AllLT()) {
				operator = "<"
			} else if i-1 < len(ctx.AllGT())+len(ctx.AllLT())+len(ctx.AllGE()) {
				operator = ">="
			} else if i-1 < len(ctx.AllGT())+len(ctx.AllLT())+len(ctx.AllGE())+len(ctx.AllLE()) {
				operator = "<="
			}

			// Convert to comparable types
			left, right := v.toComparable(result, value)

			switch operator {
			case ">":
				result = left > right
			case "<":
				result = left < right
			case ">=":
				result = left >= right
			case "<=":
				result = left <= right
			}
		}
	}

	return result, nil
}

// VisitPrimaryExpression visits primary expressions
func (v *TemplateVisitor) VisitPrimaryExpression(ctx *generated.PrimaryExpressionContext) (any, error) {

	// Handle NOT operator
	if ctx.NOT() != nil && ctx.PrimaryExpression() != nil {
		value, err := v.VisitPrimaryExpression(ctx.PrimaryExpression().(*generated.PrimaryExpressionContext))
		if err != nil {
			return nil, err
		}
		return !isTrue(value), nil
	}

	// Handle NOT_STMT operator
	if ctx.NOT_STMT() != nil && ctx.PrimaryExpressionStmt() != nil {
		value, err := v.VisitPrimaryExpressionStmt(ctx.PrimaryExpressionStmt().(*generated.PrimaryExpressionStmtContext))
		if err != nil {
			return nil, err
		}
		return !isTrue(value), nil
	}

	// Handle basic expressions
	if ctx.NUMBER() != nil {
		text := ctx.NUMBER().GetText()
		if strings.Contains(text, ".") {
			if num, err := strconv.ParseFloat(text, 64); err == nil {
				return num, nil
			}
		}
		if num, err := strconv.Atoi(text); err == nil {
			return num, nil
		}
		return text, nil
	}

	if ctx.STRING() != nil {
		text := ctx.STRING().GetText()
		return strings.Trim(text, `"'`), nil
	}

	if ctx.BOOLEAN() != nil {
		text := ctx.BOOLEAN().GetText()
		return text == "true", nil
	}

	if ctx.BOOLEAN_STMT() != nil {
		text := ctx.BOOLEAN_STMT().GetText()
		return text == "true", nil
	}

	if ctx.IDENTIFIER() != nil {
		name := ctx.IDENTIFIER().GetText()
		if value, exists := v.context.GetVariable(name); exists {
			return value, nil
		}
		return name, nil
	}

	if ctx.VarExpr() != nil {
		return v.VisitVarExpr(ctx.VarExpr().(*generated.VarExprContext))
	}

	if ctx.FunctionCall() != nil {
		return v.VisitFunctionCall(ctx.FunctionCall().(*generated.FunctionCallContext))
	}

	if ctx.LPAREN() != nil && ctx.Expression() != nil {
		return v.VisitExpression(ctx.Expression().(*generated.ExpressionContext))
	}

	return "", nil
}

// toComparable converts two values to comparable types
func (v *TemplateVisitor) toComparable(left, right any) (float64, float64) {
	toFloat := func(val any) float64 {
		switch v := val.(type) {
		case int:
			return float64(v)
		case float64:
			return v
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		case bool:
			if v {
				return 1
			}
			return 0
		}
		return 0
	}

	return toFloat(left), toFloat(right)
}

// VisitLogicalOrExpressionStmt visits logical OR expressions in statement mode
func (v *TemplateVisitor) VisitLogicalOrExpressionStmt(ctx *generated.LogicalOrExpressionStmtContext) (any, error) {
	if len(ctx.AllLogicalAndExpressionStmt()) == 1 {
		return v.VisitLogicalAndExpressionStmt(ctx.LogicalAndExpressionStmt(0).(*generated.LogicalAndExpressionStmtContext))
	}

	// Handle multiple OR expressions
	var result any
	for i, expr := range ctx.AllLogicalAndExpressionStmt() {
		value, err := v.VisitLogicalAndExpressionStmt(expr.(*generated.LogicalAndExpressionStmtContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// OR logic: if any operand is true, result is true
			if isTrue(result) || isTrue(value) {
				result = true
			} else {
				result = false
			}
		}
	}

	return result, nil
}

// VisitLogicalAndExpressionStmt visits logical AND expressions in statement mode
func (v *TemplateVisitor) VisitLogicalAndExpressionStmt(ctx *generated.LogicalAndExpressionStmtContext) (any, error) {
	if len(ctx.AllEqualityExpressionStmt()) == 1 {
		return v.VisitEqualityExpressionStmt(ctx.EqualityExpressionStmt(0).(*generated.EqualityExpressionStmtContext))
	}

	// Handle multiple AND expressions
	var result any
	for i, expr := range ctx.AllEqualityExpressionStmt() {
		value, err := v.VisitEqualityExpressionStmt(expr.(*generated.EqualityExpressionStmtContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// AND logic: if any operand is false, result is false
			if !isTrue(result) || !isTrue(value) {
				result = false
			} else {
				result = true
			}
		}
	}

	return result, nil
}

// VisitEqualityExpressionStmt visits equality expressions in statement mode
func (v *TemplateVisitor) VisitEqualityExpressionStmt(ctx *generated.EqualityExpressionStmtContext) (any, error) {
	if len(ctx.AllComparisonExpressionStmt()) == 1 {
		return v.VisitComparisonExpressionStmt(ctx.ComparisonExpressionStmt(0).(*generated.ComparisonExpressionStmtContext))
	}

	// Handle multiple equality expressions
	var result any
	for i, expr := range ctx.AllComparisonExpressionStmt() {
		value, err := v.VisitComparisonExpressionStmt(expr.(*generated.ComparisonExpressionStmtContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// Get the operator
			var operator string
			if i-1 < len(ctx.AllEQ_STMT()) {
				operator = "=="
			} else if i-1 < len(ctx.AllEQ_STMT())+len(ctx.AllNE_STMT()) {
				operator = "!="
			}

			switch operator {
			case "==":
				result = result == value
			case "!=":
				result = result != value
			}
		}
	}

	return result, nil
}

// VisitComparisonExpressionStmt visits comparison expressions in statement mode
func (v *TemplateVisitor) VisitComparisonExpressionStmt(ctx *generated.ComparisonExpressionStmtContext) (any, error) {
	if len(ctx.AllPrimaryExpressionStmt()) == 1 {
		return v.VisitPrimaryExpressionStmt(ctx.PrimaryExpressionStmt(0).(*generated.PrimaryExpressionStmtContext))
	}

	// Handle multiple comparison expressions
	var result any
	for i, expr := range ctx.AllPrimaryExpressionStmt() {
		value, err := v.VisitPrimaryExpressionStmt(expr.(*generated.PrimaryExpressionStmtContext))
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
		} else {
			// Get the operator
			var operator string
			if i-1 < len(ctx.AllGT_STMT()) {
				operator = ">"
			} else if i-1 < len(ctx.AllGT_STMT())+len(ctx.AllLT_STMT()) {
				operator = "<"
			} else if i-1 < len(ctx.AllGT_STMT())+len(ctx.AllLT_STMT())+len(ctx.AllGE_STMT()) {
				operator = ">="
			} else if i-1 < len(ctx.AllGT_STMT())+len(ctx.AllLT_STMT())+len(ctx.AllGE_STMT())+len(ctx.AllLE_STMT()) {
				operator = "<="
			}

			// Convert to comparable types
			left, right := v.toComparable(result, value)

			switch operator {
			case ">":
				result = left > right
			case "<":
				result = left < right
			case ">=":
				result = left >= right
			case "<=":
				result = left <= right
			}
		}
	}

	return result, nil
}

// VisitPrimaryExpressionStmt visits primary expressions in statement mode
func (v *TemplateVisitor) VisitPrimaryExpressionStmt(ctx *generated.PrimaryExpressionStmtContext) (any, error) {

	// Handle NOT operator
	if ctx.NOT_STMT() != nil && ctx.PrimaryExpressionStmt() != nil {
		value, err := v.VisitPrimaryExpressionStmt(ctx.PrimaryExpressionStmt().(*generated.PrimaryExpressionStmtContext))
		if err != nil {
			return nil, err
		}
		return !isTrue(value), nil
	}

	// Handle basic expressions
	if ctx.NUMBER_STMT() != nil {
		text := ctx.NUMBER_STMT().GetText()
		if strings.Contains(text, ".") {
			if num, err := strconv.ParseFloat(text, 64); err == nil {
				return num, nil
			}
		}
		if num, err := strconv.Atoi(text); err == nil {
			return num, nil
		}
		return text, nil
	}

	if ctx.STRING_STMT() != nil {
		text := ctx.STRING_STMT().GetText()
		return strings.Trim(text, `"'`), nil
	}

	if ctx.BOOLEAN_STMT() != nil {
		text := ctx.BOOLEAN_STMT().GetText()
		return text == "true", nil
	}

	if ctx.IDENTIFIER_STMT() != nil {
		name := ctx.IDENTIFIER_STMT().GetText()
		if value, exists := v.context.GetVariable(name); exists {
			return value, nil
		}
		return name, nil
	}

	if ctx.LPAREN_STMT() != nil && ctx.Expression() != nil {
		return v.VisitExpression(ctx.Expression().(*generated.ExpressionContext))
	}

	return "", nil
}

// callMacro calls a macro with the given arguments
func (v *TemplateVisitor) callMacro(macro *Macro, expressions []generated.IExpressionContext) (any, error) {
	// Collect arguments
	var args []any
	for _, expr := range expressions {
		if arg, err := v.VisitExpression(expr.(*generated.ExpressionContext)); err != nil {
			return nil, err
		} else {
			args = append(args, arg)
		}
	}

	// Check argument count
	if len(args) != len(macro.Parameters) {
		return nil, fmt.Errorf("macro '%s' expects %d arguments, got %d", macro.Name, len(macro.Parameters), len(args))
	}

	// Create child context for macro execution
	childContext := &Context{
		Variables: make(map[string]any),
		Functions: v.context.Functions,
		Macros:    v.context.Macros,
		Parent:    v.context,
	}

	// Set macro parameters as variables in child context
	for i, param := range macro.Parameters {
		childContext.SetVariable(param, args[i])
	}

	// Save original visitor context
	originalContext := v.context
	v.context = childContext
	defer func() { v.context = originalContext }()

	// Execute macro template
	return v.VisitTemplate(macro.Template)
}

// VisitMacroStatement visits macro statements
func (v *TemplateVisitor) VisitMacroStatement(ctx *generated.MacroStatementContext) (any, error) {
	// Extract macro name and parameters
	macroName := ctx.IDENTIFIER_STMT(0).GetText()

	var parameters []string
	for i := 1; i < len(ctx.AllIDENTIFIER_STMT()); i++ {
		parameters = append(parameters, ctx.IDENTIFIER_STMT(i).GetText())
	}

	// Create macro definition
	macro := &Macro{
		Name:       macroName,
		Parameters: parameters,
		Template:   ctx.Template().(*generated.TemplateContext),
	}

	// Register macro in context
	v.context.SetMacro(macroName, macro)

	return "", nil
}

// isTrue checks if a value is truthy
func isTrue(value any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return true
	}
}

func Parse(template string) {
	stream := antlr.NewInputStream(template)
	lex := generated.NewunsafeLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	parser := generated.NewunsafeParser(tokens)

	fmt.Println(parser.Template().ToStringTree(nil, parser))
}
