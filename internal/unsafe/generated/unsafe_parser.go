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
		"", "'{#'", "'{{'", "'{%'", "", "'#}'", "", "'}}'", "", "", "'{'", "'}'",
		"'#'", "'%'", "'$'", "'|'", "','", "':='", "'.'", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "'%}'", "", "", "'if'", "'else'",
		"'loop'", "'in'", "'macro'", "'end'", "'endif'", "'endloop'", "'endmacro'",
	}
	staticData.SymbolicNames = []string{
		"", "COMMENT_START", "EXPR_START", "STMT_START", "TEXT", "COMMENT_END",
		"COMMENT_CONTENT", "EXPR_END", "LPAREN", "RPAREN", "LBRACE", "RBRACE",
		"HASH", "MOD", "DOLLAR", "PIPE", "COMMA", "COLON_EQUAL", "DOT", "EQ",
		"NE", "GT", "LT", "GE", "LE", "AND", "OR", "NOT", "STRING", "NUMBER",
		"BOOLEAN", "IDENTIFIER", "WS", "STMT_END", "LPAREN_STMT", "RPAREN_STMT",
		"IF", "ELSE", "LOOP", "IN", "MACRO", "END", "ENDIF", "ENDLOOP", "ENDMACRO",
		"EQ_STMT", "NE_STMT", "GT_STMT", "LT_STMT", "GE_STMT", "LE_STMT", "AND_STMT",
		"OR_STMT", "NOT_STMT", "STRING_STMT", "NUMBER_STMT", "BOOLEAN_STMT",
		"IDENTIFIER_STMT", "WS_STMT",
	}
	staticData.RuleNames = []string{
		"template", "content", "statement", "commentStatement", "variableStatement",
		"ifStatement", "elseStatement", "loopStatement", "expressionStatement",
		"macroStatement", "pipelineExpression", "expression", "logicalOrExpression",
		"logicalAndExpression", "equalityExpression", "comparisonExpression",
		"primaryExpression", "varExpr", "functionCall", "logicalOrExpressionStmt",
		"logicalAndExpressionStmt", "equalityExpressionStmt", "comparisonExpressionStmt",
		"primaryExpressionStmt",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 58, 276, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 1, 0, 1, 0, 5, 0, 51, 8, 0, 10,
		0, 12, 0, 54, 9, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3,
		2, 64, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 82, 8, 5, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1,
		7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 1, 9, 5, 9, 115, 8, 9, 10, 9, 12, 9, 118, 9, 9, 3, 9, 120,
		8, 9, 1, 9, 3, 9, 123, 8, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10,
		1, 10, 1, 10, 5, 10, 134, 8, 10, 10, 10, 12, 10, 137, 9, 10, 1, 11, 1,
		11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 3, 11, 154, 8, 11, 1, 12, 1, 12, 1, 12, 5, 12, 159,
		8, 12, 10, 12, 12, 12, 162, 9, 12, 1, 13, 1, 13, 1, 13, 5, 13, 167, 8,
		13, 10, 13, 12, 13, 170, 9, 13, 1, 14, 1, 14, 1, 14, 5, 14, 175, 8, 14,
		10, 14, 12, 14, 178, 9, 14, 1, 15, 1, 15, 1, 15, 5, 15, 183, 8, 15, 10,
		15, 12, 15, 186, 9, 15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16,
		1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3,
		16, 205, 8, 16, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 5, 18, 212, 8, 18, 10,
		18, 12, 18, 215, 9, 18, 1, 18, 1, 18, 1, 18, 3, 18, 220, 8, 18, 1, 18,
		1, 18, 5, 18, 224, 8, 18, 10, 18, 12, 18, 227, 9, 18, 1, 18, 3, 18, 230,
		8, 18, 1, 19, 1, 19, 1, 19, 5, 19, 235, 8, 19, 10, 19, 12, 19, 238, 9,
		19, 1, 20, 1, 20, 1, 20, 5, 20, 243, 8, 20, 10, 20, 12, 20, 246, 9, 20,
		1, 21, 1, 21, 1, 21, 5, 21, 251, 8, 21, 10, 21, 12, 21, 254, 9, 21, 1,
		22, 1, 22, 1, 22, 5, 22, 259, 8, 22, 10, 22, 12, 22, 262, 9, 22, 1, 23,
		1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 274,
		8, 23, 1, 23, 0, 0, 24, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24,
		26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 0, 7, 1, 0, 41, 42, 2, 0, 41,
		41, 43, 43, 2, 0, 41, 41, 44, 44, 1, 0, 19, 20, 1, 0, 21, 24, 1, 0, 45,
		46, 1, 0, 47, 50, 302, 0, 52, 1, 0, 0, 0, 2, 55, 1, 0, 0, 0, 4, 63, 1,
		0, 0, 0, 6, 65, 1, 0, 0, 0, 8, 69, 1, 0, 0, 0, 10, 75, 1, 0, 0, 0, 12,
		87, 1, 0, 0, 0, 14, 92, 1, 0, 0, 0, 16, 103, 1, 0, 0, 0, 18, 107, 1, 0,
		0, 0, 20, 130, 1, 0, 0, 0, 22, 153, 1, 0, 0, 0, 24, 155, 1, 0, 0, 0, 26,
		163, 1, 0, 0, 0, 28, 171, 1, 0, 0, 0, 30, 179, 1, 0, 0, 0, 32, 204, 1,
		0, 0, 0, 34, 206, 1, 0, 0, 0, 36, 229, 1, 0, 0, 0, 38, 231, 1, 0, 0, 0,
		40, 239, 1, 0, 0, 0, 42, 247, 1, 0, 0, 0, 44, 255, 1, 0, 0, 0, 46, 273,
		1, 0, 0, 0, 48, 51, 3, 2, 1, 0, 49, 51, 3, 4, 2, 0, 50, 48, 1, 0, 0, 0,
		50, 49, 1, 0, 0, 0, 51, 54, 1, 0, 0, 0, 52, 50, 1, 0, 0, 0, 52, 53, 1,
		0, 0, 0, 53, 1, 1, 0, 0, 0, 54, 52, 1, 0, 0, 0, 55, 56, 5, 4, 0, 0, 56,
		3, 1, 0, 0, 0, 57, 64, 3, 6, 3, 0, 58, 64, 3, 10, 5, 0, 59, 64, 3, 14,
		7, 0, 60, 64, 3, 16, 8, 0, 61, 64, 3, 18, 9, 0, 62, 64, 3, 8, 4, 0, 63,
		57, 1, 0, 0, 0, 63, 58, 1, 0, 0, 0, 63, 59, 1, 0, 0, 0, 63, 60, 1, 0, 0,
		0, 63, 61, 1, 0, 0, 0, 63, 62, 1, 0, 0, 0, 64, 5, 1, 0, 0, 0, 65, 66, 5,
		1, 0, 0, 66, 67, 5, 6, 0, 0, 67, 68, 5, 5, 0, 0, 68, 7, 1, 0, 0, 0, 69,
		70, 5, 3, 0, 0, 70, 71, 5, 31, 0, 0, 71, 72, 5, 17, 0, 0, 72, 73, 3, 22,
		11, 0, 73, 74, 5, 33, 0, 0, 74, 9, 1, 0, 0, 0, 75, 76, 5, 3, 0, 0, 76,
		77, 5, 36, 0, 0, 77, 78, 3, 22, 11, 0, 78, 79, 5, 33, 0, 0, 79, 81, 3,
		0, 0, 0, 80, 82, 3, 12, 6, 0, 81, 80, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82,
		83, 1, 0, 0, 0, 83, 84, 5, 3, 0, 0, 84, 85, 7, 0, 0, 0, 85, 86, 5, 33,
		0, 0, 86, 11, 1, 0, 0, 0, 87, 88, 5, 3, 0, 0, 88, 89, 5, 37, 0, 0, 89,
		90, 5, 33, 0, 0, 90, 91, 3, 0, 0, 0, 91, 13, 1, 0, 0, 0, 92, 93, 5, 3,
		0, 0, 93, 94, 5, 38, 0, 0, 94, 95, 5, 57, 0, 0, 95, 96, 5, 39, 0, 0, 96,
		97, 5, 57, 0, 0, 97, 98, 5, 33, 0, 0, 98, 99, 3, 0, 0, 0, 99, 100, 5, 3,
		0, 0, 100, 101, 7, 1, 0, 0, 101, 102, 5, 33, 0, 0, 102, 15, 1, 0, 0, 0,
		103, 104, 5, 2, 0, 0, 104, 105, 3, 20, 10, 0, 105, 106, 5, 7, 0, 0, 106,
		17, 1, 0, 0, 0, 107, 108, 5, 3, 0, 0, 108, 109, 5, 40, 0, 0, 109, 122,
		5, 57, 0, 0, 110, 119, 5, 34, 0, 0, 111, 116, 5, 57, 0, 0, 112, 113, 5,
		16, 0, 0, 113, 115, 5, 57, 0, 0, 114, 112, 1, 0, 0, 0, 115, 118, 1, 0,
		0, 0, 116, 114, 1, 0, 0, 0, 116, 117, 1, 0, 0, 0, 117, 120, 1, 0, 0, 0,
		118, 116, 1, 0, 0, 0, 119, 111, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120,
		121, 1, 0, 0, 0, 121, 123, 5, 35, 0, 0, 122, 110, 1, 0, 0, 0, 122, 123,
		1, 0, 0, 0, 123, 124, 1, 0, 0, 0, 124, 125, 5, 33, 0, 0, 125, 126, 3, 0,
		0, 0, 126, 127, 5, 3, 0, 0, 127, 128, 7, 2, 0, 0, 128, 129, 5, 33, 0, 0,
		129, 19, 1, 0, 0, 0, 130, 135, 3, 22, 11, 0, 131, 132, 5, 15, 0, 0, 132,
		134, 3, 22, 11, 0, 133, 131, 1, 0, 0, 0, 134, 137, 1, 0, 0, 0, 135, 133,
		1, 0, 0, 0, 135, 136, 1, 0, 0, 0, 136, 21, 1, 0, 0, 0, 137, 135, 1, 0,
		0, 0, 138, 154, 5, 29, 0, 0, 139, 154, 5, 28, 0, 0, 140, 154, 5, 30, 0,
		0, 141, 154, 5, 31, 0, 0, 142, 154, 5, 55, 0, 0, 143, 154, 5, 54, 0, 0,
		144, 154, 5, 56, 0, 0, 145, 154, 3, 34, 17, 0, 146, 147, 5, 8, 0, 0, 147,
		148, 3, 22, 11, 0, 148, 149, 5, 9, 0, 0, 149, 154, 1, 0, 0, 0, 150, 154,
		3, 36, 18, 0, 151, 154, 3, 24, 12, 0, 152, 154, 3, 38, 19, 0, 153, 138,
		1, 0, 0, 0, 153, 139, 1, 0, 0, 0, 153, 140, 1, 0, 0, 0, 153, 141, 1, 0,
		0, 0, 153, 142, 1, 0, 0, 0, 153, 143, 1, 0, 0, 0, 153, 144, 1, 0, 0, 0,
		153, 145, 1, 0, 0, 0, 153, 146, 1, 0, 0, 0, 153, 150, 1, 0, 0, 0, 153,
		151, 1, 0, 0, 0, 153, 152, 1, 0, 0, 0, 154, 23, 1, 0, 0, 0, 155, 160, 3,
		26, 13, 0, 156, 157, 5, 26, 0, 0, 157, 159, 3, 26, 13, 0, 158, 156, 1,
		0, 0, 0, 159, 162, 1, 0, 0, 0, 160, 158, 1, 0, 0, 0, 160, 161, 1, 0, 0,
		0, 161, 25, 1, 0, 0, 0, 162, 160, 1, 0, 0, 0, 163, 168, 3, 28, 14, 0, 164,
		165, 5, 25, 0, 0, 165, 167, 3, 28, 14, 0, 166, 164, 1, 0, 0, 0, 167, 170,
		1, 0, 0, 0, 168, 166, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 27, 1, 0,
		0, 0, 170, 168, 1, 0, 0, 0, 171, 176, 3, 30, 15, 0, 172, 173, 7, 3, 0,
		0, 173, 175, 3, 30, 15, 0, 174, 172, 1, 0, 0, 0, 175, 178, 1, 0, 0, 0,
		176, 174, 1, 0, 0, 0, 176, 177, 1, 0, 0, 0, 177, 29, 1, 0, 0, 0, 178, 176,
		1, 0, 0, 0, 179, 184, 3, 32, 16, 0, 180, 181, 7, 4, 0, 0, 181, 183, 3,
		32, 16, 0, 182, 180, 1, 0, 0, 0, 183, 186, 1, 0, 0, 0, 184, 182, 1, 0,
		0, 0, 184, 185, 1, 0, 0, 0, 185, 31, 1, 0, 0, 0, 186, 184, 1, 0, 0, 0,
		187, 205, 5, 29, 0, 0, 188, 205, 5, 28, 0, 0, 189, 205, 5, 30, 0, 0, 190,
		205, 5, 31, 0, 0, 191, 205, 5, 55, 0, 0, 192, 205, 5, 54, 0, 0, 193, 205,
		5, 56, 0, 0, 194, 205, 3, 34, 17, 0, 195, 196, 5, 8, 0, 0, 196, 197, 3,
		22, 11, 0, 197, 198, 5, 9, 0, 0, 198, 205, 1, 0, 0, 0, 199, 205, 3, 36,
		18, 0, 200, 201, 5, 27, 0, 0, 201, 205, 3, 32, 16, 0, 202, 203, 5, 53,
		0, 0, 203, 205, 3, 46, 23, 0, 204, 187, 1, 0, 0, 0, 204, 188, 1, 0, 0,
		0, 204, 189, 1, 0, 0, 0, 204, 190, 1, 0, 0, 0, 204, 191, 1, 0, 0, 0, 204,
		192, 1, 0, 0, 0, 204, 193, 1, 0, 0, 0, 204, 194, 1, 0, 0, 0, 204, 195,
		1, 0, 0, 0, 204, 199, 1, 0, 0, 0, 204, 200, 1, 0, 0, 0, 204, 202, 1, 0,
		0, 0, 205, 33, 1, 0, 0, 0, 206, 207, 5, 14, 0, 0, 207, 208, 5, 31, 0, 0,
		208, 35, 1, 0, 0, 0, 209, 213, 5, 31, 0, 0, 210, 212, 3, 22, 11, 0, 211,
		210, 1, 0, 0, 0, 212, 215, 1, 0, 0, 0, 213, 211, 1, 0, 0, 0, 213, 214,
		1, 0, 0, 0, 214, 230, 1, 0, 0, 0, 215, 213, 1, 0, 0, 0, 216, 217, 5, 31,
		0, 0, 217, 219, 5, 8, 0, 0, 218, 220, 3, 22, 11, 0, 219, 218, 1, 0, 0,
		0, 219, 220, 1, 0, 0, 0, 220, 225, 1, 0, 0, 0, 221, 222, 5, 16, 0, 0, 222,
		224, 3, 22, 11, 0, 223, 221, 1, 0, 0, 0, 224, 227, 1, 0, 0, 0, 225, 223,
		1, 0, 0, 0, 225, 226, 1, 0, 0, 0, 226, 228, 1, 0, 0, 0, 227, 225, 1, 0,
		0, 0, 228, 230, 5, 9, 0, 0, 229, 209, 1, 0, 0, 0, 229, 216, 1, 0, 0, 0,
		230, 37, 1, 0, 0, 0, 231, 236, 3, 40, 20, 0, 232, 233, 5, 52, 0, 0, 233,
		235, 3, 40, 20, 0, 234, 232, 1, 0, 0, 0, 235, 238, 1, 0, 0, 0, 236, 234,
		1, 0, 0, 0, 236, 237, 1, 0, 0, 0, 237, 39, 1, 0, 0, 0, 238, 236, 1, 0,
		0, 0, 239, 244, 3, 42, 21, 0, 240, 241, 5, 51, 0, 0, 241, 243, 3, 42, 21,
		0, 242, 240, 1, 0, 0, 0, 243, 246, 1, 0, 0, 0, 244, 242, 1, 0, 0, 0, 244,
		245, 1, 0, 0, 0, 245, 41, 1, 0, 0, 0, 246, 244, 1, 0, 0, 0, 247, 252, 3,
		44, 22, 0, 248, 249, 7, 5, 0, 0, 249, 251, 3, 44, 22, 0, 250, 248, 1, 0,
		0, 0, 251, 254, 1, 0, 0, 0, 252, 250, 1, 0, 0, 0, 252, 253, 1, 0, 0, 0,
		253, 43, 1, 0, 0, 0, 254, 252, 1, 0, 0, 0, 255, 260, 3, 46, 23, 0, 256,
		257, 7, 6, 0, 0, 257, 259, 3, 46, 23, 0, 258, 256, 1, 0, 0, 0, 259, 262,
		1, 0, 0, 0, 260, 258, 1, 0, 0, 0, 260, 261, 1, 0, 0, 0, 261, 45, 1, 0,
		0, 0, 262, 260, 1, 0, 0, 0, 263, 274, 5, 55, 0, 0, 264, 274, 5, 54, 0,
		0, 265, 274, 5, 56, 0, 0, 266, 274, 5, 57, 0, 0, 267, 268, 5, 34, 0, 0,
		268, 269, 3, 22, 11, 0, 269, 270, 5, 35, 0, 0, 270, 274, 1, 0, 0, 0, 271,
		272, 5, 53, 0, 0, 272, 274, 3, 46, 23, 0, 273, 263, 1, 0, 0, 0, 273, 264,
		1, 0, 0, 0, 273, 265, 1, 0, 0, 0, 273, 266, 1, 0, 0, 0, 273, 267, 1, 0,
		0, 0, 273, 271, 1, 0, 0, 0, 274, 47, 1, 0, 0, 0, 23, 50, 52, 63, 81, 116,
		119, 122, 135, 153, 160, 168, 176, 184, 204, 213, 219, 225, 229, 236, 244,
		252, 260, 273,
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
	unsafeParserEOF             = antlr.TokenEOF
	unsafeParserCOMMENT_START   = 1
	unsafeParserEXPR_START      = 2
	unsafeParserSTMT_START      = 3
	unsafeParserTEXT            = 4
	unsafeParserCOMMENT_END     = 5
	unsafeParserCOMMENT_CONTENT = 6
	unsafeParserEXPR_END        = 7
	unsafeParserLPAREN          = 8
	unsafeParserRPAREN          = 9
	unsafeParserLBRACE          = 10
	unsafeParserRBRACE          = 11
	unsafeParserHASH            = 12
	unsafeParserMOD             = 13
	unsafeParserDOLLAR          = 14
	unsafeParserPIPE            = 15
	unsafeParserCOMMA           = 16
	unsafeParserCOLON_EQUAL     = 17
	unsafeParserDOT             = 18
	unsafeParserEQ              = 19
	unsafeParserNE              = 20
	unsafeParserGT              = 21
	unsafeParserLT              = 22
	unsafeParserGE              = 23
	unsafeParserLE              = 24
	unsafeParserAND             = 25
	unsafeParserOR              = 26
	unsafeParserNOT             = 27
	unsafeParserSTRING          = 28
	unsafeParserNUMBER          = 29
	unsafeParserBOOLEAN         = 30
	unsafeParserIDENTIFIER      = 31
	unsafeParserWS              = 32
	unsafeParserSTMT_END        = 33
	unsafeParserLPAREN_STMT     = 34
	unsafeParserRPAREN_STMT     = 35
	unsafeParserIF              = 36
	unsafeParserELSE            = 37
	unsafeParserLOOP            = 38
	unsafeParserIN              = 39
	unsafeParserMACRO           = 40
	unsafeParserEND             = 41
	unsafeParserENDIF           = 42
	unsafeParserENDLOOP         = 43
	unsafeParserENDMACRO        = 44
	unsafeParserEQ_STMT         = 45
	unsafeParserNE_STMT         = 46
	unsafeParserGT_STMT         = 47
	unsafeParserLT_STMT         = 48
	unsafeParserGE_STMT         = 49
	unsafeParserLE_STMT         = 50
	unsafeParserAND_STMT        = 51
	unsafeParserOR_STMT         = 52
	unsafeParserNOT_STMT        = 53
	unsafeParserSTRING_STMT     = 54
	unsafeParserNUMBER_STMT     = 55
	unsafeParserBOOLEAN_STMT    = 56
	unsafeParserIDENTIFIER_STMT = 57
	unsafeParserWS_STMT         = 58
)

// unsafeParser rules.
const (
	unsafeParserRULE_template                 = 0
	unsafeParserRULE_content                  = 1
	unsafeParserRULE_statement                = 2
	unsafeParserRULE_commentStatement         = 3
	unsafeParserRULE_variableStatement        = 4
	unsafeParserRULE_ifStatement              = 5
	unsafeParserRULE_elseStatement            = 6
	unsafeParserRULE_loopStatement            = 7
	unsafeParserRULE_expressionStatement      = 8
	unsafeParserRULE_macroStatement           = 9
	unsafeParserRULE_pipelineExpression       = 10
	unsafeParserRULE_expression               = 11
	unsafeParserRULE_logicalOrExpression      = 12
	unsafeParserRULE_logicalAndExpression     = 13
	unsafeParserRULE_equalityExpression       = 14
	unsafeParserRULE_comparisonExpression     = 15
	unsafeParserRULE_primaryExpression        = 16
	unsafeParserRULE_varExpr                  = 17
	unsafeParserRULE_functionCall             = 18
	unsafeParserRULE_logicalOrExpressionStmt  = 19
	unsafeParserRULE_logicalAndExpressionStmt = 20
	unsafeParserRULE_equalityExpressionStmt   = 21
	unsafeParserRULE_comparisonExpressionStmt = 22
	unsafeParserRULE_primaryExpressionStmt    = 23
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
	p.SetState(52)
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
			p.SetState(50)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case unsafeParserTEXT:
				{
					p.SetState(48)
					p.Content()
				}

			case unsafeParserCOMMENT_START, unsafeParserEXPR_START, unsafeParserSTMT_START:
				{
					p.SetState(49)
					p.Statement()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

		}
		p.SetState(54)
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
		p.SetState(55)
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
	CommentStatement() ICommentStatementContext
	IfStatement() IIfStatementContext
	LoopStatement() ILoopStatementContext
	ExpressionStatement() IExpressionStatementContext
	MacroStatement() IMacroStatementContext
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

func (s *StatementContext) CommentStatement() ICommentStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommentStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommentStatementContext)
}

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
	p.SetState(63)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(57)
			p.CommentStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(58)
			p.IfStatement()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(59)
			p.LoopStatement()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(60)
			p.ExpressionStatement()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(61)
			p.MacroStatement()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(62)
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

// ICommentStatementContext is an interface to support dynamic dispatch.
type ICommentStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMENT_START() antlr.TerminalNode
	COMMENT_CONTENT() antlr.TerminalNode
	COMMENT_END() antlr.TerminalNode

	// IsCommentStatementContext differentiates from other interfaces.
	IsCommentStatementContext()
}

type CommentStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommentStatementContext() *CommentStatementContext {
	var p = new(CommentStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_commentStatement
	return p
}

func InitEmptyCommentStatementContext(p *CommentStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_commentStatement
}

func (*CommentStatementContext) IsCommentStatementContext() {}

func NewCommentStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommentStatementContext {
	var p = new(CommentStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_commentStatement

	return p
}

func (s *CommentStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *CommentStatementContext) COMMENT_START() antlr.TerminalNode {
	return s.GetToken(unsafeParserCOMMENT_START, 0)
}

func (s *CommentStatementContext) COMMENT_CONTENT() antlr.TerminalNode {
	return s.GetToken(unsafeParserCOMMENT_CONTENT, 0)
}

func (s *CommentStatementContext) COMMENT_END() antlr.TerminalNode {
	return s.GetToken(unsafeParserCOMMENT_END, 0)
}

func (s *CommentStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommentStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommentStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitCommentStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) CommentStatement() (localctx ICommentStatementContext) {
	localctx = NewCommentStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, unsafeParserRULE_commentStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(65)
		p.Match(unsafeParserCOMMENT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(66)
		p.Match(unsafeParserCOMMENT_CONTENT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(67)
		p.Match(unsafeParserCOMMENT_END)
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

// IVariableStatementContext is an interface to support dynamic dispatch.
type IVariableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STMT_START() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	COLON_EQUAL() antlr.TerminalNode
	Expression() IExpressionContext
	STMT_END() antlr.TerminalNode

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

func (s *VariableStatementContext) STMT_START() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_START, 0)
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

func (s *VariableStatementContext) STMT_END() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_END, 0)
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
	p.EnterRule(localctx, 8, unsafeParserRULE_variableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(70)
		p.Match(unsafeParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(71)
		p.Match(unsafeParserCOLON_EQUAL)
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
		p.Match(unsafeParserSTMT_END)
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
	AllSTMT_START() []antlr.TerminalNode
	STMT_START(i int) antlr.TerminalNode
	IF() antlr.TerminalNode
	Expression() IExpressionContext
	AllSTMT_END() []antlr.TerminalNode
	STMT_END(i int) antlr.TerminalNode
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

func (s *IfStatementContext) AllSTMT_START() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTMT_START)
}

func (s *IfStatementContext) STMT_START(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_START, i)
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

func (s *IfStatementContext) AllSTMT_END() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTMT_END)
}

func (s *IfStatementContext) STMT_END(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_END, i)
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
	p.EnterRule(localctx, 10, unsafeParserRULE_ifStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Match(unsafeParserIF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(77)
		p.Expression()
	}
	{
		p.SetState(78)
		p.Match(unsafeParserSTMT_END)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(79)
		p.Template()
	}
	p.SetState(81)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(80)
			p.ElseStatement()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(83)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(84)
		_la = p.GetTokenStream().LA(1)

		if !(_la == unsafeParserEND || _la == unsafeParserENDIF) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(85)
		p.Match(unsafeParserSTMT_END)
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
	STMT_START() antlr.TerminalNode
	ELSE() antlr.TerminalNode
	STMT_END() antlr.TerminalNode
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

func (s *ElseStatementContext) STMT_START() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_START, 0)
}

func (s *ElseStatementContext) ELSE() antlr.TerminalNode {
	return s.GetToken(unsafeParserELSE, 0)
}

func (s *ElseStatementContext) STMT_END() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_END, 0)
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
	p.EnterRule(localctx, 12, unsafeParserRULE_elseStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(87)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(88)
		p.Match(unsafeParserELSE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(89)
		p.Match(unsafeParserSTMT_END)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(90)
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
	AllSTMT_START() []antlr.TerminalNode
	STMT_START(i int) antlr.TerminalNode
	LOOP() antlr.TerminalNode
	AllIDENTIFIER_STMT() []antlr.TerminalNode
	IDENTIFIER_STMT(i int) antlr.TerminalNode
	IN() antlr.TerminalNode
	AllSTMT_END() []antlr.TerminalNode
	STMT_END(i int) antlr.TerminalNode
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

func (s *LoopStatementContext) AllSTMT_START() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTMT_START)
}

func (s *LoopStatementContext) STMT_START(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_START, i)
}

func (s *LoopStatementContext) LOOP() antlr.TerminalNode {
	return s.GetToken(unsafeParserLOOP, 0)
}

func (s *LoopStatementContext) AllIDENTIFIER_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserIDENTIFIER_STMT)
}

func (s *LoopStatementContext) IDENTIFIER_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER_STMT, i)
}

func (s *LoopStatementContext) IN() antlr.TerminalNode {
	return s.GetToken(unsafeParserIN, 0)
}

func (s *LoopStatementContext) AllSTMT_END() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTMT_END)
}

func (s *LoopStatementContext) STMT_END(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_END, i)
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
	p.EnterRule(localctx, 14, unsafeParserRULE_loopStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(93)
		p.Match(unsafeParserLOOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(94)
		p.Match(unsafeParserIDENTIFIER_STMT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(95)
		p.Match(unsafeParserIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(96)
		p.Match(unsafeParserIDENTIFIER_STMT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(97)
		p.Match(unsafeParserSTMT_END)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(98)
		p.Template()
	}
	{
		p.SetState(99)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(100)
		_la = p.GetTokenStream().LA(1)

		if !(_la == unsafeParserEND || _la == unsafeParserENDLOOP) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(101)
		p.Match(unsafeParserSTMT_END)
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
	EXPR_START() antlr.TerminalNode
	PipelineExpression() IPipelineExpressionContext
	EXPR_END() antlr.TerminalNode

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

func (s *ExpressionStatementContext) EXPR_START() antlr.TerminalNode {
	return s.GetToken(unsafeParserEXPR_START, 0)
}

func (s *ExpressionStatementContext) PipelineExpression() IPipelineExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPipelineExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPipelineExpressionContext)
}

func (s *ExpressionStatementContext) EXPR_END() antlr.TerminalNode {
	return s.GetToken(unsafeParserEXPR_END, 0)
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
	p.EnterRule(localctx, 16, unsafeParserRULE_expressionStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(unsafeParserEXPR_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(104)
		p.PipelineExpression()
	}
	{
		p.SetState(105)
		p.Match(unsafeParserEXPR_END)
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
	AllSTMT_START() []antlr.TerminalNode
	STMT_START(i int) antlr.TerminalNode
	MACRO() antlr.TerminalNode
	AllIDENTIFIER_STMT() []antlr.TerminalNode
	IDENTIFIER_STMT(i int) antlr.TerminalNode
	AllSTMT_END() []antlr.TerminalNode
	STMT_END(i int) antlr.TerminalNode
	Template() ITemplateContext
	END() antlr.TerminalNode
	ENDMACRO() antlr.TerminalNode
	LPAREN_STMT() antlr.TerminalNode
	RPAREN_STMT() antlr.TerminalNode
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

func (s *MacroStatementContext) AllSTMT_START() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTMT_START)
}

func (s *MacroStatementContext) STMT_START(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_START, i)
}

func (s *MacroStatementContext) MACRO() antlr.TerminalNode {
	return s.GetToken(unsafeParserMACRO, 0)
}

func (s *MacroStatementContext) AllIDENTIFIER_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserIDENTIFIER_STMT)
}

func (s *MacroStatementContext) IDENTIFIER_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER_STMT, i)
}

func (s *MacroStatementContext) AllSTMT_END() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserSTMT_END)
}

func (s *MacroStatementContext) STMT_END(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserSTMT_END, i)
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

func (s *MacroStatementContext) LPAREN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN_STMT, 0)
}

func (s *MacroStatementContext) RPAREN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN_STMT, 0)
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
	p.EnterRule(localctx, 18, unsafeParserRULE_macroStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(108)
		p.Match(unsafeParserMACRO)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(109)
		p.Match(unsafeParserIDENTIFIER_STMT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(122)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == unsafeParserLPAREN_STMT {
		{
			p.SetState(110)
			p.Match(unsafeParserLPAREN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(119)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == unsafeParserIDENTIFIER_STMT {
			{
				p.SetState(111)
				p.Match(unsafeParserIDENTIFIER_STMT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(116)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == unsafeParserCOMMA {
				{
					p.SetState(112)
					p.Match(unsafeParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(113)
					p.Match(unsafeParserIDENTIFIER_STMT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(118)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(121)
			p.Match(unsafeParserRPAREN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(124)
		p.Match(unsafeParserSTMT_END)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(125)
		p.Template()
	}
	{
		p.SetState(126)
		p.Match(unsafeParserSTMT_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(127)
		_la = p.GetTokenStream().LA(1)

		if !(_la == unsafeParserEND || _la == unsafeParserENDMACRO) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(128)
		p.Match(unsafeParserSTMT_END)
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

// IPipelineExpressionContext is an interface to support dynamic dispatch.
type IPipelineExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllPIPE() []antlr.TerminalNode
	PIPE(i int) antlr.TerminalNode

	// IsPipelineExpressionContext differentiates from other interfaces.
	IsPipelineExpressionContext()
}

type PipelineExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPipelineExpressionContext() *PipelineExpressionContext {
	var p = new(PipelineExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_pipelineExpression
	return p
}

func InitEmptyPipelineExpressionContext(p *PipelineExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_pipelineExpression
}

func (*PipelineExpressionContext) IsPipelineExpressionContext() {}

func NewPipelineExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PipelineExpressionContext {
	var p = new(PipelineExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_pipelineExpression

	return p
}

func (s *PipelineExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *PipelineExpressionContext) AllExpression() []IExpressionContext {
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

func (s *PipelineExpressionContext) Expression(i int) IExpressionContext {
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

func (s *PipelineExpressionContext) AllPIPE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserPIPE)
}

func (s *PipelineExpressionContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserPIPE, i)
}

func (s *PipelineExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PipelineExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PipelineExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitPipelineExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) PipelineExpression() (localctx IPipelineExpressionContext) {
	localctx = NewPipelineExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, unsafeParserRULE_pipelineExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(130)
		p.Expression()
	}
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserPIPE {
		{
			p.SetState(131)
			p.Match(unsafeParserPIPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(132)
			p.Expression()
		}

		p.SetState(137)
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

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode
	STRING() antlr.TerminalNode
	BOOLEAN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	NUMBER_STMT() antlr.TerminalNode
	STRING_STMT() antlr.TerminalNode
	BOOLEAN_STMT() antlr.TerminalNode
	VarExpr() IVarExprContext
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode
	FunctionCall() IFunctionCallContext
	LogicalOrExpression() ILogicalOrExpressionContext
	LogicalOrExpressionStmt() ILogicalOrExpressionStmtContext

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

func (s *ExpressionContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(unsafeParserNUMBER, 0)
}

func (s *ExpressionContext) STRING() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTRING, 0)
}

func (s *ExpressionContext) BOOLEAN() antlr.TerminalNode {
	return s.GetToken(unsafeParserBOOLEAN, 0)
}

func (s *ExpressionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *ExpressionContext) NUMBER_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserNUMBER_STMT, 0)
}

func (s *ExpressionContext) STRING_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTRING_STMT, 0)
}

func (s *ExpressionContext) BOOLEAN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserBOOLEAN_STMT, 0)
}

func (s *ExpressionContext) VarExpr() IVarExprContext {
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

func (s *ExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN, 0)
}

func (s *ExpressionContext) Expression() IExpressionContext {
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

func (s *ExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN, 0)
}

func (s *ExpressionContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *ExpressionContext) LogicalOrExpression() ILogicalOrExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalOrExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogicalOrExpressionContext)
}

func (s *ExpressionContext) LogicalOrExpressionStmt() ILogicalOrExpressionStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalOrExpressionStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogicalOrExpressionStmtContext)
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
	p.EnterRule(localctx, 22, unsafeParserRULE_expression)
	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(138)
			p.Match(unsafeParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(139)
			p.Match(unsafeParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(140)
			p.Match(unsafeParserBOOLEAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(141)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(142)
			p.Match(unsafeParserNUMBER_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(143)
			p.Match(unsafeParserSTRING_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(144)
			p.Match(unsafeParserBOOLEAN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(145)
			p.VarExpr()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(146)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(147)
			p.Expression()
		}
		{
			p.SetState(148)
			p.Match(unsafeParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(150)
			p.FunctionCall()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(151)
			p.LogicalOrExpression()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(152)
			p.LogicalOrExpressionStmt()
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

// ILogicalOrExpressionContext is an interface to support dynamic dispatch.
type ILogicalOrExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLogicalAndExpression() []ILogicalAndExpressionContext
	LogicalAndExpression(i int) ILogicalAndExpressionContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsLogicalOrExpressionContext differentiates from other interfaces.
	IsLogicalOrExpressionContext()
}

type LogicalOrExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalOrExpressionContext() *LogicalOrExpressionContext {
	var p = new(LogicalOrExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalOrExpression
	return p
}

func InitEmptyLogicalOrExpressionContext(p *LogicalOrExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalOrExpression
}

func (*LogicalOrExpressionContext) IsLogicalOrExpressionContext() {}

func NewLogicalOrExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalOrExpressionContext {
	var p = new(LogicalOrExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_logicalOrExpression

	return p
}

func (s *LogicalOrExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalOrExpressionContext) AllLogicalAndExpression() []ILogicalAndExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILogicalAndExpressionContext); ok {
			len++
		}
	}

	tst := make([]ILogicalAndExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILogicalAndExpressionContext); ok {
			tst[i] = t.(ILogicalAndExpressionContext)
			i++
		}
	}

	return tst
}

func (s *LogicalOrExpressionContext) LogicalAndExpression(i int) ILogicalAndExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalAndExpressionContext); ok {
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

	return t.(ILogicalAndExpressionContext)
}

func (s *LogicalOrExpressionContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserOR)
}

func (s *LogicalOrExpressionContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserOR, i)
}

func (s *LogicalOrExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalOrExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalOrExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitLogicalOrExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) LogicalOrExpression() (localctx ILogicalOrExpressionContext) {
	localctx = NewLogicalOrExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, unsafeParserRULE_logicalOrExpression)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(155)
		p.LogicalAndExpression()
	}
	p.SetState(160)
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
				p.SetState(156)
				p.Match(unsafeParserOR)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(157)
				p.LogicalAndExpression()
			}

		}
		p.SetState(162)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext())
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

// ILogicalAndExpressionContext is an interface to support dynamic dispatch.
type ILogicalAndExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllEqualityExpression() []IEqualityExpressionContext
	EqualityExpression(i int) IEqualityExpressionContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsLogicalAndExpressionContext differentiates from other interfaces.
	IsLogicalAndExpressionContext()
}

type LogicalAndExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalAndExpressionContext() *LogicalAndExpressionContext {
	var p = new(LogicalAndExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalAndExpression
	return p
}

func InitEmptyLogicalAndExpressionContext(p *LogicalAndExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalAndExpression
}

func (*LogicalAndExpressionContext) IsLogicalAndExpressionContext() {}

func NewLogicalAndExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalAndExpressionContext {
	var p = new(LogicalAndExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_logicalAndExpression

	return p
}

func (s *LogicalAndExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalAndExpressionContext) AllEqualityExpression() []IEqualityExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEqualityExpressionContext); ok {
			len++
		}
	}

	tst := make([]IEqualityExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEqualityExpressionContext); ok {
			tst[i] = t.(IEqualityExpressionContext)
			i++
		}
	}

	return tst
}

func (s *LogicalAndExpressionContext) EqualityExpression(i int) IEqualityExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEqualityExpressionContext); ok {
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

	return t.(IEqualityExpressionContext)
}

func (s *LogicalAndExpressionContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserAND)
}

func (s *LogicalAndExpressionContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserAND, i)
}

func (s *LogicalAndExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalAndExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalAndExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitLogicalAndExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) LogicalAndExpression() (localctx ILogicalAndExpressionContext) {
	localctx = NewLogicalAndExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, unsafeParserRULE_logicalAndExpression)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(163)
		p.EqualityExpression()
	}
	p.SetState(168)
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
				p.SetState(164)
				p.Match(unsafeParserAND)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(165)
				p.EqualityExpression()
			}

		}
		p.SetState(170)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext())
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

// IEqualityExpressionContext is an interface to support dynamic dispatch.
type IEqualityExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllComparisonExpression() []IComparisonExpressionContext
	ComparisonExpression(i int) IComparisonExpressionContext
	AllEQ() []antlr.TerminalNode
	EQ(i int) antlr.TerminalNode
	AllNE() []antlr.TerminalNode
	NE(i int) antlr.TerminalNode

	// IsEqualityExpressionContext differentiates from other interfaces.
	IsEqualityExpressionContext()
}

type EqualityExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEqualityExpressionContext() *EqualityExpressionContext {
	var p = new(EqualityExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_equalityExpression
	return p
}

func InitEmptyEqualityExpressionContext(p *EqualityExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_equalityExpression
}

func (*EqualityExpressionContext) IsEqualityExpressionContext() {}

func NewEqualityExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EqualityExpressionContext {
	var p = new(EqualityExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_equalityExpression

	return p
}

func (s *EqualityExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *EqualityExpressionContext) AllComparisonExpression() []IComparisonExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IComparisonExpressionContext); ok {
			len++
		}
	}

	tst := make([]IComparisonExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IComparisonExpressionContext); ok {
			tst[i] = t.(IComparisonExpressionContext)
			i++
		}
	}

	return tst
}

func (s *EqualityExpressionContext) ComparisonExpression(i int) IComparisonExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonExpressionContext); ok {
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

	return t.(IComparisonExpressionContext)
}

func (s *EqualityExpressionContext) AllEQ() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserEQ)
}

func (s *EqualityExpressionContext) EQ(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserEQ, i)
}

func (s *EqualityExpressionContext) AllNE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserNE)
}

func (s *EqualityExpressionContext) NE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserNE, i)
}

func (s *EqualityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EqualityExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitEqualityExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) EqualityExpression() (localctx IEqualityExpressionContext) {
	localctx = NewEqualityExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, unsafeParserRULE_equalityExpression)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(171)
		p.ComparisonExpression()
	}
	p.SetState(176)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(172)
				_la = p.GetTokenStream().LA(1)

				if !(_la == unsafeParserEQ || _la == unsafeParserNE) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(173)
				p.ComparisonExpression()
			}

		}
		p.SetState(178)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext())
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

// IComparisonExpressionContext is an interface to support dynamic dispatch.
type IComparisonExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPrimaryExpression() []IPrimaryExpressionContext
	PrimaryExpression(i int) IPrimaryExpressionContext
	AllGT() []antlr.TerminalNode
	GT(i int) antlr.TerminalNode
	AllLT() []antlr.TerminalNode
	LT(i int) antlr.TerminalNode
	AllGE() []antlr.TerminalNode
	GE(i int) antlr.TerminalNode
	AllLE() []antlr.TerminalNode
	LE(i int) antlr.TerminalNode

	// IsComparisonExpressionContext differentiates from other interfaces.
	IsComparisonExpressionContext()
}

type ComparisonExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonExpressionContext() *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_comparisonExpression
	return p
}

func InitEmptyComparisonExpressionContext(p *ComparisonExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_comparisonExpression
}

func (*ComparisonExpressionContext) IsComparisonExpressionContext() {}

func NewComparisonExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_comparisonExpression

	return p
}

func (s *ComparisonExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonExpressionContext) AllPrimaryExpression() []IPrimaryExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			len++
		}
	}

	tst := make([]IPrimaryExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPrimaryExpressionContext); ok {
			tst[i] = t.(IPrimaryExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExpressionContext) PrimaryExpression(i int) IPrimaryExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
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

	return t.(IPrimaryExpressionContext)
}

func (s *ComparisonExpressionContext) AllGT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserGT)
}

func (s *ComparisonExpressionContext) GT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserGT, i)
}

func (s *ComparisonExpressionContext) AllLT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserLT)
}

func (s *ComparisonExpressionContext) LT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserLT, i)
}

func (s *ComparisonExpressionContext) AllGE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserGE)
}

func (s *ComparisonExpressionContext) GE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserGE, i)
}

func (s *ComparisonExpressionContext) AllLE() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserLE)
}

func (s *ComparisonExpressionContext) LE(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserLE, i)
}

func (s *ComparisonExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitComparisonExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) ComparisonExpression() (localctx IComparisonExpressionContext) {
	localctx = NewComparisonExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, unsafeParserRULE_comparisonExpression)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(179)
		p.PrimaryExpression()
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(180)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&31457280) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(181)
				p.PrimaryExpression()
			}

		}
		p.SetState(186)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext())
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

// IPrimaryExpressionContext is an interface to support dynamic dispatch.
type IPrimaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode
	STRING() antlr.TerminalNode
	BOOLEAN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	NUMBER_STMT() antlr.TerminalNode
	STRING_STMT() antlr.TerminalNode
	BOOLEAN_STMT() antlr.TerminalNode
	VarExpr() IVarExprContext
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode
	FunctionCall() IFunctionCallContext
	NOT() antlr.TerminalNode
	PrimaryExpression() IPrimaryExpressionContext
	NOT_STMT() antlr.TerminalNode
	PrimaryExpressionStmt() IPrimaryExpressionStmtContext

	// IsPrimaryExpressionContext differentiates from other interfaces.
	IsPrimaryExpressionContext()
}

type PrimaryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExpressionContext() *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_primaryExpression
	return p
}

func InitEmptyPrimaryExpressionContext(p *PrimaryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_primaryExpression
}

func (*PrimaryExpressionContext) IsPrimaryExpressionContext() {}

func NewPrimaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_primaryExpression

	return p
}

func (s *PrimaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExpressionContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(unsafeParserNUMBER, 0)
}

func (s *PrimaryExpressionContext) STRING() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTRING, 0)
}

func (s *PrimaryExpressionContext) BOOLEAN() antlr.TerminalNode {
	return s.GetToken(unsafeParserBOOLEAN, 0)
}

func (s *PrimaryExpressionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER, 0)
}

func (s *PrimaryExpressionContext) NUMBER_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserNUMBER_STMT, 0)
}

func (s *PrimaryExpressionContext) STRING_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTRING_STMT, 0)
}

func (s *PrimaryExpressionContext) BOOLEAN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserBOOLEAN_STMT, 0)
}

func (s *PrimaryExpressionContext) VarExpr() IVarExprContext {
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

func (s *PrimaryExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN, 0)
}

func (s *PrimaryExpressionContext) Expression() IExpressionContext {
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

func (s *PrimaryExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN, 0)
}

func (s *PrimaryExpressionContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *PrimaryExpressionContext) NOT() antlr.TerminalNode {
	return s.GetToken(unsafeParserNOT, 0)
}

func (s *PrimaryExpressionContext) PrimaryExpression() IPrimaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionContext)
}

func (s *PrimaryExpressionContext) NOT_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserNOT_STMT, 0)
}

func (s *PrimaryExpressionContext) PrimaryExpressionStmt() IPrimaryExpressionStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionStmtContext)
}

func (s *PrimaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitPrimaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) PrimaryExpression() (localctx IPrimaryExpressionContext) {
	localctx = NewPrimaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, unsafeParserRULE_primaryExpression)
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(187)
			p.Match(unsafeParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(188)
			p.Match(unsafeParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(189)
			p.Match(unsafeParserBOOLEAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(190)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(191)
			p.Match(unsafeParserNUMBER_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(192)
			p.Match(unsafeParserSTRING_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(193)
			p.Match(unsafeParserBOOLEAN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(194)
			p.VarExpr()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(195)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(196)
			p.Expression()
		}
		{
			p.SetState(197)
			p.Match(unsafeParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(199)
			p.FunctionCall()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(200)
			p.Match(unsafeParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(201)
			p.PrimaryExpression()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(202)
			p.Match(unsafeParserNOT_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(203)
			p.PrimaryExpressionStmt()
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
	p.EnterRule(localctx, 34, unsafeParserRULE_varExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(206)
		p.Match(unsafeParserDOLLAR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(207)
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
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
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

func (s *FunctionCallContext) AllExpression() []IExpressionContext {
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

func (s *FunctionCallContext) Expression(i int) IExpressionContext {
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
	p.EnterRule(localctx, 36, unsafeParserRULE_functionCall)
	var _la int

	var _alt int

	p.SetState(229)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(209)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(213)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 14, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(210)
					p.Expression()
				}

			}
			p.SetState(215)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 14, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(216)
			p.Match(unsafeParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(217)
			p.Match(unsafeParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(219)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&279223198237606144) != 0 {
			{
				p.SetState(218)
				p.Expression()
			}

		}
		p.SetState(225)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == unsafeParserCOMMA {
			{
				p.SetState(221)
				p.Match(unsafeParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(222)
				p.Expression()
			}

			p.SetState(227)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(228)
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

// ILogicalOrExpressionStmtContext is an interface to support dynamic dispatch.
type ILogicalOrExpressionStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLogicalAndExpressionStmt() []ILogicalAndExpressionStmtContext
	LogicalAndExpressionStmt(i int) ILogicalAndExpressionStmtContext
	AllOR_STMT() []antlr.TerminalNode
	OR_STMT(i int) antlr.TerminalNode

	// IsLogicalOrExpressionStmtContext differentiates from other interfaces.
	IsLogicalOrExpressionStmtContext()
}

type LogicalOrExpressionStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalOrExpressionStmtContext() *LogicalOrExpressionStmtContext {
	var p = new(LogicalOrExpressionStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalOrExpressionStmt
	return p
}

func InitEmptyLogicalOrExpressionStmtContext(p *LogicalOrExpressionStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalOrExpressionStmt
}

func (*LogicalOrExpressionStmtContext) IsLogicalOrExpressionStmtContext() {}

func NewLogicalOrExpressionStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalOrExpressionStmtContext {
	var p = new(LogicalOrExpressionStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_logicalOrExpressionStmt

	return p
}

func (s *LogicalOrExpressionStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalOrExpressionStmtContext) AllLogicalAndExpressionStmt() []ILogicalAndExpressionStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILogicalAndExpressionStmtContext); ok {
			len++
		}
	}

	tst := make([]ILogicalAndExpressionStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILogicalAndExpressionStmtContext); ok {
			tst[i] = t.(ILogicalAndExpressionStmtContext)
			i++
		}
	}

	return tst
}

func (s *LogicalOrExpressionStmtContext) LogicalAndExpressionStmt(i int) ILogicalAndExpressionStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalAndExpressionStmtContext); ok {
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

	return t.(ILogicalAndExpressionStmtContext)
}

func (s *LogicalOrExpressionStmtContext) AllOR_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserOR_STMT)
}

func (s *LogicalOrExpressionStmtContext) OR_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserOR_STMT, i)
}

func (s *LogicalOrExpressionStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalOrExpressionStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalOrExpressionStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitLogicalOrExpressionStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) LogicalOrExpressionStmt() (localctx ILogicalOrExpressionStmtContext) {
	localctx = NewLogicalOrExpressionStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, unsafeParserRULE_logicalOrExpressionStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(231)
		p.LogicalAndExpressionStmt()
	}
	p.SetState(236)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserOR_STMT {
		{
			p.SetState(232)
			p.Match(unsafeParserOR_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(233)
			p.LogicalAndExpressionStmt()
		}

		p.SetState(238)
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

// ILogicalAndExpressionStmtContext is an interface to support dynamic dispatch.
type ILogicalAndExpressionStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllEqualityExpressionStmt() []IEqualityExpressionStmtContext
	EqualityExpressionStmt(i int) IEqualityExpressionStmtContext
	AllAND_STMT() []antlr.TerminalNode
	AND_STMT(i int) antlr.TerminalNode

	// IsLogicalAndExpressionStmtContext differentiates from other interfaces.
	IsLogicalAndExpressionStmtContext()
}

type LogicalAndExpressionStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalAndExpressionStmtContext() *LogicalAndExpressionStmtContext {
	var p = new(LogicalAndExpressionStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalAndExpressionStmt
	return p
}

func InitEmptyLogicalAndExpressionStmtContext(p *LogicalAndExpressionStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_logicalAndExpressionStmt
}

func (*LogicalAndExpressionStmtContext) IsLogicalAndExpressionStmtContext() {}

func NewLogicalAndExpressionStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalAndExpressionStmtContext {
	var p = new(LogicalAndExpressionStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_logicalAndExpressionStmt

	return p
}

func (s *LogicalAndExpressionStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *LogicalAndExpressionStmtContext) AllEqualityExpressionStmt() []IEqualityExpressionStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEqualityExpressionStmtContext); ok {
			len++
		}
	}

	tst := make([]IEqualityExpressionStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEqualityExpressionStmtContext); ok {
			tst[i] = t.(IEqualityExpressionStmtContext)
			i++
		}
	}

	return tst
}

func (s *LogicalAndExpressionStmtContext) EqualityExpressionStmt(i int) IEqualityExpressionStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEqualityExpressionStmtContext); ok {
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

	return t.(IEqualityExpressionStmtContext)
}

func (s *LogicalAndExpressionStmtContext) AllAND_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserAND_STMT)
}

func (s *LogicalAndExpressionStmtContext) AND_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserAND_STMT, i)
}

func (s *LogicalAndExpressionStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalAndExpressionStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalAndExpressionStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitLogicalAndExpressionStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) LogicalAndExpressionStmt() (localctx ILogicalAndExpressionStmtContext) {
	localctx = NewLogicalAndExpressionStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, unsafeParserRULE_logicalAndExpressionStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(239)
		p.EqualityExpressionStmt()
	}
	p.SetState(244)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserAND_STMT {
		{
			p.SetState(240)
			p.Match(unsafeParserAND_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(241)
			p.EqualityExpressionStmt()
		}

		p.SetState(246)
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

// IEqualityExpressionStmtContext is an interface to support dynamic dispatch.
type IEqualityExpressionStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllComparisonExpressionStmt() []IComparisonExpressionStmtContext
	ComparisonExpressionStmt(i int) IComparisonExpressionStmtContext
	AllEQ_STMT() []antlr.TerminalNode
	EQ_STMT(i int) antlr.TerminalNode
	AllNE_STMT() []antlr.TerminalNode
	NE_STMT(i int) antlr.TerminalNode

	// IsEqualityExpressionStmtContext differentiates from other interfaces.
	IsEqualityExpressionStmtContext()
}

type EqualityExpressionStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEqualityExpressionStmtContext() *EqualityExpressionStmtContext {
	var p = new(EqualityExpressionStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_equalityExpressionStmt
	return p
}

func InitEmptyEqualityExpressionStmtContext(p *EqualityExpressionStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_equalityExpressionStmt
}

func (*EqualityExpressionStmtContext) IsEqualityExpressionStmtContext() {}

func NewEqualityExpressionStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EqualityExpressionStmtContext {
	var p = new(EqualityExpressionStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_equalityExpressionStmt

	return p
}

func (s *EqualityExpressionStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *EqualityExpressionStmtContext) AllComparisonExpressionStmt() []IComparisonExpressionStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IComparisonExpressionStmtContext); ok {
			len++
		}
	}

	tst := make([]IComparisonExpressionStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IComparisonExpressionStmtContext); ok {
			tst[i] = t.(IComparisonExpressionStmtContext)
			i++
		}
	}

	return tst
}

func (s *EqualityExpressionStmtContext) ComparisonExpressionStmt(i int) IComparisonExpressionStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonExpressionStmtContext); ok {
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

	return t.(IComparisonExpressionStmtContext)
}

func (s *EqualityExpressionStmtContext) AllEQ_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserEQ_STMT)
}

func (s *EqualityExpressionStmtContext) EQ_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserEQ_STMT, i)
}

func (s *EqualityExpressionStmtContext) AllNE_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserNE_STMT)
}

func (s *EqualityExpressionStmtContext) NE_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserNE_STMT, i)
}

func (s *EqualityExpressionStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityExpressionStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EqualityExpressionStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitEqualityExpressionStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) EqualityExpressionStmt() (localctx IEqualityExpressionStmtContext) {
	localctx = NewEqualityExpressionStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, unsafeParserRULE_equalityExpressionStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(247)
		p.ComparisonExpressionStmt()
	}
	p.SetState(252)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == unsafeParserEQ_STMT || _la == unsafeParserNE_STMT {
		{
			p.SetState(248)
			_la = p.GetTokenStream().LA(1)

			if !(_la == unsafeParserEQ_STMT || _la == unsafeParserNE_STMT) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(249)
			p.ComparisonExpressionStmt()
		}

		p.SetState(254)
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

// IComparisonExpressionStmtContext is an interface to support dynamic dispatch.
type IComparisonExpressionStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPrimaryExpressionStmt() []IPrimaryExpressionStmtContext
	PrimaryExpressionStmt(i int) IPrimaryExpressionStmtContext
	AllGT_STMT() []antlr.TerminalNode
	GT_STMT(i int) antlr.TerminalNode
	AllLT_STMT() []antlr.TerminalNode
	LT_STMT(i int) antlr.TerminalNode
	AllGE_STMT() []antlr.TerminalNode
	GE_STMT(i int) antlr.TerminalNode
	AllLE_STMT() []antlr.TerminalNode
	LE_STMT(i int) antlr.TerminalNode

	// IsComparisonExpressionStmtContext differentiates from other interfaces.
	IsComparisonExpressionStmtContext()
}

type ComparisonExpressionStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonExpressionStmtContext() *ComparisonExpressionStmtContext {
	var p = new(ComparisonExpressionStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_comparisonExpressionStmt
	return p
}

func InitEmptyComparisonExpressionStmtContext(p *ComparisonExpressionStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_comparisonExpressionStmt
}

func (*ComparisonExpressionStmtContext) IsComparisonExpressionStmtContext() {}

func NewComparisonExpressionStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonExpressionStmtContext {
	var p = new(ComparisonExpressionStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_comparisonExpressionStmt

	return p
}

func (s *ComparisonExpressionStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonExpressionStmtContext) AllPrimaryExpressionStmt() []IPrimaryExpressionStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPrimaryExpressionStmtContext); ok {
			len++
		}
	}

	tst := make([]IPrimaryExpressionStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPrimaryExpressionStmtContext); ok {
			tst[i] = t.(IPrimaryExpressionStmtContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExpressionStmtContext) PrimaryExpressionStmt(i int) IPrimaryExpressionStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionStmtContext); ok {
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

	return t.(IPrimaryExpressionStmtContext)
}

func (s *ComparisonExpressionStmtContext) AllGT_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserGT_STMT)
}

func (s *ComparisonExpressionStmtContext) GT_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserGT_STMT, i)
}

func (s *ComparisonExpressionStmtContext) AllLT_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserLT_STMT)
}

func (s *ComparisonExpressionStmtContext) LT_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserLT_STMT, i)
}

func (s *ComparisonExpressionStmtContext) AllGE_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserGE_STMT)
}

func (s *ComparisonExpressionStmtContext) GE_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserGE_STMT, i)
}

func (s *ComparisonExpressionStmtContext) AllLE_STMT() []antlr.TerminalNode {
	return s.GetTokens(unsafeParserLE_STMT)
}

func (s *ComparisonExpressionStmtContext) LE_STMT(i int) antlr.TerminalNode {
	return s.GetToken(unsafeParserLE_STMT, i)
}

func (s *ComparisonExpressionStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExpressionStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonExpressionStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitComparisonExpressionStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) ComparisonExpressionStmt() (localctx IComparisonExpressionStmtContext) {
	localctx = NewComparisonExpressionStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, unsafeParserRULE_comparisonExpressionStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(255)
		p.PrimaryExpressionStmt()
	}
	p.SetState(260)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2111062325329920) != 0 {
		{
			p.SetState(256)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2111062325329920) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(257)
			p.PrimaryExpressionStmt()
		}

		p.SetState(262)
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

// IPrimaryExpressionStmtContext is an interface to support dynamic dispatch.
type IPrimaryExpressionStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER_STMT() antlr.TerminalNode
	STRING_STMT() antlr.TerminalNode
	BOOLEAN_STMT() antlr.TerminalNode
	IDENTIFIER_STMT() antlr.TerminalNode
	LPAREN_STMT() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN_STMT() antlr.TerminalNode
	NOT_STMT() antlr.TerminalNode
	PrimaryExpressionStmt() IPrimaryExpressionStmtContext

	// IsPrimaryExpressionStmtContext differentiates from other interfaces.
	IsPrimaryExpressionStmtContext()
}

type PrimaryExpressionStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExpressionStmtContext() *PrimaryExpressionStmtContext {
	var p = new(PrimaryExpressionStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_primaryExpressionStmt
	return p
}

func InitEmptyPrimaryExpressionStmtContext(p *PrimaryExpressionStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = unsafeParserRULE_primaryExpressionStmt
}

func (*PrimaryExpressionStmtContext) IsPrimaryExpressionStmtContext() {}

func NewPrimaryExpressionStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExpressionStmtContext {
	var p = new(PrimaryExpressionStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = unsafeParserRULE_primaryExpressionStmt

	return p
}

func (s *PrimaryExpressionStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExpressionStmtContext) NUMBER_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserNUMBER_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) STRING_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserSTRING_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) BOOLEAN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserBOOLEAN_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) IDENTIFIER_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserIDENTIFIER_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) LPAREN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserLPAREN_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) Expression() IExpressionContext {
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

func (s *PrimaryExpressionStmtContext) RPAREN_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserRPAREN_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) NOT_STMT() antlr.TerminalNode {
	return s.GetToken(unsafeParserNOT_STMT, 0)
}

func (s *PrimaryExpressionStmtContext) PrimaryExpressionStmt() IPrimaryExpressionStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionStmtContext)
}

func (s *PrimaryExpressionStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExpressionStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExpressionStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case unsafeParserVisitor:
		return t.VisitPrimaryExpressionStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *unsafeParser) PrimaryExpressionStmt() (localctx IPrimaryExpressionStmtContext) {
	localctx = NewPrimaryExpressionStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, unsafeParserRULE_primaryExpressionStmt)
	p.SetState(273)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case unsafeParserNUMBER_STMT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(263)
			p.Match(unsafeParserNUMBER_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserSTRING_STMT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(264)
			p.Match(unsafeParserSTRING_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserBOOLEAN_STMT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(265)
			p.Match(unsafeParserBOOLEAN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserIDENTIFIER_STMT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(266)
			p.Match(unsafeParserIDENTIFIER_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserLPAREN_STMT:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(267)
			p.Match(unsafeParserLPAREN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(268)
			p.Expression()
		}
		{
			p.SetState(269)
			p.Match(unsafeParserRPAREN_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case unsafeParserNOT_STMT:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(271)
			p.Match(unsafeParserNOT_STMT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(272)
			p.PrimaryExpressionStmt()
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
