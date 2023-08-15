package utils

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

// Returns false if passed string is empty, begins/ends in whitespace
func ValidateBasic(s string) error {

	//ensure that the string is not " " or "  "
	trimmed := strings.TrimSpace(s)

	// ensure name is not empty
	if s == "" || len(trimmed) == 0 {
		return errors.New("cannot be empty or contain only whitespace")
	}

	// Ensure name of blacklist is not beginning or ending in whitespace
	for i, character := range s {

		if unicode.IsSpace(character) {
			if i == 0 || i == (len(s)-1) {
				return errors.New("cannot begin or end in whitespace")
			}
		}
	}

	return nil
}

func ValidateURLforIPAddr(url string) bool {
	response, err := http.Get(url)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	// Extract IP addresses from the response body
	ipList := strings.Split(string(body), "\n")

	for _, ipAddress := range ipList {
		if len(ipAddress) == 0 {
			continue
		}
		if len(ipAddress) > 0 && string(ipAddress[0]) == "#" {
			continue
		}
		if net.ParseIP(ipAddress) == nil {
			return false
		}
	}
	return true
}

func ValidEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}

	return match
}
