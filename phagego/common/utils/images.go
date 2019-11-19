package utils

import "regexp"

func GetImgSuffix(s string) string {
	re, _ := regexp.Compile(".(jpg|jpeg|png|gif)")
	return re.FindString(s)
}
