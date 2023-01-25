package store

import (
	"context"
	"time"

	newJWT "github.com/dgrijalva/jwt-go"

	"github.com/Risuii/config/jwt"
	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/helpers/response"
	"github.com/Risuii/models/store"
	"github.com/Risuii/models/token"
)

type (
	StoreUseCase interface {
		CreateStore(ctx context.Context, userid int64, params store.Store) response.Response
		Read(ctx context.Context, userID int64) (response.Response, token.Token)
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

	_, err := su.repository.FindByName(ctx, params.NameStore)
	if err == nil {
		return response.Error(response.StatusConflicted, exception.ErrConflicted)
	}

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

func (su *storeUseCaseimpl) Read(ctx context.Context, userID int64) (response.Response, token.Token) {

	store, err := su.repository.FindByUserID(ctx, userID)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound), token.Token{}
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer), token.Token{}
	}

	claims := &jwt.JWTclaim{
		StoreID: store[0].ID,
		StandardClaims: newJWT.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
		},
	}

	tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, claims)

	tokenString, err := tokenAlgo.SignedString(jwt.JWT_KEY)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer), token.Token{}
	}

	newToken := token.Token{
		Token: tokenString,
	}

	return response.Success(response.StatusOK, store), newToken
}

func (su *storeUseCaseimpl) UpdateStore(ctx context.Context, id int64, params store.Store) response.Response {
	stores, err := su.repository.FindByID(ctx, id)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	stores = store.Store{
		UserID:      params.UserID,
		NameStore:   params.NameStore,
		Description: params.Description,
		UpdateAt:    time.Now(),
	}

	err = su.repository.Update(ctx, id, stores)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, stores)
}

func (su *storeUseCaseimpl) DeleteStore(ctx context.Context, id int64) response.Response {
	_, err := su.repository.FindByID(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	err = su.repository.Delete(ctx, id)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	msg := "Success Delete Store"

	return response.Success(response.StatusOK, msg)
}
