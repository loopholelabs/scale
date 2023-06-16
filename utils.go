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

package scale

func packUint32(offset uint32, length uint32) uint64 {
	return uint64(offset)<<32 | uint64(length)
}

func unpackUint32(packed uint64) (uint32, uint32) {
	return uint32(packed >> 32), uint32(packed)
}
