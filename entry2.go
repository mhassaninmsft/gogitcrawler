// main.go
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/mhassaninmsft/gogitcrawler/database"
	"github.com/mhassaninmsft/gogitcrawler/githubapi"
	"github.com/mhassaninmsft/gogitcrawler/models"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: github-contributors <repo_owner> <repo_name> <github_token>")
		os.Exit(1)
	}

	repoOwner := os.Args[1]
	repoName := os.Args[2]
	githubToken := os.Args[3]

	// Initialize the GitHub API client
	apiClient := githubapi.NewClient(githubToken)

	// Get the list of contributors
	contributors, err := apiClient.GetContributors(repoOwner, repoName)
	if err != nil {
		fmt.Printf("Error fetching contributors: %v\n", err)
		os.Exit(1)
	}

	// Initialize the PostgreSQL database
	db, err := database.Connect("postgres://mhassanin:magical_password@localhost/gitrepos?sslmode=disable")
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		os.Exit(1)
	}

	// Use goroutines and WaitGroup for parallelism
	var wg sync.WaitGroup
	wg.Add(len(contributors))

	// Iterate through the contributors and fetch their public repositories
	for _, contributor := range contributors {
		go func(contributor *models.Contributor) {
			defer wg.Done()

			repos, err := apiClient.GetPublicRepos(contributor.Login)
			if err != nil {
				fmt.Printf("Error fetching public repos for %s: %v\n", contributor.Login, err)
				return
			}

			contributor.Repos = repos

			// Save the contributor and their repositories to the database
			if err := db.SaveContributor(contributor); err != nil {
				fmt.Printf("Error saving contributor %s to the database: %v\n", contributor.Login, err)
				return
			}
			fmt.Printf("Saved contributor %s and their %d public repositories to the database.\n", contributor.Login, len(repos))
		}(contributor)
	}

	wg.Wait()
	fmt.Println("All contributors and their repositories have been saved to the database.")
}
