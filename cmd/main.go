package main

import (
	"backlog"
	"backlog/router"
	"log"
	"net/http"
)

func main() {
	r := &router.Router{}
	r.HandlerFunc("/userstories/open", http.MethodGet, backlog.GetOpenUserStories)
	r.HandlerFunc("/userstories", http.MethodPost, backlog.AddUserStory)
	r.HandlerFunc(`/userstories/\d`, http.MethodPut, backlog.UpdateUserStory)
	log.Fatal(http.ListenAndServe(":8443", r))
}
