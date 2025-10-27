package main

import (
	"fmt"
	"net/http" // Necesario para http.ListenAndServe
	"os"        // Necesario para leer variables de entorno
)

func main() {
	fmt.Println("--- Arrancando Backend Empresa ---")

	// Intenta leer las variables de entorno para asegurar que están configuradas
	dbURL := os.Getenv("DATABASE_URL")
	distribuidorHost := os.Getenv("DISTRIBUIDOR_HOST")
	fmt.Printf("DATABASE_URL: %s\n", dbURL)
	fmt.Printf("DISTRIBUIDOR_HOST: %s\n", distribuidorHost)

	// Crear un servidor HTTP simple que no hace nada, solo escucha en el puerto correcto.
	// Esto mantendrá el contenedor corriendo.
	// Asegúrate que el puerto coincida con el EXPOSE del Dockerfile (8081)
	fmt.Println("Servidor escuchando en el puerto 8081...")
	err := http.ListenAndServe(":8081", nil) // nil significa usar el manejador por defecto (que no hace nada)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
