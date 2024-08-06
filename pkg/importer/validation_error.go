package importer

import (
	"errors"
	"strings"
)

var (
	ErrorSameMatterUnderDifferentFolders               = errors.New("matter is under the different folders")
	ErrorSameCustodianEmailUnderDifferentCustodianName = errors.New("custodian email is under different custodian names")
	ErrorRequiredLastIssued                            = errors.New("last issued field is required")
	ErrorAttachmentFileNotFound                        = errors.New("attachment file not found")
	ErrorInvalidEmailAddress                           = errors.New("invalid email address")
	ErrorHoldNameTooLong                               = errors.New("hold name too long")
	ErrorCustodianNotFound                             = errors.New("custodian not found")
)

type ValidationError struct {
	ErrType error
	Errors  []error
}

func newValidationError(err error) *ValidationError {
	return &ValidationError{
		ErrType: err,
		Errors:  []error{},
	}
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(e.ErrType.Error())
	sb.WriteString("\n")

	for _, err := range e.Errors {
		sb.WriteString(" - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *ValidationError) hasErrors() bool {
	return len(e.Errors) > 0
}

func (e *ValidationError) add(err error) {
	if err == nil {
		return
	}
	e.Errors = append(e.Errors, err)
}
