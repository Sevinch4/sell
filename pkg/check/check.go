package check

import "errors"


func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password length should be more than 6")
	}

	return nil
}