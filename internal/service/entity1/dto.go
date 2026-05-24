package entity1

// CreateRequest contains data for creating entity1.
type CreateRequest struct {
	Field1 bool
	Field2 int64
	Field3 string
}

// UpdateRequest contains data for updating entity1.
type UpdateRequest struct {
	ID     int64
	Field3 string
}
