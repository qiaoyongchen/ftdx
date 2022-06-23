package l

import (
	"testing"
)

func TestOBJ(t *testing.T) {
	var z1 obj = &integerobj{
		value: 0,
	}

	var z2 obj = &integerobj{
		value: 0,
	}
	t.Log(z1 == z2)
}
