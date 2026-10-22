package main

import (
	"net/http"
	"stuffRegistry/httpRouting"
)

func main() {
	router := httpRouting.GetRouter()
	err := http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
}
