package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
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
}

type distribuidor struct {
	Distribuidor_ID string
	listaVentas     []ventaCombustible
	FactorUtilidad  float64
	mutex           sync.Mutex
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

func main() {
	fmt.Println("--- Arrancando Backend Distribuidor ---")

	rand.Seed(time.Now().UnixNano())

	// Crear distribuidor con factor de utilidad aleatorio
	factorUtilidad := 1.10 + rand.Float64()*0.20 // entre 1.10 (10%) y 1.30 (30%) de margen

	dist := &distribuidor{
		Distribuidor_ID: "DIST_001",
		FactorUtilidad:  factorUtilidad,
		listaVentas:     make([]ventaCombustible, 0),
	}

	fmt.Printf("Distribuidor %s iniciado con factor de utilidad: %.2f (%.0f%% margen)\n",
		dist.Distribuidor_ID, factorUtilidad, (factorUtilidad-1)*100)

	// Configuración de surtidores
	numSurtidores := 5 // Cambia este número para crear N surtidores

	var wg sync.WaitGroup

	// Crear y ejecutar surtidores en paralelo (infinitamente)
	for i := 1; i <= numSurtidores; i++ {
		wg.Add(1)
		surtidorID := fmt.Sprintf("surtidor-%d", i)
		go surtidor(dist, surtidorID, &wg)
	}

	// Variables de entorno
	dbURL := os.Getenv("DATABASE_URL")
	empresaHost := os.Getenv("EMPRESA_HOST")
	surtidorHost := os.Getenv("SURTIDOR_HOST")
	fmt.Printf("DATABASE_URL: %s\n", dbURL)
	fmt.Printf("EMPRESA_HOST: %s\n", empresaHost)
	fmt.Printf("SURTIDOR_HOST: %s\n", surtidorHost)

	fmt.Println("\nPresiona Ctrl+C para detener...")

	// Esperar indefinidamente (los surtidores nunca terminan)
	wg.Wait()
}
