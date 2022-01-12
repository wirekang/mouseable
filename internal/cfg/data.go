package cfg

import (
	"strconv"
)

type dataValue string

func (d dataValue) String() string {
	return string(d)
}

func (d dataValue) Int() int {
	i, _ := strconv.ParseInt(string(d), 10, 32)
	return int(i)
}

func (d dataValue) Float() float64 {
	i, _ := strconv.ParseFloat(string(d), 64)
	return i
}
