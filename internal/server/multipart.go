package server

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
)

// MultipartWriter wraps multipart.Writer with additional functionality
type MultipartWriter struct {
	*multipart.Writer
}

// NewMultipartWriter creates a new MultipartWriter
func NewMultipartWriter(w io.Writer) *MultipartWriter {
	return &MultipartWriter{
		Writer: multipart.NewWriter(w),
	}
}

// CreateFormFile creates a new form file field
func (mw *MultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", "application/octet-stream")
	return mw.CreatePart(h)
}

// escapeQuotes escapes quotes in strings for multipart headers
func escapeQuotes(s string) string {
	return strings.Replace(s, `"`, `\"`, -1)
}