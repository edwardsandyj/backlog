package backlog

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetOpenUserStories(t *testing.T) {
	t.Log("GetOpenUserStories")
	t.Log("want open user stories as JSON")
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/userstories/open", nil)
	ds = &Datastore{
		UserStories: []UserStory{
			{1, "Find Airbnbs", false},
			{2, "Get car repaired", false},
		},
	}
	GetOpenUserStories(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("Output --> Got %d want %d", recorder.Code, http.StatusOK)
	}
	want := "[{\"id\":1,\"description\":\"Find Airbnbs\",\"closed\":false},{\"id\":2,\"description\":\"Get car repaired\",\"closed\":false}]"

	if ans := recorder.Body.String(); ans != want {
		t.Errorf("Output --> Got %s want %s", ans, want)
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
			UserStories: []UserStory{
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
		if !reflect.DeepEqual(testcase.ds.UserStories, testcase.want) {
			t.Errorf("=> Got %v want %v", testcase.ds.UserStories, testcase.want)
		}
	}
}
