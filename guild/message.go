package guild

import (
	smail "github.com/xhit/go-simple-mail/v2"
)

// Message wraps an email struct
type Message struct {
	email             *smail.Email
	textAsAlternative bool
}

func NewMessage() *Message {
	return &Message{
		email: smail.NewMSG(),
	}
}

func (c *Message) SetPriorityHigh() {
	c.email.SetPriority(smail.PriorityHigh)
}

func (c *Message) SetPriorityLow() {
	c.email.SetPriority(smail.PriorityLow)
}

func (c *Message) SetPriority(isHigh bool) {
	if isHigh {
		c.SetPriorityHigh()
	}
}

// SetSubject sets the email subject to the provided string value.
func (c *Message) SetSubject(subject string) {
	c.email.SetSubject(subject)
}

// SetTextBody sets the raw text email body. If the text body has
// been set after the Html body, it is set as the alternative body
// of the email, otherwise it will be treated as the primary.
// Setting the text body with an empty string has no effect.
func (c *Message) SetTextBody(text string) {
	if text == "" {
		return
	}
	if c.textAsAlternative == true {
		c.email.AddAlternative(smail.TextPlain, text)
	} else {
		c.email.SetBody(smail.TextPlain, text)
	}
}

// SetHtmlBody sets the html email body, and flags the text body
// to be set as the alternative. Setting the html body with an
// empty string has no effect, other than to flag the text body
// as the primary.
func (c *Message) SetHtmlBody(html string) {
	c.textAsAlternative = false
	if html != "" {
		c.email.SetBody(smail.TextHTML, html)
		c.textAsAlternative = true
	}
}

func (c *Message) AddAttachment(filePath string, fileName string) {
	file := smail.File{
		FilePath: filePath,
		Name:     fileName,
		Inline:   false,
	}

	c.AddFile(&file)
}

func (c *Message) AddFile(file *smail.File) {
	c.email.Attach(file)
}

// Seal applies the envelope addressing to the email struct.
func (c *Message) Seal(envelope *Envelope) {
	envelope.Stamp(c.email)
}

// Message returns the underlying Email struct
func (c *Message) Message() *smail.Email {
	return c.email
}

// String returns the raw email text.
func (c *Message) String() string {
	return c.email.GetMessage()
}
