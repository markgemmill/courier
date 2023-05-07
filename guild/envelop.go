package guild

import (
	"fmt"
	emailverifier "github.com/AfterShip/email-verifier"
	mapset "github.com/deckarep/golang-set/v2"
	enormalizer "github.com/dimuska139/go-email-normalizer"
	smail "github.com/xhit/go-simple-mail/v2"
	"net/mail"
	"strings"
)

type AddressType int

const (
	FromAddress AddressType = iota
	ReplyToAddress
	ToAddress
	CcAddress
	BccAddress
)

var EmptyAddress = mail.Address{}

type Envelope struct {
	verifier       *emailverifier.Verifier
	normalizer     *enormalizer.Normalizer
	FromAddress    mail.Address
	ReplyToAddress mail.Address
	toAddresses    mapset.Set[mail.Address]
	ccAddresses    mapset.Set[mail.Address]
	bccAddresses   mapset.Set[mail.Address]
	errors         []error
}

func NewEnvelope() *Envelope {
	mgr := Envelope{
		verifier:     emailverifier.NewVerifier(),
		normalizer:   enormalizer.NewNormalizer(),
		toAddresses:  mapset.NewSet[mail.Address](),
		ccAddresses:  mapset.NewSet[mail.Address](),
		bccAddresses: mapset.NewSet[mail.Address](),
	}

	return &mgr
}

func (em *Envelope) addError(err error) {
	em.errors = append(em.errors, err)
}

func (em *Envelope) HasErrors() bool {
	return len(em.errors) > 0
}

func (em *Envelope) GetErrors() error {
	errCount := len(em.errors)
	var errDetails []string
	for _, err := range em.errors {
		errDetails = append(errDetails, err.Error())
	}
	return fmt.Errorf("%d envelope error(s): %s", errCount, strings.Join(errDetails, ", "))
}

// AcceptAddressString parses, validates and stores the given addresses. The idea here is
// it can be used to pull in email addresses from several locations and formats.
func (em *Envelope) acceptAddressString(addressString string, addressType AddressType) {
	addressString = strings.TrimSpace(addressString)
	if addressString == "" {
		em.addError(fmt.Errorf("empty string was provided to ParseEmailAddresses"))
	}

	addresses, err := mail.ParseAddressList(addressString)
	if err != nil {
		em.addError(err)
	}

	if addressType == FromAddress && len(addresses) > 1 {
		em.addError(fmt.Errorf("there can only be one from email address"))
	}

	if addressType == ReplyToAddress && len(addresses) > 1 {
		em.addError(fmt.Errorf("there can only be one reply to email address"))
	}

	// re-verify email address is valid
	for _, address := range addresses {

		result := em.verifier.ParseAddress(address.Address)
		if !result.Valid {
			em.addError(fmt.Errorf("'%s' is an invalid email address", address.Address))
		}

		// normalize the address string
		address.Address = em.normalizer.Normalize(address.Address)

		switch addressType {
		case ToAddress:
			em.toAddresses.Add(*address)
		case CcAddress:
			em.ccAddresses.Add(*address)
		case BccAddress:
			em.bccAddresses.Add(*address)
		case FromAddress:
			em.FromAddress = *address
		case ReplyToAddress:
			em.ReplyToAddress = *address
		}
	}

}

func (em *Envelope) acceptAddresses(addresses []string, addressType AddressType) {
	for _, adddress := range addresses {
		em.acceptAddressString(adddress, addressType)
	}
}

func (em *Envelope) SetFromAddress(address string) {
	em.acceptAddressString(address, FromAddress)
}

func (em *Envelope) SetReplyToAddress(address string) {
	// for reply to - we just ignore an empty string
	// this is for convenience purposes only
	if strings.TrimSpace(address) == "" {
		return
	}
	em.acceptAddressString(address, ReplyToAddress)
}

func (em *Envelope) AddToAddress(address string) {
	em.acceptAddressString(address, ToAddress)
}

func (em *Envelope) AddToAddresses(addresses []string) {
	em.acceptAddresses(addresses, ToAddress)
}

func (em *Envelope) AddCcAddress(address string) {
	em.acceptAddressString(address, CcAddress)
}

func (em *Envelope) AddCcAddresses(addresses []string) {
	em.acceptAddresses(addresses, CcAddress)
}

func (em *Envelope) AddBccAddress(address string) {
	em.acceptAddressString(address, BccAddress)
}

func (em *Envelope) AddBccAddresses(addresses []string) {
	em.acceptAddresses(addresses, BccAddress)
}

func (em *Envelope) GetToAddresses() []mail.Address {
	return em.toAddresses.ToSlice()
}

func (em *Envelope) GetCcAddresses() []mail.Address {
	return em.ccAddresses.ToSlice()
}

func (em *Envelope) GetBccAddresses() []mail.Address {
	return em.bccAddresses.ToSlice()
}

func (em *Envelope) Stamp(msg *smail.Email) {
	msg.SetFrom(em.FromAddress.String())

	if em.ReplyToAddress != EmptyAddress {
		msg.SetReplyTo(em.ReplyToAddress.String())
	}

	for _, addr := range em.GetToAddresses() {
		msg.AddTo(addr.String())
	}

	for _, addr := range em.GetCcAddresses() {
		msg.AddCc(addr.String())
	}

	for _, addr := range em.GetBccAddresses() {
		msg.AddBcc(addr.String())
	}
}
