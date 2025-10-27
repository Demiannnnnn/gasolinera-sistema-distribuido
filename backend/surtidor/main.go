package main

import (
	"fmt"
	"net/http" // Necesario para http.ListenAndServe
	"os"        // Necesario para leer variables de entorno
)

func main() {
	fmt.Println("--- Arrancando Backend Surtidor ---")

	// Intenta leer las variables de entorno para asegurar que están configuradas
	distribuidorHost := os.Getenv("DISTRIBUIDOR_HOST")
	fmt.Printf("DISTRIBUIDOR_HOST: %s\n", distribuidorHost)


	// Crear un servidor HTTP simple que no hace nada, solo escucha en el puerto correcto.
	// Esto mantendrá el contenedor corriendo.
	// Asegúrate que el puerto coincida con el EXPOSE del Dockerfile (8082)
	fmt.Println("Servidor escuchando en el puerto 8082...")
	err := http.ListenAndServe(":8082", nil) // nil significa usar el manejador por defecto (que no hace nada)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
