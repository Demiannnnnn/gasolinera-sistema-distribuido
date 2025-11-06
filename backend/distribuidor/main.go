package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// --- Variables Globales para Sockets ---
var surtidorConnections struct {
	sync.Mutex
	conns map[net.Conn]bool
}

// --- Tipos y Constantes ---
const (
	Fuel93     = "93"
	Fuel95     = "95"
	Fuel97     = "97"
	FuelDiesel = "DIESEL"
)

type ventaCombustible struct {
	DistribuidorID  string    `json:"distribuidor_id"`
	ID_venta        int       `json:"id"`
	HiloID          string    `json:"surtidor_id"` // Renombrado de HiloID
	TipoCombustible string    `json:"tipo_combustible"`
	LitrosVendidos  float64   `json:"litros_vendidos"`
	PrecioPorLitro  float64   `json:"precio_por_litro"`
	TotalVenta      float64   `json:"total_venta"`
	FechaVenta      time.Time `json:"fecha_venta"`
	NumeroFactura   string    `json:"numero_factura"`
	Persistida      bool      `json:"-"`
}

type distribuidor struct {
	Distribuidor_ID string
	listaVentas     []ventaCombustible
	FactorUtilidad  float64
	mutex           sync.Mutex
	db              *sql.DB
	preciosBase     map[string]float64
}

// --- Lógica de Sockets ---

func iniciarServidorSockets(dist *distribuidor) {
	socketPort := os.Getenv("SOCKET_PORT")
	if socketPort == "" {
		socketPort = "9001" // Puerto por defecto
	}
	listener, err := net.Listen("tcp", ":"+socketPort)
	if err != nil {
		log.Fatalf("Error al iniciar servidor de sockets: %v", err)
	}
	defer listener.Close()
	fmt.Printf("🔌 Servidor de Sockets escuchando en el puerto %s\n", socketPort)

	surtidorConnections.conns = make(map[net.Conn]bool)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión de surtidor: %v", err)
			continue
		}

		surtidorConnections.Lock()
		surtidorConnections.conns[conn] = true
		surtidorConnections.Unlock()

		go handleSurtidorConnection(conn, dist)
	}
}

func handleSurtidorConnection(conn net.Conn, dist *distribuidor) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Surtidor conectado desde %s\n", remoteAddr)
	defer func() {
		fmt.Printf("Surtidor desconectado: %s\n", remoteAddr)
		surtidorConnections.Lock()
		delete(surtidorConnections.conns, conn)
		surtidorConnections.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return // Termina si hay error o se cierra la conexión
		}
		// "ID:1;LITROS:25.50;TIPO:97"
		parts := strings.Split(strings.TrimSpace(message), ";")
		if len(parts) != 3 {
			continue
		}

		hiloIDStr := strings.Split(parts[0], ":")[1]
		litrosStr := strings.Split(parts[1], ":")[1]
		tipo := strings.Split(parts[2], ":")[1]

		litros, _ := strconv.ParseFloat(litrosStr, 64)

		// Procesar la venta recibida
	realizarVenta(dist, hiloIDStr, tipo, litros)
	}
}

// --- Lógica de Negocio ---

func conectarDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL no configurada")
	}

	var db *sql.DB
	var err error
	maxRetries := 5
	retryDelay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("Intento %d/%d: Error al abrir la conexión a la DB: %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("Intento %d/%d: Error al verificar la conexión a la DB: %v", i+1, maxRetries, err)
			db.Close() // Cerrar la conexión fallida antes de reintentar
			time.Sleep(retryDelay)
			continue
		}

		fmt.Println("✅ Conexión a la base de datos establecida con éxito.")
		return db, nil
	}

	return nil, fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos", maxRetries)
}

// Nueva función para cargar datos al iniciar
func cargarVentasDesdeDB(dist *distribuidor) error {
	if dist.db == nil {
		return fmt.Errorf("La conexión a la base de datos no está disponible")
	}

	rows, err := dist.db.Query("SELECT distribuidor_id, id, hilo_id, tipo_combustible, litros_vendidos, precio_por_litro, total_venta, fecha_venta, numero_factura FROM ventas_combustible ORDER BY id ASC")
	if err != nil {
		return fmt.Errorf("error al consultar ventas: %v", err)
	}
	defer rows.Close()

	dist.mutex.Lock()
	defer dist.mutex.Unlock()

	ventasCargadas := 0
	for rows.Next() {
		var v ventaCombustible
		// El ID de la venta no se escanea porque es autoincremental en la DB
		if err := rows.Scan(&v.DistribuidorID, &v.ID_venta, &v.HiloID, &v.TipoCombustible, &v.LitrosVendidos, &v.PrecioPorLitro, &v.TotalVenta, &v.FechaVenta, &v.NumeroFactura); err != nil {
			log.Printf("Error al escanear fila de venta: %v", err)
			continue
		}
		v.Persistida = true // Marcar como persistida ya que viene de la DB
		dist.listaVentas = append(dist.listaVentas, v)
		ventasCargadas++
	}

	fmt.Printf("📦 %d ventas cargadas desde la base de datos.\n", ventasCargadas)
	return nil
}

func realizarVenta(d *distribuidor, hiloID string, tipoCombustible string, litros float64) {
	d.mutex.Lock()
	precioBase, ok := d.preciosBase[tipoCombustible]
	if !ok {
		precioBase = 1.2 + rand.Float64()*0.8 // Fallback a un precio aleatorio si no está definido
	}
	precioFinal := precioBase * d.FactorUtilidad
	total := litros * precioFinal

	venta := ventaCombustible{
		DistribuidorID:  d.Distribuidor_ID,
		ID_venta:        len(d.listaVentas) + 1,
		HiloID:          hiloID,
		TipoCombustible: tipoCombustible,
		LitrosVendidos:  litros,
		PrecioPorLitro:  precioFinal,
		TotalVenta:      total,
		FechaVenta:      time.Now(),
		NumeroFactura:   fmt.Sprintf("F-%s-%d", d.Distribuidor_ID, len(d.listaVentas)+1),
		Persistida:      false,
	}
	d.listaVentas = append(d.listaVentas, venta)
	d.mutex.Unlock()

	fmt.Printf("Venta recibida de surtidor %s: %.2f L de %s\n", hiloID, litros, tipoCombustible)
}

func actualizarBaseDatos(d *distribuidor) {
	d.mutex.Lock()
	ventasNoPersistidas := make([]ventaCombustible, 0)
	for i, venta := range d.listaVentas {
		if !venta.Persistida {
			ventasNoPersistidas = append(ventasNoPersistidas, venta)
			d.listaVentas[i].Persistida = true
		}
	}
	d.mutex.Unlock()

	if len(ventasNoPersistidas) == 0 {
		return
	}

	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("Error iniciando transacción: %v", err)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO ventas_combustible 
        (distribuidor_id, hilo_id, tipo_combustible, litros_vendidos, precio_por_litro, total_venta, fecha_venta, numero_factura)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		log.Printf("Error preparando statement: %v", err)
		return
	}
	defer stmt.Close()

	for _, venta := range ventasNoPersistidas {
		_, err = stmt.Exec(venta.DistribuidorID, venta.HiloID, venta.TipoCombustible, venta.LitrosVendidos, venta.PrecioPorLitro, venta.TotalVenta, venta.FechaVenta, venta.NumeroFactura)
		if err != nil {
			log.Printf("Error insertando venta %s: %v", venta.NumeroFactura, err)
			continue
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error confirmando transacción: %v", err)
		return
	}
	fmt.Printf("✅ Base de datos actualizada: %d ventas insertadas\n", len(ventasNoPersistidas))
}

func updatePricesHandler(dist *distribuidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var precios map[string]float64
		if err := json.NewDecoder(r.Body).Decode(&precios); err != nil {
			fmt.Printf("Error al decodificar el cuerpo de la solicitud de precios: %v\n", err)
			http.Error(w, "Cuerpo de petición inválido", http.StatusBadRequest)
			return
		}

		fmt.Printf("Recibida solicitud de actualización de precios: %+v\n", precios)

		dist.mutex.Lock()
		for tipo, precio := range precios {
			dist.preciosBase[tipo] = precio
		}
		dist.mutex.Unlock()

		fmt.Println("Precios actualizados correctamente")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Precios actualizados correctamente"})
	}
}

// ===== ENDPOINTS HTTP =====

func statusHandler(dist *distribuidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dist.mutex.Lock()
		surtidorConnections.Lock()
		response := map[string]interface{}{
			"status":             "online",
			"distribuidor_id":    dist.Distribuidor_ID,
			"surtidores_activos": len(surtidorConnections.conns),
			"total_ventas":       len(dist.listaVentas),
			"timestamp":          time.Now(),
		}
		surtidorConnections.Unlock()
		dist.mutex.Unlock()
		json.NewEncoder(w).Encode(response)
	}
}

func ventasRecientesHandler(dist *distribuidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		limitStr := r.URL.Query().Get("limit")
		limit := 20
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}
		dist.mutex.Lock()
		totalVentas := len(dist.listaVentas)
		var ventasRecientes []ventaCombustible
		if totalVentas > 0 {
			start := totalVentas - limit
			if start < 0 {
				start = 0
			}
			ventasRecientes = make([]ventaCombustible, totalVentas-start)
			copy(ventasRecientes, dist.listaVentas[start:])
			for i, j := 0, len(ventasRecientes)-1; i < j; i, j = i+1, j-1 {
				ventasRecientes[i], ventasRecientes[j] = ventasRecientes[j], ventasRecientes[i]
			}
		}
		dist.mutex.Unlock()
		json.NewEncoder(w).Encode(ventasRecientes)
	}
}

// Endpoint para obtener estadísticas del distribuidor
func statsHandler(dist *distribuidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dist.mutex.Lock()
		defer dist.mutex.Unlock()

		surtidorConnections.Lock()
		defer surtidorConnections.Unlock()

		// Calcular estadísticas
		totalVentas := len(dist.listaVentas)
		var totalIngresos float64
		var ventasHoy int

		hoy := time.Now().Truncate(24 * time.Hour)

		for _, venta := range dist.listaVentas {
			totalIngresos += venta.TotalVenta
			if venta.FechaVenta.Truncate(24 * time.Hour).Equal(hoy) {
				ventasHoy++
			}
		}

		response := map[string]interface{}{
			"distribuidor_id":    dist.Distribuidor_ID,
			"total_ventas":       totalVentas,
			"ventas_hoy":         ventasHoy,
			"total_ingresos":     totalIngresos,
			"factor_utilidad":    dist.FactorUtilidad,
			"surtidores_activos": len(surtidorConnections.conns), // Corregido
			"precios_base":       dist.preciosBase,
			"timestamp":          time.Now(),
		}

		json.NewEncoder(w).Encode(response)
	}
}

// NUEVO endpoint para controlar los surtidores
func controlSurtidoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var payload struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Cuerpo de petición inválido", http.StatusBadRequest)
		return
	}

	var command string
	switch payload.Action {
	case "pause":
		command = "CMD:PAUSE\n"
	case "resume":
		command = "CMD:RESUME\n"
	default:
		http.Error(w, "Acción no válida. Usar 'pause' o 'resume'", http.StatusBadRequest)
		return
	}

	surtidorConnections.Lock()
	defer surtidorConnections.Unlock()

	if len(surtidorConnections.conns) == 0 {
		fmt.Println("No hay surtidores conectados para enviar comando.")
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay surtidores conectados."})
		return
	}

	sentCount := 0
	for conn := range surtidorConnections.conns {
		_, err := conn.Write([]byte(command))
		if err != nil {
			log.Printf("Error enviando comando a %s: %v", conn.RemoteAddr(), err)
		} else {
			sentCount++
		}
	}

	fmt.Printf("Comando '%s' enviado a %d surtidores.\n", payload.Action, sentCount)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      fmt.Sprintf("Comando '%s' enviado", payload.Action),
		"surtidores_afectados": sentCount,
	})
}

func setupHTTPServer(dist *distribuidor) *http.Server {
	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r.HandleFunc("/status", statusHandler(dist)).Methods("GET")
	r.HandleFunc("/ventas/recientes", ventasRecientesHandler(dist)).Methods("GET")
	r.HandleFunc("/stats", statsHandler(dist)).Methods("GET") // Reactivado
	r.HandleFunc("/surtidores/control", controlSurtidoresHandler).Methods("POST")
	r.HandleFunc("/prices", updatePricesHandler(dist)).Methods("PUT")

	handler := c.Handler(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &http.Server{Addr: ":" + port, Handler: handler}
}

func main() {
	fmt.Println("--- Arrancando Backend Distribuidor (Versión Servidor de Sockets) ---")

	db, err := conectarDB()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos después de varios intentos: %v", err)
	}
	defer db.Close()

	factorUtilidad := 1.10 + rand.Float64()*0.20
	id_distribuidor := os.Getenv("DISTRIBUIDOR_ID")
	if id_distribuidor == "" {
		id_distribuidor = "DISTRIBUIDOR_DEFAULT"
	}

	dist := &distribuidor{
		Distribuidor_ID: id_distribuidor,
		FactorUtilidad:  factorUtilidad,
		listaVentas:     make([]ventaCombustible, 0),
		db:              db,
		preciosBase:     map[string]float64{
			"93":       1.0,
			"95":       1.1,
			"97":       1.2,
			"DIESEL":   0.9,
			"KEROSENE": 0.8,
		},
	}

	// Cargar datos existentes desde la base de datos
	if db != nil {
		if err := cargarVentasDesdeDB(dist); err != nil {
			log.Printf("Error al cargar ventas desde la DB: %v", err)
		}
	}

	fmt.Printf("Distribuidor %s iniciado con factor de utilidad: %.2f\n", dist.Distribuidor_ID, factorUtilidad)

	// Iniciar servidor de sockets en una goroutine
	go iniciarServidorSockets(dist)

	// Iniciar servidor HTTP en una goroutine
	server := setupHTTPServer(dist)
	go func() {
		fmt.Printf("🌐 Servidor HTTP iniciado en puerto %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en servidor HTTP: %v", err)
		}
	}()

	// Goroutine para actualizar base de datos
	if db != nil {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				actualizarBaseDatos(dist)
			}
		}()
	}

	fmt.Println("\n🚀 Sistema iniciado correctamente!")
	// Mantener la aplicación corriendo
	select {}
}