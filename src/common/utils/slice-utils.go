package app_utils

import "reflect"

func Includes[T any](slice []T, el T) bool {
	for i := 0; i < len(slice); i++ {
		elememt := slice[i]
		reflectEl := reflect.ValueOf(el)
		reflectSliceEl := reflect.ValueOf(elememt)
		if reflectEl.Interface() == reflectSliceEl.Interface() {
			return true
		}
	}
	return false
}
