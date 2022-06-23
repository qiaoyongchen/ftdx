package l

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {
	str := "SIG := COUNT(C>HHV(C<O,10.2),10);O==C;"
	l := newlexer(str)

	for i := 0; i < 100; i++ {
		tk := l.nexttoken()
		fmt.Println(tk)
		if tk.tYpe == EOF {
			break
		}
	}
}

func TestCall1(t *testing.T) {
	str := "a(v1,v2);"
	l := newlexer(str)

	for i := 0; i < 100; i++ {
		tk := l.nexttoken()
		fmt.Println(tk)
		if tk.tYpe == EOF {
			break
		}
	}
}

func TestAssign1(t *testing.T) {
	str := "a>0.001;"
	l := newlexer(str)

	for i := 0; i < 100; i++ {
		tk := l.nexttoken()
		fmt.Println(tk)
		if tk.tYpe == EOF {
			break
		}
	}
}
