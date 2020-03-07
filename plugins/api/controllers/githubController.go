package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

const GITHUB_API = "https://api.github.com"

type GitHubUser struct {
	Login         string `json:"login,omitempty"`
	ID            int64  `json:"id,omitempty"`
	RepositoryURL string `json:"html_url,omitempty"`
}

type gitHubUserBadCredentials struct {
	Message          string `json:"message,omitempty"`
	DocumentationURL string `json:"documentation_url,omitempty"`
}

func GetGithubUser(authToken string) (*GitHubUser, *gitHubUserBadCredentials) {
	request, err1 := http.NewRequest("GET", GITHUB_API+"/user", nil)
	request.Header.Add("Authorization", "Bearer "+authToken)
	if err1 != nil {
		log.Fatalln(err1)
	}

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

	var githubUser *GitHubUser
	err3 := json.NewDecoder(response.Body).Decode(&githubUser)
	if err3 != nil {
		log.Fatalln(err3)
	}
	return githubUser, nil
}
