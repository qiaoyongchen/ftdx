package l

type tokentype = string

const (
	ILLEGAL tokentype = "ILLEGAL" // 错误标志
	EOF     tokentype = "EOF"     // 结束标志
	IDENT   tokentype = "IDENT"   // 标识符
	INT     tokentype = "INT"     // int 类型
	FLOAT   tokentype = "FLOAT"   // 小数
	TRUE    tokentype = "TRUE"    //
	FALSE   tokentype = "FALSE"   //

	// 操作符
	ASSIGN   tokentype = ":="  // 赋值
	PLUS     tokentype = "+"   // 加
	MINUS    tokentype = "-"   // 减
	ASTERISK tokentype = "*"   // 乘
	SLASH    tokentype = "/"   // 除
	EQ       tokentype = "=="  // 等于
	NOT_EQ   tokentype = "!="  // 不等于
	GT       tokentype = ">"   // 大于
	LT       tokentype = "<"   // 小于
	GTEQ     tokentype = ">="  // 大于等于
	LTEQ     tokentype = "<="  // 小于等于
	AND      tokentype = "AND" // AND
	OR       tokentype = "OR"  // OR

	// 定界符
	COMMA     tokentype = ","
	SEMICOLON tokentype = ";"
	COLON     tokentype = ":" // 赋值并取值
	LPAREN    tokentype = "("
	RPAREN    tokentype = ")"
)

type token struct {
	tYpe    tokentype // token 类型
	literal string    // 字面值
}
