package models

import "time"

type B3 struct {
	DataReferencia              time.Time
	AcaoAtualizacao             string
	DataNegocio                 string
	CodigoInstrumento           string
	PrecoNegocio                float64
	QuantidadeNegociada         int
	HoraFechamento              int
	CodigoIdentificadorNegocio  string
	TipoSessaoPregao            int
	CodigoParticipanteComprador string
	CodigoParticipanteVendedor  string
}
