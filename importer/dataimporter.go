package importer

import (
	"b3importer/importer/balance"
	"b3importer/importer/common"
	"b3importer/importer/dfc"
	"b3importer/importer/dmpl"
	"b3importer/importer/dra"
	"b3importer/importer/dre"
	"b3importer/importer/dva"
	"encoding/csv"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

type DataImporter struct {
}

func NewDataImporter() *DataImporter {
	return &DataImporter{}
}

func (me *DataImporter) ImportFile(file string) error {
	// read file

	fileReader, err := os.OpenFile(file, os.O_RDONLY, os.ModeType)
	if err != nil {
		return err
	}

	defer func() {
		if err := fileReader.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := charmap.Windows1252.NewDecoder()
	csvReader := csv.NewReader(decoder.Reader(fileReader))
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	content, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	content = content[1:] //remove csv first row

	// detect document type
	processor := getContentProcessor(file)
	if processor != nil {
		err := processor.PerformMigration()
		if err != nil {
			return err
		}
		log.Info().Msgf("processing file %s", file)
		return processor.ProcessContent(content)
	}

	return nil
}

func getContentProcessor(filename string) common.ContentProcessor {
	consolidated := strings.Contains(filename, common.CONSOLIDATED)
	switch {
	case strings.Contains(filename, string(common.BPA)):
		return balance.NewBalanceActiveProcessor(consolidated)
	case strings.Contains(filename, string(common.BPP)):
		return balance.NewBalancePassiveProcessor(consolidated)
	case strings.Contains(filename, string(common.DFCMD)):
		return dfc.NewDemCashFlowDirectProcessor(consolidated)
	case strings.Contains(filename, string(common.DFCMI)):
		return dfc.NewDemCashFlowIndirectProcessor(consolidated)
	case strings.Contains(filename, string(common.DMPL)):
		return dmpl.NewDemEquityChangeProcessor(consolidated)
	case strings.Contains(filename, string(common.DRA)):
		return dra.NewWideEarningsProcessor(consolidated)
	case strings.Contains(filename, string(common.DRE)):
		return dre.NewWideEarningsProcessor(consolidated)
	case strings.Contains(filename, string(common.DVA)):
		return dva.NewDemAddedValueProcessor(consolidated)
	default:
		break
	}
	return nil
}
