// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: xds/annotations/v3/security.proto

package v3

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on FieldSecurityAnnotation with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *FieldSecurityAnnotation) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FieldSecurityAnnotation with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FieldSecurityAnnotationMultiError, or nil if none found.
func (m *FieldSecurityAnnotation) ValidateAll() error {
	return m.validate(true)
}

func (m *FieldSecurityAnnotation) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ConfigureForUntrustedDownstream

	// no validation rules for ConfigureForUntrustedUpstream

	if len(errors) > 0 {
		return FieldSecurityAnnotationMultiError(errors)
	}

	return nil
}

// FieldSecurityAnnotationMultiError is an error wrapping multiple validation
// errors returned by FieldSecurityAnnotation.ValidateAll() if the designated
// constraints aren't met.
type FieldSecurityAnnotationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FieldSecurityAnnotationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FieldSecurityAnnotationMultiError) AllErrors() []error { return m }

// FieldSecurityAnnotationValidationError is the validation error returned by
// FieldSecurityAnnotation.Validate if the designated constraints aren't met.
type FieldSecurityAnnotationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FieldSecurityAnnotationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FieldSecurityAnnotationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FieldSecurityAnnotationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FieldSecurityAnnotationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FieldSecurityAnnotationValidationError) ErrorName() string {
	return "FieldSecurityAnnotationValidationError"
}

// Error satisfies the builtin error interface
func (e FieldSecurityAnnotationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFieldSecurityAnnotation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FieldSecurityAnnotationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FieldSecurityAnnotationValidationError{}
