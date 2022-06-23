package l

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	NULL_OBJ         = "NULL"         // null
	INTEGER_OBJ      = "INTEGER"      // 整型
	FLOAT_OBJ        = "FLOAT"        // 浮点
	BOOLEAN_OBJ      = "BOOLEAN"      // 布尔
	RETURN_VALUE_OBJ = "RETURN_VALUE" // return
	ERROR_OBJ        = "ERROR"        // error
	FUNCTION_OBJ     = "FUNCTION"     // user defined function
	STRING_OBJ       = "STRING"       // string
	BUILTIN_OBJ      = "BUILTIN"      // buildin function
	ARRAY_OBJ        = "ARRAY"        //
	HASH_OBJ         = "HASH"         //
)

type objtype string

type obj interface {
	tYpe() objtype   // 类型
	inspect() string // 检查
	isnumberic() bool
	tofloat64() *floatobj
}

// 整数类型
type integerobj struct {
	value int64
}

func (i *integerobj) inspect() string      { return fmt.Sprintf("%d", i.value) }
func (i *integerobj) tYpe() objtype        { return INTEGER_OBJ }
func (i *integerobj) isnumberic() bool     { return true }
func (i *integerobj) tofloat64() *floatobj { return &floatobj{value: float64(i.value)} }

type floatobj struct {
	value float64
}

func (f *floatobj) inspect() string      { return fmt.Sprintf("%f", f.value) }
func (f *floatobj) tYpe() objtype        { return FLOAT_OBJ }
func (f *floatobj) isnumberic() bool     { return true }
func (f *floatobj) tofloat64() *floatobj { return f }

//布尔类型ObjectType
type booleanobj struct {
	value bool
}

func (b *booleanobj) inspect() string      { return fmt.Sprintf("%t", b.value) }
func (b *booleanobj) tYpe() objtype        { return BOOLEAN_OBJ }
func (b *booleanobj) isnumberic() bool     { return false }
func (b *booleanobj) tofloat64() *floatobj { return nil }

// 空指针类型
type nullobj struct{}

func (n *nullobj) tYpe() objtype        { return NULL_OBJ }
func (n *nullobj) inspect() string      { return "null" }
func (n *nullobj) isnumberic() bool     { return false }
func (n *nullobj) tofloat64() *floatobj { return nil }

// return值(可包含任何类型的值)
type returnvalueobj struct {
	value obj
}

func (rv *returnvalueobj) tYpe() objtype        { return RETURN_VALUE_OBJ }
func (rv *returnvalueobj) Inspect() string      { return rv.value.inspect() }
func (rv *returnvalueobj) isnumberic() bool     { return rv.value.isnumberic() }
func (rv *returnvalueobj) tofloat64() *floatobj { return rv.value.tofloat64() }

// 错误类型
type errorobj struct {
	message string
}

func (e *errorobj) tYpe() objtype        { return ERROR_OBJ }
func (e *errorobj) inspect() string      { return "ERROR: " + e.message }
func (e *errorobj) isnumberic() bool     { return false }
func (e *errorobj) tofloat64() *floatobj { return nil }

// 内置函数
type builtinfnobj struct {
	fn builtinfunction
}
type builtinfunction func(args ...obj) obj

func (b *builtinfnobj) tYpe() objtype        { return BUILTIN_OBJ }
func (b *builtinfnobj) inspect() string      { return "builtin funciton" }
func (b *builtinfnobj) isnumberic() bool     { return false }
func (b *builtinfnobj) tofloat64() *floatobj { return nil }

// 数组
type arrayobj struct {
	elements []obj //包含任何类型的列表
}

func (ao *arrayobj) tYpe() objtype { return ARRAY_OBJ }
func (ao *arrayobj) inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.elements {
		elements = append(elements, e.inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (ao *arrayobj) isnumberic() bool     { return false }
func (ao *arrayobj) tofloat64() *floatobj { return nil }
