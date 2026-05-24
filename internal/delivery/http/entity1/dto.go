package entity1

import (
	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

type createRequest struct {
	Field1 bool   `json:"field1"`
	Field2 int64  `json:"field2"`
	Field3 string `json:"field3"`
}

type createResponse struct {
	ID int64 `json:"id"`
}

func (r *createRequest) toServiceRequest() *entity1service.CreateRequest {
	return &entity1service.CreateRequest{
		Field1: r.Field1,
		Field2: r.Field2,
		Field3: r.Field3,
	}
}

type updateRequest struct {
	Field3 string `json:"field3"`
}

type updateResponse struct {
	Status string `json:"status"`
}

func (r *updateRequest) toServiceRequest(id int64) *entity1service.UpdateRequest {
	return &entity1service.UpdateRequest{
		ID:     id,
		Field3: r.Field3,
	}
}

type getResponse struct {
	ID     int64  `json:"id"`
	Field1 bool   `json:"field1"`
	Field2 int64  `json:"field2"`
	Field3 string `json:"field3"`
}

func toGetResponse(e *domainentity1.Entity1) *getResponse {
	return &getResponse{
		ID:     e.ID,
		Field1: e.Field1,
		Field2: e.Field2,
		Field3: e.Field3,
	}
}

type deleteResponse struct {
	Status string `json:"status"`
}
