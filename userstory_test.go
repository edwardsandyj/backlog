package backlog

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var casesTestAddUserStory = []struct {
	name     string
	saveFunc func(story UserStory) error
	desc     []byte
	want     int
}{
	{
		name: "add new user story from JSON",
		desc: []byte(`{"Description":"Save route map offline."}`),
		want: http.StatusCreated,
	},
	{
		name: "return bad request if JSON decoding error",
		desc: []byte(""),
		want: http.StatusBadRequest,
	},
	{
		name: "return bad request if Stories.SaveUserStory returns error",
		saveFunc: func(story UserStory) error {
			return errors.New("error adding to Stories")
		},
		desc: []byte(`{"Description":"Save route map offline."}`),
		want: http.StatusBadRequest,
	},
	{
		name: "return bad request if UserStory description is empty",
		desc: []byte(`{"Description":""}`),
		want: http.StatusBadRequest,
	},
}

func TestAddUserStory(t *testing.T) {
	t.Log("AddUserStory")
	for _, testcase := range casesTestAddUserStory {
		t.Log(testcase.name)
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, "/userstories", bytes.NewBuffer(testcase.desc))
		defer func() { s = &Stories{} }()
		s = &mockedBacklog{
			SaveUserStoryFunc: testcase.saveFunc,
		}
		AddUserStory(recorder, request)
		if recorder.Code != testcase.want {
			t.Errorf("Output --> Got %d want %d", recorder.Code, testcase.want)
		}
	}

}

var casesTestUpdateUserStory = []struct {
	name     string
	saveFunc func(story UserStory) error
	desc     []byte
	want     int
}{
	{
		name: "return OK if successful",
		desc: []byte(`{"ID":1, "Description":"Save route map offline", "Closed":true }`),
		want: http.StatusOK,
	},
	{
		name: "return bad request if JSON decoding error",
		desc: []byte(""),
		want: http.StatusBadRequest,
	},
	{
		name: "return bad request if Stories.SaveUserStory returns error",
		saveFunc: func(story UserStory) error {
			return errors.New("error updating in Stories")
		},
		desc: []byte(`{"ID":1, "Description":"Save route map offline", "Closed":true}`),
		want: http.StatusBadRequest,
	},
	{
		name: "return bad request if UserStory description is empty",
		desc: []byte(`{"Description":""}`),
		want: http.StatusBadRequest,
	},
}

func TestUpdateUserStory(t *testing.T) {
	t.Log("UpdateUserStory")
	for _, testcase := range casesTestUpdateUserStory {
		t.Logf(testcase.name)
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPut, "/userstories/1", bytes.NewBuffer(testcase.desc))
		defer func() { s = &Stories{} }()
		s = &mockedBacklog{
			SaveUserStoryFunc: testcase.saveFunc,
		}
		UpdateUserStory(recorder, request)
		if recorder.Code != testcase.want {
			t.Errorf("Output --> Got %d want %d", recorder.Code, testcase.want)
		}
	}
}

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
