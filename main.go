package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
)

type Message struct {
	comments_url string
}

func main() {

	log.Info(getRepositoriesCommitUrls())
}

func getRepositories() []byte {

	response, err := http.Get("https://api.github.com/orgs/gaia-adm/repos?per_page=1000")
	if err != nil {
		log.Error("Failed to get gaia-adm repositories from github. ", err)
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Failed to read gaia-adm repositories' response JSON from github. ", err)
		panic(err)
	}

	return body
}

func getRepositoriesCommitUrls() []string {

	var ret []string
	var repositories []map[string]interface{}
	json.Unmarshal(getRepositories(), &repositories)
	for _, currRepository := range repositories {
		commitsUrl := currRepository["commits_url"]
		if typeof(commitsUrl) == "string" {
			currUrl := currRepository["commits_url"].(string)
			if len(currUrl) > 0 && strings.HasSuffix(currUrl, "/commits{/sha}") {
				ret = append(ret, strings.TrimSuffix(currUrl, "{/sha}"))
			}
		}
	}

	return ret
}

func typeof(v interface{}) string {

	return fmt.Sprintf("%T", v)
}