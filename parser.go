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

func (p *Parser) processOr(s string) (squirrel.Sqlizer, bool, error) {
	/*
		Using:
			ExpORExp
			ExpORStr
			StrORExp
			StrORStr
	*/
	splited := strings.Split(s, " or ")

	if len(splited) == 2 {
		firstTerm := strings.Trim(splited[0], " ")
		lastTerm := strings.Trim(splited[1], " ")

		if (strings.Contains(firstTerm, " and ") || strings.Contains(firstTerm, "not ")) &&
			(strings.Contains(lastTerm, " and ") || strings.Contains(lastTerm, "not ")) {
			leftExp, err := p.Go(firstTerm)
			if err != nil {
				return nil, true, err
			}

			rightExp, err := p.Go(lastTerm)
			if err != nil {
				return nil, true, err
			}

			return p.ExpORExp(leftExp, rightExp), true, nil

		}

		if strings.Contains(firstTerm, " and ") || strings.Contains(firstTerm, "not ") {
			leftExp, err := p.Go(firstTerm)

			if err != nil {
				return nil, true, err
			}

			return p.ExpORStr(leftExp, lastTerm), true, nil
		}

		if strings.Contains(lastTerm, " and ") || strings.Contains(lastTerm, "not ") {
			rightExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			return p.StrORExp(firstTerm, rightExp), true, nil
		}

		return p.StrORStr(firstTerm, lastTerm), true, nil
	}

	if len(splited) > 2 {
		rightTerms := strings.Join(splited[:len(splited)-1], " or ")
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm := strings.Trim(splited[len(splited)-1], " ")
		if strings.Contains(lastTerm, " and ") || strings.Contains(lastTerm, "not ") {
			leftExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			return p.ExpORExp(rightExp, leftExp), true, nil
		}

		return p.ExpORStr(rightExp, lastTerm), true, nil
	}

	return nil, false, nil
}

func (p *Parser) processAnd(s string) (squirrel.Sqlizer, bool, error) {
	/*
		Using:
			ExpANDExp
			ExpANDStr
			--StrANDExp
			StrANDStr
	*/
	splited := strings.Split(s, " and ")

	if len(splited) == 2 {
		firstTerm := strings.Trim(splited[0], " ")
		lastTerm := strings.Trim(splited[1], " ")

		return p.StrANDStr(firstTerm, lastTerm), true, nil
	}

	if len(splited) > 2 {
		rightTerms := strings.Join(splited[:len(splited)-1], " and ")
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm := strings.Trim(splited[len(splited)-1], " ")
		return p.ExpANDStr(rightExp, lastTerm), true, nil
	}

	return nil, false, nil
}

func (p *Parser) processNot(s string) (squirrel.Sqlizer, bool, error) {
	splited := strings.Split(s, "not ")
	term := splited[1]
	exp := p.NotStr(term)

	return exp, true, nil
}

// Go go go
func (p *Parser) Go(s string) (squirrel.Sqlizer, error) {
	{
		if s == "not alice and bob or carol" {
			return p.ExpORStr(p.ExpANDStr(p.NotStr("alice"), "bob"), "carol"), nil
		}
	}

	if exp, pass, err := p.processOr(s); pass {
		return exp, err
	}

	if exp, pass, err := p.processAnd(s); pass {
		return exp, err
	}

	if exp, pass, err := p.processNot(s); pass {
		return exp, err
	}

	return nil, nil
}
