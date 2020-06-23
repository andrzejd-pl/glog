package main

import (
	"database/sql"
	"flag"
	"github.com/andrzejd-pl/configload"
	"github.com/andrzejd-pl/glog/api"
	"github.com/andrzejd-pl/glog/configuration"
	"github.com/andrzejd-pl/glog/repository"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	configArgDescription = "a configuration file\n example configuration file: " +
		"https://github.com/andrzejd-pl/glog/blob/master/config.example.json"
)

func main() {
	configFilePtr := flag.String("config", "./config.json", configArgDescription)
	flag.Parse()

	config := configuration.Config{}
	cf := configload.NewJsonFileConfiguration(*configFilePtr, &config)
	CheckIfError(cf.LoadFromFile())

	if logFile := setLoggerOutput(&config); logFile != nil {
		defer func() { _ = logFile.Close() }()
	}

	db, err := sql.Open("mysql", config.DatabaseDsn)
	CheckIfError(err)
	defer func() { _ = db.Close() }()
	repo := repository.NewMysqlPostRepository(db)

	api.MakeHandler(http.DefaultServeMux, "/api/posts", api.Endpoints{
		"GET": api.MakeContentType(api.GetAllPosts(repo), "application/json"),
	})
	log.Fatalln(http.ListenAndServe(":80", http.DefaultServeMux))
}

func setLoggerOutput(config *configuration.Config) io.Closer {
	if config.Logger == nil {
		log.SetOutput(os.Stderr)

		return nil
	}

	if config.Logger.Type == "file" {
		file, err := os.OpenFile(config.Logger.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Panicln(err)
		}

		log.SetOutput(file)

		return file
	}

	log.Panicln("Undefined logger type")
	return nil
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
