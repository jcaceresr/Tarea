package funciones

import (
	"api/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Stock struct {
	Stock int `json:"stock_de_pasajeros" bson:"stock_de_pasajeros"`
}

func Updates(pnr string, apellido string) {
	keepRunning := true

	for keepRunning {

		originalReservation := reserva_sin_print(pnr, apellido)

		if len(originalReservation.Pasajeros) == 0 && len(originalReservation.Vuelos) == 0 {
			fmt.Println("Reserva no encontrada, ¿seguro que ingresó los datos correctos?")
			break
		}

		fmt.Println(
			"Opciones:\n",
			"1. Cambiar fecha de vuelo\n",
			"2. Adicionar ancillaries\n",
			"3. Salir",
		)

		var option int
		fmt.Print("Ingrese una opción: ")
		fmt.Scanln(&option)

		if option == 3 {
			keepRunning = false
			break
		}

		if option == 1 {
			// Mostrar vuelos que se pueden cambiar
			fmt.Println("Vuelos:")
			fmt.Printf("1. ida: %s %s - %s\n", originalReservation.Vuelos[0].Numero_vuelo, originalReservation.Vuelos[0].Hora_salida, originalReservation.Vuelos[0].Hora_llegada)
			fmt.Printf("2. vuelta: %s %s - %s\n", originalReservation.Vuelos[1].Numero_vuelo, originalReservation.Vuelos[1].Hora_salida, originalReservation.Vuelos[1].Hora_llegada)

			// Elegir vuelo a cambiar
			var vueloToReplace int
			fmt.Print("Ingrese una opción: ")
			fmt.Scanln(&vueloToReplace)
			vueloToReplace -= 1

			updatedReservation, err := cambiar_vuelos(originalReservation, vueloToReplace, pnr, apellido)
			if err != nil {
				log.Fatal(err)
				continue
			}

			actualizar_reservacion(pnr, apellido, updatedReservation)
			actualizar_stock(updatedReservation, updatedReservation.Vuelos[vueloToReplace], originalReservation.Vuelos[vueloToReplace])
		}

		if option == 2 {
			var choice string
			fmt.Print("Ingrese el PNR: ")
			fmt.Scanln(&choice)
			fmt.Print("Ingrese el apellido: ")
			fmt.Scanln(&choice)

			fmt.Println("Vuelos:")
			for i, v := range originalReservation.Vuelos {
				fmt.Printf("%d. %s %s - %s\n", i+1, v.Numero_vuelo, v.Hora_salida, v.Hora_llegada)
			}

			var vueloToReplace int
			fmt.Print("Ingrese una opción: ")
			fmt.Scanln(&vueloToReplace)
			//replaceObject := originalReservation.Vuelos[vueloToReplace-1]

			// Mostrar ancillaries disponibles

			ancillaries := GetAncillaries()

			showAncillaries(ancillaries)
			_, precio := chooseAncillaries(ancillaries)
			fmt.Println("Total ancillaries: ", precio)

			//var selectedAncillaries []model.AncillaryPasajero
			//var balances model.Balance

		}
	}
}

func cambiar_vuelos(reserva model.Reserva, vueloToReplace int, pnr string, apellido string) (model.Reserva, error) {
	vuelo := reserva.Vuelos[vueloToReplace]

	// Elegir nueva fecha
	var newDate string
	fmt.Print("Ingrese nueva fecha: ")
	fmt.Scanln(&newDate)

	// Obtener vuelos disponibles para la nueva fecha
	vuelos := findVuelos(vuelo.Origen, vuelo.Destino, newDate)
	if len(vuelos) == 0 {
		return model.Reserva{}, fmt.Errorf("No hay vuelos disponibles para la nueva fecha")
	}

	// Mostrar vuelos disponibles para la nueva fecha
	fmt.Println("Vuelos disponibles:")
	for i, v := range vuelos {
		fmt.Printf("%d. %s %s - %s\n", i+1, v.Numero_vuelo, v.Hora_salida, v.Hora_llegada)
	}
	var selectedNewVuelo string
	fmt.Print("Ingrese una opción: ")
	fmt.Scanln(&selectedNewVuelo)

	// Obtener vuelo seleccionado
	var newReserva model.Reserva
	newReserva.PNR = reserva.PNR
	newReserva.Pasajeros = append(newReserva.Pasajeros, reserva.Pasajeros...)
	newReserva.Vuelos = append(newReserva.Vuelos, reserva.Vuelos...)

	newReserva.Vuelos[vueloToReplace] = vuelos[vueloToReplace]

	return newReserva, nil
}

func findVuelos(origen string, destino string, newDate string) []model.Vuelo {
	queries := map[string]string{
		"origen":  origen,
		"destino": destino,
		"fecha":   newDate,
	}
	url := crear_url("vuelo", queries)

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var vuelos []model.Vuelo
	if err := json.Unmarshal(body, &vuelos); err != nil {
		return nil
	}

	return vuelos
}

func actualizar_reservacion(pnr string, apellido string, updatedReservation model.Reserva) {
	queries := map[string]string{
		"pnr":      pnr,
		"apellido": apellido,
	}
	url := crear_url("reserva", queries)

	updatedReservationJson, err := json.Marshal(&updatedReservation)
	if err != nil {
		log.Fatal("Error al actualizar la reserva")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(updatedReservationJson))
	if err != nil {
		log.Fatal("Error al actualizar la reserva")
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)

	defer req.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("La reserva fué modificada exitosamente!")
	} else {
		log.Fatal("Error del servidor al actualizar la reserva")
	}
}

func actualizar_stock(reserva model.Reserva, newVuelo model.Vuelo, oldVuelo model.Vuelo) {
	updateStock(reserva, newVuelo, "decrease")
	updateStock(reserva, oldVuelo, "add")
}

func updateStock(reserva model.Reserva, vuelo model.Vuelo, operator string) {
	queries := map[string]string{
		"numero_vuelo": vuelo.Numero_vuelo,
		"origen":       vuelo.Origen,
		"destino":      vuelo.Destino,
		"fecha":        vuelo.Fecha,
	}
	url := crear_url("vuelo", queries)

	var stock Stock

	if operator == "add" {
		stock = Stock{
			Stock: vuelo.Avion.Stock_de_pasajeros + len(reserva.Pasajeros),
		}
	} else {
		stock = Stock{
			Stock: vuelo.Avion.Stock_de_pasajeros - len(reserva.Pasajeros),
		}
	}

	stockJson, err := json.Marshal(&stock)
	if err != nil {
		log.Fatal("Error al actualizar el stock")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(stockJson))
	if err != nil {
		log.Fatal("Error al actualizar el stock")
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)

	defer req.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Error del servidor al actualizar el stock")
	}
}
