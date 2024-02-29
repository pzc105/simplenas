package distributedchat

func isSet(slots []byte, i int) bool {
	return (slots[i/8] & byte(1<<(i%8))) > 0
}

func setSlot(slots []byte, i int) {
	slots[i/8] |= byte(1 << (i % 8))
}

func unsetSlot(slots []byte, i int) {
	slots[i/8] &= ^byte(1 << (i % 8))
}
