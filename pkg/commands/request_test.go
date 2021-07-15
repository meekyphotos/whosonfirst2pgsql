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
	assert.Equal(t, -0.109401, req.longitude)
	assert.NotEmpty(t, req.preferredNames)
	assert.NotEmpty(t, req.variantNames)
	assert.Equal(t, "London", req.preferredNames["eng"])
	assert.Equal(t, "Lodoni", req.preferredNames["fij"])
	assert.Equal(t, "London", req.preferredNames["eng_ca"])
	assert.Equal(t, "LON", req.variantNames["eng"])
	assert.Equal(t, "GB", req.countryCode)
	assert.NotEmpty(t, req.metadata)
}
