package emqx

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 28 days, 21 hours, 55 minutes, 29 seconds
func TestPath(t *testing.T) {
	f := "1 days, 21 hours, 0 minutes, 29 seconds"
	f = strings.ReplaceAll(f, " ", "")
	values := strings.Split(f, ",")

	total := 0
	for _, value := range values {

		if strings.Contains(value, "days") {
			vstr := strings.ReplaceAll(value, "days", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v * 86400
			}
		}

		if strings.Contains(value, "hours") {
			vstr := strings.ReplaceAll(value, "hours", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v * 3600
			}
		}

		if strings.Contains(value, "minutes") {
			vstr := strings.ReplaceAll(value, "minutes", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v * 60
			}
		}

		if strings.Contains(value, "seconds") {
			vstr := strings.ReplaceAll(value, "seconds", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v
			}
		}
	}
	fmt.Println(total)
}

func TestStringToG(t *testing.T) {

	str := "4.70G"
	str = strings.ReplaceAll(str, " ", "")

	if strings.Contains(str, "G") {
		str = strings.ReplaceAll(str, "G", "")

		size, err := strconv.ParseFloat(str, 64)
		if err == nil {
			fmt.Println(size)
		}
	}
}
