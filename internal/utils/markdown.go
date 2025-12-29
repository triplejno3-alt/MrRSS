package utils

import (
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// RenderMarkdown converts markdown text to safe HTML
// It uses the gomarkdown/markdown library with security configurations
func RenderMarkdown(markdownText string) string {
	if markdownText == "" {
		return ""
	}

	// Create parser with common extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Create HTML renderer with safe options
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Parse and render markdown
	htmlBytes := markdown.ToHTML([]byte(markdownText), p, renderer)

	return string(htmlBytes)
}

// RenderMarkdownInline converts markdown to HTML without wrapping <p> tags
// Useful for inline text in existing HTML structure
func RenderMarkdownInline(markdownText string) string {
	if markdownText == "" {
		return ""
	}

	// Create parser with common extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Parse and render markdown
	htmlBytes := markdown.ToHTML([]byte(markdownText), p, renderer)

	// Remove wrapping <p> tags if present
	result := string(htmlBytes)
	result = strings.TrimPrefix(result, "<p>")
	result = strings.TrimSuffix(result, "</p>")
	result = strings.TrimSuffix(result, "<p />")

	return result
}

// SanitizeHTML removes potentially dangerous HTML tags and attributes
// This is a basic sanitizer - for production use, consider using a dedicated library like bluemonday
func SanitizeHTML(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	htmlContent = scriptRegex.ReplaceAllString(htmlContent, "")

	// Remove style tags and their content
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	htmlContent = styleRegex.ReplaceAllString(htmlContent, "")

	// Remove iframe tags
	iframeRegex := regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`)
	htmlContent = iframeRegex.ReplaceAllString(htmlContent, "")

	// Remove on* event handlers (onclick, onerror, etc.)
	eventRegex := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	htmlContent = eventRegex.ReplaceAllString(htmlContent, "")

	// Remove javascript: protocol
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	htmlContent = jsRegex.ReplaceAllString(htmlContent, "")

	return htmlContent
}

// ConvertMarkdownToHTML converts markdown to safe HTML with sanitization
// This is the main function that should be used for user-generated markdown content
func ConvertMarkdownToHTML(markdownText string) string {
	if markdownText == "" {
		return ""
	}

	// First render markdown to HTML
	htmlContent := RenderMarkdown(markdownText)

	// Then sanitize the HTML to remove any potentially dangerous content
	safeHTML := SanitizeHTML(htmlContent)

	return safeHTML
}
