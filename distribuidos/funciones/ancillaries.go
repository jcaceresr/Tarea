package funciones

import(
	"api/model"
	"fmt"
	"strings"
	"sort"
)


type Ancillaries_precio struct {
	Nombre string `json:"nombre" bson:"nombre"`
	Precio int    `json:"precio" bson:"precio"`
	Ssr    string `json:"ssr" bson:"ssr"`
}


func GetAncillaries() map[string]Ancillaries_precio {

	ancillariesData := map[string]Ancillaries_precio{
		"1": {Nombre: "Equipaje de mano", Precio: 10000, Ssr: "BGH"},
		"2": {Nombre: "Equipaje de bodega", Precio: 30000, Ssr: "BGR"},
		"3": {Nombre: "Asiento", Precio: 5000, Ssr: "STDF"},
		"4": {Nombre: "Embarque y Check In prioritario", Precio: 2000, Ssr: "PAXS"},
		"5": {Nombre: "Mascota en cabina", Precio: 40000, Ssr: "PTCR"},
		"6": {Nombre: "Mascota en bodega", Precio: 40000, Ssr: "AVIH"},
		"7": {Nombre: "Equipaje especial", Precio: 35000, Ssr: "SPML"},
		"8": {Nombre: "Acceso a Salón VIP", Precio: 15000, Ssr: "LNGE"},
		"9": {Nombre: "Wi-Fi a bordo", Precio: 20000, Ssr: "WIFI"},
	}

	return ancillariesData
}



func showAncillaries(ancillariesData map[string]Ancillaries_precio) {
	
	keys := make([]string, 0, len(ancillariesData))
	for k := range ancillariesData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := ancillariesData[k]
		fmt.Printf("%s: %s - %d\n", k, v.Nombre, v.Precio)
	}
}



func chooseAncillaries(ancillariesData map[string]Ancillaries_precio) ([]model.AncillaryInfo, int) {

	var selection string
	fmt.Print("Ingrese los ancillaries (separados por coma): ")
	fmt.Scanln(&selection)

	selectedAncillariesSplitted := strings.Split(selection, ",")

	selectedAncillaries := []model.AncillaryInfo{}

	var selectedAncillariesTotalPrice int

	for _, ancillary := range selectedAncillariesSplitted {
		ancillaryObject := ancillariesData[ancillary]

		selectedAncillariesTotalPrice += ancillaryObject.Precio

		ssr := ancillaryObject.Ssr

		found := false
		for i, selectedAncillary := range selectedAncillaries {
			if selectedAncillary.Ssr == ssr {
				selectedAncillaries[i].Cantidad++
				found = true
				break
			}
		}

		if !found {
			selectedAncillaries = append(selectedAncillaries, model.AncillaryInfo{
				Ssr:      ancillaryObject.Ssr,
				Cantidad: 1,
			})
		}
	}

	return selectedAncillaries, selectedAncillariesTotalPrice
}