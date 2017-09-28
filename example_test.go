package vcard_test

import (
	"io"
	"log"
	"os"

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
