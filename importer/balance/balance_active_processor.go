package balance

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type BalanceActiveProcessor struct {
	consolidated bool
}

func NewBalanceActiveProcessor(consolidated bool) *BalanceActiveProcessor {
	return &BalanceActiveProcessor{
		consolidated,
	}
}

func (p *BalanceActiveProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&BalanceActiveIndividual{},
		&BalanceActiveConsolidated{},
	)
}

func (p *BalanceActiveProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getActiveConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getActiveIndividualFromData(content)
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

func (p *BalanceActiveProcessor) clearOldRecords(content [][]string) error {
	var obj interface{}
	if p.consolidated {
		obj = BalanceActiveConsolidated{}
	} else {
		obj = BalanceActiveIndividual{}
	}

	return common.ClearOldRecords(content, obj)
}
