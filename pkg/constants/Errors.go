package constants

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorUserNotAuthorized = errors.New("user not authorized")

	ErrorVariantAlreadyExists = errors.New("variant already exists")
	ErrorVariantTooLong       = errors.New("variant too long: >16")
	ErrorVariantNotFound      = errors.New("variant not found")
	ErrorNoVariantsYet        = errors.New("no variants yet")
	ErrorVariantCompleted     = errors.New("variant completed")

	ErrorQuestionAlreadyExists = errors.New("question already exists")
	ErrorQuestionNotFound      = errors.New("question not found")
	ErrorQuestionLimitExceeded = errors.New("question limit exceeded")

	ErrorTestNotFound = errors.New("testing not found")
)
