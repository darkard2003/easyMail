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

func writeHeader(buf *bytes.Buffer, headers map[string]string) {
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
}

func writeAttachment(buf *bytes.Buffer, filename string, data []byte, boundary string) {
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: application/octet-stream\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
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

func getFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []byte(""), err
	}
	defer f.Close()
	b := new(bytes.Buffer)
	if _, err := b.ReadFrom(f); err != nil {
		return []byte(""), err
	}
	return b.Bytes(), nil
}

func (m *Mail) AttachTextFile(src string) error {

	if !strings.HasSuffix(src, ".txt") {
		return fmt.Errorf("file %s is not a text file", src)
	}

	b, err := getFile(src)
	if err != nil {
		return err
	}
	m.Body = string(b)
	m.isHtml = false
	return nil

}

func (m *Mail) AddHtmlFile(path string) error {
	if !strings.HasSuffix(path, ".html") {
		return fmt.Errorf("file %s is not a html file", path)
	}
	b, err := getFile(path)
	if err != nil {
		return err
	}
	m.Body = string(b)
	m.isHtml = true
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

func (m *Mail) Raw() (string, error) {
	buf, err := m.ToBytes()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (m *Mail) ToBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	headers := make(map[string]string, len(m.Headers)+7)

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
