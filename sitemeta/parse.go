package sitemeta

import (
	"context"
	"errors"
	"fmt"

	content "github.com/hirakiuc/site-meta-go/sitemeta/content"
)

//nolint
var CantParseErr error = errors.New("this site can't parse SiteMeta")

// Parse return SiteMeta instance if the url content has meta tags about twitter card or ogp.
func Parse(c context.Context, url string) (*SiteMeta, error) {
	meta := SiteMeta{
		Attrs: map[string]string{},
	}

	html, err := content.FetchHTMLContent(c, url)
	if err != nil {
		return nil, err
	}

	meta.Attrs = html.MetaAttrs()

	if meta.IsEmpty() {
		return &meta, nil
	}

	err = meta.convertEncoding(html.ContentEncoding)
	if err != nil {
		return nil, err
	}

	if !meta.IsValid() {
		return nil, fmt.Errorf("%w", CantParseErr)
	}

	return &meta, nil
}
