package script

type scriptException struct {
	context   interface{}
	exception interface{}
}

func (s *scriptException) GetException() interface{} {
	return s.exception
}

func (s *scriptException) Error() string {
	return ""
}

type ScriptException interface {
	GetException() interface{}
}

var _ error = &scriptException{}
var _ ScriptException = &scriptException{}

func MakeException(context, exception interface{}) ScriptException {
	return &scriptException{
		context,
		exception,
	}
}
