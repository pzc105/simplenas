package utils

func GetPow2_32(hint int) int {
	hint--
	hint |= (hint >> 1)
	hint |= (hint >> 2)
	hint |= (hint >> 4)
	hint |= (hint >> 8)
	hint |= (hint >> 16)
	hint++
	return hint
}

func GetPow2_64(hint int64) int64 {
	hint--
	hint |= (hint >> 1)
	hint |= (hint >> 2)
	hint |= (hint >> 4)
	hint |= (hint >> 8)
	hint |= (hint >> 16)
	hint |= (hint >> 32)
	hint++
	return hint
}
