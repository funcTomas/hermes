package tool

import "time"

func GetNowDate() int {
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()
	return year*10000 + month*100 + day
}
