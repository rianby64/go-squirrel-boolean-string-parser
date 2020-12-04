package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

/*
  Cases tested:
    p.Go("alice")                       // a
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
