
CREATE TABLE IF NOT EXISTS schema_version (
  id SERIAL PRIMARY KEY,
  
  version TEXT NOT NULL,
  applied_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

