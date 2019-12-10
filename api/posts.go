package api

import (
	"github.com/andrzejd-pl/glog/repository"
	"log"
	"net/http"
)

func GetAllPosts(repository repository.PostRepository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if repository == nil {
			http.Error(writer, ResourceNotFound, http.StatusInternalServerError)
			log.Println("post repository is nil pointer")
			return
		}

		posts, err := repository.GetAllPosts()

		if err != nil {
			http.Error(writer, DatabaseConnectionError, http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		jsonResponse, err := posts.MarshalJSON()

		if err != nil {
			http.Error(writer, JsonMarshalError, http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(jsonResponse)

		if err != nil {
			log.Println(err.Error())
		}
	}
}
