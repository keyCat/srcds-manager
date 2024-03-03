package utils

func SliceIncludesInt(arr []int, value int) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == value {
			return true
		}
	}
	return false
}
