package utils

func BinaryConverter(number int, bits int) []int {
	result := make([]int, bits)

	for i := bits - 1; i >= 0; i-- {
		result[i] = number % 2
		number /= 2
	}

	return result
}
