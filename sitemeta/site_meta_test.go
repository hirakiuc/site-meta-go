package sitemeta

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidWithEmptyAttr(t *testing.T) {
	assert := assert.New(t)

	examples := []SiteMeta{
		{Attrs: map[string]string{}},
	}

	for _, ex := range examples {
		assert.True(ex.IsValid(), "Empty meta should be valid")
	}
}

func TestIsValidWithInvalidAttrs(t *testing.T) {
	assert := assert.New(t)

	examples := []SiteMeta{
		{
			Attrs: map[string]string{
				"description": "website description",
			},
		},
		{
			Attrs: map[string]string{
				"keywords": "golang,meta",
			},
		},
	}

	for _, ex := range examples {
		msg := fmt.Sprintf("[Example %s] should be invalid", ex.String())
		assert.False(ex.IsValid(), msg)
	}
}

func TestIsValidWithValidAttrs(t *testing.T) {
	assert := assert.New(t)

	examples := []SiteMeta{
		{
			Attrs: map[string]string{
				"twitter:card": "summary",
			},
		},
	}

	for _, ex := range examples {
		msg := fmt.Sprintf("[Example %s] should be valid", ex.String())
		assert.True(ex.IsValid(), msg)
	}
}

func TestStringWithValidAttr(t *testing.T) {
	assert := assert.New(t)

	examples := []SiteMeta{
		{
			Attrs: map[string]string{
				"twitter:card": "summary",
			},
		},
	}

	for _, ex := range examples {
		assert.Equal("twitter:card - summary", ex.String(), "SiteMeta should be stringify")
	}
}
