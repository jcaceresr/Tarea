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

func Crear_reserva(fechaIda string, fechaVuelta string, origen string, destino string, cantidadPasajeros int) {
	// Seleccionar vuelos
	fmt.Println("Vuelos disponibles:")
	vuelos, precioPasajeIda, precioPasajeVuelta := elegir_vuelo(fechaIda, fechaVuelta, origen, destino)

	// Ingresar pasajeros
	pasajeros := info_pasajeros(cantidadPasajeros, len(vuelos) > 1)

	// Actualizar balances con precio de pasajes y obtener costo total
	var costoTotal int
	for i := range pasajeros {
		pasajeros[i].Balances.Vuelo_ida = precioPasajeIda
		pasajeros[i].Balances.Vuelo_vuelta = precioPasajeVuelta

		costoTotal += precioPasajeIda +
			precioPasajeVuelta +
			pasajeros[i].Balances.Ancillaries_ida +
			pasajeros[i].Balances.Ancillaries_vuelta
	}

	// Crear reserva
	reserva := model.Reserva{
		Pasajeros: pasajeros,
		Vuelos:    vuelos,
	}

	reservaJson, err := json.Marshal(reserva)
	if err != nil {
		log.Fatal("Error al crear reserva")
	}

	url := crear_url("reserva", nil)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reservaJson))
	if err != nil {
		log.Fatal("Error al crear reserva")
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error al crear reserva")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error del servidor al crear la reserva")
	}

	if resp.StatusCode != http.StatusCreated {
		log.Fatal("Error del servidor al crear la reserva")
		return
	}

	var reservaResponse PNR

	if err := json.Unmarshal(body, &reservaResponse); err != nil {
		log.Fatal("Respuesta no válida")
	}

	fmt.Println("La reserva fue generada, el PNR es:", reservaResponse.Pnr)
	fmt.Printf("El costo total de la reserva fué de $%d\n", costoTotal)
}

//pasajeros

func info_pasajeros(cantidadPasajeros int, vueloConVuelta bool) []model.Pasajero {
	pasajeros := make([]model.Pasajero, cantidadPasajeros)

	for i := 0; i < cantidadPasajeros; i++ {
		fmt.Printf("Pasajero %d:\n", i+1)

		pasajero := ancillaries_pasajeros(vueloConVuelta)
		pasajeros[i] = pasajero
	}

	return pasajeros
}

func ancillaries_pasajeros(vueloConVuelta bool) model.Pasajero {
	nombre, apellido, edad := getPassengerData()
	ancillariesData := GetAncillaries()

	fmt.Println("Ancillaries de ida:")
	showAncillaries(ancillariesData)

	ancillariesIda, totalAncillariesIda := chooseAncillaries(ancillariesData)
	fmt.Println("Total ancillaries: ", totalAncillariesIda)

	var selectedAncillaries []model.AncillaryPasajero
	var balances model.Balance

	if vueloConVuelta {
		fmt.Println("Ancillaries de vuelta:")
		showAncillaries(ancillariesData)

		ancillariesVuelta, totalAncillariesVuelta := chooseAncillaries(ancillariesData)
		fmt.Println("Total ancillaries: ", totalAncillariesVuelta)

		selectedAncillaries = []model.AncillaryPasajero{
			{
				Ida: ancillariesIda,
			},
			{
				Vuelta: ancillariesVuelta,
			},
		}

		balances = model.Balance{
			Ancillaries_ida:    totalAncillariesIda,
			Ancillaries_vuelta: totalAncillariesVuelta,
		}

	} else {
		selectedAncillaries = []model.AncillaryPasajero{
			{
				Ida: ancillariesIda,
			},
		}

		balances = model.Balance{
			Ancillaries_ida: totalAncillariesIda,
		}
	}

	return model.Pasajero{
		Nombre:      nombre,
		Apellido:    apellido,
		Edad:        edad,
		Ancillaries: selectedAncillaries,
		Balances:    balances,
	}
}

func getPassengerData() (string, string, int) {
	var nombre string
	var apellido string
	var edad int

	fmt.Print("Ingrese el nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Ingrese el apellido: ")
	fmt.Scanln(&apellido)
	fmt.Print("Ingrese la edad: ")
	fmt.Scanln(&edad)

	return nombre, apellido, edad
}

func Obtener_reserva(pnr string, apellido string) {
	queries := map[string]string{
		"pnr":      pnr,
		"apellido": apellido,
	}

	url := crear_url("reserva", queries)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Reserva no encontrada")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Reserva no encontrada")
	}

	var reserva model.Reserva

	if err := json.Unmarshal(body, &reserva); err != nil {
		log.Fatal("Respuesta no válida")
	}

	// Mostrar detalles de vuelos de ida y vuelta
	for i, vuelo := range reserva.Vuelos {
		var seccion string
		if i == 0 {
			seccion = "Ida"
		} else {
			seccion = "Vuelta"
		}

		fmt.Print(
			fmt.Sprintf("%s:\n %s %s %s\n", seccion, vuelo.Numero_vuelo, vuelo.Hora_salida, vuelo.Hora_llegada),
		)
	}

	// Mostrar detalle de pasajeros
	fmt.Print(
		"Pasajeros:\n",
	)

	for _, pasajero := range reserva.Pasajeros {
		// Datos personales
		fmt.Print(
			fmt.Sprintf("%s %d\n", pasajero.Nombre, pasajero.Edad),
		)

		// Ancillaries
		for _, ancillaries := range pasajero.Ancillaries {
			if ancillaries.Ida != nil {
				fmt.Print("Ancillaries ida: ")
				for _, ancillary := range ancillaries.Ida {
					fmt.Print(
						fmt.Sprintf("%s ", ancillary.Ssr),
					)
				}
				fmt.Println()
			}

			if ancillaries.Vuelta != nil {
				fmt.Print("Ancillaries vuelta: ")
				for _, ancillary := range ancillaries.Vuelta {
					fmt.Print(
						fmt.Sprintf("%s ", ancillary.Ssr),
					)
				}
				fmt.Println()
			}
		}
	}
}

func reserva_sin_print(pnr string, apellido string) model.Reserva {
	queries := map[string]string{
		"pnr":      pnr,
		"apellido": apellido,
	}

	url := crear_url("reserva", queries)

	resp, err := http.Get(url)
	if err != nil {
		return model.Reserva{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Reserva{}
	}

	var reserva model.Reserva
	if err := json.Unmarshal(body, &reserva); err != nil {
		return model.Reserva{}
	}

	return reserva
}
