package entity1

type CreateRequest struct {
	Field1 bool
	Field2 int64
	Field3 string
}

type UpdateRequest struct {
	ID     int64
	Field3 string
}
