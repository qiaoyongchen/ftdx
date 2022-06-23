package l

func typevalidate_isnumberic(arg obj) bool {
	return arg.isnumberic()
}

func typevalidate_isboolean(arg obj) bool {
	return arg.tYpe() == BOOLEAN_OBJ
}

func typevalidate_isinteger(arg obj) bool {
	return arg.tYpe() == INTEGER_OBJ
}

func typevalidate_isfloat(arg obj) bool {
	return arg.tYpe() == FLOAT_OBJ
}

func typevalidate_isnull(arg obj) bool {
	return arg.tYpe() == NULL_OBJ
}

func typevalidate_isarray(arg obj) bool {
	return arg.tYpe() == ARRAY_OBJ
}

func typevalidate_iserror(arg obj) bool {
	return arg.tYpe() == ERROR_OBJ
}

func typevalidate_or(fns ...func(obj) bool) func(obj) bool {
	return func(arg obj) bool {
		rst := false
		for i := 0; i < len(fns); i++ {
			fn := fns[i]
			if fn(arg) == true {
				rst = true
				break
			}
		}
		return rst
	}
}

func typevalidate_and(fns ...func(obj) bool) func(obj) bool {
	return func(arg obj) bool {
		rst := true
		for i := 0; i < len(fns); i++ {
			fn := fns[i]
			if fn(arg) == false {
				rst = false
				break
			}
		}
		return rst
	}
}

func make_2args_func(fnname string,
	arg0typevalidate func(obj) bool,
	arg1typevalidate func(obj) bool,
	calc func(arg0 obj, arg1 obj) obj) func(...obj) obj {

	fn := func(args ...obj) obj {
		selffn := make_2args_func(fnname, arg0typevalidate, arg1typevalidate, calc)
		if len(args) != 2 {
			return newerror("`%s` func len (%d) does't match 2.", fnname, len(args))
		}
		arg0 := args[0]
		arg1 := args[1]
		if arg0.tYpe() != ARRAY_OBJ && arg1.tYpe() != ARRAY_OBJ {
			if arg0typevalidate(arg0) && arg1typevalidate(arg1) {
				return calc(arg0, arg1)
			}
			return newerror("`%s` func does't match: %s(%s, %s)", fnname, fnname, arg0.inspect(), arg1.inspect())
		}

		rst := []obj{}
		if arg0.tYpe() == ARRAY_OBJ && arg1.tYpe() == ARRAY_OBJ {
			if len(arg0.(*arrayobj).elements) != len(arg1.(*arrayobj).elements) {
				return newerror("`%s` func: len of arg0 not equal len of arg1", fnname)
			}
			for i := 0; i < len(arg0.(*arrayobj).elements); i++ {
				tmprst := selffn(arg0.(*arrayobj).elements[i], arg1.(*arrayobj).elements[i])
				if tmprst.tYpe() == ERROR_OBJ {
					return tmprst
				}
				rst = append(rst, tmprst)
			}
			return &arrayobj{elements: rst}
		}
		if arg0.tYpe() == ARRAY_OBJ {
			for i := 0; i < len(arg0.(*arrayobj).elements); i++ {
				tmprst := selffn(arg0.(*arrayobj).elements[i], arg1)
				if tmprst.tYpe() == ERROR_OBJ {
					return tmprst
				}
				rst = append(rst, tmprst)
			}
		}
		if arg1.tYpe() == ARRAY_OBJ {
			for i := 0; i < len(arg1.(*arrayobj).elements); i++ {
				tmprst := selffn(arg0, arg1.(*arrayobj).elements[i])
				if tmprst.tYpe() == ERROR_OBJ {
					return tmprst
				}
				rst = append(rst, tmprst)
			}
		}
		return &arrayobj{elements: rst}
	}
	return fn
}
