package dva

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type DemAddedValueProcessor struct {
	consolidated bool
}

func NewDemAddedValueProcessor(consolidated bool) *DemAddedValueProcessor {
	return &DemAddedValueProcessor{
		consolidated,
	}
}

func (p *DemAddedValueProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&DemAddedValueConsolidated{},
		&DemAddedValueIndividual{},
	)
}

func (p *DemAddedValueProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getDVAConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getDVAIndividualFromData(content)
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

func (p *DemAddedValueProcessor) clearOldRecords(content [][]string) error {

	var obj interface{}
	if p.consolidated {
		obj = DemAddedValueConsolidated{}
	} else {
		obj = DemAddedValueIndividual{}
	}

	return common.ClearOldRecords(content, obj)

}
