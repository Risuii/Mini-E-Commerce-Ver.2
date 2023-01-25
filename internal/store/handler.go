package store

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	api.HandleFunc("/store/{id}", handler.EditStore).Methods(http.MethodPatch)
	api.HandleFunc("/store/{id}", handler.DeleteStore).Methods(http.MethodDelete)

	router.HandleFunc("/store/{userID}", handler.Store).Methods(http.MethodGet)
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
	var res response.Response

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

	res, token := handler.UseCase.Read(ctx, claims.UserID)

	http.SetCookie(w, &http.Cookie{
		Name:     "Store-token",
		Path:     "/",
		Value:    token.Token,
		HttpOnly: true,
	})

	res.JSON(w)
}

func (handler *StoreHandler) Store(w http.ResponseWriter, r *http.Request) {
	var res response.Response

	ctx := r.Context()

	params := mux.Vars(r)
	userID, _ := strconv.ParseInt(params["userID"], 10, 64)

	res, _ = handler.UseCase.Read(ctx, userID)

	res.JSON(w)
}

func (handler *StoreHandler) EditStore(w http.ResponseWriter, r *http.Request) {
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

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)

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

	userInput.UserID = claims.UserID

	res = handler.UseCase.UpdateStore(ctx, id, userInput)

	res.JSON(w)
}

func (handler *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	var res response.Response

	ctx := r.Context()

	_, err := r.Cookie("token")
	if err != nil {
		res = response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
		res.JSON(w)
		return
	}

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)

	res = handler.UseCase.DeleteStore(ctx, id)

	res.JSON(w)
}
