package parser

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

// Parser exposes the Go
type Parser interface {
	Go(string) (squirrel.Sqlizer, error)
}

// New constructor
func New(Str func(search string) squirrel.Sqlizer) Parser {
	p := &parser2{
		Str: Str,
		StrORStr: func(a, b string) squirrel.Or {
			return squirrel.Or{Str(a), Str(b)}
		},

		ExpORStr: func(a squirrel.Sqlizer, b string) squirrel.Or {
			return squirrel.Or{a, Str(b)}
		},

		StrORExp: func(a string, b squirrel.Sqlizer) squirrel.Or {
			return squirrel.Or{Str(a), b}
		},

		ExpORExp: func(a, b squirrel.Sqlizer) squirrel.Or {
			return squirrel.Or{a, b}
		},

		StrANDStr: func(a, b string) squirrel.And {
			return squirrel.And{Str(a), Str(b)}
		},

		ExpANDStr: func(a squirrel.Sqlizer, b string) squirrel.And {
			return squirrel.And{a, Str(b)}
		},

		StrANDExp: func(a string, b squirrel.Sqlizer) squirrel.And {
			return squirrel.And{Str(a), b}
		},

		ExpANDExp: func(a, b squirrel.Sqlizer) squirrel.And {
			return squirrel.And{a, b}
		},

		NotStr: func(a string) squirrel.Sqlizer {
			s, v, _ := Str(a).ToSql()
			return squirrel.Expr(fmt.Sprintf("NOT (%s)", s), v...)
		},

		NotExp: func(a squirrel.Sqlizer) squirrel.Sqlizer {
			s, v, _ := a.ToSql()
			return squirrel.Expr(fmt.Sprintf("NOT (%s)", s), v...)
		},
	}

	return p
}
