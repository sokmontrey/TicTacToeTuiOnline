package pkg

func SliceRemoveByIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func SliceRemoveByValue[T comparable](slice []T, value T) []T {
	for i, v := range slice {
		if v == value {
			return SliceRemoveByIndex(slice, i)
		}
	}
	return slice
}
