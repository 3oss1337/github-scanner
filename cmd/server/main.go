package main

import (
	"github.com/3oss1337/github-scanner/githubapi"
	"github.com/3oss1337/github-scanner/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	RepoOwner := "openapi"
	RepoName := "openai-cookbook"
	token := utils.GetGithubToken()

	client := githubapi.NewClient(token)

	repofiles := client.GetRepoFiles(RepoOwner, RepoName, "")

}
