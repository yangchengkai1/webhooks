package controller

import (
	"log"
	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"

	r "github.com/dancannon/gorethink"
	"github.com/gin-gonic/gin"
	model "github.com/yangchengkai1/webhooks/model/rethinkdb"
)

// Session -
type Session struct {
	ys *r.Session
	gs *r.Session
}

// RegisterRouter -
func RegisterRouter(router gin.IRouter) {
	var ss *Session

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost",
	})
	if err != nil {
		log.Fatal(err)
	}

	ys, err := model.Create("yuque", "yuque", session)
	if err != nil {
		log.Fatal(err)
	}

	gs, err := model.Create("github", "github", session)
	if err != nil {
		log.Fatal(err)
	}

	ss = &Session{ys: ys, gs: gs}

	router.POST("/github/webhook", ss.githubStore)
	router.POST("/yuque/webhook", ss.yuqueStore)
	router.GET("/select/value", ss.selectHandler)
	router.POST("/select/field", ss.selectItems)
	router.GET("/select/all", ss.selectAllHandler)
	router.POST("/delete", ss.deleteHandler)
	router.POST("/update", ss.updateHandler)
	router.POST("/filter", ss.filterHandler)
}

func (s Session) githubStore(c *gin.Context) {
	hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecret"))

	payload, err := hook.Parse(c.Request, github.PushEvent, github.PullRequestEvent)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		err = model.InsertPushRecord(push, s.gs)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
			return
		}
	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		err = model.InsertReleaseRecord(release, s.gs)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
			return
		}
	case github.RepositoryPayload:
		repo := payload.(github.RepositoryPayload)
		err = model.InsertRepoRecord(repo, s.gs)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (s Session) yuqueStore(c *gin.Context) {
	var yqhook model.DocDetailSerializer

	err := c.ShouldBind(&yqhook)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.InsertYuQueRecord(yqhook.Data.Body, yqhook.Data.ActionType, yqhook.Data.UpdatedAt, yqhook.Data.User.Name, s.ys)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (s Session) selectHandler(c *gin.Context) {
	var (
		term struct {
			DBName    string `json:"db_name"    binding:"required"`
			TableName string `json:"table_name" binding:"required"`
			Field     string `json:"field"      binding:"required"`
			Value     string `json:"value"      binding:"required"`
		}

		session *r.Session
	)

	err := c.ShouldBind(&term)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch term.TableName {
	case "yuque":
		session = s.ys
	case "github":
		session = s.gs
	}

	all, err := model.SelectRecord(session, term.DBName, term.TableName, term.Field, term.Value)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all": all})
}

func (s Session) selectItems(c *gin.Context) {
	var (
		term struct {
			DBName    string   `json:"db_name"    binding:"required"`
			TableName string   `json:"table_name" binding:"required"`
			Field     []string `json:"field"      binding:"required"`
		}

		session *r.Session
	)

	err := c.ShouldBind(&term)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch term.TableName {
	case "yuque":
		session = s.ys
	case "github":
		session = s.gs
	}

	all, err := model.SelectItems(session, term.DBName, term.TableName, term.Field)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all": all})
}

func (s Session) selectAllHandler(c *gin.Context) {
	var (
		term struct {
			DBName    string `json:"db_name"    binding:"required"`
			TableName string `json:"table_name" binding:"required"`
		}

		session *r.Session
	)

	err := c.ShouldBind(&term)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch term.TableName {
	case "yuque":
		session = s.ys
	case "github":
		session = s.gs
	}

	all, err := model.AllRecord(session, term.DBName, term.TableName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all": all})
}

func (s Session) deleteHandler(c *gin.Context) {
	var (
		term struct {
			DBName    string `json:"db_name"    binding:"required"`
			TableName string `json:"table_name" binding:"required"`
			Field     string `json:"field"      binding:"required"`
			Value     string `json:"value"      binding:"required"`
		}

		session *r.Session
	)

	err := c.ShouldBind(&term)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch term.TableName {
	case "yuque":
		session = s.ys
	case "github":
		session = s.gs
	}

	err = model.DelateRecord(session, term.DBName, term.TableName, term.Field, term.Value)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (s Session) updateHandler(c *gin.Context) {
	var (
		term struct {
			DBName    string `json:"db_name"    binding:"required"`
			TableName string `json:"table_name" binding:"required"`
			Field     string `json:"field"      binding:"required"`
			Value     string `json:"value"      binding:"required"`
		}

		session *r.Session
	)

	err := c.ShouldBind(&term)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch term.TableName {
	case "yuque":
		session = s.ys
	case "github":
		session = s.gs
	}

	resp, err := model.UpdateRecord(session, term.DBName, term.TableName, term.Field, term.Value)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "resp": resp})
}

func (s Session) filterHandler(c *gin.Context) {
	var (
		term struct {
			DBName    string   `json:"db_name"     binding:"required"`
			TableName string   `json:"table_name"  binding:"required"`
			Filter    []string `json:"filter"      binding:"required"`
		}

		session *r.Session
	)

	err := c.ShouldBind(&term)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch term.TableName {
	case "yuque":
		session = s.ys
	case "github":
		session = s.gs
	}

	resp, err := model.Filter(session, term.DBName, term.TableName, term.Filter)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "resp": resp})
}
