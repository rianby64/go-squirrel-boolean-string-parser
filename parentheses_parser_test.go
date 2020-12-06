package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

/*
  Cases tested:
	p.Go("alice")                // a
	p.Go("(alice)")              // a
	p.Go("((alice))")            // a
	p.Go("(((alice)))")          // a
	p.Go("(((alice))) and bob")  // a & b
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

	p.Go("alice")
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

	p.Go("(alice)")
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

	p.Go("((alice))")
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

	p.Go("(((alice)))")
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

	p.Go("(((alice))) and bob")
	assert.True(t, StrANDStrCalled)
}
