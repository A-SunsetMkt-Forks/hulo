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
		"", "'{#'", "'{{'", "'{%'", "", "'#}'", "", "'}}'", "'('", "')'", "'{'",
		"'}'", "'#'", "'%'", "'$'", "'|'", "','", "':='", "'.'", "'''", "'\"'",
		"", "", "", "", "", "'%}'", "'if'", "'else'", "'loop'", "'in'", "'macro'",
		"'template'", "'end'", "'endif'", "'endloop'", "'endmacro'",
	}
	staticData.SymbolicNames = []string{
		"", "COMMENT_LBRACE", "DOUBLE_LBRACE", "STATEMENT_LBRACE", "TEXT", "COMMENT_RBRACE",
		"COMMENT_CONTENT", "DOUBLE_RBRACE", "LPAREN", "RPAREN", "LBRACE", "RBRACE",
		"HASH", "MOD", "DOLLAR", "PIPE", "COMMA", "COLON_EQUAL", "DOT", "SINGLE_QUOTE",
		"QUOTE", "STRING", "NUMBER", "IDENTIFIER", "ESC", "WS", "STATEMENT_RBRACE",
		"IF", "ELSE", "LOOP", "IN", "MACRO", "TEMPLATE", "END", "ENDIF", "ENDLOOP",
		"ENDMACRO", "IDENTIFIER_STMT", "WS_STMT",
	}
	staticData.RuleNames = []string{
		"template", "content", "statement", "variableStatement", "ifStatement",
		"elseStatement", "loopStatement", "expressionStatement", "macroStatement",
		"templateStatement", "expression", "pipelineExpr", "primaryExpr", "varExpr",
		"functionCall",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 38, 165, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 1, 0, 1, 0,
		5, 0, 33, 8, 0, 10, 0, 12, 0, 36, 9, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 3, 2, 46, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4,
		1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 60, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1,
		5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1,
		8, 1, 8, 1, 8, 5, 8, 93, 8, 8, 10, 8, 12, 8, 96, 9, 8, 3, 8, 98, 8, 8,
		1, 8, 3, 8, 101, 8, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9,
		1, 9, 1, 9, 5, 9, 113, 8, 9, 10, 9, 12, 9, 116, 9, 9, 1, 9, 1, 9, 1, 10,
		1, 10, 1, 11, 1, 11, 1, 11, 5, 11, 125, 8, 11, 10, 11, 12, 11, 128, 9,
		11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 3, 12, 138,
		8, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 5, 14, 145, 8, 14, 10, 14, 12,
		14, 148, 9, 14, 1, 14, 1, 14, 1, 14, 3, 14, 153, 8, 14, 1, 14, 1, 14, 5,
		14, 157, 8, 14, 10, 14, 12, 14, 160, 9, 14, 1, 14, 3, 14, 163, 8, 14, 1,
		14, 0, 0, 15, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 0,
		3, 1, 0, 33, 34, 2, 0, 33, 33, 35, 35, 2, 0, 33, 33, 36, 36, 170, 0, 34,
		1, 0, 0, 0, 2, 37, 1, 0, 0, 0, 4, 45, 1, 0, 0, 0, 6, 47, 1, 0, 0, 0, 8,
		53, 1, 0, 0, 0, 10, 65, 1, 0, 0, 0, 12, 70, 1, 0, 0, 0, 14, 81, 1, 0, 0,
		0, 16, 85, 1, 0, 0, 0, 18, 108, 1, 0, 0, 0, 20, 119, 1, 0, 0, 0, 22, 121,
		1, 0, 0, 0, 24, 137, 1, 0, 0, 0, 26, 139, 1, 0, 0, 0, 28, 162, 1, 0, 0,
		0, 30, 33, 3, 2, 1, 0, 31, 33, 3, 4, 2, 0, 32, 30, 1, 0, 0, 0, 32, 31,
		1, 0, 0, 0, 33, 36, 1, 0, 0, 0, 34, 32, 1, 0, 0, 0, 34, 35, 1, 0, 0, 0,
		35, 1, 1, 0, 0, 0, 36, 34, 1, 0, 0, 0, 37, 38, 5, 4, 0, 0, 38, 3, 1, 0,
		0, 0, 39, 46, 3, 8, 4, 0, 40, 46, 3, 12, 6, 0, 41, 46, 3, 14, 7, 0, 42,
		46, 3, 16, 8, 0, 43, 46, 3, 18, 9, 0, 44, 46, 3, 6, 3, 0, 45, 39, 1, 0,
		0, 0, 45, 40, 1, 0, 0, 0, 45, 41, 1, 0, 0, 0, 45, 42, 1, 0, 0, 0, 45, 43,
		1, 0, 0, 0, 45, 44, 1, 0, 0, 0, 46, 5, 1, 0, 0, 0, 47, 48, 5, 3, 0, 0,
		48, 49, 5, 23, 0, 0, 49, 50, 5, 17, 0, 0, 50, 51, 3, 20, 10, 0, 51, 52,
		5, 26, 0, 0, 52, 7, 1, 0, 0, 0, 53, 54, 5, 3, 0, 0, 54, 55, 5, 27, 0, 0,
		55, 56, 3, 20, 10, 0, 56, 57, 5, 26, 0, 0, 57, 59, 3, 0, 0, 0, 58, 60,
		3, 10, 5, 0, 59, 58, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0, 60, 61, 1, 0, 0, 0,
		61, 62, 5, 3, 0, 0, 62, 63, 7, 0, 0, 0, 63, 64, 5, 26, 0, 0, 64, 9, 1,
		0, 0, 0, 65, 66, 5, 3, 0, 0, 66, 67, 5, 28, 0, 0, 67, 68, 5, 26, 0, 0,
		68, 69, 3, 0, 0, 0, 69, 11, 1, 0, 0, 0, 70, 71, 5, 3, 0, 0, 71, 72, 5,
		29, 0, 0, 72, 73, 5, 23, 0, 0, 73, 74, 5, 30, 0, 0, 74, 75, 5, 23, 0, 0,
		75, 76, 5, 26, 0, 0, 76, 77, 3, 0, 0, 0, 77, 78, 5, 3, 0, 0, 78, 79, 7,
		1, 0, 0, 79, 80, 5, 26, 0, 0, 80, 13, 1, 0, 0, 0, 81, 82, 5, 2, 0, 0, 82,
		83, 3, 20, 10, 0, 83, 84, 5, 7, 0, 0, 84, 15, 1, 0, 0, 0, 85, 86, 5, 3,
		0, 0, 86, 87, 5, 31, 0, 0, 87, 100, 5, 23, 0, 0, 88, 97, 5, 8, 0, 0, 89,
		94, 5, 23, 0, 0, 90, 91, 5, 16, 0, 0, 91, 93, 5, 23, 0, 0, 92, 90, 1, 0,
		0, 0, 93, 96, 1, 0, 0, 0, 94, 92, 1, 0, 0, 0, 94, 95, 1, 0, 0, 0, 95, 98,
		1, 0, 0, 0, 96, 94, 1, 0, 0, 0, 97, 89, 1, 0, 0, 0, 97, 98, 1, 0, 0, 0,
		98, 99, 1, 0, 0, 0, 99, 101, 5, 9, 0, 0, 100, 88, 1, 0, 0, 0, 100, 101,
		1, 0, 0, 0, 101, 102, 1, 0, 0, 0, 102, 103, 5, 26, 0, 0, 103, 104, 3, 0,
		0, 0, 104, 105, 5, 3, 0, 0, 105, 106, 7, 2, 0, 0, 106, 107, 5, 26, 0, 0,
		107, 17, 1, 0, 0, 0, 108, 109, 5, 3, 0, 0, 109, 110, 5, 32, 0, 0, 110,
		114, 5, 23, 0, 0, 111, 113, 3, 20, 10, 0, 112, 111, 1, 0, 0, 0, 113, 116,
		1, 0, 0, 0, 114, 112, 1, 0, 0, 0, 114, 115, 1, 0, 0, 0, 115, 117, 1, 0,
		0, 0, 116, 114, 1, 0, 0, 0, 117, 118, 5, 26, 0, 0, 118, 19, 1, 0, 0, 0,
		119, 120, 3, 22, 11, 0, 120, 21, 1, 0, 0, 0, 121, 126, 3, 28, 14, 0, 122,
		123, 5, 15, 0, 0, 123, 125, 3, 28, 14, 0, 124, 122, 1, 0, 0, 0, 125, 128,
		1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 23, 1, 0,
		0, 0, 128, 126, 1, 0, 0, 0, 129, 138, 5, 22, 0, 0, 130, 138, 5, 21, 0,
		0, 131, 138, 5, 23, 0, 0, 132, 138, 3, 26, 13, 0, 133, 134, 5, 8, 0, 0,
		134, 135, 3, 20, 10, 0, 135, 136, 5, 9, 0, 0, 136, 138, 1, 0, 0, 0, 137,
		129, 1, 0, 0, 0, 137, 130, 1, 0, 0, 0, 137, 131, 1, 0, 0, 0, 137, 132,
		1, 0, 0, 0, 137, 133, 1, 0, 0, 0, 138, 25, 1, 0, 0, 0, 139, 140, 5, 14,
		0, 0, 140, 141, 5, 23, 0, 0, 141, 27, 1, 0, 0, 0, 142, 146, 5, 23, 0, 0,
		143, 145, 3, 24, 12, 0, 144, 143, 1, 0, 0, 0, 145, 148, 1, 0, 0, 0, 146,
		144, 1, 0, 0, 0, 146, 147, 1, 0, 0, 0, 147, 163, 1, 0, 0, 0, 148, 146,
		1, 0, 0, 0, 149, 150, 5, 23, 0, 0, 150, 152, 5, 8, 0, 0, 151, 153, 3, 24,
		12, 0, 152, 151, 1, 0, 0, 0, 152, 153, 1, 0, 0, 0, 153, 158, 1, 0, 0, 0,
		154, 155, 5, 16, 0, 0, 155, 157, 3, 24, 12, 0, 156, 154, 1, 0, 0, 0, 157,
		160, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 158, 159, 1, 0, 0, 0, 159, 161,
		1, 0, 0, 0, 160, 158, 1, 0, 0, 0, 161, 163, 5, 9, 0, 0, 162, 142, 1, 0,
		0, 0, 162, 149, 1, 0, 0, 0, 163, 29, 1, 0, 0, 0, 14, 32, 34, 45, 59, 94,
		97, 100, 114, 126, 137, 146, 152, 158, 162,
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
	unsafeParserEOF              = antlr.TokenEOF
	unsafeParserCOMMENT_LBRACE   = 1
	unsafeParserDOUBLE_LBRACE    = 2
	unsafeParserSTATEMENT_LBRACE = 3
	unsafeParserTEXT             = 4
	unsafeParserCOMMENT_RBRACE   = 5
	unsafeParserCOMMENT_CONTENT  = 6
	unsafeParserDOUBLE_RBRACE    = 7
	unsafeParserLPAREN           = 8
	unsafeParserRPAREN           = 9
	unsafeParserLBRACE           = 10
	unsafeParserRBRACE           = 11
	unsafeParserHASH             = 12
	unsafeParserMOD              = 13
	unsafeParserDOLLAR           = 14
	unsafeParserPIPE             = 15
	unsafeParserCOMMA            = 16
	unsafeParserCOLON_EQUAL      = 17
	unsafeParserDOT              = 18
	unsafeParserSINGLE_QUOTE     = 19
	unsafeParserQUOTE            = 20
	unsafeParserSTRING           = 21
	unsafeParserNUMBER           = 22
	unsafeParserIDENTIFIER       = 23
	unsafeParserESC              = 24
	unsafeParserWS               = 25
	unsafeParserSTATEMENT_RBRACE = 26
	unsafeParserIF               = 27
	unsafeParserELSE             = 28
	unsafeParserLOOP             = 29
	unsafeParserIN               = 30
	unsafeParserMACRO            = 31
	unsafeParserTEMPLATE         = 32
	unsafeParserEND              = 33
	unsafeParserENDIF            = 34
	unsafeParserENDLOOP          = 35
	unsafeParserENDMACRO         = 36
	unsafeParserIDENTIFIER_STMT  = 37
	unsafeParserWS_STMT          = 38
)

// unsafeParser rules.
const (
	unsafeParserRULE_template            = 0
	unsafeParserRULE_content             = 1
	unsafeParserRULE_statement           = 2
	unsafeParserRULE_variableStatement   = 3
	unsafeParserRULE_ifStatement         = 4
	unsafeParserRULE_elseStatement       = 5
	unsafeParserRULE_loopStatement       = 6
	unsafeParserRULE_expressionStatement = 7
	unsafeParserRULE_macroStatement      = 8
	unsafeParserRULE_templateStatement   = 9
	unsafeParserRULE_expression          = 10
	unsafeParserRULE_pipelineExpr        = 11
	unsafeParserRULE_primaryExpr         = 12
	unsafeParserRULE_varExpr             = 13
	unsafeParserRULE_functionCall        = 14
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
	p.SetState(34)
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
			p.SetState(32)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case unsafeParserTEXT:
				{
					p.SetState(30)
					p.Content()
				}

			case unsafeParserDOUBLE_LBRACE, unsafeParserSTATEMENT_LBRACE:
				{
					p.SetState(31)
					p.Statement()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

		}
		p.SetState(36)
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
		p.SetState(37)
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
	MacroStatement() IMacroStatementContext
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

func (s *StatementContext) MacroStatement() IMacroStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMacroStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMacroStatementContext)
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
	p.SetState(45)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(39)
			p.IfStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(40)
			p.LoopStatement()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(41)
			p.ExpressionStatement()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(42)
			p.MacroStatement()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(43)
			p.TemplateStatement()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(44)
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
	STATEMENT_LBRACE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	COLON_EQUAL() antlr.TerminalNode
	Expression() IExpressionContext
	STATEMENT_RBRACE() antlr.TerminalNode

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

func (s *VariableStatementContext) STATEMENT_LBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_LBRACE, 0)
}

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

func (s *VariableStatementContext) STATEMENT_RBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_RBRACE, 0)
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
		p.SetState(47)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(48)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(49)
		p.Match(unsafeParserCOLON_EQUAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(50)
		p.Expression()
	}
	{
		p.SetState(51)
		p.Match(unsafeParserSTATEMENT_RBRACE)
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

// IIfStatementContext is an interface to support dynamic dispatch.
type IIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSTATEMENT_LBRACE() []antlr.TerminalNode
	STATEMENT_LBRACE(i int) antlr.TerminalNode
	IF() antlr.TerminalNode
	Expression() IExpressionContext
	AllSTATEMENT_RBRACE() []antlr.TerminalNode
	STATEMENT_RBRACE(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode
	ENDIF() antlr.TerminalNode
	ElseStatement() IElseStatementContext

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

func (s *IfStatementContext) AllSTATEMENT_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTATEMENT_LBRACE)
}

func (s *IfStatementContext) STATEMENT_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_LBRACE, i)
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

func (s *IfStatementContext) AllSTATEMENT_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTATEMENT_RBRACE)
}

func (s *IfStatementContext) STATEMENT_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_RBRACE, i)
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

func (s *IfStatementContext) ENDIF() antlr.TerminalNode {
	return s.GetToken(unsafeParserENDIF, 0)
}

func (s *IfStatementContext) ElseStatement() IElseStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseStatementContext)
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(53)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(54)
		p.Match(unsafeParserIF)
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
		p.Match(unsafeParserSTATEMENT_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(57)
		p.Template()
	}
	p.SetState(59)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(58)
			p.ElseStatement()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(61)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(62)
		_la = p.GetTokenStream().LA(1)

		if !(_la == unsafeParserEND || _la == unsafeParserENDIF) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(63)
		p.Match(unsafeParserSTATEMENT_RBRACE)
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

// IElseStatementContext is an interface to support dynamic dispatch.
type IElseStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STATEMENT_LBRACE() antlr.TerminalNode
	ELSE() antlr.TerminalNode
	STATEMENT_RBRACE() antlr.TerminalNode
	Template() ITemplateContext

	// IsElseStatementContext differentiates from other interfaces.
	IsElseStatementContext()
}

type ElseStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElseStatementContext() *ElseStatementContext {
	var p = new(ElseStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_elseStatement
	return p
}

func InitEmptyElseStatementContext(p *ElseStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_elseStatement
}

func (*ElseStatementContext) IsElseStatementContext() {}

func NewElseStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseStatementContext {
	var p = new(ElseStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_elseStatement

	return p
}

func (s *ElseStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ElseStatementContext) STATEMENT_LBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_LBRACE, 0)
}

func (s *ElseStatementContext) ELSE() antlr.TerminalNode {
	return s.GetToken(unsafeParserELSE, 0)
}

func (s *ElseStatementContext) STATEMENT_RBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_RBRACE, 0)
}

func (s *ElseStatementContext) Template() ITemplateContext {
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

func (s *ElseStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElseStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ElseStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitElseStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) ElseStatement() (localctx IElseStatementContext) {
	localctx = NewElseStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, unsafeParserRULE_elseStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(65)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(66)
		p.Match(unsafeParserELSE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(67)
		p.Match(unsafeParserSTATEMENT_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(68)
		p.Template()
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
	AllSTATEMENT_LBRACE() []antlr.TerminalNode
	STATEMENT_LBRACE(i int) antlr.TerminalNode
	LOOP() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	IN() antlr.TerminalNode
	AllSTATEMENT_RBRACE() []antlr.TerminalNode
	STATEMENT_RBRACE(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode
	ENDLOOP() antlr.TerminalNode

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

func (s *LoopStatementContext) AllSTATEMENT_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTATEMENT_LBRACE)
}

func (s *LoopStatementContext) STATEMENT_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_LBRACE, i)
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

func (s *LoopStatementContext) AllSTATEMENT_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTATEMENT_RBRACE)
}

func (s *LoopStatementContext) STATEMENT_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_RBRACE, i)
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

func (s *LoopStatementContext) ENDLOOP() antlr.TerminalNode {
	return s.GetToken(unsafeParserENDLOOP, 0)
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
	p.EnterRule(localctx, 12, unsafeParserRULE_loopStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(71)
		p.Match(unsafeParserLOOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(72)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(73)
		p.Match(unsafeParserIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(74)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(75)
		p.Match(unsafeParserSTATEMENT_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Template()
	}
	{
		p.SetState(77)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(78)
		_la = p.GetTokenStream().LA(1)

		if !(_la == unsafeParserEND || _la == unsafeParserENDLOOP) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(79)
		p.Match(unsafeParserSTATEMENT_RBRACE)
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
	p.EnterRule(localctx, 14, unsafeParserRULE_expressionStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(81)
		p.Match(unsafeParserDOUBLE_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(82)
		p.Expression()
	}
	{
		p.SetState(83)
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

// IMacroStatementContext is an interface to support dynamic dispatch.
type IMacroStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSTATEMENT_LBRACE() []antlr.TerminalNode
	STATEMENT_LBRACE(i int) antlr.TerminalNode
	MACRO() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllSTATEMENT_RBRACE() []antlr.TerminalNode
	STATEMENT_RBRACE(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode
	ENDMACRO() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsMacroStatementContext differentiates from other interfaces.
	IsMacroStatementContext()
}

type MacroStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMacroStatementContext() *MacroStatementContext {
	var p = new(MacroStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_macroStatement
	return p
}

func InitEmptyMacroStatementContext(p *MacroStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_macroStatement
}

func (*MacroStatementContext) IsMacroStatementContext() {}

func NewMacroStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroStatementContext {
	var p = new(MacroStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_macroStatement

	return p
}

func (s *MacroStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *MacroStatementContext) AllSTATEMENT_LBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTATEMENT_LBRACE)
}

func (s *MacroStatementContext) STATEMENT_LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_LBRACE, i)
}

func (s *MacroStatementContext) MACRO() antlr.TerminalNode {
	return s.GetToken(unsafeParserMACRO, 0)
}

func (s *MacroStatementContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserIDENTIFIER)
}

func (s *MacroStatementContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, i)
}

func (s *MacroStatementContext) AllSTATEMENT_RBRACE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTATEMENT_RBRACE)
}

func (s *MacroStatementContext) STATEMENT_RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_RBRACE, i)
}

func (s *MacroStatementContext) Template() ITemplateContext {
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

func (s *MacroStatementContext) END() antlr.TerminalNode {
	return s.GetToken(unsafeParserEND, 0)
}

func (s *MacroStatementContext) ENDMACRO() antlr.TerminalNode {
	return s.GetToken(unsafeParserENDMACRO, 0)
}

func (s *MacroStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN, 0)
}

func (s *MacroStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN, 0)
}

func (s *MacroStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserCOMMA)
}

func (s *MacroStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserCOMMA, i)
}

func (s *MacroStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MacroStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MacroStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitMacroStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) MacroStatement() (localctx IMacroStatementContext) {
	localctx = NewMacroStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, unsafeParserRULE_macroStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(86)
		p.Match(unsafeParserMACRO)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(87)
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

	if _la == unsafeParserLPAREN {
		{
			p.SetState(88)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(97)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == unsafeParserIDENTIFIER {
			{
				p.SetState(89)
				p.Match(unsafeParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(94)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == unsafeParserCOMMA {
				{
					p.SetState(90)
					p.Match(unsafeParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(91)
					p.Match(unsafeParserIDENTIFIER)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(96)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(99)
			p.Match(unsafeParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(102)
		p.Match(unsafeParserSTATEMENT_RBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(103)
		p.Template()
	}
	{
		p.SetState(104)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(105)
		_la = p.GetTokenStream().LA(1)

		if !(_la == unsafeParserEND || _la == unsafeParserENDMACRO) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(106)
		p.Match(unsafeParserSTATEMENT_RBRACE)
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
	STATEMENT_LBRACE() antlr.TerminalNode
	TEMPLATE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	STATEMENT_RBRACE() antlr.TerminalNode
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

func (s *TemplateStatementContext) STATEMENT_LBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_LBRACE, 0)
}

func (s *TemplateStatementContext) TEMPLATE() antlr.TerminalNode {
	return s.GetToken(unsafeParserTEMPLATE, 0)
}

func (s *TemplateStatementContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *TemplateStatementContext) STATEMENT_RBRACE() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTATEMENT_RBRACE, 0)
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
	p.EnterRule(localctx, 18, unsafeParserRULE_templateStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(108)
		p.Match(unsafeParserSTATEMENT_LBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(109)
		p.Match(unsafeParserTEMPLATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(110)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(114)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserIDENTIFIER {
		{
			p.SetState(111)
			p.Expression()
		}

		p.SetState(116)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(117)
		p.Match(unsafeParserSTATEMENT_RBRACE)
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
	p.EnterRule(localctx, 20, unsafeParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
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
	AllFunctionCall() []IFunctionCallContext
	FunctionCall(i int) IFunctionCallContext
	AllPIPE() []antlr.TerminalNode
	PIPE(i int) antlr.TerminalNode

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

func (s *PipelineExprContext) AllPIPE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserPIPE)
}

func (s *PipelineExprContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserPIPE, i)
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
	p.EnterRule(localctx, 22, unsafeParserRULE_pipelineExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(121)
		p.FunctionCall()
	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserPIPE {
		{
			p.SetState(122)
			p.Match(unsafeParserPIPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(123)
			p.FunctionCall()
		}

		p.SetState(128)
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
	VarExpr() IVarExprContext
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

func (s *PrimaryExprContext) VarExpr() IVarExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarExprContext)
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
	p.EnterRule(localctx, 24, unsafeParserRULE_primaryExpr)
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case unsafeParserNUMBER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(129)
			p.Match(unsafeParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(130)
			p.Match(unsafeParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(131)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserDOLLAR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(132)
			p.VarExpr()
		}

	case unsafeParserLPAREN:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(133)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(134)
			p.Expression()
		}
		{
			p.SetState(135)
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

// IVarExprContext is an interface to support dynamic dispatch.
type IVarExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOLLAR() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsVarExprContext differentiates from other interfaces.
	IsVarExprContext()
}

type VarExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVarExprContext() *VarExprContext {
	var p = new(VarExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_varExpr
	return p
}

func InitEmptyVarExprContext(p *VarExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_varExpr
}

func (*VarExprContext) IsVarExprContext() {}

func NewVarExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarExprContext {
	var p = new(VarExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_varExpr

	return p
}

func (s *VarExprContext) GetParser() antlr.Parser { return s.parser }

func (s *VarExprContext) DOLLAR() antlr.TerminalNode {
	return s.GetToken(unsafeParserDOLLAR, 0)
}

func (s *VarExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *VarExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VarExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitVarExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) VarExpr() (localctx IVarExprContext) {
	localctx = NewVarExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, unsafeParserRULE_varExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(139)
		p.Match(unsafeParserDOLLAR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(140)
		p.Match(unsafeParserIDENTIFIER)
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
	p.EnterRule(localctx, 28, unsafeParserRULE_functionCall)
	var _la int

	var _alt int

	p.SetState(162)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(142)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(146)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(143)
					p.PrimaryExpr()
				}

			}
			p.SetState(148)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(149)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(150)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(152)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&14696704) != 0 {
			{
				p.SetState(151)
				p.PrimaryExpr()
			}

		}
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == unsafeParserCOMMA {
			{
				p.SetState(154)
				p.Match(unsafeParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(155)
				p.PrimaryExpr()
			}

			p.SetState(160)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(161)
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
