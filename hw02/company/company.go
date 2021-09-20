package company

import "github.com/willsem/tfs-go-hw/hw02/company/operation"

type Operation struct {
	Company        *string              `json:"company,omitempty"`
	OperationField *operation.Operation `json:"operation,omitempty"`
	operation.Operation
}
