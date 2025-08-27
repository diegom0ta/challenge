package repository

import (
	"database/sql"
	"fmt"
	"time"

	"challenge/models"
)

// B3Repository handles B3 data operations
type B3Repository struct {
	db *sql.DB
}

// NewB3Repository creates a new B3 repository
func NewB3Repository(db *sql.DB) *B3Repository {
	return &B3Repository{db: db}
}

// Insert inserts a single B3 record into the database
func (r *B3Repository) Insert(b3 *models.B3) error {
	query := `
		INSERT INTO b3 (data_referencia, acao_atualizacao, data_negocio, codigo_instrumento, 
		               preco_negocio, quantidade_negociada, hora_fechamento, codigo_identificador_negocio,
		               tipo_sessao_pregao, codigo_participante_comprador, codigo_participante_vendedor)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		b3.DataReferencia,
		b3.AcaoAtualizacao,
		b3.DataNegocio,
		b3.CodigoInstrumento,
		b3.PrecoNegocio,
		b3.QuantidadeNegociada,
		b3.HoraFechamento,
		b3.CodigoIdentificadorNegocio,
		b3.TipoSessaoPregao,
		b3.CodigoParticipanteComprador,
		b3.CodigoParticipanteVendedor,
	)

	if err != nil {
		return fmt.Errorf("failed to insert B3 record: %v", err)
	}

	return nil
}

// InsertBatch inserts multiple B3 records in a transaction
func (r *B3Repository) InsertBatch(b3Records []*models.B3) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		err := tx.Rollback() // Will be ignored if transaction is committed
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	// Prepare the statement
	stmt, err := tx.Prepare(`
		INSERT INTO b3 (data_referencia, acao_atualizacao, data_negocio, codigo_instrumento, 
		               preco_negocio, quantidade_negociada, hora_fechamento, codigo_identificador_negocio,
		               tipo_sessao_pregao, codigo_participante_comprador, codigo_participante_vendedor)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Insert all records
	for i, b3 := range b3Records {
		_, err = stmt.Exec(
			b3.DataReferencia,
			b3.AcaoAtualizacao,
			b3.DataNegocio,
			b3.CodigoInstrumento,
			b3.PrecoNegocio,
			b3.QuantidadeNegociada,
			b3.HoraFechamento,
			b3.CodigoIdentificadorNegocio,
			b3.TipoSessaoPregao,
			b3.CodigoParticipanteComprador,
			b3.CodigoParticipanteVendedor,
		)
		if err != nil {
			return fmt.Errorf("failed to insert B3 record %d: %v", i+1, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Printf("Successfully inserted %d B3 records\n", len(b3Records))
	return nil
}

// ParseCSVRowToB3 parses a CSV row into a B3 model (expects 11 columns)
func ParseCSVRowToB3(row []string) (*models.B3, error) {
	if len(row) != 11 {
		return nil, fmt.Errorf("expected exactly 11 columns, got %d", len(row))
	}

	// Parse data_referencia (timestamp)
	dataReferencia, err := time.Parse("2006-01-02", row[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse data_referencia: %v", err)
	}

	// Parse acao_atualizacao (string)
	var acaoAtualizacao string
	if _, err := fmt.Sscanf(row[1], "%s", &acaoAtualizacao); err != nil {
		return nil, fmt.Errorf("failed to parse acao_atualizacao: %v", err)
	}

	// Parse data_negocio (string)
	var dataNegocio string
	if _, err := fmt.Sscanf(row[2], "%s", &dataNegocio); err != nil {
		return nil, fmt.Errorf("failed to parse data_negocio: %v", err)
	}

	// Parse preco_negocio (float)
	var precoNegocio float64
	if _, err := fmt.Sscanf(row[4], "%f", &precoNegocio); err != nil {
		return nil, fmt.Errorf("failed to parse preco_negocio: %v", err)
	}

	// Parse quantidade_negociada (int)
	var quantidadeNegociada int
	if _, err := fmt.Sscanf(row[5], "%d", &quantidadeNegociada); err != nil {
		return nil, fmt.Errorf("failed to parse quantidade_negociada: %v", err)
	}

	// Parse hora_fechamento (integer)
	var horaFechamento int
	if _, err := fmt.Sscanf(row[6], "%d", &horaFechamento); err != nil {
		return nil, fmt.Errorf("failed to parse hora_fechamento: %v", err)
	}

	// Parse tipo_sessao_pregao (int)
	var tipoSessaoPregao int
	if _, err := fmt.Sscanf(row[8], "%d", &tipoSessaoPregao); err != nil {
		return nil, fmt.Errorf("failed to parse tipo_sessao_pregao: %v", err)
	}

	return &models.B3{
		DataReferencia:              dataReferencia,
		AcaoAtualizacao:             acaoAtualizacao,
		DataNegocio:                 dataNegocio,
		CodigoInstrumento:           row[3],
		PrecoNegocio:                precoNegocio,
		QuantidadeNegociada:         quantidadeNegociada,
		HoraFechamento:              horaFechamento,
		CodigoIdentificadorNegocio:  row[7],
		TipoSessaoPregao:            tipoSessaoPregao,
		CodigoParticipanteComprador: row[9],
		CodigoParticipanteVendedor:  row[10],
	}, nil
}
