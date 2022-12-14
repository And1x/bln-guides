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
			f.Errors.Add(field, "Required field! Please enter something!")
		}
	}
}

// IsPositiveNumber checks that input is a number that is > 0
func (f *Form) IsPositiveNumber(field string) {
	value := f.Get(field)
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(field, "Please enter a number > 0!")
		return
	}
	if valueInt <= 0 {
		f.Errors.Add(field, "Please enter a number > 0!")
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
		f.Errors.Add(field, fmt.Sprintf("Too long! Please enter soemthing < %d chars!", max))
	}
}

// MinLength checks that form input is not to short eg. password
func (f *Form) MinLength(field string, min int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < min {
		f.Errors.Add(field, fmt.Sprintf("Too short! Please enter something > %d chars!", min))
	}
}

// ValidMail check if a mail is valid
func (f *Form) ValidMail(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		_, err := mail.ParseAddress(value)
		if err != nil {
			f.Errors.Add(field, "Invalid! Please enter a valid address format!")
		}
	}
}

// Valid returns if posted Form is valid or not
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
