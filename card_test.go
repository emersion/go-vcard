package vcard

import (
	"reflect"
	"testing"
)

var testCard = Card{
	"VERSION": []*Field{{Value: "4.0"}},
	"UID": []*Field{{Value: "urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b1"}},
	"FN": []*Field{{
		Value: "J. Doe",
		Params: Params{"PID": {"1.1"}},
	}},
	"N": []*Field{{Value: "Doe;J.;;;"}},
	"EMAIL": []*Field{{
		Value: "jdoe@example.com",
		Params: Params{"PID": {"1.1"}},
	}},
	"CLIENTPIDMAP": []*Field{{Value: "1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0556"}},
}

var testCardHandmade = Card{
	"VERSION": []*Field{{Value: "4.0"}},
	"N": []*Field{{Value: "Bloggs;Joe;;;"}},
	"FN": []*Field{{Value: "Joe Bloggs"}},
	"EMAIL": []*Field{{
		Value: "me@joebloggs.com",
		Params: Params{"TYPE": {"home"}, "PREF": {"1"}},
	}},
	"TEL": []*Field{{
		Value: "tel:+44 20 1234 5678",
		Params: Params{"TYPE": {"\"cell", "home\""}, "PREF": {"1"}},
	}},
	"ADR": []*Field{{
		Value: ";;1 Trafalgar Square;London;;WC2N;United Kingdom",
		Params: Params{"TYPE": {"home"}, "PREF": {"1"}},
	}},
	"URL": []*Field{{
		Value: "http://joebloggs.com",
		Params: Params{"TYPE": {"home"}, "PREF": {"1"}},
	}},
	"IMPP": []*Field{{
		Value: "skype:joe.bloggs",
		Params: Params{"TYPE": {"home"}, "PREF": {"1"}},
	}},
	"X-SOCIALPROFILE": []*Field{{
		Value: "twitter:https://twitter.com/joebloggs",
		Params: Params{"TYPE": {"home"}, "PREF": {"1"}},
	}},
}

var testCardGoogle = Card{
	"VERSION": []*Field{{Value: "3.0"}},
	"N": []*Field{{Value: "Bloggs;Joe;;;"}},
	"FN": []*Field{{Value: "Joe Bloggs"}},
	"EMAIL": []*Field{{
		Value: "me@joebloggs.com",
		Params: Params{"TYPE": {"INTERNET", "HOME"}},
	}},
	"TEL": []*Field{{
		Value: "+44 20 1234 5678",
		Params: Params{"TYPE": {"CELL"}},
	}},
	"ADR": []*Field{{
		Value: ";;1 Trafalgar Square;London;;WC2N;United Kingdom",
		Params: Params{"TYPE": {"HOME"}},
	}},
	"URL": []*Field{
		{Value: "http\\://joebloggs.com", Group: "item1"},
		{Value: "http\\://twitter.com/test", Group: "item2"},
	},
	"X-SKYPE": []*Field{{Value: "joe.bloggs"}},
	"X-ABLABEL": []*Field{
		{Value: "_$!<HomePage>!$_", Group: "item1"},
		{Value: "Twitter", Group: "item2"},
	},
}

var testCardApple = Card{
	"VERSION": []*Field{{Value: "3.0"}},
	"N": []*Field{{Value: "Bloggs;Joe;;;"}},
	"FN": []*Field{{Value: "Joe Bloggs"}},
	"EMAIL": []*Field{{
		Value: "me@joebloggs.com",
		Params: Params{"TYPE": {"INTERNET", "HOME", "pref"}},
	}},
	"TEL": []*Field{{
		Value: "+44 20 1234 5678",
		Params: Params{"TYPE": {"CELL", "VOICE", "pref"}},
	}},
	"ADR": []*Field{{
		Value: ";;1 Trafalgar Square;London;;WC2N;United Kingdom",
		Params: Params{"TYPE": {"HOME", "pref"}},
	}},
	"URL": []*Field{{
		Value: "http://joebloggs.com",
		Params: Params{"TYPE": {"pref"}},
		Group: "item1",
	}},
	"X-ABLABEL": []*Field{
		{Value: "_$!<HomePage>!$_", Group: "item1"},
	},
	"IMPP": []*Field{{
		Value: "skype:joe.bloggs",
		Params: Params{"X-SERVICE-TYPE": {"Skype"}, "TYPE": {"HOME", "pref"}},
	}},
	"X-SOCIALPROFILE": []*Field{{
		Value: "https://twitter.com/joebloggs",
		Params: Params{"TYPE": {"twitter"}},
	}},
}

func TestMaybeGet(t *testing.T) {
	l := []string{"a", "b", "c"}

	expected := []string{"a", "b", "c", "", ""}
	for i, exp := range expected {
		if v := maybeGet(l, i); v != exp {
			t.Errorf("maybeGet(l, %v): expected %q but got %q", i, exp, v)
		}
	}
}

func TestCard(t *testing.T) {
	testCardFullName := testCard["FN"][0]
	if field := testCard.Get(FieldFormattedName); !reflect.DeepEqual(testCardFullName, field) {
		t.Errorf("Expected card FN field to be %+v but got %+v", testCardFullName, field)
	}
	if v := testCard.Value(FieldFormattedName); v != testCardFullName.Value {
		t.Errorf("Expected card FN field to be %q but got %q", testCardFullName.Value, v)
	}

	testCardName := &Name{
		Field: testCard["N"][0],
		FamilyName: "Doe",
		GivenName: "J.",
	}
	if name := testCard.Name(); !reflect.DeepEqual(testCardName, name) {
		t.Errorf("Expected card name to be %+v but got %+v", testCardName, name)
	}
}

func TestCard_Preferred(t *testing.T) {
	cards := []Card{
		Card{
			"EMAIL": []*Field{
				{Value: "me@example.org", Params: Params{"TYPE": {"home"}}},
				{Value: "me@example.com", Params: Params{"TYPE": {"work"}, "PREF": {"1"}}},
			},
		},
		// Apple Contacts
		Card{
			"EMAIL": []*Field{
				{Value: "me@example.org", Params: Params{"TYPE": {"home"}}},
				{Value: "me@example.com", Params: Params{"TYPE": {"work", "pref"}}},
			},
		},
	}

	for _, card := range cards {
		if pref := card.Preferred(FieldEmail); pref != card["EMAIL"][1] {
			t.Errorf("Expected card preferred email to be %+v but got %+v", card["EMAIL"][1], pref)
		}
		if v := card.PreferredValue(FieldEmail); v != "me@example.com" {
			t.Errorf("Expected card preferred email to be %q but got %q", "me@example.com", v)
		}
	}
}
