package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/markgemmill/courier"
	"github.com/markgemmill/courier/params"
)

type SendCmd struct {
	// courier options
	Host     string `name:"host" short:"H" required:""`
	Port     int    `name:"port" short:"P" required:""`
	User     string `name:"user-name" short:"u" optional:""`
	Password string `name:"user-pwd" short:"p" optional:""`
	// envelop options
	SendFrom string   `name:"send-from" short:"F" required:""`
	ReplyTo  string   `name:"reply-to" optional:""`
	SendTo   []string `name:"send-to" short:"T"`
	SendCc   []string `name:"send-cc" short:"C"`
	// message options
	HighPriority bool              `name:"high-priority"`
	Subject      string            `name:"subject" short:"S"`
	Html         bool              `name:"html" optional:""`
	Template     string            `name:"template" short:"t" group:"templating" enum:"none,go,pongo" default:"none"`
	Params       map[string]string `name:"params" group:"templating"`
	Message      string            `arg:""`
	Attachment   []string          `name:"attach" short:"A" type:"existingfile"`
}

func (cmd *SendCmd) AfterApply() error {
	if cmd.Template != "none" && len(cmd.Params) == 0 {
		return fmt.Errorf("Template option '%s' requires at least one parameter value.", cmd.Template)
	}
	return nil
}

type CLI struct {
	Send *SendCmd `cmd:""`
}

func (cmd *SendCmd) Run(ctx *kong.Context) error {
	message := params.Parameters{
		CourierParams: params.CourierParams{
			Host:     cmd.Host,
			Port:     cmd.Port,
			User:     cmd.User,
			Password: cmd.Password,
		},
		EnvelopParams: params.EnvelopParams{
			SendFrom: cmd.SendFrom,
			ReplyTo:  cmd.ReplyTo,
			SendTo:   cmd.SendTo,
			SendCc:   cmd.SendCc,
		},
		MessageParams: params.MessageParams{
			HighPriority: cmd.HighPriority,
			Subject:      cmd.Subject,
			Attachments:  cmd.Attachment,
		},
	}

	message.TemplateType = cmd.Template
	params.SetMessage(&message, cmd.Message, cmd.Html)
	params.SetTemplateData(&message, cmd.Params)

	return courier.Deliver(message)
}

func main() {
	fmt.Println("Running courier test...")
	ctx := kong.Parse(&CLI{},
		kong.Name("courier"),
		kong.Description(""),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"version": "0.1.0",
		})
	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err, "")
}
