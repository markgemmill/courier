package courier

import (
	"github.com/markgemmill/courier/guild"
	"github.com/markgemmill/courier/params"
)

// Deliver is the only function that is needed to send an email.
func Deliver(p params.Parameters) error {
	courier := guild.NewCourier(
		p.Host,
		p.Port,
		p.User,
		p.Password,
	)

	var scribe guild.Scribe

	if p.TemplateType == "pongo" {
		scribe = guild.NewPongoScribe()
	} else if p.TemplateType == "go" {
		scribe = guild.NewTemplateScribe()
	} else {
		scribe = guild.NewSimpleTextScribe()
	}

	// set base template / text
	// this is content that will be reused
	scribe.SetPriority(p.HighPriority)
	scribe.SetSubjectTemplate(p.Subject)
	scribe.SetTextBodyTemplate(guild.LoadTemplateString(p.TextMessage))
	scribe.SetHtmlBodyTemplate(guild.LoadTemplateString(p.HtmlMessage))

	for _, filepath := range p.Attachments {
		scribe.Include(filepath, "")
	}

	// start a new correspondence session
	_, err := scribe.Open()
	if err != nil {
		return err
	}
	defer func() {
		scribe.Close()
	}()

	envelope := guild.NewEnvelope()
	envelope.SetFromAddress(p.SendFrom)
	envelope.SetReplyToAddress(p.ReplyTo)
	envelope.AddToAddresses(p.SendTo)
	envelope.AddCcAddresses(p.SendCc)
	envelope.AddCcAddresses(p.SendBcc)

	if envelope.HasErrors() {
		return envelope.GetErrors()
	}

	if p.TemplateType == "pongo" {
		scribe.Compose(guild.MakePongoContext(p.TemplateData))
	} else if p.TemplateType == "go" {
		scribe.Compose(p.TemplateData)
	} else {
		scribe.Compose()
	}

	scribe.Seal(envelope)

	if scribe.HasErrors() {
		return scribe.GetErrors()
	}

	err = courier.Deliver(scribe.Message())
	if err != nil {
		return err
	}

	return nil

}
