package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Estructura para recibir los precios desde el frontend

type PriceUpdateRequest struct {
	Precios map[string]float64 `json:"precios"`
}

// Función para obtener la URL del distribuidor basado en su ID

func getDistributorURL(distributorID string) string {
	switch distributorID {
	case "DISTRIBUIDOR_1":
		return os.Getenv("DISTRIBUIDOR_1_HOST")
	case "DISTRIBUIDOR_2":
		return os.Getenv("DISTRIBUIDOR_2_HOST")
	case "DISTRIBUIDOR_3":
		return os.Getenv("DISTRIBUIDOR_3_HOST")
	default:
		return ""
	}
}

// Manejador para actualizar los precios de un distribuidor

func updatePriceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	distributorID := vars["id"]
	fmt.Printf("Solicitud de actualización de precios para el distribuidor: %s\n", distributorID)

	// Obtener la URL del distribuidor
	distributorURL := getDistributorURL(distributorID)
	if distributorURL == "" {
		fmt.Printf("Distribuidor no encontrado: %s\n", distributorID)
		http.Error(w, "Distribuidor no encontrado", http.StatusNotFound)
		return
	}
	fmt.Printf("URL del distribuidor: %s\n", distributorURL)

	// Decodificar el cuerpo de la solicitud
	var priceUpdate PriceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&priceUpdate); err != nil {
		fmt.Printf("Error al decodificar el cuerpo de la solicitud: %v\n", err)
		http.Error(w, "Cuerpo de la solicitud inválido", http.StatusBadRequest)
		return
	}
	fmt.Printf("Precios recibidos: %+v\n", priceUpdate.Precios)

	// Crear el cuerpo de la solicitud para el distribuidor
	jsonBody, err := json.Marshal(priceUpdate.Precios)
	if err != nil {
		fmt.Printf("Error al procesar la solicitud: %v\n", err)
		http.Error(w, "Error al procesar la solicitud", http.StatusInternalServerError)
		return
	}

	// Crear la solicitud PUT al distribuidor
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/prices", distributorURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error al crear la solicitud para el distribuidor: %v\n", err)
		http.Error(w, "Error al crear la solicitud para el distribuidor", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error al comunicarse con el distribuidor: %v\n", err)
		http.Error(w, "Error al comunicarse con el distribuidor", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Verificar la respuesta del distribuidor
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("El distribuidor devolvió un estado no OK: %d\n", resp.StatusCode)
		http.Error(w, "El distribuidor no pudo actualizar los precios", http.StatusInternalServerError)
		return
	}

	// Responder al frontend
	fmt.Println("Precios actualizados correctamente en el distribuidor")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Precios actualizados correctamente"})
}

func main() {
	fmt.Println("--- Arrancando Backend Empresa ---")

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r.HandleFunc("/distribuidores/{id}/prices", updatePriceHandler).Methods("PUT")

	handler := c.Handler(r)

	fmt.Println("Servidor escuchando en el puerto 8081...")
	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}