package model

import (
	r "github.com/dancannon/gorethink"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// InsertPushRecord -
func InsertPushRecord(push github.PushPayload, session *r.Session) error {
	var data = map[string]interface{}{
		"repo_owner": push.Repository.Owner.Login,
		"repo_name":  push.Repository.FullName,
		"message":    push.HeadCommit.Message,
		"URL":        push.HeadCommit.URL,
		"updated_at": push.Repository.UpdatedAt,
	}

	_, err := r.DB("github").Table("github").Insert(data).RunWrite(session)

	return err
}

// InsertReleaseRecord -
func InsertReleaseRecord(release github.ReleasePayload, session *r.Session) error {
	var data = map[string]interface{}{
		"action":       release.Action,
		"tag_name":     release.Release.TagName,
		"published_at": release.Release.PublishedAt,
		"repo_owner":   release.Repository.Owner.Login,
		"repo_name":    release.Repository.FullName,
		"sender":       release.Sender.Login,
	}

	_, err := r.DB("github").Table("github").Insert(data).RunWrite(session)

	return err
}

// InsertRepoRecord -
func InsertRepoRecord(repo github.RepositoryPayload, session *r.Session) error {
	var data = map[string]interface{}{
		"action":     repo.Action,
		"repo_owner": repo.Repository.Owner.Login,
		"repo_name":  repo.Repository.FullName,
		"created_at": repo.Repository.CreatedAt,
		"updated_at": repo.Repository.UpdatedAt,
		"pushed_at":  repo.Repository.PushedAt,
		"sender":     repo.Sender.Login,
	}

	_, err := r.DB("github").Table("github").Insert(data).RunWrite(session)

	return err
}
