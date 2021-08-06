package internal

import "adlq/logger"

func CheckErr(err error, msg string) {
	if err != nil {
		logger.Error.Fatalln(msg, ":", err)
	}
}
