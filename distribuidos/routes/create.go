package routes

import (
	getcollection "api/Collection"
	database "api/databases"
	model "api/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
)

// funcion que al ser llamada genera un pnr de 6 caracteres aleatorios entre letras de la a-z y numeros del 0-9
func generatePNR() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreatePasajero(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, "Posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.Pasajero)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	//)
	postPayload := model.Pasajero{
		//agrega un pnr de 6 caracteres aleatorios entre letras de la a-z y numeros del 0-9,
		Nombre:      post.Nombre,
		Apellido:    post.Apellido,
		Edad:        post.Edad,
		Ancillaries: post.Ancillaries,
		Balances:    post.Balances,
	}

	result, err := postCollection.InsertOne(ctx, postPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": map[string]interface{}{"data": result}})

}

func CreateVuelo(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollectionVuelos(DB, "Vuelos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.Vuelo)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	postPayload := model.Vuelo{
		Numero_vuelo: post.Numero_vuelo,
		Origen:       post.Origen,
		Destino:      post.Destino,
		Hora_salida:  post.Hora_salida,
		Hora_llegada: post.Hora_llegada,
		Fecha:        post.Fecha,
		Avion:        post.Avion,
		Ancillaries:  post.Ancillaries,
	}

	result, err := postCollection.InsertOne(ctx, postPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": map[string]interface{}{"data": result}})

}

func CreateReserva(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollectionReservas(DB, "Reservas")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.Reserva)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	fmt.Printf("%+v", post.Pasajeros[0].Ancillaries)

	pnr := generatePNR()
	postPayload := model.Reserva{
		PNR:       pnr,
		Vuelos:    post.Vuelos,
		Pasajeros: post.Pasajeros,
	}

	result, err := postCollection.InsertOne(ctx, postPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": map[string]interface{}{"data": result}})

}
