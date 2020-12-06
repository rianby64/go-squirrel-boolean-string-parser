package parser

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

var (
	// ErrorParentheses defines it
	ErrorParentheses = fmt.Errorf("parenthesis error at it")
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

	Str func(a string) squirrel.Sqlizer
}

func (p *Parser) testParentheses(s string) bool {
	q := 0
	for i := 0; i < len(s); i++ {
		t := s[i : i+1]
		if t == "(" {
			q++
		} else if t == ")" {
			q--
		}

		if q < 0 {
			return false
		}
	}

	return q == 0
}

func (p *Parser) simplify(s string) (string, error) {
	st := strings.Trim(s, " ")
	if !p.testParentheses(st) {
		return "", ErrorParentheses
	}

	l := len(st) - 1
	first := st[:1]
	last := st[l:]

	if first == "(" && last == ")" {
		middle := st[1:l]
		r, err := p.simplify(middle)
		if err == ErrorParentheses {
			return st, nil
		}

		return r, err
	}

	return st, nil
}

func (p *Parser) processOr(s string) (squirrel.Sqlizer, bool, error) {
	/*
		Using:
			ExpORExp
			ExpORStr
			StrORExp
			StrORStr
	*/
	st, _ := p.simplify(s)
	splited := strings.Split(st, " or ")

	if len(splited) == 2 {
		firstTerm, _ := p.simplify(splited[0])
		lastTerm, _ := p.simplify(splited[1])

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

		lastTerm, _ := p.simplify(splited[len(splited)-1])
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
			StrANDExp
			StrANDStr
	*/
	st, _ := p.simplify(s)
	splited := strings.Split(st, " and ")

	if len(splited) == 2 {
		firstTerm, _ := p.simplify(splited[0])
		lastTerm, _ := p.simplify(splited[1])

		if (strings.Contains(firstTerm, " or ") || strings.Contains(firstTerm, "not ")) &&
			(strings.Contains(lastTerm, " or ") || strings.Contains(lastTerm, "not ")) {
			leftExp, err := p.Go(firstTerm)
			if err != nil {
				return nil, true, err
			}

			rightExp, err := p.Go(lastTerm)
			if err != nil {
				return nil, true, err
			}

			return p.ExpANDExp(leftExp, rightExp), true, nil

		}

		if strings.Contains(firstTerm, " or ") || strings.Contains(firstTerm, "not ") {
			leftExp, err := p.Go(firstTerm)

			if err != nil {
				return nil, true, err
			}

			return p.ExpANDStr(leftExp, lastTerm), true, nil
		}

		if strings.Contains(lastTerm, " or ") || strings.Contains(lastTerm, "not ") {
			rightExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			return p.StrANDExp(firstTerm, rightExp), true, nil
		}

		return p.StrANDStr(firstTerm, lastTerm), true, nil
	}

	if len(splited) > 2 {
		rightTerms := strings.Join(splited[:len(splited)-1], " and ")
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm, _ := p.simplify(splited[len(splited)-1])
		return p.ExpANDStr(rightExp, lastTerm), true, nil
	}

	return nil, false, nil
}

func (p *Parser) processNot(s string) (squirrel.Sqlizer, bool, error) {
	st, _ := p.simplify(s)
	splited := strings.Split(st, "not ")
	if len(splited) > 1 {
		term, _ := p.simplify(splited[1])
		exp := p.NotStr(term)

		return exp, true, nil
	}

	return nil, false, nil
}

// Go go go
func (p *Parser) Go(s string) (squirrel.Sqlizer, error) {
	if exp, pass, err := p.processOr(s); pass {
		return exp, err
	}

	if exp, pass, err := p.processAnd(s); pass {
		return exp, err
	}

	if exp, pass, err := p.processNot(s); pass {
		return exp, err
	}

	st, _ := p.simplify(s)
	return p.Str(st), nil
}
