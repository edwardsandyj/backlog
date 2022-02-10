package backlog

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetOpenUserStories(t *testing.T) {
	t.Log("GetOpenUserStories")
	t.Log("want OpenUserStories data as JSON")
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/userstories/open", nil)
	defer func() { s = &Stories{} }()
	s = &mockedBacklog{}
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
	s     *Stories
	story UserStory
	want  []UserStory
	err   error
}{
	{
		name:  "save new story in datastore",
		s:     &Stories{},
		story: UserStory{Description: "Find Airbnbs"},
		want: []UserStory{
			{1, "Find Airbnbs", false},
		},
	},
	{
		name: "update existing story in datastore",
		s: &Stories{
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
		s:     &Stories{},
		story: UserStory{1, "Find Airbnbs", true},
		err:   ErrUserStoryNotFound,
	},
}

func TestSaveUserStory(t *testing.T) {
	t.Log("SaveUserStory")
	for _, testcase := range casesTestSaveUserStory {
		t.Log(testcase.name)
		testcase.s.SaveUserStory(testcase.story)
		if !reflect.DeepEqual(testcase.s.userstories, testcase.want) {
			t.Errorf("=> Got %v want %v", testcase.s.userstories, testcase.want)
		}
	}
}
