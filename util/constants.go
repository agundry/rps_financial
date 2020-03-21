package util

import (
	"time"
)

type Name = string
const (
	AUSTIN Name = "AUSTIN"
	SAM Name = "SAM"
)

func GetEpochSeconds() int64 {
	return time.Now().Unix()
}



