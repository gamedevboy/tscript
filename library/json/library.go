package json

import (
    "fmt"
    "math"
    "strconv"
    "strings"
    "text/scanner"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
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

var Library = &library{}

func init() {
    Library.Encode = func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return ""
        }

        switch val := args[0].(type) {
        case script.Value:
            return script.String(value.ToJsonString(val))
        default:
            return script.String(value.ToJsonString(script.InterfaceToValue(val)))
        }
    }
    Library.Decode = func(this interface{}, args ...interface{}) interface{} {
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
            switch val {
            case "[":
                sa := Library.context.(runtime.ScriptContext).NewScriptArray(0)

            arrayParserLoop:
                for {
                    if s.Peek() == scanner.EOF {
                        break
                    }

                    sa.(script.Array).Push(r.Set(scanValue()))

                    if s.Scan() == scanner.EOF {
                        break
                    }

                    switch s.TokenText() {
                    case "]":
                        break arrayParserLoop
                    case ",":
                        continue arrayParserLoop
                    default:
                        panic(fmt.Errorf("Json.Decode Excepting , or ] not %v", s.TokenText()))
                    }
                }
                return sa
            case "{":
                sm := Library.context.(runtime.ScriptContext).NewScriptMap(0)
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
                        sm.(script.Object).ScriptSet(strings.Replace(val, "\"", "", -1), r.Set(scanValue()))
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
                case script.Float:
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
                    return script.String(strings.Replace(val, "\"", "", -1))
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