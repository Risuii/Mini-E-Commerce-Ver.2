package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"

	"github.com/Risuii/config"
)

func main() {
	cfg := config.New()

	_, err := sql.Open("mysql", cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	// validator := validator.New()
	router := mux.NewRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	port := os.Getenv("PORT")

	fmt.Println("SERVER ON")
	fmt.Println("PORT :", port)
	log.Fatal(server.ListenAndServe())
}
