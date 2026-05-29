package entity1

import (
	"context"
	"errors"
	"testing"

	"github.com/LullNil/go-cleanarch/domain"
	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
)

func TestCreateEntity1_InvalidInput(t *testing.T) {
	t.Parallel()

	service := New(newFakeRepo(), newFakeCache(), nil)

	_, err := service.CreateEntity1(context.Background(), &CreateCommand{
		Field3: " ",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestGetEntity1Details_CacheHitDoesNotCallRepo(t *testing.T) {
	t.Parallel()

	cached := &domainentity1.Entity1{
		ID:     10,
		Field1: true,
		Field2: 42,
		Field3: "cached",
	}
	repo := newFakeRepo()
	cache := newFakeCache(cached)
	service := New(repo, cache, nil)

	got, err := service.GetEntity1Details(context.Background(), cached.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Field3 != cached.Field3 {
		t.Fatalf("expected cached field3 %q, got %q", cached.Field3, got.Field3)
	}
	if repo.getCalls != 0 {
		t.Fatalf("expected repo not to be called, got %d calls", repo.getCalls)
	}
	if cache.getCalls != 1 {
		t.Fatalf("expected cache to be called once, got %d calls", cache.getCalls)
	}
}

func TestGetEntity1Details_CacheMissLoadsRepoAndSetsCache(t *testing.T) {
	t.Parallel()

	stored := &domainentity1.Entity1{
		ID:     20,
		Field1: true,
		Field2: 99,
		Field3: "stored",
	}
	repo := newFakeRepo(stored)
	cache := newFakeCache()
	service := New(repo, cache, nil)

	got, err := service.GetEntity1Details(context.Background(), stored.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Field3 != stored.Field3 {
		t.Fatalf("expected repo field3 %q, got %q", stored.Field3, got.Field3)
	}
	if repo.getCalls != 1 {
		t.Fatalf("expected repo to be called once, got %d calls", repo.getCalls)
	}
	if cache.setCalls != 1 {
		t.Fatalf("expected cache set once, got %d calls", cache.setCalls)
	}
	if _, ok := cache.items[stored.ID]; !ok {
		t.Fatalf("expected entity to be cached")
	}
}

func TestDeleteEntity1_DeletesRepoAndCache(t *testing.T) {
	t.Parallel()

	stored := &domainentity1.Entity1{
		ID:     30,
		Field1: true,
		Field2: 100,
		Field3: "delete-me",
	}
	repo := newFakeRepo(stored)
	cache := newFakeCache(stored)
	service := New(repo, cache, nil)

	if err := service.DeleteEntity1(context.Background(), stored.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.deleteCalls != 1 {
		t.Fatalf("expected repo delete once, got %d calls", repo.deleteCalls)
	}
	if cache.deleteCalls != 1 {
		t.Fatalf("expected cache delete once, got %d calls", cache.deleteCalls)
	}
	if _, ok := cache.items[stored.ID]; ok {
		t.Fatalf("expected cache entry to be deleted")
	}
}

func TestCheckEntity1Access_UsesAuthClient(t *testing.T) {
	t.Parallel()

	auth := &fakeAuthClient{allowed: true}
	service := New(newFakeRepo(), newFakeCache(), auth)

	err := service.CheckEntity1Access(context.Background(), &AccessCommand{
		SubjectID: "user-1",
		Entity1ID: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if auth.calls != 1 {
		t.Fatalf("expected auth client to be called once, got %d calls", auth.calls)
	}
}

func TestCheckEntity1Access_Denied(t *testing.T) {
	t.Parallel()

	service := New(newFakeRepo(), newFakeCache(), &fakeAuthClient{})

	err := service.CheckEntity1Access(context.Background(), &AccessCommand{
		SubjectID: "user-1",
		Entity1ID: 10,
	})
	if !errors.Is(err, domain.ErrPermissionDenied) {
		t.Fatalf("expected ErrPermissionDenied, got %v", err)
	}
}
