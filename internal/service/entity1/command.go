package entity1

// CreateCommand contains data for creating entity1.
type CreateCommand struct {
	Field1 bool
	Field2 int64
	Field3 string
}

// UpdateCommand contains data for updating entity1.
type UpdateCommand struct {
	ID     int64
	Field3 string
}

// AccessCommand contains data for checking entity1 access.
type AccessCommand struct {
	SubjectID string
	Entity1ID int64
}
