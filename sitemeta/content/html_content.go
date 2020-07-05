package content

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	// DefaultEncoding describe a default encoding.
	DefaultEncoding string = "UTF-8"

	// Expected number of Charset header components
	CharSetComponent int = 2
)

//nolint
var notHTMLContentErr error = errors.New("target content seems like not html")

// HTMLContent describe the html page of the URL.
type HTMLContent struct {
	ContentEncoding string
	URL             string
	doc             *goquery.Document
}

func extractCharset(contentType string) string {
	exp := regexp.MustCompile(`charset=((\w|\d|\-)*)`)
	group := exp.FindStringSubmatch(contentType)

	if len(group) < CharSetComponent {
		return DefaultEncoding
	}

	// TBD Normalize target encoding.
	return strings.ToUpper(group[1])
}

func fetchContent(c context.Context, url string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Jar: jar}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c)

	ch := make(chan error)
	defer close(ch)

	var resp *http.Response

	go func(r **http.Response) {
		//nolint
		// NOTE: this response will be closed at caller code.
		*r, err = client.Do(req)
		if err != nil {
			ch <- err
			return
		}

		ch <- nil
	}(&resp)

	// Wait for the content
	err = <-ch
	if err != nil {
		return nil, err
	}

	return resp, nil
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
func FetchHTMLContent(c context.Context, url string) (*HTMLContent, error) {
	content := NewHTMLContent(url)

	isValid, err := content.isValidContent()
	if err != nil {
		return nil, err
	} else if !isValid {
		return nil, fmt.Errorf("%w", notHTMLContentErr)
	}

	res, err := fetchContent(c, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content.doc, err = goquery.NewDocumentFromResponse(res)
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
	defer res.Body.Close()

	contentType := strings.TrimSpace(res.Header["Content-Type"][0])

	isHTMLContent := strings.Contains(contentType, "text/html")
	if isHTMLContent {
		// Update ContentEncoding with charset which in Content-Type response header.
		content.ContentEncoding = extractCharset(contentType)
	}

	return isHTMLContent, nil
}

func (content *HTMLContent) extractCharset() {
	selections := (content.doc).Find("meta[charset]")

	if selections.Length() > 0 {
		s := selections.First()

		charset, found := s.Attr("charset")
		if found {
			content.ContentEncoding = strings.ToUpper(charset)

			return
		}
	}

	selections = (content.doc).Find("meta[http-equiv=content-type]")

	if selections.Length() > 0 {
		s := selections.First()

		str, found := s.Attr("content")
		if found {
			content.ContentEncoding = extractCharset(str)

			return
		}
	}
}

// MetaAttrs return the Meta attributes of the HTMLContent.
func (content *HTMLContent) MetaAttrs() map[string]string {
	// Extract Charset from doc
	content.extractCharset()

	selection := (content.doc).Find("meta[name]").Add("meta[property]")

	if selection.Length() == 0 {
		return map[string]string{}
	}

	result := make([]MetaAttr, selection.Length())

	for _, node := range selection.Nodes {
		attr := parseMetaNode(node)
		if attr == nil {
			continue
		}

		result = append(result, *attr)
	}

	return attrsToMap(result)
}

func attrsToMap(attrs []MetaAttr) map[string]string {
	result := map[string]string{}

	for _, attr := range attrs {
		if attr.IsValid() {
			result[attr.Name] = attr.Content
		}
	}

	return result
}

func parseMetaNode(node *html.Node) *MetaAttr {
	v := MetaAttr{}

	for _, attr := range node.Attr {
		name := strings.ToLower(attr.Key)

		// Find name attr or property attr
		if name == "name" || name == "property" {
			v.Name = strings.ToLower(attr.Val)

			continue
		}

		// Find content attribute
		if name == "content" {
			v.Content = attr.Val

			continue
		}
	}

	if !v.IsValid() {
		return nil
	}

	return &v
}
