package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateCep(cep string) (string, error) {
	regex := regexp.MustCompile(`\d{2}\.?\d{3}-?\d{3}`)
	if !regex.MatchString(cep) {
		return "", fmt.Errorf("invalid zipcode")
	}
	parsedCep := strings.ReplaceAll(cep, ".", "")
	parsedCep = strings.ReplaceAll(parsedCep, "-", "")
	return parsedCep, nil
}
