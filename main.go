package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"github.com/joho/godotenv"
)

var githubOauth2Config = &oauth2.Config {
	ClientID: "",
	ClientSecret: "",
	Endpoint: github.Endpoint,
}

type githubResponse struct {
	Data struct {
		Viewer struct {
			ID string `json:"id"`
		} `json:"viewer"`
	} `json:"data"`
}

var gstate = "0000"

func main() {
	godotenv.Load(".env")
	githubOauth2Config.ClientID = os.Getenv("CLIENT_ID")
	githubOauth2Config.ClientSecret = os.Getenv("CLIENT_SECRET")

	http.HandleFunc("/", index)
	http.HandleFunc("/oauth2/github", startGithubOauth2)
	http.HandleFunc("/oauth2/receive", receiveGithubOauth2)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>oauth2 test</title>
	</head>
	<body>
		<form action="/oauth2/github" method="post">
			<input type="submit" value="Login with GitHub">
		</form>
	</body>
</html>`)
}

func startGithubOauth2(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, githubOauth2Config.AuthCodeURL(gstate), http.StatusSeeOther)
}

func receiveGithubOauth2(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	if state != gstate {
		http.Error(w, "State is incorrect", http.StatusBadRequest)
		return
	}

	token, err := githubOauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Couldnt login", http.StatusInternalServerError)
		return
	}

	ts := githubOauth2Config.TokenSource(r.Context(), token)
	client := oauth2.NewClient(r.Context(), ts)

	requestBody := strings.NewReader(`{ "query": "query { viewer {id} }"}`)
	resp, err := client.Post("https://api.github.com/graphql", "application/json", requestBody)
	if err != nil {
		http.Error(w, "Couldnt get user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var gr githubResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Github invalid response", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(gr.Data.Viewer.ID))
}
