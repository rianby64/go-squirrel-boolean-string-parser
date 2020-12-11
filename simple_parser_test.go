package parser

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

/*
  Cases tested:
    p.Go("not alice")                       // !a
	p.Go("alice and bob")                   // a & b
	p.Go("alice and bob and carol")         // (a & b) & c
	p.Go("alice and bob and carol and dan") // ((a & b) & c) & d
	p.Go("alice or bob")                    // a | b
	p.Go("alice or bob or carol")           // (a | b) | c
	p.Go("alice or bob or carol or dan")    // ((a | b) | c) | d
	p.Go("alice and bob or carol or dan")   // (a & b) | (c & d)
	p.Go("alice or bob or carol and dan")   // (a | b) | (c & d)
	p.Go("alice or bob and carol and dan")  // a | ((b & c) & d)
	p.Go("alice or bob and carol or dan")   // (a | (b & c)) | d
	p.Go("alice and bob and carol or dan")  // ((a & b) & c) | d
	p.Go("alice or bob and carol")          // a | (b & c)
	p.Go("alice and bob or carol and dan")  // (a & b) | (c & d)
	p.Go("not alice and bob or carol")      // (!a & b) | c
	p.Go("not alice and not bob")           // !a & !b
	p.Go("not alice and bob")               // !a & b
	p.Go("alice and not bob")               // a & !b
*/

func Test_parser_one_not(t *testing.T) {
	NotStrCalled := false
	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := parser2{
		NotStr: NotStr,
	}

	_, err := p.Go("not alice")
	assert.Nil(t, err)
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

	p := parser2{
		StrANDStr: StrANDStr,
	}

	_, err := p.Go("alice and bob")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
}

func Test_parser_three_ands(t *testing.T) {
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

	p := parser2{
		StrANDStr: StrANDStr,
		ExpANDStr: ExpANDStr,
	}

	_, err := p.Go("alice and bob and carol")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
	assert.True(t, ExpANDStrCalled)
}

func Test_parser_four_ands(t *testing.T) {
	StrANDStrCalled := false
	ExpANDStrCalled := 0

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
		if ExpANDStrCalled == 0 {
			assert.Equal(t, "carol", b)
		}

		if ExpANDStrCalled == 1 {
			assert.Equal(t, "dan", b)
		}

		ExpANDStrCalled++

		return squirrel.And{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := parser2{
		StrANDStr: StrANDStr,
		ExpANDStr: ExpANDStr,
	}

	_, err := p.Go("alice and bob and carol and dan")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
	assert.Equal(t, 2, ExpANDStrCalled)
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

	p := parser2{
		StrORStr: StrORStr,
	}

	_, err := p.Go("alice or bob")
	assert.Nil(t, err)
	assert.True(t, StrORStrCalled)
}

func Test_parser_three_ors(t *testing.T) {
	StrORStrCalled := false
	ExpORStrCalled := false

	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrORStrCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		assert.Equal(t, "carol", b)

		ExpORStrCalled = true

		return squirrel.Or{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := parser2{
		StrORStr: StrORStr,
		ExpORStr: ExpORStr,
	}

	_, err := p.Go("alice or bob or carol")
	assert.Nil(t, err)
	assert.True(t, StrORStrCalled)
	assert.True(t, ExpORStrCalled)
}

func Test_parser_four_ors(t *testing.T) {
	StrORStrCalled := false
	ExpORStrCalled := 0

	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrORStrCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		if ExpORStrCalled == 0 {
			assert.Equal(t, "carol", b)
		}

		if ExpORStrCalled == 1 {
			assert.Equal(t, "dan", b)
		}

		ExpORStrCalled++

		return squirrel.Or{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := parser2{
		StrORStr: StrORStr,
		ExpORStr: ExpORStr,
	}

	_, err := p.Go("alice or bob or carol or dan")
	assert.Nil(t, err)
	assert.True(t, StrORStrCalled)
	assert.Equal(t, 2, ExpORStrCalled)
}

func Test_parser_five_terms_case1(t *testing.T) {
	StrANDStrCalled := false
	ExpORStrCalled := 0

	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrANDStrCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		if ExpORStrCalled == 0 {
			assert.Equal(t, "carol", b)
		}

		if ExpORStrCalled == 1 {
			assert.Equal(t, "dan", b)
		}

		ExpORStrCalled++

		return squirrel.Or{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := parser2{
		StrANDStr: StrANDStr,
		ExpORStr:  ExpORStr,
	}

	_, err := p.Go("alice and bob or carol or dan")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
	assert.Equal(t, 2, ExpORStrCalled)
}

func Test_parser_five_terms_case2(t *testing.T) {
	StrANDStrCalled := false
	StrORStrCalled := false
	ExpORExpCalled := false

	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "carol", a)
		assert.Equal(t, "dan", b)

		StrANDStrCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	StrORStr := func(a, b string) squirrel.Or {
		assert.Equal(t, "alice", a)
		assert.Equal(t, "bob", b)

		StrORStrCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORExp := func(a, b squirrel.Sqlizer) squirrel.Or {
		ExpORExpCalled = true

		return squirrel.Or{
			a,
			b,
		}
	}

	p := parser2{
		StrORStr:  StrORStr,
		StrANDStr: StrANDStr,
		ExpORExp:  ExpORExp,
	}

	_, err := p.Go("alice or bob or carol and dan")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
	assert.True(t, StrORStrCalled)
	assert.True(t, ExpORExpCalled)
}

func Test_parser_five_terms_case3(t *testing.T) {
	StrANDStrCalled := false
	ExpANDStrCalled := false
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

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "dan", b)

		ExpANDStrCalled = true

		return squirrel.And{
			a,
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

	p := parser2{
		StrANDStr: StrANDStr,
		ExpANDStr: ExpANDStr,
		StrORExp:  StrORExp,
	}

	_, err := p.Go("alice or bob and carol and dan")
	assert.Nil(t, err)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, StrANDStrCalled)
	assert.True(t, StrORExpCalled)
}

func Test_parser_five_terms_case4(t *testing.T) {
	StrANDStrCalled := false
	ExpORStrCalled := false
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

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		assert.Equal(t, "dan", b)

		ExpORStrCalled = true

		return squirrel.Or{
			a,
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

	p := parser2{
		StrANDStr: StrANDStr,
		ExpORStr:  ExpORStr,
		StrORExp:  StrORExp,
	}

	// alice or bob and carol or dan
	// alice or (bob and carol or dan)
	// alice or ((bob and carol) or dan)
	_, err := p.Go("alice or bob and carol or dan")
	assert.Nil(t, err)
	assert.True(t, ExpORStrCalled)
	assert.True(t, StrANDStrCalled)
	assert.True(t, StrORExpCalled)
}

func Test_parser_case1(t *testing.T) {
	StrANDStrCalled := false
	ExpANDStrCalled := false
	ExpORStrCalled := false

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

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		assert.Equal(t, "dan", b)

		ExpORStrCalled = true

		return squirrel.Or{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := parser2{
		StrANDStr: StrANDStr,
		ExpANDStr: ExpANDStr,
		ExpORStr:  ExpORStr,
	}

	_, err := p.Go("alice and bob and carol or dan")
	assert.Nil(t, err)
	assert.True(t, StrANDStrCalled)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, ExpORStrCalled)
}

func Test_parser_case2(t *testing.T) {
	StrORExpCalled := false
	StrANDStrCalled := false

	StrORExp := func(a string, b squirrel.Sqlizer) squirrel.Or {
		assert.Equal(t, "alice", a)

		StrORExpCalled = true

		return squirrel.Or{
			squirrel.Expr("col = '%s'", a),
			b,
		}
	}

	StrANDStr := func(a, b string) squirrel.And {
		assert.Equal(t, "bob", a)
		assert.Equal(t, "carol", b)

		StrANDStrCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	p := parser2{
		StrORExp:  StrORExp,
		StrANDStr: StrANDStr,
	}

	_, err := p.Go("alice or bob and carol")
	assert.Nil(t, err)
	assert.True(t, StrORExpCalled)
	assert.True(t, StrANDStrCalled)
}

func Test_parser_case3(t *testing.T) {
	StrANDStrCalled := 0
	ExpORExpCalled := false

	StrANDStr := func(a, b string) squirrel.And {
		if StrANDStrCalled == 0 {
			assert.Equal(t, "alice", a)
			assert.Equal(t, "bob", b)
		}

		if StrANDStrCalled == 1 {
			assert.Equal(t, "carol", a)
			assert.Equal(t, "dan", b)
		}

		StrANDStrCalled++

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORExp := func(a, b squirrel.Sqlizer) squirrel.Or {
		ExpORExpCalled = true

		return squirrel.Or{
			a,
			b,
		}
	}

	p := parser2{
		StrANDStr: StrANDStr,
		ExpORExp:  ExpORExp,
	}

	_, err := p.Go("alice and bob or carol and dan")
	assert.Nil(t, err)
	assert.Equal(t, 2, StrANDStrCalled)
	assert.True(t, ExpORExpCalled)
}

func Test_parser_case4(t *testing.T) {
	ExpANDStrCalled := false
	ExpORStrCalled := false
	NotStrCalled := false

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "bob", b)

		ExpANDStrCalled = true

		return squirrel.And{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	ExpORStr := func(a squirrel.Sqlizer, b string) squirrel.Or {
		assert.Equal(t, "carol", b)

		ExpORStrCalled = true

		return squirrel.Or{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := parser2{
		NotStr:    NotStr,
		ExpANDStr: ExpANDStr,
		ExpORStr:  ExpORStr,
	}

	_, err := p.Go("not alice and bob or carol")
	assert.Nil(t, err)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, ExpORStrCalled)
	assert.True(t, NotStrCalled)
}

func Test_parser_case5(t *testing.T) {
	ExpANDExpCalled := false
	NotStrCalled := 0

	ExpANDExp := func(a, b squirrel.Sqlizer) squirrel.And {
		ExpANDExpCalled = true

		return squirrel.And{
			a,
			b,
		}
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

	p := parser2{
		NotStr:    NotStr,
		ExpANDExp: ExpANDExp,
	}

	_, err := p.Go("not alice and not bob")
	assert.Nil(t, err)
	assert.True(t, ExpANDExpCalled)
	assert.Equal(t, 2, NotStrCalled)
}

func Test_parser_case6(t *testing.T) {
	ExpANDStrCalled := false
	NotStrCalled := false

	ExpANDStr := func(a squirrel.Sqlizer, b string) squirrel.And {
		assert.Equal(t, "bob", b)

		ExpANDStrCalled = true

		return squirrel.And{
			a,
			squirrel.Expr("col = '%s'", b),
		}
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "alice", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := parser2{
		NotStr:    NotStr,
		ExpANDStr: ExpANDStr,
	}

	_, err := p.Go("not alice and bob")
	assert.Nil(t, err)
	assert.True(t, ExpANDStrCalled)
	assert.True(t, NotStrCalled)
}

func Test_parser_case7(t *testing.T) {
	StrANDExpCalled := false
	NotStrCalled := false

	StrANDExp := func(a string, b squirrel.Sqlizer) squirrel.And {
		assert.Equal(t, "alice", a)

		StrANDExpCalled = true

		return squirrel.And{
			squirrel.Expr("col = '%s'", a),
			b,
		}
	}

	NotStr := func(a string) squirrel.Sqlizer {
		assert.Equal(t, "bob", a)

		NotStrCalled = true

		r := squirrel.Expr("col <> '%s'", a)
		return r
	}

	p := parser2{
		NotStr:    NotStr,
		StrANDExp: StrANDExp,
	}

	_, err := p.Go("alice and not bob")
	assert.Nil(t, err)
	assert.True(t, StrANDExpCalled)
	assert.True(t, NotStrCalled)
}
