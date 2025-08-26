-- Create B3 table
CREATE TABLE IF NOT EXISTS b3 (
    id SERIAL PRIMARY KEY,
    data_negocio TIMESTAMP NOT NULL,
    codigo_instrumento VARCHAR(50) NOT NULL,
    preco_negocio DECIMAL(15,2) NOT NULL,
    quantidade_negociada INTEGER NOT NULL,
    hora_fechamento TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on data_negocio for better query performance
CREATE INDEX IF NOT EXISTS idx_b3_data_negocio ON b3(data_negocio);

-- Create index on codigo_instrumento for better query performance
CREATE INDEX IF NOT EXISTS idx_b3_codigo_instrumento ON b3(codigo_instrumento);

-- Create index on hora_fechamento for better query performance
CREATE INDEX IF NOT EXISTS idx_b3_hora_fechamento ON b3(hora_fechamento);
