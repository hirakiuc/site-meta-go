package sitemeta

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseWithInvalidUrl(t *testing.T) {
	url := "invalid url"

	ctx := context.Background()

	_, err := Parse(ctx, url)
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

	ctx := context.Background()

	result, err := Parse(ctx, ts.URL)
	if err == nil {
		t.Errorf("Error should throw.")
	}

	if result != nil {
		t.Errorf("Result should be nil. %v", result)
	}
}

func TestParseWithValidContentUrl(t *testing.T) {
	handlers := []http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "Sample Response")
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf8")
			fmt.Fprintf(w, "Sample Response")
		}),
	}

	for _, handler := range handlers {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		ctx := context.Background()

		result, err := Parse(ctx, ts.URL)
		if err != nil {
			t.Errorf("Error should not be thrown. %v", err)
		}

		if result == nil {
			t.Errorf("Result should not be nil.")
		}

		if !result.IsEmpty() {
			t.Errorf("Result should be empty.")
		}
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

		ctx := context.Background()

		result, err := Parse(ctx, ts.URL)
		if err != nil {
			t.Errorf("Error should not be thrown. %v", err)
		}

		if result == nil {
			t.Errorf("SiteMeta should not be nil. %v", result)
		}

		if !result.IsEmpty() {
			t.Errorf("result should be empty.")
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

		ctx := context.Background()

		result, err := Parse(ctx, ts.URL)
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

		ctx := context.Background()

		result, err := Parse(ctx, ts.URL)
		if err != nil {
			t.Errorf("Error should not thrown. %v", err)
		}

		if result == nil {
			t.Errorf("SiteMeta should not be nil.")
		}
	}
}

type FileExample struct {
	FileName string
	Encoding string
}

func makeExampleHandler(example FileExample) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentTypeValue := fmt.Sprintf("text/html; charset=%s", example.Encoding)
		w.Header().Set("Content-Type", contentTypeValue)

		fpath := fmt.Sprintf("../test/files/%s", example.FileName)

		file, err := os.Open(fpath)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		body, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		_, err = w.Write(body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	})
}

func TestParseWithNonUTF8ContentUrl(t *testing.T) {
	examples := []struct {
		FileName string
		Encoding string
	}{
		{FileName: "utf8.html", Encoding: "utf-8"},
		{FileName: "eucjp.html", Encoding: "euc-jp"},
		{FileName: "sjis.html", Encoding: "SHIFT_JIS"},
		{FileName: "iso2022jp.html", Encoding: "ISO-2022-JP"},
	}

	for _, example := range examples {
		handler := makeExampleHandler(example)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		ctx := context.Background()

		result, err := Parse(ctx, ts.URL)
		if err != nil {
			t.Errorf("Error should not thrown. %v", err)
		}

		if result == nil {
			t.Errorf("SiteMeta should not be nil.")
		}
	}
}

func TestParseWithInvalidEncodingUrl(t *testing.T) {
	examples := []struct {
		FileName string
		Encoding string
	}{
		{FileName: "utf8.html", Encoding: "InvalidEncoding"},
	}

	for _, example := range examples {
		handler := makeExampleHandler(example)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		ctx := context.Background()

		result, err := Parse(ctx, ts.URL)
		if err == nil {
			t.Errorf("Error should be thrown.")
		}

		if result != nil {
			t.Errorf("SiteMeta should be nil.")
		}
	}
}