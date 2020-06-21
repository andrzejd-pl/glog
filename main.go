package main

import (
	"database/sql"
	"flag"
	"github.com/andrzejd-pl/configuration_loader"
	"github.com/andrzejd-pl/glog/api"
	"github.com/andrzejd-pl/glog/configuration"
	"github.com/andrzejd-pl/glog/repository"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	configFilePtr := flag.String(
		"config",
		"./config.json",
		"a configuration file\n example configuration file: "+
			"https://github.com/andrzejd-pl/glog/blob/master/config.example.json",
	)
	flag.Parse()

	config := configuration.Config{}
	cf := configuration_loader.NewJsonFileConfiguration(*configFilePtr, &config)
	CheckIfError(cf.LoadFromFile())

	db, err := sql.Open("mysql", config.DatabaseDsn)
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
