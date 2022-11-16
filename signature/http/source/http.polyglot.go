// Code generated by polyglot-go 0.5.0, DO NOT EDIT.
// source: signature/http/source/http.proto

package http

import (
	"errors"
	"github.com/loopholelabs/polyglot-go"
)

var (
	NilDecode = errors.New("cannot decode into a nil root struct")
)

type HttpContext struct {
	Request  *HttpRequest
	Response *HttpResponse
}

func NewHttpContext() *HttpContext {
	return &HttpContext{
		Request:  NewHttpRequest(),
		Response: NewHttpResponse(),
	}
}

func (x *HttpContext) error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *HttpContext) internalEncode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		x.Request.internalEncode(b)
		x.Response.internalEncode(b)
	}
}

func (x *HttpContext) internalDecode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	d := polyglot.GetDecoder(b)
	defer d.Return()
	return x.decode(d)
}

func (x *HttpContext) decode(d *polyglot.Decoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	if x.Request == nil {
		x.Request = NewHttpRequest()
	}
	err = x.Request.decode(d)
	if err != nil {
		return err
	}
	if x.Response == nil {
		x.Response = NewHttpResponse()
	}
	err = x.Response.decode(d)
	if err != nil {
		return err
	}
	return nil
}

type HttpRequestHeadersMap map[string]*HttpStringList

func NewHttpRequestHeadersMap(size uint32) map[string]*HttpStringList {
	return make(map[string]*HttpStringList, size)
}

func (x HttpRequestHeadersMap) internalEncode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {
		polyglot.Encoder(b).Map(uint32(len(x)), polyglot.StringKind, polyglot.AnyKind)
		for k, v := range x {
			polyglot.Encoder(b).String(k)
			v.internalEncode(b)
		}
	}
}

func (x HttpRequestHeadersMap) decode(d *polyglot.Decoder, size uint32) error {
	if size == 0 {
		return nil
	}
	var k string
	var v *HttpStringList
	var err error
	for i := uint32(0); i < size; i++ {
		k, err = d.String()
		if err != nil {
			return err
		}
		v = NewHttpStringList()
		err = v.decode(d)
		if err != nil {
			return err
		}
		x[k] = v
	}
	return nil
}

type HttpRequest struct {
	Headers       HttpRequestHeadersMap
	URI           string
	Method        string
	ContentLength int64
	Protocol      string
	IP            string
	Body          []byte
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{}
}

func (x *HttpRequest) error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *HttpRequest) internalEncode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).String(x.URI).String(x.Method).Int64(x.ContentLength).String(x.Protocol).String(x.IP).Bytes(x.Body)
		x.Headers.internalEncode(b)
	}
}

func (x *HttpRequest) internalDecode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	d := polyglot.GetDecoder(b)
	defer d.Return()
	return x.decode(d)
}

func (x *HttpRequest) decode(d *polyglot.Decoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.URI, err = d.String()
	if err != nil {
		return err
	}
	x.Method, err = d.String()
	if err != nil {
		return err
	}
	x.ContentLength, err = d.Int64()
	if err != nil {
		return err
	}
	x.Protocol, err = d.String()
	if err != nil {
		return err
	}
	x.IP, err = d.String()
	if err != nil {
		return err
	}
	x.Body, err = d.Bytes(nil)
	if err != nil {
		return err
	}
	if !d.Nil() {
		HeadersSize, err := d.Map(polyglot.StringKind, polyglot.AnyKind)
		if err != nil {
			return err
		}
		x.Headers = NewHttpRequestHeadersMap(HeadersSize)
		err = x.Headers.decode(d, HeadersSize)
		if err != nil {
			return err
		}
	}
	return nil
}

type HttpResponseHeadersMap map[string]*HttpStringList

func NewHttpResponseHeadersMap(size uint32) map[string]*HttpStringList {
	return make(map[string]*HttpStringList, size)
}

func (x HttpResponseHeadersMap) internalEncode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {
		polyglot.Encoder(b).Map(uint32(len(x)), polyglot.StringKind, polyglot.AnyKind)
		for k, v := range x {
			polyglot.Encoder(b).String(k)
			v.internalEncode(b)
		}
	}
}

func (x HttpResponseHeadersMap) decode(d *polyglot.Decoder, size uint32) error {
	if size == 0 {
		return nil
	}
	var k string
	var v *HttpStringList
	var err error
	for i := uint32(0); i < size; i++ {
		k, err = d.String()
		if err != nil {
			return err
		}
		v = NewHttpStringList()
		err = v.decode(d)
		if err != nil {
			return err
		}
		x[k] = v
	}
	return nil
}

type HttpResponse struct {
	Headers    HttpResponseHeadersMap
	StatusCode int32
	Body       []byte
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{}
}

func (x *HttpResponse) error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *HttpResponse) internalEncode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Int32(x.StatusCode).Bytes(x.Body)
		x.Headers.internalEncode(b)
	}
}

func (x *HttpResponse) internalDecode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	d := polyglot.GetDecoder(b)
	defer d.Return()
	return x.decode(d)
}

func (x *HttpResponse) decode(d *polyglot.Decoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	x.StatusCode, err = d.Int32()
	if err != nil {
		return err
	}
	x.Body, err = d.Bytes(nil)
	if err != nil {
		return err
	}
	if !d.Nil() {
		HeadersSize, err := d.Map(polyglot.StringKind, polyglot.AnyKind)
		if err != nil {
			return err
		}
		x.Headers = NewHttpResponseHeadersMap(HeadersSize)
		err = x.Headers.decode(d, HeadersSize)
		if err != nil {
			return err
		}
	}
	return nil
}

type HttpStringList struct {
	Value []string
}

func NewHttpStringList() *HttpStringList {
	return &HttpStringList{}
}

func (x *HttpStringList) error(b *polyglot.Buffer, err error) {
	polyglot.Encoder(b).Error(err)
}

func (x *HttpStringList) internalEncode(b *polyglot.Buffer) {
	if x == nil {
		polyglot.Encoder(b).Nil()
	} else {

		polyglot.Encoder(b).Slice(uint32(len(x.Value)), polyglot.StringKind)
		for _, v := range x.Value {
			polyglot.Encoder(b).String(v)
		}
	}
}

func (x *HttpStringList) internalDecode(b []byte) error {
	if x == nil {
		return NilDecode
	}
	d := polyglot.GetDecoder(b)
	defer d.Return()
	return x.decode(d)
}

func (x *HttpStringList) decode(d *polyglot.Decoder) error {
	if d.Nil() {
		return nil
	}

	var err error

	var sliceSize uint32
	sliceSize, err = d.Slice(polyglot.StringKind)
	if err != nil {
		return err
	}
	if uint32(len(x.Value)) != sliceSize {
		x.Value = make([]string, sliceSize)
	}
	for i := uint32(0); i < sliceSize; i++ {
		x.Value[i], err = d.String()
		if err != nil {
			return err
		}
	}
	return nil
}
