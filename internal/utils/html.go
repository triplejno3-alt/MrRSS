package utils

import (
	"regexp"
	"strings"
)

// selfClosingTags is the list of HTML self-closing tags to handle
const selfClosingTags = "img|br|hr|input|meta|link"

// CleanHTML sanitizes HTML content by fixing common malformed patterns
// that can cause rendering issues.
func CleanHTML(html string) string {
	if html == "" {
		return html
	}

	// Fix malformed opening tags like <p--> to <p>
	// This pattern matches tags like <p-->, <div-->, etc.
	malformedTagRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)\s*--+>`)
	html = malformedTagRegex.ReplaceAllString(html, "<$1>")

	// Fix malformed self-closing tags like <img-->, <br--> to <img>, <br>
	// Some feeds have broken self-closing tags with or without attributes
	// Pattern 1: Tags with attributes (e.g., <img src="..." -->)
	// Use [^<>]+ to avoid matching angle brackets and nested tags
	malformedSelfClosingWithAttrs := regexp.MustCompile(`<(` + selfClosingTags + `)\s+([^<>]+?)--+>`)
	html = malformedSelfClosingWithAttrs.ReplaceAllString(html, "<$1 $2>")
	
	// Pattern 2: Tags without attributes (e.g., <br-->)
	malformedSelfClosingNoAttrs := regexp.MustCompile(`<(` + selfClosingTags + `)\s*--+>`)
	html = malformedSelfClosingNoAttrs.ReplaceAllString(html, "<$1>")

	// Trim any leading/trailing whitespace
	html = strings.TrimSpace(html)

	return html
}
