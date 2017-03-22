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
		Params: map[string]string{"PID": "1.1"},
	}},
	"N": []*Field{{Value: "Doe;J.;;;"}},
	"EMAIL": []*Field{{
		Value: "jdoe@example.com",
		Params: map[string]string{"PID": "1.1"},
	}},
	"CLIENTPIDMAP": []*Field{{Value: "1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0556"}},
}

var testCardHandmade = Card{
	"VERSION": []*Field{{Value: "4.0"}},
	"N": []*Field{{Value: "Bloggs;Joe;;;"}},
	"FN": []*Field{{Value: "Joe Bloggs"}},
	"EMAIL": []*Field{{
		Value: "me@joebloggs.com",
		Params: map[string]string{"TYPE": "home", "PREF": "1"},
	}},
	"TEL": []*Field{{
		Value: "tel:+44 20 1234 5678",
		Params: map[string]string{"TYPE": "\"cell,home\"", "PREF": "1"},
	}},
	"ADR": []*Field{{
		Value: ";;1 Trafalgar Square;London;;WC2N;United Kingdom",
		Params: map[string]string{"TYPE": "home", "PREF": "1"},
	}},
	"URL": []*Field{{
		Value: "http://joebloggs.com",
		Params: map[string]string{"TYPE": "home", "PREF": "1"},
	}},
	"IMPP": []*Field{{
		Value: "skype:joe.bloggs",
		Params: map[string]string{"TYPE": "home", "PREF": "1"},
	}},
	"X-SOCIALPROFILE": []*Field{{
		Value: "twitter:https://twitter.com/joebloggs",
		Params: map[string]string{"TYPE": "home", "PREF": "1"},
	}},
}

var testCardGoogle = Card{
	"VERSION": []*Field{{Value: "3.0"}},
	"N": []*Field{{Value: "Bloggs;Joe;;;"}},
	"FN": []*Field{{Value: "Joe Bloggs"}},
	"EMAIL": []*Field{{
		Value: "me@joebloggs.com",
		Params: map[string]string{"TYPE": "HOME"},
	}},
	"TEL": []*Field{{
		Value: "+44 20 1234 5678",
		Params: map[string]string{"TYPE": "CELL"},
	}},
	"ADR": []*Field{{
		Value: ";;1 Trafalgar Square;London;;WC2N;United Kingdom",
		Params: map[string]string{"TYPE": "HOME"},
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
		Params: map[string]string{"TYPE": "pref"},
	}},
	"TEL": []*Field{{
		Value: "+44 20 1234 5678",
		Params: map[string]string{"TYPE": "pref"},
	}},
	"ADR": []*Field{{
		Value: ";;1 Trafalgar Square;London;;WC2N;United Kingdom",
		Params: map[string]string{"TYPE": "pref"},
	}},
	"URL": []*Field{{
		Value: "http://joebloggs.com",
		Params: map[string]string{"TYPE": "pref"},
		Group: "item1",
	}},
	"X-ABLABEL": []*Field{
		{Value: "_$!<HomePage>!$_", Group: "item1"},
	},
	"IMPP": []*Field{{
		Value: "skype:joe.bloggs",
		Params: map[string]string{"X-SERVICE-TYPE": "Skype", "TYPE": "pref"},
	}},
	"X-SOCIALPROFILE": []*Field{{
		Value: "https://twitter.com/joebloggs",
		Params: map[string]string{"TYPE": "twitter"},
	}},
}

func TestCard(t *testing.T) {
	testCardFullName := &Field{
		Value: "J. Doe",
		Params: map[string]string{"PID": "1.1"},
	}

	if field := testCard.Get(FieldFormattedName); !reflect.DeepEqual(testCardFullName, field) {
		t.Errorf("Expected card FN field to be %+v but got %+v", testCardFullName, field)
	}
	if v := testCard.Value(FieldFormattedName); v != testCardFullName.Value {
		t.Errorf("Expected card FN field to be %q but got %q", testCardFullName.Value, v)
	}

	testCardName := &Name{
		FamilyName: "Doe",
		GivenName: "J.",
	}
	if name := testCard.Name(); !reflect.DeepEqual(testCardName, name) {
		t.Errorf("Expected card name to be %+v but got %+v", testCardName, name)
	}
}

func TestCard_Preferred(t *testing.T) {
	card := Card{
		"EMAIL": []*Field{
			{
				Value: "me@example.org",
				Params: map[string]string{"TYPE": "home"},
			},
			{
				Value: "me@example.com",
				Params: map[string]string{"TYPE": "work", "PREF": "1"},
			},
		},
	}

	if pref := card.Preferred(FieldEmail); pref != card["EMAIL"][1] {
		t.Errorf("Expected card preferred email to be %+v but got %+v", card["EMAIL"][1], pref)
	}
	if v := card.PreferredValue(FieldEmail); v != "me@example.com" {
		t.Errorf("Expected card preferred email to be %q but got %q", "me@example.com", v)
	}
}
