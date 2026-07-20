package validator

import (
	"regexp"
	"unicode"
)

var (
	EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z]{2,}$")
)

type Validator struct {
	errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Errors() map[string]string {
	return v.errors
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func Matches(rx *regexp.Regexp, values string) bool {
	if rx == nil {
		return false
	}

	return rx.MatchString(values)
}

func (v *Validator) ValidateEmail(email string) {
	v.Check(email != "", "email", "must not be empty")
	v.Check(Matches(EmailRegex, email), "email", "email not valid")
}

func (v *Validator) ValidatePassword(password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be atleast 8 characters")
	v.Check(len(password) <= 72, "password", "must be atleast 72 characters long")
	v.Check(StrongPassword(password), "password", "must contain at least one uppercase letter, lowercase letter, number, and special character")
}

func (v *Validator) ValidateUser(email, password string) {
	v.ValidateEmail(email)
	v.ValidatePassword(password)
}

func StrongPassword(password string) bool {
	var hasLower, hasUpper, hasDigit, hasSymbol bool

	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsSymbol(r) || unicode.IsPunct(r):
			hasSymbol = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSymbol
}
