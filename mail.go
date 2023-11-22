package easyMail

import (
	"fmt"
	"strings"
)

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
	IsHtml  bool
	Attachments []string
	Headers map[string]string
}

func NewMail(
	from string,
	to []string,
	subject string,
	body string,
	isHtml bool,
) *Mail {
	return &Mail{
		From:    from,
		To:      to,
		Headers: make(map[string]string),
		Subject: subject,
		Body:    body,
		IsHtml:  isHtml,
	}
}

func (m *Mail) AddHeader(key string, value string) {
	m.Headers[key] = value
}

func (m *Mail) AddAttachment(path string) error {
	m.Attachments = append(m.Attachments, path)
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