package calculate

import (
	"github.com/willsem/tfs-go-hw/hw02/company"
)

type OperationsResult struct {
	Company              string               `json:"company"`
	ValidOperationsCount int                  `json:"valid_operations_count"`
	Balance              int                  `json:"balance"`
	InvalidOperations    []InvalidOperationID `json:"invalid_operations,omitempty"`
}

func Operations(operations []company.Operation) []OperationsResult {
	resultMap := make(map[string]OperationsResult)
	for _, operation := range operations {
		name := *operation.Company
		c, ok := resultMap[name]
		if !ok {
			c = OperationsResult{
				Company:              name,
				ValidOperationsCount: 0,
				Balance:              0,
				InvalidOperations:    make([]InvalidOperationID, 0),
			}
		}

		if operation.CheckValid() {
			c.ValidOperationsCount++
			c.Balance += operation.GetValue().Value
		} else if operation.GetID() != nil {
			c.InvalidOperations = append(c.InvalidOperations,
				NewInvalidOperationID(operation.GetID().Value))
		}

		resultMap[name] = c
	}

	result := make([]OperationsResult, 0, len(resultMap))
	for _, value := range resultMap {
		if len(value.InvalidOperations) == 0 {
			value.InvalidOperations = nil
		}
		result = append(result, value)
	}

	return result
}
