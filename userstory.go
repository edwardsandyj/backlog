package backlog

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrUserStoryNotFound = errors.New("User Story was not found")

// User story that is still open or closed
type UserStory struct {
	ID          int    `json:"id"`          // story identifier
	Description string `json:"description"` // story description
	Closed      bool   `json:"closed"`      // story completion status
}

// Datastore holds all user stories saved in the application
type Datastore struct {
	UserStories []UserStory
	lastID      int // mark IDs that are already used in UserStories
}

var ds = &Datastore{}

// SaveUserStory creates new or updates an existing UserStory in Datastore
func (ds *Datastore) SaveUserStory(story UserStory) error {
	if story.ID == 0 {
		ds.lastID++
		story.ID = ds.lastID
		ds.UserStories = append(ds.UserStories, story)
		return nil
	}

	for i, t := range ds.UserStories {
		if t.ID == story.ID {
			ds.UserStories[i] = story
			return nil
		}
	}
	return ErrUserStoryNotFound
}

// GetOpenUserStories returns all unfinished user stories in Datastore
func (ds *Datastore) GetOpenUserStories() []UserStory {
	var openstories []UserStory
	for _, story := range ds.UserStories {
		if !story.Closed {
			openstories = append(openstories, story)
		}
	}
	return openstories
}

// GetOpenUserStories writes unfinished user stories in response as JSON
func GetOpenUserStories(writer http.ResponseWriter, recorder *http.Request) {
	text := ds.GetOpenUserStories()
	jason, _ := json.Marshal(text)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jason)
}

func main() {
	// Registering handlers for router in http package
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("HTTP server started..."))
	})
	// Start using default router
	http.ListenAndServe(":8080", nil)
}
