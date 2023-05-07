package guild

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateTestPongoScribe() *PongoScribe {
	scribe := NewPongoScribe()
	scribe.SetSubjectTemplate("Subject Is {{ subject }}")
	scribe.SetTextBodyTemplate("This is the {{ subject }} body.")
	return scribe
}

func CreateEnvelope() *Envelope {
	envelope := NewEnvelope()
	envelope.SetFromAddress("sender@email.com")
	envelope.AddToAddress("receiver@email.com")
	return envelope
}

func TestPongoScribe_Templates(t *testing.T) {
	tst := assert.New(t)

	scribe := CreateTestPongoScribe()

	tst.NotNil(scribe.subjectTemplate)
	tst.NotNil(scribe.textTemplate)
	tst.Nil(scribe.htmlTemplate)

	ctx := pongo2.Context{"subject": "TRY THIS"}

	result, err := scribe.subjectTemplate.Execute(ctx)
	tst.Nil(err)
	tst.Equal("Subject Is TRY THIS", result)

	result, err = scribe.textTemplate.Execute(ctx)
	tst.Nil(err)
	tst.Equal("This is the TRY THIS body.", result)

}

func TestPongoScribe_RenderTextOnlyEmail(t *testing.T) {
	tst := assert.New(t)

	scribe := CreateTestPongoScribe()
	envelope := CreateEnvelope()

	_, err := scribe.Open()
	tst.Nil(err)

	defer func() {
		scribe.Close()
	}()

	scribe.Compose(pongo2.Context{
		"subject": "TEST",
	})
	scribe.Seal(envelope)
	msg := scribe.message.String()

	tst.Contains(msg, "From: <sender@email.com>")
	tst.Contains(msg, "To: <receiver@email.com>")
	tst.Contains(msg, "Subject: Subject Is TEST")
	tst.Contains(msg, "This is the TEST body.")

}

func TestPongoScribe_RenderWithHtmlEmail(t *testing.T) {
	tst := assert.New(t)

	scribe := CreateTestPongoScribe()
	envelope := CreateEnvelope()

	scribe.SetHtmlBodyTemplate("<p>This is the {{ subject }} html body.</p>")

	_, err := scribe.Open()
	tst.Nil(err)

	defer func() {
		scribe.Close()
	}()

	scribe.Compose(pongo2.Context{
		"subject": "TEST",
	})
	scribe.Seal(envelope)
	msg := scribe.message.String()

	fmt.Println(msg)
	tst.Contains(msg, "From: <sender@email.com>")
	tst.Contains(msg, "To: <receiver@email.com>")
	tst.Contains(msg, "Subject: Subject Is TEST")
	tst.Contains(msg, "This is the TEST body.", "Missing text body.")
	tst.Contains(msg, "<p>This is the TEST html body.</p>", "Missing html body.")

}

func TestPongoScribe_RenderWithStaticTemplates(t *testing.T) {
	tst := assert.New(t)

	scribe := CreateTestPongoScribe()
	envelope := CreateEnvelope()

	scribe.SetSubjectTemplate("STATIC SUBJECT")
	scribe.SetTextBodyTemplate("STATIC TEXT BODY")
	scribe.SetHtmlBodyTemplate("<p>STATIC HTML BODY</p>")

	_, err := scribe.Open()
	tst.Nil(err)

	defer func() {
		scribe.Close()
	}()

	scribe.Compose(pongo2.Context{
		"subject": "TEST",
	})
	scribe.Seal(envelope)
	msg := scribe.message.String()

	fmt.Println(msg)
	tst.Contains(msg, "From: <sender@email.com>")
	tst.Contains(msg, "To: <receiver@email.com>")
	tst.Contains(msg, "Subject: STATIC SUBJECT", "Missing subject.")
	tst.Contains(msg, "STATIC TEXT BODY", "Missing text body.")
	tst.Contains(msg, "<p>STATIC HTML BODY</p>", "Missing html body.")

}

func TestPongoScribe_RenderHtmlWithStyles(t *testing.T) {
	tst := assert.New(t)

	scribe := CreateTestPongoScribe()
	envelope := CreateEnvelope()

	scribe.SetSubjectTemplate("STATIC SUBJECT")
	scribe.SetTextBodyTemplate("STATIC TEXT BODY")
	scribe.SetHtmlBodyTemplate(`<html>
  <head>
    <style>
      p {
        color: red;
     }
    </style>
  </head>
  <body>
    <p>STATIC HTML BODY</p>
  </body>
</html>`)

	_, err := scribe.Open()
	tst.Nil(err)

	defer func() {
		scribe.Close()
	}()

	scribe.Compose(pongo2.Context{
		"subject": "TEST",
	})
	scribe.Seal(envelope)
	mail := scribe.Message()
	msg := mail.GetMessage()

	fmt.Println(msg)

	tst.False(scribe.HasErrors())
	tst.Contains(msg, "From: <sender@email.com>")
	tst.Contains(msg, "To: <receiver@email.com>")
	tst.Contains(msg, "Subject: STATIC SUBJECT", "Missing subject.")
	tst.Contains(msg, "STATIC TEXT BODY", "Missing text body.")
	tst.Contains(msg, `<p style=3D"color:red">STATIC HTML BODY</p>`, "Missing html body.")

}
