package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mhassaninmsft/gogitcrawler/util"
)
type Contributor struct {
	Login string `json:"login"`
}
// import ("github.com/mhassaninmsft/gogitcrawler/util" "")
func main_old() {
	println("Hello, World!")
	var x = util.AddTwo(1, 2)
	println(x)
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: go run main.go <repository>")
	// 	os.Exit(1)
	// }

	// repo := os.Args[1]
	GITHUB_TOKEN, found := os.LookupEnv("GITHUB_TOKEN")
	if !found {
		fmt.Println("GITHUB_TOKEN environment variable not found")
		os.Exit(1)
	}
	url := "https://api.github.com/repos/google/material-design-icons/contributors"
	authorization_token := fmt.Sprintf("Bearer %s", GITHUB_TOKEN)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", authorization_token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error:", resp.Status)
		os.Exit(1)
	}

	var contributors []Contributor
	err = json.NewDecoder(resp.Body).Decode(&contributors)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		os.Exit(1)
	}

	for _, c := range contributors {
		fmt.Println(c.Login)
	}
}