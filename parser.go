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

// Error definitions
var (
	ErrorNotDefinedStrORStr  = fmt.Errorf("not defined StrORStr")
	ErrorNotDefinedExpORStr  = fmt.Errorf("not defined ExpORStr")
	ErrorNotDefinedStrORExp  = fmt.Errorf("not defined StrORExp")
	ErrorNotDefinedExpORExp  = fmt.Errorf("not defined ExpORExp")
	ErrorNotDefinedStrANDStr = fmt.Errorf("not defined StrANDStr")
	ErrorNotDefinedExpANDStr = fmt.Errorf("not defined ExpANDStr")
	ErrorNotDefinedStrANDExp = fmt.Errorf("not defined StrANDExp")
	ErrorNotDefinedExpANDExp = fmt.Errorf("not defined ExpANDExp")
	ErrorNotDefinedNotStr    = fmt.Errorf("not defined NotStr")
	ErrorNotDefinedNotExp    = fmt.Errorf("not defined NotExp")
	ErrorNotDefinedStr       = fmt.Errorf("not defined Str")
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

			if p.ExpORExp == nil {
				return nil, true, ErrorNotDefinedExpORExp
			}

			return p.ExpORExp(leftExp, rightExp), true, nil

		}

		if firstTermContainsOperator {
			leftExp, err := p.Go(firstTerm)

			if err != nil {
				return nil, true, err
			}

			if p.ExpORStr == nil {
				return nil, true, ErrorNotDefinedExpORStr
			}

			return p.ExpORStr(leftExp, lastTerm), true, nil
		}

		if lastTermContainsOperator {
			rightExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			if p.StrORExp == nil {
				return nil, true, ErrorNotDefinedStrORExp
			}

			return p.StrORExp(firstTerm, rightExp), true, nil
		}

		if p.StrORStr == nil {
			return nil, true, ErrorNotDefinedStrORStr
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

			if p.ExpORExp == nil {
				return nil, true, ErrorNotDefinedExpORExp
			}

			return p.ExpORExp(rightExp, leftExp), true, nil
		}

		if p.ExpORStr == nil {
			return nil, true, ErrorNotDefinedExpORStr
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

			if p.ExpANDExp == nil {
				return nil, true, ErrorNotDefinedExpANDExp
			}

			return p.ExpANDExp(leftExp, rightExp), true, nil

		}

		if firstTermContainsOperator {
			leftExp, err := p.Go(firstTerm)

			if err != nil {
				return nil, true, err
			}

			if p.ExpANDStr == nil {
				return nil, true, ErrorNotDefinedExpANDStr
			}

			return p.ExpANDStr(leftExp, lastTerm), true, nil
		}

		if lastTermContainsOperator {
			rightExp, err := p.Go(lastTerm)

			if err != nil {
				return nil, true, err
			}

			if p.StrANDExp == nil {
				return nil, true, ErrorNotDefinedStrANDExp
			}

			return p.StrANDExp(firstTerm, rightExp), true, nil
		}

		if p.StrANDStr == nil {
			return nil, true, ErrorNotDefinedStrANDStr
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

			if p.ExpANDExp == nil {
				return nil, true, ErrorNotDefinedExpANDExp
			}

			return p.ExpANDExp(rightExp, leftExp), true, nil
		}

		if p.ExpANDStr == nil {
			return nil, true, ErrorNotDefinedExpANDStr
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

			if p.NotExp == nil {
				return nil, true, ErrorNotDefinedNotExp
			}

			return p.NotExp(exp), true, nil
		}

		if p.NotStr == nil {
			return nil, true, ErrorNotDefinedNotStr
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
	if p.Str == nil {
		return nil, ErrorNotDefinedStr
	}

	return p.Str(st), nil
}
