// Code generated by go-swagger; DO NOT EDIT.

package access

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDeleteAccessApikeyIDParams creates a new DeleteAccessApikeyIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteAccessApikeyIDParams() *DeleteAccessApikeyIDParams {
	return &DeleteAccessApikeyIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteAccessApikeyIDParamsWithTimeout creates a new DeleteAccessApikeyIDParams object
// with the ability to set a timeout on a request.
func NewDeleteAccessApikeyIDParamsWithTimeout(timeout time.Duration) *DeleteAccessApikeyIDParams {
	return &DeleteAccessApikeyIDParams{
		timeout: timeout,
	}
}

// NewDeleteAccessApikeyIDParamsWithContext creates a new DeleteAccessApikeyIDParams object
// with the ability to set a context for a request.
func NewDeleteAccessApikeyIDParamsWithContext(ctx context.Context) *DeleteAccessApikeyIDParams {
	return &DeleteAccessApikeyIDParams{
		Context: ctx,
	}
}

// NewDeleteAccessApikeyIDParamsWithHTTPClient creates a new DeleteAccessApikeyIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteAccessApikeyIDParamsWithHTTPClient(client *http.Client) *DeleteAccessApikeyIDParams {
	return &DeleteAccessApikeyIDParams{
		HTTPClient: client,
	}
}

/*
DeleteAccessApikeyIDParams contains all the parameters to send to the API endpoint

	for the delete access apikey ID operation.

	Typically these are written to a http.Request.
*/
type DeleteAccessApikeyIDParams struct {

	/* ID.

	   API Key ID
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete access apikey ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteAccessApikeyIDParams) WithDefaults() *DeleteAccessApikeyIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete access apikey ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteAccessApikeyIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) WithTimeout(timeout time.Duration) *DeleteAccessApikeyIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) WithContext(ctx context.Context) *DeleteAccessApikeyIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) WithHTTPClient(client *http.Client) *DeleteAccessApikeyIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) WithID(id string) *DeleteAccessApikeyIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete access apikey ID params
func (o *DeleteAccessApikeyIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteAccessApikeyIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
