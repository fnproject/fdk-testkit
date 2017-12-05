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

// CallWrapper call wrapper
// swagger:model CallWrapper

type CallWrapper struct {

	// Call object.
	// Required: true
	Call *Call `json:"call"`
}

/* polymorph CallWrapper call false */

// Validate validates this call wrapper
func (m *CallWrapper) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCall(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CallWrapper) validateCall(formats strfmt.Registry) error {

	if err := validate.Required("call", "body", m.Call); err != nil {
		return err
	}

	if m.Call != nil {

		if err := m.Call.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("call")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CallWrapper) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CallWrapper) UnmarshalBinary(b []byte) error {
	var res CallWrapper
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
