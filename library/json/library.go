package json

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"text/scanner"

	"tklibs/script"
	"tklibs/script/runtime"
	"tklibs/script/runtime/native"
	"tklibs/script/value"
)

type library struct {
	context interface{}
	Encode  native.FunctionType
	Decode  native.FunctionType
}

func (*library) GetName() string {
	return "json"
}

func (l *library) SetScriptContext(context interface{}) {
	l.context = context
}

func NewLibrary() *library {
	ret := &library{}
	ret.init()
	return ret
}

func (l *library) init() {
	l.Encode = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return ""
		}

		switch val := args[0].(type) {
		case script.Value:
			return script.String(value.ToJsonString(val.Get()))
		default:
			return script.String(value.ToJsonString(val))
		}
	}
	l.Decode = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		r := script.Value{}

		if len(args) < 1 {
			return script.Null
		}

		s := &scanner.Scanner{}
		switch val := args[0].(type) {
		case string:
			s.Init(strings.NewReader(val))
		case script.String:
			s.Init(strings.NewReader(string(val)))
		}

		s.Mode = scanner.ScanStrings | scanner.ScanInts | scanner.ScanFloats | scanner.ScanIdents |
			scanner.SkipComments | scanner.ScanComments

		var scanValue func() interface{}

		scanValue = func() interface{} {
			tokenType := s.Scan()
			if tokenType == scanner.EOF {
				return script.Null
			}
			val := s.TokenText()
			scriptContext := l.context.(runtime.ScriptContext)
			switch val {
			case "[":
				sa := scriptContext.NewScriptArray(0)

			arrayParserLoop:
				for {
					if s.Peek() == scanner.EOF {
						break
					}

					if s.Peek() == ']' {
						s.Scan()
						break arrayParserLoop
					}

					nextValue := scanValue()
					if nextValue != nil {
						sa.(script.Array).Push(r.Set(nextValue))
					}

					if s.Scan() == scanner.EOF {
						break
					}

					switch s.TokenText() {
					case "]":
						break arrayParserLoop
					case ",":
						continue arrayParserLoop
					default:
						token := s.TokenText()
						panic(fmt.Errorf("Json.Decode Excepting , or ] not %v", token))
					}
				}
				return sa
			case "{":
				sm := scriptContext.NewScriptMap(0).(script.Map)
			mapParserLoop:
				for {
					tokenType := s.Scan()
					if tokenType == scanner.EOF {
						break
					}

					val := s.TokenText()
					switch tokenType {
					case scanner.String:
						// found key
						if s.Scan() == scanner.EOF { // skip ":"
							break mapParserLoop
						}

						if s.TokenText() != ":" {
							panic(fmt.Errorf("excepting : not %v", s.TokenText()))
						}

						sm.Set(script.String(*scriptContext.GetStringPool().Insert(strings.Trim(val, "\""))), scanValue())
					default:
						switch val {
						case "}":
							break mapParserLoop
						case ",":
							continue mapParserLoop
						default:
							panic("")
						}
					}
				}
				return sm
			case "-":
				switch val := scanValue().(type) {
				case script.Int:
					return -val
				case script.Int64:
					return -val
				case script.Float:
					return -val
				case script.Float64:
					return -val
				default:
					panic("unknow type")
				}
			default:
				switch tokenType {
				case scanner.Int:
					if v, err := strconv.ParseInt(val, 10, 64); err == nil {
						if v > math.MaxInt32 || v < math.MinInt32 {
							return script.Int64(v)
						}

						return script.Int(v)
					} else {
						panic(err)
					}
				case scanner.Float:
					if v, err := strconv.ParseFloat(val, 64); err == nil {
						return script.Float64(v)
					} else {
						panic(err)
					}
				case scanner.String:
					return script.String(strings.Trim(val, "\""))
				case scanner.Ident:
					switch val {
					case "true":
						return script.Bool(true)
					case "false":
						return script.Bool(false)
					case "null":
						return script.Null
					}

					panic(fmt.Errorf("Json.Decode Unknown token: %v", val))
				}
			}

			return script.Null
		}

		return scanValue()
	}
}
