package dre

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type DemEarningsProcessor struct {
	consolidated bool
}

func NewWideEarningsProcessor(consolidated bool) *DemEarningsProcessor {
	return &DemEarningsProcessor{
		consolidated,
	}
}

func (p *DemEarningsProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&DemEarningsConsolidated{},
		&DemEarningsIndividual{},
	)
}

func (p *DemEarningsProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getDREConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getDREIndividualFromData(content)
		if conv_error != nil {
			return conv_error
		}
	}

	// get pairs to delete (cnpj, date)
	if err := p.clearOldRecords(content); err != nil {
		return err
	}

	// inserto to db
	if err := data.Instance().Create(values); err != nil {
		return err
	}

	return nil
}

func (p *DemEarningsProcessor) clearOldRecords(content [][]string) error {

	var obj interface{}
	if p.consolidated {
		obj = DemEarningsConsolidated{}
	} else {
		obj = DemEarningsIndividual{}
	}

	return common.ClearOldRecords(content, obj)

}
