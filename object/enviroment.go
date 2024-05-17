package object

type Enviroment struct {
	store map[string]Object
	outer *Enviroment
}

func NewEnviroment() *Enviroment {
	s := make(map[string]Object)
	return &Enviroment{store: s, outer: nil}
}

func (e Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Enviroment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func NewEnclosedEnviroment(e *Enviroment) *Enviroment {
	env := NewEnviroment()
	env.outer = e
	return env
}
