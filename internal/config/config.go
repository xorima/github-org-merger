package config

import (
	"os"
)

var AppConfig = NewConfig()

type Config struct {
	GithubToken string
	SourceOrg   Org
}

type Org struct {
	Name string
}

func NewConfig() *Config {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		panic("GITHUB_TOKEN is not set")
	}
	return &Config{
		GithubToken: githubToken,
	}
}
