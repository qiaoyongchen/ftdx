package l

// 一个环境就是一个map
// 用于一个key 和 一个 object 进行关联
type env struct {
	store map[string]obj
	outer *env // 外层环境
}

func newenv() *env {
	s := make(map[string]obj)
	return &env{store: s}
}

// 通过传入A *env 新建 B *env
// A 在 B 的外层
// 通过这种方式模拟闭包: A 是函数定义时的外环境, B 是函数执行时的内环境
func newenclosedenv(outer *env) *env {
	env := newenv()
	env.outer = outer
	return env
}

// get : 先从自己找,找不到再向外层找
func (e *env) get(name string) (obj, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.get(name)
	}
	return obj, ok
}

// set
func (e *env) set(name string, val obj) obj {
	e.store[name] = val
	return val
}
