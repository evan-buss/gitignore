package main

// arrayContains determines if the given []string array contains the given value
func arrayContains(array []string, value string) bool {
	for _, val := range array {
		if val == value {
			return true
		}
	}
	return false
}
