package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

/*
  Cases tested:
	p.splitParentheses("(alice or bob) and carol")
	p.splitParentheses("alice or (bob and carol)")
	p.splitParentheses("(alice or bob) and (carol or dan)")
	p.splitParentheses("(alice or bob) and (carol or dan) or (elen and (frank or glenn))")
	p.splitParentheses("(alice or bob) and carol or dan or (elen and (frank or glenn))")
	p.splitParentheses("(alice or bob) and not (carol or dan) or (elen and not (frank or glenn))")
*/
func Test_splitParentheses_case1(t *testing.T) {
	p := Parser{}
	terms, err := p.splitParentheses("(alice or bob) and carol")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice or bob)",
		" and carol",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitParentheses_case2(t *testing.T) {
	p := Parser{}
	terms, err := p.splitParentheses("alice or (bob and carol)")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"alice or ",
		"(bob and carol)",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitParentheses_case3(t *testing.T) {
	p := Parser{}
	terms, err := p.splitParentheses("(alice or bob) and (carol or dan)")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice or bob)",
		" and ",
		"(carol or dan)",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitParentheses_case4(t *testing.T) {
	p := Parser{}
	terms, err := p.splitParentheses("(alice or bob) and (carol or dan) or (elen and (frank or glenn))")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice or bob)",
		" and ",
		"(carol or dan)",
		" or ",
		"(elen and (frank or glenn))",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitParentheses_case5(t *testing.T) {
	p := Parser{}
	terms, err := p.splitParentheses("(alice or bob) and carol or dan or (elen and (frank or glenn))")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice or bob)",
		" and carol or dan or ",
		"(elen and (frank or glenn))",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitParentheses_case6(t *testing.T) {
	p := Parser{}
	terms, err := p.splitParentheses("(alice or bob) and not (carol or dan) or (elen and not (frank or glenn))")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice or bob)",
		" and not ",
		"(carol or dan)",
		" or ",
		"(elen and not (frank or glenn))",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitOr_case1(t *testing.T) {
	p := Parser{}
	terms, err := p.splitOr("(alice and bob) or carol")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice and bob)",
		"carol",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitOr_case2(t *testing.T) {
	p := Parser{}
	terms, err := p.splitOr("zero or (alice and bob) or carol or (dan and elen) or (frank and glenn)")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"zero",
		"(alice and bob)",
		"carol",
		"(dan and elen)",
		"(frank and glenn)",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitAnd_case3(t *testing.T) {
	p := Parser{}
	terms, err := p.splitAnd("(alice or bob) and carol")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"(alice or bob)",
		"carol",
	}
	assert.Equal(t, expected, terms)
}

func Test_splitAnd_case4(t *testing.T) {
	p := Parser{}
	terms, err := p.splitAnd("zero and (alice or bob) and carol and (dan or elen) and (frank or glenn)")

	if err != nil {
		t.Error(err)
		return
	}

	expected := []string{
		"zero",
		"(alice or bob)",
		"carol",
		"(dan or elen)",
		"(frank or glenn)",
	}
	assert.Equal(t, expected, terms)
}

/*
Cases tested:
	p.Go("(alice or bob) and carol")                                  // (a | b) & c
	p.Go("alice or (bob and carol)")                                  // a | (b & c)
	p.Go("(alice or bob) and (((carol or dan) and frank) or glenn)")  // (a | b) & (((c | d) & f) | g)
	p.Go("not (alice or bob) and carol")                              // !(a | b) & c
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

	_, err := p.Go("(alice or bob) and carol")
	assert.Nil(t, err)
	assert.True(t, StrORStrCalled)
	assert.True(t, ExpANDStrCalled)
}

func Test_parenthesis_parser_case2(t *testing.T) {
	StrANDStrCalled := false
	StrORExpCalled := false

	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "bob", a)
		assert.Equal(t, "carol", b)

		StrANDStrCalled = true

		return squirrel.And{
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
		StrANDStr: StrANDStr,
		StrORExp:  StrORExp,
	}

	_, err := p.Go("alice or (bob and carol)")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
	assert.True(t, StrORExpCalled)
}

func Test_parenthesis_parser_case3(t *testing.T) {
	StrORStrCalled := 0
	ExpANDStrCalled := false
	ExpANDExpCalled := false
	ExpORStrCalled := false

	StrORStr := func(a, b string) squirrel.Or {
		if StrORStrCalled == 0 {
			assert.Equal(t, "alice", a)
			assert.Equal(t, "bob", b)
		}

		if StrORStrCalled == 1 {
			assert.Equal(t, "carol", a)
			assert.Equal(t, "dan", b)
		}

		StrORStrCalled++

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "frank", b)

		ExpANDStrCalled = true

		return squirrel.And{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		assert.Equal(t, "glenn", b)

		ExpORStrCalled = true

		return squirrel.Or{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpANDExp := func(a, b squirrel.Sqlizer) squirrel.And {
		ExpANDExpCalled = true

		return squirrel.And{
			a,
			b,
		}
	}

	p := Parser{
		ExpANDExp: ExpANDExp,
		StrORStr:  StrORStr,
		ExpANDStr: ExpANDStr,
		ExpORStr:  ExpORStr,
	}

	_, err := p.Go("(alice or bob) and (((carol or dan) and frank) or glenn)")
	assert.Nil(t, err)
	assert.Equal(t, 2, StrORStrCalled)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, ExpANDExpCalled)
	assert.True(t, ExpORStrCalled)
}

func Test_parenthesis_parser_case4(t *testing.T) {
	StrORStrCalled := false
	ExpANDStrCalled := false
	NotExpCalled := false

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

	NotExp := func(a squirrel.Sqlizer) squirrel.Sqlizer {
		NotExpCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := Parser{
		NotExp:    NotExp,
		StrORStr:  StrORStr,
		ExpANDStr: ExpANDStr,
	}

	_, err := p.Go("not (alice or bob) and carol")
	assert.Nil(t, err)
	assert.True(t, StrORStrCalled)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, NotExpCalled)
}
