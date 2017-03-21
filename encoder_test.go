package vcard

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncoder(t *testing.T) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(testCard); err != nil {
		t.Fatal("Expected no error when formatting card, got:", err)
	}

	card, err := NewDecoder(&b).Decode()
	if err != nil {
		t.Fatal("Expected no error when parsing formatted card, got:", err)
	}

	if !reflect.DeepEqual(card, testCard) {
		t.Errorf("Invalid parsed card: expected \n%+v\n but got \n%+v", testCard, card)
	}
}
