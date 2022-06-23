package l

import (
	"bytes"
	"strings"
)

// 节点
type node interface {
	tokenliteral() string
	string() string
}

// 语句
type statement interface {
	node
	statementnode()
}

// 表达式
type expression interface {
	node
	expressionnode()
}

// 程序
type program struct {
	statements []statement
}

func (p *program) tokenliteral() string {
	if len(p.statements) > 0 {
		return p.statements[0].tokenliteral()
	} else {
		return ""
	}
}

func (p *program) string() string {
	var out bytes.Buffer
	for _, s := range p.statements {
		out.WriteString(s.string())
	}
	return out.String()
}

// 标识符
// 可以看成表达式
type identifier struct {
	token token
	value string
}

func (i *identifier) expressionnode()      {}
func (i *identifier) tokenliteral() string { return i.token.literal }
func (i *identifier) string() string       { return i.value }

// 表达式语句
type expressionstatement struct {
	token      token
	expression expression
}

func (es *expressionstatement) statementnode()       {}
func (es *expressionstatement) tokenliteral() string { return es.token.literal }
func (es *expressionstatement) string() string {
	if es.expression != nil {
		return es.expression.string()
	}
	return ""
}

// 整型
// 整型可以看成表达式
type integerliteral struct {
	token token
	value int64
}

func (il *integerliteral) expressionnode()      {}
func (il *integerliteral) tokenliteral() string { return il.token.literal }
func (il *integerliteral) string() string       { return il.token.literal }

// 整型
// 整型可以看成表达式
type floatliteral struct {
	token token
	value float64
}

func (fl *floatliteral) expressionnode()      {}
func (fl *floatliteral) tokenliteral() string { return fl.token.literal }
func (fl *floatliteral) string() string       { return fl.token.literal }

// 中缀表达式
type infixexpression struct {
	token    token
	operator string
	left     expression
	right    expression
}

func (i *infixexpression) expressionnode()      {}
func (i *infixexpression) tokenliteral() string { return i.token.literal }
func (i *infixexpression) string() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.left.string())
	out.WriteString(" " + i.operator + " ")
	out.WriteString(i.right.string())
	out.WriteString(")")
	return out.String()
}

// 前缀表达式
type prefixexpression struct {
	token    token // '!' / '-'
	operator string
	right    expression
}

func (pe *prefixexpression) expressionnode()      {}
func (pe *prefixexpression) tokenliteral() string { return pe.token.literal }
func (pe *prefixexpression) string() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.operator)
	out.WriteString(pe.right.string())
	out.WriteString(")")
	return out.String()
}

// bool 表达式
type boolean struct {
	token token
	value bool
}

func (b *boolean) expressionnode()      {}
func (b *boolean) tokenliteral() string { return b.token.literal }
func (b *boolean) string() string       { return b.token.literal }

// 函数调用表达式
type callexpression struct {
	token     token        // '('
	function  expression   // 函数
	arguments []expression // 参数列表
}

func (ce *callexpression) expressionnode()      {}
func (ce *callexpression) tokenliteral() string { return ce.token.literal }
func (ce *callexpression) string() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.arguments {
		args = append(args, a.string())
	}
	out.WriteString(ce.function.string())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
