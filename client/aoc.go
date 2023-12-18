package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/aoccli/helpers"
)

type AOCClient struct {
	domain     string
	httpClient *http.Client
}

func NewAOCClient(token string) (AOCClient, error) {
	aocUrl := "https://adventofcode.com"
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
		domain: aocUrl,
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
		"level":  []string{strconv.Itoa(problem)},
		"answer": []string{solution},
	}
	response, err := ac.httpClient.PostForm(helpers.GetDaySubmitUrl(ac.domain, year, day), postBody)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return "", err
	}

	return doc.Find("article").First().Find("p").First().Text(), nil
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
		return "", fmt.Errorf("error parsing HTML: %s", err)
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

func (ac AOCClient) GetLeaderboard(leaderboardId string) (LeaderboardResponse, error) {

	response, err := ac.httpClient.Get(helpers.GetLeaderboardUrl(ac.domain, leaderboardId))
	if err != nil {
		return LeaderboardResponse{}, fmt.Errorf("error when getting leadeboard: %s", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return LeaderboardResponse{}, fmt.Errorf("got status %d when fetching leader board.", response.StatusCode)
	}

	defer response.Body.Close()

	var leaderboard LeaderboardResponse
	err = json.NewDecoder(response.Body).Decode(&leaderboard)
	if err != nil {
		return LeaderboardResponse{}, fmt.Errorf("error when parsing leadeboard: %s", err)
	}

	return leaderboard, nil
}

type (
	LeaderboardResponse struct {
		OwnerId   int               `json:"owner_id"`
		Event     string            `json:"event"`
		MemberMap map[string]Member `json:"members"`
	}
	Member struct {
		Name               string                     `json:"name"`
		Stars              int                        `json:"stars"`
		LocalScore         int                        `json:"local_score"`
		CompletionDayLevel map[string]CompletionLevel `json:"completion_day_level"`
	}
	CompletionLevel map[string]any
)

func (l LeaderboardResponse) String() string {

	var members []Member
	for _, m := range l.MemberMap {
		members = append(members, m)
	}
	sort.Slice(members, func(i, j int) bool {
		return members[i].LocalScore > members[j].LocalScore
	})

	s := "\t 1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25\n"

	pos := 1
	for _, m := range members {
		starString := ""
		for i := 1; i <= 25; i++ {
			if c, ok := m.CompletionDayLevel[strconv.Itoa(i)]; ok {
				if len(c) == 1 {
					starString = starString + " ☆ "
				} else {
					starString = starString + " ★ "
				}

			} else {
				starString = starString + " - "
			}
		}
		s = s + fmt.Sprintf("%d) %d\t%s \t %s (%d stars)\n", pos, m.LocalScore, starString, m.Name, m.Stars)
		pos = pos + 1
	}

	return s
}
