package main

import (
	"fmt"
	"net/http"

	"github.com/pathvars"
)

var matcher = pathvars.NewExtractor("/user/{userid}")

func main() {
	r, _ := http.NewRequest("GET", "/user/123", nil)
	values, ok := matcher.Extract(r.URL)
	fmt.Println("OK:", ok)
	fmt.Println("User ID:", values["userid"])
}
