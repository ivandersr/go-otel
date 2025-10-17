package utils_test

import (
	"testing"

	"github.com/ivandersr/go-weather/pkg/utils"
)

func TestSanitizeInput(t *testing.T) {
	t.Run("should remove accents", func(t *testing.T) {
		input := "água"
		expected := "agua"
		result := utils.SanitizeInput(input)
		if result != expected {
			t.Errorf("expected %s but got %s", expected, result)
		}
	})

	t.Run("should replace special characters", func(t *testing.T) {
		input := "água,!"
		expected := "agua%2C%21"
		result := utils.SanitizeInput(input)
		if result != expected {
			t.Errorf("expected %s but got %s", expected, result)
		}
	})

	t.Run("should remove spaces", func(t *testing.T) {
		input := "água doce"
		expected := "agua+doce"
		result := utils.SanitizeInput(input)
		if result != expected {
			t.Errorf("expected %s but got %s", expected, result)
		}
	})
}
