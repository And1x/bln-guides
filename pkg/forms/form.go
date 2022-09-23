package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// todo: Rewrite error field with better messages

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

// IsPositiveNumber checks that input is a number that is > 0
func (f *Form) IsPositiveNumber(field string) {
	value := f.Get(field)
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(field, "This field got not a number")
		return
	}
	if valueInt <= 0 {
		f.Errors.Add(field, "This field needs number > 0")
		return
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

// ValidMail check if a mail is valid
func (f *Form) ValidMail(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

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
