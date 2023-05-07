package guild

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/vanng822/go-premailer/premailer"
	smail "github.com/xhit/go-simple-mail/v2"
	"strings"
)

// PongoScribe
type PongoScribe struct {
	highPriority    bool
	message         *Message
	subjectTemplate *pongo2.Template
	textTemplate    *pongo2.Template
	htmlTemplate    *pongo2.Template
	attachments     []*smail.File
	errors          []error
}

func (ps *PongoScribe) createTemplate(tmplStr string) *pongo2.Template {
	tmpl, err := pongo2.FromString(tmplStr)
	if err != nil {
		ps.addError(err)
	}
	return tmpl
}

func (ps *PongoScribe) SetPriority(isHigh bool) {
	ps.highPriority = isHigh
}

func (ps *PongoScribe) SetSubjectTemplate(subject string) {
	ps.subjectTemplate = ps.createTemplate(subject)
}

func (ps *PongoScribe) SetTextBodyTemplate(text string) {
	if emptyString(text) {
		return
	}
	ps.textTemplate = ps.createTemplate(text)
}

func (ps *PongoScribe) SetHtmlBodyTemplate(html string) {
	if emptyString(html) {
		return
	}
	ps.htmlTemplate = ps.createTemplate(html)
}

func (ps *PongoScribe) SetBodyTemplate(text string, html bool) {
	if html == true {
		ps.SetHtmlBodyTemplate(text)
	} else {
		ps.SetTextBodyTemplate(text)
	}
}

func (ps *PongoScribe) Include(filepath string, name ...string) {
	file := smail.File{FilePath: filepath}
	ps.attachments = append(ps.attachments, &file)
}

func (ps *PongoScribe) Open() (*Message, error) {
	if ps.message != nil {
		return nil, fmt.Errorf("Cannot start a new message while there is an existing one.")
	}
	ps.message = NewMessage()
	return ps.message, nil
}

func (ps *PongoScribe) Close() {
	ps.message = nil
	ps.errors = nil
}

func (ps *PongoScribe) addError(err error) {
	ps.errors = append(ps.errors, err)
}

func (ps *PongoScribe) HasErrors() bool {
	return len(ps.errors) > 0
}

func (ps *PongoScribe) GetErrors() error {
	errCount := len(ps.errors)
	var errDetails []string
	for _, err := range ps.errors {
		errDetails = append(errDetails, err.Error())
	}
	return fmt.Errorf("Scribe encountered %d error(s): %s", errCount, strings.Join(errDetails, ", "))
}

func (ps *PongoScribe) renderSubject(ctx pongo2.Context) string {
	if ps.subjectTemplate == nil {
		return ""
	}

	subject, err := ps.subjectTemplate.Execute(ctx)
	if err != nil {
		ps.addError(err)
		return ""
	}
	return subject
}

func (ps *PongoScribe) renderText(ctx pongo2.Context) string {
	if ps.textTemplate == nil {
		return ""
	}
	text, err := ps.textTemplate.Execute(ctx)
	if err != nil {
		ps.addError(err)
		return ""
	}
	return text
}

func (ps *PongoScribe) renderHtml(ctx pongo2.Context) string {
	if ps.htmlTemplate == nil {
		return ""
	}
	html, err := ps.htmlTemplate.Execute(ctx)
	if err != nil {
		ps.addError(err)
		return ""
	}

	premHtml, err := premailer.NewPremailerFromString(html, premailer.NewOptions())
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

func (ps *PongoScribe) Compose(ctx ...any) Scribe {
	if len(ctx) == 0 {
		ps.addError(fmt.Errorf("PongoScribe.Compose must receive a single Context argument."))
		return ps
	}

	var pctx pongo2.Context

	if len(ctx) == 1 {
		pctx = ctx[0].(pongo2.Context)
	} else {
		ps.addError(fmt.Errorf("PongoScribe.Compose only accepts a single Context argument."))
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

func (ps *PongoScribe) Seal(envelope *Envelope) Scribe {
	ps.message.Seal(envelope)
	return ps
}

func (ps *PongoScribe) Message() *smail.Email {
	if ps.HasErrors() {
		panic(ps.GetErrors())
	}
	return ps.message.Message()
}

func NewPongoScribe() *PongoScribe {
	return &PongoScribe{message: nil}
}

func MakePongoContext(srcCtx map[string]any) pongo2.Context {
	ctx := pongo2.Context{}
	for k, v := range srcCtx {
		ctx[k] = v
	}
	return ctx
}
