package main

import (
	"strconv"
)

//IntToString is to convert int to string
func IntToString(val int) string {
	return strconv.FormatInt(int64(val), 10)
}

//Int64ToString is to convert int64 to string
func Int64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}
