package main

import (
	"fmt"
	"net/http" // Necesario para http.ListenAndServe
	"os"        // Necesario para leer variables de entorno
)

func main() {
	fmt.Println("--- Arrancando Backend Distribuidor ---")

	// Intenta leer las variables de entorno para asegurar que están configuradas
	dbURL := os.Getenv("DATABASE_URL")
	empresaHost := os.Getenv("EMPRESA_HOST")
	surtidorHost := os.Getenv("SURTIDOR_HOST")
	fmt.Printf("DATABASE_URL: %s\n", dbURL)
	fmt.Printf("EMPRESA_HOST: %s\n", empresaHost)
	fmt.Printf("SURTIDOR_HOST: %s\n", surtidorHost)

	// Crear un servidor HTTP simple que no hace nada, solo escucha en el puerto correcto.
	// Esto mantendrá el contenedor corriendo.
	// Asegúrate que el puerto coincida con el EXPOSE del Dockerfile (8080)
	fmt.Println("Servidor escuchando en el puerto 8080...")
	err := http.ListenAndServe(":8080", nil) // nil significa usar el manejador por defecto (que no hace nada)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
