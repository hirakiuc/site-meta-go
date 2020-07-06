package sitemeta

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWithInvalidUrl(t *testing.T) {
	assert := assert.New(t)

	url := "invalid url"

	ctx := context.Background()

	_, err := Parse(ctx, url)
	assert.NotNil(err, "Invalid error should return error")
}

func TestParseWithInvalidContentUrl(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		fmt.Fprintf(w, "Sample Response")
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	ctx := context.Background()

	result, err := Parse(ctx, ts.URL)

	assert.NotNil(err, "Error should throw")
	assert.Nil(result, "Result should be nil")
}

func TestParseWithValidContentUrl(t *testing.T) {
	assert := assert.New(t)

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

		assert.Nil(err, "Error should not be thrown")
		assert.NotNil(result, "result should not be nil.")
		assert.True(result.IsEmpty(), "result should be empty")
	}
}

func TestParseWithValidContentUrlWithSiteMetaTag(t *testing.T) {
	assert := assert.New(t)

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

		assert.Nil(err, "Error should not be thrown")
		assert.NotNil(result, "SiteMeta should not be nil")
		assert.True(result.IsEmpty(), "result should be empty")
	}
}

func TestParseWithValidContentUrlWithTwitterCard(t *testing.T) {
	assert := assert.New(t)

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

		assert.Nil(err, "error should not be thrown")
		assert.NotNil(result, "SIteMeta should not be nil")
	}
}

func TestParseWithValidContentUrlWithOgp(t *testing.T) {
	assert := assert.New(t)

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

		assert.Nil(err, "error should not be thrown")
		assert.NotNil(result, "SiteMeta should not be nil.")
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
	assert := assert.New(t)

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

		assert.Nil(err, "error should not be thrown")
		assert.NotNil(result, "SiteMeta should not be nil.")
	}
}

func TestParseWithInvalidEncodingUrl(t *testing.T) {
	assert := assert.New(t)

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

		assert.NotNil(err, "error should not be thrown")
		assert.Nil(result, "SiteMeta should be nil")
	}
}
