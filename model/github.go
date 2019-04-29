package model

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// CreateGitTable -
func CreateGitTable(session *r.Session) error {
	_, err := r.DB("github").TableCreate("GitHub").RunWrite(session)
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

// SelectRecord -
func SelectRecord(session *r.Session, id, tablename string) error {
	ew := r.DB("test").Table(tablename).Get(id).Values()
	fmt.Println(ew, "2")

	res, err := r.DB("test").Table(tablename).Get(id).Run(session)
	fmt.Println(res, "2")
	return err
}
