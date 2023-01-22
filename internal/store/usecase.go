package store

import (
	"context"
	"time"

	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/helpers/response"
	"github.com/Risuii/models/store"
)

type (
	StoreUseCase interface {
		CreateStore(ctx context.Context, userid int64, params store.Store) response.Response
		Read(ctx context.Context, id int64) response.Response
		UpdateStore(ctx context.Context, id int64, params store.Store) response.Response
		DeleteStore(ctx context.Context, id int64) response.Response
	}

	storeUseCaseimpl struct {
		repository StoreRepository
	}
)

func NewStoreUseCaseImpl(repo StoreRepository) StoreUseCase {
	return &storeUseCaseimpl{
		repository: repo,
	}
}

func (su *storeUseCaseimpl) CreateStore(ctx context.Context, userid int64, params store.Store) response.Response {
	store := store.Store{
		UserID:      userid,
		NameStore:   params.NameStore,
		Description: params.Description,
		CreatedAt:   time.Now(),
	}

	ID, err := su.repository.Create(ctx, store)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	store.ID = ID

	return response.Success(response.StatusCreated, store)
}

func (su *storeUseCaseimpl) Read(ctx context.Context, id int64) response.Response {
	return nil
}

func (su *storeUseCaseimpl) UpdateStore(ctx context.Context, id int64, params store.Store) response.Response {
	return nil
}

func (su *storeUseCaseimpl) DeleteStore(ctx context.Context, id int64) response.Response {
	return nil
}
