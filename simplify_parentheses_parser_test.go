package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

/*
  Cases tested:
	p.Go("alice")                           // a
	p.Go("(alice)")                         // a
	p.Go("((alice))")                       // a
	p.Go("(((alice)))")                     // a
	p.Go("(((alice))) and bob")             // a & b
	p.Go("((alice) or bob)")                // a | b
	p.Go("((not (alice)) or bob)")          // !a | b
	p.Go("(alice or (not (bob)))")          // a | !b
	p.Go("((not (alice)) or (not (bob)))")  // !a | !b
	p.Go("((not (alice)) or bob)")          // !a & b
	p.Go("(alice or (not (bob)))")          // a & !b
	p.Go("((not (alice)) or (not (bob)))")  // !a & !b
*/

func Test_parser_parenthesis_case1(t *testing.T) {
	StrCalled := false
	Str := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		StrCalled = true

		r := squirrel.Expr("%s", a)
		return r
	}

	p := Parser{
		Str: Str,
	}

	_, err := p.Go("alice")
	assert.Nil(t, err)
	assert.True(t, StrCalled)
}

func Test_parser_parenthesis_case2(t *testing.T) {
	StrCalled := false
	Str := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		StrCalled = true

		r := squirrel.Expr("%s", a)
		return r
	}

	p := Parser{
		Str: Str,
	}

	_, err := p.Go("(alice)")
	assert.Nil(t, err)
	assert.True(t, StrCalled)
}

func Test_parser_parenthesis_case3(t *testing.T) {
	StrCalled := false
	Str := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		StrCalled = true

		r := squirrel.Expr("%s", a)
		return r
	}

	p := Parser{
		Str: Str,
	}

	_, err := p.Go("((alice))")
	assert.Nil(t, err)
	assert.True(t, StrCalled)
}

func Test_parser_parenthesis_case4(t *testing.T) {
	StrCalled := false
	Str := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		StrCalled = true

		r := squirrel.Expr("%s", a)
		return r
	}

	p := Parser{
		Str: Str,
	}

	_, err := p.Go("(((alice)))")
	assert.Nil(t, err)
	assert.True(t, StrCalled)
}

func Test_parser_parenthesis_case5(t *testing.T) {
	StrANDStrCalled := false
	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrANDStrCalled = true

		r := squirrel.And{
			squirrel.Expr("%s", a),
			squirrel.Expr("%s", b),
		}
		return r
	}

	p := Parser{
		StrANDStr: StrANDStr,
	}

	_, err := p.Go("(((alice))) and bob")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
}

func Test_parser_parenthesis_case6(t *testing.T) {
	StrANDStrCalled := false
	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrANDStrCalled = true

		r := squirrel.And{
			squirrel.Expr("%s", a),
			squirrel.Expr("%s", b),
		}
		return r
	}

	p := Parser{
		StrANDStr: StrANDStr,
	}

	_, err := p.Go("((alice) and bob)")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
}

func Test_parser_parenthesis_case7(t *testing.T) {
	StrORStrCalled := false
	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrORStrCalled = true

		r := squirrel.Or{
			squirrel.Expr("%s", a),
			squirrel.Expr("%s", b),
		}
		return r
	}

	p := Parser{
		StrORStr: StrORStr,
	}

	_, err := p.Go("((alice) or bob)")
	assert.Nil(t, err)
	assert.True(t, StrORStrCalled)
}

func Test_parser_parenthesis_case8(t *testing.T) {
	ExpORStrCalled := false
	NotStrCalled := false

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		assert.Equal(t, "bob", b)

		ExpORStrCalled = true

		r := squirrel.Or{
			a,
			squirrel.Expr("%s", b),
		}
		return r
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		ExpORStr: ExpORStr,
		NotStr:   NotStr,
	}

	_, err := p.Go("((not (alice)) or bob)")
	assert.Nil(t, err)
	assert.True(t, ExpORStrCalled)
	assert.True(t, NotStrCalled)
}

func Test_parser_parenthesis_case9(t *testing.T) {
	StrORExpCalled := false
	NotStrCalled := false

	StrORExp := func(a string, b squirrel.Sqlizer) squirrel.Or {
		assert.Equal(t, "alice", a)

		StrORExpCalled = true

		r := squirrel.Or{
			squirrel.Expr("%s", a),
			b,
		}
		return r
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "bob", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		StrORExp: StrORExp,
		NotStr:   NotStr,
	}

	_, err := p.Go("(alice or (not (bob)))")
	assert.Nil(t, err)
	assert.True(t, StrORExpCalled)
	assert.True(t, NotStrCalled)
}

func Test_parser_parenthesis_case10(t *testing.T) {
	ExpORExpCalled := false
	NotStrCalled := 0

	ExpORExp := func(a, b squirrel.Sqlizer) squirrel.Or {
		ExpORExpCalled = true

		r := squirrel.Or{
			squirrel.Expr("%s", a),
			b,
		}
		return r
	}

	NotStr := func(a string) squirrel.Sqlizer {
		if NotStrCalled == 0 {
			assert.Equal(t, "alice", a)
		}

		if NotStrCalled == 1 {
			assert.Equal(t, "bob", a)
		}

		NotStrCalled++

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		ExpORExp: ExpORExp,
		NotStr:   NotStr,
	}

	_, err := p.Go("((not (alice)) or (not (bob)))")
	assert.Nil(t, err)
	assert.True(t, ExpORExpCalled)
	assert.Equal(t, 2, NotStrCalled)
}

func Test_parser_parenthesis_case11(t *testing.T) {
	ExpANDStrCalled := false
	NotStrCalled := false

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "bob", b)

		ExpANDStrCalled = true

		r := squirrel.And{
			a,
			squirrel.Expr("%s", b),
		}
		return r
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		ExpANDStr: ExpANDStr,
		NotStr:    NotStr,
	}

	_, err := p.Go("((not (alice)) and bob)")
	assert.Nil(t, err)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, NotStrCalled)
}

func Test_parser_parenthesis_case12(t *testing.T) {
	StrANDExpCalled := false
	NotStrCalled := false

	StrANDExp := func(a string, b squirrel.Sqlizer) squirrel.And {
		assert.Equal(t, "alice", a)

		StrANDExpCalled = true

		r := squirrel.And{
			squirrel.Expr("%s", a),
			b,
		}
		return r
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "bob", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		StrANDExp: StrANDExp,
		NotStr:    NotStr,
	}

	_, err := p.Go("(alice and (not (bob)))")
	assert.Nil(t, err)
	assert.True(t, StrANDExpCalled)
	assert.True(t, NotStrCalled)
}

func Test_parser_parenthesis_case13(t *testing.T) {
	ExpANDExpCalled := false
	NotStrCalled := 0

	ExpANDExp := func(a, b squirrel.Sqlizer) squirrel.And {
		ExpANDExpCalled = true

		r := squirrel.And{
			squirrel.Expr("%s", a),
			b,
		}
		return r
	}

	NotStr := func(a string) squirrel.Sqlizer {
		if NotStrCalled == 0 {
			assert.Equal(t, "alice", a)
		}

		if NotStrCalled == 1 {
			assert.Equal(t, "bob", a)
		}

		NotStrCalled++

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		ExpANDExp: ExpANDExp,
		NotStr:    NotStr,
	}

	_, err := p.Go("((not (alice)) and (not (bob)))")
	assert.Nil(t, err)
	assert.True(t, ExpANDExpCalled)
	assert.Equal(t, 2, NotStrCalled)
}
