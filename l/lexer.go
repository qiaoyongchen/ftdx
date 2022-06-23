package l

import (
	"strings"
)

// 词法解析器
type lexer struct {
	cur   int    // 当前读取的位置
	next  int    // 下一个要读取的位置
	ch    byte   // 当前要读取的字符
	input string // 要解析的源码
}

func newlexer(input string) *lexer {
	input = strings.ToUpper(input)
	l := &lexer{input: input}
	l.readchar()
	return l
}

// 读取一个字符
func (l *lexer) readchar() {
	if l.next >= len(l.input) {
		l.ch = byte(0)
	} else {
		l.ch = l.input[l.next]
	}
	l.cur = l.next
	l.next += 1
}

// 预读一个字符，但是不挪动指针
func (l *lexer) peekChar() byte {
	if l.next >= len(l.input) {
		return byte(0)
	} else {
		return l.input[l.next]
	}
}

func (l *lexer) nexttoken() token {
	l.skipwhitespace()
	var tk token
	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readchar()
			tk = token{tYpe: ASSIGN, literal: string(ch) + string(l.ch)}
		} else {
			tk = token{tYpe: COLON, literal: string(l.ch)}
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readchar()
			tk = token{tYpe: EQ, literal: string(ch) + string(l.ch)}
		} else {
			tk = token{tYpe: ILLEGAL, literal: string(l.ch)} // 错误
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readchar()
			tk = token{tYpe: NOT_EQ, literal: string(ch) + string(l.ch)}
		} else {
			tk = token{tYpe: ILLEGAL, literal: string(l.ch)} // 错误
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readchar()
			tk = token{tYpe: GTEQ, literal: string(ch) + string(l.ch)}
		} else {
			tk = token{tYpe: GT, literal: string(l.ch)}
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readchar()
			tk = token{tYpe: LTEQ, literal: string(ch) + string(l.ch)}
		} else {
			tk = token{tYpe: LT, literal: string(l.ch)}
		}
	case '+':
		tk = token{tYpe: PLUS, literal: string(l.ch)}
	case '-':
		tk = token{tYpe: MINUS, literal: string(l.ch)}
	case '*':
		tk = token{tYpe: ASTERISK, literal: string(l.ch)}
	case '/':
		tk = token{tYpe: SLASH, literal: string(l.ch)}
	case ';':
		tk = token{tYpe: SEMICOLON, literal: string(l.ch)}
	case ',':
		tk = token{tYpe: COMMA, literal: string(l.ch)}
	case '(':
		tk = token{tYpe: LPAREN, literal: string(l.ch)}
	case ')':
		tk = token{tYpe: RPAREN, literal: string(l.ch)}
	case byte(0):
		tk = token{tYpe: EOF, literal: ""}
	default:
		if isdigit(l.ch) {
			literal_ := l.readnumber()
			if l.ch != '.' {
				tk = token{tYpe: INT, literal: literal_}
				return tk
			} else {
				l.readchar()
				if isdigit(l.ch) {
					tk = token{tYpe: FLOAT, literal: literal_ + "." + l.readnumber()}
					return tk
				} else {
					tk = token{tYpe: ILLEGAL, literal: ""}
					return tk
				}
			}
		} else if isletter(l.ch) {
			letter := l.readidentifier()
			// 处理特殊的标识符 AND
			if letter == AND {
				tk = token{tYpe: AND, literal: letter}
				return tk
			}

			// 处理特殊的标识符 OR
			if letter == OR {
				tk = token{tYpe: OR, literal: letter}
				return tk
			}

			// 其他标识符
			tk = token{tYpe: IDENT, literal: letter}
			return tk
		} else {
			tk = token{tYpe: IDENT, literal: l.readidentifier()}
		}
	}
	l.readchar()
	return tk
}

// 解析标识符
func (l *lexer) readidentifier() string {
	position := l.cur
	for isletter(l.ch) {
		l.readchar()
	}
	return string(l.input[position:l.cur])
}

// 跳过空白字符
func (l *lexer) skipwhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readchar()
	}
}

// 读取数字
func (l *lexer) readnumber() string {
	position := l.cur
	for isdigit(l.ch) {
		l.readchar()
	}
	return string(l.input[position:l.cur])
}

// 检查是否为字母
func isletter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || '0' <= ch && ch <= '9'
}

// 是否为数字
func isdigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
