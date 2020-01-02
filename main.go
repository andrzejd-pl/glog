package main

import (
	"database/sql"
	"github.com/andrzejd-pl/glog/api"
	"github.com/andrzejd-pl/glog/repository"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(database)/glog")
	CheckIfError(err)
	defer func() { _ = db.Close() }()
	repo := repository.NewMysqlPostRepository(db)

	api.MakeHandler(http.DefaultServeMux, "/api/posts", api.Endpoints{
		"GET": api.MakeContentType(api.GetAllPosts(repo), "application/json"),
	})
	log.Fatalln(http.ListenAndServe(":80", http.DefaultServeMux))
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
