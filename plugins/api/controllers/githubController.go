package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

const GITHUB_API = "https://api.github.com"

type gitHubUser struct {
	Login         string `json:"login,omitempty"`
	ID            int64  `json:"id,omitempty"`
	RepositoryURL string `json:"html_url,omitempty"`
}

type gitHubUserBadCredentials struct {
	Message          string `json:"message,omitempty"`
	DocumentationURL string `json:"documentation_url,omitempty"`
}

func GetGithubUser(authToken string) (*gitHubUser, *gitHubUserBadCredentials) {
	request, err1 := http.NewRequest("GET", GITHUB_API+"/user", nil)
	request.Header.Add("Authorization", "Bearer "+authToken)
	if err1 != nil {
		log.Fatalln(err1)
	}
	log.Println("Request creating successfull")

	response, err2 := httpClient.Do(request)
	if err2 != nil {
		log.Fatalln(err2)
	}
	if response.StatusCode == 401 {
		var badCredentials *gitHubUserBadCredentials
		errDecode := json.NewDecoder(response.Body).Decode(&badCredentials)
		if errDecode != nil {
			log.Fatalln(errDecode)
		}
		return nil, badCredentials
	}

	log.Println("Executing request successfull")

	var githubUser *gitHubUser
	err3 := json.NewDecoder(response.Body).Decode(&githubUser)
	if err3 != nil {
		log.Fatalln(err3)
	}

	log.Println(githubUser.Login)
	log.Println(githubUser.ID)
	log.Println(githubUser.RepositoryURL)
	return githubUser, nil
}
