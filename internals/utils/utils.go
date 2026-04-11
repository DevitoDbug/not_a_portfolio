// Package utils - has all the utility functions used in all other packages
// Does not import from other packages
package utils

func GetAllowedOrigins(environment string) []string {
	if environment == "development" {
		return []string{"http://*", "https://*"}
	}

	return []string{"https://*"}
}
