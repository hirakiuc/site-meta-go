package sitemeta

import (
	"errors"
	"fmt"
	"strings"

	iconv "github.com/djimenez/iconv-go"

	content "github.com/hirakiuc/site-meta-go/internal/content"
	logging "github.com/hirakiuc/site-meta-go/internal/logger"
)

const (
	// DefaultEncoding describe a default encoding.
	DefaultEncoding string = "utf-8"
)

// SiteMeta describe meta data of the website, like ogp, TwitterCard.
type SiteMeta struct {
	Attrs map[string]string
}

var logger = logging.GetLogger()

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
		logger.Printf("IsValid: key:%s", key)
		if strings.HasPrefix(key, "twitter:") || strings.HasPrefix(key, "og:") {
			continue
		}

		return false
	}

	return true
}

func (meta *SiteMeta) convertEncoding(toEncoding string) error {
	logger.Printf("Convert encoding: from:%s to:%s", toEncoding, DefaultEncoding)
	if toEncoding == DefaultEncoding {
		logger.Printf("  Skipped.")
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
			logger.Printf("Failed to convert encoding: %s %v\n", key, err)
			return err
		}

		newValue, err := converter.ConvertString(value)
		if err != nil {
			logger.Printf("Failed to convert encoding: %s %v\n", value, err)
			return err
		}
		result[newKey] = newValue
	}
	meta.Attrs = result

	return nil
}

// Parse return SiteMeta instance if the url content has meta tags about twitter card or ogp.
func Parse(url string) (*SiteMeta, error) {
	data := newSiteMeta()

	html, err := content.FetchHTMLContent(url)
	if err != nil {
		return nil, err
	}
	data.Attrs = html.MetaAttrs()

	err = data.convertEncoding(html.ContentEncoding)
	if err != nil {
		return nil, err
	}

	if data.IsValid() == false {
		return nil, errors.New("This site can't parse SiteMeta")
	}

	return &data, nil
}
