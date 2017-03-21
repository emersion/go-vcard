package vcard

import (
	"strings"
	"reflect"
	"testing"
)

var testCardString = `BEGIN:VCARD
VERSION:4.0
UID:urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b1
FN;PID=1.1:J. Doe
N:Doe;J.;;;
EMAIL;PID=1.1:jdoe@example.com
CLIENTPIDMAP:1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0556
END:VCARD`

func TestDecoder(t *testing.T) {
	r := strings.NewReader(testCardString)
	dec := NewDecoder(r)
	card, err := dec.Decode()
	if err != nil {
		t.Fatal("Expected no error when decoding card, got:", err)
	}

	if !reflect.DeepEqual(card, testCard) {
		t.Errorf("Invalid parsed card: expected \n%+v\n but got \n%+v", testCard, card)
	}
}
