package util

import "log"

func CheckError(errorMessage string, err error) {
	if err != nil {
		log.Fatal(errorMessage, err)
	}
}
