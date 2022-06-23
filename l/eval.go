package l

import (
	"fmt"
	"strings"
)

var (
	NULLOBJ  = &nullobj{}                // null
	TRUEOBJ  = &booleanobj{value: true}  // true
	FALSEOBJ = &booleanobj{value: false} // false
)

// 通过 GO 类型 系统的true/false值
// 返回全局构造的 TRUEOBJ/FALSEOBJ
func bool2boolobj(input bool) *booleanobj {
	if input {
		return TRUEOBJ
	}
	return FALSEOBJ
}

// 执行 Node (Statement | Expression)
// 新增一个执行中环境,用于关联变量
func eval(node node, env *env) obj {
	switch node := node.(type) {

	// 语句列表
	case *program:
		return evalprogram(node, env)

	// 表达式语句
	case *expressionstatement:
		return eval(node.expression, env)

	// 整型
	case *integerliteral:
		return &integerobj{value: node.value}

	// 浮点类型
	case *floatliteral:
		return &floatobj{value: node.value}

	// 布尔类型
	case *boolean:
		return bool2boolobj(node.value)

	//前缀表达式
	case *prefixexpression:

		// 这里传进来的可能是很多奇怪的东西(boolen, integer, null ....)
		// 需要兼容这些, 所以先把执行出来结果再进行前缀操作
		right := eval(node.right, env)
		if iserror(right) {
			return right
		}
		return evalprefix(node.operator, right)

	// 中缀表达式
	// 先分别求出左，右表达式再进行计算
	case *infixexpression:

		// 先处理赋值
		// 赋值比较特殊 左边必须是标识符
		// TODO ':' 未做处理 ':'在通达信中的意思是赋值并输出值, 以后再处理，这个贼简单，后面再说
		if node.operator == ":=" || node.operator == ":" {
			right := eval(node.right, env)
			id, ok := node.left.(*identifier)
			if !ok {
				return newerror("assign name is not identifier:" + node.left.tokenliteral() + ", " + node.left.string())
			}
			return env.set(id.value, right)
		}

		// 中缀操作符左边求值
		left := eval(node.left, env)
		if iserror(left) {
			return left
		}

		// 中缀操作符右边求值
		right := eval(node.right, env)
		if iserror(right) {
			return right
		}

		// 执行中缀操作
		return evalInfixexpression(node.operator, left, right)

	// 执行标识符的时候,需要传入环境
	// 即在环境中取值然后执行
	case *identifier:
		return evalidentifer(node, env)

	// 函数调用
	case *callexpression:
		// 解析出object.Function类型
		function := eval(node.function, env)
		if iserror(function) {
			return function
		}

		// 运行参数表达式,解析[]object.Object做为参数
		args := evalexpressions(node.arguments, env)
		if len(args) == 1 && iserror(args[0]) {
			return args[0]
		}
		return applyfunction(function, args)

	}
	return nil
}

// 暂时只支持内置函数不支持自定义函数
// 自定义函数还没实现,不过通达信好像也不支持
func applyfunction(fn obj, args []obj) obj {
	switch fn := fn.(type) {
	case *builtinfnobj: // 内置函数
		return fn.fn(args...)
	default:
		return newerror("not a function %s", fn.tYpe())
	}
}

// 解析表达式列表
// 用于解析函数的参数列表
// 和数组中表达式列表
func evalexpressions(exps []expression, env *env) []obj {
	// 解析表达式的结果列表
	var result []obj
	// 挨个解析表达式,并加入到结果列表中
	for _, e := range exps {
		// 执行表达式
		evaluated := eval(e, env)
		// 执行错误直接返回
		if iserror(evaluated) {
			return []obj{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// 执行程序
func evalprogram(program *program, env *env) obj {
	var result obj
	for _, statement := range program.statements {
		result = eval(statement, env)
		switch result := result.(type) {
		case *errorobj:
			return result
		}
	}
	return result
}

// 解析语句列表
// 返回最后一个语句的值
func evalstatements(stmts []statement, env *env) obj {
	var result obj
	for _, statement := range stmts {
		result = eval(statement, env)
	}
	return result
}

// 解析前缀表达式
func evalprefix(operator string, right obj) obj {
	// 暂时似乎没其他的前缀操作符,暂时无法处理,返回一个错误
	if operator != "-" {
		return newerror("unknown operator: %s %s", operator, right.tYpe())
	}

	minus := func(right obj) obj {
		if right.tYpe() == INTEGER_OBJ {
			value := right.(*integerobj).value
			return &integerobj{value: -value}
		}
		if right.tYpe() == FLOAT_OBJ {
			value := right.(*floatobj).value
			return &floatobj{value: -value}
		}
		return newerror("can't use: -%s", right.tYpe())
	}

	if right.tYpe() == INTEGER_OBJ || right.tYpe() == FLOAT_OBJ {
		return minus(right)
	}

	if right.tYpe() == ARRAY_OBJ {
		rst := []obj{}
		elements := right.(*arrayobj).elements
		for _, element := range elements {
			tmp := minus(element)
			if tmp.tYpe() == ERROR_OBJ {
				return tmp
			}
			rst = append(rst, tmp)
		}
		return &arrayobj{elements: rst}
	}

	return newerror("can't use: -%s", right.tYpe())
}

// 解析中缀表达式
func evalInfixexpression(operator string, left obj, right obj) obj {
	funcmap := map[string]func(arg1 obj, arg2 obj) obj{
		"AND": func(arg1 obj, arg2 obj) obj { return bool2boolobj(istruthy(arg1) && istruthy(arg2)) },

		"OR": func(arg1 obj, arg2 obj) obj { return bool2boolobj(istruthy(arg1) || istruthy(arg2)) },

		"==": func(left obj, right obj) obj {
			if left.tYpe() == BOOLEAN_OBJ && right.tYpe() == BOOLEAN_OBJ {
				return bool2boolobj(left.(*booleanobj).value == right.(*booleanobj).value)
			}
			if left.isnumberic() && right.isnumberic() {
				return bool2boolobj(left.tofloat64().value == right.tofloat64().value)
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"!=": func(left obj, right obj) obj {
			if left.tYpe() == INTEGER_OBJ && right.tYpe() == INTEGER_OBJ {
				return bool2boolobj(left.(*integerobj).value != right.(*integerobj).value)
			}
			if left.isnumberic() && right.isnumberic() {
				return bool2boolobj(left.tofloat64().value != right.tofloat64().value)
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"+": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return &floatobj{value: left.tofloat64().value + right.tofloat64().value}
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"*": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return &floatobj{value: left.tofloat64().value * right.tofloat64().value}
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"/": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				if right.tofloat64().value == 0 {
					return &floatobj{value: 0}
				}
				return &floatobj{value: left.tofloat64().value / right.tofloat64().value}
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"-": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return &floatobj{value: left.tofloat64().value - right.tofloat64().value}
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		">": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return bool2boolobj(left.tofloat64().value > right.tofloat64().value)
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"<": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return bool2boolobj(left.tofloat64().value < right.tofloat64().value)
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		">=": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return bool2boolobj(left.tofloat64().value >= right.tofloat64().value)
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},

		"<=": func(left obj, right obj) obj {
			if left.isnumberic() && right.isnumberic() {
				return bool2boolobj(left.tofloat64().value <= right.tofloat64().value)
			}
			return newerror("can't use: %s %s %s", left.tYpe(), operator, right.tYpe())
		},
	}

	f, exist := funcmap[operator]
	if !exist {
		return newerror("unknown operator: %s %s %s", left.tYpe(), operator, right.tYpe())
	}

	if left.tYpe() != ARRAY_OBJ && right.tYpe() != ARRAY_OBJ {
		return f(left, right)
	}

	rst := []obj{}
	if left.tYpe() == ARRAY_OBJ && right.tYpe() == ARRAY_OBJ {
		if len(left.(*arrayobj).elements) != len(right.(*arrayobj).elements) {
			return newerror("len of arg0 not equal len of arg1")
		}
		for i := 0; i < len(left.(*arrayobj).elements); i++ {
			rst = append(rst, f(left.(*arrayobj).elements[i], right.(*arrayobj).elements[i]))
		}
		return &arrayobj{elements: rst}
	}

	// 如果左边是数组则扩充右边
	if left.tYpe() == ARRAY_OBJ {
		leftelements := left.(*arrayobj).elements
		for _, leftelement := range leftelements {
			rst = append(rst, f(leftelement, right))
		}
	}

	// 如果右边是数组则扩充右边
	if right.tYpe() == ARRAY_OBJ {
		rightelements := right.(*arrayobj).elements
		for _, rightelement := range rightelements {
			rst = append(rst, f(left, rightelement))
		}
	}

	return &arrayobj{elements: rst}
}

// 检查object是不是boolean
// 需要兼容其他类型
func istruthy(obj obj) bool {
	switch obj {
	case NULLOBJ:
		return false
	case TRUEOBJ:
		return true
	case FALSEOBJ:
		return false
	default:
		if obj.tYpe() == INTEGER_OBJ {
			if obj.(*integerobj).value == 0 {
				return false
			}
		}
		if obj.tYpe() == FLOAT_OBJ {
			if obj.(*floatobj).value == 0.0 {
				return false
			}
		}
		return true
	}
}

// 生成错误
func newerror(format string, a ...interface{}) *errorobj {
	return &errorobj{message: fmt.Sprintf(format, a...)}
}

// 检查是不是错误
func iserror(obj obj) bool {
	if obj != nil {
		return obj.tYpe() == ERROR_OBJ
	}
	return false
}

// 运行标识符表达式
// 从环境中取值然后执行
// 添加内置函数后还需要查看标识符是不是内置函数的函数名
func evalidentifer(node *identifier, env *env) obj {
	// 所有的全局变量或内置方法,转大写,照抄通达信,哈哈
	val := strings.ToUpper(node.value)

	// 先搜索执行环境,查看执行环境中是否保存该值
	if val, ok := env.get(val); ok {
		return val
	}

	// 再搜索内置方法
	if builtin, ok := builtins[val]; ok {
		return builtin
	}
	return newerror("idenfier not found: " + node.value)
}
