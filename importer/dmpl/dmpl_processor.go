package dmpl

import (
	"b3importer/data"
	"b3importer/importer/common"
)

type DemEquityChangeProcessor struct {
	consolidated bool
}

func NewDemEquityChangeProcessor(consolidated bool) *DemEquityChangeProcessor {
	return &DemEquityChangeProcessor{
		consolidated,
	}
}

func (p *DemEquityChangeProcessor) PerformMigration() error {
	return data.Instance().Migrate(
		&DemEquityChangeConsolidated{},
		&DemEquityChangeIndividual{},
	)
}

func (p *DemEquityChangeProcessor) ProcessContent(content [][]string) error {
	var values interface{}
	var conv_error error
	if p.consolidated {
		values, conv_error = getDMPLConsolidatedFromData(content)
		if conv_error != nil {
			return conv_error
		}
	} else {
		values, conv_error = getDMPLIndividualFromData(content)
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

func (p *DemEquityChangeProcessor) clearOldRecords(content [][]string) error {

	var obj interface{}
	if p.consolidated {
		obj = DemEquityChangeConsolidated{}
	} else {
		obj = DemEquityChangeIndividual{}
	}

	return common.ClearOldRecords(content, obj)
}
