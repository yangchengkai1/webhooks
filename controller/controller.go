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

	model.CreateGitTable(s.session)
	model.CreateYuQueTable(s.session)

	router.POST("/GitHub/webhook", s.githubHandler)
	router.POST("/yuque/webhook", s.yuqueHandler)

	router.POST("/update", s.updateByID)
	router.POST("/select", s.selectByID)
	router.POST("/delate", s.delateByID)

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

func (s Session) selectByID(c *gin.Context) {
	var github struct {
		ID        string `json:"id"`
		TableName string `json:"tablename"`
	}

	err := c.ShouldBind(&github)
	if err != nil {
		c.Error(err)
		return
	}

	err = model.SelectRecord(s.session, github.ID, github.TableName)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{github.TableName: "1024"})
}

func (s Session) delateByID(c *gin.Context) {
	var github struct {
		ID        string `json:"id"`
		TableName string `json:"tablename"`
	}

	err := c.ShouldBind(&github)
	if err != nil {
		c.Error(err)
		return
	}

	result, err := model.DelateRecord(s.session, github.ID, github.TableName)
	if err != nil {
		c.Error(err)
		return
	}

	fmt.Println(result)
	c.JSON(http.StatusOK, gin.H{github.TableName: "1024"})
}

func (s Session) updateByID(c *gin.Context) {
	var github struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Update    string `json:"update"`
		TableName string `json:"tablename"`
	}

	err := c.ShouldBind(&github)
	if err != nil {
		c.Error(err)
		return
	}

	err = model.UpdateRecord(s.session, github.ID, github.TableName, github.Title, github.Update)
	if err != nil {
		c.Error(err)
		return
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
