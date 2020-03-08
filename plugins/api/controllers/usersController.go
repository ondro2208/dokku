package controllers

import (
	"context"
	"encoding/json"
	sshKeys "github.com/dokku/dokku/plugins/ssh-keys"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ServiceID int64              `json:"serviceId,omitempty" bson:"serviceId,omitempty"`
	UserName  string             `json:"userName,omitempty" bson:"userName,omitempty"`
	SshKey    string             `json:"sshKey,omitempty" bson:"sshKey,omitempty"`
}

type PutUser struct {
	AuthToken string `json:"auth_token"`
	NewUser   User   `json:"user"`
}

func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	collection := Client.Database("tmp_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var results []User
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(results)
}

func UsersEndpoint(response http.ResponseWriter, request *http.Request) {
	var idParam string = mux.Vars(request)["id"]
	response.Header().Set("content-type", "application/json")

	//decode request body
	var putUser PutUser
	_ = json.NewDecoder(request.Body).Decode(&putUser)

	//check user authenticated
	githubUser, badCredentials := GetGithubUser(putUser.AuthToken)
	if githubUser == nil {
		json.NewEncoder(response).Encode(badCredentials)
		response.WriteHeader(401)
		return
	}

	id, _ := primitive.ObjectIDFromHex(idParam)
	newUser := putUser.NewUser
	collection := Client.Database("tmp_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var existingUser User
	collection.FindOne(ctx, User{ID: id}).Decode(&existingUser)

	//check updated user is authorized to edit such user
	if existingUser.ServiceID != githubUser.ID {
		response.WriteHeader(401)
		return
	}

	userToUpdate := mergeUsers(existingUser, newUser)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set", userToUpdate}},
	)

	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		log.Println("matched and replaced an existing document")
		return
	}
	if result.UpsertedCount != 0 {
		log.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	//TODO
	//if new ssh remove old and add new
	if len(userToUpdate.SshKey) > 0 && existingUser.SshKey != newUser.SshKey {
		//remove old from dokku

		//add new to dokku
	}

	json.NewEncoder(response).Encode(newUser)
}

func mergeUsers(existingUser User, newUser User) (userToUpdate User) {
	var id = existingUser.ID
	var serviceID = existingUser.ServiceID

	var sshKey string
	if len(newUser.SshKey) > 0 && newUser.SshKey != existingUser.SshKey {
		sshKey = newUser.SshKey
	} else {
		sshKey = existingUser.SshKey
	}

	var userName string
	if len(newUser.UserName) > 0 && newUser.UserName != existingUser.UserName {
		userName = newUser.UserName
	} else {
		userName = existingUser.UserName
	}

	return User{ID: id, ServiceID: serviceID, SshKey: sshKey, UserName: userName}
}

func GetSshKeys(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var sshKeysOutput = sshKeys.ListSshKeys()

	jsonData, err := json.Marshal(sshKeyNames{SshKeys: sshKeysOutput})
	if err != nil {
		log.Fatal(err)
	}

	response.Write(jsonData)
}

func DeleteSshKeys(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var idParam string = mux.Vars(request)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)
	collection := Client.Database("tmp_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var existingUser User
	collection.FindOne(ctx, User{ID: id}).Decode(&existingUser)
	var sshKeyName = existingUser.UserName

	isAdded, err := sshKeys.RemoveSshKey(sshKeyName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(isAdded)

	jsonData, err := json.Marshal(map[string]string{"message": "SSH key deleted succesfully"})
	if err != nil {
		log.Fatal(err)
	}

	response.Write(jsonData)
}

func PutSshKeys(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var idParam string = mux.Vars(request)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)
	collection := Client.Database("tmp_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var existingUser User
	collection.FindOne(ctx, User{ID: id}).Decode(&existingUser)
	var sshKeyName = existingUser.UserName

	var sshKeyValue sshKey
	_ = json.NewDecoder(request.Body).Decode(&sshKeyValue)

	isAdded, err := sshKeys.AddSshKey(sshKeyName, sshKeyValue.SshKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(isAdded)

	jsonData, err := json.Marshal(map[string]string{"message": "SSH key added succesfully"})
	if err != nil {
		log.Fatal(err)
	}

	response.Write(jsonData)
}

type sshKeyNames struct {
	SshKeys []string `json:"sshKeyNames"`
}

type sshKey struct {
	SshKey string `json:"sshKey"`
}
