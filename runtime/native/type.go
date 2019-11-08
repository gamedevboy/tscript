package native

type Type interface {
    New(args ...interface{}) interface{}
}
