package helpers

import "fmt"

func GetDayUrl(domain string, year string, day string) string {
	return fmt.Sprintf("%s/%s/day/%s", domain, year, day)
}

func GetDayInputUrl(domain string, year string, day string) string {
	return fmt.Sprintf("%s/input", GetDayUrl(domain, year, day))
}

func GetDaySubmitUrl(domain string, year string, day string) string {
	return fmt.Sprintf("%s/answer", GetDayUrl(domain, year, day))
}

func GetLeaderboardUrl(domain string, leaderboardId string) string {
	return fmt.Sprintf("%s/leaderboard/private/view/%s.json", domain, leaderboardId)
}
