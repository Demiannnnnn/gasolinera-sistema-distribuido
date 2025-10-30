package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	Fuel93       = "93"
	Fuel95       = "95"
	Fuel97       = "97"
	FuelDiesel   = "DIESEL"
	FuelKerosene = "KEROSENE"
	FuelPetroleo = "PETROLEO"
)

func GetCombustibleTipo() []string {
	return []string{
		Fuel93,
		Fuel95,
		Fuel97,
		FuelDiesel,
		FuelKerosene,
		FuelPetroleo,
	}
}

// recibe un numero aleatorio y retorna un tipo de combustible
func CombustibleAleatorio(n int) string {
	fuels := GetCombustibleTipo()
	if len(fuels) == 0 {
		return ""
	}
	idx := n % len(fuels)
	if idx < 0 {
		idx += len(fuels)
	}
	return fuels[idx]
}

type ventaCombustible struct {
	DistribuidorID  string    `json:"distribuidor_id"`
	ID_venta        int       `json:"id"`
	HiloID          string    `json:"hilo_id"`
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
}

func conectarDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL no configurada")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando a DB: %v", err)
	}

	// Verificar conexión
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando conexión DB: %v", err)
	}

	return db, nil
}

func actualizarBaseDatos(d *distribuidor) {
	//revisar ventas no persistidas (que no hayan sido ingresadas)
	d.mutex.Lock()
	ventasNoPersistidas := make([]ventaCombustible, 0)

	// Filtrar ventas no persistidas
	for i, venta := range d.listaVentas {
		if !venta.Persistida {
			ventasNoPersistidas = append(ventasNoPersistidas, venta)
			d.listaVentas[i].Persistida = true // marcar como persistida
		}
	}
	d.mutex.Unlock()

	if len(ventasNoPersistidas) == 0 {
		fmt.Printf("No hay ventas nuevas para persistir\n")
		return
	}

	// Insertar en lote
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

	ventasInsertadas := 0
	for _, venta := range ventasNoPersistidas {
		_, err = stmt.Exec(
			venta.DistribuidorID,
			venta.HiloID,
			venta.TipoCombustible,
			venta.LitrosVendidos,
			venta.PrecioPorLitro,
			venta.TotalVenta,
			venta.FechaVenta,
			venta.NumeroFactura,
		)
		if err != nil {
			log.Printf("Error insertando venta %s: %v", venta.NumeroFactura, err)
			continue
		}
		ventasInsertadas++
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error confirmando transacción: %v", err)
		return
	}

	fmt.Printf("✅ Base de datos actualizada: %d ventas insertadas\n", ventasInsertadas)
}

func realizarVenta(d *distribuidor, hiloID string) {
	// Generar venta aleatoria
	tipoCombustible := CombustibleAleatorio(rand.Intn(100))
	litros := 5.0 + rand.Float64()*50.0    // entre 5 y 55 litros
	precioBase := 1.2 + rand.Float64()*0.8 // precio base entre 1.2 y 2.0

	// Aplicar factor de utilidad del distribuidor
	precioFinal := precioBase * d.FactorUtilidad
	total := litros * precioFinal

	// Bloquear para acceso seguro a la lista
	d.mutex.Lock()
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

	// Guardar en la lista del distribuidor
	d.listaVentas = append(d.listaVentas, venta)
	d.mutex.Unlock()

	fmt.Printf("Venta realizada por %s (hilo %s): %.2f L de %s a $%.2f/L = $%.2f\n",
		d.Distribuidor_ID, hiloID, litros, tipoCombustible, precioFinal, total)
}

func surtidor(d *distribuidor, surtidorID string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Surtidor %s iniciado - vendiendo indefinidamente...\n", surtidorID)

	for {
		realizarVenta(d, surtidorID)
		// Simular tiempo entre ventas (1-3 segundos)
		time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
	}
}

// ===== ENDPOINTS HTTP =====

// Endpoint para verificar estado del distribuidor
func statusHandler(dist *distribuidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dist.mutex.Lock()
		response := map[string]interface{}{
			"status":          "online",
			"distribuidor_id": dist.Distribuidor_ID,
			"factor_utilidad": dist.FactorUtilidad,
			"total_ventas":    len(dist.listaVentas),
			"timestamp":       time.Now(),
		}
		dist.mutex.Unlock()

		json.NewEncoder(w).Encode(response)
	}
}

// Endpoint para obtener ventas recientes
func ventasRecientesHandler(dist *distribuidor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Obtener parámetro de límite (opcional)
		limitStr := r.URL.Query().Get("limit")
		limit := 20 // por defecto
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		dist.mutex.Lock()
		totalVentas := len(dist.listaVentas)

		// Obtener las últimas ventas
		var ventasRecientes []ventaCombustible
		if totalVentas > 0 {
			start := totalVentas - limit
			if start < 0 {
				start = 0
			}
			ventasRecientes = make([]ventaCombustible, totalVentas-start)
			copy(ventasRecientes, dist.listaVentas[start:])

			// Invertir el orden para que las más recientes aparezcan primero
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

		// Calcular estadísticas
		totalVentas := len(dist.listaVentas)
		var totalIngresos float64
		var ventasHoy int
		surtidoresActivos := 5 // Hardcoded por ahora

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
			"surtidores_activos": surtidoresActivos,
			"timestamp":          time.Now(),
		}

		dist.mutex.Unlock()

		json.NewEncoder(w).Encode(response)
	}
}

// Configurar el servidor HTTP
func setupHTTPServer(dist *distribuidor) *http.Server {
	r := mux.NewRouter()

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Rutas
	r.HandleFunc("/status", statusHandler(dist)).Methods("GET")
	r.HandleFunc("/ventas/recientes", ventasRecientesHandler(dist)).Methods("GET")
	r.HandleFunc("/stats", statsHandler(dist)).Methods("GET")

	// Aplicar CORS
	handler := c.Handler(r)

	// Obtener puerto de variable de entorno o usar puerto por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto si no se especifica
	}

	// Configurar servidor
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return server
}

func main() {
	fmt.Println("--- Arrancando Backend Distribuidor ---")

	// Conectar a la base de datos
	db, err := conectarDB()
	if err != nil {
		log.Printf("Advertencia: No se pudo conectar a la base de datos: %v", err)
		// Continuar sin DB para desarrollo
	}
	if db != nil {
		defer db.Close()
	}

	// Crear distribuidor con factor de utilidad aleatorio
	factorUtilidad := 1.10 + rand.Float64()*0.20 // entre 1.10 (10%) y 1.30 (30%) de margen
	id_distribuidor := os.Getenv("DISTRIBUIDOR_ID")
	if id_distribuidor == "" {
		id_distribuidor = "DISTRIBUIDOR_DEFAULT"
	}

	dist := &distribuidor{
		Distribuidor_ID: id_distribuidor,
		FactorUtilidad:  factorUtilidad,
		listaVentas:     make([]ventaCombustible, 0),
		db:              db,
	}

	fmt.Printf("Distribuidor %s iniciado con factor de utilidad: %.2f (%.0f%% margen)\n",
		dist.Distribuidor_ID, factorUtilidad, (factorUtilidad-1)*100)

	// Configurar servidor HTTP
	server := setupHTTPServer(dist)

	// Configuración de surtidores
	numSurtidores := 5 // Cambia este número para crear N surtidores

	var wg sync.WaitGroup

	// Crear y ejecutar surtidores en paralelo (infinitamente)
	for i := 1; i <= numSurtidores; i++ {
		wg.Add(1)
		surtidorID := fmt.Sprintf("surtidor-%d", i)
		go surtidor(dist, surtidorID, &wg)
	}

	// Goroutine para actualizar base de datos cada 5 segundos (solo si hay DB)
	if db != nil {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				actualizarBaseDatos(dist)
			}
		}()
	}

	// Iniciar servidor HTTP en goroutine
	go func() {
		fmt.Printf("🌐 Servidor HTTP iniciado en puerto %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error en servidor HTTP: %v", err)
		}
	}()

	// Variables de entorno
	dbURL := os.Getenv("DATABASE_URL")
	empresaHost := os.Getenv("EMPRESA_HOST")
	fmt.Printf("DATABASE_URL: %s\n", dbURL)
	fmt.Printf("EMPRESA_HOST: %s\n", empresaHost)

	fmt.Println("\n🚀 Sistema iniciado correctamente!")
	fmt.Println("📊 Endpoints disponibles:")
	fmt.Printf("   - GET %s/status\n", server.Addr)
	fmt.Printf("   - GET %s/ventas/recientes\n", server.Addr)
	fmt.Printf("   - GET %s/stats\n", server.Addr)
	fmt.Println("\nPresiona Ctrl+C para detener...")

	// Esperar indefinidamente (los surtidores nunca terminan)
	wg.Wait()
}
