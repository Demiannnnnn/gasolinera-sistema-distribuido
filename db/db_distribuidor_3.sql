-- Inicialización específica para la base `db_distribuidor_3`
-- Rellena este archivo con el esquema y datos necesarios para el distribuidor 3
-- Referencia: ver el PDF del proyecto (proporcionar extracto si quieres que lo implemente exactamente)

CREATE TABLE IF NOT EXISTS schema_version (
  id SERIAL PRIMARY KEY,
  version TEXT NOT NULL,
  applied_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- TODO: Añadir tablas específicas del distribuidor (tanques, surtidores, transacciones, etc.)
