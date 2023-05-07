package guild

import (
	smail "github.com/xhit/go-simple-mail/v2"
	"time"
)

// Courier defines an smtp server/client responsible for
// "delivering" messages.
type Courier struct {
	server *smail.SMTPServer
}

func NewCourier(host string, port int, user, password string) *Courier {
	_courier := Courier{}

	server := smail.NewSMTPClient()
	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	// TODO: handle different encryptions
	server.Encryption = smail.EncryptionNone

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	_courier.server = server

	return &_courier
}

// Deliver provides a one-off connection and deliver of an email.
func (cr *Courier) Deliver(msg *smail.Email) error {

	client, err := cr.server.Connect()
	if err != nil {
		return err
	}

	defer func() {
		_ = client.Close()
	}()

	err = msg.Send(client)
	if err != nil {
		return err
	}
	return nil

}
