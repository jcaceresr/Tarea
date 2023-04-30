package routes

import (
	getcollection "api/Collection"
	database "api/databases"
	model "api/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadOnePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, "Posts")

	postId := c.Param("postId")
	var result model.Posts

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	err := postCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&result)

	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success!", "Data": res})
}

func GetReserva(c *gin.Context) {
	var DB = database.ConnectDB()
	var collection = getcollection.GetCollectionReservas(DB, "Reservas")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pnr := c.Query("pnr")
	apellido := c.Query("apellido")

	var reserva model.Reserva
	err := collection.FindOne(ctx, bson.M{"pnr": pnr, "pasajeros.apellido": apellido}).Decode(&reserva)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Reserva no encontrada"})
		return
	}

	c.JSON(http.StatusOK, reserva)
}

func GetVuelos(c *gin.Context) {
	var DB = database.ConnectDB()
	var vuelosCollection = getcollection.GetCollectionVuelos(DB, "Vuelos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	origen := c.Query("origen")
	destino := c.Query("destino")
	fecha := c.Query("fecha")

	filter := bson.M{"origen": origen, "destino": destino, "fecha": fecha}

	fmt.Println("antes de find")

	cursor, err := vuelosCollection.Find(ctx, filter)

	fmt.Println("despues de find")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	var vuelos []model.Vuelo

	if err := cursor.All(ctx, &vuelos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	fmt.Println("antes de retornar")

	c.JSON(http.StatusOK, vuelos)
}

func GetStatics(c *gin.Context) {
	// crea un map vacio que almacenara "origen-destino" : ganancia de las reservas
	var ganancias = make(map[string]int)
	var DB = database.ConnectDB()
	var collection = getcollection.GetCollectionReservas(DB, "Reservas")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// obtiene todas las reservas
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	//rellena las keys con el formato "origen-destino" y los valores con 0
	for cursor.Next(ctx) {
		var reserva model.Reserva
		cursor.Decode(&reserva)
		//La propiedad vuelo es un array de vuelos, por lo que se itera sobre el array
		for _, vuelo := range reserva.Vuelos {
			ganancias[vuelo.Origen+"-"+vuelo.Destino] = 0
		}

	}
	/*para rellenar la ganancia de cada ruta en el map ganancia se debe acceder a cada uno de los pasajeros de cada reserva,
	a su campo balances y sumar las propiedades vuelo_ida y vuelo_vuelta, luego agregarlo a la respectiva ruta*/
	// itera nuevamente sobre todas las reservas
	cursor, err = collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	for cursor.Next(ctx) {
		var reserva model.Reserva
		cursor.Decode(&reserva)
		// itera sobre los vuelos de cada reserva
		for _, vuelo := range reserva.Vuelos {
			// itera sobre los pasajeros de cada reserva
			for _, pasajero := range reserva.Pasajeros {
				// suma la ganancia de los vuelos de ida y vuelta para el pasajero actual
				ganancia := pasajero.Balances.Vuelo_ida + pasajero.Balances.Vuelo_vuelta
				// agrega la ganancia a la ruta correspondiente en el map ganancias
				ruta := vuelo.Origen + "-" + vuelo.Destino
				ganancias[ruta] += ganancia
			}
		}
	}
	// recorre el map ganancias e imprime cada ruta con su ganancia separados por un salto de linea
	for ruta, ganancia := range ganancias {
		fmt.Println(ruta, ganancia)
	}
	//imprime el valor maximo y el valor minimo del map ganancias
	var max int
	var min int
	for _, ganancia := range ganancias {
		if ganancia > max {
			max = ganancia
		}
		if ganancia < min {
			min = ganancia
		}
	}
	fmt.Println("maximo: ", max)
	fmt.Println("minimo: ", min)

	// retorna el map ganancias en formato JSON
	c.JSON(http.StatusOK, gin.H{"data": ganancias})

}
