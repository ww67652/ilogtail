package utils

func ContainElement[T, E comparable](arr []T, val E, eqCondition func(T, E) bool) bool {
	if arr == nil || len(arr) == 0 {
		return false
	}
	for _, element := range arr {
		if eqCondition != nil && eqCondition(element, val) {
			return true
		}
	}
	return false
}

func RemoveElement[T comparable](arr []T, val T, eqCondition func(T, T) bool) (bool, []T) {
	if arr == nil || len(arr) == 0 {
		return false, arr
	}
	for i, element := range arr {
		if element == val || eqCondition != nil && eqCondition(element, val) {
			return true, append(arr[:i], arr[i+1:]...)
		}
	}
	return false, arr
}

// ReplaceElement 如果存在就替换，不存在则返回原数组
func ReplaceElement[T comparable](arr []T, val T, eqCondition func(T, any) bool) {
	if eqCondition == nil {
		return
	}
	if arr == nil || len(arr) == 0 {
		return
	}
	for i, element := range arr {
		if eqCondition(element, val) {
			arr[i] = val
		}
	}
}

func Filter[T comparable](arr []T, fn func(T) bool) []T {
	finalArr := make([]T, 0)
	for _, element := range arr {
		if fn(element) {
			finalArr = append(finalArr, element)
		}
	}
	return finalArr
}

func Map[T, E any](arr []T, fn func(T) E) []E {
	finalArr := make([]E, 0)
	for _, element := range arr {
		finalArr = append(finalArr, fn(element))
	}
	return finalArr
}
