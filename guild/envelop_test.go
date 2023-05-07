package guild

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMailer_ParseEmailAddresses_WithEmptyValue(t *testing.T) {
	tst := assert.New(t)

	em := NewEnvelope()
	em.SetFromAddress("")

	tst.True(em.HasErrors())
}

type AddressData struct {
	Src          string
	Count        int
	FirstAddress string
	LastAddress  string
}

var data = []AddressData{
	{
		Src:          "myname@emaildomain.com",
		Count:        1,
		FirstAddress: "myname@emaildomain.com",
		LastAddress:  "myname@emaildomain.com",
	},
	{
		Src:          "myname@emaildomain.com,other@newname.ca",
		Count:        2,
		FirstAddress: "myname@emaildomain.com",
		LastAddress:  "other@newname.ca",
	},
}

func TestMailer_ParseEmailAddresses(t *testing.T) {
	tst := assert.New(t)
	em := NewEnvelope()

	for _, d := range data {
		em.AddToAddress(d.Src)
		tst.False(em.HasErrors())
		tst.Equal(d.Count, em.toAddresses.Cardinality())
		tst.Equal(d.FirstAddress, em.GetToAddresses()[0].Address)
	}
}
