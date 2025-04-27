package object

type Function struct {
	value func(args ...Value) Value
}

func BuiltinFunction(value func(args ...Value) Value) Value {
	return nil
}

func (*Function) Name() {

}
