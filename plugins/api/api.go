package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dokku/dokku/plugins/apps"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var client *mongo.Client
var httpClient = &http.Client{}

const GITHUB_API = "https://api.github.com"

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ServiceID int64              `json:"serviceId,omitempty" bson:"serviceId,omitempty"`
	UserName  string             `json:"userName,omitempty" bson:"userName,omitempty"`
}

type RegisterUser struct {
	AuthToken string `json:"auth_token,omitempty"`
}

type LoginUser struct {
	AuthToken string `json:"auth_token,omitempty"`
}

type GitHubUser struct {
	Login   string `json:"login,omitempty"`
	ID      int64  `json:"id,omitempty"`
	HtmlUrl string `json:"html_url,omitempty"`
}

type GitHubUserBadCredentials struct {
	Message          string `json:"message,omitempty"`
	DocumentationURL string `json:"documentation_url,omitempty"`
}

func ApiRoute() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	router.HandleFunc("/register", RegisterUserEndpoint).Methods("POST")
	router.HandleFunc("/login", LoginUserEndpoint).Methods("POST")
	router.HandleFunc("/apps", AppsEndpoint).Methods("POST")
	http.ListenAndServe(":3000", router)
}

func AppsEndpoint(response http.ResponseWriter, request *http.Request) {
	apps.Apps_create("test1app")
}

func RegisterUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var registerBody RegisterUser
	_ = json.NewDecoder(request.Body).Decode(&registerBody)
	fmt.Println(registerBody)

	var user User
	githubUser, badCredentials := getGithubUser(registerBody.AuthToken)
	if githubUser == nil {
		json.NewEncoder(response).Encode(badCredentials)
		return
	}
	collection := client.Database("tmp_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ServiceID: githubUser.ID}).Decode(&user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		result, _ := collection.InsertOne(ctx, User{ServiceID: githubUser.ID, UserName: githubUser.Login})
		collection.FindOne(ctx, User{ID: result.InsertedID.(primitive.ObjectID)}).Decode(&user)
	}
	json.NewEncoder(response).Encode(user)
}

func LoginUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var loginUser LoginUser
	_ = json.NewDecoder(request.Body).Decode(&loginUser)
	fmt.Println(loginUser)

	var user User
	githubUser, badCredentials := getGithubUser(loginUser.AuthToken)
	if githubUser == nil {
		json.NewEncoder(response).Encode(badCredentials)
		return
	}
	collection := client.Database("tmp_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ServiceID: githubUser.ID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func getGithubUser(authToken string) (*GitHubUser, *GitHubUserBadCredentials) {
	request, err1 := http.NewRequest("GET", GITHUB_API+"/user", nil)
	request.Header.Add("Authorization", "Bearer "+authToken)
	if err1 != nil {
		log.Fatalln(err1)
	}
	fmt.Println("Request creating successfull")

	response, err2 := httpClient.Do(request)
	if err2 != nil {
		log.Fatalln(err2)
	}
	if response.StatusCode == 401 {
		var badCredentials *GitHubUserBadCredentials
		errDecode := json.NewDecoder(response.Body).Decode(&badCredentials)
		if errDecode != nil {
			log.Fatalln(errDecode)
		}
		return nil, badCredentials
	}

	fmt.Println("Executing request successfull")

	var githubUser *GitHubUser
	err3 := json.NewDecoder(response.Body).Decode(&githubUser)
	if err3 != nil {
		log.Fatalln(err3)
	}

	fmt.Println(githubUser.Login)
	fmt.Println(githubUser.ID)
	fmt.Println(githubUser.HtmlUrl)
	return githubUser, nil
}
