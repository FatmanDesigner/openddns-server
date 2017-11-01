package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type OAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Error       string `json:"error"`
}

// HttpServe OpenDDNS server at a given port
func HttpServe(port int) {
	log.Printf("Serving on port %d", port)

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/api/generate-secret", generateSecretHandler)
	http.HandleFunc("/oauth/github", oauthGithubCallback)

	fs := http.FileServer(http.Dir("/Users/khanhhua/dev/project-openddns/open-ddns-server/web-ui/dist"))
	http.Handle("/assets/", fs)
	http.Handle("/", fs)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func pingHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Bad Request")
		return
	}

	appid := req.URL.Query().Get("appid")
	if appid == "" {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Bad Request")
		return
	}

	if req.ContentLength == 0 {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Bad Request")
		return
	}

	body := make([]byte, req.ContentLength)
	req.Body.Read(body)

	_splat := strings.SplitAfterN(string(body), "\n", 2)
	secret, domainName := strings.TrimSpace(_splat[0]), strings.TrimSpace(_splat[1])

	if len(secret) == 0 {
		res.WriteHeader(http.StatusForbidden)
		io.WriteString(res, "Request not authorized")
		return
	}

	// TODO: Honor appid and secret as security measure

	if len(domainName) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Invalid domain name")
		return
	}

	// TODO: Validate domainName as per RFC 1034 - IETF (see https://www.ietf.org/rfc/rfc1034.txt)
	// TODO: Take into account req.Headers().Get("X-Forwarded-For")
	// .. when router is behind proxies
	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	ip := net.ParseIP(host)
	Register(domainName, ip.String())

	res.WriteHeader(http.StatusOK)
	io.WriteString(res, "OK")
}

func generateSecretHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "OK")
}

func oauthGithubCallback(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("Handling GitHub OAuth callback...")

	code := req.URL.Query().Get("code")

	// Interacting with GitHub OAuth get access token
	accessTokenURL := "https://github.com/login/oauth/access_token"
	data := map[string]string{
		"client_id":     "4b4ec20bb4bb4d306500",
		"client_secret": "caa8824df01a9538a85234b3e71ae3a19f873e3c",
		"code":          code,
	}
	body, _ := json.Marshal(data)
	req2, _ := http.NewRequest("POST", accessTokenURL, bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Accept", "application/json")
	client := &http.Client{}
	res2, err := client.Do(req2)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res2.Body.Close()

	bodyAccessToken, err := ioutil.ReadAll(res2.Body)
	log.Printf(string(bodyAccessToken))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var oauth OAuth
	if json.Unmarshal(bodyAccessToken, &oauth) != nil {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Could not parse OAuth response: "+string(bodyAccessToken))
		return
	}

	if len(oauth.Error) != 0 {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, oauth.Error)
		return
	}
	log.Printf("Authorization from GitHub user: access_token = %s token_type = %s",
		oauth.AccessToken, oauth.TokenType)

	// Interacting with GitHub user api to get User Profile
	userProfileURL := "https://api.github.com/user?access_token=" + oauth.AccessToken
	res3, err := http.Get(userProfileURL)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Could not call https://api.github.com/user")
		return
	}
	defer res3.Body.Close()

	bodyUserProfile, err := ioutil.ReadAll(res3.Body)

	var userProfile map[string]*json.RawMessage
	if json.Unmarshal(bodyUserProfile, &userProfile) != nil {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Could not parse user profile response: "+string(bodyUserProfile))
		return
	}

	GitHubUserID := string([]byte(*userProfile["id"]))
	log.Printf("Unmarshaled UserID: %s", GitHubUserID)

	log.Printf("User profile: User Login = %s User ID = %s",
		string([]byte(*userProfile["login"])), GitHubUserID)

	_, _, ok := GenerateApp(string(GitHubUserID))
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Could not create appid and secret")

		return
	}

	res.WriteHeader(http.StatusOK)
	io.WriteString(res, "OK")
}
