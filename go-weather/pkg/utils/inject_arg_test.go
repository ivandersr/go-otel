package utils_test

import (
	"testing"

	"github.com/ivandersr/go-weather/pkg/utils"
)

func TestInjectArg(t *testing.T) {
	t.Run("should inject arg", func(t *testing.T) {
		arg := "12345678"
		template := "https://viacep.com.br/ws/{arg}/json/"
		expected := "https://viacep.com.br/ws/12345678/json/"
		result := utils.InjectArg(arg, template)
		if result != expected {
			t.Errorf("expected %s but got %s", expected, result)
		}
	})
}
