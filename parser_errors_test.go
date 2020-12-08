package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_error_case1(t *testing.T) {
	p := Parser{}

	assert.False(t, p.testExpression(""))

	assert.True(t, p.testExpression("alice and bob"))
	assert.True(t, p.testExpression("not alice and bob"))
	assert.True(t, p.testExpression("alice and not bob"))
	assert.True(t, p.testExpression("not alice and not bob"))

	assert.True(t, p.testExpression("alice or bob"))
	assert.True(t, p.testExpression("not alice or bob"))
	assert.True(t, p.testExpression("alice or not bob"))
	assert.True(t, p.testExpression("not alice or not bob"))

	assert.True(t, p.testExpression("alice"))
	assert.True(t, p.testExpression("not alice"))

	assert.False(t, p.testExpression("alice and and bob"))
	assert.False(t, p.testExpression("alice or or bob"))
	assert.False(t, p.testExpression("alice and or bob"))
	assert.False(t, p.testExpression("alice or and bob"))
	assert.False(t, p.testExpression("alice not"))
	assert.False(t, p.testExpression("alice not bob"))
}
