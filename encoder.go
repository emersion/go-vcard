package vcard

import (
	"io"
	"strings"
)

// An Encoder formats cards.
type Encoder struct {
	w io.Writer
}

// NewEncoder creates a new Encoder that writes cards to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode formats a card.
func (enc *Encoder) Encode(c Card) error {
	begin := "BEGIN:VCARD\r\n"
	if _, err := io.WriteString(enc.w, begin); err != nil {
		return err
	}

	version := c.Get("VERSION")
	if _, err := io.WriteString(enc.w, formatLine("VERSION", version.Value, version.Params)+"\r\n"); err != nil {
		return err
	}

	for k, fields := range c {
		if strings.EqualFold(k, "VERSION") {
			continue
		}
		for _, f := range fields {
			if _, err := io.WriteString(enc.w, formatLine(k, f.Value, f.Params)+"\r\n"); err != nil {
				return err
			}
		}
	}

	end := "END:VCARD\r\n"
	_, err := io.WriteString(enc.w, end)
	return err
}

func formatLine(k, v string, params map[string]string) string {
	kparams := make([]string, 1 + len(params))
	kparams[0] = k
	i := 1
	for pk, pv := range params {
		kparams[i] = formatParam(pk, pv)
		i++
	}

	return strings.Join(kparams, ";")+":"+v
}

func formatParam(k, v string) string {
	return k+"="+v
}
