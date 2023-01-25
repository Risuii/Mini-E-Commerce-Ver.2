package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"

	"github.com/Risuii/config"
	"github.com/Risuii/config/bcrypt"
	"github.com/Risuii/helpers/constant"
	"github.com/Risuii/internal/account"
	"github.com/Risuii/internal/item"
	"github.com/Risuii/internal/store"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("mysql", cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	validator := validator.New()
	router := mux.NewRouter()
	bcrypt := bcrypt.NewBcrypt(cfg.Bcrypt.HashCost)

	userRepo := account.NewAccountRepositoryImpl(db, constant.TableAccount)
	storeRepo := store.NewStoreRepository(db, constant.TableStores)
	itemRepo := item.NewItemRepositoryImpl(db, constant.TableItems)
	userUseCase := account.NewAccountUseCaseImpl(userRepo, bcrypt)
	storeUseCase := store.NewStoreUseCaseImpl(storeRepo)
	itemUseCase := item.NewItemUseCaseImpl(itemRepo)

	account.NewAbsensiHandler(router, validator, userUseCase)
	store.NewStoreHandler(router, validator, storeUseCase)
	item.NewItemHandler(router, validator, itemUseCase)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	port := os.Getenv("PORT")

	fmt.Println("SERVER ON")
	fmt.Println("PORT :", port)
	log.Fatal(server.ListenAndServe())
}
