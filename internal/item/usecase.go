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
		GetAllItems(ctx context.Context, storeID int64) response.Response
		GetOneItem(ctx context.Context, id int64, storeID int64) response.Response
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
	data, err := iu.repository.FindByName(ctx, params.Name)

	if err == nil {
		data = item.Item{
			ID:          data.ID,
			StoreID:     data.ID,
			Name:        data.Name,
			Description: data.Description,
			Quantity:    data.Quantity + params.Quantity,
			UpdateAt:    data.UpdateAt,
		}

		err := iu.repository.UpdateKuantitas(ctx, data.ID, data)
		if err != nil {
			return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
		}

		return response.Success(response.StatusOK, data)
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

func (iu *itemUseCaseImpl) GetAllItems(ctx context.Context, storeID int64) response.Response {

	data, err := iu.repository.GetAllItem(ctx, storeID)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, data)
}

func (iu *itemUseCaseImpl) GetOneItem(ctx context.Context, id int64, storeID int64) response.Response {

	data, err := iu.repository.FindByIDWithStoreID(ctx, id, storeID)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}
	return response.Success(response.StatusOK, data)
}

func (iu *itemUseCaseImpl) UpdateItem(ctx context.Context, id int64, params item.Item) response.Response {
	data, err := iu.repository.FindByID(ctx, id)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	data = item.Item{
		ID:          data.ID,
		StoreID:     data.StoreID,
		Name:        params.Name,
		Description: params.Description,
		Quantity:    params.Quantity,
		UpdateAt:    time.Now(),
	}

	if err := iu.repository.UpdateItem(ctx, id, data); err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, data)
}

func (iu *itemUseCaseImpl) DeleteItem(ctx context.Context, id int64) response.Response {

	data, err := iu.repository.FindByID(ctx, id)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	if err := iu.repository.DeleteItem(ctx, data.ID); err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	msg := "Success Delete Data"

	return response.Success(response.StatusOK, msg)
}
