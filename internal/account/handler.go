package account

import (
	"encoding/json"
	"net/http"

	newJWT "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Risuii/config/jwt"
	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/helpers/response"
	"github.com/Risuii/models/account"
)

type AccountHandler struct {
	Validate *validator.Validate
	UseCase  AccountUseCase
}

func NewAbsensiHandler(router *mux.Router, validate *validator.Validate, usecase AccountUseCase) {
	handler := &AccountHandler{
		Validate: validate,
		UseCase:  usecase,
	}

	api := router.PathPrefix("/account").Subrouter()

	router.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	api.HandleFunc("/update", handler.Update).Methods(http.MethodPatch)
}

func (handler *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput account.Account

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		res = response.Error(response.StatusUnprocessableEntity, exception.ErrUnprocessableEntity)
		res.JSON(w)
		return
	}

	err := handler.Validate.StructCtx(ctx, userInput)
	if err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res = handler.UseCase.Register(ctx, userInput)

	res.JSON(w)
}

func (handler *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput account.AccountLogin

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		res = response.Error(response.StatusUnprocessableEntity, exception.ErrUnprocessableEntity)
		res.JSON(w)
		return
	}

	err := handler.Validate.StructCtx(ctx, userInput)
	if err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res, token := handler.UseCase.Login(ctx, userInput)

	if token.Token == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "",
			Path:     "",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token.Token,
		HttpOnly: true,
	})

	res.JSON(w)
}

func (handler *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	var userInput account.Account

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

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		res = response.Error(response.StatusUnprocessableEntity, exception.ErrUnprocessableEntity)
		res.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, userInput)
	if err != nil {
		res = response.Error(response.StatusBadRequest, exception.ErrBadRequest)
		res.JSON(w)
		return
	}

	res = handler.UseCase.Update(ctx, claims.ID, userInput)

	res.JSON(w)
}
