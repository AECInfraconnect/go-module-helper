package helper

// InArray checks if a value exists in an array and returns its position
func InArray(needle interface{}, haystack interface{}) (bool, int) {
	switch arr := haystack.(type) {
	case []string:
		if str, ok := needle.(string); ok {
			for i, v := range arr {
				if v == str {
					return true, i
				}
			}
		}
	case []int:
		if num, ok := needle.(int); ok {
			for i, v := range arr {
				if v == num {
					return true, i
				}
			}
		}
	case []interface{}:
		for i, v := range arr {
			if v == needle {
				return true, i
			}
		}
	}
	return false, -1
}

// Contains checks if a string exists in a string slice
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
