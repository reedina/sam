package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/reedina/sam/ctrl"
	"github.com/reedina/sam/model"
	"github.com/rs/cors"

	//Initialize pq driver
	_ "github.com/lib/pq"
)

//App  (TYPE)
type App struct {
	Router *mux.Router
}

//InitializeApplication - Init router, db connection and restful routes
func (a *App) InitializeApplication(user, password, url, dbname string) {

	model.ConnectDB(user, password, url, dbname)
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

//InitializeRoutes - Declare all application routes
func (a *App) InitializeRoutes() {

	// Build Request
	a.Router.HandleFunc("/api/buildRequest/{email}", ctrl.GetBuildRequestProfile).Methods("GET")
	a.Router.HandleFunc("/api/buildRequest", ctrl.CreateBuild).Methods("POST")
}

//RunApplication - Start the HTTP server
func (a *App) RunApplication(addr string) {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	log.Fatal(http.ListenAndServe(addr, c.Handler(a.Router)))
}
