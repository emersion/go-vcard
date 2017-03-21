package vcard

type Card map[string][]*Field

func (c Card) Get(k string) (v string, params map[string]string) {
	fields := c[k]
	if len(fields) == 0 {
		return
	}
	return fields[0].Value, fields[0].Params
}

type Field struct {
	Value string
	Params map[string]string
}
