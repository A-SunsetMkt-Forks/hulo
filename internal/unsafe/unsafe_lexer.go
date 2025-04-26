// Code generated from unsafeLexer.g4 by ANTLR 4.13.2. DO NOT EDIT.

package unsafe

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type unsafeLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var UnsafeLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func unsafelexerLexerInit() {
	staticData := &UnsafeLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE", "TEMPLATE_MODE",
	}
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
		"DOUBLE_LBRACE", "TEXT", "DOUBLE_RBRACE", "IF", "LOOP", "END", "IN",
		"LPAREN", "RPAREN", "DOLLAR", "PIPE", "COMMA", "IDENTIFIER", "STRING",
		"NUMBER", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 16, 123, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3,
		7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9,
		7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7,
		14, 2, 15, 7, 15, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 4, 1, 41, 8, 1, 11,
		1, 12, 1, 42, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1,
		7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12, 1, 12, 5,
		12, 77, 8, 12, 10, 12, 12, 12, 80, 9, 12, 1, 13, 1, 13, 1, 13, 1, 13, 5,
		13, 86, 8, 13, 10, 13, 12, 13, 89, 9, 13, 1, 13, 1, 13, 1, 14, 4, 14, 94,
		8, 14, 11, 14, 12, 14, 95, 1, 14, 1, 14, 4, 14, 100, 8, 14, 11, 14, 12,
		14, 101, 3, 14, 104, 8, 14, 1, 14, 1, 14, 3, 14, 108, 8, 14, 1, 14, 4,
		14, 111, 8, 14, 11, 14, 12, 14, 112, 3, 14, 115, 8, 14, 1, 15, 4, 15, 118,
		8, 15, 11, 15, 12, 15, 119, 1, 15, 1, 15, 0, 0, 16, 2, 1, 4, 2, 6, 3, 8,
		4, 10, 5, 12, 6, 14, 7, 16, 8, 18, 9, 20, 10, 22, 11, 24, 12, 26, 13, 28,
		14, 30, 15, 32, 16, 2, 0, 1, 8, 1, 0, 123, 123, 3, 0, 65, 90, 95, 95, 97,
		122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 2, 0, 34, 34, 92, 92, 1, 0,
		48, 57, 2, 0, 69, 69, 101, 101, 2, 0, 43, 43, 45, 45, 3, 0, 9, 10, 13,
		13, 32, 32, 132, 0, 2, 1, 0, 0, 0, 0, 4, 1, 0, 0, 0, 1, 6, 1, 0, 0, 0,
		1, 8, 1, 0, 0, 0, 1, 10, 1, 0, 0, 0, 1, 12, 1, 0, 0, 0, 1, 14, 1, 0, 0,
		0, 1, 16, 1, 0, 0, 0, 1, 18, 1, 0, 0, 0, 1, 20, 1, 0, 0, 0, 1, 22, 1, 0,
		0, 0, 1, 24, 1, 0, 0, 0, 1, 26, 1, 0, 0, 0, 1, 28, 1, 0, 0, 0, 1, 30, 1,
		0, 0, 0, 1, 32, 1, 0, 0, 0, 2, 34, 1, 0, 0, 0, 4, 40, 1, 0, 0, 0, 6, 44,
		1, 0, 0, 0, 8, 49, 1, 0, 0, 0, 10, 52, 1, 0, 0, 0, 12, 57, 1, 0, 0, 0,
		14, 61, 1, 0, 0, 0, 16, 64, 1, 0, 0, 0, 18, 66, 1, 0, 0, 0, 20, 68, 1,
		0, 0, 0, 22, 70, 1, 0, 0, 0, 24, 72, 1, 0, 0, 0, 26, 74, 1, 0, 0, 0, 28,
		81, 1, 0, 0, 0, 30, 93, 1, 0, 0, 0, 32, 117, 1, 0, 0, 0, 34, 35, 5, 123,
		0, 0, 35, 36, 5, 123, 0, 0, 36, 37, 1, 0, 0, 0, 37, 38, 6, 0, 0, 0, 38,
		3, 1, 0, 0, 0, 39, 41, 8, 0, 0, 0, 40, 39, 1, 0, 0, 0, 41, 42, 1, 0, 0,
		0, 42, 40, 1, 0, 0, 0, 42, 43, 1, 0, 0, 0, 43, 5, 1, 0, 0, 0, 44, 45, 5,
		125, 0, 0, 45, 46, 5, 125, 0, 0, 46, 47, 1, 0, 0, 0, 47, 48, 6, 2, 1, 0,
		48, 7, 1, 0, 0, 0, 49, 50, 5, 105, 0, 0, 50, 51, 5, 102, 0, 0, 51, 9, 1,
		0, 0, 0, 52, 53, 5, 108, 0, 0, 53, 54, 5, 111, 0, 0, 54, 55, 5, 111, 0,
		0, 55, 56, 5, 112, 0, 0, 56, 11, 1, 0, 0, 0, 57, 58, 5, 101, 0, 0, 58,
		59, 5, 110, 0, 0, 59, 60, 5, 100, 0, 0, 60, 13, 1, 0, 0, 0, 61, 62, 5,
		105, 0, 0, 62, 63, 5, 110, 0, 0, 63, 15, 1, 0, 0, 0, 64, 65, 5, 40, 0,
		0, 65, 17, 1, 0, 0, 0, 66, 67, 5, 41, 0, 0, 67, 19, 1, 0, 0, 0, 68, 69,
		5, 36, 0, 0, 69, 21, 1, 0, 0, 0, 70, 71, 5, 124, 0, 0, 71, 23, 1, 0, 0,
		0, 72, 73, 5, 44, 0, 0, 73, 25, 1, 0, 0, 0, 74, 78, 7, 1, 0, 0, 75, 77,
		7, 2, 0, 0, 76, 75, 1, 0, 0, 0, 77, 80, 1, 0, 0, 0, 78, 76, 1, 0, 0, 0,
		78, 79, 1, 0, 0, 0, 79, 27, 1, 0, 0, 0, 80, 78, 1, 0, 0, 0, 81, 87, 5,
		34, 0, 0, 82, 86, 8, 3, 0, 0, 83, 84, 5, 92, 0, 0, 84, 86, 9, 0, 0, 0,
		85, 82, 1, 0, 0, 0, 85, 83, 1, 0, 0, 0, 86, 89, 1, 0, 0, 0, 87, 85, 1,
		0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 90, 1, 0, 0, 0, 89, 87, 1, 0, 0, 0, 90,
		91, 5, 34, 0, 0, 91, 29, 1, 0, 0, 0, 92, 94, 7, 4, 0, 0, 93, 92, 1, 0,
		0, 0, 94, 95, 1, 0, 0, 0, 95, 93, 1, 0, 0, 0, 95, 96, 1, 0, 0, 0, 96, 103,
		1, 0, 0, 0, 97, 99, 5, 46, 0, 0, 98, 100, 7, 4, 0, 0, 99, 98, 1, 0, 0,
		0, 100, 101, 1, 0, 0, 0, 101, 99, 1, 0, 0, 0, 101, 102, 1, 0, 0, 0, 102,
		104, 1, 0, 0, 0, 103, 97, 1, 0, 0, 0, 103, 104, 1, 0, 0, 0, 104, 114, 1,
		0, 0, 0, 105, 107, 7, 5, 0, 0, 106, 108, 7, 6, 0, 0, 107, 106, 1, 0, 0,
		0, 107, 108, 1, 0, 0, 0, 108, 110, 1, 0, 0, 0, 109, 111, 7, 4, 0, 0, 110,
		109, 1, 0, 0, 0, 111, 112, 1, 0, 0, 0, 112, 110, 1, 0, 0, 0, 112, 113,
		1, 0, 0, 0, 113, 115, 1, 0, 0, 0, 114, 105, 1, 0, 0, 0, 114, 115, 1, 0,
		0, 0, 115, 31, 1, 0, 0, 0, 116, 118, 7, 7, 0, 0, 117, 116, 1, 0, 0, 0,
		118, 119, 1, 0, 0, 0, 119, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120,
		121, 1, 0, 0, 0, 121, 122, 6, 15, 2, 0, 122, 33, 1, 0, 0, 0, 13, 0, 1,
		42, 78, 85, 87, 95, 101, 103, 107, 112, 114, 119, 3, 5, 1, 0, 4, 0, 0,
		6, 0, 0,
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

// unsafeLexerInit initializes any static state used to implement unsafeLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewunsafeLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func UnsafeLexerInit() {
	staticData := &UnsafeLexerLexerStaticData
	staticData.once.Do(unsafelexerLexerInit)
}

// NewunsafeLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewunsafeLexer(input antlr.CharStream) *unsafeLexer {
	UnsafeLexerInit()
	l := new(unsafeLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &UnsafeLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "unsafeLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// unsafeLexer tokens.
const (
	unsafeLexerDOUBLE_LBRACE = 1
	unsafeLexerTEXT          = 2
	unsafeLexerDOUBLE_RBRACE = 3
	unsafeLexerIF            = 4
	unsafeLexerLOOP          = 5
	unsafeLexerEND           = 6
	unsafeLexerIN            = 7
	unsafeLexerLPAREN        = 8
	unsafeLexerRPAREN        = 9
	unsafeLexerDOLLAR        = 10
	unsafeLexerPIPE          = 11
	unsafeLexerCOMMA         = 12
	unsafeLexerIDENTIFIER    = 13
	unsafeLexerSTRING        = 14
	unsafeLexerNUMBER        = 15
	unsafeLexerWS            = 16
)

// unsafeLexerTEMPLATE_MODE is the unsafeLexer mode.
const unsafeLexerTEMPLATE_MODE = 1
