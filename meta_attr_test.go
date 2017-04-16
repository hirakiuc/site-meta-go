package sitemeta

import "testing"

func TestIsValidWithInvalidMeta(t *testing.T) {
	examples := []MetaAttr{
		{Name: "", Content: ""},
	}

	for _, ex := range examples {
		if ex.IsValid() != false {
			t.Errorf("[Example %s] should be invalid", ex.String())
		}
	}
}

func TestIsValidWithValidMeta(t *testing.T) {
	examples := []MetaAttr{
		{Name: "twitter:card", Content: "summary"},
		{Name: "twitter:title", Content: "example title"},
		{Name: "og:title", Content: "example title"},
		{Name: "og:type", Content: "video.movie"},
	}

	for _, ex := range examples {
		if ex.IsValid() != true {
			t.Errorf("[Example %s] should be valid", ex.String())
		}
	}
}
