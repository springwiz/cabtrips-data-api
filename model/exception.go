// Package model contains all the entities
package model

// Exception entity
type Exception struct {
	CodeInt uint64

	CodeString string

	CodeDescription string
}

func (e *Exception) Error() string {
	return e.CodeString
}

// NewException creates a new exception
func NewException(codeString string, codeDescription string) *Exception {
	NewException := &Exception{
		CodeInt:         1,
		CodeString:      codeString,
		CodeDescription: codeDescription,
	}
	return NewException
}
