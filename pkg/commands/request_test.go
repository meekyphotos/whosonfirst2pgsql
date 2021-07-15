package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestInit(t *testing.T) {
	v, err := fastjson.Parse(TestData)

	assert.Nil(t, err)

	req := fromJsonValue(v)

	assert.Equal(t, int64(101750367), req.id)
	assert.Equal(t, 51.500526, req.latitude)
}
