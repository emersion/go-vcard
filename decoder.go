package vcard

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// A Decoder parses cards.
type Decoder struct {
	r *bufio.Reader
}

// NewDecoder creates a new Decoder reading cards from an io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (dec *Decoder) readLine() (string, error) {
	l, err := dec.r.ReadString('\n')
	if err != nil {
		return l, err
	}
	l = strings.TrimRight(l, "\r\n")

	for {
		next, err := dec.r.Peek(1)
		if err == io.EOF {
			break
		} else if err != nil {
			return l, err
		}

		if ch := next[0]; ch != ' ' && ch != '\t' {
			break
		}

		if _, err := dec.r.Discard(1); err != nil {
			return l, err
		}

		folded, err := dec.r.ReadString('\n')
		if err != nil {
			return l, err
		}
		l += strings.TrimRight(folded, "\r\n")
	}

	return l, nil
}

// Decode parses a single card.
func (dec *Decoder) Decode() (Card, error) {
	card := make(Card)

	hasHeader := false
	for {
		l, err := dec.readLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return card, err
		}

		k, f, err := parseLine(l)
		if err != nil {
			continue
		}

		if !hasHeader {
			if k == "BEGIN" {
				if strings.ToUpper(f.Value) != "VCARD" {
					return card, errors.New("vcard: invalid BEGIN value")
				}
				hasHeader = true
				continue
			} else {
				return card, errors.New("vcard: no BEGIN field found")
			}
		} else if k == "END" {
			if strings.ToUpper(f.Value) != "VCARD" {
				return card, errors.New("vcard: invalid END value")
			}
			break
		}

		card[k] = append(card[k], f)
	}

	return card, nil
}

func parseLine(l string) (key string, field *Field, err error) {
	kv := strings.SplitN(l, ":", 2)
	if len(kv) < 2 {
		return l, nil, errors.New("vcard: invalid field")
	}

	field = new(Field)
	field.Value = parseValue(kv[1])
	kparams := strings.Split(kv[0], ";")

	key, field.Group = parseKey(kparams[0])

	if len(kparams) > 1 {
		field.Params = make(Params)
		for i := 1; i < len(kparams); i++ {
			pk, pvs := parseParam(kparams[i])
			for i, pv := range pvs {
				pvs[i] = parseValue(pv)
			}
			field.Params[pk] = append(field.Params[pk], pvs...)
		}
	}

	return
}

func parseKey(s string) (key, group string) {
	parts := strings.SplitN(s, ".", 2)
	if len(parts) < 2 {
		key = s
	} else {
		group = parts[0]
		key = parts[1]
	}
	key = strings.ToUpper(key)
	return
}

func parseParam(s string) (k string, vs []string) {
	kv := strings.SplitN(s, "=", 2)
	if len(kv) < 2 {
		return s, nil
	}
	return strings.ToUpper(kv[0]), strings.Split(kv[1], ",")
}

var valueParser = strings.NewReplacer("\\\\", "\\", "\\n", "\n", "\\,", ",")

func parseValue(v string) string {
	return valueParser.Replace(v)
}
