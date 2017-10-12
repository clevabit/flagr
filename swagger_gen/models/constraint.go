// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Constraint constraint
// swagger:model constraint

type Constraint struct {

	// id
	// Read Only: true
	// Minimum: 1
	ID int64 `json:"id,omitempty"`

	// operator
	// Required: true
	// Min Length: 1
	Operator *string `json:"operator"`

	// property
	// Required: true
	// Min Length: 1
	Property *string `json:"property"`

	// value
	// Required: true
	// Min Length: 1
	Value *string `json:"value"`
}

/* polymorph constraint id false */

/* polymorph constraint operator false */

/* polymorph constraint property false */

/* polymorph constraint value false */

// Validate validates this constraint
func (m *Constraint) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateOperator(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateProperty(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Constraint) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.MinimumInt("id", "body", int64(m.ID), 1, false); err != nil {
		return err
	}

	return nil
}

var constraintTypeOperatorPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["EQ","NEQ","LT","LTE","GT","GTE","EREG","NEREG","IN","NOTIN","CONTAINS","NOTCONTAINS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		constraintTypeOperatorPropEnum = append(constraintTypeOperatorPropEnum, v)
	}
}

const (
	// ConstraintOperatorEQ captures enum value "EQ"
	ConstraintOperatorEQ string = "EQ"
	// ConstraintOperatorNEQ captures enum value "NEQ"
	ConstraintOperatorNEQ string = "NEQ"
	// ConstraintOperatorLT captures enum value "LT"
	ConstraintOperatorLT string = "LT"
	// ConstraintOperatorLTE captures enum value "LTE"
	ConstraintOperatorLTE string = "LTE"
	// ConstraintOperatorGT captures enum value "GT"
	ConstraintOperatorGT string = "GT"
	// ConstraintOperatorGTE captures enum value "GTE"
	ConstraintOperatorGTE string = "GTE"
	// ConstraintOperatorEREG captures enum value "EREG"
	ConstraintOperatorEREG string = "EREG"
	// ConstraintOperatorNEREG captures enum value "NEREG"
	ConstraintOperatorNEREG string = "NEREG"
	// ConstraintOperatorIN captures enum value "IN"
	ConstraintOperatorIN string = "IN"
	// ConstraintOperatorNOTIN captures enum value "NOTIN"
	ConstraintOperatorNOTIN string = "NOTIN"
	// ConstraintOperatorCONTAINS captures enum value "CONTAINS"
	ConstraintOperatorCONTAINS string = "CONTAINS"
	// ConstraintOperatorNOTCONTAINS captures enum value "NOTCONTAINS"
	ConstraintOperatorNOTCONTAINS string = "NOTCONTAINS"
)

// prop value enum
func (m *Constraint) validateOperatorEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, constraintTypeOperatorPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Constraint) validateOperator(formats strfmt.Registry) error {

	if err := validate.Required("operator", "body", m.Operator); err != nil {
		return err
	}

	if err := validate.MinLength("operator", "body", string(*m.Operator), 1); err != nil {
		return err
	}

	// value enum
	if err := m.validateOperatorEnum("operator", "body", *m.Operator); err != nil {
		return err
	}

	return nil
}

func (m *Constraint) validateProperty(formats strfmt.Registry) error {

	if err := validate.Required("property", "body", m.Property); err != nil {
		return err
	}

	if err := validate.MinLength("property", "body", string(*m.Property), 1); err != nil {
		return err
	}

	return nil
}

func (m *Constraint) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("value", "body", m.Value); err != nil {
		return err
	}

	if err := validate.MinLength("value", "body", string(*m.Value), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Constraint) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Constraint) UnmarshalBinary(b []byte) error {
	var res Constraint
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
