package l

import (
	"fmt"
	"strconv"
)

const (
	_                  int = iota
	LOWEST                 // 最低
	EQUALS                 // ==
	ORP                    // or
	ANDP                   // and
	LESSGREATER            // > <
	PLUS_AND_MINUS         // + -
	ASTERISK_AND_SLASH     // * /
	PREFIX                 // 前缀操作符
	CALL                   // 方法调用 ...
)

// 各token 对应的优先级
var precedences = map[tokentype]int{
	EQ:       EQUALS,
	NOT_EQ:   EQUALS,
	ASSIGN:   EQUALS,
	COLON:    EQUALS,
	OR:       ORP,
	AND:      ANDP,
	GT:       LESSGREATER,
	GTEQ:     LESSGREATER,
	LT:       LESSGREATER,
	LTEQ:     LESSGREATER,
	PLUS:     PLUS_AND_MINUS,
	MINUS:    PLUS_AND_MINUS,
	ASTERISK: ASTERISK_AND_SLASH,
	SLASH:    ASTERISK_AND_SLASH,
	LPAREN:   CALL,
}

type (
	prefixParseFn func() expression           // 前缀表达式(!, -)
	infixParseFn  func(expression) expression // 中缀表达式(+,-,*,/...)
)

// 解析器
// 词法 解析为 语法树
type parser struct {
	l      *lexer
	errors []string

	cur  token // 当前token
	next token // 下一个token

	prefixfns map[tokentype]prefixParseFn // 前缀表达式
	infixfns  map[tokentype]infixParseFn  // 中缀表达式
}

// 注册前缀表达式
func (p *parser) register_prefix(tkt tokentype, fn prefixParseFn) {
	p.prefixfns[tkt] = fn
}

// 注册中缀表达式
func (p *parser) register_infix(tkt tokentype, fn infixParseFn) {
	p.infixfns[tkt] = fn
}

func newparser(l *lexer) *parser {
	p := &parser{
		l:      l,
		errors: []string{},
	}
	// 注册前缀表达式的解析函数
	p.prefixfns = make(map[tokentype]prefixParseFn)
	p.register_prefix(IDENT, p.parseidentifier)         // 标识符
	p.register_prefix(INT, p.parseIntegerLiteral)       // 整型
	p.register_prefix(FLOAT, p.parseFloatLiteral)       // 浮点类型
	p.register_prefix(MINUS, p.parseprefixexpression)   // -(取负)
	p.register_prefix(TRUE, p.parseboolean)             // true
	p.register_prefix(FALSE, p.parseboolean)            // false
	p.register_prefix(LPAREN, p.parseGroupedExpression) // ( : (a+b)模式
	// 注册中缀表达式的解析函数
	p.infixfns = make(map[tokentype]infixParseFn)
	p.register_infix(ASSIGN, p.parseinfixexpression)   //':=' 赋值
	p.register_infix(COLON, p.parseinfixexpression)    //':' 赋值并显示
	p.register_infix(PLUS, p.parseinfixexpression)     //'+'
	p.register_infix(MINUS, p.parseinfixexpression)    //'-'(减)
	p.register_infix(SLASH, p.parseinfixexpression)    //'/'(除)
	p.register_infix(ASTERISK, p.parseinfixexpression) //'*'
	p.register_infix(EQ, p.parseinfixexpression)       //'=='
	p.register_infix(NOT_EQ, p.parseinfixexpression)   //'!='
	p.register_infix(LT, p.parseinfixexpression)       //'<'
	p.register_infix(LTEQ, p.parseinfixexpression)     //'<='
	p.register_infix(GT, p.parseinfixexpression)       //'>'
	p.register_infix(GTEQ, p.parseinfixexpression)     //'>='
	p.register_infix(LPAREN, p.parsecallexpression)    //'(' : a(v1,v2) 模式
	p.register_infix(AND, p.parseinfixexpression)      // AND
	p.register_infix(OR, p.parseinfixexpression)       // OR

	p.nexttoken()
	p.nexttoken()
	return p
}

func (p *parser) nexttoken() {
	p.cur = p.next
	p.next = p.l.nexttoken()
}

// 解析语法树入口
func (p *parser) parseprogram() *program {
	program := &program{}
	program.statements = []statement{}
	for !p.curis(EOF) {
		stmt := p.parsestatement()
		if stmt != nil {
			program.statements = append(program.statements, stmt)
		}
		p.nexttoken()
	}
	return program
}

// 返回错误
func (p *parser) Errors() []string {
	return p.errors
}

// 预读失败，一般为语法错误
func (p *parser) peekerror(t tokentype) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.next.tYpe)

	p.errors = append(p.errors, msg)
}

// 解析标识符
func (p *parser) parseidentifier() expression {
	return &identifier{token: p.cur, value: p.cur.literal}
}

// 检查语句的类型
// 没有其他类型的语句，就是解析表达式语句
func (p *parser) parsestatement() statement { return p.parseexpressionstatement() }

// 解析表达式类型语句
func (p *parser) parseexpressionstatement() *expressionstatement {
	stmt := &expressionstatement{token: p.cur}
	// 以最低优先级解析表达式
	stmt.expression = p.parseexpression(LOWEST)
	if p.nextis(SEMICOLON) {
		p.nexttoken()
	}
	return stmt
}

// 解析表达式
func (p *parser) parseexpression(precedence int) expression {
	prefix := p.prefixfns[p.cur.tYpe]
	if prefix == nil {
		p.noprefixparsefnerror(p.cur.tYpe)
		return nil
	}

	leftexp := prefix()
	for !p.nextis(SEMICOLON) && precedence < p.peekprecedence() {
		infix := p.infixfns[p.next.tYpe]
		if infix == nil {
			return leftexp
		}
		p.nexttoken()
		leftexp = infix(leftexp)
	}
	return leftexp
}

// 解析int类型字面量
func (p *parser) parseIntegerLiteral() expression {
	lit := &integerliteral{token: p.cur}
	value, err := strconv.ParseInt(p.cur.literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.cur.literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.value = value
	return lit
}

func (p *parser) parseFloatLiteral() expression {
	lit := &floatliteral{token: p.cur}
	value, err := strconv.ParseFloat(p.cur.literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.cur.literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.value = value
	return lit
}

// 解析中缀类型表达式
func (p *parser) parseinfixexpression(left expression) expression {
	expression := &infixexpression{
		token:    p.cur,
		operator: p.cur.literal,
		left:     left,
	}
	precedence := p.curprecedence()
	p.nexttoken()
	if expression.operator == "+" {
		expression.right = p.parseexpression(precedence - 1)
	} else {
		expression.right = p.parseexpression(precedence)
	}
	return expression
}

// 解析前缀类型表达式
func (p *parser) parseprefixexpression() expression {
	expression := &prefixexpression{
		token:    p.cur,
		operator: p.cur.literal,
	}
	p.nexttoken()
	// 带入PREFIX的优先级解析后面的表达式
	expression.right = p.parseexpression(PREFIX)
	return expression
}

// 下一个token的优先级
func (p *parser) peekprecedence() int {
	if p, ok := precedences[p.next.tYpe]; ok {
		return p
	}
	return LOWEST
}

// 当前token的优先级
func (p *parser) curprecedence() int {
	if p, ok := precedences[p.cur.tYpe]; ok {
		return p
	}
	return LOWEST
}

func (p *parser) noprefixparsefnerror(t tokentype) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 检查当前token的类型是否匹配
func (p *parser) curis(t tokentype) bool {
	return p.cur.tYpe == t
}

// 检查下一个token的类型是否匹配
func (p *parser) nextis(t tokentype) bool {
	return p.next.tYpe == t
}

func (p *parser) expectpeek(t tokentype) bool {
	if p.nextis(t) {
		p.nexttoken()
		return true
	} else {
		p.peekerror(t)
		return false
	}
}

// 检查 true / false 表达式
func (p *parser) parseboolean() expression {
	return &boolean{token: p.cur, value: p.curis(TRUE)}
}

// 检查 '( xxx )' 类型表达式
func (p *parser) parseGroupedExpression() expression {
	p.nexttoken()
	exp := p.parseexpression(LOWEST)
	if !p.expectpeek(RPAREN) {
		return nil
	}
	return exp
}

// 解析函数调用
func (p *parser) parsecallexpression(function expression) expression {
	// 函数调用标识符 '('
	exp := &callexpression{token: p.cur, function: function}

	// 解析函数调用参数
	// 函数参数为表达式列表
	// 例如: add(1+2, 3+4);
	exp.arguments = p.parsecallarguments()
	return exp
}

// 解析函数调用参数
func (p *parser) parsecallarguments() []expression {
	// 参数列表就是表达式列表
	args := []expression{}

	// 期望')'进行结束
	if p.nextis(RPAREN) {
		p.nexttoken()
		return args
	}

	p.nexttoken()
	// 解析调用参数
	args = append(args, p.parseexpression(LOWEST))

	for p.nextis(COMMA) {
		p.nexttoken()
		p.nexttoken()
		args = append(args, p.parseexpression(LOWEST))
	}

	if !p.expectpeek(RPAREN) {
		return nil
	}

	return args
}
