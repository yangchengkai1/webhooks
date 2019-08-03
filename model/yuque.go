package model

import (
	"fmt"

	r "github.com/dancannon/gorethink"
)

// CreateYuQueTable -
func CreateYuQueTable() (*r.Session, error) {
	yuqueSess, err := r.Connect(r.ConnectOpts{
		Address:  "localhost",
		Database: "yuque",
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = CheckTable(yuqueSess, "yuque", "yuque")
	if err == errTable {
		return yuqueSess, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = r.DB("yuque").TableCreate("yuque").RunWrite(yuqueSess)

	return yuqueSess, err
}

// InsertYuQueRecord -
func InsertYuQueRecord(body string, actionType string, updateAt string, user string, session *r.Session) error {
	var data = map[string]interface{}{
		"ActionType": actionType,
		"Body":       body,
		"UpdateAt":   updateAt,
		"User":       user,
	}

	_, err := r.Table("yuque").Insert(data).RunWrite(session)

	return err
}

// DocDetailSerializer yuque document details
type DocDetailSerializer struct {
	ID     int64  `json:"id"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	BookID int64  `json:"book_id"`
	Book   struct {
		ID               int64  `json:"id"`
		Type             string `json:"type"`
		Slug             string `json:"slug"`
		Name             string `json:"name"`
		UserID           int64  `json:"user_id"`
		Description      string `json:"description"`
		CreatorID        int64  `json:"creator_id"`
		Public           int    `json:"public"`
		ItemsCount       int64  `json:"items_count"`
		LikesCount       int64  `json:"likes_count"`
		WatchesCount     int64  `json:"watches_count"`
		ContentUpdatedAt string `json:"content_updated_at"`
		UpdatedAt        string `json:"updated_at"`
		CreatedAt        string `json:"created_at"`
		User             string `json:"user"`
	} `json:"book"`
	UserID int64 `json:"user_id"`
	User   struct {
		ID               int64  `json:"id"`
		Type             string `json:"type"`
		Login            string `json:"login"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		AvatarURL        string `json:"avatar_url"`
		LargeAvatarURL   string `json:"large_avatar_url"`
		MediumAvatarURL  string `json:"medium_avatar_url"`
		SmallAvatarURL   string `json:"small_avatar_url"`
		BooksCount       int64  `json:"books_count"`
		PublicBooksCount int64  `json:"public_books_count"`
		FollowersCount   int64  `json:"followers_count"`
		UpdatedAt        string `json:"updated_at"`
		CreatedAt        string `json:"created_at"`
	} `json:"user"`
	Format           string `json:"format"`
	Body             string `json:"body"`
	BodyDraft        string `json:"body_draft"`
	BodyHTML         string `json:"body_html"`
	Public           int    `json:"public"`
	Status           int    `json:"status"`
	LikesCount       int64  `json:"likes_count"`
	CommentsCount    int64  `json:"comments_count"`
	ContentUpdatedAt string `json:"content_updated_at"`
	DeletedAt        string `json:"deleted_at"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	PublishedAt      string `json:"published_at"`
	FirstPublishedAt string `json:"first_published_at"`
	WordCount        int64  `json:"word_count"`
	ActionType       string `json:"action_type"`
	Publish          bool   `json:"publish"`
	Path             string `json:"path"`
}
