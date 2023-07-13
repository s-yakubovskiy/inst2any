package instagram

func reverseSlice(slice []string) []string {
	newSlice := make([]string, len(slice))
	copy(newSlice, slice)

	for i, j := 0, len(newSlice)-1; i < j; i, j = i+1, j-1 {
		newSlice[i], newSlice[j] = newSlice[j], newSlice[i]
	}
	return newSlice
}
