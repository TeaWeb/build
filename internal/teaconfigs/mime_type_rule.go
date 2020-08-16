package teaconfigs

import "regexp"

// mime type
type MimeTypeRule struct {
	Value  string
	Regexp *regexp.Regexp
}
