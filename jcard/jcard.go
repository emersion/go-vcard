// Package jcard implements the JSON format for vCard.
//
// jCard is defined in RFC 7095.
package jcard

import (
	"encoding/json"
	"strings"

	"github.com/emersion/go-vcard"
)

type Card vcard.Card

var _ json.Marshaler = Card(nil)

func (card Card) MarshalJSON() ([]byte, error) {
	var props [][]interface{}

	// The "version" field must appear first
	for _, f := range card[vcard.FieldVersion] {
		props = append(props, rawProp(vcard.FieldVersion, f))
	}

	for k, fields := range card {
		if k == vcard.FieldVersion {
			continue
		}
		for _, f := range fields {
			props = append(props, rawProp(k, f))
		}
	}

	raw := [2]interface{}{"vcard", props}
	return json.Marshal(raw)
}

func rawProp(k string, f *vcard.Field) []interface{} {
	params := make(map[string]interface{})
	for k, values := range f.Params {
		if k == vcard.ParamValue {
			continue
		}

		if len(values) == 1 {
			params[k] = values[0]
		} else {
			params[k] = values
		}
	}

	if f.Group != "" {
		params["group"] = strings.ToLower(f.Group)
	}

	typ := "unknown"
	if value := f.Params.Get(vcard.ParamValue); value != "" {
		typ = value
	}

	// TODO: encode value according to type

	return []interface{}{strings.ToLower(k), params, typ, f.Value}
}
