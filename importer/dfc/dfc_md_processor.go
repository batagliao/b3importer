package dfc

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type DemCashFlowDirectProcessor struct {
	consolidated bool
}

func NewDemCashFlowDirectProcessor(consolidated bool) *DemCashFlowDirectProcessor {
	return &DemCashFlowDirectProcessor{
		consolidated,
	}
}

func (p *DemCashFlowDirectProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&DemCashFlowDirectMethodIndividual{},
		&DemCashFlowDirectMethodConsolidated{},
	)
}

func (p *DemCashFlowDirectProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getDFCDirectConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getDFCDirectIndividualFromData(content)
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

func (p *DemCashFlowDirectProcessor) clearOldRecords(content [][]string) error {
	var obj interface{}
	if p.consolidated {
		obj = DemCashFlowDirectMethodConsolidated{}
	} else {
		obj = DemCashFlowDirectMethodIndividual{}
	}

	return common.ClearOldRecords(content, obj)
}
