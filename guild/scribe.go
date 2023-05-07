package guild

import (
	"fmt"
	smail "github.com/xhit/go-simple-mail/v2"
	"regexp"
)

var stripInterTagWhiteSpace *regexp.Regexp = regexp.MustCompile(">[\r\n\t ]+<")

func flattenHtml(html string) string {
	return stripInterTagWhiteSpace.ReplaceAllString(html, "><")
}

// Scribe creates messages.
type Scribe interface {
	SetPriority(bool)
	SetSubjectTemplate(string)
	SetTextBodyTemplate(string)
	SetHtmlBodyTemplate(string)
	SetBodyTemplate(string, bool)
	Include(string, ...string)
	Open() (*Message, error)
	Close()
	Seal(*Envelope) Scribe
	Compose(...any) Scribe
	Message() *smail.Email
	HasErrors() bool
	GetErrors() error
}

// SimpleTextScribe renders plain static text for the email body.
type SimpleTextScribe struct {
	highPriority bool
	message      *Message
	subject      string
	text         string
	html         string
	attachments  []*smail.File
}

func (sts *SimpleTextScribe) SetPriority(isHigh bool) {
	sts.highPriority = isHigh
}

func (sts *SimpleTextScribe) SetSubjectTemplate(subject string) {
	sts.subject = subject
}

func (sts *SimpleTextScribe) SetTextBodyTemplate(text string) {
	sts.text = text
}

func (sts *SimpleTextScribe) SetHtmlBodyTemplate(html string) {
	sts.html = html
}

func (sts *SimpleTextScribe) SetBodyTemplate(text string, html bool) {
	if html == true {
		sts.SetHtmlBodyTemplate(text)
	} else {
		sts.SetTextBodyTemplate(text)
	}
}

func (sts *SimpleTextScribe) Include(filepath string, name ...string) {
	file := smail.File{FilePath: filepath}
	sts.attachments = append(sts.attachments, &file)
}

func (sts *SimpleTextScribe) Open() (*Message, error) {
	if sts.message != nil {
		return nil, fmt.Errorf("Cannot start a new message while there is an existing one.")
	}
	sts.message = NewMessage()
	return sts.message, nil
}

func (sts *SimpleTextScribe) Close() {
	sts.message = nil
}

func (sts *SimpleTextScribe) HasErrors() bool {
	// TODO: reallY?
	return false
}

func (sts *SimpleTextScribe) GetErrors() error {
	return nil
}

func (sts *SimpleTextScribe) Compose(ctx ...any) Scribe {
	sts.message.SetPriority(sts.highPriority)
	sts.message.SetSubject(sts.subject)
	sts.message.SetTextBody(sts.text)
	sts.message.SetHtmlBody(sts.html)
	for _, attach := range sts.attachments {
		sts.message.AddFile(attach)
	}

	return sts
}

func (sts *SimpleTextScribe) Seal(envelope *Envelope) Scribe {
	sts.message.Seal(envelope)
	return sts
}

func (sts *SimpleTextScribe) Message() *smail.Email {
	return sts.message.Message()
}

func NewSimpleTextScribe() *SimpleTextScribe {
	return &SimpleTextScribe{message: nil}
}
