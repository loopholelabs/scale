package context

//export next
//go:linkname next
func next(offset uint32, length uint32) (packed uint64)
