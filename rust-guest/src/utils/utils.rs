pub fn pack_uint32(ptr: u32, len: u32) -> u64 {
    return ((ptr as u64) << 32) | len as u64;
}

pub fn unpack_uint32(packed: u64) -> (u32, u32) {
	return ((packed >> 32) as u32, packed as u32)
}
