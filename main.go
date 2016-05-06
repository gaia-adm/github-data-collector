package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {

	for _, currRepository := range getRepositories() {
		log.Infof("Repository: %s, Commits: %d", currRepository["name"], len(getCommits(currRepository)))
	}
}

func getCommits(repository map[string]interface{}) []map[string]interface{} {

	return toMap(httpGet(getCommitsUrl(repository)))
}

func toMap(data []byte) []map[string]interface{} {

	var ret []map[string]interface{}
	json.Unmarshal(data, &ret)

	return ret
}

func getCommitsUrl(repository map[string]interface{}) string {

	var ret string
	commitsUrl := repository["commits_url"]
	if typeof(commitsUrl) == "string" {
		currUrl := commitsUrl.(string)
		if len(currUrl) > 0 && strings.HasSuffix(currUrl, "/commits{/sha}") {
			ret = strings.TrimSuffix(currUrl, "{/sha}")
		}
	}

	return ret
}

func getRepositories() []map[string]interface{} {

	return toMap(httpGet("https://api.github.com/orgs/gaia-adm/repos"))
}

func httpGet(url string) []byte {

	response, err := http.Get(url + "?per_page=1000")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func typeof(v interface{}) string {

	return fmt.Sprintf("%T", v)
}