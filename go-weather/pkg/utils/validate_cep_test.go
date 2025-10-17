package utils_test

import (
	"strings"
	"testing"

	"github.com/ivandersr/go-weather/pkg/utils"
)

func TestValidateCep(t *testing.T) {
	t.Run("should validate cep", func(t *testing.T) {
		cep := "12345-678"
		expected := "12345678"
		parsedCep, err := utils.ValidateCep(cep)
		if err != nil {
			t.Errorf("expected %s but got %s", expected, parsedCep)
		}
		if parsedCep != expected {
			t.Errorf("expected %s but got %s", expected, parsedCep)
		}
	})

	t.Run("should return error when cep is invalid", func(t *testing.T) {
		cep := "1234-5678"
		_, err := utils.ValidateCep(cep)
		if err == nil {
			t.Errorf("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "invalid zipcode") {
			t.Errorf("expected error message to contain 'invalid zipcode' but got %s", err.Error())
		}
	})

	t.Run("should return error when cep is empty string", func(t *testing.T) {
		cep := ""
		_, err := utils.ValidateCep(cep)
		if err == nil {
			t.Errorf("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "invalid zipcode") {
			t.Errorf("expected error message to contain 'invalid zipcode' but got %s", err.Error())
		}
	})
}
