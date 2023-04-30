package main

import (
	"api/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func vueloSender() {
	// Read vuelos.json file
	file, err := ioutil.ReadFile("vuelos.json")
	if err != nil {
		panic(err)
	}

	// Parse JSON array
	var vuelos []model.Vuelo
	err = json.Unmarshal(file, &vuelos)
	if err != nil {
		panic(err)
	}

	// Iterate over JSON objects and send HTTP POST request
	for _, vuelo := range vuelos {
		// Convert vuelo to JSON string
		vueloJSON, err := json.Marshal(vuelo)
		if err != nil {
			panic(err)
		}

		// Create HTTP POST request
		req, err := http.NewRequest("POST", "http://localhost:5000/api/vuelo", bytes.NewBuffer(vueloJSON))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send HTTP POST request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// Print HTTP response status code
		fmt.Println(resp.StatusCode)
	}
}

func reservaSender() {
	// Read vuelos.json file
	file, err := ioutil.ReadFile("reservas.json")
	if err != nil {
		panic(err)
	}

	// Parse JSON array
	var reservas []model.Reserva
	err = json.Unmarshal(file, &reservas)
	if err != nil {
		panic(err)
	}

	// Iterate over JSON objects and send HTTP POST request
	for _, reserva := range reservas {
		// Convert reserva to JSON string
		reservaJSON, err := json.Marshal(reserva)
		if err != nil {
			panic(err)
		}

		// Create HTTP POST request
		req, err := http.NewRequest("POST", "http://localhost:5000/api/reserva", bytes.NewBuffer(reservaJSON))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send HTTP POST request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func main() {
	vueloSender()
	reservaSender()
}
