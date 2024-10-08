// Code generated by scale-extension v0.4.8, DO NOT EDIT.
// output: local_inttest_latest_guest

package local_inttest_latest_guest

import (
	"github.com/loopholelabs/polyglot"
	"unsafe"
)

var (
	writeBuffer = polyglot.NewBuffer()
	readBuffer  []byte
)

//export ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Resize
//go:linkname ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Resize
func ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Resize(size uint32) uint32 {
	readBuffer = make([]byte, size)
	//if uint32(cap(readBuffer)) < size {
	//	readBuffer = append(make([]byte, 0, uint32(len(readBuffer))+size), readBuffer...)
	//}
	//readBuffer = readBuffer[:size]
	return uint32(uintptr(unsafe.Pointer(&readBuffer[0])))
}

// Define any interfaces we need here...
// Also define structs we can use to hold instanceId

// Define concrete types with a hidden instanceId

type _Example struct {
	instanceId uint64
}

func (d *_Example) Hello(params *Stringval) (Stringval, error) {

	// First we take the params, serialize them.
	writeBuffer.Reset()
	params.Encode(writeBuffer)
	underlying := writeBuffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	off := uint32(unsafePtr)
	l := uint32(writeBuffer.Len())

	// Now make the call to the host.
	ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Example_Hello(d.instanceId, off, l)
	// IF the return type is a model, we should read the data from the read buffer.

	ret := &Stringval{}
	r, err := DecodeStringval(ret, readBuffer)

	if err != nil {
		return Stringval{}, err
	}

	return *r, err

}

//export ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Example_Hello
//go:linkname ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Example_Hello
func ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_Example_Hello(instance uint64, offset uint32, length uint32) uint64

// Define any global functions here...

//export ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_New
//go:linkname ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_New
func ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_New(instance uint64, offset uint32, length uint32) uint64

func New(params *Stringval) (Example, error) {
	// First we take the params, serialize them.
	writeBuffer.Reset()
	params.Encode(writeBuffer)
	underlying := writeBuffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	off := uint32(unsafePtr)
	l := uint32(writeBuffer.Len())

	// Now make the call to the host.
	readBuffer = nil
	v := ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_New(0, off, l)
	// IF the return type is an interface return ifc, which contains hidden instanceId.

	// Handle error from host. In this case there'll be an error in the readBuffer
	if readBuffer != nil {
		val, err := polyglot.GetDecoder(readBuffer).Error()
		if err != nil {
			panic(err)
		}
		return nil, val
	}

	ret := &_Example{
		instanceId: v,
	}

	return ret, nil

}

//export ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_World
//go:linkname ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_World
func ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_World(instance uint64, offset uint32, length uint32) uint64

func World(params *Stringval) (Stringval, error) {
	// First we take the params, serialize them.
	writeBuffer.Reset()
	params.Encode(writeBuffer)
	underlying := writeBuffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	off := uint32(unsafePtr)
	l := uint32(writeBuffer.Len())

	// Now make the call to the host.
	ext_b30af2dd8561988edd7b281ad5c1b84487072727a8ad0e490a87be0a66b037d7_World(0, off, l)
	// IF the return type is a model, we should read the data from the read buffer.

	ret := &Stringval{}
	r, err := DecodeStringval(ret, readBuffer)

	if err != nil {
		return Stringval{}, err
	}

	return *r, err

}

// Error serializes an error into the global writeBuffer and returns a pointer to the buffer and its size
//
// Users should not use this method.
func Error(err error) (uint32, uint32) {
	writeBuffer.Reset()
	polyglot.Encoder(writeBuffer).Error(err)
	underlying := writeBuffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(writeBuffer.Len())
}
