package lexer

import (
	"container/list"
	"fmt"
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
	stringStack  []rune
}

func skipWhitespacesAndLines(cs *contentState) int {
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
		return i
	}

	cs.content = cs.content[i:]
	return i
}

func convertStr(content []rune) []rune {
	j, contentLen := 0, len(content)
	ret := make([]rune, 0, contentLen)

	for i := 1; i < contentLen-1; i++ {
		if content[i] == '\\' && i < contentLen-2 {
			i++
			switch content[i] {
			case 'n':
				ret = append(ret, '\n')
			case '"':
				ret = append(ret, '"')
			case '\'':
				ret = append(ret, '\'')
			case '\\':
				ret = append(ret, '\\')
			default:
				panic("unsupported \\ operations")
			}
		} else {
			ret = append(ret, content[i])
		}

		j++
	}

	return ret
}

func (c *Component) ParseFromRunes(file string, useImport bool, content []rune, tokenList *list.List) *list.List {
	cs := &contentState{content: content}
	cs.file = file

parseLoop:
	for len(cs.content) > 0 {
		t := parseToken(cs)

		tokenType := t.GetType()

		switch tokenType {
		case token.TokenTypeUnknown:
			break parseLoop
		// case token.TokenTypeCommet:
		// 	lastToken := tokenList.Back()
		// 	if lastToken != nil {
		// 		lt := lastToken.Value.(token.Token)
		// 		if lt.GetLine() == t.GetLine() && lt.GetFilePath() == t.GetFilePath() {
		// 			lt.SetComment(t.GetValue())
		// 			continue
		// 		}
		// 	}
		case token.TokenTypeIDENT:
			v := t.GetValue()
			if strings.IndexRune(v, '#') == 0 && useImport {
				switch v {
				case "#import":
					if len(cs.content) > 0 {
						t = parseToken(cs)

						filePath := t.GetValue() // todo check value type

						tl, err := c.ParseFile(path.Join(path.Dir(file), filePath), useImport, tokenList)
						if err != nil {
							continue
						}

						tokenList = tl
						continue
					} else {
						panic(fmt.Errorf("unsupport [%v]", v))
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
		if skipWhitespacesAndLines(cs) == 0 {
			break
		}
	}

	parseContent := cs.content

	i := 0
	length := len(parseContent)

scanLoop:
	for ; i < length; i++ {
		curRune := parseContent[i]

		switch {
		case currentTokenType == token.TokenTypeSTRING:
			switch curRune {
			case '\\':
				i++
			case '"', '\'':
				topStringRune := cs.stringStack[len(cs.stringStack)-1]
				if topStringRune == curRune {
					cs.stringStack = cs.stringStack[:len(cs.stringStack)-1]
					if 0 == len(cs.stringStack) {
						i++
						break scanLoop
					}
				} else {
					// cs.stringStack = append(cs.stringStack, curRune)
					continue
				}
			default:
				continue
			}
		case curRune == '"' || curRune == '\'':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeSTRING
				cs.stringStack = append(cs.stringStack, curRune)
			default:
				break scanLoop
			}
		case curRune == '#' || curRune == '_' || curRune == '$':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeIDENT
			case token.TokenTypeIDENT:
				continue
			default:
				break scanLoop
			}
		case curRune == '[':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeLBRACK
			default:
				break scanLoop
			}
		case curRune == '(':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeLPAREN
			default:
				break scanLoop
			}
		case curRune == '{':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeLBRACE
			default:
				break scanLoop
			}
		case curRune == ']':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeRBRACK
			default:
				break scanLoop
			}
		case curRune == ')':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeRPAREN
			default:
				break scanLoop
			}
		case curRune == '}':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeRBRACE
			default:
				break scanLoop
			}
		case curRune == ',':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeCOMMA
			default:
				break scanLoop
			}
		case curRune == '?':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeQUES
			case token.TokenTypeQUES:
				currentTokenType = token.TokenTypeNULLISH
			default:
				break scanLoop
			}
		case curRune == ':':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeCOLON
			default:
				break scanLoop
			}
		case curRune == ';':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeSEMICOLON
			default:
				break scanLoop
			}
		case curRune == '+':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeADD
			case token.TokenTypeADD:
				currentTokenType = token.TokenTypeINC
			default:
				break scanLoop
			}
		case curRune == '-':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeSUB
			case token.TokenTypeSUB:
				currentTokenType = token.TokenTypeDEC
			default:
				break scanLoop
			}
		case curRune == '*':
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
		case curRune == '/':
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

				i--
				break scanLoop
			default:
				break scanLoop
			}
		case curRune == '%':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeREM
			default:
				break scanLoop
			}
		case curRune == '>':
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
		case curRune == '<':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeLESS
			case token.TokenTypeLESS:
				currentTokenType = token.TokenTypeSHL
			default:
				break scanLoop
			}
		case curRune == '&':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeAND
			case token.TokenTypeAND:
				currentTokenType = token.TokenTypeLAND
			default:
				break scanLoop
			}
		case curRune == '|':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeOR
			case token.TokenTypeOR:
				currentTokenType = token.TokenTypeLOR
			default:
				break scanLoop
			}
		case curRune == '!':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeLNOT
			default:
				break scanLoop
			}
		case curRune == '^':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypeXOR
			default:
				break scanLoop
			}
		case curRune == '=':
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
		case curRune == '.':
			switch currentTokenType {
			case token.TokenTypeUnknown:
				currentTokenType = token.TokenTypePERIOD
			case token.TokenTypeINT:
				currentTokenType = token.TokenTypeFLOAT
			case token.TokenTypeQUES:
				currentTokenType = token.TokenTypeOptPERIOD
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
		case curRune == 'x':
			switch currentTokenType {
			case token.TokenTypeINT:
				continue
			}
			fallthrough
		case unicode.IsLetter(curRune) || curRune == '_' || curRune == '$':
			if curRune == '_' {
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
		case unicode.IsDigit(curRune):
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

func (c *Component) ParseFile(fileName string, useImport bool, tokenList *list.List) (*list.List, error) {
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

	return c.ParseFromRunes(fileName, useImport, []rune(c.content), tokenList), nil
}
