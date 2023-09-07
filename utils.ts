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

export function PackUint32(ptr: number, len: number): bigint {
	if (ptr > 0xffffffff || len > 0xffffffff) {
		throw new Error("ptr or len is too large");
	}
	return (BigInt(ptr) << BigInt(32)) | BigInt(len);
}

// Unpack a memory ref from 64bit to 2x32bits
export function UnpackUint32(packed: bigint): [number, number] {
	const ptr = Number((packed >> BigInt(32)) & BigInt(0xffffffff));
	const len = Number(packed & BigInt(0xffffffff));
	return [ptr, len];
}