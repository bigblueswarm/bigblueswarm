// Package utils provide few utilies functions
package utils

// ArrayContainsString checks if a string is in an array
func ArrayContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
