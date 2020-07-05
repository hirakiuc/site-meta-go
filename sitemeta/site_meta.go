package sitemeta

import (
	"encoding/json"
	"fmt"
	"strings"

	iconv "github.com/djimenez/iconv-go"
)

const (
	// DefaultEncoding describe a default encoding.
	DefaultEncoding string = "UTF-8"
)

// SiteMeta describe meta data of the website, like ogp, TwitterCard.
type SiteMeta struct {
	Attrs map[string]string
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

func (meta *SiteMeta) AddMeta(key string, val string) bool {
	if !strings.HasPrefix(key, "twitter:") && !strings.HasPrefix(key, "og:") {
		return false
	}

	meta.Attrs[key] = val

	return true
}

func (meta *SiteMeta) ToJSON() (string, error) {
	bytes, err := json.Marshal(meta.Attrs)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (meta *SiteMeta) IsEmpty() bool {
	return len(meta.Attrs) == 0
}

// IsValid validate that this instance keeps valid value, or not.
func (meta *SiteMeta) IsValid() bool {
	if meta.IsEmpty() {
		return true
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
