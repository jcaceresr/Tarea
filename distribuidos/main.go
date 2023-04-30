package main

import (

	"fmt"
	funciones "api/funciones"

)

func main() {

	fmt.Println("Bienvenido")

	Inicio:

	for{
        var option int

        fmt.Println("\nMenu: ")
        fmt.Println("1. Gestionar Reservas")
        fmt.Println("2. Obtener Estadisticas")
        fmt.Println("3. Salir")
        fmt.Print("Ingrese una opción: ")
        fmt.Scanln(&option)
        
        switch option{
            // Gestion de reservas
            case 1:

                    for{
                        var option int
                        
                        fmt.Println("\nSubmenu: ")
                        fmt.Println("1. Crear Reserva")
						fmt.Println("2. Obtener Reserva")
                        fmt.Println("3. Modificar Reserva")
						fmt.Println("4. Salir")
                        fmt.Println("Ingrese una opción: ")
                        fmt.Scanln(&option)


						//crear Reserva             
                        if option == 1 {

							var fechaIda string
							var fechaVuelta string
							var origen string
							var destino string
							var cantidadPasajeros int

							fmt.Println("Ingrese la fecha de ida: ")
							fmt.Scanln(&fechaIda)
							fmt.Print("Ingrese la fecha de vuelta: ")
							fmt.Scanln(&fechaVuelta)
							fmt.Print("Ingrese el origen: ")
							fmt.Scanln(&origen)
							fmt.Print("Ingrese el destino: ")
							fmt.Scanln(&destino)
							fmt.Print("Ingrese la cantidad de pasajeros: ")
							fmt.Scanln(&cantidadPasajeros)

							funciones.Crear_reserva(fechaIda, fechaVuelta, origen, destino, cantidadPasajeros)
						}

						// obtener Reserva 
						if option == 2 { 
							
							var pnr string
							var apellido string

							fmt.Println("\nIngrese PNR:")
							fmt.Scan(&pnr)
							fmt.Println("ingrese apellido:")
							fmt.Scan(&apellido)
							funciones.Obtener_reserva(pnr, apellido)
						}
							// Modificar Reserva

						
						if option == 3 {
							var pnr string
							var apellido string

							fmt.Print("Ingrese el PNR: ")
							fmt.Scanln(&pnr)

							fmt.Print("Ingrese el apellido: ")
							fmt.Scanln(&apellido)

							funciones.Updates(pnr, apellido)
						}


						if option == 4 {
								break 
						}
         
                        
                    }

                

            // obtener estadistica
            case 2:
                //aqui van las funciones de obtener estadisticas
				fmt.Println("\nfuncion de estadistica xd")
            case 3:
                fmt.Println("\nHasta Luego!")
                break Inicio
            
        
        }
	}
}
