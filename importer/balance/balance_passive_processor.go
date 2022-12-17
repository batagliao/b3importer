package balance

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type BalancePassiveProcessor struct {
	consolidated bool
}

func NewBalancePassiveProcessor(consolidated bool) *BalancePassiveProcessor {
	return &BalancePassiveProcessor{
		consolidated,
	}
}

func (p *BalancePassiveProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&BalancePassiveIndividual{},
		&BalancePassiveConsolidated{},
	)
}

func (p *BalancePassiveProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getPassiveConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getPassiveIndividualFromData(content)
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

func (p *BalancePassiveProcessor) clearOldRecords(content [][]string) error {
	var obj interface{}
	if p.consolidated {
		obj = BalancePassiveConsolidated{}
	} else {
		obj = BalancePassiveIndividual{}
	}

	return common.ClearOldRecords(content, obj)
}
