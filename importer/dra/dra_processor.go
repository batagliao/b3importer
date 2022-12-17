package dra

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type DemWideEarningsProcessor struct {
	consolidated bool
}

func NewWideEarningsProcessor(consolidated bool) *DemWideEarningsProcessor {
	return &DemWideEarningsProcessor{
		consolidated,
	}
}

func (p *DemWideEarningsProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&DemWideEarningsConsolidated{},
		&DemWideEarningsIndividual{},
	)
}

func (p *DemWideEarningsProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getDRAConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getDRAIndividualFromData(content)
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

func (p *DemWideEarningsProcessor) clearOldRecords(content [][]string) error {

	var obj interface{}
	if p.consolidated {
		obj = DemWideEarningsConsolidated{}
	} else {
		obj = DemWideEarningsIndividual{}
	}

	return common.ClearOldRecords(content, obj)

}
