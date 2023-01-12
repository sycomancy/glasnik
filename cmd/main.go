package main

import (
	"fmt"
	"regexp"
)

var data = `Stan u stambenoj zgradi, prizemlje
																																																								Stambena površina: 256.62 m2
																																																								Lokacija: Črnomerec, Bijenik`

var leftDelimiter = "ina:"
var rightDelimiter = "m2"

func main() {
	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(leftDelimiter) + `(.*?)` + regexp.QuoteMeta(rightDelimiter))
	matches := rx.FindAllStringSubmatch(data, -1)
	fmt.Println(matches[0][1])
}
