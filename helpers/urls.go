package helpers

import "fmt"

func GetDayUrl(domain string, year string, day string) string {
	return fmt.Sprintf("https://%s/%d/day/%d", domain, year, day)
}

func GetDayInputUrl(domain string, year string, day string) string {
	return fmt.Sprintf("%s/input", GetDayUrl(domain, year, day))
}

func GetDaySubmitUrl(domain string, year string, day string) string {
	return fmt.Sprintf("%s/answer", GetDayUrl(domain, year, day))
}
