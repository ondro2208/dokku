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
	http.ListenAndServe(":3000", router)
}

func AppsEndpoint(response http.ResponseWriter, request *http.Request) {
	apps.Apps_create("test1app")
}
