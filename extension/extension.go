package extension

type ModuleMemory interface {
	Write(offset uint32, v []byte) bool
	Read(offset uint32, byteCount uint32) ([]byte, bool)
}

type Resizer func(name string, size uint64) (uint64, error)

type InstallableFunc func(mem ModuleMemory, resize Resizer, params []uint64)
