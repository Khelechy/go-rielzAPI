// base.go
package controllers

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres" //postgres

    "github.com/khelechy/rielzapi/api/middlewares"
	"github.com/khelechy/rielzapi/api/models"
    "github.com/khelechy/rielzapi/api/responses"
)

type App struct {
    Router *mux.Router
    DB     *gorm.DB
}

// Initialize connect to the database and wire up routes
func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
    var err error
    DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

    a.DB, err = gorm.Open("postgres", DBURI)
    if err != nil {
        fmt.Printf("\n Cannot connect to database %s", DbName)
        log.Fatal("This is the error:", err)
    } else {
        fmt.Printf("We are connected to the database %s", DbName)
    }

    a.DB.Debug().AutoMigrate(&models.User{}, &models.House{}) //database migration

    a.Router = mux.NewRouter().StrictSlash(true)
    a.initializeRoutes()
}

func (a *App) initializeRoutes() {
    a.Router.Use(middlewares.SetContentTypeMiddleware) // setting content-type to json

    a.Router.HandleFunc("/", home).Methods("GET")
    a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
    a.Router.HandleFunc("/login", a.Login).Methods("POST")
	a.Router.HandleFunc("/api/houses", a.GetHouses).Methods("GET")
	a.Router.HandleFunc("/api/houses/{id:[0-9]+}", a.GetHouseById).Methods("GET")
	

	s := a.Router.PathPrefix("/api").Subrouter() // routes that require authentication
    s.Use(middlewares.AuthJwtVerify)

    //s.HandleFunc("/users", a.GetAllUsers).Methods("GET")
    s.HandleFunc("/houses", a.CreateHouse).Methods("POST")
	s.HandleFunc("/houses/landlord/{id:[0-9]+}", a.GetHousesByLandlord).Methods("GET")
    s.HandleFunc("/houses/{id:[0-9]+}", a.UpdateHouse).Methods("PUT")
    s.HandleFunc("/houses/{id:[0-9]+}", a.DeleteHouse).Methods("DELETE")
}

func (a *App) RunServer() {
    log.Printf("\nServer starting on port 5000")
    log.Fatal(http.ListenAndServe(":5000", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) { // this is the home route
    responses.JSON(w, http.StatusOK, "Welcome To Ivents")
}