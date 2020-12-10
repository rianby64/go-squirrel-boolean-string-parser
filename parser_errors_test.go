package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Case without parentheses
func Test_error_case1(t *testing.T) {
	assert.False(t, testExpression(""))

	assert.True(t, testExpression("alice and bob"))
	assert.True(t, testExpression("not alice and bob"))
	assert.True(t, testExpression("alice and not bob"))
	assert.True(t, testExpression("not alice and not bob"))

	assert.True(t, testExpression("alice or bob"))
	assert.True(t, testExpression("not alice or bob"))
	assert.True(t, testExpression("alice or not bob"))
	assert.True(t, testExpression("not alice or not bob"))

	assert.True(t, testExpression("alice"))
	assert.True(t, testExpression("not alice"))

	assert.False(t, testExpression("alice and and bob"))
	assert.False(t, testExpression("alice or or bob"))
	assert.False(t, testExpression("alice and or bob"))
	assert.False(t, testExpression("alice or and bob"))

	assert.False(t, testExpression("not alice and and bob"))
	assert.False(t, testExpression("not alice or or bob"))
	assert.False(t, testExpression("not alice and or bob"))
	assert.False(t, testExpression("not alice or and bob"))

	assert.False(t, testExpression("alice and and not bob"))
	assert.False(t, testExpression("alice or or not bob"))
	assert.False(t, testExpression("alice and or not bob"))
	assert.False(t, testExpression("alice or and not bob"))

	assert.False(t, testExpression("not alice and and not bob"))
	assert.False(t, testExpression("not alice or or not bob"))
	assert.False(t, testExpression("not alice and or not bob"))
	assert.False(t, testExpression("not alice or and not bob"))

	assert.False(t, testExpression("and alice and bob"))
	assert.False(t, testExpression("or alice or bob"))
	assert.False(t, testExpression("and alice or bob"))
	assert.False(t, testExpression("or alice and bob"))

	assert.False(t, testExpression("alice and bob and"))
	assert.False(t, testExpression("alice or bob or"))
	assert.False(t, testExpression("alice and bob or"))
	assert.False(t, testExpression("alice or bob and"))

	assert.False(t, testExpression("alice not"))
	assert.False(t, testExpression("alice not bob"))
}

// Case with parentheses
func Test_error_case2(t *testing.T) {
	assert.True(t, testExpression("(alice and bob)"))
	assert.True(t, testExpression("(alice and bob) and (carol and dan)"))
	assert.True(t, testExpression("(alice and bob) and (carol and dan) and (elen and gleen)"))
	assert.True(t, testExpression("not (alice and bob) and (carol and dan) and (elen and gleen)"))
	assert.True(t, testExpression("not (alice and bob) and not (carol and dan) and (elen and gleen)"))
	assert.True(t, testExpression("not (alice and bob) and not (carol and dan) and not (elen and gleen)"))

	assert.True(t, testExpression("(alice or bob)"))
	assert.True(t, testExpression("(alice or bob) or (carol or dan)"))
	assert.True(t, testExpression("(alice or bob) or (carol or dan) or (elen or gleen)"))
	assert.True(t, testExpression("not (alice or bob) or (carol or dan) or (elen or gleen)"))
	assert.True(t, testExpression("not (alice or bob) or not (carol or dan) or (elen or gleen)"))
	assert.True(t, testExpression("not (alice or bob) or not (carol or dan) or not (elen or gleen)"))

	assert.True(t, testExpression("(alice or bob)"))
	assert.True(t, testExpression("(alice or bob) and (carol or dan)"))
	assert.True(t, testExpression("(alice or bob) and (carol or dan) and (elen or gleen)"))
	assert.True(t, testExpression("not (alice or bob) and (carol or dan) and (elen or gleen)"))
	assert.True(t, testExpression("not (alice or bob) and not (carol or dan) and (elen or gleen)"))
	assert.True(t, testExpression("not (alice or bob) and (not (carol or dan)) and not (elen or gleen)"))

	assert.True(t, testExpression("not (alice or bob) and not (carol or dan) and not (elen or gleen)"))

	assert.True(t, testExpression("((not (alice or bob)) and (not (carol or dan)))"))
	assert.True(t, testExpression("(alice or bob) and (not ((carol or dan) and (elen or gleen)))"))
	assert.True(t, testExpression("((alice or bob) and ((carol or dan) and (elen or gleen)))"))
	assert.True(t, testExpression("((alice or bob) and (carol or dan)) and (elen or gleen) and ((not ((alice or bob) and (carol or dan)) and (elen or gleen)))"))
}

func Test_error_case3(t *testing.T) {
	assert.False(t, testExpression("alice and  "))
	assert.False(t, testExpression("alice and or bob "))
	assert.False(t, testExpression("(alice and bob))"))
	assert.False(t, testExpression("((alice and bob)"))
	assert.False(t, testExpression("()"))
	assert.False(t, testExpression("(())"))
}
