package mocks

import (
	"context"

	"go-starter/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type Querier struct {
	mock.Mock
	models.Querier
}

func (m *Querier) CreatePost(ctx context.Context, db models.DBTX, params models.CreatePostParams) (models.Post, error) {
	args := m.Called(ctx, db, params)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *Querier) ListPosts(ctx context.Context, db models.DBTX) ([]models.Post, error) {
	args := m.Called(ctx, db)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *Querier) GetPost(ctx context.Context, db models.DBTX, id uuid.UUID) (models.Post, error) {
	args := m.Called(ctx, db, id)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *Querier) UpdatePost(ctx context.Context, db models.DBTX, params models.UpdatePostParams) (models.Post, error) {
	args := m.Called(ctx, db, params)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *Querier) DeletePost(ctx context.Context, db models.DBTX, id uuid.UUID) (int64, error) {
	args := m.Called(ctx, db, id)
	return args.Get(0).(int64), args.Error(1)
}
