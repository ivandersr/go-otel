package cep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ivandersr/go-weather/pkg/utils"
)

type CepResult struct {
	City  string
	State string
}

type Cep struct {
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
}

func GetLocation(cepURL, cep string) (*CepResult, error) {
	url := utils.InjectArg(cep, cepURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var c Cep
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	city := utils.SanitizeInput(c.Localidade)
	state := utils.SanitizeInput(c.Estado)

	if city == "" || state == "" {
		return nil, fmt.Errorf("não foi possível encontrar a cidade com o CEP informado")
	}
	return &CepResult{
		City:  city,
		State: state,
	}, nil
}
