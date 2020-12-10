package parser

import "strings"

func splitParentheses(s string) ([]string, error) {
	st, err := simplify(s)
	if err != nil {
		return nil, err
	}

	terms := []string{}
	currPart := ""
	q := 0
	j := 0

	for i := 0; i < len(st); i++ {
		t := st[i : i+1]
		if t == openExp {
			if currPart != "" && currPart != operatorNot {
				terms = append(terms, currPart)
				currPart = ""
			}
			q++
		} else if t == closeExp {
			q--
		}

		if q < 0 {
			return nil, ErrorParentheses
		} else if q == 0 {
			sp := st[j : i+1]
			if len(sp) == 1 {
				currPart += t
			} else {
				terms = append(terms, currPart+sp)
				currPart = ""
			}
			j = i + 1
		}
	}

	if currPart != "" {
		terms = append(terms, currPart)
	}

	for i := 0; i < len(terms); i++ {
		term := terms[i]
		if term == operatorAnd+operatorNot {
			terms[i] = operatorAnd
			terms[i+1] = operatorNot + terms[i+1]
		}

		if term == operatorOr+operatorNot {
			terms[i] = operatorOr
			terms[i+1] = operatorNot + terms[i+1]
		}
	}

	return terms, nil
}

func splitOr(s string) ([]string, error) {
	return splitParenthesesBy(operatorOr, s)
}

func splitAnd(s string) ([]string, error) {
	return splitParenthesesBy(operatorAnd, s)
}

func splitParenthesesBy(operator, s string) ([]string, error) {
	terms, err := splitParentheses(s)
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
