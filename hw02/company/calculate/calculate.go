package calculate

import "github.com/willsem/tfs-go-hw/hw02/company"

type OperationsResult struct {
	Company              string               `json:"company"`
	ValidOperationsCount int                  `json:"valid_operations_count"`
	Balance              int                  `json:"balance"`
	InvalidOperations    []InvalidOperationID `json:"invalid_operations"`
}

func CalculateOperations(operations []company.Operation) []OperationsResult {
	return []OperationsResult{}
}
