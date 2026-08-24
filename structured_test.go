package vcard

import (
	"bytes"
	"strings"
	"testing"
)

func roundtripCard(t *testing.T, c Card) Card {
	t.Helper()
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(c); err != nil {
		t.Fatal(err)
	}
	out, err := NewDecoder(strings.NewReader(buf.String())).Decode()
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func versionedCard() Card {
	c := make(Card)
	c.SetValue(FieldVersion, "4.0")
	return c
}

// A literal ';' inside a structured component must survive a round trip: per
// RFC 6350 section 3.4 a ';' separates components only when unescaped.
func TestStructuredName_Roundtrip(t *testing.T) {
	tests := []Name{
		{FamilyName: "a;b", GivenName: "c"},
		{FamilyName: "de Groot", GivenName: "Rene;", AdditionalName: ";x;"},
		{FamilyName: `back\slash`, GivenName: "co,mma"},
		{FamilyName: `ends-with\`, GivenName: "b"},
		{FamilyName: `lit\;eral`, GivenName: "c"},
		{FamilyName: "line\nbreak"},
		{FamilyName: "Doe", GivenName: "J."},
	}
	for _, want := range tests {
		c := versionedCard()
		c.SetName(&want)
		got := roundtripCard(t, c).Name()
		if got == nil {
			t.Fatalf("%+v: Name() nil after round trip", want)
		}
		if got.FamilyName != want.FamilyName || got.GivenName != want.GivenName ||
			got.AdditionalName != want.AdditionalName {
			t.Errorf("round trip: got (%q,%q,%q) want (%q,%q,%q)",
				got.FamilyName, got.GivenName, got.AdditionalName,
				want.FamilyName, want.GivenName, want.AdditionalName)
		}
	}
}

func TestStructuredAddress_Roundtrip(t *testing.T) {
	want := &Address{
		StreetAddress: "12 Main St; Apt 3",
		Locality:      "A,B",
		Region:        `C\D`,
		Country:       "x;y;z",
	}
	c := versionedCard()
	c.SetAddress(want)
	got := roundtripCard(t, c).Address()
	if got == nil {
		t.Fatal("Address() nil after round trip")
	}
	if got.StreetAddress != want.StreetAddress || got.Locality != want.Locality ||
		got.Region != want.Region || got.Country != want.Country {
		t.Errorf("round trip: got %+v want %+v", got, want)
	}
}

// The encoded form escapes a literal ';' as '\;' and keeps separators literal.
func TestStructuredName_EncodedForm(t *testing.T) {
	c := versionedCard()
	c.SetName(&Name{FamilyName: "a;b", GivenName: "c"})
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(c); err != nil {
		t.Fatal(err)
	}
	const want = "N:a\\;b;c;;;"
	if !strings.Contains(buf.String(), want) {
		t.Errorf("encoded N line: want %q in\n%s", want, buf.String())
	}
}

// Only property values are escaped, not parameter values (RFC 6350 section 3.4):
// a ';' in a parameter must not gain a structured-value backslash escape.
func TestStructuredParam_Unaffected(t *testing.T) {
	c := versionedCard()
	c.Set("TEL", &Field{Value: "123", Params: Params{"TYPE": {"a;b"}}})
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(c); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), `\;`) {
		t.Errorf("parameter value was structurally escaped:\n%s", buf.String())
	}
}
