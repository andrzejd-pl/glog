package api

import (
	"log"
	"net/http"
)

const (
	DatabaseConnectionError = "database connection error"
	JsonMarshalError        = "data marshal error"
	ResourceNotFound        = "resource not found"
)

type Endpoints map[string]http.HandlerFunc

func MakeHandler(serverMux *http.ServeMux, address string, endpoints Endpoints) {
	serverMux.HandleFunc(address, func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.RemoteAddr, request.Method, request.URL)
		action, ok := endpoints[request.Method]

		if !ok {
			http.NotFound(writer, request)
			return
		}

		action(writer, request)
	})
}

func MakeContentType(handler http.HandlerFunc, contentType string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", contentType)
		handler(writer, request)
	}
}
