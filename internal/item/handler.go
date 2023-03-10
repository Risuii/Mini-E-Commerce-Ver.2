package item

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
	"github.com/Risuii/models/item"
)

type ItemHandler struct {
	validate *validator.Validate
	UseCase  ItemUseCase
}

func NewItemHandler(router *mux.Router, validate *validator.Validate, usecase ItemUseCase) {
	handler := ItemHandler{
		validate: validate,
		UseCase:  usecase,
	}

	api := router.PathPrefix("/store").Subrouter()

	api.HandleFunc("/items", handler.AddItem).Methods(http.MethodPost)
	api.HandleFunc("/items", handler.GetAllItems).Methods(http.MethodPut)
	api.HandleFunc("/items/{itemID}", handler.GetOneItem).Methods(http.MethodGet)
	api.HandleFunc("/items/{id}", handler.UpdateItem).Methods(http.MethodPatch)
	api.HandleFunc("/items/{id}", handler.DeleteItem).Methods(http.MethodDelete)
}

func (handler *ItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput item.Item

	ctx := r.Context()

	c, err := r.Cookie("Store-token")
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

	if err := handler.validate.StructCtx(ctx, userInput); err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res = handler.UseCase.AddItem(ctx, claims.StoreID, userInput)

	res.JSON(w)
}

func (handler *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	var res response.Response

	ctx := r.Context()

	c, err := r.Cookie("Store-token")
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

	res = handler.UseCase.GetAllItems(ctx, claims.StoreID)

	res.JSON(w)
}

func (handler *ItemHandler) GetOneItem(w http.ResponseWriter, r *http.Request) {
	var res response.Response

	ctx := r.Context()

	c, err := r.Cookie("Store-token")
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
	itemID, _ := strconv.ParseInt(params["itemID"], 10, 64)

	res = handler.UseCase.GetOneItem(ctx, itemID, claims.StoreID)

	res.JSON(w)
}

func (handler *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput item.Item

	ctx := r.Context()

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)

	c, err := r.Cookie("Store-token")
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

	if err := handler.validate.StructCtx(ctx, userInput); err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res = handler.UseCase.UpdateItem(ctx, id, userInput)

	res.JSON(w)
}

func (handler *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var res response.Response

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)

	ctx := r.Context()

	res = handler.UseCase.DeleteItem(ctx, id)

	res.JSON(w)
}
