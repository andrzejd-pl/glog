package repository

import (
	"database/sql"
	"math/big"
	"time"
)

const (
	getAllQuery     string = "select p.PostUUID, p.PostTitle, p.PostContent, p.PostInsertTimestamp from Posts p WHERE p.PostId IN (SELECT MAX(PostId) FROM Posts GROUP BY PostUUID);"
	timestampLayout string = "2006-01-02 15:04:05"
)

type Post struct {
	Uuid           big.Int
	Title, Content string
	Timestamp      time.Time
}

type PostRepository interface {
	GetAllPosts() (*[]Post, error)
}

type mysqlPostRepository struct {
	mysqlConnection *sql.DB
}

func NewMysqlPostRepository(connection *sql.DB) PostRepository {
	return &mysqlPostRepository{mysqlConnection: connection}
}

func (m *mysqlPostRepository) GetAllPosts() (*[]Post, error) {
	rows, err := m.mysqlConnection.Query(getAllQuery)

	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var uuid, title, content, timestampValue string
		err = rows.Scan(&uuid, &title, &content, &timestampValue)
		if err != nil {
			return &posts, err
		}

		timestamp, err := time.Parse(timestampLayout, timestampValue)

		if err != nil {
			return &posts, err
		}

		post := Post{Uuid: big.Int{}, Title: title, Content: content, Timestamp: timestamp}
		post.Uuid.SetString(uuid, 10)

		posts = append(posts, post)
	}

	return &posts, nil
}
