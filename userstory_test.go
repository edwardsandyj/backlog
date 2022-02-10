package backlog

import (
	"reflect"
	"testing"
)

func TestGetOpenUserStories(t *testing.T) {
	t.Log("GetOpenUserStories")

	ds := Datastore{
		userstories: []UserStory{
			{1, "Find Airbnbs", true},
			{2, "Get car repaired", false},
		},
	}

	want := []UserStory{ds.userstories[1]}

	t.Log("want the stories which still need to be closed")

	if ans := ds.GetOpenUserStories(); !reflect.DeepEqual(ans, want) {
		t.Errorf("Got %v want %v", ans, want)
	}
}

var casesTestSaveUserStory = []struct {
	name  string
	ds    *Datastore
	story UserStory
	want  []UserStory
	err   error
}{
	{
		name:  "save new story in datastore",
		ds:    &Datastore{},
		story: UserStory{Description: "Find Airbnbs"},
		want: []UserStory{
			{1, "Find Airbnbs", false},
		},
	},
	{
		name: "update existing story in datastore",
		ds: &Datastore{
			userstories: []UserStory{
				{1, "Find Airbnbs", false},
			},
		},
		story: UserStory{1, "Find Airbnbs", true},
		want: []UserStory{
			{1, "Find Airbnbs", true},
		},
	},
	{
		name:  "throw error when story ID does not exist",
		ds:    &Datastore{},
		story: UserStory{1, "Find Airbnbs", true},
		err:   ErrUserStoryNotFound,
	},
}

func TestSaveUserStory(t *testing.T) {
	t.Log("SaveUserStory")
	for _, testcase := range casesTestSaveUserStory {
		t.Log(testcase.name)
		testcase.ds.SaveUserStory(testcase.story)
		if !reflect.DeepEqual(testcase.ds.userstories, testcase.want) {
			t.Errorf("=> Got %v want %v", testcase.ds.userstories, testcase.want)
		}
	}
}
