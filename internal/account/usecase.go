package account

import (
	"context"
	"time"

	newJWT "github.com/dgrijalva/jwt-go"

	"github.com/Risuii/config/bcrypt"
	"github.com/Risuii/config/jwt"
	"github.com/Risuii/helpers/exception"
	"github.com/Risuii/helpers/response"
	"github.com/Risuii/models/account"
	"github.com/Risuii/models/token"
)

type (
	AccountUseCase interface {
		Register(ctx context.Context, params account.Account) response.Response
		Login(ctx context.Context, params account.AccountLogin) (response.Response, token.Token)
		Update(ctx context.Context, id int64, params account.Account) response.Response
		ReadOne(ctx context.Context, id int64) response.Response
		Delete(ctx context.Context, id int64) response.Response
	}

	accountUseCaseImpl struct {
		repo   AccountRepository
		bcrypt bcrypt.Bcrypt
	}
)

func NewAccountUseCaseImpl(repo AccountRepository, bcrypt bcrypt.Bcrypt) AccountUseCase {
	return &accountUseCaseImpl{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (au *accountUseCaseImpl) Register(ctx context.Context, params account.Account) response.Response {
	_, err := au.repo.FindByEmail(ctx, params.Email)

	if err == nil {
		return response.Error(response.StatusConflicted, exception.ErrConflicted)
	}

	hashedPassword, err := au.bcrypt.HashPassword(params.Password)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	user := account.Account{
		ID:        params.ID,
		Name:      params.Name,
		Password:  hashedPassword,
		Email:     params.Email,
		Address:   params.Address,
		CreatedAt: time.Now(),
	}

	userID, err := au.repo.Register(ctx, user)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	user.ID = userID
	user.Password = ""

	return response.Success(response.StatusCreated, user)
}

func (au *accountUseCaseImpl) Login(ctx context.Context, params account.AccountLogin) (response.Response, token.Token) {
	user, err := au.repo.FindByEmail(ctx, params.Email)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound), token.Token{}
	}
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer), token.Token{}
	}

	isPasswordValid := au.bcrypt.ComparePasswordHash(params.Password, user.Password)

	if !isPasswordValid {
		return response.Error(response.StatusUnauthorized, exception.ErrUnauthorized), token.Token{}
	}

	user.Password = ""

	claims := &jwt.JWTclaim{
		ID:     user.ID,
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		StandardClaims: newJWT.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
		},
	}

	tokenAlgo := newJWT.NewWithClaims(newJWT.SigningMethodHS256, claims)

	tokenJWT, err := tokenAlgo.SignedString(jwt.JWT_KEY)
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer), token.Token{}
	}

	newToken := token.Token{
		Token: tokenJWT,
	}

	return response.Success(response.StatusOK, user), newToken
}

func (au *accountUseCaseImpl) Update(ctx context.Context, id int64, params account.Account) response.Response {
	user, err := au.repo.FindByID(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	hashedPassword, err := au.bcrypt.HashPassword(params.Password)
	if err != nil {
		return response.Error(response.StatusUnauthorized, exception.ErrUnauthorized)
	}

	user.Name = params.Name
	user.Password = hashedPassword
	user.Email = params.Email
	user.Address = params.Address
	user.UpdateAt = time.Now()

	if err := au.repo.Update(ctx, id, user); err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, user)
}

func (au *accountUseCaseImpl) ReadOne(ctx context.Context, id int64) response.Response {
	user, err := au.repo.FindByID(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, user)
}

func (au *accountUseCaseImpl) Delete(ctx context.Context, id int64) response.Response {
	user, err := au.repo.FindByID(ctx, id)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}
	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	if err := au.repo.Delete(ctx, user.ID); err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	msg := "Success Delete Account"

	return response.Success(response.StatusOK, msg)
}
