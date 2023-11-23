package easyMail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Mail struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	isHtml      bool
	Attachments map[string][]byte
	Headers     map[string]string
}

func (m *Mail) getContentType(boundary string) string {
	if m.hasAttachments() {
		return "multipart/mixed; boundary=" + boundary
	}

	if m.isHtml {
		return "text/html; charset=utf-8"
	}

	return "text/plain; charset=utf-8"
}

func (m *Mail) hasAttachments() bool {
	return len(m.Attachments) > 0
}

func (m *Mail) hasCc() bool {
	return len(m.Cc) > 0
}

func (m *Mail) hasBcc() bool {
	return len(m.Bcc) > 0
}

func writeHeader(buf *bytes.Buffer, headers map[string]string) {
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
}

func writeAttachment(buf *bytes.Buffer, filename string, data []byte, boundary string) {
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString(fmt.Sprintf("Content-Type: application/octet-stream\r\n"))
	buf.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n"))
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename))
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b64, data)
	buf.Write(b64)
	buf.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
}

func NewMail(
	from string,
	to []string,
	subject string,
	body string,
	isHtml bool,
) *Mail {
	return &Mail{
		From:        from,
		To:          to,
		Cc:          make([]string, 0),
		Bcc:         make([]string, 0),
		Headers:     make(map[string]string),
		Attachments: make(map[string][]byte),
		Subject:     subject,
		Body:        body,
		isHtml:      isHtml,
	}
}

func (m *Mail) AttachHtml(html string) {
	m.Body = html
	m.isHtml = true
}

func (m *Mail) AttachText(text string) {
	m.Body = text
	m.isHtml = false
}

func (m *Mail) AddHtmlFile(path string) error {
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	f, e := os.Open(fullpath)
	if e != nil {
		return e
	}
	defer f.Close()
	b := new(bytes.Buffer)
	if _, err := b.ReadFrom(f); err != nil {
		return err
	}
	m.AttachHtml(b.String())
	return nil
}

func (m *Mail) AddHeader(key string, value string) {
	m.Headers[key] = value
}

func (m *Mail) AddAttachment(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, filename := filepath.Split(path)
	m.Attachments[filename] = b
	return nil
}

func (m *Mail) AddCc(cc string) {
	m.Cc = append(m.Cc, cc)
}

func (m *Mail) AddBcc(bcc string) {
	m.Bcc = append(m.Bcc, bcc)
}

func (m *Mail) AddTo(to string) {
	m.To = append(m.To, to)
}

func (m *Mail) AddToAll(to []string) {
	m.To = append(m.To, to...)
}

func (m *Mail) AddCcAll(cc []string) {
	m.Cc = append(m.Cc, cc...)
}

func (m *Mail) AddBccAll(bcc []string) {
	m.Bcc = append(m.Bcc, bcc...)
}

func (m *Mail) AddAttachmentAll(paths []string) error {
	for _, path := range paths {
		err := m.AddAttachment(path)
		if err != nil {
			return err
		}
	}

	return nil

}

func (m *Mail) ToRFC822() (string, error) {
	headers := make([]string, 0, len(m.Headers)+4)

	// Add common headers
	headers = append(headers, fmt.Sprintf("From: %s", m.From))
	headers = append(headers, fmt.Sprintf("To: %s", strings.Join(m.To, ", ")))
	headers = append(headers, fmt.Sprintf("Subject: %s", m.Subject))

	// Add Content-Type header based on IsHtml flag

	// Add custom headers
	for k, v := range m.Headers {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}

	// Join all headers together, separate with CRLF, add CRLF at the end to separate headers from body
	headersStr := strings.Join(headers, "\r\n") + "\r\n\r\n"

	// Return headers + body
	return headersStr + m.Body, nil
}

func (m *Mail) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	headers := make(map[string]string, len(m.Headers)+4)

	for k, v := range m.Headers {
		headers[k] = v
	}

	headers["From"] = m.From
	headers["To"] = strings.Join(m.To, ", ")
	headers["Subject"] = m.Subject

	headers["Cc"] = strings.Join(m.Cc, ", ")
	headers["Bcc"] = strings.Join(m.Bcc, ", ")

	headers["MIME-Version"] = "1.0"

	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	headers["Content-Type"] = m.getContentType(boundary)
	writeHeader(buf, headers)

	if m.hasAttachments() {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}

	buf.WriteString(m.Body)

	if m.hasAttachments() {
		for filename, data := range m.Attachments {
			writeAttachment(buf, filename, data, boundary)
		}
		buf.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	}

	return buf.Bytes(), nil
}
