// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ModelsDeployFunctionResponse models deploy function response
//
// swagger:model models.DeployFunctionResponse
type ModelsDeployFunctionResponse struct {

	// identifier
	Identifier string `json:"identifier,omitempty"`

	// root domain
	RootDomain string `json:"root_domain,omitempty"`

	// subdomain
	Subdomain string `json:"subdomain,omitempty"`
}

// Validate validates this models deploy function response
func (m *ModelsDeployFunctionResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this models deploy function response based on context it is used
func (m *ModelsDeployFunctionResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelsDeployFunctionResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelsDeployFunctionResponse) UnmarshalBinary(b []byte) error {
	var res ModelsDeployFunctionResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}