package api

import (
	"context"
	"github.com/dokku/dokku/plugins/api/controllers"
	"github.com/dokku/dokku/plugins/apps"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func ApiRoute() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	controllers.Client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	router.HandleFunc("/register", controllers.RegisterUserEndpoint).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUserEndpoint).Methods("POST")
	router.HandleFunc("/apps", AppsEndpoint).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.UsersEndpoint).Methods("PUT")
	router.HandleFunc("/users/sshKeys", controllers.GetSshKeys).Methods("GET")
	router.HandleFunc("/users/{id}/sshKeys", controllers.PutSshKeys).Methods("PUT")
	router.HandleFunc("/users/{id}/sshKeys", controllers.DeleteSshKeys).Methods("DELETE")
	http.ListenAndServe(":3000", router)
}

func AppsEndpoint(response http.ResponseWriter, request *http.Request) {
	apps.Apps_create("test1app")
}
