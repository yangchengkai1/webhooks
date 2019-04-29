package main

import (
	"fmt"
	"mytest/webhook/yuque"
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gin-gonic/gin"
	"github.com/yangchengkai1/Webhooks/model"
	"gopkg.in/go-playground/webhooks.v5/github"
)

var yqhook struct {
	Data yuque.DocDetailSerializer `json:"data"`
}

// Session -
type Session struct {
	session *r.Session
}

func main() {
	router := gin.Default()

	yqSession, err := r.Connect(r.ConnectOpts{
		Address:  "localhost",
		Database: "yuque",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	yuque := &Session{yqSession}

	gitSession, err := r.Connect(r.ConnectOpts{
		Address:  "localhost",
		Database: "github",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	git := &Session{gitSession}

	model.CreateGitTable(yuque.session)
	model.CreateYuQueTable(git.session)

	router.POST("/GitHub/webhook", git.githubHandler)
	router.POST("/yuque/webhook", yuque.yuqueHandler)

	router.Run(":8080")
}

func (s Session) githubHandler(c *gin.Context) {
	hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecret"))

	payload, err := hook.Parse(c.Request, github.PushEvent, github.PullRequestEvent)
	if err != nil {
		c.Error(err)
		return
	}

	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		err = model.InsertGitRecord(push, s.session)
		if err != nil {
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"GitHub": "1024"})
}

func (s Session) yuqueHandler(c *gin.Context) {
	err := c.ShouldBind(&yqhook)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.InsertYuQueRecord(yqhook.Data.Body, yqhook.Data.ActionType, yqhook.Data.UpdatedAt, s.session)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"yuque": "1024"})
}
