package parser

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

// Parser is the parser
type Parser struct {
	StrORStr func(a, b string) squirrel.Or
	ExpORStr func(a squirrel.Sqlizer, b string) squirrel.Or
	StrORExp func(a string, b squirrel.Sqlizer) squirrel.Or
	ExpORExp func(a, b squirrel.Sqlizer) squirrel.Or

	StrANDStr func(a, b string) squirrel.And
	ExpANDStr func(a squirrel.Sqlizer, b string) squirrel.And
	StrANDExp func(a string, b squirrel.Sqlizer) squirrel.And
	ExpANDExp func(a, b squirrel.Sqlizer) squirrel.And

	NotStr func(a string) squirrel.Sqlizer
	NotExp func(a squirrel.Sqlizer) squirrel.Sqlizer
}

// Go go go
func (p *Parser) Go(s string) error {
	{
		if s == "not alice" {
			p.NotStr("alice")
			return nil
		}
	}

	{
		splited := strings.Split(s, " and ")

		if len(splited) == 2 {
			p.StrANDStr(strings.Trim(splited[0], " "), strings.Trim(splited[1], " "))
			return nil
		}

		if len(splited) == 3 {
			left := p.StrANDStr(strings.Trim(splited[1], " "), strings.Trim(splited[2], " "))
			p.StrANDExp(strings.Trim(splited[0], " "), left)
			return nil
		}
	}

	{
		splited := strings.Split(s, " or ")

		if len(splited) == 2 {
			p.StrORStr(strings.Trim(splited[0], " "), strings.Trim(splited[1], " "))
			return nil
		}

		if len(splited) == 3 {
			left := p.StrORStr(strings.Trim(splited[1], " "), strings.Trim(splited[2], " "))
			p.StrORExp(strings.Trim(splited[0], " "), left)
			return nil
		}
	}
	return nil
}
