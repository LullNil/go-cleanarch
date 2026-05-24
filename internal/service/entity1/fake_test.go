package entity1

import (
	"context"

	"github.com/LullNil/go-cleanarch/domain"
	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
)

type fakeRepo struct {
	nextID int64
	items  map[int64]*domainentity1.Entity1

	saveCalls   int
	getCalls    int
	updateCalls int
	deleteCalls int
}

func newFakeRepo(items ...*domainentity1.Entity1) *fakeRepo {
	repo := &fakeRepo{
		nextID: 1,
		items:  make(map[int64]*domainentity1.Entity1),
	}

	for _, item := range items {
		repo.items[item.ID] = cloneEntity(item)
		if item.ID >= repo.nextID {
			repo.nextID = item.ID + 1
		}
	}

	return repo
}

func (r *fakeRepo) Save(_ context.Context, e *domainentity1.Entity1) (int64, error) {
	r.saveCalls++

	id := r.nextID
	r.nextID++

	clone := cloneEntity(e)
	clone.ID = id
	r.items[id] = clone

	return id, nil
}

func (r *fakeRepo) GetByID(_ context.Context, id int64) (*domainentity1.Entity1, error) {
	r.getCalls++

	e, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return cloneEntity(e), nil
}

func (r *fakeRepo) Update(_ context.Context, e *domainentity1.Entity1) error {
	r.updateCalls++

	if _, ok := r.items[e.ID]; !ok {
		return domain.ErrNotFound
	}

	r.items[e.ID] = cloneEntity(e)
	return nil
}

func (r *fakeRepo) Delete(_ context.Context, id int64) error {
	r.deleteCalls++

	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}

	delete(r.items, id)
	return nil
}

type fakeCache struct {
	items map[int64]*domainentity1.Entity1

	getCalls    int
	setCalls    int
	deleteCalls int
}

func newFakeCache(items ...*domainentity1.Entity1) *fakeCache {
	cache := &fakeCache{
		items: make(map[int64]*domainentity1.Entity1),
	}

	for _, item := range items {
		cache.items[item.ID] = cloneEntity(item)
	}

	return cache
}

func (c *fakeCache) Get(_ context.Context, id int64) (*domainentity1.Entity1, error) {
	c.getCalls++

	e, ok := c.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return cloneEntity(e), nil
}

func (c *fakeCache) Set(_ context.Context, e *domainentity1.Entity1) error {
	c.setCalls++
	c.items[e.ID] = cloneEntity(e)
	return nil
}

func (c *fakeCache) Delete(_ context.Context, id int64) error {
	c.deleteCalls++
	delete(c.items, id)
	return nil
}

func cloneEntity(e *domainentity1.Entity1) *domainentity1.Entity1 {
	if e == nil {
		return nil
	}

	clone := *e
	return &clone
}
