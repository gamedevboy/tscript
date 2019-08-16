package lexer

import (
    "container/list"
    "io/ioutil"
    "path"
    "runtime"
    "strings"
    "unicode"

    "tklibs/script"
    "tklibs/script/compiler/token"
)

type Component struct {
    script.ComponentType

    content    string
    parseIndex int

    importFiles list.List
}

func NewLexer(owner interface{}) *Component {
    return &Component{ComponentType: script.MakeComponentType(owner)}
}

func (c *Component) ReadFile(fileName string) error {
    buf, err := ioutil.ReadFile(fileName)

    if err != nil {
        return err
    }

    c.content = string(buf)

    return nil
}

type contentState struct {
    line, column int
    content      []rune
    file         string
}

func skipWhitespacesAndLines(cs *contentState) bool {
    i, contentLen := 0, len(cs.content)

    for ; i < contentLen; i++ {
        if !(unicode.IsSpace(cs.content[i])) || cs.content[i] == '\r' || cs.content[i] == '\n' {
            break
        }

        cs.column++
    }

    for ; i < contentLen; i++ {
        if !(cs.content[i] == '\r' || cs.content[i] == '\n') {
            break
        }

        if cs.content[i] == '\n' {
            cs.line++
            cs.column = 0
        }
    }

    if i == 0 {
        return false
    }

    cs.content = cs.content[i:]
    return true
}

func convertStr(content []rune) []rune {
    j, contentLen := 0, len(content)
    ret := make([]rune, contentLen-2)

    for i := 1; i < contentLen-1; i++ {
        if content[i] == '\\' && i < contentLen-2 {
            i++
            switch content[i] {
            case 'n':
                ret[j] = '\n'
            case '\\':
                ret[j] = '\\'
            default:
                panic("unsupported \\ operations")
            }
        } else {
            ret[j] = content[i]
        }

        j++
    }

    return ret
}

func (c *Component) ParseFromRunes(file string, content []rune, tokenList *list.List) *list.List {
    cs := &contentState{content: content}
    cs.file = file

parseLoop:
    for len(cs.content) > 0 {
        t := parseToken(cs)

        tokenType := t.GetType()

        switch tokenType {
        case token.TokenTypeUnknown:
            break parseLoop
        case token.TokenTypeCommet:
            continue
        case token.TokenTypeIDENT:
            v := t.GetValue()
            if strings.IndexRune(v, '#') == 0 {
                switch v {
                case "#import":
                    if len(cs.content) > 0 {
                        t = parseToken(cs)

                        filePath := t.GetValue() //todo check value type

                        tl, err := c.ParseFile(path.Join(path.Dir(file), filePath), tokenList)
                        if err != nil {
                            continue
                        }

                        tokenList = tl
                        continue
                    } else {
                        panic("")
                    }
                }
            }
        }

        tokenList.PushBack(t)
    }

    return tokenList
}

func parseToken(cs *contentState) token.Token {
    currentTokenType := token.TokenTypeUnknown

    for {
        if !skipWhitespacesAndLines(cs) {
            break
        }
    }

    parseContent := cs.content

    i := 0
    length := len(parseContent)

scanLoop:
    for ; i < length; i++ {
        switch {
        case currentTokenType == token.TokenTypeSTRING:
            switch parseContent[i] {
            case '"':
                fallthrough
            case '\'':
                i++
                break scanLoop
            default:
                continue
            }
        case parseContent[i] == '"' || parseContent[i] == '\'':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeSTRING
            default:
                break scanLoop
            }
        case parseContent[i] == '#' || parseContent[i] == '_' || parseContent[i] == '$':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeIDENT
            case token.TokenTypeIDENT:
                continue
            default:
                break scanLoop
            }
        case parseContent[i] == '[':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeLBRACK
            default:
                break scanLoop
            }
        case parseContent[i] == '(':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeLPAREN
            default:
                break scanLoop
            }
        case parseContent[i] == '{':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeLBRACE
            default:
                break scanLoop
            }
        case parseContent[i] == ']':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeRBRACK
            default:
                break scanLoop
            }
        case parseContent[i] == ')':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeRPAREN
            default:
                break scanLoop
            }
        case parseContent[i] == '}':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeRBRACE
            default:
                break scanLoop
            }
        case parseContent[i] == ',':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeCOMMA
            default:
                break scanLoop
            }
        case parseContent[i] == '?':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeQUES
            default:
                break scanLoop
            }
        case parseContent[i] == ':':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeCOLON
            default:
                break scanLoop
            }
        case parseContent[i] == ';':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeSEMICOLON
            default:
                break scanLoop
            }
        case parseContent[i] == '+':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeADD
            case token.TokenTypeADD:
                currentTokenType = token.TokenTypeINC
            default:
                break scanLoop
            }
        case parseContent[i] == '-':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeSUB
            case token.TokenTypeSUB:
                currentTokenType = token.TokenTypeDEC
            default:
                break scanLoop
            }
        case parseContent[i] == '*':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeMUL
            case token.TokenTypeDIV:
                currentTokenType = token.TokenTypeCommet
                for i < length-1 {
                    for ; i < length-1; i++ {
                        if parseContent[i] == '\r' || parseContent[i] == '\n' {
                            break
                        } else {
                            if parseContent[i] == '*' && parseContent[i+1] == '/' {
                                i += 2
                                break scanLoop
                            }
                        }
                    }

                    for ; i < length; i++ {
                        if !(parseContent[i] == '\r' || parseContent[i] == '\n') {
                            break
                        }

                        if parseContent[i] == '\n' {
                            cs.line++
                            cs.column = 0
                        }
                    }
                }
            default:
                break scanLoop
            }
        case parseContent[i] == '/':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeDIV
            case token.TokenTypeDIV:
                currentTokenType = token.TokenTypeCommet
                for ; i < length; i++ {
                    if parseContent[i] == '\r' || parseContent[i] == '\n' {
                        break
                    }
                }

                for ; i < length; i++ {
                    if !(parseContent[i] == '\r' || parseContent[i] == '\n') {
                        break
                    }

                    if parseContent[i] == '\n' {
                        cs.line++
                        cs.column = 0
                    }
                }

                break scanLoop
            default:
                break scanLoop
            }
        case parseContent[i] == '%':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeREM
            default:
                break scanLoop
            }
        case parseContent[i] == '>':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeGREATER
            case token.TokenTypeASSIGN:
                currentTokenType = token.TokenTypeLAMBDA
            case token.TokenTypeGREATER:
                currentTokenType = token.TokenTypeSHR
            default:
                break scanLoop
            }
        case parseContent[i] == '<':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeLESS
            case token.TokenTypeLESS:
                currentTokenType = token.TokenTypeSHL
            default:
                break scanLoop
            }
        case parseContent[i] == '&':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeAND
            case token.TokenTypeAND:
                currentTokenType = token.TokenTypeLAND
            default:
                break scanLoop
            }
        case parseContent[i] == '|':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeOR
            case token.TokenTypeOR:
                currentTokenType = token.TokenTypeLOR
            default:
                break scanLoop
            }
        case parseContent[i] == '!':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeLNOT
            default:
                break scanLoop
            }
        case parseContent[i] == '^':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeXOR
            default:
                break scanLoop
            }
        case parseContent[i] == '=':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeASSIGN
            case token.TokenTypeADD:
                currentTokenType = token.TokenTypeADDASSIGN
            case token.TokenTypeSUB:
                currentTokenType = token.TokenTypeSUBASSIGN
            case token.TokenTypeMUL:
                currentTokenType = token.TokenTypeMULASSIGN
            case token.TokenTypeDIV:
                currentTokenType = token.TokenTypeDIVASSIGN
            case token.TokenTypeREM:
                currentTokenType = token.TokenTypeREMASSIGN
            case token.TokenTypeAND:
                currentTokenType = token.TokenTypeANDASSIGN
            case token.TokenTypeOR:
                currentTokenType = token.TokenTypeORASSIGN
            case token.TokenTypeXOR:
                currentTokenType = token.TokenTypeXORASSIGN
            case token.TokenTypeSHL:
                currentTokenType = token.TokenTypeSHLASSIGN
            case token.TokenTypeSHR:
                currentTokenType = token.TokenTypeSHRASSIGN
            case token.TokenTypeGREATER:
                currentTokenType = token.TokenTypeGEQ
            case token.TokenTypeLESS:
                currentTokenType = token.TokenTypeLEQ
            case token.TokenTypeLNOT:
                currentTokenType = token.TokenTypeNEQ
            case token.TokenTypeASSIGN:
                currentTokenType = token.TokenTypeEQL
            default:
                break scanLoop
            }
        case parseContent[i] == '.':
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypePERIOD
            case token.TokenTypeINT:
                currentTokenType = token.TokenTypeFLOAT
            case token.TokenTypePERIOD:
                if i < length-1 && parseContent[i+1] == '.' {
                    currentTokenType = token.TokenTypeELLIPSIS
                    i++
                } else {
                    panic("invalid dot")
                }
            default:
                break scanLoop
            }
        case parseContent[i] == 'x':
            switch currentTokenType {
            case token.TokenTypeINT:
                continue
            }
            fallthrough
        case unicode.IsLetter(parseContent[i]) || parseContent[i] == '_' || parseContent[i] == '$':
            if parseContent[i] == '_' {
                runtime.Breakpoint()
            }
            switch currentTokenType {
            case token.TokenTypeUnknown:
                currentTokenType = token.TokenTypeIDENT
            case token.TokenTypeFLOAT:
                currentTokenType = token.TokenTypeINT
                i--
                break scanLoop
            case token.TokenTypeIDENT:
                continue
            default:
                break scanLoop
            }
        case unicode.IsDigit(parseContent[i]):
            switch currentTokenType {
            case token.TokenTypeUnknown, token.TokenTypeSUB:
                currentTokenType = token.TokenTypeINT
            case token.TokenTypeINT, token.TokenTypeFLOAT, token.TokenTypeIDENT:
                continue
            default:
                break scanLoop
            }
        default:
            break scanLoop
        }
    }

    cs.content = parseContent[i:]

    defer func() {
        cs.column += i
    }()

    retContent := parseContent[:i]

    if currentTokenType == token.TokenTypeSTRING {
        retContent = convertStr(parseContent[:i])
    }

    return token.CreateToken(retContent, currentTokenType, cs.line+1, cs.column+1, cs.file)
}

func (c *Component) ParseFile(fileName string, tokenList *list.List) (*list.List, error) {
    for it := c.importFiles.Front(); it != nil; it = it.Next() {
        if it.Value.(string) == fileName {
            return tokenList, nil
        }
    }

    c.importFiles.PushBack(fileName)

    err := c.ReadFile(fileName)
    if err != nil {
        return tokenList, err
    }

    return c.ParseFromRunes(fileName, []rune(c.content), tokenList), nil
}