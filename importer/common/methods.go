package common

import (
	"b3importer/data"
	"time"
)

func ClearOldRecords(content [][]string, obj interface{}) error {
	// structure to contain data already seen with O(1) access
	var pairs map[string]map[time.Time]bool = make(map[string]map[time.Time]bool)

	for _, rec := range content {
		cnpj := rec[0]
		data_ref, err := time.Parse(DATE_LAYOUT, rec[1])
		if err != nil {
			return err
		}

		// check if we have cnpj already
		first_step_pair, exist_cnpj := pairs[cnpj]
		if !exist_cnpj {
			err := data.Instance().ClearCompanyOldData(obj, cnpj, data_ref)
			if err != nil {
				return err
			}
			pairs[cnpj] = make(map[time.Time]bool)
			pairs[cnpj][data_ref] = true
			continue
		}

		// here cnpj exists
		_, exists_data_ref := first_step_pair[data_ref]
		if !exists_data_ref {
			err := data.Instance().ClearCompanyOldData(obj, cnpj, data_ref)
			if err != nil {
				return err
			}
			pairs[cnpj][data_ref] = true
			continue
		}
	}
	return nil
}
