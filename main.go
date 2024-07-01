package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type Address struct {
	Cep        string
	Logradouro string
	Bairro     string
	Cidade     string
	Estado     string
	Service    string
}

func main() {
	c1 := make(chan Address)
	c2 := make(chan Address)

	go func() {
		req, err := http.NewRequest(http.MethodGet, "https://brasilapi.com.br/api/cep/v1/80035050", nil)
		if err != nil {
			panic(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var brasilAPI BrasilAPI
		if err := json.NewDecoder(resp.Body).Decode(&brasilAPI); err != nil {
			panic(err)
		}

		address := parseBrasilAPI(brasilAPI)

		c1 <- address
	}()

	go func() {
		req, err := http.NewRequest(http.MethodGet, "http://viacep.com.br/ws/80035050/json/", nil)
		if err != nil {
			panic(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var viaCEP ViaCEP
		if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
			panic(err)
		}

		address := parseViaCEP(viaCEP)

		c2 <- address
	}()

	select {
	case msg := <-c1:
		printAddress(msg)
	case msg := <-c2:
		printAddress(msg)
	case <-time.After(time.Second):
		panic("timeout")
	}
}

func printAddress(address Address) {
	fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nEstado: %s\nServiço: %s\n", address.Cep, address.Logradouro, address.Bairro, address.Cidade, address.Estado, address.Service)
}

func parseViaCEP(viaCEP ViaCEP) Address {
	return Address{
		Cep:        viaCEP.Cep,
		Logradouro: viaCEP.Logradouro,
		Bairro:     viaCEP.Bairro,
		Cidade:     viaCEP.Localidade,
		Estado:     viaCEP.Uf,
		Service:    "viaCEP",
	}
}

func parseBrasilAPI(brasilAPI BrasilAPI) Address {
	return Address{
		Cep:        brasilAPI.Cep,
		Logradouro: brasilAPI.Street,
		Bairro:     brasilAPI.Neighborhood,
		Cidade:     brasilAPI.City,
		Estado:     brasilAPI.State,
		Service:    "brasilAPI",
	}
}
