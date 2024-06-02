package utils

import (
	"errors"
	"fmt"
)

func Error(msg string, funcName string) error {
	return errors.New(fmt.Sprintf("Error in %v: %v", funcName, msg))
}

func MethodNotImplemented() error {
	return errors.New("Method is not implemented!")
}

func UnmarshalError(msg string, funcName string) error {
	return errors.New(fmt.Sprintf("Unmarshal error in %v: %v", funcName, msg))
}

func InvalidAiResponseFormat() error {
	return errors.New(fmt.Sprintf("Can't parse content to string", "getSongsListFromResponse"))
}

func IndexOfSubstring(str, subStr string) int {
	for i := 0; i < len(str); i++ {
		if str[i:i+len(subStr)] == subStr {
			return i
		}
	}
	return -1
}
