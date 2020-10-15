package math

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"tklibs/script"
	"tklibs/script/runtime/native"
)

type library struct {
	context interface{}
	ToInt,
	MaxInt32,
	Max,
	Rand,
	Log native.FunctionType
}

func (*library) GetName() string {
	return "math"
}

func (l *library) SetScriptContext(context interface{}) {
	l.context = context
}

func NewLibrary() *library {
	ret := &library{}
	ret.init()
	return ret
}

func toInt(a interface{}) script.Int {
	switch v :=a.(type) {
	case script.Int:
		return v
	case script.Float:
		return script.Int(v)
	case script.Bool:
		if v {
			return 1
		}
		return 0
	case script.String:
		val, _ := strconv.Atoi(string(v))
		return script.Int(val)
	default:
		return 0
	}
}

func (l *library) init() {
	l.ToInt = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return script.Null
		}
		return toInt(args[0])
	}

	l.Rand = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		return script.Float(rand.New(rand.NewSource(time.Now().UnixNano())).Float32())
	}

	l.Log = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		return script.Float(math.Log(float64(args[1].(script.Float))) / math.Log(float64(args[0].(script.Float))))
	}

	l.Max = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return 0
		}
		if len(args) < 2 {
			return toInt(args[0])
		}
		a1 := toInt(args[0])
		a2 := toInt(args[1])
		if a1 < a2 {
			return a2
		}
		return a1
	}

	l.MaxInt32 = func(context interface{}, this interface{}, args ...interface{}) interface{} {
		return math.MaxInt32
	}
}
