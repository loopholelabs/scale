package context

//export next
//go:linkname next
func next(offset uint32, byteCount uint32)

//export modifyGlobal
func modifyGlobal(offset uint32, byteCount uint32) {
	globalOffset = offset
	globalLength = byteCount
}
