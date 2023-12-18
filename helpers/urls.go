package helpers

import "fmt"

func GetDayUrl(url string, year string, day string) string {
	return fmt.Sprintf("%s/%s/day/%s", url, year, day)
}

func GetDayInputUrl(url string, year string, day string) string {
	return fmt.Sprintf("%s/input", GetDayUrl(url, year, day))
}

func GetDaySubmitUrl(url string, year string, day string) string {
	return fmt.Sprintf("%s/answer", GetDayUrl(url, year, day))
}

func GetLeaderboardUrl(url string, leaderboardId string) string {
	return fmt.Sprintf("%s/leaderboard/private/view/%s.json", url, leaderboardId)
}
