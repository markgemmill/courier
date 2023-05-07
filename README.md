# courier 

courier is an email library with a different metaphor and for simplicity. 

Basic usage:

```go
package main


import (
	"github.com/markgemmill/courier"
	"github.com/markgemmill/courier/params"
)

func main () {
	
    message := params.Parameters{
        CourierParams: params.CourierParams{
            Host:     "mail.smtphost.com",
            Port:     2525,
            User:     "smtpUserId",
            Password: "smtpUserPwd",
        },
        EnvelopParams: params.EnvelopParams{
            SendFrom: "sender@senderhost.org",
            ReplyTo:  "reply@senderhost.org",
            SendTo:   []string{"getit@emailtown.com"},
            SendCc:   []string{"also@emailtown.com"},
        },
        MessageParams: params.MessageParams{
            HighPriority: true,
			TemplateType: "pongo",
            Subject:      "A Subject Worth Sending",
            Attachments:  []string{"/pth/to/attachment.doc"},
        },
    }

    params.SetMessage(&message, "<h1>Hello {{ name }}</h1>", true)
    params.SetTemplateData(&message, map[string]string{"name": "Courier!"})

    err := courier.Deliver(message)

    if err != nil {
        panic(err)
    }

}
```

For more verbose usage:

```go
package main

import (
    "github.com/markgemmill/courier/guild"
    "github.com/markgemmill/courier/params"
)


func main() {
        
    // create the envelope
    envelope := guild.NewEnvelope()
    envelope.SetFromAddress("henry@england.com")
    envelope.AddToAddress("charels@france.com")
    envelope.AddCcAddress("walsingham@england.com")

    if envelope.HasErrors() {
        panic(envelope.GetErrors())
    }

    // create a scribe to define the message content and format
    scribe := guild.NewSimpleTextScribe() 
	scribe.SetPriority(true)
    scribe.SetSubjectTemplate("Give Up, or Else")
    scribe.SetTextBodyTemplate("It's mine and I'll come get it if I have to!")


    // start the scribe writing the message
    _, err := scribe.Open()
	if err != nil {
		panic(err)
    }

    scribe.Compose().Seal(envelope)

    if scribe.HasErrors() {
        panic(scribe.GetErrors())
    }

    courier := guild.NewCourier("smtpservice.org", 2525, "royalc", "seekrit")

    err = courier.Deliver(scribe.Message())
	
	if err != nil {
		panic(err)
    }
}
```

courier is built on top of the following libraries:

* [go-simple-mail](https://github.com/xhit/go-simple-mail/v2)

* [email-verifier](https://github.com/AfterShip/email-verifier)

* [go-email-normalizer](https://github.com/dimuska139/go-email-normalizer)

* [premailer](https://github.com/vanng822/go-premailer/premailer)

