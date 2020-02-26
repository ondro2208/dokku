package controllers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var Client *mongo.Client
var httpClient = &http.Client{}

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

func RegisterUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var registerBody RegisterUser
	_ = json.NewDecoder(request.Body).Decode(&registerBody)
	log.Println(registerBody)

	var user User
	githubUser, badCredentials := GetGithubUser(registerBody.AuthToken)
	if githubUser == nil {
		json.NewEncoder(response).Encode(badCredentials)
		return
	}
	collection := Client.Database("tmp_api").Collection("users")
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
	log.Println(loginUser)

	var user User
	githubUser, badCredentials := GetGithubUser(loginUser.AuthToken)
	if githubUser == nil {
		json.NewEncoder(response).Encode(badCredentials)
		return
	}
	collection := Client.Database("tmp_api").Collection("users")
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
