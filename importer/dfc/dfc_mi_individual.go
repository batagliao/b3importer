package dfc

import (
	"b3importer/importer/common"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DemCashFlowIndirectMethodIndividual struct {
	gorm.Model
	CNPJCia          string    `gorm:"index"`
	DataReferencia   time.Time `gorm:"index"`
	Versao           int
	DenominacaoCia   string
	CodigoCVM        string
	GrupoDFP         string
	Moeda            string
	EscalaMoeda      string
	OrdemExercicio   string
	DataIniExercicio time.Time
	DataFimExercicio time.Time
	CodigoConta      string `gorm:"index"`
	DescricaoConta   string `gorm:"index"`
	ValorConta       float64
	ContaFixa        bool
}

func (DemCashFlowIndirectMethodIndividual) TableName() string {
	return "DemFluxoCaixaMetodoIndiretoIndividual"
}

func getDFCIndirectIndividualFromData(data [][]string) ([]*DemCashFlowIndirectMethodIndividual, error) {
	total := len(data)
	result := []*DemCashFlowIndirectMethodIndividual{}
	for i, rec := range data {
		log.Info().Msgf("reading record %d/%d", i+1, total)
		value, err := parseDFCIndirectIndividualRecord(rec)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}

func parseDFCIndirectIndividualRecord(rec []string) (*DemCashFlowIndirectMethodIndividual, error) {
	obj := DemCashFlowIndirectMethodIndividual{}
	obj.CNPJCia = strings.TrimSpace(rec[0])

	dtRef, err := time.Parse(common.DATE_LAYOUT, rec[1])
	if err != nil {
		return nil, err
	}
	obj.DataReferencia = dtRef

	versao, err := strconv.Atoi(rec[2])
	if err != nil {
		return nil, err
	}
	obj.Versao = versao

	obj.DenominacaoCia = strings.TrimSpace(rec[3])
	obj.CodigoCVM = strings.TrimSpace(rec[4])
	obj.GrupoDFP = strings.TrimSpace(rec[5])
	obj.Moeda = strings.TrimSpace(rec[6])
	obj.EscalaMoeda = strings.TrimSpace(rec[7])
	obj.OrdemExercicio = strings.TrimSpace(rec[8])

	dtIni, err := time.Parse(common.DATE_LAYOUT, rec[9])
	if err != nil {
		return nil, err
	}
	obj.DataIniExercicio = dtIni

	dtFim, err := time.Parse(common.DATE_LAYOUT, rec[10])
	if err != nil {
		return nil, err
	}
	obj.DataFimExercicio = dtFim
	obj.CodigoConta = strings.TrimSpace(rec[11])
	obj.DescricaoConta = strings.TrimSpace(rec[12])

	valor, err := strconv.ParseFloat(strings.TrimSpace(rec[13]), 64)
	if err != nil {
		return nil, err
	}
	obj.ValorConta = valor

	return &obj, nil
}
