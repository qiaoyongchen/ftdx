package l

import "testing"

func TestInt(t *testing.T) {
	l := newlexer("5;")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestEQ(t *testing.T) {
	l := newlexer("a==b;b==c;")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestAssign(t *testing.T) {
	l := newlexer("a:=b;b:=c;")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestAssignAndShow(t *testing.T) {
	l := newlexer("a:b+c;b:c;")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestCalc(t *testing.T) {
	l := newlexer("a*b + b *c;")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestCall(t *testing.T) {
	l := newlexer("a(v1,v2);")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestLPAREN(t *testing.T) {
	l := newlexer("(1+2);(3+5)")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func TestAND(t *testing.T) {
	l := newlexer("true AND true")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}

func Test11(t *testing.T) {
	l := newlexer("open > 0.001;")
	p := newparser(l)
	program := p.parseprogram()
	t.Log(p.errors)
	t.Log(len(program.statements))
	t.Log(program.statements[0])
	t.Log(program.tokenliteral())
	t.Log(program.string())
}
