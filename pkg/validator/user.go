package validator

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type UserValidator struct{}

func (v UserValidator) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your name"
		}
		return errMsg[0]
	case "min, max":
		return "Your name should be between 3 to 20 characters"
	case "fullName":
		return "Name should not include any special characters or numbers"
	default:
		return "Something went wrong, please try again later"
	}
}

func (v UserValidator) Email(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your email"
		}
		return errMsg[0]
	case "min", "max", "email":
		return "Please enter a valid email"
	default:
		return "Something went wrong, please try again later"
	}
}

func (v UserValidator) Password(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your password"
	case "min", "max":
		return "Your password should be between 3 and 50 characters"
	case "eqfield":
		return "Your password does not match"
	default:
		return "Something went wrong, please try again later"
	}
}

func (v UserValidator) Login(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Wrong JSON format"
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				return v.Email(err.Tag())
			case "Password":
				return v.Password(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

func (v UserValidator) CreateUser(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				return v.Name(err.Tag())
			case "Email":
				return v.Email(err.Tag())
			case "Password":
				return v.Password(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}
