package utils

func Contains(slice []string, item string) bool {
	for _, el := range slice {
		if el == item {
			return true
		}
	}

	return false
}
