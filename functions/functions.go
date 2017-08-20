package functions

import (
	"log"
	"strconv"
	"strings"
)

func StringToFloat64(str string) (float64, error) {
	val, err := strconv.ParseFloat(strings.Replace(str, ",", ".", 1), 64)
	return val, err
}

func Float64ToString(flo float64) string {
	return strconv.FormatFloat(flo, 'f', 4, 64)
}

func ChErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
