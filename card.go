// Package vcard implements the vCard format, defined in RFC 6350.
package vcard

// A Card is an address book entry.
type Card map[string][]*Field

// Get returns the first field of a card. If there is no such field, it returns
// zero values.
func (c Card) Get(k string) (v string, params map[string]string) {
	fields := c[k]
	if len(fields) == 0 {
		return
	}
	return fields[0].Value, fields[0].Params
}

// A field contains a value and some parameters.
type Field struct {
	Value string
	Params map[string]string
}
