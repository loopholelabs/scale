// Code generated by go-swagger; DO NOT EDIT.

package registry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new registry API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for registry API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DeleteRegistryFunctionOrganizationNameTag(params *DeleteRegistryFunctionOrganizationNameTagParams, opts ...ClientOption) (*DeleteRegistryFunctionOrganizationNameTagOK, error)

	GetRegistryFunction(params *GetRegistryFunctionParams, opts ...ClientOption) (*GetRegistryFunctionOK, error)

	GetRegistryFunctionNameTag(params *GetRegistryFunctionNameTagParams, opts ...ClientOption) (*GetRegistryFunctionNameTagOK, error)

	GetRegistryFunctionOrganization(params *GetRegistryFunctionOrganizationParams, opts ...ClientOption) (*GetRegistryFunctionOrganizationOK, error)

	GetRegistryFunctionOrganizationNameTag(params *GetRegistryFunctionOrganizationNameTagParams, opts ...ClientOption) (*GetRegistryFunctionOrganizationNameTagOK, error)

	PostDomain(params *PostDomainParams, opts ...ClientOption) (*PostDomainOK, error)

	PostRegistryFunction(params *PostRegistryFunctionParams, opts ...ClientOption) (*PostRegistryFunctionOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
DeleteRegistryFunctionOrganizationNameTag Deletes a function from the given `organization` given its `name` and `tag`. If the session is scoped to an organization it must be the same as the organization of the function.
*/
func (a *Client) DeleteRegistryFunctionOrganizationNameTag(params *DeleteRegistryFunctionOrganizationNameTagParams, opts ...ClientOption) (*DeleteRegistryFunctionOrganizationNameTagOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteRegistryFunctionOrganizationNameTagParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteRegistryFunctionOrganizationNameTag",
		Method:             "DELETE",
		PathPattern:        "/registry/function/{organization}/{name}/{tag}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteRegistryFunctionOrganizationNameTagReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteRegistryFunctionOrganizationNameTagOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteRegistryFunctionOrganizationNameTag: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetRegistryFunction Lists all the functions in the default organization.
*/
func (a *Client) GetRegistryFunction(params *GetRegistryFunctionParams, opts ...ClientOption) (*GetRegistryFunctionOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRegistryFunctionParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetRegistryFunction",
		Method:             "GET",
		PathPattern:        "/registry/function",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetRegistryFunctionReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetRegistryFunctionOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetRegistryFunction: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetRegistryFunctionNameTag Retrieves a function from the default organization given its `name` and `tag`. If the session is scoped to the same `organization`, functions that are not public will be returned, otherwise only public functions will be returned.
*/
func (a *Client) GetRegistryFunctionNameTag(params *GetRegistryFunctionNameTagParams, opts ...ClientOption) (*GetRegistryFunctionNameTagOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRegistryFunctionNameTagParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetRegistryFunctionNameTag",
		Method:             "GET",
		PathPattern:        "/registry/function/{name}/{tag}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetRegistryFunctionNameTagReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetRegistryFunctionNameTagOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetRegistryFunctionNameTag: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetRegistryFunctionOrganization Lists all the functions in the given `organization`. If the session is scoped to the same `organization`, functions that are not public will be returned, otherwise only public functions from the `organization` will be returned.
*/
func (a *Client) GetRegistryFunctionOrganization(params *GetRegistryFunctionOrganizationParams, opts ...ClientOption) (*GetRegistryFunctionOrganizationOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRegistryFunctionOrganizationParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetRegistryFunctionOrganization",
		Method:             "GET",
		PathPattern:        "/registry/function/{organization}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetRegistryFunctionOrganizationReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetRegistryFunctionOrganizationOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetRegistryFunctionOrganization: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetRegistryFunctionOrganizationNameTag Retrieves a function from the given `organization` given its `name` and `tag`. If the session is scoped to the same `organization`, functions that are not public will be returned, otherwise only public functions will be returned.
*/
func (a *Client) GetRegistryFunctionOrganizationNameTag(params *GetRegistryFunctionOrganizationNameTagParams, opts ...ClientOption) (*GetRegistryFunctionOrganizationNameTagOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetRegistryFunctionOrganizationNameTagParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetRegistryFunctionOrganizationNameTag",
		Method:             "GET",
		PathPattern:        "/registry/function/{organization}/{name}/{tag}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetRegistryFunctionOrganizationNameTagReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetRegistryFunctionOrganizationNameTagOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetRegistryFunctionOrganizationNameTag: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PostDomain Creates a new domain. If the session is scoped to an organization, the domain will be created in that `organization`, otherwise the domain will be created to the user's default `organization`.
*/
func (a *Client) PostDomain(params *PostDomainParams, opts ...ClientOption) (*PostDomainOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostDomainParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostDomain",
		Method:             "POST",
		PathPattern:        "/domain",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostDomainReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostDomainOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostDomain: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PostRegistryFunction Uploads a function to the Scale Registry. If the session is scoped to an organization, the function will be uploaded to that `organization`, otherwise the function will be uploaded to the user's default `organization`.
*/
func (a *Client) PostRegistryFunction(params *PostRegistryFunctionParams, opts ...ClientOption) (*PostRegistryFunctionOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostRegistryFunctionParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostRegistryFunction",
		Method:             "POST",
		PathPattern:        "/registry/function",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"multipart/form-data"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostRegistryFunctionReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostRegistryFunctionOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostRegistryFunction: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}