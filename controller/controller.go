package main

import (
	"log"
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

	yuqueSess, err := model.CreateYuQueTable()
	if err == nil {
		log.Fatal("failed")
	}

	yuque := &Session{yuqueSess}

	githubSess, err := model.CreateGitTable()
	if err == nil {
		log.Fatal("failed")
	}

	git := &Session{githubSess}

	router.POST("/GitHub/webhook", git.githubHandler)
	router.POST("/yuque/webhook", yuque.yuqueHandler)
	router.POST("/GitHub/select", git.selectHandler)
	router.POST("yuque/select", yuque.selectHandler)

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

	err = model.InsertYuQueRecord(yqhook.Data.Body, yqhook.Data.ActionType, yqhook.Data.UpdatedAt, yqhook.Data.User.Name, s.session)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"yuque": "1024"})
}

func (s Session) selectHandler(c *gin.Context) {
	var github struct {
		DBName    string `json:"dbname"`
		TableName string `json:"tablename"`
		Field     string `json:"field"`
		Value     string `json:"value"`
		Update    string `json:"update"`
	}

	err := c.ShouldBind(&github)
	if err != nil {
		c.Error(err)
		return
	}
	/*
		up, err := model.UpdateRecord(s.session, github.DBName, github.TableName, github.Field, github.Update)
		if err != nil {
			c.Error(err)
			return
		}*/

	all, err := model.SelectRecord(s.session, github.DBName, github.TableName, github.Field, github.Value)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"all": all})
}
