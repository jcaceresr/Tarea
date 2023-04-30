package routes

import (
	getcollection "api/Collection"
	database "api/databases"
	model "api/model"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollectionVuelos(DB, "Vuelos")

	postId := c.Param("postId")
	var post model.Posts

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	edited := bson.M{"Nombre": post.Nombre, "Apellido": post.Apellido, "Edad": post.Edad}

	result, err := postCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": edited})

	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if result.MatchedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "Data": res})
}

// funcion update para vuelos

func UpdateVuelo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var vueloCollection = getcollection.GetCollectionVuelos(DB, "Vuelos")

	defer cancel()

	numeroVuelo := c.Query("numero_vuelo")
	origen := c.Query("origen")
	destino := c.Query("destino")
	fecha := c.Query("fecha")

	// Verificar si existe el vuelo
	var vuelo model.Vuelo
	if err := vueloCollection.FindOne(ctx, bson.M{"numero_vuelo": numeroVuelo, "origen": origen, "destino": destino, "fecha": fecha}).Decode(&vuelo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vuelo no encontrado"})
		return
	}

	// Decodificar el cuerpo de la solicitud
	var updateData struct {
		Stock_de_pasajeros int `json:"stock_de_pasajeros"`
	}
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de solicitud incorrectos"})
		return
	}

	// Actualizar el stock de pasajeros en el vuelo
	update := bson.M{"$set": bson.M{"avion.stock_de_pasajeros": updateData.Stock_de_pasajeros}}
	result, err := vueloCollection.UpdateOne(ctx, bson.M{"numero_vuelo": numeroVuelo, "origen": origen, "destino": destino, "fecha": fecha}, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error de servidor"})
		return
	}

	if result.MatchedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se encontró el vuelo"})
		return
	}

	// Devolver el vuelo actualizado
	updatedVuelo := model.Vuelo{
		Numero_vuelo: vuelo.Numero_vuelo,
		Origen:       vuelo.Origen,
		Destino:      vuelo.Destino,
		Hora_salida:  vuelo.Hora_salida,
		Hora_llegada: vuelo.Hora_llegada,
	}

	c.JSON(http.StatusOK, updatedVuelo)
}

// funcion update para reservas
func UpdateReserva(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var reservaCollection = getcollection.GetCollectionReservas(DB, "Reservas")

	pnr := c.Query("pnr")

	var reserva model.Reserva

	defer cancel()

	if err := c.BindJSON(&reserva); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	edited := bson.M{"vuelos": reserva.Vuelos, "pasajeros": reserva.Pasajeros}

	result, err := reservaCollection.UpdateOne(ctx, bson.M{"pnr": pnr}, bson.M{"$set": edited})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if result.MatchedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"PNR": pnr})
}
