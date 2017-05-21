package content

import (
	"fmt"
	"strings"
)

// MetaAttr keeps a pair of name and content.
type MetaAttr struct {
	Name    string
	Content string
}

// String return a description about this instance.
func (attr MetaAttr) String() string {
	return fmt.Sprintf("%s - %s", attr.Name, attr.Content)
}

// IsValid validate that this instance keeps valid value, or not.
func (attr MetaAttr) IsValid() bool {
	if len(attr.Name) == 0 || len(attr.Content) == 0 {
		return false
	}

	prefixes := []string{`twitter:`, `og:`}
	for _, prefix := range prefixes {
		if strings.HasPrefix(attr.Name, prefix) {
			return true
		}
	}

	return false
}
