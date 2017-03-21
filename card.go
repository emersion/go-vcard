// Package vcard implements the vCard format, defined in RFC 6350.
package vcard

import (
	"strings"
)

func maybeGet(l []string, i int) string {
	if i < len(l) {
		return l[i]
	}
	return ""
}

// A Card is an address book entry.
type Card map[string][]*Field

// Get returns the first field of a card. If there is no such field, it returns
// nil.
func (c Card) Get(k string) *Field {
	fields := c[k]
	if len(fields) == 0 {
		return nil
	}
	return fields[0]
}

// A field contains a value and some parameters.
type Field struct {
	Value string
	Params map[string]string
	Group string
}

type Name struct {
	FamilyName string
	GivenName string
	AdditionalName string
	HonorificPrefix string
	HonorificSuffix string

	Params map[string]string
	Group string
}

func newName(field *Field) *Name {
	components := strings.Split(field.Value, ";")
	return &Name{
		maybeGet(components, 0),
		maybeGet(components, 1),
		maybeGet(components, 2),
		maybeGet(components, 3),
		maybeGet(components, 4),
		field.Params,
		field.Group,
	}
}

type Gender string

const (
	GenderFemale Gender = "F"
	GenderMale = "M"
	GenderOther = "O"
	GenderNone = "N"
	GenderUnknown = "U"
)
