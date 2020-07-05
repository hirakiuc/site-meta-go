package sitemeta

import (
	"testing"
)

func TestIsValidWithEmptyAttr(t *testing.T) {
	examples := []SiteMeta{
		{Attrs: map[string]string{}},
	}

	for _, ex := range examples {
		if !ex.IsValid() {
			t.Errorf("[Example %s] should be invalid", ex.String())
		}
	}
}

func TestIsValidWithInvalidAttrs(t *testing.T) {
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
		if ex.IsValid() != false {
			t.Errorf("[Example %s] should be invalid", ex.String())
		}
	}
}

func TestIsValidWithValidAttrs(t *testing.T) {
	examples := []SiteMeta{
		{
			Attrs: map[string]string{
				"twitter:card": "summary",
			},
		},
	}

	for _, ex := range examples {
		if ex.IsValid() != true {
			t.Errorf("[Example %s] should be valid", ex.String())
		}
	}
}

func TestStringWithValidAttr(t *testing.T) {
	examples := []SiteMeta{
		{
			Attrs: map[string]string{
				"twitter:card": "summary",
			},
		},
	}

	for _, ex := range examples {
		str := ex.String()
		if str != "twitter:card - summary" {
			t.Errorf("[Example %s] should be stringify", str)
		}
	}
}
