package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

/*
Cases tested:
	p.Go("(alice or bob) and carol")         // (a | b) & c
	p.Go("alice or (bob and carol)")         // a | (b & c)
*/

func Test_parenthesis_parser_case1(t *testing.T) {
	StrORStrCalled := false
	ExpANDStrCalled := false

	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrORStrCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "carol", b)

		ExpANDStrCalled = true

		return squirrel.And{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := Parser{
		StrORStr:  StrORStr,
		ExpANDStr: ExpANDStr,
	}

	p.Go("(alice or bob) and carol")
	assert.True(t, StrORStrCalled)
	assert.True(t, ExpANDStrCalled)
}

func Test_parenthesis_parser_case2(t *testing.T) {
	StrANDStrCalled := false
	ExpANDStrCalled := false

	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrANDStrCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "carol", b)

		ExpANDStrCalled = true

		return squirrel.And{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := Parser{
		StrANDStr: StrANDStr,
		ExpANDStr: ExpANDStr,
	}

	p.Go("alice or (bob and carol)")
	assert.True(t, StrANDStrCalled)
	assert.True(t, ExpANDStrCalled)
}
