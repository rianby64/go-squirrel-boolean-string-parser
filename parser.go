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

func isTerm(s string) bool {
	first := s[:1]
	last := s[len(s)-1:]

	if first == "(" && last == ")" {
		return true
	}

	return false
}

func containsOperator(s string) bool {
	return strings.Contains(s, " and ") ||
		strings.Contains(s, " or ") ||
		strings.Contains(s, "not ")
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

func (p *Parser) splitParentheses(s string) ([]string, error) {
	st, err := p.simplify(s)
	if err != nil {
		return nil, err
	}

	parts := []string{}
	currPart := ""
	q := 0
	j := 0

	for i := 0; i < len(st); i++ {
		t := st[i : i+1]
		if t == "(" {
			if currPart != "" && currPart != "not " {
				parts = append(parts, currPart)
				currPart = ""
			}
			q++
		} else if t == ")" {
			q--
		}

		if q < 0 {
			return nil, ErrorParentheses
		} else if q == 0 {
			sp := st[j : i+1]
			if len(sp) == 1 {
				currPart += t
			} else {
				parts = append(parts, currPart+sp)
				currPart = ""
			}
			j = i + 1
		}
	}

	if currPart != "" {
		parts = append(parts, currPart)
	}

	return parts, nil
}

func (p *Parser) splitOr(s string) ([]string, error) {
	return p.splitParenthesesBy(" or ", s)
}

func (p *Parser) splitAnd(s string) ([]string, error) {
	return p.splitParenthesesBy(" and ", s)
}

func (p *Parser) splitParenthesesBy(operator, s string) ([]string, error) {

	terms, err := p.splitParentheses(s)
	if err != nil {
		return nil, err
	}

	split := []string{}
	for _, term := range terms {
		if isTerm(term) {
			split = append(split, term)
		} else {
			parts := strings.Split(term, operator)

			for _, part := range parts {
				if part != "" {
					split = append(split, part)
				}
			}
		}
	}

	restored := strings.Join(split, operator)
	if restored != s {
		return []string{s}, nil
	}

	return split, nil
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

	st, err := p.simplify(s)
	if err != nil {
		return nil, true, err
	}

	terms, err := p.splitOr(st)
	if err != nil {
		return nil, true, err
	}

	if len(terms) == 2 {
		firstTerm, err := p.simplify(terms[0])
		if err == ErrorParentheses {
			return nil, false, nil
		}

		lastTerm, err := p.simplify(terms[1])
		if err == ErrorParentheses {
			return nil, false, nil
		}

		firstTermContainsOperator := containsOperator(firstTerm)
		lastTermContainsOperator := containsOperator(lastTerm)

		if firstTermContainsOperator && lastTermContainsOperator {
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

		if firstTermContainsOperator {
			leftExp, err := p.Go(firstTerm)

			if err != nil {
				return nil, true, err
			}

			return p.ExpORStr(leftExp, lastTerm), true, nil
		}

		if lastTermContainsOperator {
			rightExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			return p.StrORExp(firstTerm, rightExp), true, nil
		}

		return p.StrORStr(firstTerm, lastTerm), true, nil
	}

	if len(terms) > 2 {
		rightTerms := strings.Join(terms[:len(terms)-1], " or ")
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm, _ := p.simplify(terms[len(terms)-1])
		lastTermContainsOperator := containsOperator(lastTerm)
		if lastTermContainsOperator {
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

	st, err := p.simplify(s)
	if err != nil {
		return nil, true, err
	}

	terms, err := p.splitAnd(st)
	if err != nil {
		return nil, true, err
	}

	if len(terms) == 2 {
		firstTerm, err := p.simplify(terms[0])
		if err == ErrorParentheses {
			return nil, false, nil
		}

		lastTerm, err := p.simplify(terms[1])
		if err == ErrorParentheses {
			return nil, false, nil
		}

		firstTermContainsOperator := containsOperator(firstTerm)
		lastTermContainsOperator := containsOperator(lastTerm)

		if firstTermContainsOperator && lastTermContainsOperator {
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

		if firstTermContainsOperator {
			leftExp, err := p.Go(firstTerm)

			if err != nil {
				return nil, true, err
			}

			return p.ExpANDStr(leftExp, lastTerm), true, nil
		}

		if lastTermContainsOperator {
			rightExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			return p.StrANDExp(firstTerm, rightExp), true, nil
		}

		return p.StrANDStr(firstTerm, lastTerm), true, nil
	}

	if len(terms) > 2 {
		rightTerms := strings.Join(terms[:len(terms)-1], " and ")
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm, _ := p.simplify(terms[len(terms)-1])
		return p.ExpANDStr(rightExp, lastTerm), true, nil
	}

	return nil, false, nil
}

func (p *Parser) processNot(s string) (squirrel.Sqlizer, bool, error) {
	st, _ := p.simplify(s)
	terms := strings.Split(st, "not ")
	if len(terms) > 1 {
		term, _ := p.simplify(terms[1])
		if containsOperator(term) {
			exp, err := p.Go(term)

			if err != nil {
				return nil, true, err
			}

			return p.NotExp(exp), true, nil
		}

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
