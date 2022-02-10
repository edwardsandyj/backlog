package backlog

import "errors"

var ErrUserStoryNotFound = errors.New("User Story was not found")

// User story that is still open or closed
type UserStory struct {
	ID          int    // story identifier
	Description string // story description
	Closed      bool   // story completion status
}

// Datastore holds all user stories saved in the application
type Datastore struct {
	userstories []UserStory
	lastID      int // mark IDs that are already used in userstories
}

// SaveUserStory creates new or updates an existing UserStory in Datastore
func (ds *Datastore) SaveUserStory(story UserStory) error {
	if story.ID == 0 {
		ds.lastID++
		story.ID = ds.lastID
		ds.userstories = append(ds.userstories, story)
		return nil
	}

	for i, t := range ds.userstories {
		if t.ID == story.ID {
			ds.userstories[i] = story
			return nil
		}
	}
	return ErrUserStoryNotFound
}

// GetOpenUserStories returns all unfinished user stories in Datastore
func (ds *Datastore) GetOpenUserStories() []UserStory {
	var openstories []UserStory
	for _, story := range ds.userstories {
		if !story.Closed {
			openstories = append(openstories, story)
		}
	}
	return openstories
}
