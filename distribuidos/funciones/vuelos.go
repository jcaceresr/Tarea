package funciones

import (
	"api/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func elegir_vuelo(fechaIda string, fechaVuelta string, origen string, destino string) ([]model.Vuelo, int, int) {
	// Vuelos de ida
	fmt.Println("Ida:")
	vueloIda, precioPasajeIda, err := chooseVuelo(origen, destino, fechaIda)
	if err != nil {
		fmt.Println("No se encontraron vuelos de ida, intente con otra fecha")
		return nil, 0, 0
	}

	// Vuelos de vuelta
	var vueloVuelta *model.Vuelo
	var precioPasajeVuelta int

	if fechaVuelta != "no" {
		fmt.Println("Vuelta:")
		vueloVuelta, precioPasajeVuelta, err = chooseVuelo(destino, origen, fechaVuelta)
		if err != nil {
			fmt.Println("No se encontraron vuelos de vuelta, se reservará solo el vuelo de ida")
		}
	}

	var vuelosReserva []model.Vuelo
	vuelosReserva = append(vuelosReserva, *vueloIda)
	if vueloVuelta != nil {
		vuelosReserva = append(vuelosReserva, *vueloVuelta)
	}

	return vuelosReserva, precioPasajeIda, precioPasajeVuelta
}

func chooseVuelo(origen string, destino string, fecha string) (*model.Vuelo, int, error) {
	// Obtener vuelos
	queries := map[string]string{
		"origen":  origen,
		"destino": destino,
		"fecha":   fecha,
	}

	url := crear_url("vuelo", queries)
	vuelos := requestVuelos(url)
	var vuelo model.Vuelo
	var choice int

	if len(vuelos) > 0 {
		// Mostrar opciones
		mostrar_vuelo(vuelos)

		fmt.Println("Ingrese una opción: ")
		fmt.Scanln(&choice)
	} else {
		return nil, 0, fmt.Errorf("no se encontraron vuelos")
	}

	vuelo = vuelos[choice-1]
	precioPasaje := valor_pasaje(vuelo)

	return &vuelo, precioPasaje, nil
}

func requestVuelos(url string) []model.Vuelo {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	var vuelos []model.Vuelo

	if err := json.Unmarshal(body, &vuelos); err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	return vuelos
}

func mostrar_vuelo(vuelos []model.Vuelo) {
	for _, vuelo := range vuelos {
		precioPasaje := valor_pasaje(vuelo)

		fmt.Printf("%s %s %s $%d\n", vuelo.Numero_vuelo, vuelo.Hora_salida, vuelo.Hora_llegada, precioPasaje)
	}
}
