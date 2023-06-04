package helpers

import (
	"fmt"
	"regexp"
	"time"
)

// func to get time duration from given time until now
func GetTimeDiffNow(timeToCheck time.Time) time.Duration {
	timeToCheck = time.Date(
		timeToCheck.Year(), timeToCheck.Month(), timeToCheck.Day(),
		timeToCheck.Hour(), timeToCheck.Minute(), 0, 0, time.Local)

	return time.Now().Sub(timeToCheck)
}

func PrepareTextLinks(text string) string {

	rexp, err := regexp.Compile(`\bhttp(\S+)`)
	if err != nil {
		fmt.Printf("Rexp err: %v\n", err)
		return text
	}

	links := rexp.FindAllString(text, -1)

	for _, link := range links {
		strToReplace := fmt.Sprintf("<a href='%s' target='_blank'>%s</a>", link, link)

		re, linkErr := regexp.Compile(regexp.QuoteMeta(link))

		if linkErr != nil {
			fmt.Printf("Compile link err")
			continue
		}

		text = re.ReplaceAllString(text, strToReplace)
	}

	return text
}
