// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CallsWrapper calls wrapper
// swagger:model CallsWrapper

type CallsWrapper struct {

	// calls
	// Required: true
	Calls CallsWrapperCalls `json:"calls"`

	// error
	Error *ErrorBody `json:"error,omitempty"`

	// cursor to send with subsequent request to receive the next page, if non-empty
	// Read Only: true
	NextCursor string `json:"next_cursor,omitempty"`
}

/* polymorph CallsWrapper calls false */

/* polymorph CallsWrapper error false */

/* polymorph CallsWrapper next_cursor false */

// Validate validates this calls wrapper
func (m *CallsWrapper) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCalls(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateError(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CallsWrapper) validateCalls(formats strfmt.Registry) error {

	if err := validate.Required("calls", "body", m.Calls); err != nil {
		return err
	}

	return nil
}

func (m *CallsWrapper) validateError(formats strfmt.Registry) error {

	if swag.IsZero(m.Error) { // not required
		return nil
	}

	if m.Error != nil {

		if err := m.Error.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("error")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CallsWrapper) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CallsWrapper) UnmarshalBinary(b []byte) error {
	var res CallsWrapper
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
