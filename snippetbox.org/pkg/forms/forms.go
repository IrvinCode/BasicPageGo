package forms

import (
	"strings"
	"unicode/utf8"
	"regexp"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9])")

type NewSnippet struct {
	Title string
	Content string
	Expires string
	Failures map[string]string
}

func (f *NewSnippet) Valid() bool {
	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
	} else if utf8.RuneCountInString(f.Title) > 100 {
		f.Failures["Title"] = "Title cannot be longer than 100 characters"
	}

	if strings.TrimSpace(f.Content) == "" {
		f.Failures["Content"] = "Content is required"
	}

	permitted := map[string]bool{"3600": true, "86400": true, "31536000": true}
	if strings.TrimSpace(f.Expires) == "" {
		f.Failures["Expires"] = "Expiry time is required"
	} else if !permitted[f.Expires] {
		f.Failures["Expires"] = "Expiry time must be 3600, 86400 or 31536000 seconds"
	}

	return len(f.Failures) == 0
}

type SignupUser struct {
	Name string
	Email string
	Password string
	Failures map[string]string
}

func (f *SignupUser) Valid() bool {
	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		f.Failures["Name"] = "Name is required"
	}

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	} else if len(f.Email) > 254 || !rxEmail.MatchString(f.Email) {
		f.Failures["Email"] = "Email is not a valid address"
	}

	if utf8.RuneCountInString(f.Password) < 8 {
		f.Failures["Password"] = "Password cannot be shorter than 8 characters"
	}

	return len(f.Failures) == 0
}

type LoginUser struct {
	Email string
	Password string
	Failures map[string]string
}

func (f *LoginUser) Valid() bool {
	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	}

	if strings.TrimSpace(f.Password) == "" {
		f.Failures["Password"] = "Password is required"
	}

	return len(f.Failures) == 0
}

type SignupAdmin struct {
	Name string
	Email string
	Password string
	Failures map[string]string
}


func (f *SignupAdmin) Valid() bool {
	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		f.Failures["Name"] = "Name is required"
	}

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	} else if len(f.Email) > 254 || !rxEmail.MatchString(f.Email) {
		f.Failures["Email"] = "Email is not a valid address"
	}

	if utf8.RuneCountInString(f.Password) < 8 {
		f.Failures["Password"] = "Password cannot be shorter than 8 characters"
	}

	return len(f.Failures) == 0
}

type DeleteSnippet struct {
	Id string
	Failures map[string]string
}

func (f *DeleteSnippet) Valid() bool {
	return len(f.Failures) == 0
}