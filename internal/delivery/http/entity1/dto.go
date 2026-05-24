package entity1

type createRequest struct {
	Field1 bool   `json:"field1"`
	Field2 int64  `json:"field2"`
	Field3 string `json:"field3"`
}

type updateRequest struct {
	Field3 string `json:"field3"`
}
