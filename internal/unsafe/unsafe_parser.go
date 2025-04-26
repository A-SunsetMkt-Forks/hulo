// Code generated from UnsafeParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package unsafe // UnsafeParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type UnsafeParser struct {
	*antlr.BaseParser
}

var UnsafeParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func unsafeparserParserInit() {
	staticData := &UnsafeParserParserStaticData
	staticData.LiteralNames = []string{
		"", "'{{'", "", "'}}'", "'if'", "'loop'", "'end'", "'in'", "'('", "')'",
		"'$'", "'|'", "','",
	}
	staticData.SymbolicNames = []string{
		"", "DOUBLE_LBRACE", "TEXT", "DOUBLE_RBRACE", "IF", "LOOP", "END", "IN",
		"LPAREN", "RPAREN", "DOLLAR", "PIPE", "COMMA", "IDENTIFIER", "STRING",
		"NUMBER", "WS",
	}
	staticData.RuleNames = []string{
		"template", "content", "statement", "ifStatement", "loopStatement",
		"expressionStatement", "expression", "pipelineExpr", "primaryExpr",
		"functionCall",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 16, 100, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 1,
		0, 5, 0, 23, 8, 0, 10, 0, 12, 0, 26, 9, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2,
		3, 2, 33, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1,
		5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 5, 7, 64, 8, 7, 10, 7, 12,
		7, 67, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 76, 8, 8,
		1, 9, 1, 9, 5, 9, 80, 8, 9, 10, 9, 12, 9, 83, 9, 9, 1, 9, 1, 9, 1, 9, 3,
		9, 88, 8, 9, 1, 9, 1, 9, 5, 9, 92, 8, 9, 10, 9, 12, 9, 95, 9, 9, 1, 9,
		3, 9, 98, 8, 9, 1, 9, 0, 0, 10, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 0, 0,
		101, 0, 24, 1, 0, 0, 0, 2, 27, 1, 0, 0, 0, 4, 32, 1, 0, 0, 0, 6, 34, 1,
		0, 0, 0, 8, 43, 1, 0, 0, 0, 10, 54, 1, 0, 0, 0, 12, 58, 1, 0, 0, 0, 14,
		60, 1, 0, 0, 0, 16, 75, 1, 0, 0, 0, 18, 97, 1, 0, 0, 0, 20, 23, 3, 2, 1,
		0, 21, 23, 3, 4, 2, 0, 22, 20, 1, 0, 0, 0, 22, 21, 1, 0, 0, 0, 23, 26,
		1, 0, 0, 0, 24, 22, 1, 0, 0, 0, 24, 25, 1, 0, 0, 0, 25, 1, 1, 0, 0, 0,
		26, 24, 1, 0, 0, 0, 27, 28, 5, 2, 0, 0, 28, 3, 1, 0, 0, 0, 29, 33, 3, 6,
		3, 0, 30, 33, 3, 8, 4, 0, 31, 33, 3, 10, 5, 0, 32, 29, 1, 0, 0, 0, 32,
		30, 1, 0, 0, 0, 32, 31, 1, 0, 0, 0, 33, 5, 1, 0, 0, 0, 34, 35, 5, 1, 0,
		0, 35, 36, 5, 4, 0, 0, 36, 37, 3, 12, 6, 0, 37, 38, 5, 3, 0, 0, 38, 39,
		3, 0, 0, 0, 39, 40, 5, 1, 0, 0, 40, 41, 5, 6, 0, 0, 41, 42, 5, 3, 0, 0,
		42, 7, 1, 0, 0, 0, 43, 44, 5, 1, 0, 0, 44, 45, 5, 5, 0, 0, 45, 46, 5, 13,
		0, 0, 46, 47, 5, 7, 0, 0, 47, 48, 5, 13, 0, 0, 48, 49, 5, 3, 0, 0, 49,
		50, 3, 0, 0, 0, 50, 51, 5, 1, 0, 0, 51, 52, 5, 6, 0, 0, 52, 53, 5, 3, 0,
		0, 53, 9, 1, 0, 0, 0, 54, 55, 5, 1, 0, 0, 55, 56, 3, 12, 6, 0, 56, 57,
		5, 3, 0, 0, 57, 11, 1, 0, 0, 0, 58, 59, 3, 14, 7, 0, 59, 13, 1, 0, 0, 0,
		60, 65, 3, 16, 8, 0, 61, 62, 5, 11, 0, 0, 62, 64, 3, 18, 9, 0, 63, 61,
		1, 0, 0, 0, 64, 67, 1, 0, 0, 0, 65, 63, 1, 0, 0, 0, 65, 66, 1, 0, 0, 0,
		66, 15, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 68, 76, 5, 15, 0, 0, 69, 76, 5,
		14, 0, 0, 70, 76, 5, 13, 0, 0, 71, 72, 5, 8, 0, 0, 72, 73, 3, 12, 6, 0,
		73, 74, 5, 9, 0, 0, 74, 76, 1, 0, 0, 0, 75, 68, 1, 0, 0, 0, 75, 69, 1,
		0, 0, 0, 75, 70, 1, 0, 0, 0, 75, 71, 1, 0, 0, 0, 76, 17, 1, 0, 0, 0, 77,
		81, 5, 13, 0, 0, 78, 80, 3, 16, 8, 0, 79, 78, 1, 0, 0, 0, 80, 83, 1, 0,
		0, 0, 81, 79, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 98, 1, 0, 0, 0, 83, 81,
		1, 0, 0, 0, 84, 85, 5, 13, 0, 0, 85, 87, 5, 8, 0, 0, 86, 88, 3, 16, 8,
		0, 87, 86, 1, 0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 93, 1, 0, 0, 0, 89, 90,
		5, 12, 0, 0, 90, 92, 3, 16, 8, 0, 91, 89, 1, 0, 0, 0, 92, 95, 1, 0, 0,
		0, 93, 91, 1, 0, 0, 0, 93, 94, 1, 0, 0, 0, 94, 96, 1, 0, 0, 0, 95, 93,
		1, 0, 0, 0, 96, 98, 5, 9, 0, 0, 97, 77, 1, 0, 0, 0, 97, 84, 1, 0, 0, 0,
		98, 19, 1, 0, 0, 0, 9, 22, 24, 32, 65, 75, 81, 87, 93, 97,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// UnsafeParserInit initializes any static state used to implement UnsafeParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewUnsafeParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func UnsafeParserInit() {
	staticData := &UnsafeParserParserStaticData
	staticData.once.Do(unsafeparserParserInit)
}

// NewUnsafeParser produces a new parser instance for the optional input antlr.TokenStream.
func NewUnsafeParser(input antlr.TokenStream) *UnsafeParser {
	UnsafeParserInit()
	this := new(UnsafeParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &UnsafeParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "UnsafeParser.g4"

	return this
}

// UnsafeParser tokens.
const (
	UnsafeParserEOF           = antlr.TokenEOF
	UnsafeParserDOUBLE_LBRACE = 1
	UnsafeParserTEXT          = 2
	UnsafeParserDOUBLE_RBRACE = 3
	UnsafeParserIF            = 4
	UnsafeParserLOOP          = 5
	UnsafeParserEND           = 6
	UnsafeParserIN            = 7
	UnsafeParserLPAREN        = 8
	UnsafeParserRPAREN        = 9
	UnsafeParserDOLLAR        = 10
	UnsafeParserPIPE          = 11
	UnsafeParserCOMMA         = 12
	UnsafeParserIDENTIFIER    = 13
	UnsafeParserSTRING        = 14
	UnsafeParserNUMBER        = 15
	UnsafeParserWS            = 16
)

// UnsafeParser rules.
const (
	UnsafeParserRULE_template            = 0
	UnsafeParserRULE_content             = 1
	UnsafeParserRULE_statement           = 2
	UnsafeParserRULE_ifStatement         = 3
	UnsafeParserRULE_loopStatement       = 4
	UnsafeParserRULE_expressionStatement = 5
	UnsafeParserRULE_expression          = 6
	UnsafeParserRULE_pipelineExpr        = 7
	UnsafeParserRULE_primaryExpr         = 8
	UnsafeParserRULE_functionCall        = 9
)

// ITemplateContext is an interface to support dynamic dispatch.
type ITemplateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllContent() []IContentContext
	Content(i int) IContentContext
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsTemplateContext differentiates from other interfaces.
	IsTemplateContext()
}

type TemplateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTemplateContext() *TemplateContext {
	var p = new(TemplateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_template
	return p
}

func InitEmptyTemplateContext(p *TemplateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_template
}

func (*TemplateContext) IsTemplateContext() {}

func NewTemplateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TemplateContext {
	var p = new(TemplateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_template

	return p
}

func (s *TemplateContext) GetParser() antlr.Parser { return s.parser }

func (s *TemplateContext) AllContent() []IContentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IContentContext); ok {
			len++
		}
	}

	tst := make([]IContentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IContentContext); ok {
			tst[i] = t.(IContentContext)
			i++
		}
	}

	return tst
}

func (s *TemplateContext) Content(i int) IContentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContentContext)
}

func (s *TemplateContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *TemplateContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *TemplateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TemplateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TemplateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitTemplate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) Template() (localctx ITemplateContext) {
	localctx = NewTemplateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, UnsafeParserRULE_template)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(24)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(22)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case UnsafeParserTEXT:
				{
					p.SetState(20)
					p.Content()
				}

			case UnsafeParserDOUBLE_LBRACE:
				{
					p.SetState(21)
					p.Statement()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

		}
		p.SetState(26)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IContentContext is an interface to support dynamic dispatch.
type IContentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TEXT() antlr.TerminalNode

	// IsContentContext differentiates from other interfaces.
	IsContentContext()
}

type ContentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContentContext() *ContentContext {
	var p = new(ContentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_content
	return p
}

func InitEmptyContentContext(p *ContentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_content
}

func (*ContentContext) IsContentContext() {}

func NewContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContentContext {
	var p = new(ContentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_content

	return p
}

func (s *ContentContext) GetParser() antlr.Parser { return s.parser }

func (s *ContentContext) TEXT() antlr.TerminalNode {
	return s.GetToken(UnsafeParserTEXT, 0)
}

func (s *ContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ContentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitContent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) Content() (localctx IContentContext) {
	localctx = NewContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, UnsafeParserRULE_content)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(27)
		p.Match(UnsafeParserTEXT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IfStatement() IIfStatementContext
	LoopStatement() ILoopStatementContext
	ExpressionStatement() IExpressionStatementContext

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_statement
	return p
}

func InitEmptyStatementContext(p *StatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_statement
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) IfStatement() IIfStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfStatementContext)
}

func (s *StatementContext) LoopStatement() ILoopStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoopStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoopStatementContext)
}

func (s *StatementContext) ExpressionStatement() IExpressionStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionStatementContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, UnsafeParserRULE_statement)
	p.SetState(32)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(29)
			p.IfStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(30)
			p.LoopStatement()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(31)
			p.ExpressionStatement()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIfStatementContext is an interface to support dynamic dispatch.
type IIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDOUBLE_LBRACE() []antlr.TerminalNode
	DOUBLE_LBRACE(i int) antlr.TerminalNode
	IF() antlr.TerminalNode
	Expression() IExpressionContext
	AllDOUBLE_RBRACE() []antlr.TerminalNode
	DOUBLE_RBRACE(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode

	// IsIfStatementContext differentiates from other interfaces.
	IsIfStatementContext()
}

type IfStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIfStatementContext() *IfStatementContext {
	var p = new(IfStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_ifStatement
	return p
}

func InitEmptyIfStatementContext(p *IfStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_ifStatement
}

func (*IfStatementContext) IsIfStatementContext() {}

func NewIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfStatementContext {
	var p = new(IfStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_ifStatement

	return p
}

func (s *IfStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *IfStatementContext) AllDOUBLE_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserDOUBLE_LBRACE)
}

func (s *IfStatementContext) DOUBLE_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserDOUBLE_LBRACE, i)
}

func (s *IfStatementContext) IF() antlr.TerminalNode {
	return s.GetToken(UnsafeParserIF, 0)
}

func (s *IfStatementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *IfStatementContext) AllDOUBLE_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserDOUBLE_RBRACE)
}

func (s *IfStatementContext) DOUBLE_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserDOUBLE_RBRACE, i)
}

func (s *IfStatementContext) Template() ITemplateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITemplateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITemplateContext)
}

func (s *IfStatementContext) END() antlr.TerminalNode {
	return s.GetToken(UnsafeParserEND, 0)
}

func (s *IfStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IfStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitIfStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) IfStatement() (localctx IIfStatementContext) {
	localctx = NewIfStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, UnsafeParserRULE_ifStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(34)
		p.Match(UnsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(35)
		p.Match(UnsafeParserIF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(36)
		p.Expression()
	}
	{
		p.SetState(37)
		p.Match(UnsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(38)
		p.Template()
	}
	{
		p.SetState(39)
		p.Match(UnsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(40)
		p.Match(UnsafeParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(41)
		p.Match(UnsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILoopStatementContext is an interface to support dynamic dispatch.
type ILoopStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDOUBLE_LBRACE() []antlr.TerminalNode
	DOUBLE_LBRACE(i int) antlr.TerminalNode
	LOOP() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	IN() antlr.TerminalNode
	AllDOUBLE_RBRACE() []antlr.TerminalNode
	DOUBLE_RBRACE(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode

	// IsLoopStatementContext differentiates from other interfaces.
	IsLoopStatementContext()
}

type LoopStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLoopStatementContext() *LoopStatementContext {
	var p = new(LoopStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_loopStatement
	return p
}

func InitEmptyLoopStatementContext(p *LoopStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_loopStatement
}

func (*LoopStatementContext) IsLoopStatementContext() {}

func NewLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoopStatementContext {
	var p = new(LoopStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_loopStatement

	return p
}

func (s *LoopStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *LoopStatementContext) AllDOUBLE_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserDOUBLE_LBRACE)
}

func (s *LoopStatementContext) DOUBLE_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserDOUBLE_LBRACE, i)
}

func (s *LoopStatementContext) LOOP() antlr.TerminalNode {
	return s.GetToken(UnsafeParserLOOP, 0)
}

func (s *LoopStatementContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserIDENTIFIER)
}

func (s *LoopStatementContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserIDENTIFIER, i)
}

func (s *LoopStatementContext) IN() antlr.TerminalNode {
	return s.GetToken(UnsafeParserIN, 0)
}

func (s *LoopStatementContext) AllDOUBLE_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserDOUBLE_RBRACE)
}

func (s *LoopStatementContext) DOUBLE_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserDOUBLE_RBRACE, i)
}

func (s *LoopStatementContext) Template() ITemplateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITemplateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITemplateContext)
}

func (s *LoopStatementContext) END() antlr.TerminalNode {
	return s.GetToken(UnsafeParserEND, 0)
}

func (s *LoopStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LoopStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitLoopStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) LoopStatement() (localctx ILoopStatementContext) {
	localctx = NewLoopStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, UnsafeParserRULE_loopStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		p.Match(UnsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(44)
		p.Match(UnsafeParserLOOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(45)
		p.Match(UnsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(46)
		p.Match(UnsafeParserIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(47)
		p.Match(UnsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(48)
		p.Match(UnsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(49)
		p.Template()
	}
	{
		p.SetState(50)
		p.Match(UnsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(51)
		p.Match(UnsafeParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(52)
		p.Match(UnsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionStatementContext is an interface to support dynamic dispatch.
type IExpressionStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOUBLE_LBRACE() antlr.TerminalNode
	Expression() IExpressionContext
	DOUBLE_RBRACE() antlr.TerminalNode

	// IsExpressionStatementContext differentiates from other interfaces.
	IsExpressionStatementContext()
}

type ExpressionStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionStatementContext() *ExpressionStatementContext {
	var p = new(ExpressionStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_expressionStatement
	return p
}

func InitEmptyExpressionStatementContext(p *ExpressionStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_expressionStatement
}

func (*ExpressionStatementContext) IsExpressionStatementContext() {}

func NewExpressionStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionStatementContext {
	var p = new(ExpressionStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_expressionStatement

	return p
}

func (s *ExpressionStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionStatementContext) DOUBLE_LBRACE() antlr.TerminalNode {
	return s.GetToken(UnsafeParserDOUBLE_LBRACE, 0)
}

func (s *ExpressionStatementContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionStatementContext) DOUBLE_RBRACE() antlr.TerminalNode {
	return s.GetToken(UnsafeParserDOUBLE_RBRACE, 0)
}

func (s *ExpressionStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitExpressionStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) ExpressionStatement() (localctx IExpressionStatementContext) {
	localctx = NewExpressionStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, UnsafeParserRULE_expressionStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.Match(UnsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(55)
		p.Expression()
	}
	{
		p.SetState(56)
		p.Match(UnsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PipelineExpr() IPipelineExprContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) PipelineExpr() IPipelineExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPipelineExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPipelineExprContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, UnsafeParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.PipelineExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPipelineExprContext is an interface to support dynamic dispatch.
type IPipelineExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrimaryExpr() IPrimaryExprContext
	AllPIPE() []antlr.TerminalNode
	PIPE(i int) antlr.TerminalNode
	AllFunctionCall() []IFunctionCallContext
	FunctionCall(i int) IFunctionCallContext

	// IsPipelineExprContext differentiates from other interfaces.
	IsPipelineExprContext()
}

type PipelineExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPipelineExprContext() *PipelineExprContext {
	var p = new(PipelineExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_pipelineExpr
	return p
}

func InitEmptyPipelineExprContext(p *PipelineExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_pipelineExpr
}

func (*PipelineExprContext) IsPipelineExprContext() {}

func NewPipelineExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PipelineExprContext {
	var p = new(PipelineExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_pipelineExpr

	return p
}

func (s *PipelineExprContext) GetParser() antlr.Parser { return s.parser }

func (s *PipelineExprContext) PrimaryExpr() IPrimaryExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExprContext)
}

func (s *PipelineExprContext) AllPIPE() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserPIPE)
}

func (s *PipelineExprContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserPIPE, i)
}

func (s *PipelineExprContext) AllFunctionCall() []IFunctionCallContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFunctionCallContext); ok {
			len++
		}
	}

	tst := make([]IFunctionCallContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFunctionCallContext); ok {
			tst[i] = t.(IFunctionCallContext)
			i++
		}
	}

	return tst
}

func (s *PipelineExprContext) FunctionCall(i int) IFunctionCallContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *PipelineExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PipelineExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PipelineExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitPipelineExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) PipelineExpr() (localctx IPipelineExprContext) {
	localctx = NewPipelineExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, UnsafeParserRULE_pipelineExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(60)
		p.PrimaryExpr()
	}
	p.SetState(65)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == UnsafeParserPIPE {
		{
			p.SetState(61)
			p.Match(UnsafeParserPIPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(62)
			p.FunctionCall()
		}

		p.SetState(67)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimaryExprContext is an interface to support dynamic dispatch.
type IPrimaryExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode
	STRING() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode

	// IsPrimaryExprContext differentiates from other interfaces.
	IsPrimaryExprContext()
}

type PrimaryExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExprContext() *PrimaryExprContext {
	var p = new(PrimaryExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_primaryExpr
	return p
}

func InitEmptyPrimaryExprContext(p *PrimaryExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_primaryExpr
}

func (*PrimaryExprContext) IsPrimaryExprContext() {}

func NewPrimaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExprContext {
	var p = new(PrimaryExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_primaryExpr

	return p
}

func (s *PrimaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExprContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(UnsafeParserNUMBER, 0)
}

func (s *PrimaryExprContext) STRING() antlr.TerminalNode {
	return s.GetToken(UnsafeParserSTRING, 0)
}

func (s *PrimaryExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(UnsafeParserIDENTIFIER, 0)
}

func (s *PrimaryExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(UnsafeParserLPAREN, 0)
}

func (s *PrimaryExprContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PrimaryExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(UnsafeParserRPAREN, 0)
}

func (s *PrimaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitPrimaryExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) PrimaryExpr() (localctx IPrimaryExprContext) {
	localctx = NewPrimaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, UnsafeParserRULE_primaryExpr)
	p.SetState(75)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UnsafeParserNUMBER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(68)
			p.Match(UnsafeParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UnsafeParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(69)
			p.Match(UnsafeParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UnsafeParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(70)
			p.Match(UnsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UnsafeParserLPAREN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(71)
			p.Match(UnsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(72)
			p.Expression()
		}
		{
			p.SetState(73)
			p.Match(UnsafeParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	AllPrimaryExpr() []IPrimaryExprContext
	PrimaryExpr(i int) IPrimaryExprContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UnsafeParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UnsafeParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(UnsafeParserIDENTIFIER, 0)
}

func (s *FunctionCallContext) AllPrimaryExpr() []IPrimaryExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPrimaryExprContext); ok {
			len++
		}
	}

	tst := make([]IPrimaryExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPrimaryExprContext); ok {
			tst[i] = t.(IPrimaryExprContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) PrimaryExpr(i int) IPrimaryExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExprContext)
}

func (s *FunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(UnsafeParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(UnsafeParserRPAREN, 0)
}

func (s *FunctionCallContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(UnsafeParserCOMMA)
}

func (s *FunctionCallContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(UnsafeParserCOMMA, i)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UnsafeParserVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UnsafeParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, UnsafeParserRULE_functionCall)
	var _la int

	p.SetState(97)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(77)
			p.Match(UnsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(81)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&57600) != 0 {
			{
				p.SetState(78)
				p.PrimaryExpr()
			}

			p.SetState(83)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(84)
			p.Match(UnsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(85)
			p.Match(UnsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(87)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&57600) != 0 {
			{
				p.SetState(86)
				p.PrimaryExpr()
			}

		}
		p.SetState(93)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == UnsafeParserCOMMA {
			{
				p.SetState(89)
				p.Match(UnsafeParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(90)
				p.PrimaryExpr()
			}

			p.SetState(95)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(96)
			p.Match(UnsafeParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
