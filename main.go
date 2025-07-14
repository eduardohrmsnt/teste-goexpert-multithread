package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type CepBrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Service      string `json:"service"`
	CepResponse
}

type CepViaCepResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Service      string `json:"unidade"`
	CepResponse
}

type CepResponse struct {
	ApiRetorno string `json:"apiRetorno"`
}

func main() {
	router := chi.NewRouter()

	router.Get("/obter-cep/{cep}", ObterCep)

	http.ListenAndServe(":8080", router)
}

func ObterCep(writer http.ResponseWriter, req *http.Request) {
	cep := chi.URLParam(req, "cep")

	responseBrasilApi := make(chan CepBrasilApiResponse)
	responseViaCep := make(chan CepViaCepResponse)

	go func() {
		res, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)

		if err != nil {
			println("error")
			return
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			println("error reading body:", err.Error())
			return
		}

		var response CepBrasilApiResponse
		json.Unmarshal(body, &response)

		response.ApiRetorno = "BrasilApi"
		responseBrasilApi <- response
	}()

	go func() {
		res, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

		if err != nil {
			println("error")
			return
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			println("error reading body:", err.Error())
			return
		}

		var response CepViaCepResponse
		json.Unmarshal(body, &response)

		response.ApiRetorno = "ViaCep"
		responseViaCep <- response
	}()

	select {
	case response1 := <-responseBrasilApi:
		fmt.Printf("%v Brasil API", response1)
		response1.ApiRetorno = "Brasil API"
		returnResponse(writer, response1)

	case response2 := <-responseViaCep:
		fmt.Printf("%v ViaCep", response2)
		response2.ApiRetorno = "ViaCep"
		returnResponse(writer, response2)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func returnResponse(writer http.ResponseWriter, response any) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}
