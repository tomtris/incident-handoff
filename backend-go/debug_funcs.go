package main

import (
	"fmt"
	"io"
	"net/http"
)

func printBody(r http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println(string(body))
}
