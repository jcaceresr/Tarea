package main

import (
	routes "api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/", routes.CreatePasajero)

	router.POST("/api/vuelo", routes.CreateVuelo) //crea un vuelo prueba 1

	router.POST("/api/reserva", routes.CreateReserva) //crea una reserva prueba 1")

	router.GET("/api/reserva", routes.GetReserva)

	router.GET("/api/vuelo", routes.GetVuelos)

	router.PUT("/api/vuelo", routes.UpdateVuelo)

	router.PUT("/api/reserva", routes.UpdateReserva)

	router.GET("/api/estadisticas", routes.GetStatics)

	// called as localhost:3000/getOne/{id}
	router.GET("getOne/:postId", routes.ReadOnePost)

	// called as localhost:3000/update/{id}
	router.PUT("/update/:postId", routes.UpdatePost)

	// called as localhost:3000/delete/{id}
	router.DELETE("/delete/:postId", routes.DeletePost)

	router.Run("localhost: 5000")
}
