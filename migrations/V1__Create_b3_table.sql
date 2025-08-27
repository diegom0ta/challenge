CREATE TABLE IF NOT EXISTS b3 (
    id SERIAL PRIMARY KEY,
    data_referencia TIMESTAMP,
    acao_atualizacao VARCHAR(50),
    data_negocio VARCHAR(50) NOT NULL,
    codigo_instrumento VARCHAR(50) NOT NULL,
    preco_negocio DECIMAL(15,2) NOT NULL,
    quantidade_negociada INTEGER NOT NULL,
    hora_fechamento INTEGER,
    codigo_identificador_negocio VARCHAR(50),
    tipo_sessao_pregao INTEGER,
    codigo_participante_comprador VARCHAR(50),
    codigo_participante_vendedor VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_b3_data_negocio ON b3(data_negocio);
CREATE INDEX IF NOT EXISTS idx_b3_codigo_instrumento ON b3(codigo_instrumento);