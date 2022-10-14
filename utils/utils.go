package utils

func PackUint32(offset uint32, length uint32) uint64 {
	return uint64(offset)<<32 | uint64(length)
}

func UnpackUint32(packed uint64) (uint32, uint32) {
	return uint32(packed >> 32), uint32(packed)
}
