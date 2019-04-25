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

	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost",
		Database: "test",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	s := &Session{session}

	model.CreateGitTable(session)
	model.CreateYuqueTable(session)

	router.POST("/GitHub/webhook", s.handlerGitHub)
	router.POST("/yuque/webhook", s.handlerYuque)
	router.Run(":8080")
}

func (s Session) handlerGitHub(c *gin.Context) {
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

	c.JSON(200, gin.H{"yuque": "1024"})
}

func (s Session) handlerYuque(c *gin.Context) {
	err := c.ShouldBind(&yqhook)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.InsertYuqueRecord(yqhook.Data.Body, yqhook.Data.ActionType, yqhook.Data.UpdatedAt, s.session)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"yuque": "1024"})
}
