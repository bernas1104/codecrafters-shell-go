package utils

func Contains(slice []string, item string) bool {
	for _, el := range slice {
		if el == item {
			return true
		}
	}

	return false
}

func AnyMultiple(slice []string, items []string) bool {
	for _, el := range slice {
		for _, it := range items {
			if el == it {
				return true
			}
		}
	}

	return false
}

func FindRedirectIndex(slice []string) int {
	for idx, el := range slice {
		if el == ">" || el == "1>" {
			return idx
		}
	}

	return -1
}
