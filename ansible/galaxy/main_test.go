package galaxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRoleURL(t *testing.T) {
	var cases = []struct {
		user    string
		role    string
		expURL  string
		expUser string
		expRepo string
	}{
		{
			user:    "geerlingguy",
			role:    "pip",
			expURL:  "https://github.com/geerlingguy/ansible-role-pip",
			expUser: "geerlingguy",
			expRepo: "ansible-role-pip",
		},
		{
			user:    "williamyeh",
			role:    "oracle-java",
			expURL:  "",
			expUser: "",
			expRepo: "",
		},
		{
			user:    "abessifi",
			role:    "java",
			expURL:  "https://github.com/abessifi/ansible-java",
			expUser: "abessifi",
			expRepo: "ansible-java",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.expURL, func(t *testing.T) {
			t.Parallel()
			t.Log(tt.expURL)
			retURL, retUser, retRole := FindRoleURL(tt.user, tt.role)
			assert.Equal(t, retURL, tt.expURL)
			assert.Equal(t, retUser, tt.expUser)
			assert.Equal(t, retRole, tt.expRepo)
		})
	}
}
