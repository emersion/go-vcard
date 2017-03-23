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

func TestFormatLine_withGroup(t *testing.T) {
	l := formatLine("FN", &Field{
		Value: "Akiyama Mio",
		Group: "item1",
	})

	expected := "item1.FN:Akiyama Mio"
	if l != expected {
		t.Errorf("Excpected formatted line with group to be %q, but got %q", expected, l)
	}
}

var testValue = []struct{
	v string
	formatted string
}{
	{"Hello World!", "Hello World!"},
	{"this is a single value, with a comma encoded", "this is a single value\\, with a comma encoded"},
	{"Mythical Manager\nHyjinx Software Division", "Mythical Manager\\nHyjinx Software Division"},
	{"aa\\\nbb", "aa\\\\\\nbb"},
}

func TestFormatValue(t *testing.T) {
	for _, test := range testValue {
		if formatted := formatValue(test.v); formatted != test.formatted {
			t.Errorf("formatValue(%q): expected %q, got %q", test.v, test.formatted, formatted)
		}
	}
}
