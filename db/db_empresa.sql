-- Inicialización específica para la base `db_empresa`
-- Rellena este archivo con el esquema y datos necesarios para la empresa
-- Referencia: ver el PDF del proyecto (proporcionar extracto si quieres que lo implemente exactamente)

-- Ejemplo mínimo: crear tabla de control de versión de esquema
CREATE TABLE IF NOT EXISTS schema_version (
  id SERIAL PRIMARY KEY,
  version TEXT NOT NULL,
  applied_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- TODO: Añadir tablas reales para la empresa (precios, productos, etc.) según el PDF
