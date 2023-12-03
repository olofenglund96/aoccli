package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/aoccli/helpers"
)

type AOCClient struct {
	domain     string
	httpClient *http.Client
}

func NewAOCClient(domain string, token string) (AOCClient, error) {
	aocUrl := fmt.Sprintf("https://%s", domain)
	cUrl, err := url.Parse(aocUrl)
	if err != nil {
		return AOCClient{}, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return AOCClient{}, err
	}

	jar.SetCookies(cUrl, []*http.Cookie{
		{
			Name:  "session",
			Value: token,
		},
	})

	return AOCClient{
		domain: domain,
		httpClient: &http.Client{
			Jar: jar,
		},
	}, nil
}

func (ac AOCClient) GetDayInput(year int, day int) (string, error) {
	response, err := ac.httpClient.Get(helpers.GetDayInputUrl(ac.domain, year, day))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	return string(body), nil
}

func (ac AOCClient) SubmitProblem(year int, day int, problem int, solution string) (string, error) {
	postBody := url.Values{
		"level":  []string{strconv.Itoa(problem)},
		"answer": []string{solution},
	}
	response, err := ac.httpClient.PostForm(helpers.GetDaySubmitUrl(ac.domain, year, day), postBody)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	return string(body), nil
}
