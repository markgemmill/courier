package guild

import (
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	smail "github.com/xhit/go-simple-mail/v2"
	"strings"
	"text/template"
)

// TemplateScribe
type TemplateScribe struct {
	highPriority    bool
	message         *Message
	subjectTemplate *template.Template
	textTemplate    *template.Template
	htmlTemplate    *template.Template
	attachments     []*smail.File
	errors          []error
}

func (ps *TemplateScribe) createTemplate(name, tmplStr string) *template.Template {
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		ps.addError(err)
	}
	return tmpl
}

func (ps *TemplateScribe) SetPriority(isHigh bool) {
	ps.highPriority = isHigh
}

func (ps *TemplateScribe) SetSubjectTemplate(subject string) {
	ps.subjectTemplate = ps.createTemplate("subject", subject)
}

func (ps *TemplateScribe) SetTextBodyTemplate(text string) {
	if emptyString(text) {
		return
	}
	ps.textTemplate = ps.createTemplate("text", text)
}

func (ps *TemplateScribe) SetHtmlBodyTemplate(html string) {
	if emptyString(html) {
		return
	}
	ps.htmlTemplate = ps.createTemplate("html", html)
}

func (ps *TemplateScribe) SetBodyTemplate(text string, html bool) {
	if html {
		ps.SetHtmlBodyTemplate(text)
	} else {
		ps.SetTextBodyTemplate(text)
	}
}

func (ps *TemplateScribe) Include(filepath string, name ...string) {
	file := smail.File{FilePath: filepath}
	ps.attachments = append(ps.attachments, &file)
}

func (ps *TemplateScribe) Open() (*Message, error) {
	if ps.message != nil {
		return nil, fmt.Errorf("Cannot start a new message while there is an existing one.")
	}
	ps.message = NewMessage()
	return ps.message, nil
}

func (ps *TemplateScribe) Close() {
	ps.message = nil
	ps.errors = nil
}

func (ps *TemplateScribe) addError(err error) {
	ps.errors = append(ps.errors, err)
}

func (ps *TemplateScribe) HasErrors() bool {
	return len(ps.errors) > 0
}

func (ps *TemplateScribe) GetErrors() error {
	errCount := len(ps.errors)
	var errDetails []string
	for _, err := range ps.errors {
		errDetails = append(errDetails, err.Error())
	}
	return fmt.Errorf("Scribe encountered %d error(s): %s", errCount, strings.Join(errDetails, ", "))
}

func (ps *TemplateScribe) renderSubject(ctx any) string {
	if ps.subjectTemplate == nil {
		return ""
	}
	subject := strings.Builder{}
	err := ps.subjectTemplate.Execute(&subject, ctx)
	if err != nil {
		ps.addError(err)
		return ""
	}
	return subject.String()
}

func (ps *TemplateScribe) renderText(ctx any) string {
	if ps.textTemplate == nil {
		return ""
	}
	text := strings.Builder{}
	err := ps.textTemplate.Execute(&text, ctx)
	if err != nil {
		ps.addError(err)
		return ""
	}
	return text.String()
}

func (ps *TemplateScribe) renderHtml(ctx any) string {
	if ps.htmlTemplate == nil {
		return ""
	}
	html := strings.Builder{}
	err := ps.htmlTemplate.Execute(&html, ctx)
	if err != nil {
		ps.addError(err)
		return ""
	}

	premHtml, err := premailer.NewPremailerFromString(html.String(), premailer.NewOptions())
	if err != nil {
		ps.addError(err)
		return ""
	}

	renderedHtml, err := premHtml.Transform()
	if err != nil {
		ps.addError(err)
		return ""
	}
	return flattenHtml(renderedHtml)
}

func (ps *TemplateScribe) Compose(ctx ...any) Scribe {
	if len(ctx) == 0 {
		ps.addError(fmt.Errorf("TemplateScribe.Compose must receive a single Context argument."))
		return ps
	}

	var pctx any

	if len(ctx) == 1 {
		pctx = ctx[0]
	} else {
		ps.addError(fmt.Errorf("TemplateScribe.Compose only accepts a single Context argument."))
		return ps
	}

	subject := ps.renderSubject(pctx)
	html := ps.renderHtml(pctx)
	text := ps.renderText(pctx)

	ps.message.SetPriority(ps.highPriority)
	ps.message.SetSubject(subject)
	ps.message.SetHtmlBody(html)
	ps.message.SetTextBody(text)

	for _, attach := range ps.attachments {
		ps.message.AddFile(attach)
	}

	return ps
}

func (ps *TemplateScribe) Seal(envelope *Envelope) Scribe {
	ps.message.Seal(envelope)
	return ps
}

func (ps *TemplateScribe) Message() *smail.Email {
	if ps.HasErrors() {
		panic(ps.GetErrors())
	}
	return ps.message.Message()
}

func NewTemplateScribe() *TemplateScribe {
	return &TemplateScribe{message: nil}
}

func MakeTemplateContext(srcCtx map[string]any) map[string]any {
	return srcCtx
}
