package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_only_one_and(t *testing.T) {
	assert.True(t, testExpression("and"))
	assert.True(t, testExpression("and and"))
}
