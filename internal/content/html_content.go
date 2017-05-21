package content

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	logging "github.com/hirakiuc/site-meta-go/internal/logger"
)

const (
	// DefaultEncoding describe a default encoding.
	DefaultEncoding string = "utf-8"
)

// HTMLContent describe the html page of the URL.
type HTMLContent struct {
	ContentEncoding string
	URL             string
	doc             *goquery.Document
}

var logger = logging.GetLogger()

//------
func extractCharset(contentType string) string {
	logger.Printf("ContentType: %s\n", contentType)

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

func parseMetaAttr(key string, s *goquery.Selection) *MetaAttr {
	attr := MetaAttr{}

	if val, found := s.Attr(key); found == true {
		attr.Name = strings.TrimSpace(val)
	} else if val, found := s.Attr("property"); found == true {
		attr.Name = strings.TrimSpace(val)
	} else {
		return nil
	}

	if val, found := s.Attr("content"); found == true {
		attr.Content = strings.TrimSpace(val)
	}

	if attr.IsValid() == false {
		return nil
	}

	return &attr
}

func attrsToMap(attrs []MetaAttr) map[string]string {
	result := map[string]string{}
	for _, attr := range attrs {
		result[attr.Name] = attr.Content
	}
	return result
}

//------

// NewHTMLContent return a new HTMLContent instance.
func NewHTMLContent(url string) HTMLContent {
	return HTMLContent{
		ContentEncoding: DefaultEncoding,
		URL:             url,
		doc:             nil,
	}
}

// FetchHTMLContent fetch the url content, check and return HTMLContent instance.
func FetchHTMLContent(url string) (*HTMLContent, error) {
	content := NewHTMLContent(url)

	isValid, err := content.isValidContent()
	if err != nil {
		logger.Printf("Error: %v", err)
		return nil, err
	} else if isValid == false {
		return nil, errors.New("Target Content seems like not html")
	}

	content.doc, err = goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (content *HTMLContent) isValidContent() (bool, error) {
	res, err := http.Head(content.URL)
	if err != nil {
		return false, err
	}

	contentType := strings.TrimSpace(res.Header["Content-Type"][0])
	logger.Printf("contentType:%s", contentType)

	isHTMLContent := (strings.Index(contentType, "text/html") != -1)
	if isHTMLContent == true {
		// Update ContentEncoding with charset which in Content-Type response header.
		content.ContentEncoding = extractCharset(contentType)
	}

	return isHTMLContent, nil
}

func (content *HTMLContent) extractCharset() {
	logger.Printf("extractCharset")
	selections := (content.doc).Find("meta[charset]")
	if selections.Length() > 0 {
		s := selections.First()
		charset, found := s.Attr("charset")
		if found {
			logger.Printf("extractCharset:%s", charset)
			content.ContentEncoding = charset
			return
		}
	}

	selections = (content.doc).Find("meta[http-equiv=content-type]")
	if selections.Length() > 0 {
		s := selections.First()
		str, found := s.Attr("content")
		if found {
			logger.Printf("extractCharset:%s", str)
			content.ContentEncoding = extractCharset(str)
			return
		}
	}
}

// MetaAttrs return the Meta attributes of the HTMLContent.
func (content *HTMLContent) MetaAttrs() map[string]string {
	// Extract Charset from doc
	content.extractCharset()

	// Extract MetaAttrs from doc
	result := []MetaAttr{}
	for _, key := range []string{"name", "property"} {
		selector := fmt.Sprintf("meta[%s]", key)

		selections := (content.doc).Find(selector)
		if selections.Length() == 0 {
			continue
		}

		attrs := []MetaAttr{}
		selections.Each(func(_ int, s *goquery.Selection) {
			if attr := parseMetaAttr(key, s); attr != nil {
				attrs = append(attrs, *attr)
			}
		})

		result = append(result, attrs...)
	}

	return attrsToMap(result)
}
