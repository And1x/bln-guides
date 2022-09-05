package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// New initialzes a new Form struct with specific data as param
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks that form is not empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank!")
		}
	}
}

// MaxLength checks that form input is not too long
func (f *Form) MaxLength(field string, max int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > max {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (max. %d chars)", max))
	}
}

// MinLength checks that form input is not to short eg. password
func (f *Form) MinLength(field string, min int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < min {
		f.Errors.Add(field, fmt.Sprintf("This field is too shor (min. %d chars)", min))
	}
}

func (f *Form) ValidMail(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if value == "" {
			return
		}

		_, err := mail.ParseAddress(value)
		if err != nil {
			f.Errors.Add(field, "This Address is not valid")
		}
	}
}

// Valid returns if posted Form is valid or not
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
