package iam

import "errors"

func errInvalid(msg string) error {
	return errors.New(msg)
}
