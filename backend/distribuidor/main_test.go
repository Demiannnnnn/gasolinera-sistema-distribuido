package main

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func nuevoDistribuidor(factor float64) *distribuidor {
	return &distribuidor{
		Distribuidor_ID: "DIST_TEST",
		FactorUtilidad:  factor,
		listaVentas:     make([]ventaCombustible, 0),
	}
}

func casiIgual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestGetCombustibleTipo(t *testing.T) {
	got := GetCombustibleTipo()
	want := []string{"93", "95", "97", "DIESEL", "KEROSENE", "PETROLEO"}
	if len(got) != len(want) {
		t.Fatalf("esperaba %d combustibles, obtuve %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("posición %d: esperaba %q, obtuve %q", i, want[i], got[i])
		}
	}
}

func TestCombustibleAleatorio(t *testing.T) {
	casos := []struct {
		n    int
		want string
	}{
		{0, Fuel93},
		{1, Fuel95},
		{5, FuelPetroleo},
		{6, Fuel93},  // da la vuelta
		{13, Fuel95}, // 13 % 6 = 1
		{-1, FuelPetroleo},
		{-6, Fuel93},
		{-7, FuelPetroleo},
	}
	for _, c := range casos {
		if got := CombustibleAleatorio(c.n); got != c.want {
			t.Errorf("CombustibleAleatorio(%d) = %q, esperaba %q", c.n, got, c.want)
		}
	}
}

func TestCombustibleAleatorioSiempreValido(t *testing.T) {
	validos := map[string]bool{}
	for _, f := range GetCombustibleTipo() {
		validos[f] = true
	}
	for n := -100; n <= 100; n++ {
		if !validos[CombustibleAleatorio(n)] {
			t.Fatalf("combustible inválido para n=%d", n)
		}
	}
}

func TestRealizarVentaAplicaFactorYRangos(t *testing.T) {
	d := nuevoDistribuidor(1.25)
	for i := 0; i < 200; i++ {
		realizarVenta(d, "surtidor-1")
	}
	if len(d.listaVentas) != 200 {
		t.Fatalf("esperaba 200 ventas, obtuve %d", len(d.listaVentas))
	}
	for i, v := range d.listaVentas {
		if v.ID_venta != i+1 {
			t.Errorf("venta %d: ID esperado %d, obtuve %d", i, i+1, v.ID_venta)
		}
		if v.DistribuidorID != "DIST_TEST" || v.HiloID != "surtidor-1" {
			t.Errorf("venta %d: identificadores incorrectos %+v", i, v)
		}
		if v.LitrosVendidos < 5 || v.LitrosVendidos > 55 {
			t.Errorf("venta %d: litros fuera de rango: %f", i, v.LitrosVendidos)
		}
		// precio base entre 1.2 y 2.0, multiplicado por el factor 1.25
		if v.PrecioPorLitro < 1.2*1.25 || v.PrecioPorLitro > 2.0*1.25 {
			t.Errorf("venta %d: precio fuera de rango: %f", i, v.PrecioPorLitro)
		}
		if !casiIgual(v.TotalVenta, v.LitrosVendidos*v.PrecioPorLitro) {
			t.Errorf("venta %d: total %f != litros*precio %f", i, v.TotalVenta, v.LitrosVendidos*v.PrecioPorLitro)
		}
		if want := "F-DIST_TEST-" + itoa(i+1); v.NumeroFactura != want {
			t.Errorf("venta %d: factura %q, esperaba %q", i, v.NumeroFactura, want)
		}
		if v.Persistida {
			t.Errorf("venta %d: no debería estar persistida", i)
		}
	}
}

func TestRealizarVentaConcurrenteSinIDsDuplicados(t *testing.T) {
	d := nuevoDistribuidor(1.1)
	var wg sync.WaitGroup
	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				realizarVenta(d, "s")
			}
		}()
	}
	wg.Wait()
	if len(d.listaVentas) != 500 {
		t.Fatalf("esperaba 500 ventas, obtuve %d", len(d.listaVentas))
	}
	vistos := map[int]bool{}
	facturas := map[string]bool{}
	for _, v := range d.listaVentas {
		if vistos[v.ID_venta] || facturas[v.NumeroFactura] {
			t.Fatalf("ID o factura duplicada: %d %s", v.ID_venta, v.NumeroFactura)
		}
		vistos[v.ID_venta] = true
		facturas[v.NumeroFactura] = true
	}
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func TestStatusHandler(t *testing.T) {
	d := nuevoDistribuidor(1.2)
	realizarVenta(d, "s1")
	realizarVenta(d, "s2")

	rec := httptest.NewRecorder()
	statusHandler(d)(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "online" || body["distribuidor_id"] != "DIST_TEST" {
		t.Errorf("respuesta inesperada: %v", body)
	}
	if body["total_ventas"].(float64) != 2 {
		t.Errorf("total_ventas = %v", body["total_ventas"])
	}
	if !casiIgual(body["factor_utilidad"].(float64), 1.2) {
		t.Errorf("factor_utilidad = %v", body["factor_utilidad"])
	}
}

func obtenerRecientes(t *testing.T, d *distribuidor, query string) []ventaCombustible {
	t.Helper()
	rec := httptest.NewRecorder()
	ventasRecientesHandler(d)(rec, httptest.NewRequest(http.MethodGet, "/ventas/recientes"+query, nil))
	var ventas []ventaCombustible
	if err := json.Unmarshal(rec.Body.Bytes(), &ventas); err != nil {
		t.Fatalf("JSON inválido %q: %v", rec.Body.String(), err)
	}
	return ventas
}

func TestVentasRecientesSinVentas(t *testing.T) {
	d := nuevoDistribuidor(1.1)
	if v := obtenerRecientes(t, d, ""); len(v) != 0 {
		t.Errorf("esperaba lista vacía, obtuve %d", len(v))
	}
}

func TestVentasRecientesLimiteYOrden(t *testing.T) {
	d := nuevoDistribuidor(1.1)
	for i := 0; i < 30; i++ {
		realizarVenta(d, "s")
	}

	porDefecto := obtenerRecientes(t, d, "")
	if len(porDefecto) != 20 {
		t.Fatalf("límite por defecto: esperaba 20, obtuve %d", len(porDefecto))
	}
	if porDefecto[0].ID_venta != 30 || porDefecto[19].ID_venta != 11 {
		t.Errorf("orden inesperado: primero %d, último %d", porDefecto[0].ID_venta, porDefecto[19].ID_venta)
	}

	cinco := obtenerRecientes(t, d, "?limit=5")
	if len(cinco) != 5 || cinco[0].ID_venta != 30 || cinco[4].ID_venta != 26 {
		t.Errorf("limit=5 inesperado: %d elementos", len(cinco))
	}

	todas := obtenerRecientes(t, d, "?limit=100")
	if len(todas) != 30 || todas[29].ID_venta != 1 {
		t.Errorf("limit mayor al total: esperaba 30, obtuve %d", len(todas))
	}

	for _, q := range []string{"?limit=0", "?limit=-3", "?limit=abc"} {
		if v := obtenerRecientes(t, d, q); len(v) != 20 {
			t.Errorf("%s debería usar el límite por defecto, obtuve %d", q, len(v))
		}
	}
}

func TestVentasRecientesNoExponePersistida(t *testing.T) {
	d := nuevoDistribuidor(1.1)
	realizarVenta(d, "s")
	rec := httptest.NewRecorder()
	ventasRecientesHandler(d)(rec, httptest.NewRequest(http.MethodGet, "/ventas/recientes", nil))
	if strings.Contains(strings.ToLower(rec.Body.String()), "persistida") {
		t.Errorf("el campo Persistida no debería serializarse: %s", rec.Body.String())
	}
}

func TestStatsHandler(t *testing.T) {
	d := nuevoDistribuidor(1.15)
	ahora := time.Now()
	d.listaVentas = []ventaCombustible{
		{TotalVenta: 10.5, FechaVenta: ahora},
		{TotalVenta: 20.25, FechaVenta: ahora},
		{TotalVenta: 100, FechaVenta: ahora.Add(-72 * time.Hour)},
	}

	rec := httptest.NewRecorder()
	statsHandler(d)(rec, httptest.NewRequest(http.MethodGet, "/stats", nil))
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["total_ventas"].(float64) != 3 {
		t.Errorf("total_ventas = %v", body["total_ventas"])
	}
	if body["ventas_hoy"].(float64) != 2 {
		t.Errorf("ventas_hoy = %v", body["ventas_hoy"])
	}
	if !casiIgual(body["total_ingresos"].(float64), 130.75) {
		t.Errorf("total_ingresos = %v", body["total_ingresos"])
	}
	if body["surtidores_activos"].(float64) != 5 {
		t.Errorf("surtidores_activos = %v", body["surtidores_activos"])
	}
}

func TestSetupHTTPServerPuertos(t *testing.T) {
	casos := []struct {
		port, id, want string
	}{
		{"", "DISTRIBUIDOR_1", ":8080"},
		{"", "DISTRIBUIDOR_2", ":8083"},
		{"", "DISTRIBUIDOR_3", ":8084"},
		{"", "OTRO", ":8080"},
		{"9000", "DISTRIBUIDOR_2", ":9000"},
	}
	for _, c := range casos {
		t.Setenv("PORT", c.port)
		t.Setenv("DISTRIBUIDOR_ID", c.id)
		if got := setupHTTPServer(nuevoDistribuidor(1.1)).Addr; got != c.want {
			t.Errorf("PORT=%q ID=%q: Addr=%q, esperaba %q", c.port, c.id, got, c.want)
		}
	}
}

func TestRutasHTTP(t *testing.T) {
	t.Setenv("PORT", "0")
	srv := setupHTTPServer(nuevoDistribuidor(1.1))
	for _, ruta := range []string{"/status", "/ventas/recientes", "/stats"} {
		rec := httptest.NewRecorder()
		srv.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, ruta, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d", ruta, rec.Code)
		}
	}
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/status", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /status = %d, esperaba 405", rec.Code)
	}
	rec = httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/no-existe", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET /no-existe = %d, esperaba 404", rec.Code)
	}
}

func TestConectarDBSinURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	db, err := conectarDB()
	if err == nil || db != nil {
		t.Fatalf("esperaba error sin DATABASE_URL, obtuve db=%v err=%v", db, err)
	}
}

const insertRegex = `INSERT INTO ventas_combustible`

func TestActualizarBaseDatosInsertaSoloNoPersistidas(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	d := nuevoDistribuidor(1.1)
	d.db = db
	d.listaVentas = []ventaCombustible{
		{NumeroFactura: "F-1", Persistida: true},
		{NumeroFactura: "F-2", DistribuidorID: "DIST_TEST", HiloID: "s1", TipoCombustible: "95", LitrosVendidos: 10, PrecioPorLitro: 2, TotalVenta: 20},
		{NumeroFactura: "F-3"},
	}

	mock.ExpectBegin()
	prep := mock.ExpectPrepare(insertRegex)
	prep.ExpectExec().WithArgs("DIST_TEST", "s1", "95", 10.0, 2.0, 20.0, sqlmock.AnyArg(), "F-2").
		WillReturnResult(sqlmock.NewResult(1, 1))
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "F-3").
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	actualizarBaseDatos(d)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
	for _, v := range d.listaVentas {
		if !v.Persistida {
			t.Errorf("venta %s debería quedar marcada como persistida", v.NumeroFactura)
		}
	}
}

func TestActualizarBaseDatosSinPendientesNoTocaDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	d := nuevoDistribuidor(1.1)
	d.db = db
	d.listaVentas = []ventaCombustible{{Persistida: true}}

	actualizarBaseDatos(d) // sin expectativas: cualquier llamada a la DB falla

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestActualizarBaseDatosErrorEnBeginNoEntraEnPanico(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	d := nuevoDistribuidor(1.1)
	d.db = db
	d.listaVentas = []ventaCombustible{{NumeroFactura: "F-1"}}
	mock.ExpectBegin().WillReturnError(errors.New("db caída"))

	actualizarBaseDatos(d)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

// Documenta el comportamiento actual. Ojo: en PostgreSQL real, un error dentro de
// la transacción la aborta y las inserciones siguientes también fallan.
func TestActualizarBaseDatosErrorEnUnaVentaContinuaConLasDemas(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	d := nuevoDistribuidor(1.1)
	d.db = db
	d.listaVentas = []ventaCombustible{{NumeroFactura: "F-1"}, {NumeroFactura: "F-2"}}

	mock.ExpectBegin()
	prep := mock.ExpectPrepare(insertRegex)
	prep.ExpectExec().WillReturnError(errors.New("violación de restricción"))
	prep.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	actualizarBaseDatos(d)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
