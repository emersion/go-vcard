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

func TestCard(t *testing.T) {
	testCardFullName := &Field{
		Value: "J. Doe",
		Params: map[string]string{"PID": "1.1"},
	}

	if v, params := testCard.Get("FN"); v != "J. Doe" {
		t.Errorf("Expected card FN field value to be %q but got %q", testCardFullName.Value, v)
	} else if !reflect.DeepEqual(testCardFullName.Params, params) {
		t.Errorf("Expected card FN field params to be %v but got %v", testCardFullName.Params, params)
	}
}
