package dfc

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type DemCashFlowIndirectProcessor struct {
	consolidated bool
}

func NewDemCashFlowIndirectProcessor(consolidated bool) *DemCashFlowIndirectProcessor {
	return &DemCashFlowIndirectProcessor{
		consolidated,
	}
}

func (p *DemCashFlowIndirectProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&DemCashFlowIndirectMethodIndividual{},
		&DemCashFlowIndirectMethodConsolidated{},
	)
}

func (p *DemCashFlowIndirectProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getDFCIndirectConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getDFCIndirectIndividualFromData(content)
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

func (p *DemCashFlowIndirectProcessor) clearOldRecords(content [][]string) error {
	var obj interface{}
	if p.consolidated {
		obj = DemCashFlowIndirectMethodConsolidated{}
	} else {
		obj = DemCashFlowIndirectMethodIndividual{}
	}
	return common.ClearOldRecords(content, obj)
}
