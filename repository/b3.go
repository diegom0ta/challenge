package repository

import (
	"database/sql"
	"fmt"
	"time"

	"challenge/models"
)

type B3Repository struct {
	db *sql.DB
}

func NewB3Repository(db *sql.DB) *B3Repository {
	return &B3Repository{db: db}
}

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

func (r *B3Repository) InsertBatch(b3Records []*models.B3) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

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

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Printf("Successfully inserted %d B3 records\n", len(b3Records))
	return nil
}

func (r *B3Repository) GetAggregatedData(ticker string, startDate *time.Time) (float64, int, error) {
	var query string
	var args []interface{}

	if startDate != nil {
		startDateStr := startDate.Format("2006-01-02")
		query = `
			SELECT 
				MAX(preco_negocio) as max_price,
				MAX(daily_volume) as max_daily_volume
			FROM (
				SELECT 
					preco_negocio,
					SUM(quantidade_negociada) OVER (PARTITION BY data_negocio) as daily_volume
				FROM b3 
				WHERE codigo_instrumento = $1 
				AND data_negocio >= $2
			) subq
		`
		args = []interface{}{ticker, startDateStr}
	} else {
		today := time.Now()
		sevenDaysAgo := today.AddDate(0, 0, -7).Format("2006-01-02")
		yesterday := today.AddDate(0, 0, -1).Format("2006-01-02")

		query = `
			SELECT 
				MAX(preco_negocio) as max_price,
				MAX(daily_volume) as max_daily_volume
			FROM (
				SELECT 
					preco_negocio,
					SUM(quantidade_negociada) OVER (PARTITION BY data_negocio) as daily_volume
				FROM b3 
				WHERE codigo_instrumento = $1 
				AND data_negocio >= $2
				AND data_negocio <= $3
			) subq
		`
		args = []interface{}{ticker, sevenDaysAgo, yesterday}
	}

	var maxPrice sql.NullFloat64
	var maxDailyVolume sql.NullInt64

	err := r.db.QueryRow(query, args...).Scan(&maxPrice, &maxDailyVolume)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get aggregated data: %v", err)
	}

	var price float64
	var volume int

	if maxPrice.Valid {
		price = maxPrice.Float64
	}

	if maxDailyVolume.Valid {
		volume = int(maxDailyVolume.Int64)
	}

	return price, volume, nil
}

func ParseCSVRowToB3(row []string) (*models.B3, error) {
	if len(row) != 11 {
		return nil, fmt.Errorf("expected exactly 11 columns, got %d", len(row))
	}

	dataReferencia, err := time.Parse("2006-01-02", row[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse data_referencia: %v", err)
	}

	acaoAtualizacao := row[1]

	dataNegocio := row[2]

	var precoNegocio float64
	if _, err := fmt.Sscanf(row[4], "%f", &precoNegocio); err != nil {
		return nil, fmt.Errorf("failed to parse preco_negocio: %v", err)
	}

	var quantidadeNegociada int
	if _, err := fmt.Sscanf(row[5], "%d", &quantidadeNegociada); err != nil {
		return nil, fmt.Errorf("failed to parse quantidade_negociada: %v", err)
	}

	var horaFechamento int
	if _, err := fmt.Sscanf(row[6], "%d", &horaFechamento); err != nil {
		return nil, fmt.Errorf("failed to parse hora_fechamento: %v", err)
	}

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
