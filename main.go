package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rishit1234567889/carZone/driver"
	carHandler "github.com/rishit1234567889/carZone/handler/car"
	engineHandler "github.com/rishit1234567889/carZone/handler/engine"
	"github.com/rishit1234567889/carZone/middleware"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	loginHandler "github.com/rishit1234567889/carZone/handler/login"
	carService "github.com/rishit1234567889/carZone/service/car"
	engineService "github.com/rishit1234567889/carZone/service/engine"
	carStore "github.com/rishit1234567889/carZone/store/car"
	engineStore "github.com/rishit1234567889/carZone/store/engine"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()

	carstore := carStore.New(db)
	carservice := carService.NewCarService(carstore)

	enginestore := engineStore.New(db)
	engineservice := engineService.NewEngineService(enginestore)

	carhandler := carHandler.NewCarHandler(carservice)
	enginehandler := engineHandler.NewEngineHandler(engineservice)

	router := mux.NewRouter()

	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatal("Error while executing schema file ", err)
	}

	router.HandleFunc("/login", loginHandler.LoginHandler).Methods("POST")

	// Middleware
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	//router.HandleFunc("/ping", Ping).Methods("GET")

	protected.HandleFunc("/cars/{id}", carhandler.GetCarByID).Methods("GET")
	protected.HandleFunc("/cars", carhandler.GetCarByBrand).Methods("GET")
	protected.HandleFunc("/cars", carhandler.CreateCar).Methods("POST") // Need to correct function call
	protected.HandleFunc("/cars/{id}", carhandler.UpdateCar).Methods("PUT")
	protected.HandleFunc("/cars/{id}", carhandler.DeleteCar).Methods("DELETE")

	protected.HandleFunc("/engine/{id}", enginehandler.GetEngineByID).Methods("GET")
	protected.HandleFunc("/engine", enginehandler.CreateEngine).Methods("POST")
	protected.HandleFunc("/engine/{id}", enginehandler.UpdateEngine).Methods("PUT")
	protected.HandleFunc("/engine/{id}", enginehandler.DeleteEngine).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))

}

func executeSchemaFile(db *sql.DB, schemaFile string) error {
	sqlFile, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(sqlFile)); err != nil {
		return err
	}
	return nil
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
