/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

// Package signature implements the Signature type, that must be exported by Signatures
package signature

type Signature interface {
	Version() string                // Version of the Signature
	RuntimeContext() RuntimeContext // RuntimeContext of the Signature
	Resize(size uint32) uint32      // Resize resizes the global buffers to the given size and returns the offset
}

// RuntimeContext is the interface that must be implemented by the Context of a Signature
// in order for it to be used by the runtime.
type RuntimeContext interface {
	Read(b []byte) error // Read updates the Context by decoding the given bytes
	Write() []byte       // Write encodes the Context and returns the encoded bytes
}

// GuestContext is the interface that must be implemented by the Context of a Signature
// in order for it to be used by the guest.
type GuestContext interface {
	ToWriteBuffer() (uint32, uint32) // ToWriteBuffer serializes the Context to a global buffer and returns the offset and length
	FromReadBuffer() error           // FromReadBuffer deserializes the Context from the global buffer
}
