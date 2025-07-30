// Code generated from unsafeParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // unsafeParser
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

type unsafeParser struct {
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
		"", "'{{'", "", "'}}'", "'if'", "'loop'", "'end'", "'in'", "'define'",
		"'template'", "'('", "')'", "'$'", "'|'", "','", "':='",
	}
	staticData.SymbolicNames = []string{
		"", "DOUBLE_LBRACE", "TEXT", "DOUBLE_RBRACE", "IF", "LOOP", "END", "IN",
		"DEFINE", "TEMPLATE", "LPAREN", "RPAREN", "DOLLAR", "PIPE", "COMMA",
		"COLON_EQUAL", "IDENTIFIER", "STRING", "NUMBER", "WS",
	}
	staticData.RuleNames = []string{
		"template", "content", "statement", "variableStatement", "ifStatement",
		"loopStatement", "expressionStatement", "defineStatement", "templateStatement",
		"expression", "pipelineExpr", "primaryExpr", "functionCall",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 19, 147, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 1, 0, 1, 0, 5, 0, 29, 8, 0, 10, 0, 12,
		0, 32, 9, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 42,
		8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5,
		1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7,
		5, 7, 79, 8, 7, 10, 7, 12, 7, 82, 9, 7, 3, 7, 84, 8, 7, 1, 7, 3, 7, 87,
		8, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 5, 8,
		99, 8, 8, 10, 8, 12, 8, 102, 9, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10,
		1, 10, 5, 10, 111, 8, 10, 10, 10, 12, 10, 114, 9, 10, 1, 11, 1, 11, 1,
		11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 123, 8, 11, 1, 12, 1, 12, 5, 12,
		127, 8, 12, 10, 12, 12, 12, 130, 9, 12, 1, 12, 1, 12, 1, 12, 3, 12, 135,
		8, 12, 1, 12, 1, 12, 5, 12, 139, 8, 12, 10, 12, 12, 12, 142, 9, 12, 1,
		12, 3, 12, 145, 8, 12, 1, 12, 0, 0, 13, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 0, 0, 152, 0, 30, 1, 0, 0, 0, 2, 33, 1, 0, 0, 0, 4, 41,
		1, 0, 0, 0, 6, 43, 1, 0, 0, 0, 8, 47, 1, 0, 0, 0, 10, 56, 1, 0, 0, 0, 12,
		67, 1, 0, 0, 0, 14, 71, 1, 0, 0, 0, 16, 94, 1, 0, 0, 0, 18, 105, 1, 0,
		0, 0, 20, 107, 1, 0, 0, 0, 22, 122, 1, 0, 0, 0, 24, 144, 1, 0, 0, 0, 26,
		29, 3, 2, 1, 0, 27, 29, 3, 4, 2, 0, 28, 26, 1, 0, 0, 0, 28, 27, 1, 0, 0,
		0, 29, 32, 1, 0, 0, 0, 30, 28, 1, 0, 0, 0, 30, 31, 1, 0, 0, 0, 31, 1, 1,
		0, 0, 0, 32, 30, 1, 0, 0, 0, 33, 34, 5, 2, 0, 0, 34, 3, 1, 0, 0, 0, 35,
		42, 3, 8, 4, 0, 36, 42, 3, 10, 5, 0, 37, 42, 3, 12, 6, 0, 38, 42, 3, 14,
		7, 0, 39, 42, 3, 16, 8, 0, 40, 42, 3, 6, 3, 0, 41, 35, 1, 0, 0, 0, 41,
		36, 1, 0, 0, 0, 41, 37, 1, 0, 0, 0, 41, 38, 1, 0, 0, 0, 41, 39, 1, 0, 0,
		0, 41, 40, 1, 0, 0, 0, 42, 5, 1, 0, 0, 0, 43, 44, 5, 16, 0, 0, 44, 45,
		5, 15, 0, 0, 45, 46, 3, 18, 9, 0, 46, 7, 1, 0, 0, 0, 47, 48, 5, 1, 0, 0,
		48, 49, 5, 4, 0, 0, 49, 50, 3, 18, 9, 0, 50, 51, 5, 3, 0, 0, 51, 52, 3,
		0, 0, 0, 52, 53, 5, 1, 0, 0, 53, 54, 5, 6, 0, 0, 54, 55, 5, 3, 0, 0, 55,
		9, 1, 0, 0, 0, 56, 57, 5, 1, 0, 0, 57, 58, 5, 5, 0, 0, 58, 59, 5, 16, 0,
		0, 59, 60, 5, 7, 0, 0, 60, 61, 5, 16, 0, 0, 61, 62, 5, 3, 0, 0, 62, 63,
		3, 0, 0, 0, 63, 64, 5, 1, 0, 0, 64, 65, 5, 6, 0, 0, 65, 66, 5, 3, 0, 0,
		66, 11, 1, 0, 0, 0, 67, 68, 5, 1, 0, 0, 68, 69, 3, 18, 9, 0, 69, 70, 5,
		3, 0, 0, 70, 13, 1, 0, 0, 0, 71, 72, 5, 1, 0, 0, 72, 73, 5, 8, 0, 0, 73,
		86, 5, 16, 0, 0, 74, 83, 5, 10, 0, 0, 75, 80, 5, 16, 0, 0, 76, 77, 5, 14,
		0, 0, 77, 79, 5, 16, 0, 0, 78, 76, 1, 0, 0, 0, 79, 82, 1, 0, 0, 0, 80,
		78, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81, 84, 1, 0, 0, 0, 82, 80, 1, 0, 0,
		0, 83, 75, 1, 0, 0, 0, 83, 84, 1, 0, 0, 0, 84, 85, 1, 0, 0, 0, 85, 87,
		5, 11, 0, 0, 86, 74, 1, 0, 0, 0, 86, 87, 1, 0, 0, 0, 87, 88, 1, 0, 0, 0,
		88, 89, 5, 3, 0, 0, 89, 90, 3, 0, 0, 0, 90, 91, 5, 1, 0, 0, 91, 92, 5,
		6, 0, 0, 92, 93, 5, 3, 0, 0, 93, 15, 1, 0, 0, 0, 94, 95, 5, 1, 0, 0, 95,
		96, 5, 9, 0, 0, 96, 100, 5, 16, 0, 0, 97, 99, 3, 18, 9, 0, 98, 97, 1, 0,
		0, 0, 99, 102, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 101,
		103, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 103, 104, 5, 3, 0, 0, 104, 17, 1,
		0, 0, 0, 105, 106, 3, 20, 10, 0, 106, 19, 1, 0, 0, 0, 107, 112, 3, 22,
		11, 0, 108, 109, 5, 13, 0, 0, 109, 111, 3, 24, 12, 0, 110, 108, 1, 0, 0,
		0, 111, 114, 1, 0, 0, 0, 112, 110, 1, 0, 0, 0, 112, 113, 1, 0, 0, 0, 113,
		21, 1, 0, 0, 0, 114, 112, 1, 0, 0, 0, 115, 123, 5, 18, 0, 0, 116, 123,
		5, 17, 0, 0, 117, 123, 5, 16, 0, 0, 118, 119, 5, 10, 0, 0, 119, 120, 3,
		18, 9, 0, 120, 121, 5, 11, 0, 0, 121, 123, 1, 0, 0, 0, 122, 115, 1, 0,
		0, 0, 122, 116, 1, 0, 0, 0, 122, 117, 1, 0, 0, 0, 122, 118, 1, 0, 0, 0,
		123, 23, 1, 0, 0, 0, 124, 128, 5, 16, 0, 0, 125, 127, 3, 22, 11, 0, 126,
		125, 1, 0, 0, 0, 127, 130, 1, 0, 0, 0, 128, 126, 1, 0, 0, 0, 128, 129,
		1, 0, 0, 0, 129, 145, 1, 0, 0, 0, 130, 128, 1, 0, 0, 0, 131, 132, 5, 16,
		0, 0, 132, 134, 5, 10, 0, 0, 133, 135, 3, 22, 11, 0, 134, 133, 1, 0, 0,
		0, 134, 135, 1, 0, 0, 0, 135, 140, 1, 0, 0, 0, 136, 137, 5, 14, 0, 0, 137,
		139, 3, 22, 11, 0, 138, 136, 1, 0, 0, 0, 139, 142, 1, 0, 0, 0, 140, 138,
		1, 0, 0, 0, 140, 141, 1, 0, 0, 0, 141, 143, 1, 0, 0, 0, 142, 140, 1, 0,
		0, 0, 143, 145, 5, 11, 0, 0, 144, 124, 1, 0, 0, 0, 144, 131, 1, 0, 0, 0,
		145, 25, 1, 0, 0, 0, 13, 28, 30, 41, 80, 83, 86, 100, 112, 122, 128, 134,
		140, 144,
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

// unsafeParserInit initializes any static state used to implement unsafeParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewunsafeParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func UnsafeParserInit() {
	staticData := &UnsafeParserParserStaticData
	staticData.once.Do(unsafeparserParserInit)
}

// NewunsafeParser produces a new parser instance for the optional input antlr.TokenStream.
func NewunsafeParser(input antlr.TokenStream) *unsafeParser {
	UnsafeParserInit()
	this := new(unsafeParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &UnsafeParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "unsafeParser.g4"

	return this
}

// unsafeParser tokens.
const (
	unsafeParserEOF           = antlr.TokenEOF
	unsafeParserDOUBLE_LBRACE = 1
	unsafeParserTEXT          = 2
	unsafeParserDOUBLE_RBRACE = 3
	unsafeParserIF            = 4
	unsafeParserLOOP          = 5
	unsafeParserEND           = 6
	unsafeParserIN            = 7
	unsafeParserDEFINE        = 8
	unsafeParserTEMPLATE      = 9
	unsafeParserLPAREN        = 10
	unsafeParserRPAREN        = 11
	unsafeParserDOLLAR        = 12
	unsafeParserPIPE          = 13
	unsafeParserCOMMA         = 14
	unsafeParserCOLON_EQUAL   = 15
	unsafeParserIDENTIFIER    = 16
	unsafeParserSTRING        = 17
	unsafeParserNUMBER        = 18
	unsafeParserWS            = 19
)

// unsafeParser rules.
const (
	unsafeParserRULE_template            = 0
	unsafeParserRULE_content             = 1
	unsafeParserRULE_statement           = 2
	unsafeParserRULE_variableStatement   = 3
	unsafeParserRULE_ifStatement         = 4
	unsafeParserRULE_loopStatement       = 5
	unsafeParserRULE_expressionStatement = 6
	unsafeParserRULE_defineStatement     = 7
	unsafeParserRULE_templateStatement   = 8
	unsafeParserRULE_expression          = 9
	unsafeParserRULE_pipelineExpr        = 10
	unsafeParserRULE_primaryExpr         = 11
	unsafeParserRULE_functionCall        = 12
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
	p.RuleIndex = unsafeParserRULE_template
	return p
}

func InitEmptyTemplateContext(p *TemplateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_template
}

func (*TemplateContext) IsTemplateContext() {}

func NewTemplateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TemplateContext {
	var p = new(TemplateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_template

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
	case unsafeParserVisitor:
		return t.VisitTemplate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) Template() (localctx ITemplateContext) {
	localctx = NewTemplateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, unsafeParserRULE_template)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(30)
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
			p.SetState(28)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case unsafeParserTEXT:
				{
					p.SetState(26)
					p.Content()
				}

			case unsafeParserDOUBLE_LBRACE, unsafeParserIDENTIFIER:
				{
					p.SetState(27)
					p.Statement()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

		}
		p.SetState(32)
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
	p.RuleIndex = unsafeParserRULE_content
	return p
}

func InitEmptyContentContext(p *ContentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_content
}

func (*ContentContext) IsContentContext() {}

func NewContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContentContext {
	var p = new(ContentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_content

	return p
}

func (s *ContentContext) GetParser() antlr.Parser { return s.parser }

func (s *ContentContext) TEXT() antlr.TerminalNode {
	return s.GetToken(unsafeParserTEXT, 0)
}

func (s *ContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ContentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitContent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) Content() (localctx IContentContext) {
	localctx = NewContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, unsafeParserRULE_content)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(33)
		p.Match(unsafeParserTEXT)
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
	DefineStatement() IDefineStatementContext
	TemplateStatement() ITemplateStatementContext
	VariableStatement() IVariableStatementContext

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
	p.RuleIndex = unsafeParserRULE_statement
	return p
}

func InitEmptyStatementContext(p *StatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_statement
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_statement

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

func (s *StatementContext) DefineStatement() IDefineStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDefineStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDefineStatementContext)
}

func (s *StatementContext) TemplateStatement() ITemplateStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITemplateStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITemplateStatementContext)
}

func (s *StatementContext) VariableStatement() IVariableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableStatementContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, unsafeParserRULE_statement)
	p.SetState(41)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(35)
			p.IfStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(36)
			p.LoopStatement()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(37)
			p.ExpressionStatement()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(38)
			p.DefineStatement()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(39)
			p.TemplateStatement()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(40)
			p.VariableStatement()
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

// IVariableStatementContext is an interface to support dynamic dispatch.
type IVariableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON_EQUAL() antlr.TerminalNode
	Expression() IExpressionContext

	// IsVariableStatementContext differentiates from other interfaces.
	IsVariableStatementContext()
}

type VariableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableStatementContext() *VariableStatementContext {
	var p = new(VariableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_variableStatement
	return p
}

func InitEmptyVariableStatementContext(p *VariableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_variableStatement
}

func (*VariableStatementContext) IsVariableStatementContext() {}

func NewVariableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableStatementContext {
	var p = new(VariableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_variableStatement

	return p
}

func (s *VariableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableStatementContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *VariableStatementContext) COLON_EQUAL() antlr.TerminalNode {
	return s.GetToken(unsafeParserCOLON_EQUAL, 0)
}

func (s *VariableStatementContext) Expression() IExpressionContext {
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

func (s *VariableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitVariableStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) VariableStatement() (localctx IVariableStatementContext) {
	localctx = NewVariableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, unsafeParserRULE_variableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(44)
		p.Match(unsafeParserCOLON_EQUAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(45)
		p.Expression()
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
	p.RuleIndex = unsafeParserRULE_ifStatement
	return p
}

func InitEmptyIfStatementContext(p *IfStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_ifStatement
}

func (*IfStatementContext) IsIfStatementContext() {}

func NewIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfStatementContext {
	var p = new(IfStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_ifStatement

	return p
}

func (s *IfStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *IfStatementContext) AllDOUBLE_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserDOUBLE_LBRACE)
}

func (s *IfStatementContext) DOUBLE_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_LBRACE, i)
}

func (s *IfStatementContext) IF() antlr.TerminalNode {
	return s.GetToken(unsafeParserIF, 0)
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
	return s.GetTokens(unsafeParserDOUBLE_RBRACE)
}

func (s *IfStatementContext) DOUBLE_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_RBRACE, i)
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
	return s.GetToken(unsafeParserEND, 0)
}

func (s *IfStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IfStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitIfStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) IfStatement() (localctx IIfStatementContext) {
	localctx = NewIfStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, unsafeParserRULE_ifStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(47)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(48)
		p.Match(unsafeParserIF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(49)
		p.Expression()
	}
	{
		p.SetState(50)
		p.Match(unsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(51)
		p.Template()
	}
	{
		p.SetState(52)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(53)
		p.Match(unsafeParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(54)
		p.Match(unsafeParserDOUBLE_RBRACE)
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
	p.RuleIndex = unsafeParserRULE_loopStatement
	return p
}

func InitEmptyLoopStatementContext(p *LoopStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_loopStatement
}

func (*LoopStatementContext) IsLoopStatementContext() {}

func NewLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoopStatementContext {
	var p = new(LoopStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_loopStatement

	return p
}

func (s *LoopStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *LoopStatementContext) AllDOUBLE_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserDOUBLE_LBRACE)
}

func (s *LoopStatementContext) DOUBLE_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_LBRACE, i)
}

func (s *LoopStatementContext) LOOP() antlr.TerminalNode {
	return s.GetToken(unsafeParserLOOP, 0)
}

func (s *LoopStatementContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserIDENTIFIER)
}

func (s *LoopStatementContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, i)
}

func (s *LoopStatementContext) IN() antlr.TerminalNode {
	return s.GetToken(unsafeParserIN, 0)
}

func (s *LoopStatementContext) AllDOUBLE_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserDOUBLE_RBRACE)
}

func (s *LoopStatementContext) DOUBLE_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_RBRACE, i)
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
	return s.GetToken(unsafeParserEND, 0)
}

func (s *LoopStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LoopStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitLoopStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) LoopStatement() (localctx ILoopStatementContext) {
	localctx = NewLoopStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, unsafeParserRULE_loopStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(56)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(57)
		p.Match(unsafeParserLOOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(58)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(59)
		p.Match(unsafeParserIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(60)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(61)
		p.Match(unsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(62)
		p.Template()
	}
	{
		p.SetState(63)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(64)
		p.Match(unsafeParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(65)
		p.Match(unsafeParserDOUBLE_RBRACE)
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
	p.RuleIndex = unsafeParserRULE_expressionStatement
	return p
}

func InitEmptyExpressionStatementContext(p *ExpressionStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_expressionStatement
}

func (*ExpressionStatementContext) IsExpressionStatementContext() {}

func NewExpressionStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionStatementContext {
	var p = new(ExpressionStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_expressionStatement

	return p
}

func (s *ExpressionStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionStatementContext) DOUBLE_LBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_LBRACE, 0)
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
	return s.GetToken(unsafeParserDOUBLE_RBRACE, 0)
}

func (s *ExpressionStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitExpressionStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) ExpressionStatement() (localctx IExpressionStatementContext) {
	localctx = NewExpressionStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, unsafeParserRULE_expressionStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(67)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(68)
		p.Expression()
	}
	{
		p.SetState(69)
		p.Match(unsafeParserDOUBLE_RBRACE)
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

// IDefineStatementContext is an interface to support dynamic dispatch.
type IDefineStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDOUBLE_LBRACE() []antlr.TerminalNode
	DOUBLE_LBRACE(i int) antlr.TerminalNode
	DEFINE() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOUBLE_RBRACE() []antlr.TerminalNode
	DOUBLE_RBRACE(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsDefineStatementContext differentiates from other interfaces.
	IsDefineStatementContext()
}

type DefineStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDefineStatementContext() *DefineStatementContext {
	var p = new(DefineStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_defineStatement
	return p
}

func InitEmptyDefineStatementContext(p *DefineStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_defineStatement
}

func (*DefineStatementContext) IsDefineStatementContext() {}

func NewDefineStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefineStatementContext {
	var p = new(DefineStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_defineStatement

	return p
}

func (s *DefineStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DefineStatementContext) AllDOUBLE_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserDOUBLE_LBRACE)
}

func (s *DefineStatementContext) DOUBLE_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_LBRACE, i)
}

func (s *DefineStatementContext) DEFINE() antlr.TerminalNode {
	return s.GetToken(unsafeParserDEFINE, 0)
}

func (s *DefineStatementContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserIDENTIFIER)
}

func (s *DefineStatementContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, i)
}

func (s *DefineStatementContext) AllDOUBLE_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserDOUBLE_RBRACE)
}

func (s *DefineStatementContext) DOUBLE_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_RBRACE, i)
}

func (s *DefineStatementContext) Template() ITemplateContext {
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

func (s *DefineStatementContext) END() antlr.TerminalNode {
	return s.GetToken(unsafeParserEND, 0)
}

func (s *DefineStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN, 0)
}

func (s *DefineStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN, 0)
}

func (s *DefineStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserCOMMA)
}

func (s *DefineStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserCOMMA, i)
}

func (s *DefineStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DefineStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DefineStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitDefineStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) DefineStatement() (localctx IDefineStatementContext) {
	localctx = NewDefineStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, unsafeParserRULE_defineStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(71)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(72)
		p.Match(unsafeParserDEFINE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(73)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(86)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == unsafeParserLPAREN {
		{
			p.SetState(74)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(83)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == unsafeParserIDENTIFIER {
			{
				p.SetState(75)
				p.Match(unsafeParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(80)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == unsafeParserCOMMA {
				{
					p.SetState(76)
					p.Match(unsafeParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(77)
					p.Match(unsafeParserIDENTIFIER)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(82)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(85)
			p.Match(unsafeParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(88)
		p.Match(unsafeParserDOUBLE_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(89)
		p.Template()
	}
	{
		p.SetState(90)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(91)
		p.Match(unsafeParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(92)
		p.Match(unsafeParserDOUBLE_RBRACE)
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

// ITemplateStatementContext is an interface to support dynamic dispatch.
type ITemplateStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOUBLE_LBRACE() antlr.TerminalNode
	TEMPLATE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	DOUBLE_RBRACE() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsTemplateStatementContext differentiates from other interfaces.
	IsTemplateStatementContext()
}

type TemplateStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTemplateStatementContext() *TemplateStatementContext {
	var p = new(TemplateStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_templateStatement
	return p
}

func InitEmptyTemplateStatementContext(p *TemplateStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_templateStatement
}

func (*TemplateStatementContext) IsTemplateStatementContext() {}

func NewTemplateStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TemplateStatementContext {
	var p = new(TemplateStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_templateStatement

	return p
}

func (s *TemplateStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *TemplateStatementContext) DOUBLE_LBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_LBRACE, 0)
}

func (s *TemplateStatementContext) TEMPLATE() antlr.TerminalNode {
	return s.GetToken(unsafeParserTEMPLATE, 0)
}

func (s *TemplateStatementContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *TemplateStatementContext) DOUBLE_RBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserDOUBLE_RBRACE, 0)
}

func (s *TemplateStatementContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *TemplateStatementContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
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

	return t.(IExpressionContext)
}

func (s *TemplateStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TemplateStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TemplateStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitTemplateStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) TemplateStatement() (localctx ITemplateStatementContext) {
	localctx = NewTemplateStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, unsafeParserRULE_templateStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(95)
		p.Match(unsafeParserTEMPLATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(96)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&459776) != 0 {
		{
			p.SetState(97)
			p.Expression()
		}

		p.SetState(102)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(103)
		p.Match(unsafeParserDOUBLE_RBRACE)
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
	p.RuleIndex = unsafeParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_expression

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
	case unsafeParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, unsafeParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(105)
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
	p.RuleIndex = unsafeParserRULE_pipelineExpr
	return p
}

func InitEmptyPipelineExprContext(p *PipelineExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_pipelineExpr
}

func (*PipelineExprContext) IsPipelineExprContext() {}

func NewPipelineExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PipelineExprContext {
	var p = new(PipelineExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_pipelineExpr

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
	return s.GetTokens(unsafeParserPIPE)
}

func (s *PipelineExprContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserPIPE, i)
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
	case unsafeParserVisitor:
		return t.VisitPipelineExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) PipelineExpr() (localctx IPipelineExprContext) {
	localctx = NewPipelineExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, unsafeParserRULE_pipelineExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)
		p.PrimaryExpr()
	}
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserPIPE {
		{
			p.SetState(108)
			p.Match(unsafeParserPIPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(109)
			p.FunctionCall()
		}

		p.SetState(114)
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
	p.RuleIndex = unsafeParserRULE_primaryExpr
	return p
}

func InitEmptyPrimaryExprContext(p *PrimaryExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_primaryExpr
}

func (*PrimaryExprContext) IsPrimaryExprContext() {}

func NewPrimaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExprContext {
	var p = new(PrimaryExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_primaryExpr

	return p
}

func (s *PrimaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExprContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(unsafeParserNUMBER, 0)
}

func (s *PrimaryExprContext) STRING() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTRING, 0)
}

func (s *PrimaryExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *PrimaryExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN, 0)
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
	return s.GetToken(unsafeParserRPAREN, 0)
}

func (s *PrimaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitPrimaryExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) PrimaryExpr() (localctx IPrimaryExprContext) {
	localctx = NewPrimaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, unsafeParserRULE_primaryExpr)
	p.SetState(122)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case unsafeParserNUMBER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(115)
			p.Match(unsafeParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(116)
			p.Match(unsafeParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(117)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserLPAREN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(118)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(119)
			p.Expression()
		}
		{
			p.SetState(120)
			p.Match(unsafeParserRPAREN)
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
	p.RuleIndex = unsafeParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
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
	return s.GetToken(unsafeParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN, 0)
}

func (s *FunctionCallContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserCOMMA)
}

func (s *FunctionCallContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserCOMMA, i)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, unsafeParserRULE_functionCall)
	var _la int

	var _alt int

	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(124)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(128)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(125)
					p.PrimaryExpr()
				}

			}
			p.SetState(130)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(131)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(132)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(134)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&459776) != 0 {
			{
				p.SetState(133)
				p.PrimaryExpr()
			}

		}
		p.SetState(140)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == unsafeParserCOMMA {
			{
				p.SetState(136)
				p.Match(unsafeParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(137)
				p.PrimaryExpr()
			}

			p.SetState(142)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(143)
			p.Match(unsafeParserRPAREN)
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
