package utils

func ContainsUint64(list []uint64, x uint64) bool {
	for _, e := range list {
		if e == x {
			return true
		}
	}
	return false
}

func RemoveFromUint64Slice(list []uint64, x uint64) []uint64 {
	for i, e := range list {
		if e == x {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}
