//go:build tinygo
// +build tinygo

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

package context

// next is the host function that is called by the runtime to execute the next
// function in the chain.
//
//export next
//go:linkname next
func next(offset uint32, length uint32) (packed uint64)

////export debug
////go:linkname debug
//func debug(offset uint32, length uint32)
