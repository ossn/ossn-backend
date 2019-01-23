package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

)

var (
	// Redirect urls
	FrontendURL = os.Getenv("FRONTEND_URL")
	LoginURL    string
	BackendURL  = os.Getenv("BACKEND_URL")
	src         = rand.NewSource(time.Now().UnixNano())

	// HTTTP Client
	timeout    = time.Duration(7 * time.Second)
	httpClient = &http.Client{Timeout: timeout}
	githubKey  = os.Getenv("GITHUB_KEY")
)

func init() {
	CheckEnvVariable(&FrontendURL, "FRONTEND_URL")
	CheckEnvVariable(&BackendURL, "BACKEND_URL")
	CheckEnvVariable(&githubKey, "GITHUB_KEY")
	LoginURL = FrontendURL + "login?token="
}

type GithubRes struct {
	Data struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
	} `json:"data"`
}

// RandStringBytesMaskImprSrc creates a random string of size n
//
// Note: This isn't thread safe
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// HandleError is a generic http error handler.
func HandleError(w http.ResponseWriter, r *http.Request, code int, err error) {
	http.Redirect(w, r, FrontendURL+"?status=error&action=login&code="+strconv.Itoa(code), http.StatusTemporaryRedirect)
	fmt.Println("Login Error: " + err.Error())
}

// CheckEnvVariable validates that a string is present.
// IF the string is present it will panic.
func CheckEnvVariable(str *string, name string) {
	if len(*str) == 0 {
		log.Fatal("Please set " + name)
	} else if strings.HasSuffix(name, "URL") && !strings.HasSuffix(*str, "/") {
		*str += "/"
	}
}

// GetGithubURL validates if a github username exists and
// emits the response to the channel.
func GetGithubURL(oidcID, username *string, c chan string) {
	defer func() {
		if r := recover(); r != nil {
			c <- ""
		}
		close(c)
	}()
	if strings.HasPrefix(*oidcID, "github") {
		str := "{\"query\":\"{ user(login: \"" + *username + "\") {id}}\"}"
		b, err := json.Marshal(str)
		if err != nil {
			fmt.Println(err)
			c <- ""
			return
		}
		req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(b))
		if err != nil {
			fmt.Println(err)
			c <- ""
			return
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "bearer "+githubKey)

		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			c <- ""
			return
		}

		defer res.Body.Close()
		resBody := &GithubRes{}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			c <- ""
			return
		}

		err = json.Unmarshal(data, resBody)
		if err != nil {
			fmt.Println(err)
			c <- ""
			return
		}
		c <- "https://github.com/" + *username
		return
	}
	c <- ""
}
