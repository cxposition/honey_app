package utils

func Inlist[T comparable](list []T, key T) bool {
	for _, t := range list {
		if t == key {
			return true
		}
	}
	return false
}
