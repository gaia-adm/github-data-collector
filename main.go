package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
	"os"
	"time"
	"strconv"
)

var proxyURL *url.URL

func main() {

	proxyString, found := os.LookupEnv("http_proxy")
	if found {
		proxyURL, _ = url.Parse(proxyString)
	}

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

	transport := &http.Transport{}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	client := &http.Client{Transport: transport}

	response, err := client.Get(url + "?per_page=1000")

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	/*	for k, v := range response.Header {
			fmt.Println(k, ": ", v)
		}*/

	rlimit, _ := strconv.Atoi(response.Header.Get("X-Ratelimit-Remaining"));
	tlimit, _ := strconv.ParseInt(response.Header.Get("X-Ratelimit-Reset"),10,64);
	llimit, _ := strconv.Atoi(response.Header.Get("X-Ratelimit-Limit"));

	if rlimit == 0 {
		log.Infof("You're out of tries - all %d calls have been used, come after %v. Think about authorized access.", llimit, time.Unix(tlimit, 0))
	} else {
		log.Infof("You can ask Github another %d times during this hour", rlimit)
	}

	return body
}

func typeof(v interface{}) string {

	return fmt.Sprintf("%T", v)
}