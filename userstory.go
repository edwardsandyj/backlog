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

// Datastore holds a slice of all user stories kept by the application
type Stories struct {
	userstories []UserStory
	lastID      int // mark IDs that are already used in UserStories
}

// Backlog defines the services on a Stories object
type Backlog interface {
	GetOpenUserStories() []UserStory
}

var s Backlog = &Stories{}

// *Mock* GetOpenUserStories call that implements interface Backlog
type mockedBacklog struct{}

func (mb *mockedBacklog) GetOpenUserStories() []UserStory {
	return []UserStory{
		{1, "Find Airbnbs", false},
		{2, "Get car repaired", false},
	}
}

// SaveUserStory creates new or updates an existing UserStory in Datastore
func (s *Stories) SaveUserStory(story UserStory) error {
	if story.ID == 0 {
		s.lastID++
		story.ID = s.lastID
		s.userstories = append(s.userstories, story)
		return nil
	}

	for i, t := range s.userstories {
		if t.ID == story.ID {
			s.userstories[i] = story
			return nil
		}
	}
	return ErrUserStoryNotFound
}

// GetOpenUserStories returns all unfinished user stories in Datastore
func (s *Stories) GetOpenUserStories() []UserStory {
	var openstories []UserStory
	for _, story := range s.userstories {
		if !story.Closed {
			openstories = append(openstories, story)
		}
	}
	return openstories
}

// GetOpenUserStories writes unfinished user stories in response as JSON
func GetOpenUserStories(writer http.ResponseWriter, recorder *http.Request) {
	text := s.GetOpenUserStories()
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
