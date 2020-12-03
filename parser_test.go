package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func Test_parser_one_not(t *testing.T) {

	NotStrCalled := false
	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		NotStr: NotStr,
	}

	p.Go("not alice")
	assert.True(t, NotStrCalled)
}

func Test_parser_two_ands(t *testing.T) {

	StrANDStrCalled := false
	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrANDStrCalled = true

		r := squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
		return r
	}

	p := Parser{
		StrANDStr: StrANDStr,
	}

	p.Go("alice and bob")
	assert.True(t, StrANDStrCalled)
}

func Test_parser_three_ands(t *testing.T) {

	StrANDStrCalled := false
	StrANDExpCalled := false

	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "bob", a)
		assert.Equal(t, "mark", b)

		StrANDStrCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	StrANDExp := func(a string, b squirrel.Sqlizer) squirrel.And {
		assert.Equal(t, "alice", a)

		StrANDExpCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			b,
		}
	}

	p := Parser{
		StrANDStr: StrANDStr,
		StrANDExp: StrANDExp,
	}

	p.Go("alice and bob and mark")
	assert.True(t, StrANDStrCalled)
	assert.True(t, StrANDExpCalled)
}

func Test_parser_two_ors(t *testing.T) {

	StrORStrCalled := false
	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrORStrCalled = true

		r := squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
		return r
	}

	p := Parser{
		StrORStr: StrORStr,
	}

	p.Go("alice or bob")
	assert.True(t, StrORStrCalled)
}

func Test_parser_three_ors(t *testing.T) {

	StrORStrCalled := false
	StrORExpCalled := false

	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "bob", a)
		assert.Equal(t, "mark", b)

		StrORStrCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	StrORExp := func(a string, b squirrel.Sqlizer) squirrel.Or {
		assert.Equal(t, "alice", a)

		StrORExpCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			b,
		}
	}

	p := Parser{
		StrORStr: StrORStr,
		StrORExp: StrORExp,
	}

	p.Go("alice or bob or mark")
	assert.True(t, StrORStrCalled)
	assert.True(t, StrORExpCalled)
}
