// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ModelsGetAPIKeyResponse models get API key response
//
// swagger:model models.GetAPIKeyResponse
type ModelsGetAPIKeyResponse struct {

	// created at
	CreatedAt string `json:"created_at,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// organization
	Organization string `json:"organization,omitempty"`
}

// Validate validates this models get API key response
func (m *ModelsGetAPIKeyResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this models get API key response based on context it is used
func (m *ModelsGetAPIKeyResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelsGetAPIKeyResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelsGetAPIKeyResponse) UnmarshalBinary(b []byte) error {
	var res ModelsGetAPIKeyResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}