package model

import (
	"time"

	r "github.com/dancannon/gorethink"
)

// InsertYuQueRecord -
func InsertYuQueRecord(body string, actionType string, updateAt time.Time, user string, session *r.Session) error {
	var data = map[string]interface{}{
		"ActionType": actionType,
		"Body":       body,
		"UpdateAt":   updateAt,
		"User":       user,
	}

	_, err := r.DB("yuque").Table("yuque").Insert(data).RunWrite(session)

	return err
}

// DocDetailSerializer yuque document details
type DocDetailSerializer struct {
	Data struct {
		ID     int    `json:"id"`
		Slug   string `json:"slug"`
		Title  string `json:"title"`
		BookID int    `json:"book_id"`
		Book   struct {
			ID               int       `json:"id"`
			Type             string    `json:"type"`
			Slug             string    `json:"slug"`
			Name             string    `json:"name"`
			UserID           int       `json:"user_id"`
			Description      string    `json:"description"`
			CreatorID        int       `json:"creator_id"`
			Public           int       `json:"public"`
			ItemsCount       int       `json:"items_count"`
			LikesCount       int       `json:"likes_count"`
			WatchesCount     int       `json:"watches_count"`
			ContentUpdatedAt time.Time `json:"content_updated_at"`
			UpdatedAt        time.Time `json:"updated_at"`
			CreatedAt        time.Time `json:"created_at"`
			User             string    `json:"user"`
		} `json:"book"`
		UserID int `json:"user_id"`
		User   struct {
			ID               int       `json:"id"`
			Type             string    `json:"type"`
			Login            string    `json:"login"`
			Name             string    `json:"name"`
			Description      string    `json:"description"`
			AvatarURL        string    `json:"avatar_url"`
			LargeAvatarURL   string    `json:"large_avatar_url"`
			MediumAvatarURL  string    `json:"medium_avatar_url"`
			SmallAvatarURL   string    `json:"small_avatar_url"`
			BooksCount       int       `json:"books_count"`
			PublicBooksCount int       `json:"public_books_count"`
			FollowersCount   int       `json:"followers_count"`
			UpdatedAt        time.Time `json:"updated_at"`
			CreatedAt        time.Time `json:"created_at"`
		} `json:"user"`
		Format           string    `json:"format"`
		Body             string    `json:"body"`
		BodyDraft        string    `json:"body_draft"`
		BodyHTML         string    `json:"body_html"`
		Public           int       `json:"public"`
		Status           int       `json:"status"`
		LikesCount       int       `json:"likes_count"`
		CommentsCount    int       `json:"comments_count"`
		ContentUpdatedAt time.Time `json:"content_updated_at"`
		DeletedAt        time.Time `json:"deleted_at"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		PublishedAt      time.Time `json:"published_at"`
		FirstPublishedAt time.Time `json:"first_published_at"`
		WordCount        int       `json:"word_count"`
		ActionType       string    `json:"action_type"`
		Publish          bool      `json:"publish"`
		Path             string    `json:"path"`
	} `json:"data"`
}
