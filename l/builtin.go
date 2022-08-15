package l

import (
	"math"
)

// 内置函数
var builtins = map[string]*builtinfnobj{
	"IF":        {fn: builtin_fn_if},
	"REF":       {fn: builtin_fn_ref},
	"COUNT":     {fn: builtin_fn_count},
	"MAX":       {fn: builtin_fn_max},
	"MIN":       {fn: builtin_fn_min},
	"ABS":       {fn: builtin_fn_abs},
	"EVERY":     {fn: builtin_fn_every},
	"EXIST":     {fn: builtin_fn_exist},
	"MA":        {fn: builtin_fn_ma},
	"BARSCOUNT": {fn: builtin_fn_barscount},
	"HHV":       {fn: builtin_fn_hhv},
	"LLV":       {fn: builtin_fn_llv},
	"SMA":       {fn: builtin_fn_sma},
}

func builtin_fn_if(args ...obj) obj {
	if len(args) != 3 {
		return newerror("`if` func len does't match %d", len(args))
	}
	arg0, arg1, arg2 := args[0], args[1], args[2]
	if arg0.tYpe() != ARRAY_OBJ && arg1.tYpe() != ARRAY_OBJ && arg2.tYpe() != ARRAY_OBJ {
		if istruthy(args[0]) {
			return args[1]
		}
		return args[2]
	}
	l := 0
	for _, arg := range args {
		if arg.tYpe() != ARRAY_OBJ {
			continue
		}
		if l == 0 {
			l = len(arg.(*arrayobj).elements)
			continue
		}
		if l == len(arg.(*arrayobj).elements) {
			continue
		}
		return newerror("`if` len does't match %d %d %d", len(arg0.(*arrayobj).elements), len(arg1.(*arrayobj).elements), len(arg2.(*arrayobj).elements))
	}
	if l == 0 {
		return &arrayobj{}
	}
	var getobj = func(argindex int, elementindex int) obj {
		arg := args[argindex]
		if arg.tYpe() != ARRAY_OBJ {
			return arg
		}
		return arg.(*arrayobj).elements[elementindex]
	}

	rst := &arrayobj{elements: []obj{}}
	for i := 0; i < l; i++ {
		if istruthy(getobj(0, i)) {
			rst.elements = append(rst.elements, getobj(1, i), getobj(2, i))
		}
	}
	return rst
}

func builtin_fn_ref(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`ref` func len does't match %d", len(args))
	}
	arg0 := args[0]
	arg1 := args[1]
	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`ref` func 1th argument must be array")
	}
	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`ref` func 2th argument must be integer")
	}

	arg0real := arg0.(*arrayobj).elements
	arg1real := arg1.(*integerobj).value

	rst := &arrayobj{elements: make([]obj, len(arg0real))}
	for i := 0; i < len(arg0real); i++ {
		ref_n := i - int(arg1real)
		if ref_n < 0 {
			rst.elements[i] = zerovalue(arg0real[i].tYpe())
		} else {
			rst.elements[i] = arg0real[ref_n]
		}
	}
	return rst
}

func builtin_fn_count(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`count` func len does't match %d", len(args))
	}
	arg0 := args[0]
	arg1 := args[1]
	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`count` func 1th argument must be array")
	}
	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`count` func 2th argument must be integer")
	}

	arg0real := arg0.(*arrayobj).elements
	arg1real := arg1.(*integerobj).value

	rst := &arrayobj{elements: make([]obj, len(arg0real))}
	for i := 0; i < len(arg0real); i++ {
		ref_n := i - int(arg1real)
		count := 0
		for j := ref_n; j < i; j++ {
			if j < 0 {
				continue
			}
			if istruthy(arg0real[j]) {
				count++
			}
		}
		rst.elements[i] = &integerobj{value: int64(count)}
	}
	return rst
}

func builtin_fn_every(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`every` func len does't match %d", len(args))
	}

	arg0 := args[0]
	arg1 := args[1]
	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`every` func 1th argument must be array")
	}
	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`every` func 2th argument must be integer")
	}

	arg0real := arg0.(*arrayobj).elements
	arg1real := arg1.(*integerobj).value

	rst := &arrayobj{elements: make([]obj, len(arg0real))}
	for i := 0; i < len(arg0real); i++ {
		ref_n := i - int(arg1real)
		alwaystrue := true
		for j := ref_n; j < i; j++ {
			if j < 0 || !istruthy(arg0real[j]) {
				alwaystrue = false
				break
			}
		}
		rst.elements[i] = &booleanobj{value: alwaystrue}
	}
	return rst
}

func builtin_fn_exist(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`exist` func len does't match %d", len(args))
	}

	arg0 := args[0]
	arg1 := args[1]
	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`exist` func 1th argument must be array")
	}
	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`exist` func 2th argument must be integer")
	}

	arg0real := arg0.(*arrayobj).elements
	arg1real := arg1.(*integerobj).value

	rst := &arrayobj{elements: make([]obj, len(arg0real))}
	for i := 0; i < len(arg0real); i++ {
		ref_n := i - int(arg1real)
		alwaysfalse := false
		for j := ref_n; j < i; j++ {
			if j < 0 {
				break
			}
			if istruthy(arg0real[i]) {
				alwaysfalse = true
				break
			}
		}
		rst.elements[i] = &booleanobj{value: alwaysfalse}
	}
	return rst
}

func builtin_fn_abs(args ...obj) obj {
	if len(args) != 1 {
		return newerror("`abs` func len does't match %d", len(args))
	}
	arg0 := args[0]
	if arg0.isnumberic() {
		return &floatobj{value: math.Abs(arg0.tofloat64().value)}
	}

	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`abs` func 1th argument must be array")
	}

	rst := &arrayobj{elements: make([]obj, len(arg0.(*arrayobj).elements))}
	for i, ele := range arg0.(*arrayobj).elements {
		rst.elements[i] = &floatobj{value: math.Abs(ele.tofloat64().value)}
	}
	return rst
}

func builtin_fn_max(args ...obj) obj {
	fn := make_2args_func(
		"max",
		typevalidate_or(typevalidate_isarray, typevalidate_isnumberic),
		typevalidate_or(typevalidate_isarray, typevalidate_isnumberic),
		func(arg0 obj, arg1 obj) obj {
			return &floatobj{value: math.Max(arg0.tofloat64().value, arg1.tofloat64().value)}
		})
	return fn(args...)
}

func builtin_fn_min(args ...obj) obj {
	fn := make_2args_func(
		"min",
		typevalidate_or(typevalidate_isarray, typevalidate_isnumberic),
		typevalidate_or(typevalidate_isarray, typevalidate_isnumberic),
		func(arg0 obj, arg1 obj) obj {
			return &floatobj{value: math.Min(arg0.tofloat64().value, arg1.tofloat64().value)}
		})
	return fn(args...)
}

func builtin_fn_ma(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`ma` func args length does'nt match, now %d, want %d", len(args), 2)
	}

	arg0 := args[0] // 第一个参数
	arg1 := args[1] // 第二个参数

	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`ma` func 0th argument must be array")
	}

	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`ma` func 1th argument must be integer")
	}

	rst_array_obj := arrayobj{elements: make([]obj, len(arg0.(*arrayobj).elements))}

	for i := range arg0.(*arrayobj).elements {
		start_i := i - int(arg1.(*integerobj).value)
		if start_i < 0 {
			start_i = 0
		}
		end_i := i
		var totol_i float64 = 0
		for j := start_i; j <= end_i; j++ {
			totol_i += arg0.(*arrayobj).elements[j].tofloat64().value
		}
		rst_array_obj.elements[i] = &floatobj{value: totol_i / arg1.tofloat64().value}
	}

	return &rst_array_obj
}

func builtin_fn_barscount(args ...obj) obj {
	if len(args) != 1 {
		return newerror("`barscount` func args length does'nt match, now %d, want %d", len(args), 1)
	}

	arg0 := args[0] // 第一个参数

	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`barscount` func 0th argument must be array")
	}

	return &integerobj{value: int64(len(arg0.(*arrayobj).elements))}
}

func builtin_fn_hhv(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`hhv` func args length does'nt match, now %d, want %d", len(args), 2)
	}

	arg0 := args[0] // 第一个参数
	arg1 := args[1] // 第二个参数

	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`hhv` func 0th argument must be array")
	}

	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`hhv` func 1th argument must be integer")
	}

	rst_array_obj := arrayobj{elements: make([]obj, len(arg0.(*arrayobj).elements))}

	for i := range arg0.(*arrayobj).elements {
		start_i := i - int(arg1.(*integerobj).value)
		if start_i < 0 {
			start_i = 0
		}
		end_i := i
		var max float64 = arg0.(*arrayobj).elements[start_i].tofloat64().value
		for j := start_i; j <= end_i; j++ {
			if arg0.(*arrayobj).elements[j].tofloat64().value > max {
				max = arg0.(*arrayobj).elements[j].tofloat64().value
			}
		}
		rst_array_obj.elements[i] = &floatobj{value: max}
	}

	return &rst_array_obj
}

func builtin_fn_llv(args ...obj) obj {
	if len(args) != 2 {
		return newerror("`hhv` func args length does'nt match, now %d, want %d", len(args), 2)
	}

	arg0 := args[0] // 第一个参数
	arg1 := args[1] // 第二个参数

	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`hhv` func 0th argument must be array")
	}

	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`hhv` func 1th argument must be integer")
	}

	rst_array_obj := arrayobj{elements: make([]obj, len(arg0.(*arrayobj).elements))}

	for i := range arg0.(*arrayobj).elements {
		start_i := i - int(arg1.(*integerobj).value)
		if start_i < 0 {
			start_i = 0
		}
		end_i := i
		var min float64 = arg0.(*arrayobj).elements[start_i].tofloat64().value
		for j := start_i; j <= end_i; j++ {
			if arg0.(*arrayobj).elements[j].tofloat64().value < min {
				min = arg0.(*arrayobj).elements[j].tofloat64().value
			}
		}
		rst_array_obj.elements[i] = &floatobj{value: min}
	}

	return &rst_array_obj
}

func builtin_fn_sma(args ...obj) obj {
	if len(args) != 3 {
		return newerror("`sma` func args length does'nt match, now %d, want %d", len(args), 3)
	}

	arg0 := args[0] // 第一个参数
	arg1 := args[1] // 第二个参数
	arg2 := args[2] // 第三个参数

	if arg0.tYpe() != ARRAY_OBJ {
		return newerror("`sma` func 0th argument must be array")
	}

	if arg1.tYpe() != INTEGER_OBJ {
		return newerror("`sma` func 1th argument must be integer")
	}

	if arg2.tYpe() != INTEGER_OBJ {
		return newerror("`sma` func 2th argument must be integer")
	}

	rst_array_obj := arrayobj{elements: make([]obj, len(arg0.(*arrayobj).elements))}
	rst_len := len(rst_array_obj.elements)

	if arg1.(*integerobj).value < 2 {
		for i := 0; i < rst_len; i++ {
			rst_array_obj.elements[i] = &integerobj{value: 0}
		}
		return &rst_array_obj
	}

	pre_sma_value := &floatobj{value: 0}
	for i := 0; i < rst_len; i++ {
		if i == 0 {
			rst_array_obj.elements[i] = arg0.(*arrayobj).elements[i].tofloat64()
		} else {
			arg0_i_value := arg0.(*arrayobj).elements[i]

			if !arg0_i_value.isnumberic() {
				return newerror("`sma` func 0th argument this position is not numberic: %d", i)
			}

			arg0_i_float_value := arg0.(*arrayobj).elements[i].tofloat64().value

			v := (arg0_i_float_value*arg2.tofloat64().value + (arg1.tofloat64().value-arg2.tofloat64().value)*pre_sma_value.value) / arg1.tofloat64().value
			rst_array_obj.elements[i] = &floatobj{value: v}
		}

		pre_sma_value = rst_array_obj.elements[i].(*floatobj)
	}

	return &rst_array_obj
}

// 根据类型取各类型零值
func zerovalue(t objtype) obj {
	if t == ARRAY_OBJ {
		return &arrayobj{}
	}
	if t == INTEGER_OBJ {
		return &integerobj{value: 0}
	}
	if t == FLOAT_OBJ {
		return &floatobj{value: 0}
	}
	if t == BOOLEAN_OBJ {
		return &booleanobj{value: false}
	}
	return NULLOBJ
}
