package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

func (f *Form) Match(field1 string, field2 string) {
	value1 := f.Get(field1)
	value2 := f.Get(field2)
	if value1 == "" || value2 == "" {
		return
	}

	if value1 != value2 {
		f.Errors.Add(field1, fmt.Sprintf("%s value is different from %s value", field1, field2))
		f.Errors.Add(field2, fmt.Sprintf("%s value is different from %s value", field1, field2))
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
