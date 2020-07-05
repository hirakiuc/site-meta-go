package sitemeta

import (
	"context"
	"errors"
	"fmt"
	"strings"

	iconv "github.com/djimenez/iconv-go"

	content "github.com/hirakiuc/site-meta-go/sitemeta/content"
)

const (
	// DefaultEncoding describe a default encoding.
	DefaultEncoding string = "UTF-8"
)

//nolint
var CantParseErr error = errors.New("this site can't parse SiteMeta")

// SiteMeta describe meta data of the website, like ogp, TwitterCard.
type SiteMeta struct {
	Attrs map[string]string
}

func newSiteMeta() SiteMeta {
	return SiteMeta{Attrs: map[string]string{}}
}

// String return a description about this instance.
func (meta *SiteMeta) String() string {
	attrs := []string{}

	for key, value := range meta.Attrs {
		str := fmt.Sprintf("%s - %s", key, value)
		attrs = append(attrs, str)
	}

	return strings.Join(attrs, "\n")
}

// IsValid validate that this instance keeps valid value, or not.
func (meta *SiteMeta) IsValid() bool {
	if len(meta.Attrs) == 0 {
		return false
	}

	for key := range meta.Attrs {
		if strings.HasPrefix(key, "twitter:") || strings.HasPrefix(key, "og:") {
			continue
		}

		return false
	}

	return true
}

func (meta *SiteMeta) convertEncoding(toEncoding string) error {
	if toEncoding == DefaultEncoding {
		return nil
	}

	converter, err := iconv.NewConverter(toEncoding, DefaultEncoding)
	if err != nil {
		return err
	}
	defer converter.Close()

	result := map[string]string{}

	for key, value := range meta.Attrs {
		newKey, err := converter.ConvertString(key)
		if err != nil {
			return err
		}

		newValue, err := converter.ConvertString(value)
		if err != nil {
			return err
		}

		result[newKey] = newValue
	}

	meta.Attrs = result

	return nil
}

// Parse return SiteMeta instance if the url content has meta tags about twitter card or ogp.
func Parse(c context.Context, url string) (*SiteMeta, error) {
	data := newSiteMeta()

	html, err := content.FetchHTMLContent(c, url)
	if err != nil {
		return nil, err
	}

	data.Attrs = html.MetaAttrs()

	err = data.convertEncoding(html.ContentEncoding)
	if err != nil {
		return nil, err
	}

	if !data.IsValid() {
		return nil, fmt.Errorf("%w", CantParseErr)
	}

	return &data, nil
}
