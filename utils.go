package parser

import "strings"

func isTerm(s string) bool {
	l := len(s)
	first := s[:1]
	last := s[l-1:]

	if l > 4 {
		firstPart := s[:4]
		if firstPart == operatorNot {
			first = s[4:5]
		}
	}

	if first == openExp && last == closeExp {
		return true
	}

	return false
}

func containsOperator(s string) bool {
	return strings.Contains(s, operatorAnd) ||
		strings.Contains(s, operatorOr) ||
		strings.Contains(s, operatorNot)
}

func testParentheses(s string) bool {
	q := 0
	for i := 0; i < len(s); i++ {
		t := s[i : i+1]
		if t == openExp {
			q++
		} else if t == closeExp {
			q--
		}

		if q < 0 {
			return false
		}
	}

	return q == 0
}

func simplify(s string) (string, error) {
	st := strings.Trim(s, separator)
	if st == "" {
		return "", ErrorExpression
	}

	st = strings.ReplaceAll(st, "not(", "not (")

	if !testParentheses(st) {
		return "", ErrorParentheses
	}

	l := len(st) - 1
	first := st[:1]
	last := st[l:]

	if first == openExp && last == closeExp {
		middle := st[1:l]
		r, err := simplify(middle)
		if err == ErrorParentheses {
			return st, nil
		}

		return r, err
	}

	return st, nil
}

func testExpression(s string) bool {
	if s == "" {
		return false
	}

	st, err := simplify(s)
	if err != nil {
		return false
	}

	var parts []string

	if !containsOperator(st) {
		l := len(st)
		if l >= 3 {
			wrongStart := st[:3]
			wrongEnd := st[l-3:]

			if wrongStart == "or " {
				return false
			}

			if wrongEnd == " or" {
				return false
			}
		}

		if l >= 3 {
			wrongEnd := st[l-3:]

			if wrongEnd == "not" {
				return false
			}
		}

		if l >= 4 {
			wrongStart := st[:4]
			wrongEnd := st[l-4:]

			if wrongStart == "and " {
				return false
			}

			if wrongEnd == " and" {
				return false
			}
		}

		return true
	}

	parts, _ = splitOr(st)
	if len(parts) > 1 {
		for _, part := range parts {
			if !testExpression(part) {
				return false
			}
		}

		return true
	}

	parts, _ = splitAnd(st)
	if len(parts) > 1 {
		for _, part := range parts {
			if !testExpression(part) {
				return false
			}
		}

		return true
	}

	if len(st) >= len(operatorNot) && st[:len(operatorNot)] == operatorNot {
		return testExpression(st[4:])
	}

	return false
}
