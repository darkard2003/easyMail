package easyMail

import (
	"fmt"
	"strings"
)

type Mail struct {
	From    string
	To      []string
	Headers map[string]string
	Subject string
	Body    string
	IsHtml  bool
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

func (m *Mail) GetFormattedMails() []string {
	mails := []string{}

	var contentType string
	if m.IsHtml {
		contentType = "text/html; charset=UTF-8"
	} else {
		contentType = "text/plain; charset=UTF-8"
	}

	for _, reciver := range m.To {
		msg := "From: " + m.From + "\n" +
			"To: " + reciver + "\n" +
			"Subject: " + m.Subject + "\n" +
			contentType + "\n" +
			"\n" +
			m.Body +
			"\n" +
			"\n" +
			"."
		mails = append(mails, msg)
	}

	return mails
}

func (m *Mail) ToRFC822() string {
	headers := make([]string, 0, len(m.Headers)+3)

	// Add common headers
	headers = append(headers, fmt.Sprintf("From: %s", m.From))
	headers = append(headers, fmt.Sprintf("To: %s", strings.Join(m.To, ", ")))
	headers = append(headers, fmt.Sprintf("Subject: %s", m.Subject))

	// Add custom headers
	for k, v := range m.Headers {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}

	// Join all headers together, separate with CRLF, add CRLF at the end to separate headers from body
	headersStr := strings.Join(headers, "\r\n") + "\r\n\r\n"

	// Return headers + body
	return headersStr + m.Body
}
