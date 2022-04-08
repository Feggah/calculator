package utils

var Operators = []byte{'+', '-', '*', '/', '%'}

func ByteInSlice(char byte, list []byte) bool {
	for _, b := range list {
		if b == char {
			return true
		}
	}
	return false
}
