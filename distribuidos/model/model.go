package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Posts struct {
	ID          primitive.ObjectID
	Nombre      string            `json:"nombre" bson:"nombre"`
	Apellido    string            `json:"apellido" bson:"apellido"`
	Edad        int               `json:"edad" bson:"edad"`
	Ancillaries AncillaryPasajero `json:"ancillaries" bson:"ancillaries"`
	Balances    Balance           `json:"balances" bson:"balances"`
}

type Reserva struct {
	PNR       string     `json:"pnr"`
	Vuelos    []Vuelo    `json:"vuelos"`
	Pasajeros []Pasajero `json:"pasajeros"`
}

type Pasajero struct {
	//un campo pnr que se genera aleatoriamente al crear un pasajero
	Nombre      string              `json:"nombre" bson:"nombre"`
	Apellido    string              `json:"apellido" bson:"apellido"`
	Edad        int                 `json:"edad" bson:"edad"`
	Ancillaries []AncillaryPasajero `json:"ancillaries,omitempty" bson:"ancillaries,omitempty"`
	Balances    Balance             `json:"balances" bson:"balances"`
}

type AncillaryPasajero struct {
	Ida    []AncillaryInfo `json:"ida,omitempty" bson:"ida,omitempty"`
	Vuelta []AncillaryInfo `json:"vuelta,omitempty" bson:"vuelta,omitempty"`
}

type AncillaryInfo struct {
	Ssr      string `json:"ssr"`
	Cantidad int    `json:"cantidad"`
}

type Balance struct {
	Ancillaries_ida    int `json:"ancillaries_ida"`
	Vuelo_ida          int `json:"vuelo_ida"`
	Ancillaries_vuelta int `json:"ancillaries_vuelta"`
	Vuelo_vuelta       int `json:"vuelo_vuelta"`
}

type Avion struct {
	Modelo             string `json:"modelo"`
	Numero_de_serie    string `json:"numero_de_serie"`
	Stock_de_pasajeros int    `json:"stock_de_pasajeros"`
}

type AncillaryVuelos struct {
	Nombre string `json:"nombre"`
	Stock  int    `json:"stock"`
	Ssr    string `json:"ssr"`
}

type Vuelo struct {
	Numero_vuelo string            `json:"numero_vuelo"`
	Origen       string            `json:"origen"`
	Destino      string            `json:"destino"`
	Hora_salida  string            `json:"hora_salida"`
	Hora_llegada string            `json:"hora_llegada"`
	Avion        *Avion            `json:",omitempty"`
	Fecha        string            `json:"fecha,omitempty"`
	Ancillaries  []AncillaryVuelos `json:",omitempty"`
}
