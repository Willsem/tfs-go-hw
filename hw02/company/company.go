package company

import "github.com/willsem/tfs-go-hw/hw02/company/operation"

type Operation struct {
	Company        *string              `json:"company,omitempty"`
	OperationField *operation.Operation `json:"operation,omitempty"`
	operation.Operation
}

func (op Operation) GetType() *operation.Type {
	if op.Type != nil {
		return op.Type
	}

	if op.OperationField != nil && op.OperationField.Type != nil {
		return op.OperationField.Type
	}

	return nil
}

func (op Operation) GetValue() *operation.Value {
	if op.Value != nil {
		return op.Value
	}

	if op.OperationField != nil && op.OperationField.Value != nil {
		return op.OperationField.Value
	}

	return nil
}

func (op Operation) GetID() *operation.ID {
	if op.ID != nil {
		return op.ID
	}

	if op.OperationField != nil && op.OperationField.ID != nil {
		return op.OperationField.ID
	}

	return nil
}

func (op Operation) GetCreatedAt() *operation.Datetime {
	if op.CreatedAt != nil {
		return op.CreatedAt
	}

	if op.OperationField != nil && op.OperationField.CreatedAt != nil {
		return op.OperationField.CreatedAt
	}

	return nil
}

func (op Operation) CheckValid() bool {
	t := op.GetType()
	v := op.GetValue()
	id := op.GetID()
	c := op.GetCreatedAt()

	if t == nil || v == nil || id == nil || c == nil {
		return false
	}

	if t.Err != nil || v.Err != nil || id.Err != nil || c.Err != nil {
		return false
	}

	return true
}
