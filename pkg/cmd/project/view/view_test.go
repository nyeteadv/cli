package view

import (
	"bytes"
	"os"
	"testing"

	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestNewCmdview(t *testing.T) {
	tests := []struct {
		name        string
		cli         string
		wants       viewOpts
		wantsErr    bool
		wantsErrMsg string
	}{
		{
			name:        "user-and-org",
			cli:         "--user monalisa --org github",
			wantsErr:    true,
			wantsErrMsg: "only one of `--user` or `--org` may be used",
		},
		{
			name:        "not-a-number",
			cli:         "x",
			wantsErr:    true,
			wantsErrMsg: "invalid number: x",
		},
		{
			name: "number",
			cli:  "123",
			wants: viewOpts{
				number: 123,
			},
		},
		{
			name: "user",
			cli:  "--user monalisa",
			wants: viewOpts{
				userOwner: "monalisa",
			},
		},
		{
			name: "org",
			cli:  "--org github",
			wants: viewOpts{
				orgOwner: "github",
			},
		},
		{
			name: "web",
			cli:  "--web",
			wants: viewOpts{
				web: true,
			},
		},
		{
			name: "json",
			cli:  "--format json",
			wants: viewOpts{
				format: "json",
			},
		},
	}

	os.Setenv("GH_TOKEN", "auth-token")
	defer os.Unsetenv("GH_TOKEN")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()
			f := &cmdutil.Factory{
				IOStreams: ios,
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts viewOpts
			cmd := NewCmdView(f, func(config viewConfig) error {
				gotOpts = config.opts
				return nil
			})

			cmd.SetArgs(argv)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantsErrMsg, err.Error())
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.number, gotOpts.number)
			assert.Equal(t, tt.wants.userOwner, gotOpts.userOwner)
			assert.Equal(t, tt.wants.orgOwner, gotOpts.orgOwner)
			assert.Equal(t, tt.wants.format, gotOpts.format)
			assert.Equal(t, tt.wants.web, gotOpts.web)
		})
	}
}

func TestBuildURLViewer(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/graphql").
		Reply(200).
		JSON(`
			{"data":
				{"viewer":
					{
						"login":"theviewer"
					}
				}
			}
		`)

	client, err := api.NewGraphQLClient(api.ClientOptions{AuthToken: "token"})
	assert.NoError(t, err)

	url, err := buildURL(viewConfig{
		opts: viewOpts{
			number:    1,
			userOwner: "@me",
		},
		client: client,
	})
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/users/theviewer/projects/1", url)
}

func TestBuildURLUser(t *testing.T) {
	url, err := buildURL(viewConfig{
		opts: viewOpts{
			userOwner: "monalisa",
			number:    1,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/users/monalisa/projects/1", url)
}

func TestBuildURLOrg(t *testing.T) {
	url, err := buildURL(viewConfig{
		opts: viewOpts{
			orgOwner: "github",
			number:   1,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/orgs/github/projects/1", url)
}

func TestRunView_User(t *testing.T) {
	defer gock.Off()

	// get user ID
	gock.New("https://api.github.com").
		Post("/graphql").
		MatchType("json").
		JSON(map[string]interface{}{
			"query": "query UserLogin.*",
			"variables": map[string]interface{}{
				"login": "monalisa",
			},
		}).
		Reply(200).
		JSON(map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"id": "an ID",
				},
			},
		})

	gock.New("https://api.github.com").
		Post("/graphql").
		Reply(200).
		JSON(`
			{"data":
				{"user":
					{
						"login":"monalisa",
						"projectV2": {
							"number": 1,
							"items": {
								"totalCount": 10
							},
							"readme": null,
							"fields": {
								"nodes": [
									{
										"name": "Title"
									}
								]
							}
						}
					}
				}
			}
		`)

	client, err := api.NewGraphQLClient(api.ClientOptions{AuthToken: "token"})
	assert.NoError(t, err)

	ios, _, _, _ := iostreams.Test()
	config := viewConfig{
		tp: tableprinter.New(ios),
		opts: viewOpts{
			userOwner: "monalisa",
			number:    1,
		},
		client: client,
	}

	err = runView(config)
	assert.NoError(t, err)

}

func TestRunView_Viewer(t *testing.T) {
	defer gock.Off()

	// get viewer ID
	gock.New("https://api.github.com").
		Post("/graphql").
		MatchType("json").
		JSON(map[string]interface{}{
			"query": "query ViewerLogin.*",
		}).
		Reply(200).
		JSON(map[string]interface{}{
			"data": map[string]interface{}{
				"viewer": map[string]interface{}{
					"id": "an ID",
				},
			},
		})

	gock.New("https://api.github.com").
		Post("/graphql").
		Reply(200).
		JSON(`
			{"data":
				{"viewer":
					{
						"login":"monalisa",
						"projectV2": {
							"number": 1,
							"items": {
								"totalCount": 10
							},
							"readme": null,
							"fields": {
								"nodes": [
									{
										"name": "Title"
									}
								]
							}
						}
					}
				}
			}
		`)

	client, err := api.NewGraphQLClient(api.ClientOptions{AuthToken: "token"})
	assert.NoError(t, err)

	ios, _, _, _ := iostreams.Test()
	config := viewConfig{
		tp: tableprinter.New(ios),
		opts: viewOpts{
			userOwner: "@me",
			number:    1,
		},
		client: client,
	}

	err = runView(config)
	assert.NoError(t, err)
}

func TestRunView_Org(t *testing.T) {
	defer gock.Off()

	// get org ID
	gock.New("https://api.github.com").
		Post("/graphql").
		MatchType("json").
		JSON(map[string]interface{}{
			"query": "query OrgLogin.*",
			"variables": map[string]interface{}{
				"login": "github",
			},
		}).
		Reply(200).
		JSON(map[string]interface{}{
			"data": map[string]interface{}{
				"organization": map[string]interface{}{
					"id": "an ID",
				},
			},
		})

	gock.New("https://api.github.com").
		Post("/graphql").
		Reply(200).
		JSON(`
			{"data":
				{"organization":
					{
						"login":"monalisa",
						"projectV2": {
							"number": 1,
							"items": {
								"totalCount": 10
							},
							"readme": null,
							"fields": {
								"nodes": [
									{
										"name": "Title"
									}
								]
							}
						}
					}
				}
			}
		`)

	client, err := api.NewGraphQLClient(api.ClientOptions{AuthToken: "token"})
	assert.NoError(t, err)

	ios, _, _, _ := iostreams.Test()
	config := viewConfig{
		tp: tableprinter.New(ios),
		opts: viewOpts{
			orgOwner: "github",
			number:   1,
		},
		client: client,
	}

	err = runView(config)
	assert.NoError(t, err)
}

func TestRunViewWeb(t *testing.T) {
	buf := bytes.Buffer{}
	config := viewConfig{
		opts: viewOpts{
			userOwner: "monalisa",
			web:       true,
			number:    8,
		},
		URLOpener: func(url string) error {
			buf.WriteString(url)
			return nil
		},
	}

	err := runView(config)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/users/monalisa/projects/8", buf.String())
}
