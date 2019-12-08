package main

import (
	"database/sql"
	"fmt"
	"github.com/andrzejd-pl/glog/repository"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:my-secret-pw@/glog")
	CheckIfError(err)
	defer func() { _ = db.Close() }()
	repo := repository.NewMysqlPostRepository(db)
	posts, err := repo.GetAllPosts()
	CheckIfError(err)

	for _, post := range *posts {
		fmt.Println("Title: ", post.Title)
		fmt.Println("Modify: ", post.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println("Content: ", post.Content)
		fmt.Println("----------------------")
	}
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
