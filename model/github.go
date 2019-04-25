package model

import (
	r "github.com/dancannon/gorethink"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// CreateGitTable -
func CreateGitTable(session *r.Session) error {
	_, err := r.DB("test").TableCreate("GitHub").RunWrite(session)
	return err
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
