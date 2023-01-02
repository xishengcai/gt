package gt

import (
	"testing"
)

type gitRepo []struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	AvatarURL   string `json:"avatar_url"`
	ReposURL    string `json:"repos_url"`
	EventsURL   string `json:"events_url"`
	MembersURL  string `json:"members_url"`
	Description string `json:"description"`
	FollowCount int    `json:"follow_count"`
}

func TestGet(t *testing.T) {
	result := gitRepo{}

	url := "https://gitee.com/api/v5/user/orgs?access_token=611895eac337ae09a08e982d8c744495&page=1&per_page=20"
	err := NewDefaultClient().
		SetURL(url).
		Get().
		Do().
		InTo(&result, JSON)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", result)
}
