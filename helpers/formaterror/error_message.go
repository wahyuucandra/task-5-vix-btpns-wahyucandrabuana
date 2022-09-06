package formaterror

import (
	"errors"
	"strings"
)

func ErrorMessage (err string) error {
	
	if strings.Contains(err, "pkey") {
		return errors.New("user id already taken")
	}else if strings.Contains(err, "email_key") {
		return errors.New("email already taken")
	}else if strings.Contains(err, "title") {
		return errors.New("title already taken")
	}else if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect password")
	}
	return errors.New("incorrect details")
}