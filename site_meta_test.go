package sitemeta

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsValidWithEmptyAttr(t *testing.T) {
	examples := []SiteMeta{
		{Attrs: []MetaAttr{}},
	}

	for _, ex := range examples {
		if ex.IsValid() != false {
			t.Errorf("[Example %s] should be invalid", ex.String())
		}
	}
}

func TestIsValidWithInvalidAttrs(t *testing.T) {
	examples := []SiteMeta{
		{
			Attrs: []MetaAttr{
				{Name: "description", Content: "website description"},
			},
		},
		{
			Attrs: []MetaAttr{
				{Name: "keywords", Content: "golang,meta"},
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
			Attrs: []MetaAttr{
				{Name: "twitter:card", Content: "summary"},
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
			Attrs: []MetaAttr{
				{Name: "twitter:card", Content: "summary"},
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

func TestParseWithInvalidUrl(t *testing.T) {
	url := "invalid url"
	_, err := Parse(url)
	if err == nil {
		t.Errorf("Invalid error should return error")
	}
}

func TestParseWithInvalidContentUrl(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		fmt.Fprintf(w, "Sample Response")
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	result, err := Parse(ts.URL)
	if err == nil {
		t.Errorf("Error should throw.")
	}

	if result != nil {
		t.Errorf("Result should be nil. %v", result)
	}
}

func TestParseWithValidContentUrl(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Sample Response")
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	result, err := Parse(ts.URL)
	if err != nil {
		t.Errorf("Error should not thrown. %v", err)
	}
	if result != nil {
		t.Errorf("Result should be nil. %v", result)
	}
}

func TestParseWithValidContentUrlWithSiteMetaTag(t *testing.T) {
	handlers := []http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><meta name="viewport" content="width=device-width"></html>`)
		}),
	}

	for _, handler := range handlers {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		result, err := Parse(ts.URL)
		if err != nil {
			t.Errorf("Error should not thrown. %v", err)
		}
		if result != nil {
			t.Errorf("SiteMeta should be nil. %v", result)
		}
	}
}

func TestParseWithValidContentUrlWithTwitterCard(t *testing.T) {
	handlers := []http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><meta name="twitter:card" content="summary"></html>`)
		}),
	}

	for _, handler := range handlers {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		result, err := Parse(ts.URL)
		if err != nil {
			t.Errorf("Error should not thrown. %v", err)
		}
		if result == nil {
			t.Errorf("SiteMeta should not be nil. %v", err)
		}
	}
}

func TestParseWithValidContentUrlWithOgp(t *testing.T) {
	handlers := []http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><meta property="og:type" content="video.movie" /></html>`)
		}),
	}

	for _, handler := range handlers {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		result, err := Parse(ts.URL)
		if err != nil {
			t.Errorf("Error should not thrown. %v", err)
		}
		if result == nil {
			t.Errorf("SiteMeta should not be nil. %v", err)
		}
	}
}
