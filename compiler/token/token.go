package token

import (
    "fmt"
)

type TokenType uint8

const (
    TokenTypeUnknown   TokenType = iota
    TokenTypeCommet     // commet
    TokenTypeIDENT      // main, var, etc
    TokenTypeINT        // 12345
    TokenTypeFLOAT      // 123.45
    TokenTypeSTRING     // "abc"
    TokenTypeNULLISH    // ??
    TokenTypeCOMMA      // ,
    TokenTypeELLIPSIS   // ...
    TokenTypeORASSIGN   // |=
    TokenTypeXORASSIGN  // ^=
    TokenTypeANDASSIGN  // &=
    TokenTypeSHRASSIGN  // >>=
    TokenTypeSHLASSIGN  // <<=
    TokenTypeREMASSIGN  // %=
    TokenTypeMULASSIGN  // *=
    TokenTypeDIVASSIGN  // /=
    TokenTypeADDASSIGN  // +=
    TokenTypeSUBASSIGN  // -=
    TokenTypeASSIGN     // =
    TokenTypeQUES       // ?
    TokenTypeLOR        // ||
    TokenTypeLAND       // &&
    TokenTypeOR         // |
    TokenTypeXOR        // ^
    TokenTypeAND        // &
    TokenTypeNEQ        // !=
    TokenTypeEQL        // ==
    TokenTypeGEQ        // >=
    TokenTypeGREATER    // >
    TokenTypeLEQ        // <=
    TokenTypeLESS       // <
    TokenTypeSHL        // <<
    TokenTypeSHR        // >
    TokenTypeSUB        // -
    TokenTypeADD        // +
    TokenTypeREM        // %
    TokenTypeDIV        // /
    TokenTypeMUL        // *
    TokenTypeNOT        // ~
    TokenTypeLNOT       // !
    TokenTypeINC        // ++
    TokenTypeDEC        // --
    TokenTypeLPAREN     // (
    TokenTypeLBRACK     // [
    TokenTypeLBRACE     // {
    TokenTypePERIOD     // .
    TokenTypeOptPERIOD  // ?.
    TokenTypeRPAREN     // )
    TokenTypeRBRACK     // ]
    TokenTypeRBRACE     // }
    TokenTypeSEMICOLON  // ;
    TokenTypeCOLON      // :
    TokenTypeLAMBDA     // =>
)

func (t TokenType) WithAssign() bool {
    switch t {
    case TokenTypeADDASSIGN,
        TokenTypeSUBASSIGN,
        TokenTypeMULASSIGN,
        TokenTypeDIVASSIGN,
        TokenTypeANDASSIGN,
        TokenTypeORASSIGN,
        TokenTypeXORASSIGN,
        TokenTypeREMASSIGN,
        TokenTypeSHLASSIGN,
        TokenTypeSHRASSIGN:
        return true
    default:
        return false
    }
}

func (t TokenType) String() string {
    switch t {
    case TokenTypeCOMMA:
        return "," // ,
    case TokenTypeELLIPSIS:
        return "..." // ...
    case TokenTypeORASSIGN:
        return "|=" // |=
    case TokenTypeXORASSIGN:
        return "^=" // ^=
    case TokenTypeANDASSIGN:
        return "&=" // &=
    case TokenTypeSHRASSIGN:
        return ">>=" // >>=
    case TokenTypeSHLASSIGN:
        return "<<=" // <<=
    case TokenTypeREMASSIGN:
        return "%=" // %=
    case TokenTypeMULASSIGN:
        return "*=" // *=
    case TokenTypeDIVASSIGN:
        return "/=" // /=
    case TokenTypeADDASSIGN:
        return "+=" // +=
    case TokenTypeSUBASSIGN:
        return "-=" // -=
    case TokenTypeASSIGN:
        return "=" // =
    case TokenTypeQUES:
        return "?" // ?
    case TokenTypeNULLISH: // ??
        return "??"
    case TokenTypeLOR:
        return "||" // ||
    case TokenTypeLAND:
        return "&&" // &&
    case TokenTypeOR:
        return "|" // |
    case TokenTypeXOR:
        return "^" // ^
    case TokenTypeAND:
        return "&" // &
    case TokenTypeNEQ:
        return "!=" // !=
    case TokenTypeEQL:
        return "==" // ==
    case TokenTypeGEQ:
        return ">=" // >=
    case TokenTypeGREATER:
        return ">" // >
    case TokenTypeLEQ:
        return "<=" // <=
    case TokenTypeLESS:
        return "<" // <
    case TokenTypeSHL:
        return "<<" // <<
    case TokenTypeSHR:
        return ">>" // >>
    case TokenTypeADD:
        return "+" // +
    case TokenTypeSUB:
        return "-" // -
    case TokenTypeREM:
        return "%" // %
    case TokenTypeMUL:
        return "*" // *
    case TokenTypeDIV:
        return "/" // /
    case TokenTypeNOT:
        return "~" // ~
    case TokenTypeLNOT:
        return "!" // !
    case TokenTypeINC:
        return "++" // ++
    case TokenTypeDEC:
        return "__" // --
    case TokenTypeLPAREN:
        return "(" // (
    case TokenTypeLBRACK:
        return "[" // [
    case TokenTypeLBRACE:
        return "{" // {
    case TokenTypePERIOD:
        return "." // .
    case TokenTypeRPAREN:
        return ")" // )
    case TokenTypeRBRACK:
        return "]" // ]
    case TokenTypeRBRACE:
        return "}" // }
    case TokenTypeSEMICOLON:
        return ";" // ;
    case TokenTypeCOLON:
        return ":" // :
    default:
        return ""
    }
}

type Token interface {
    GetType() TokenType
    GetValue() string
    GetLine() int
    GetColumn() int
    GetFilePath() string
}

type token struct {
    tokenType    TokenType
    value        string
    line, column int
    file         string
}

func (t *token) GetType() TokenType {
    return t.tokenType
}

func (t *token) GetValue() string {
    return t.value
}

func (t *token) String() string {
    return fmt.Sprint("Value: ", t.value, " Type: ", t.tokenType, " Line: ", t.line, " Column: ", t.column)
}

func (t *token) GetLine() int {
    return t.line
}

func (t *token) GetColumn() int {
    return t.column
}

func (t *token) GetFilePath() string {
    return t.file
}

func CreateToken(content []rune, tokenType TokenType, line, column int, file string) Token {
    return &token{tokenType, string(content), line, column, file}
}
