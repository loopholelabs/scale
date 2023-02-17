// Code generated by go-swagger; DO NOT EDIT.

package access

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/loopholelabs/scale/go/client/models"
)

// GetAccessApikeyNameReader is a Reader for the GetAccessApikeyName structure.
type GetAccessApikeyNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAccessApikeyNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAccessApikeyNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetAccessApikeyNameUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetAccessApikeyNameNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAccessApikeyNameInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetAccessApikeyNameOK creates a GetAccessApikeyNameOK with default headers values
func NewGetAccessApikeyNameOK() *GetAccessApikeyNameOK {
	return &GetAccessApikeyNameOK{}
}

/*
GetAccessApikeyNameOK describes a response with status code 200, with default header values.

OK
*/
type GetAccessApikeyNameOK struct {
	Payload *models.ModelsGetAPIKeyResponse
}

// IsSuccess returns true when this get access apikey name o k response has a 2xx status code
func (o *GetAccessApikeyNameOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get access apikey name o k response has a 3xx status code
func (o *GetAccessApikeyNameOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get access apikey name o k response has a 4xx status code
func (o *GetAccessApikeyNameOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get access apikey name o k response has a 5xx status code
func (o *GetAccessApikeyNameOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get access apikey name o k response a status code equal to that given
func (o *GetAccessApikeyNameOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get access apikey name o k response
func (o *GetAccessApikeyNameOK) Code() int {
	return 200
}

func (o *GetAccessApikeyNameOK) Error() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameOK  %+v", 200, o.Payload)
}

func (o *GetAccessApikeyNameOK) String() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameOK  %+v", 200, o.Payload)
}

func (o *GetAccessApikeyNameOK) GetPayload() *models.ModelsGetAPIKeyResponse {
	return o.Payload
}

func (o *GetAccessApikeyNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ModelsGetAPIKeyResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAccessApikeyNameUnauthorized creates a GetAccessApikeyNameUnauthorized with default headers values
func NewGetAccessApikeyNameUnauthorized() *GetAccessApikeyNameUnauthorized {
	return &GetAccessApikeyNameUnauthorized{}
}

/*
GetAccessApikeyNameUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type GetAccessApikeyNameUnauthorized struct {
	Payload string
}

// IsSuccess returns true when this get access apikey name unauthorized response has a 2xx status code
func (o *GetAccessApikeyNameUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get access apikey name unauthorized response has a 3xx status code
func (o *GetAccessApikeyNameUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get access apikey name unauthorized response has a 4xx status code
func (o *GetAccessApikeyNameUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get access apikey name unauthorized response has a 5xx status code
func (o *GetAccessApikeyNameUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get access apikey name unauthorized response a status code equal to that given
func (o *GetAccessApikeyNameUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get access apikey name unauthorized response
func (o *GetAccessApikeyNameUnauthorized) Code() int {
	return 401
}

func (o *GetAccessApikeyNameUnauthorized) Error() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameUnauthorized  %+v", 401, o.Payload)
}

func (o *GetAccessApikeyNameUnauthorized) String() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameUnauthorized  %+v", 401, o.Payload)
}

func (o *GetAccessApikeyNameUnauthorized) GetPayload() string {
	return o.Payload
}

func (o *GetAccessApikeyNameUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAccessApikeyNameNotFound creates a GetAccessApikeyNameNotFound with default headers values
func NewGetAccessApikeyNameNotFound() *GetAccessApikeyNameNotFound {
	return &GetAccessApikeyNameNotFound{}
}

/*
GetAccessApikeyNameNotFound describes a response with status code 404, with default header values.

Not Found
*/
type GetAccessApikeyNameNotFound struct {
	Payload string
}

// IsSuccess returns true when this get access apikey name not found response has a 2xx status code
func (o *GetAccessApikeyNameNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get access apikey name not found response has a 3xx status code
func (o *GetAccessApikeyNameNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get access apikey name not found response has a 4xx status code
func (o *GetAccessApikeyNameNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get access apikey name not found response has a 5xx status code
func (o *GetAccessApikeyNameNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get access apikey name not found response a status code equal to that given
func (o *GetAccessApikeyNameNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get access apikey name not found response
func (o *GetAccessApikeyNameNotFound) Code() int {
	return 404
}

func (o *GetAccessApikeyNameNotFound) Error() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameNotFound  %+v", 404, o.Payload)
}

func (o *GetAccessApikeyNameNotFound) String() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameNotFound  %+v", 404, o.Payload)
}

func (o *GetAccessApikeyNameNotFound) GetPayload() string {
	return o.Payload
}

func (o *GetAccessApikeyNameNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAccessApikeyNameInternalServerError creates a GetAccessApikeyNameInternalServerError with default headers values
func NewGetAccessApikeyNameInternalServerError() *GetAccessApikeyNameInternalServerError {
	return &GetAccessApikeyNameInternalServerError{}
}

/*
GetAccessApikeyNameInternalServerError describes a response with status code 500, with default header values.

Internal Server Error
*/
type GetAccessApikeyNameInternalServerError struct {
	Payload string
}

// IsSuccess returns true when this get access apikey name internal server error response has a 2xx status code
func (o *GetAccessApikeyNameInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get access apikey name internal server error response has a 3xx status code
func (o *GetAccessApikeyNameInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get access apikey name internal server error response has a 4xx status code
func (o *GetAccessApikeyNameInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get access apikey name internal server error response has a 5xx status code
func (o *GetAccessApikeyNameInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get access apikey name internal server error response a status code equal to that given
func (o *GetAccessApikeyNameInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get access apikey name internal server error response
func (o *GetAccessApikeyNameInternalServerError) Code() int {
	return 500
}

func (o *GetAccessApikeyNameInternalServerError) Error() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAccessApikeyNameInternalServerError) String() string {
	return fmt.Sprintf("[GET /access/apikey/{name}][%d] getAccessApikeyNameInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAccessApikeyNameInternalServerError) GetPayload() string {
	return o.Payload
}

func (o *GetAccessApikeyNameInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
