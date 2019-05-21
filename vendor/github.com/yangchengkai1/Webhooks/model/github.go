package model

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// CreateGitTable -
func CreateGitTable() (*r.Session, error) {
	githubSess, err := r.Connect(r.ConnectOpts{
		Address:  "localhost",
		Database: "github",
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = CheckTable(githubSess, "github", "GitHub")
	if err == errTable {
		return githubSess, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = r.DB("github").TableCreate("GitHub").RunWrite(githubSess)
	return githubSess, err
}

// InsertGitRecord -
func InsertGitRecord(push github.PushPayload, session *r.Session) error {
	var data = map[string]interface{}{
		"RepositoryOwner": push.Repository.Owner.Login,
		"FullName":        push.Repository.FullName,
		"Message":         push.HeadCommit.Message,
		"URL":             push.HeadCommit.URL,
		"UpdatedAt":       push.Repository.UpdatedAt,
	}

	_, err := r.Table("GitHub").Insert(data).RunWrite(session)

	return err
}
