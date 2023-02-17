// Code generated by go-swagger; DO NOT EDIT.

package registry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/loopholelabs/scale/go/client/models"
)

// GetRegistryFunctionOrganizationNameTagReader is a Reader for the GetRegistryFunctionOrganizationNameTag structure.
type GetRegistryFunctionOrganizationNameTagReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetRegistryFunctionOrganizationNameTagReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetRegistryFunctionOrganizationNameTagOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetRegistryFunctionOrganizationNameTagBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewGetRegistryFunctionOrganizationNameTagUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetRegistryFunctionOrganizationNameTagNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetRegistryFunctionOrganizationNameTagInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetRegistryFunctionOrganizationNameTagOK creates a GetRegistryFunctionOrganizationNameTagOK with default headers values
func NewGetRegistryFunctionOrganizationNameTagOK() *GetRegistryFunctionOrganizationNameTagOK {
	return &GetRegistryFunctionOrganizationNameTagOK{}
}

/*
GetRegistryFunctionOrganizationNameTagOK describes a response with status code 200, with default header values.

OK
*/
type GetRegistryFunctionOrganizationNameTagOK struct {
	Payload *models.ModelsGetFunctionResponse
}

// IsSuccess returns true when this get registry function organization name tag o k response has a 2xx status code
func (o *GetRegistryFunctionOrganizationNameTagOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get registry function organization name tag o k response has a 3xx status code
func (o *GetRegistryFunctionOrganizationNameTagOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get registry function organization name tag o k response has a 4xx status code
func (o *GetRegistryFunctionOrganizationNameTagOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get registry function organization name tag o k response has a 5xx status code
func (o *GetRegistryFunctionOrganizationNameTagOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get registry function organization name tag o k response a status code equal to that given
func (o *GetRegistryFunctionOrganizationNameTagOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get registry function organization name tag o k response
func (o *GetRegistryFunctionOrganizationNameTagOK) Code() int {
	return 200
}

func (o *GetRegistryFunctionOrganizationNameTagOK) Error() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagOK  %+v", 200, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagOK) String() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagOK  %+v", 200, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagOK) GetPayload() *models.ModelsGetFunctionResponse {
	return o.Payload
}

func (o *GetRegistryFunctionOrganizationNameTagOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ModelsGetFunctionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetRegistryFunctionOrganizationNameTagBadRequest creates a GetRegistryFunctionOrganizationNameTagBadRequest with default headers values
func NewGetRegistryFunctionOrganizationNameTagBadRequest() *GetRegistryFunctionOrganizationNameTagBadRequest {
	return &GetRegistryFunctionOrganizationNameTagBadRequest{}
}

/*
GetRegistryFunctionOrganizationNameTagBadRequest describes a response with status code 400, with default header values.

Bad Request
*/
type GetRegistryFunctionOrganizationNameTagBadRequest struct {
	Payload string
}

// IsSuccess returns true when this get registry function organization name tag bad request response has a 2xx status code
func (o *GetRegistryFunctionOrganizationNameTagBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get registry function organization name tag bad request response has a 3xx status code
func (o *GetRegistryFunctionOrganizationNameTagBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get registry function organization name tag bad request response has a 4xx status code
func (o *GetRegistryFunctionOrganizationNameTagBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get registry function organization name tag bad request response has a 5xx status code
func (o *GetRegistryFunctionOrganizationNameTagBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get registry function organization name tag bad request response a status code equal to that given
func (o *GetRegistryFunctionOrganizationNameTagBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get registry function organization name tag bad request response
func (o *GetRegistryFunctionOrganizationNameTagBadRequest) Code() int {
	return 400
}

func (o *GetRegistryFunctionOrganizationNameTagBadRequest) Error() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagBadRequest  %+v", 400, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagBadRequest) String() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagBadRequest  %+v", 400, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagBadRequest) GetPayload() string {
	return o.Payload
}

func (o *GetRegistryFunctionOrganizationNameTagBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetRegistryFunctionOrganizationNameTagUnauthorized creates a GetRegistryFunctionOrganizationNameTagUnauthorized with default headers values
func NewGetRegistryFunctionOrganizationNameTagUnauthorized() *GetRegistryFunctionOrganizationNameTagUnauthorized {
	return &GetRegistryFunctionOrganizationNameTagUnauthorized{}
}

/*
GetRegistryFunctionOrganizationNameTagUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type GetRegistryFunctionOrganizationNameTagUnauthorized struct {
	Payload string
}

// IsSuccess returns true when this get registry function organization name tag unauthorized response has a 2xx status code
func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get registry function organization name tag unauthorized response has a 3xx status code
func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get registry function organization name tag unauthorized response has a 4xx status code
func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get registry function organization name tag unauthorized response has a 5xx status code
func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get registry function organization name tag unauthorized response a status code equal to that given
func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get registry function organization name tag unauthorized response
func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) Code() int {
	return 401
}

func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) Error() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagUnauthorized  %+v", 401, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) String() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagUnauthorized  %+v", 401, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) GetPayload() string {
	return o.Payload
}

func (o *GetRegistryFunctionOrganizationNameTagUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetRegistryFunctionOrganizationNameTagNotFound creates a GetRegistryFunctionOrganizationNameTagNotFound with default headers values
func NewGetRegistryFunctionOrganizationNameTagNotFound() *GetRegistryFunctionOrganizationNameTagNotFound {
	return &GetRegistryFunctionOrganizationNameTagNotFound{}
}

/*
GetRegistryFunctionOrganizationNameTagNotFound describes a response with status code 404, with default header values.

Not Found
*/
type GetRegistryFunctionOrganizationNameTagNotFound struct {
	Payload string
}

// IsSuccess returns true when this get registry function organization name tag not found response has a 2xx status code
func (o *GetRegistryFunctionOrganizationNameTagNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get registry function organization name tag not found response has a 3xx status code
func (o *GetRegistryFunctionOrganizationNameTagNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get registry function organization name tag not found response has a 4xx status code
func (o *GetRegistryFunctionOrganizationNameTagNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get registry function organization name tag not found response has a 5xx status code
func (o *GetRegistryFunctionOrganizationNameTagNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get registry function organization name tag not found response a status code equal to that given
func (o *GetRegistryFunctionOrganizationNameTagNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get registry function organization name tag not found response
func (o *GetRegistryFunctionOrganizationNameTagNotFound) Code() int {
	return 404
}

func (o *GetRegistryFunctionOrganizationNameTagNotFound) Error() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagNotFound  %+v", 404, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagNotFound) String() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagNotFound  %+v", 404, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagNotFound) GetPayload() string {
	return o.Payload
}

func (o *GetRegistryFunctionOrganizationNameTagNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetRegistryFunctionOrganizationNameTagInternalServerError creates a GetRegistryFunctionOrganizationNameTagInternalServerError with default headers values
func NewGetRegistryFunctionOrganizationNameTagInternalServerError() *GetRegistryFunctionOrganizationNameTagInternalServerError {
	return &GetRegistryFunctionOrganizationNameTagInternalServerError{}
}

/*
GetRegistryFunctionOrganizationNameTagInternalServerError describes a response with status code 500, with default header values.

Internal Server Error
*/
type GetRegistryFunctionOrganizationNameTagInternalServerError struct {
	Payload string
}

// IsSuccess returns true when this get registry function organization name tag internal server error response has a 2xx status code
func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get registry function organization name tag internal server error response has a 3xx status code
func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get registry function organization name tag internal server error response has a 4xx status code
func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get registry function organization name tag internal server error response has a 5xx status code
func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get registry function organization name tag internal server error response a status code equal to that given
func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get registry function organization name tag internal server error response
func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) Code() int {
	return 500
}

func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) Error() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagInternalServerError  %+v", 500, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) String() string {
	return fmt.Sprintf("[GET /registry/function/{organization}/{name}/{tag}][%d] getRegistryFunctionOrganizationNameTagInternalServerError  %+v", 500, o.Payload)
}

func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) GetPayload() string {
	return o.Payload
}

func (o *GetRegistryFunctionOrganizationNameTagInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
