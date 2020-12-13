package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func Test_not_followed_by_parentheses_case1(t *testing.T) {
	parts, err := splitParentheses("alice and not (bob)")
	assert.Nil(t, err)
	assert.Equal(t, []string{"alice and ", "not (bob)"}, parts)
}

func Test_not_followed_by_parentheses_case2(t *testing.T) {
	parts, err := splitAnd("alice and not (bob)")
	assert.Nil(t, err)
	assert.Equal(t, []string{"alice", "not (bob)"}, parts)
}

func Test_not_followed_by_parentheses_case3(t *testing.T) {
	parts, err := splitParentheses("alice or not (bob)")
	assert.Nil(t, err)
	assert.Equal(t, []string{"alice or ", "not (bob)"}, parts)
}

func Test_not_followed_by_parentheses_case4(t *testing.T) {
	parts, err := splitOr("alice or not (bob)")
	assert.Nil(t, err)
	assert.Equal(t, []string{"alice", "not (bob)"}, parts)
}

func Test_not_followed_by_parentheses_case5(t *testing.T) {
	parts, err := splitParentheses("alice and not (bob or not carol)")
	assert.Nil(t, err)
	assert.Equal(t, []string{"alice and ", "not (bob or not carol)"}, parts)
}

func Test_not_followed_by_parentheses_case6(t *testing.T) {
	assert.True(t, testExpression("alice and not (bob or not carol)"))
}

func Test_not_followed_by_parentheses(t *testing.T) {
	p := New(func(s string) squirrel.Sqlizer {
		return squirrel.Expr("col = %s", s)
	})

	cases := []struct {
		input  string
		values []interface{}
		sql    string
	}{
		{
			"alice and not(bob)",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"alice and not bob",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"alice and not (bob)",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"alice and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"(alice) and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"(not alice) and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND NOT (col = %s))",
		},
		{
			"(not (alice)) and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND NOT (col = %s))",
		},
		{
			"not alice and bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND col = %s)",
		},
		{
			"not (alice) and bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND col = %s)",
		},
		{
			"(not (alice)) and bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND col = %s)",
		},
		{
			"alice or not bob",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"alice or not (bob)",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"alice or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"(alice) or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"(not alice) or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR NOT (col = %s))",
		},
		{
			"(not (alice)) or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR NOT (col = %s))",
		},
		{
			"not alice or bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR col = %s)",
		},
		{
			"not (alice) or bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR col = %s)",
		},
		{
			"(not (alice)) or bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR col = %s)",
		},
	}

	for _, curr := range cases {
		exp, err := p.Go(curr.input)
		assert.Nil(t, err)

		sql, v, err := exp.ToSql()
		assert.Nil(t, err)
		assert.Equal(t, curr.values, v)
		assert.Equal(t, curr.sql, sql)
	}
}

func Test_complex_not_followed_by_parentheses(t *testing.T) {
	p := New(func(s string) squirrel.Sqlizer {
		return squirrel.Expr("col = %s", s)
	})

	cases := []struct {
		input  string
		values []interface{}
		sql    string
	}{
		{
			"alice and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"(alice) and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s AND NOT (col = %s))",
		},
		{
			"(not alice) and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND NOT (col = %s))",
		},
		{
			"(not (alice)) and (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND NOT (col = %s))",
		},
		{
			"not alice and bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND col = %s)",
		},
		{
			"not (alice) and bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND col = %s)",
		},
		{
			"(not (alice)) and bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) AND col = %s)",
		},
		{
			"alice or not bob",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"alice or not (bob)",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"alice or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"(alice) or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(col = %s OR NOT (col = %s))",
		},
		{
			"(not alice) or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR NOT (col = %s))",
		},
		{
			"(not (alice)) or (not (bob))",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR NOT (col = %s))",
		},
		{
			"not alice or bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR col = %s)",
		},
		{
			"not (alice) or bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR col = %s)",
		},
		{
			"(not (alice)) or bob",
			[]interface{}{"alice", "bob"},
			"(NOT (col = %s) OR col = %s)",
		},
		{
			"alice and not (bob or not carol)",
			[]interface{}{"alice", "bob", "carol"},
			"(col = %s AND NOT ((col = %s OR NOT (col = %s))))",
		},
	}

	for _, curr := range cases {
		exp, err := p.Go(curr.input)
		assert.Nil(t, err)

		sql, v, err := exp.ToSql()
		assert.Nil(t, err)
		assert.Equal(t, curr.values, v)
		assert.Equal(t, curr.sql, sql)
	}
}
