// Code generated from unsafeLexer.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated

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
		"DEFAULT_MODE", "COMMENT_MODE", "EXPRESSION_MODE", "STATEMENT_MODE",
	}
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
		"COMMENT_LBRACE", "DOUBLE_LBRACE", "STATEMENT_LBRACE", "TEXT", "COMMENT_RBRACE",
		"COMMENT_CONTENT", "DOUBLE_RBRACE", "LPAREN", "RPAREN", "LBRACE", "RBRACE",
		"HASH", "MOD", "DOLLAR", "PIPE", "COMMA", "COLON_EQUAL", "DOT", "SINGLE_QUOTE",
		"QUOTE", "STRING", "NUMBER", "IDENTIFIER", "ESC", "WS", "STATEMENT_RBRACE",
		"IF", "ELSE", "LOOP", "IN", "MACRO", "TEMPLATE", "END", "ENDIF", "ENDLOOP",
		"ENDMACRO", "IDENTIFIER_STMT", "WS_STMT",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 38, 284, 6, -1, 6, -1, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2,
		7, 2, 2, 3, 7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8,
		7, 8, 2, 9, 7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13,
		2, 14, 7, 14, 2, 15, 7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2,
		19, 7, 19, 2, 20, 7, 20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24,
		7, 24, 2, 25, 7, 25, 2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7,
		29, 2, 30, 7, 30, 2, 31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34,
		2, 35, 7, 35, 2, 36, 7, 36, 2, 37, 7, 37, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 4, 3,
		97, 8, 3, 11, 3, 12, 3, 98, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 4, 5, 107,
		8, 5, 11, 5, 12, 5, 108, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7,
		1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12, 1, 12,
		1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1,
		17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 5, 20, 148, 8, 20,
		10, 20, 12, 20, 151, 9, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 5, 20, 158,
		8, 20, 10, 20, 12, 20, 161, 9, 20, 1, 20, 1, 20, 3, 20, 165, 8, 20, 1,
		21, 4, 21, 168, 8, 21, 11, 21, 12, 21, 169, 1, 21, 1, 21, 4, 21, 174, 8,
		21, 11, 21, 12, 21, 175, 3, 21, 178, 8, 21, 1, 21, 1, 21, 3, 21, 182, 8,
		21, 1, 21, 4, 21, 185, 8, 21, 11, 21, 12, 21, 186, 3, 21, 189, 8, 21, 1,
		22, 1, 22, 5, 22, 193, 8, 22, 10, 22, 12, 22, 196, 9, 22, 1, 23, 1, 23,
		1, 23, 1, 24, 4, 24, 202, 8, 24, 11, 24, 12, 24, 203, 1, 24, 1, 24, 1,
		25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 27,
		1, 27, 1, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 1,
		30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31,
		1, 31, 1, 31, 1, 31, 1, 31, 1, 32, 1, 32, 1, 32, 1, 32, 1, 33, 1, 33, 1,
		33, 1, 33, 1, 33, 1, 33, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34,
		1, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1,
		36, 1, 36, 5, 36, 273, 8, 36, 10, 36, 12, 36, 276, 9, 36, 1, 37, 4, 37,
		279, 8, 37, 11, 37, 12, 37, 280, 1, 37, 1, 37, 1, 108, 0, 38, 4, 1, 6,
		2, 8, 3, 10, 4, 12, 5, 14, 6, 16, 7, 18, 8, 20, 9, 22, 10, 24, 11, 26,
		12, 28, 13, 30, 14, 32, 15, 34, 16, 36, 17, 38, 18, 40, 19, 42, 20, 44,
		21, 46, 22, 48, 23, 50, 24, 52, 25, 54, 26, 56, 27, 58, 28, 60, 29, 62,
		30, 64, 31, 66, 32, 68, 33, 70, 34, 72, 35, 74, 36, 76, 37, 78, 38, 4,
		0, 1, 2, 3, 10, 1, 0, 123, 123, 3, 0, 10, 10, 39, 39, 92, 92, 3, 0, 10,
		10, 34, 34, 92, 92, 1, 0, 48, 57, 2, 0, 69, 69, 101, 101, 2, 0, 43, 43,
		45, 45, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97,
		122, 8, 0, 34, 34, 39, 39, 92, 92, 98, 98, 102, 102, 110, 110, 114, 114,
		116, 116, 3, 0, 9, 10, 13, 13, 32, 32, 297, 0, 4, 1, 0, 0, 0, 0, 6, 1,
		0, 0, 0, 0, 8, 1, 0, 0, 0, 0, 10, 1, 0, 0, 0, 1, 12, 1, 0, 0, 0, 1, 14,
		1, 0, 0, 0, 2, 16, 1, 0, 0, 0, 2, 18, 1, 0, 0, 0, 2, 20, 1, 0, 0, 0, 2,
		22, 1, 0, 0, 0, 2, 24, 1, 0, 0, 0, 2, 26, 1, 0, 0, 0, 2, 28, 1, 0, 0, 0,
		2, 30, 1, 0, 0, 0, 2, 32, 1, 0, 0, 0, 2, 34, 1, 0, 0, 0, 2, 36, 1, 0, 0,
		0, 2, 38, 1, 0, 0, 0, 2, 40, 1, 0, 0, 0, 2, 42, 1, 0, 0, 0, 2, 44, 1, 0,
		0, 0, 2, 46, 1, 0, 0, 0, 2, 48, 1, 0, 0, 0, 2, 50, 1, 0, 0, 0, 2, 52, 1,
		0, 0, 0, 3, 54, 1, 0, 0, 0, 3, 56, 1, 0, 0, 0, 3, 58, 1, 0, 0, 0, 3, 60,
		1, 0, 0, 0, 3, 62, 1, 0, 0, 0, 3, 64, 1, 0, 0, 0, 3, 66, 1, 0, 0, 0, 3,
		68, 1, 0, 0, 0, 3, 70, 1, 0, 0, 0, 3, 72, 1, 0, 0, 0, 3, 74, 1, 0, 0, 0,
		3, 76, 1, 0, 0, 0, 3, 78, 1, 0, 0, 0, 4, 80, 1, 0, 0, 0, 6, 85, 1, 0, 0,
		0, 8, 90, 1, 0, 0, 0, 10, 96, 1, 0, 0, 0, 12, 100, 1, 0, 0, 0, 14, 106,
		1, 0, 0, 0, 16, 112, 1, 0, 0, 0, 18, 117, 1, 0, 0, 0, 20, 119, 1, 0, 0,
		0, 22, 121, 1, 0, 0, 0, 24, 123, 1, 0, 0, 0, 26, 125, 1, 0, 0, 0, 28, 127,
		1, 0, 0, 0, 30, 129, 1, 0, 0, 0, 32, 131, 1, 0, 0, 0, 34, 133, 1, 0, 0,
		0, 36, 135, 1, 0, 0, 0, 38, 138, 1, 0, 0, 0, 40, 140, 1, 0, 0, 0, 42, 142,
		1, 0, 0, 0, 44, 164, 1, 0, 0, 0, 46, 167, 1, 0, 0, 0, 48, 190, 1, 0, 0,
		0, 50, 197, 1, 0, 0, 0, 52, 201, 1, 0, 0, 0, 54, 207, 1, 0, 0, 0, 56, 212,
		1, 0, 0, 0, 58, 215, 1, 0, 0, 0, 60, 220, 1, 0, 0, 0, 62, 225, 1, 0, 0,
		0, 64, 228, 1, 0, 0, 0, 66, 234, 1, 0, 0, 0, 68, 243, 1, 0, 0, 0, 70, 247,
		1, 0, 0, 0, 72, 253, 1, 0, 0, 0, 74, 261, 1, 0, 0, 0, 76, 270, 1, 0, 0,
		0, 78, 278, 1, 0, 0, 0, 80, 81, 5, 123, 0, 0, 81, 82, 5, 35, 0, 0, 82,
		83, 1, 0, 0, 0, 83, 84, 6, 0, 0, 0, 84, 5, 1, 0, 0, 0, 85, 86, 5, 123,
		0, 0, 86, 87, 5, 123, 0, 0, 87, 88, 1, 0, 0, 0, 88, 89, 6, 1, 1, 0, 89,
		7, 1, 0, 0, 0, 90, 91, 5, 123, 0, 0, 91, 92, 5, 37, 0, 0, 92, 93, 1, 0,
		0, 0, 93, 94, 6, 2, 2, 0, 94, 9, 1, 0, 0, 0, 95, 97, 8, 0, 0, 0, 96, 95,
		1, 0, 0, 0, 97, 98, 1, 0, 0, 0, 98, 96, 1, 0, 0, 0, 98, 99, 1, 0, 0, 0,
		99, 11, 1, 0, 0, 0, 100, 101, 5, 35, 0, 0, 101, 102, 5, 125, 0, 0, 102,
		103, 1, 0, 0, 0, 103, 104, 6, 4, 3, 0, 104, 13, 1, 0, 0, 0, 105, 107, 9,
		0, 0, 0, 106, 105, 1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 109, 1, 0, 0,
		0, 108, 106, 1, 0, 0, 0, 109, 110, 1, 0, 0, 0, 110, 111, 6, 5, 4, 0, 111,
		15, 1, 0, 0, 0, 112, 113, 5, 125, 0, 0, 113, 114, 5, 125, 0, 0, 114, 115,
		1, 0, 0, 0, 115, 116, 6, 6, 3, 0, 116, 17, 1, 0, 0, 0, 117, 118, 5, 40,
		0, 0, 118, 19, 1, 0, 0, 0, 119, 120, 5, 41, 0, 0, 120, 21, 1, 0, 0, 0,
		121, 122, 5, 123, 0, 0, 122, 23, 1, 0, 0, 0, 123, 124, 5, 125, 0, 0, 124,
		25, 1, 0, 0, 0, 125, 126, 5, 35, 0, 0, 126, 27, 1, 0, 0, 0, 127, 128, 5,
		37, 0, 0, 128, 29, 1, 0, 0, 0, 129, 130, 5, 36, 0, 0, 130, 31, 1, 0, 0,
		0, 131, 132, 5, 124, 0, 0, 132, 33, 1, 0, 0, 0, 133, 134, 5, 44, 0, 0,
		134, 35, 1, 0, 0, 0, 135, 136, 5, 58, 0, 0, 136, 137, 5, 61, 0, 0, 137,
		37, 1, 0, 0, 0, 138, 139, 5, 46, 0, 0, 139, 39, 1, 0, 0, 0, 140, 141, 5,
		39, 0, 0, 141, 41, 1, 0, 0, 0, 142, 143, 5, 34, 0, 0, 143, 43, 1, 0, 0,
		0, 144, 149, 3, 40, 18, 0, 145, 148, 3, 50, 23, 0, 146, 148, 8, 1, 0, 0,
		147, 145, 1, 0, 0, 0, 147, 146, 1, 0, 0, 0, 148, 151, 1, 0, 0, 0, 149,
		147, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 152, 1, 0, 0, 0, 151, 149,
		1, 0, 0, 0, 152, 153, 3, 40, 18, 0, 153, 165, 1, 0, 0, 0, 154, 159, 3,
		42, 19, 0, 155, 158, 3, 50, 23, 0, 156, 158, 8, 2, 0, 0, 157, 155, 1, 0,
		0, 0, 157, 156, 1, 0, 0, 0, 158, 161, 1, 0, 0, 0, 159, 157, 1, 0, 0, 0,
		159, 160, 1, 0, 0, 0, 160, 162, 1, 0, 0, 0, 161, 159, 1, 0, 0, 0, 162,
		163, 3, 42, 19, 0, 163, 165, 1, 0, 0, 0, 164, 144, 1, 0, 0, 0, 164, 154,
		1, 0, 0, 0, 165, 45, 1, 0, 0, 0, 166, 168, 7, 3, 0, 0, 167, 166, 1, 0,
		0, 0, 168, 169, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0,
		170, 177, 1, 0, 0, 0, 171, 173, 5, 46, 0, 0, 172, 174, 7, 3, 0, 0, 173,
		172, 1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 173, 1, 0, 0, 0, 175, 176,
		1, 0, 0, 0, 176, 178, 1, 0, 0, 0, 177, 171, 1, 0, 0, 0, 177, 178, 1, 0,
		0, 0, 178, 188, 1, 0, 0, 0, 179, 181, 7, 4, 0, 0, 180, 182, 7, 5, 0, 0,
		181, 180, 1, 0, 0, 0, 181, 182, 1, 0, 0, 0, 182, 184, 1, 0, 0, 0, 183,
		185, 7, 3, 0, 0, 184, 183, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 184,
		1, 0, 0, 0, 186, 187, 1, 0, 0, 0, 187, 189, 1, 0, 0, 0, 188, 179, 1, 0,
		0, 0, 188, 189, 1, 0, 0, 0, 189, 47, 1, 0, 0, 0, 190, 194, 7, 6, 0, 0,
		191, 193, 7, 7, 0, 0, 192, 191, 1, 0, 0, 0, 193, 196, 1, 0, 0, 0, 194,
		192, 1, 0, 0, 0, 194, 195, 1, 0, 0, 0, 195, 49, 1, 0, 0, 0, 196, 194, 1,
		0, 0, 0, 197, 198, 5, 92, 0, 0, 198, 199, 7, 8, 0, 0, 199, 51, 1, 0, 0,
		0, 200, 202, 7, 9, 0, 0, 201, 200, 1, 0, 0, 0, 202, 203, 1, 0, 0, 0, 203,
		201, 1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 205, 1, 0, 0, 0, 205, 206,
		6, 24, 4, 0, 206, 53, 1, 0, 0, 0, 207, 208, 5, 37, 0, 0, 208, 209, 5, 125,
		0, 0, 209, 210, 1, 0, 0, 0, 210, 211, 6, 25, 3, 0, 211, 55, 1, 0, 0, 0,
		212, 213, 5, 105, 0, 0, 213, 214, 5, 102, 0, 0, 214, 57, 1, 0, 0, 0, 215,
		216, 5, 101, 0, 0, 216, 217, 5, 108, 0, 0, 217, 218, 5, 115, 0, 0, 218,
		219, 5, 101, 0, 0, 219, 59, 1, 0, 0, 0, 220, 221, 5, 108, 0, 0, 221, 222,
		5, 111, 0, 0, 222, 223, 5, 111, 0, 0, 223, 224, 5, 112, 0, 0, 224, 61,
		1, 0, 0, 0, 225, 226, 5, 105, 0, 0, 226, 227, 5, 110, 0, 0, 227, 63, 1,
		0, 0, 0, 228, 229, 5, 109, 0, 0, 229, 230, 5, 97, 0, 0, 230, 231, 5, 99,
		0, 0, 231, 232, 5, 114, 0, 0, 232, 233, 5, 111, 0, 0, 233, 65, 1, 0, 0,
		0, 234, 235, 5, 116, 0, 0, 235, 236, 5, 101, 0, 0, 236, 237, 5, 109, 0,
		0, 237, 238, 5, 112, 0, 0, 238, 239, 5, 108, 0, 0, 239, 240, 5, 97, 0,
		0, 240, 241, 5, 116, 0, 0, 241, 242, 5, 101, 0, 0, 242, 67, 1, 0, 0, 0,
		243, 244, 5, 101, 0, 0, 244, 245, 5, 110, 0, 0, 245, 246, 5, 100, 0, 0,
		246, 69, 1, 0, 0, 0, 247, 248, 5, 101, 0, 0, 248, 249, 5, 110, 0, 0, 249,
		250, 5, 100, 0, 0, 250, 251, 5, 105, 0, 0, 251, 252, 5, 102, 0, 0, 252,
		71, 1, 0, 0, 0, 253, 254, 5, 101, 0, 0, 254, 255, 5, 110, 0, 0, 255, 256,
		5, 100, 0, 0, 256, 257, 5, 108, 0, 0, 257, 258, 5, 111, 0, 0, 258, 259,
		5, 111, 0, 0, 259, 260, 5, 112, 0, 0, 260, 73, 1, 0, 0, 0, 261, 262, 5,
		101, 0, 0, 262, 263, 5, 110, 0, 0, 263, 264, 5, 100, 0, 0, 264, 265, 5,
		109, 0, 0, 265, 266, 5, 97, 0, 0, 266, 267, 5, 99, 0, 0, 267, 268, 5, 114,
		0, 0, 268, 269, 5, 111, 0, 0, 269, 75, 1, 0, 0, 0, 270, 274, 7, 6, 0, 0,
		271, 273, 7, 7, 0, 0, 272, 271, 1, 0, 0, 0, 273, 276, 1, 0, 0, 0, 274,
		272, 1, 0, 0, 0, 274, 275, 1, 0, 0, 0, 275, 77, 1, 0, 0, 0, 276, 274, 1,
		0, 0, 0, 277, 279, 7, 9, 0, 0, 278, 277, 1, 0, 0, 0, 279, 280, 1, 0, 0,
		0, 280, 278, 1, 0, 0, 0, 280, 281, 1, 0, 0, 0, 281, 282, 1, 0, 0, 0, 282,
		283, 6, 37, 4, 0, 283, 79, 1, 0, 0, 0, 21, 0, 1, 2, 3, 98, 108, 147, 149,
		157, 159, 164, 169, 175, 177, 181, 186, 188, 194, 203, 274, 280, 5, 5,
		1, 0, 5, 2, 0, 5, 3, 0, 4, 0, 0, 6, 0, 0,
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
	unsafeLexerCOMMENT_LBRACE   = 1
	unsafeLexerDOUBLE_LBRACE    = 2
	unsafeLexerSTATEMENT_LBRACE = 3
	unsafeLexerTEXT             = 4
	unsafeLexerCOMMENT_RBRACE   = 5
	unsafeLexerCOMMENT_CONTENT  = 6
	unsafeLexerDOUBLE_RBRACE    = 7
	unsafeLexerLPAREN           = 8
	unsafeLexerRPAREN           = 9
	unsafeLexerLBRACE           = 10
	unsafeLexerRBRACE           = 11
	unsafeLexerHASH             = 12
	unsafeLexerMOD              = 13
	unsafeLexerDOLLAR           = 14
	unsafeLexerPIPE             = 15
	unsafeLexerCOMMA            = 16
	unsafeLexerCOLON_EQUAL      = 17
	unsafeLexerDOT              = 18
	unsafeLexerSINGLE_QUOTE     = 19
	unsafeLexerQUOTE            = 20
	unsafeLexerSTRING           = 21
	unsafeLexerNUMBER           = 22
	unsafeLexerIDENTIFIER       = 23
	unsafeLexerESC              = 24
	unsafeLexerWS               = 25
	unsafeLexerSTATEMENT_RBRACE = 26
	unsafeLexerIF               = 27
	unsafeLexerELSE             = 28
	unsafeLexerLOOP             = 29
	unsafeLexerIN               = 30
	unsafeLexerMACRO            = 31
	unsafeLexerTEMPLATE         = 32
	unsafeLexerEND              = 33
	unsafeLexerENDIF            = 34
	unsafeLexerENDLOOP          = 35
	unsafeLexerENDMACRO         = 36
	unsafeLexerIDENTIFIER_STMT  = 37
	unsafeLexerWS_STMT          = 38
)

// unsafeLexer modes.
const (
	unsafeLexerCOMMENT_MODE = iota + 1
	unsafeLexerEXPRESSION_MODE
	unsafeLexerSTATEMENT_MODE
)
