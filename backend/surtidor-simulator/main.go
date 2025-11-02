package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Usamos un estado global para la pausa, protegido por un mutex para acceso concurrente seguro.
var isPaused = false
var pauseMutex = &sync.Mutex{}

// leerComandos escucha en la conexión por comandos del distribuidor.
func leerComandos(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		comando, err := reader.ReadString('\n')
		if err != nil {
			// Si hay un error (ej. conexión cerrada), terminamos la goroutine.
			return
		}

		comando = strings.TrimSpace(comando)
		pauseMutex.Lock()
		if comando == "CMD:PAUSE" {
			isPaused = true
			fmt.Println("[SIMULADOR] Estado cambiado a PAUSADO")
		} else if comando == "CMD:RESUME" {
			isPaused = false
			fmt.Println("[SIMULADOR] Estado cambiado a REANUDADO")
		}
		pauseMutex.Unlock()
	}
}

// simularUnSurtidor representa el ciclo de vida de un único surtidor.
func simularUnSurtidor(id int, distributorAddress string) {
	// Bucle principal de reconexión. Si algo falla, vuelve aquí.
	for {
		var conn net.Conn
		var err error

		// 1. Bucle de conexión inicial
		for {
			conn, err = net.Dial("tcp", distributorAddress)
			if err == nil {
				break // Conexión exitosa
			}
			fmt.Printf("[Surtidor %d] Error al conectar a %s: %v. Reintentando en 5 segundos...\n", id, distributorAddress, err)
			time.Sleep(5 * time.Second)
		}

		fmt.Printf("[Surtidor %d] Conectado a %s\n", id, distributorAddress)

		// Inicia una goroutine para escuchar comandos del distribuidor.
		// Se crea una nueva por cada nueva conexión.
		go leerComandos(conn)

		// 2. Bucle de operación del surtidor
		operationLoop(id, conn)

		// Si operationLoop termina (por desconexión), cerramos la conexión y el bucle for principal reintentará todo.
		conn.Close()
		fmt.Printf("[Surtidor %d] Conexión perdida. Intentando reconectar...\n", id)
		time.Sleep(2 * time.Second) // Esperar un poco antes de reintentar
	}
}

// operationLoop maneja el envío de datos mientras la conexión está activa.
func operationLoop(id int, conn net.Conn) {
	for {
		pauseMutex.Lock()
		paused := isPaused
		pauseMutex.Unlock()

		if paused {
			time.Sleep(1 * time.Second)
			continue // Salta el resto del ciclo si está en pausa.
		}

		sleepTime := time.Duration(rand.Intn(15)+5) * time.Second
		time.Sleep(sleepTime)

		litros := 5 + rand.Float64()*(50-5)
		tipoCombustible := []string{"93", "95", "97", "Diesel"}[rand.Intn(4)]

		surtidorIDStr := fmt.Sprintf("surtidor-%d", id)
		mensaje := fmt.Sprintf("ID:%s;LITROS:%.2f;TIPO:%s\n", surtidorIDStr, litros, tipoCombustible)
		_, err := conn.Write([]byte(mensaje))
		if err != nil {
			// Si hay un error de escritura, salimos de este bucle para que la función superior pueda reconectar.
			return
		}
		fmt.Printf("[Surtidor %d] Transacción enviada: %.2fL de %s\n", id, litros, tipoCombustible)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	distributorAddress := os.Getenv("DISTRIBUTOR_ADDRESS")
	if distributorAddress == "" {
		fmt.Println("Error: La variable de entorno DISTRIBUTOR_ADDRESS es requerida.")
		os.Exit(1)
	}

	numSurtidoresStr := os.Getenv("NUM_SURTIDORES")
	numSurtidores, err := strconv.Atoi(numSurtidoresStr)
	if err != nil || numSurtidores <= 0 {
		numSurtidores = 5
	}

	fmt.Printf("Iniciando simulador para %d surtidores hacia %s\n", numSurtidores, distributorAddress)

	for i := 0; i < numSurtidores; i++ {
		go simularUnSurtidor(i+1, distributorAddress)
	}

	select {}
}
