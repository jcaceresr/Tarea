package funciones

import (
	"api/model"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

type PNR struct {
	Pnr string `json:"pnr,omitempty" bson:"pnr,omitempty"`
}

func crear_url(ruta string, query map[string]string) string {
	url, err := url.Parse("http://localhost:5000/api/" + ruta)
	if err != nil {
		log.Fatal("URL no válida")
	}

	if query == nil {
		return url.String()
	}

	values := url.Query()

	for key, value := range query {
		values.Add(key, value)
	}

	url.RawQuery = values.Encode()

	return url.String()
}

func valor_pasaje(vuelo model.Vuelo) int {

	date := formato_fecha(vuelo.Fecha)

	salida, _ := time.Parse("2022-04-28 10:10", fmt.Sprintf("%s %s", date, vuelo.Hora_salida))
	llegada, _ := time.Parse("2022-04-28 10:10", fmt.Sprintf("%s %s", date, vuelo.Hora_llegada))

	dif := llegada.Sub(salida)
	minutos := int(dif.Minutes())

	precioPasaje := minutos * 590

	return precioPasaje
}

func formato_fecha(date string) string {

	//transforma fecha de / a -
	separated := strings.Split(date, "/")
	formato := strings.Join(reversa(separated), "-")

	return formato
}

func reversa(fecha []string) []string {

	//da vuelta al fecha
	fecha_final := make([]string, len(fecha))

	for i, j := range fecha {
		fecha_final[len(fecha)-1-i] = j
	}

	return fecha_final
}
