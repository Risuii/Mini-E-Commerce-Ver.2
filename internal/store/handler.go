package store

import (
	"encoding/json"
	"net/http"

	newJWT "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Risuii/config/jwt"
	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/helpers/response"
	"github.com/Risuii/models/store"
)

type StoreHandler struct {
	validate *validator.Validate
	UseCase  StoreUseCase
}

func NewStoreHandler(router *mux.Router, validate *validator.Validate, usecase StoreUseCase) {
	handler := &StoreHandler{
		validate: validate,
		UseCase:  usecase,
	}

	api := router.PathPrefix("/account").Subrouter()

	api.HandleFunc("/store", handler.CreateStore).Methods(http.MethodPost)
	api.HandleFunc("/store", handler.GetStore).Methods(http.MethodGet)
	api.HandleFunc("/store", handler.EditStore).Methods(http.MethodPatch)
	api.HandleFunc("/store", handler.DeleteStore).Methods(http.MethodDelete)
}

func (handler *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput store.Store

	ctx := r.Context()

	c, err := r.Cookie("token")
	if err != nil {
		res = response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
		res.JSON(w)
		return
	}

	tokenString := c.Value
	claims := &jwt.JWTclaim{}

	newJWT.ParseWithClaims(tokenString, claims, func(t *newJWT.Token) (interface{}, error) {
		return jwt.JWT_KEY, nil
	})

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		res = response.Error(response.StatusUnprocessableEntity, exception.ErrUnprocessableEntity)
		res.JSON(w)
		return
	}

	err = handler.validate.StructCtx(ctx, userInput)
	if err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res = handler.UseCase.CreateStore(ctx, claims.UserID, userInput)

	res.JSON(w)
}

func (handler *StoreHandler) GetStore(w http.ResponseWriter, r *http.Request) {

}

func (handler *StoreHandler) EditStore(w http.ResponseWriter, r *http.Request) {

}

func (handler *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {

}
