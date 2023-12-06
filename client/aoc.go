package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

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

func (ac AOCClient) GetDayInput(year string, day string) (string, error) {
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

func (ac AOCClient) SubmitProblem(year string, day string, problem int, solution string) (string, error) {
	postBody := url.Values{
		"level":  []string{string(problem)},
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

func (ac AOCClient) GetDayTestInput(year string, day string) (string, error) {
	response, err := ac.httpClient.Get(helpers.GetDayUrl(ac.domain, year, day))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	return parseTestInput(response.Body)
}

func parseTestInput(responseBody io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return "", err
	}

	dayDescription := doc.Find(".day-desc").First()

	foundExampleText := false
	var code string
	dayDescription.Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		selectorText := s.Text()
		if foundExampleText && s.Is("pre") {
			code = s.Find("code").First().Text()
			return false
		}

		if strings.Contains(selectorText, "example") && strings.Index(selectorText, ":") == len(selectorText)-1 {
			foundExampleText = true
			return true
		}

		return true
	})

	if code == "" {
		fmt.Println("Failed to find test input.. Please get it yourself 🥰")
	}

	return code, nil
}
