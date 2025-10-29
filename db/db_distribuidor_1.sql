-- Inicialización específica para la base `db_distribuidor_1`
-- Rellena este archivo con el esquema y datos necesarios para el distribuidor 1
-- Referencia: ver el PDF del proyecto (proporcionar extracto si quieres que lo implemente exactamente)

-- TODO: Reemplazar con esquema específico del PDF

CREATE TABLE IF NOT EXISTS ventas_combustible (
    id SERIAL PRIMARY KEY,
    distribuidor_id VARCHAR(50) NOT NULL,
    hilo_id VARCHAR(50) NOT NULL,
    tipo_combustible VARCHAR(20) NOT NULL,
    litros_vendidos DECIMAL(10,3) NOT NULL,
    precio_por_litro DECIMAL(10,3) NOT NULL,
    total_venta DECIMAL(10,2) NOT NULL,
    fecha_venta TIMESTAMP NOT NULL DEFAULT NOW(),
    numero_factura VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ventas_distribuidor ON ventas_combustible(distribuidor_id);
CREATE INDEX IF NOT EXISTS idx_ventas_fecha ON ventas_combustible(fecha_venta);

-- TODO: Añadir tablas específicas del distribuidor (tanques, surtidores, transacciones, etc.)
