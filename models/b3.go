package models

import "time"

type B3 struct {
	DataNegocio         time.Time
	CodigoInstrumento   string
	PrecoNegocio        float64
	QuantidadeNegociada int
	HoraFechamento      time.Time
}
