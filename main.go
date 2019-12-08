package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const query string = "select p.PostTitle, p.PostContent, p.PostInsertTimestamp from Posts p WHERE p.PostId IN (SELECT MAX(PostId) FROM Posts GROUP BY PostUUID);"

func main() {
	db, err := sql.Open("mysql", "root:my-secret-pw@/glog")
	CheckIfError(err)
	defer db.Close()

	rows, err := db.Query(query)
	CheckIfError(err)
	for rows.Next() {
		var title, content, timestamp string
		err = rows.Scan(&title, &content, &timestamp)
		CheckIfError(err)
		fmt.Println(title + ": " + content)
	}
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
