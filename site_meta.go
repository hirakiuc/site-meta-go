package sitemeta

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

// SiteMeta describe meta data of the website, like ogp, TwitterCard.
type SiteMeta struct {
	Attrs []MetaAttr

	isHTMLContent   bool
	contentEncoding string
}

const (
	// DefaultEncoding describe a default encoding.
	DefaultEncoding string = "utf-8"
)

func init() {
	initLogger()
}

func newSiteMeta() SiteMeta {
	return SiteMeta{
		Attrs: []MetaAttr{},

		isHTMLContent:   false,
		contentEncoding: DefaultEncoding,
	}
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

func extractCharset(contentType string) string {
	logger.Printf("contentType: %s\n", contentType)

	exp := regexp.MustCompile(`charset=((\w|\d|\-)*)`)
	group := exp.FindStringSubmatch(contentType)

	if len(group) < 2 {
		logger.Printf("len(group) less than 2")
		return DefaultEncoding
	}

	// TODO: Normalize target encoding.
	logger.Printf("Found charset: %s\n", group[1])
	return group[1]
}

func (meta *SiteMeta) isValidContent(url string) (bool, error) {
	res, err := http.Head(url)
	if err != nil {
		return false, err
	}

	contentType := strings.TrimSpace(res.Header["Content-Type"][0])
	if strings.Index(contentType, "text/html") != -1 {
		meta.isHTMLContent = true
	}

	meta.contentEncoding = extractCharset(contentType)

	return meta.isHTMLContent, nil
}

// String return a description about this instance.
func (meta *SiteMeta) String() string {
	attrs := []string{}
	for _, attr := range meta.Attrs {
		if attr.IsValid() == true {
			attrs = append(attrs, attr.String())
		}
	}
	return strings.Join(attrs, "\n")
}

// IsValid validate that this instance keeps valid value, or not.
func (meta *SiteMeta) IsValid() bool {
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

func (meta *SiteMeta) convertEncoding() error {
	logger.Printf("Target encoding: %s\n", meta.contentEncoding)
	converter, err := iconv.NewConverter(meta.contentEncoding, DefaultEncoding)
	if err != nil {
		return err
	}
	defer converter.Close()

	for idx, attr := range meta.Attrs {
		attr.Name, err = converter.ConvertString(attr.Name)
		if err != nil {
			logger.Printf("Failed to convert encoding: %s %v\n", attr.Name, err)
			return err
		}

		attr.Content, err = converter.ConvertString(attr.Content)
		if err != nil {
			logger.Printf("Failed to convert encoding: %s %v\n", attr.Name, err)
			return err
		}
		meta.Attrs[idx] = attr
	}
	return nil
}

// Parse return SiteMeta instance if the url content has meta tags about twitter card or ogp.
func Parse(url string) (*SiteMeta, error) {
	data := newSiteMeta()

	result, err := data.isValidContent(url)
	if err != nil {
		return nil, err
	}
	if result == false {
		return nil, errors.New("Target Content seems like not html")
	}

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

	if err := data.convertEncoding(); err != nil {
		return nil, err
	}
	if data.IsValid() == false {
		return nil, nil
	}

	return &data, nil
}
