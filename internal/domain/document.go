package domain

import (
	"errors"
	"strings"
)

// DocumentFormat represents the output format for generated documents
type DocumentFormat int

const (
	FormatUnknown DocumentFormat = iota
	FormatMarkdown
	FormatPDF
	FormatDocx
)

var formatStrings = map[DocumentFormat]string{
	FormatUnknown:  "unknown",
	FormatMarkdown: "markdown",
	FormatPDF:      "pdf",
	FormatDocx:     "docx",
}

var stringToFormat = map[string]DocumentFormat{
	"unknown":  FormatUnknown,
	"markdown": FormatMarkdown,
	"md":       FormatMarkdown,
	"pdf":      FormatPDF,
	"docx":     FormatDocx,
	"word":     FormatDocx,
}

// String returns the string representation of the format
func (f DocumentFormat) String() string {
	if str, ok := formatStrings[f]; ok {
		return str
	}
	return "unknown"
}

// ParseDocumentFormat parses a string into a DocumentFormat
func ParseDocumentFormat(s string) (DocumentFormat, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if format, ok := stringToFormat[s]; ok {
		return format, nil
	}
	return FormatUnknown, errors.New("invalid document format: " + s)
}

// Document represents a generated brag document
type Document struct {
	Format   DocumentFormat
	Template string
	Content  []byte
	Metadata DocumentMetadata
}

// DocumentMetadata contains metadata about the generated document
type DocumentMetadata struct {
	Title       string
	Author      string
	BragCount   int
	Categories  []string
	Tags        []string
	GeneratedAt string
}

// IsValid performs basic structural validation
func (d *Document) IsValid() bool {
	return d != nil &&
		d.Format != FormatUnknown &&
		len(d.Content) > 0
}
