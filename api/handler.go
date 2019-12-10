package api

import "net/http"

type Endpoints map[string]http.HandlerFunc

func MakeHandler(serverMux *http.ServeMux, address string, endpoints Endpoints) {
	serverMux.HandleFunc(address, func(writer http.ResponseWriter, request *http.Request) {
		action, ok := endpoints[request.Method]

		if !ok {
			http.NotFound(writer, request)
			return
		}

		action(writer, request)
	})
}
