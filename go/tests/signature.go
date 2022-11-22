package tests

//import (
//	"errors"
//	"github.com/loopholelabs/polyglot-go"
//	signature "github.com/loopholelabs/scale-signature"
//	"unsafe"
//)
//
//var (
//	NilDecode = errors.New("cannot decode into a nil root struct")
//)
//
//type TestContext struct {
//	Data string
//}
//
//func NewTestContext() *TestContext {
//	return &TestContext{}
//}
//
//func (x *TestContext) error(b *polyglot.Buffer, err error) {
//	polyglot.Encoder(b).Error(err)
//}
//
//func (x *TestContext) internalEncode(b *polyglot.Buffer) {
//	if x == nil {
//		polyglot.Encoder(b).Nil()
//	} else {
//		polyglot.Encoder(b).String(x.Data)
//	}
//}
//
//func (x *TestContext) internalDecode(b []byte) error {
//	if x == nil {
//		return NilDecode
//	}
//	d := polyglot.GetDecoder(b)
//	defer d.Return()
//	return x.decode(d)
//}
//
//func (x *TestContext) decode(d *polyglot.Decoder) error {
//	if d.Nil() {
//		return nil
//	}
//
//	var err error
//	x.Data, err = d.String()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//type HttpRequestHeadersMap map[string]*HttpStringList
//
//func NewHttpRequestHeadersMap(size uint32) map[string]*HttpStringList {
//	return make(map[string]*HttpStringList, size)
//}
//
//func (x HttpRequestHeadersMap) internalEncode(b *polyglot.Buffer) {
//	if x == nil {
//		polyglot.Encoder(b).Nil()
//	} else {
//		polyglot.Encoder(b).Map(uint32(len(x)), polyglot.StringKind, polyglot.AnyKind)
//		for k, v := range x {
//			polyglot.Encoder(b).String(k)
//			v.internalEncode(b)
//		}
//	}
//}
//
//func (x HttpRequestHeadersMap) decode(d *polyglot.Decoder, size uint32) error {
//	if size == 0 {
//		return nil
//	}
//	var k string
//	var v *HttpStringList
//	var err error
//	for i := uint32(0); i < size; i++ {
//		k, err = d.String()
//		if err != nil {
//			return err
//		}
//		v = NewHttpStringList()
//		err = v.decode(d)
//		if err != nil {
//			return err
//		}
//		x[k] = v
//	}
//	return nil
//}
//
//type HttpRequest struct {
//	Headers       HttpRequestHeadersMap
//	URI           string
//	Method        string
//	ContentLength int64
//	Protocol      string
//	IP            string
//	Body          []byte
//}
//
//func NewHttpRequest() *HttpRequest {
//	return &HttpRequest{}
//}
//
//func (x *HttpRequest) error(b *polyglot.Buffer, err error) {
//	polyglot.Encoder(b).Error(err)
//}
//
//func (x *HttpRequest) internalEncode(b *polyglot.Buffer) {
//	if x == nil {
//		polyglot.Encoder(b).Nil()
//	} else {
//		polyglot.Encoder(b).String(x.URI).String(x.Method).Int64(x.ContentLength).String(x.Protocol).String(x.IP).Bytes(x.Body)
//		x.Headers.internalEncode(b)
//	}
//}
//
//func (x *HttpRequest) internalDecode(b []byte) error {
//	if x == nil {
//		return NilDecode
//	}
//	d := polyglot.GetDecoder(b)
//	defer d.Return()
//	return x.decode(d)
//}
//
//func (x *HttpRequest) decode(d *polyglot.Decoder) error {
//	if d.Nil() {
//		return nil
//	}
//
//	var err error
//	x.URI, err = d.String()
//	if err != nil {
//		return err
//	}
//	x.Method, err = d.String()
//	if err != nil {
//		return err
//	}
//	x.ContentLength, err = d.Int64()
//	if err != nil {
//		return err
//	}
//	x.Protocol, err = d.String()
//	if err != nil {
//		return err
//	}
//	x.IP, err = d.String()
//	if err != nil {
//		return err
//	}
//	x.Body, err = d.Bytes(nil)
//	if err != nil {
//		return err
//	}
//	if !d.Nil() {
//		HeadersSize, err := d.Map(polyglot.StringKind, polyglot.AnyKind)
//		if err != nil {
//			return err
//		}
//		x.Headers = NewHttpRequestHeadersMap(HeadersSize)
//		err = x.Headers.decode(d, HeadersSize)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//type HttpResponseHeadersMap map[string]*HttpStringList
//
//func NewHttpResponseHeadersMap(size uint32) map[string]*HttpStringList {
//	return make(map[string]*HttpStringList, size)
//}
//
//func (x HttpResponseHeadersMap) internalEncode(b *polyglot.Buffer) {
//	if x == nil {
//		polyglot.Encoder(b).Nil()
//	} else {
//		polyglot.Encoder(b).Map(uint32(len(x)), polyglot.StringKind, polyglot.AnyKind)
//		for k, v := range x {
//			polyglot.Encoder(b).String(k)
//			v.internalEncode(b)
//		}
//	}
//}
//
//func (x HttpResponseHeadersMap) decode(d *polyglot.Decoder, size uint32) error {
//	if size == 0 {
//		return nil
//	}
//	var k string
//	var v *HttpStringList
//	var err error
//	for i := uint32(0); i < size; i++ {
//		k, err = d.String()
//		if err != nil {
//			return err
//		}
//		v = NewHttpStringList()
//		err = v.decode(d)
//		if err != nil {
//			return err
//		}
//		x[k] = v
//	}
//	return nil
//}
//
//type HttpResponse struct {
//	Headers    HttpResponseHeadersMap
//	StatusCode int32
//	Body       []byte
//}
//
//func NewHttpResponse() *HttpResponse {
//	return &HttpResponse{}
//}
//
//func (x *HttpResponse) error(b *polyglot.Buffer, err error) {
//	polyglot.Encoder(b).Error(err)
//}
//
//func (x *HttpResponse) internalEncode(b *polyglot.Buffer) {
//	if x == nil {
//		polyglot.Encoder(b).Nil()
//	} else {
//		polyglot.Encoder(b).Int32(x.StatusCode).Bytes(x.Body)
//		x.Headers.internalEncode(b)
//	}
//}
//
//func (x *HttpResponse) internalDecode(b []byte) error {
//	if x == nil {
//		return NilDecode
//	}
//	d := polyglot.GetDecoder(b)
//	defer d.Return()
//	return x.decode(d)
//}
//
//func (x *HttpResponse) decode(d *polyglot.Decoder) error {
//	if d.Nil() {
//		return nil
//	}
//
//	var err error
//	x.StatusCode, err = d.Int32()
//	if err != nil {
//		return err
//	}
//	x.Body, err = d.Bytes(nil)
//	if err != nil {
//		return err
//	}
//	if !d.Nil() {
//		HeadersSize, err := d.Map(polyglot.StringKind, polyglot.AnyKind)
//		if err != nil {
//			return err
//		}
//		x.Headers = NewHttpResponseHeadersMap(HeadersSize)
//		err = x.Headers.decode(d, HeadersSize)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//type HttpStringList struct {
//	Value []string
//}
//
//func NewHttpStringList() *HttpStringList {
//	return &HttpStringList{}
//}
//
//func (x *HttpStringList) error(b *polyglot.Buffer, err error) {
//	polyglot.Encoder(b).Error(err)
//}
//
//func (x *HttpStringList) internalEncode(b *polyglot.Buffer) {
//	if x == nil {
//		polyglot.Encoder(b).Nil()
//	} else {
//
//		polyglot.Encoder(b).Slice(uint32(len(x.Value)), polyglot.StringKind)
//		for _, v := range x.Value {
//			polyglot.Encoder(b).String(v)
//		}
//	}
//}
//
//func (x *HttpStringList) internalDecode(b []byte) error {
//	if x == nil {
//		return NilDecode
//	}
//	d := polyglot.GetDecoder(b)
//	defer d.Return()
//	return x.decode(d)
//}
//
//func (x *HttpStringList) decode(d *polyglot.Decoder) error {
//	if d.Nil() {
//		return nil
//	}
//
//	var err error
//	var sliceSize uint32
//	sliceSize, err = d.Slice(polyglot.StringKind)
//	if err != nil {
//		return err
//	}
//	if uint32(len(x.Value)) != sliceSize {
//		x.Value = make([]string, sliceSize)
//	}
//	for i := uint32(0); i < sliceSize; i++ {
//		x.Value[i], err = d.String()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//
//var _ signature.RuntimeContext = (*RuntimeContext)(nil)
//var _ signature.GuestContext = (*GuestContext)(nil)
//var _ signature.Signature = (*Context)(nil)
//var _ signature.Context = (*Context)(nil)
//
//var (
//	writeBuffer = polyglot.NewBuffer()
//	readBuffer  []byte
//)
//
//// Context is a context object for an incoming request. It is meant to be used
//// inside the Scale function.
//type Context struct {
//	generated *
//	buffer    *polyglot.Buffer
//}
//
//type GuestContext Context
//type RuntimeContext Context
//
//// New creates a new empty Context
//func New() *Context {
//	return &Context{
//		generated: NewHttpContext(),
//		buffer:    polyglot.NewBuffer(),
//	}
//}
//
//func (x *Context) GuestContext() signature.GuestContext {
//	return (*GuestContext)(x)
//}
//
//func (x *Context) RuntimeContext() signature.RuntimeContext {
//	return (*RuntimeContext)(x)
//}
//
//func (x *Context) Version() string {
//	return signatureFile.Version
//}
//
//func (x *Context) Name() string {
//	return signatureFile.Name
//}
//
//func (x *Context) Resize(size uint32) uint32 {
//	return Resize(size)
//}
//
//// ToWriteBuffer serializes the Context into the global writeBuffer and returns the pointer to the buffer and its size
////
//// This method should only be used to read the Context from the Scale Runtime.
//// Users should not use this method.
//func (x *GuestContext) ToWriteBuffer() (uint32, uint32) {
//	writeBuffer.Reset()
//	x.generated.internalEncode(writeBuffer)
//	underlying := writeBuffer.Bytes()
//	ptr := &underlying[0]
//	unsafePtr := uintptr(unsafe.Pointer(ptr))
//	return uint32(unsafePtr), uint32(writeBuffer.Len())
//}
//
//// FromReadBuffer deserializes the data into the Context from the global readBuffer
////
//// It assumes that the readBuffer has been filled with the data from the Scale Runtime after
//// a call to the Resize method
//func (x *GuestContext) FromReadBuffer() error {
//	return x.generated.internalDecode(readBuffer)
//}
//
//// Error serializes an error into the global writeBuffer and returns a pointer to the buffer and its size
////
//// This method should only be used to write an error to the Scale Runtime, in place of the ToWriteBuffer method.
//// Users should not use this method.
//func (x *GuestContext) ErrorWriteBuffer(err error) (uint32, uint32) {
//	writeBuffer.Reset()
//	x.generated.error(writeBuffer, err)
//	underlying := writeBuffer.Bytes()
//	ptr := &underlying[0]
//	unsafePtr := uintptr(unsafe.Pointer(ptr))
//	return uint32(unsafePtr), uint32(writeBuffer.Len())
//}
//
//// Read reads the context from the given byte slice and returns an error if one occurred
////
//// This method is meant to be used by the Scale Runtime to deserialize the Context
//func (x *RuntimeContext) Read(b []byte) error {
//	return x.generated.internalDecode(b)
//}
//
//// Write writes the context into a byte slice and returns it
//func (x *RuntimeContext) Write() []byte {
//	x.buffer.Reset()
//	x.generated.internalEncode(x.buffer)
//	return x.buffer.Bytes()
//}
//
//// Next calls the next host function after writing the Context into the global writeBuffer,
//// then it reads the result from the global readBuffer back into the Context
//func (x *Context) Next() (*Context, error) {
//	next(x.GuestContext().ToWriteBuffer())
//	return x, x.GuestContext().FromReadBuffer()
//}
//
//// Generated is not meant to be used directly. It is meant to be used by the Scale Runtime.
//func (x *Context) Generated() *HttpContext {
//	return x.generated
//}
