package item

import (
	"context"
	"time"

	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/helpers/response"
	"github.com/Risuii/models/item"
)

type (
	ItemUseCase interface {
		AddItem(ctx context.Context, storeID int64, params item.Item) response.Response
		GetAllItem() response.Response
		GetOneItem(ctx context.Context, id int64) response.Response
		UpdateItem(ctx context.Context, id int64, params item.Item) response.Response
		DeleteItem(ctx context.Context, id int64) response.Response
	}

	itemUseCaseImpl struct {
		repository ItemRepository
	}
)

func NewItemUseCaseImpl(repo ItemRepository) ItemUseCase {
	return &itemUseCaseImpl{
		repository: repo,
	}
}

func (iu *itemUseCaseImpl) AddItem(ctx context.Context, storeID int64, params item.Item) response.Response {
	_, err := iu.repository.FindByName(ctx, params.Name)

	if err == nil {
		return response.Error(response.StatusConflicted, exception.ErrConflicted)
	}

	item := item.Item{
		StoreID:     storeID,
		Name:        params.Name,
		Description: params.Description,
		Quantity:    params.Quantity,
		CreatedAt:   time.Now(),
	}

	ID, err := iu.repository.AddItem(ctx, item)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	item.ID = ID

	return response.Success(response.StatusCreated, item)
}

func (iu *itemUseCaseImpl) GetAllItem() response.Response {
	return nil
}

func (iu *itemUseCaseImpl) GetOneItem(ctx context.Context, id int64) response.Response {
	return nil
}

func (iu *itemUseCaseImpl) UpdateItem(ctx context.Context, id int64, params item.Item) response.Response {
	return nil
}

func (iu *itemUseCaseImpl) DeleteItem(ctx context.Context, id int64) response.Response {
	return nil
}
