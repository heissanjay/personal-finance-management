package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/heissanjay/personal-finance-management/internal/auth"
	"github.com/heissanjay/personal-finance-management/internal/config"
	"github.com/heissanjay/personal-finance-management/internal/router"
	"github.com/heissanjay/personal-finance-management/internal/storage/postgres"
	"github.com/heissanjay/personal-finance-management/internal/user"
	_ "github.com/lib/pq"
)

func main() {

	config.LoadConfig()
	dbURL := "postgres://" + config.Config.DBUser + ":" + config.Config.DBPassword + "@" + config.Config.DBHost + ":" + config.Config.DBPort + "/" + config.Config.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to Database", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	datastore := postgres.NewPostgresDB(db)
	datastore.InitDB()

	authService := auth.NewAuthService(datastore)
	userHandler := user.NewHandler(authService)

	router := router.NewRouter(userHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
