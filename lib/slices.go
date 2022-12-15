package lib

func Find[T comparable](slice []T, needle T) (int, bool) {
	if len(slice) == 0 {
		return -1, false
	}
	for i := range slice {
		if slice[i] == needle {
			return i, true
		}
	}
	return -1, false
}
