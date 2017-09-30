package vcard_test

import (
	"io"
	"log"
	"os"
	"bufio"

	"github.com/emersion/go-vcard"
)

func ExampleNewDecoder() {
	f, err := os.Open("cards.vcf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := vcard.NewDecoder(f)
	for {
		card, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		log.Println(card.PreferredValue(vcard.FieldFormattedName))
	}
}

func ExampleNewEncoder() {
	destFile, err := os.Create("new_cards.vcf")
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	data := [][]string {
		{"John", "Webber", "Maxwell", "(+1) 199 8714",},
		{"Donald", "", "Ron", "(+44) 421 8913"},
		{"Eric", "E.", "Peter", "(+37) 221 9903",},
	}

	var (
		// card is a map of strings to []*vcard.Field objects
		card vcard.Card

		// buffer the writes to save disk I/O time
		writeBuffer = bufio.NewWriter(destFile)

		// destination where the vcard will be encoded to
		writer = vcard.NewEncoder(writeBuffer)
	)

	for _, entry := range data {
		contact := NewContact(entry)
		card[vcard.FieldFormattedName] = contact.ReadFormattedName()
		card[vcard.FieldName] = contact.ReadName()
		card[vcard.FieldTelephone] = contact.ReadTelephone()
		vcard.ToV4(card)
		err := writer.Encode(card)
		if err != nil {
			log.Fatal(err)
		}

	}
	err = writeBuffer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// encoding a vcard can be done using the following method

// Contact is our model for contact data
type Contact struct {
	FirstName  string
	MiddleName string
	LastName   string
	Telephone  string
}

// ReadName formats a Contact object into a []*vcard.Field object containing the
// name of the contact with its components
func (c *contact.Contact) ReadName() []*vcard.Field {
	return []*vcard.Field{
		&vcard.Field{
			Value: c.FirstName + ";" + c.LastName + ";" + c.MiddleName,
			Params: map[string][]string{
				vcard.ParamSortAs: []string{
					c.FirstName + " " + c.LastName,
				},
			},
		},
	}
}

// ReadFormattedName formats a Contact object into a []*vcard.Field object containing
// the formatted name of the contact
func (c *contact.Contact) ReadFormattedName() []*vcard.Field {
	return []*vcard.Field{
		&vcard.Field{
			Value: c.FirstName + " " + c.LastName + " " + c.MiddleName,
		},
	}
}

// ReadTelephone formats a telephone number string into a vcard.Field object
func (c *contact.Contact) ReadTelephone() []*vcard.Field {
	return []*vcard.Field{
		&vcard.Field{
			Value: c.Telephone,
		},
	}
}

// NewContact constructs a Contact object
func NewContact(strslc [4]string) *Contact {
	fn, mn, ln, tel := strslc[0], strslc[1], strslc[2], strslc[3]
	return &Contact{
		FirstName:  fn,
		MiddleName: mn,
		LastName:   ln,
		Telephone:  tel,
	}
}

