package parser

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

var (
	// ErrorParentheses defines it
	ErrorParentheses = fmt.Errorf("parentheses do not match")
	// ErrorOperators defines it
	ErrorOperators = fmt.Errorf("operator do not match")
	// ErrorExpression defines it
	ErrorExpression = fmt.Errorf("incorrect expression")
)

const (
	operatorAnd = " and "
	operatorOr  = " or "
	operatorNot = "not "
	openExp     = "("
	closeExp    = ")"
	separator   = " "
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

func (p *Parser) processOr(s string) (squirrel.Sqlizer, bool, error) {
	/*
		Using:
			ExpORExp
			ExpORStr
			StrORExp
			StrORStr
	*/

	st, err := simplify(s)
	if err != nil {
		return nil, true, err
	}

	terms, err := splitOr(st)
	if err != nil {
		return nil, true, err
	}

	if len(terms) == 2 {
		firstTerm, err := simplify(terms[0])
		if err == ErrorParentheses {
			return nil, false, err
		}

		lastTerm, err := simplify(terms[1])
		if err == ErrorParentheses {
			return nil, false, err
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
		rightTerms := strings.Join(terms[:len(terms)-1], operatorOr)
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm, err := simplify(terms[len(terms)-1])
		if err != nil {
			return nil, true, err
		}

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

	st, err := simplify(s)
	if err != nil {
		return nil, true, err
	}

	terms, err := splitAnd(st)
	if err != nil {
		return nil, true, err
	}

	if len(terms) == 2 {
		firstTerm, err := simplify(terms[0])
		if err == ErrorParentheses {
			return nil, false, err
		}

		lastTerm, err := simplify(terms[1])
		if err == ErrorParentheses {
			return nil, false, err
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
		rightTerms := strings.Join(terms[:len(terms)-1], operatorAnd)
		rightExp, err := p.Go(rightTerms)

		if err != nil {
			return nil, true, err
		}

		lastTerm, err := simplify(terms[len(terms)-1])
		if err != nil {
			return nil, true, err
		}

		lastTermContainsOperator := containsOperator(lastTerm)
		if lastTermContainsOperator {
			leftExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			return p.ExpANDExp(rightExp, leftExp), true, nil
		}

		return p.ExpANDStr(rightExp, lastTerm), true, nil
	}

	return nil, false, nil
}

func (p *Parser) processNot(s string) (squirrel.Sqlizer, bool, error) {
	st, _ := simplify(s)
	terms := strings.Split(st, operatorNot)
	if len(terms) > 1 {
		term, err := simplify(terms[1])
		if err != nil {
			return nil, true, err
		}

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
	if !testExpression(s) {
		return nil, ErrorExpression
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

	st, _ := simplify(s)
	return p.Str(st), nil
}
