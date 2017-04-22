package sitemeta

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SiteMeta describe meta data of the website, like ogp, TwitterCard.
type SiteMeta struct {
	Attrs []MetaAttr
}

func parseMetaAttr(s *goquery.Selection) *MetaAttr {
	attr := MetaAttr{}

	if val, found := s.Attr("name"); found == true {
		attr.Name = val
	} else if val, found := s.Attr("property"); found == true {
		attr.Name = val
	} else {
		return nil
	}

	if val, found := s.Attr("content"); found == true {
		attr.Content = val
	}

	if attr.IsValid() == false {
		return nil
	}

	return &attr
}

func isHTMLContent(url string) (bool, error) {
	res, err := http.Head(url)
	if err != nil {
		return false, err
	}

	return (res.Header["Content-Type"][0] == "text/html"), nil
}

// String return a description about this instance.
func (meta SiteMeta) String() string {
	attrs := []string{}
	for _, attr := range meta.Attrs {
		if attr.IsValid() == true {
			attrs = append(attrs, attr.String())
		}
	}
	return strings.Join(attrs, "\n")
}

// IsValid validate that this instance keeps valid value, or not.
func (meta SiteMeta) IsValid() bool {
	if len(meta.Attrs) == 0 {
		return false
	}

	for _, attr := range meta.Attrs {
		if attr.IsValid() == false {
			return false
		}
	}

	return true
}

// Parse return SiteMeta instance if the url content has meta tags about twitter card or ogp.
func Parse(url string) (*SiteMeta, error) {
	result, err := isHTMLContent(url)
	if err != nil {
		return nil, err
	}
	if result == false {
		return nil, errors.New("Target Content seems like not html")
	}

	data := SiteMeta{Attrs: []MetaAttr{}}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		attr := parseMetaAttr(s)
		if attr != nil && attr.IsValid() == true {
			data.Attrs = append(data.Attrs, *attr)
		}
	})

	if data.IsValid() == false {
		return nil, nil
	}

	return &data, nil
}
