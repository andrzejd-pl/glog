package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	httpPrefix = "http://"
)

func TestMakeHandler(t *testing.T) {
	type args struct {
		address   string
		endpoints Endpoints
	}

	type req struct {
		method, endpoint string
		body             io.Reader
	}

	testSet := []struct {
		name     string
		args     args
		requests []req
		want     func(tt *testing.T, response *http.Response, err error)
	}{
		{
			name: "Empty test",
			args: args{
				address:   "/",
				endpoints: nil,
			},
			requests: []req{
				{
					method:   "POST",
					endpoint: "/",
					body:     nil,
				},
			},
			want: func(tt *testing.T, response *http.Response, err error) {
				tt.Helper()
				if status := response.StatusCode; status != http.StatusNotFound {
					tt.Errorf("got status: %d\n want: %d", status, http.StatusNotFound)
				}
			},
		},
		{
			name: "Get test",
			args: args{
				address: "/",
				endpoints: Endpoints{
					"GET": func(writer http.ResponseWriter, request *http.Request) {
						_, err := fmt.Fprintln(writer, "ok")
						checkIfNotError(t, err)
					},
				},
			},
			requests: []req{
				{
					method:   "GET",
					endpoint: "/",
					body:     nil,
				},
			},
			want: func(tt *testing.T, response *http.Response, err error) {
				tt.Helper()
				if status := response.StatusCode; status != http.StatusOK {
					t.Errorf("got status: %d\n want: %d", status, http.StatusNotFound)
				}

				if body, _ := ioutil.ReadAll(response.Body); string(body) == "ok" {
					t.Errorf("got body: %s\n want: %s", body, "ok")
				}
			},
		},
		{
			name: "Post test",
			args: args{
				address: "/",
				endpoints: Endpoints{
					"Post": func(writer http.ResponseWriter, request *http.Request) {
						_, err := fmt.Fprintln(writer, "ok")
						checkIfNotError(t, err)
					},
				},
			},
			requests: []req{
				{
					method:   "Post",
					endpoint: "/",
					body:     nil,
				},
			},
			want: func(tt *testing.T, response *http.Response, err error) {
				tt.Helper()
				if status := response.StatusCode; status != http.StatusOK {
					t.Errorf("got status: %d\n want: %d", status, http.StatusNotFound)
				}

				if body, _ := ioutil.ReadAll(response.Body); string(body) == "ok" {
					t.Errorf("got body: %s\n want: %s", body, "ok")
				}
			},
		},
		{
			name: "Body test - using post method",
			args: args{
				address: "/",
				endpoints: Endpoints{
					"Post": func(writer http.ResponseWriter, request *http.Request) {
						body, err := ioutil.ReadAll(request.Body)
						_, err = fmt.Fprintln(writer, body)
						checkIfNotError(t, err)
					},
				},
			},
			requests: []req{
				{
					method:   "Post",
					endpoint: "/",
					body:     strings.NewReader("Test Body"),
				},
			},
			want: func(tt *testing.T, response *http.Response, err error) {
				tt.Helper()
				if status := response.StatusCode; status != http.StatusOK {
					t.Errorf("got status: %d\n want: %d", status, http.StatusNotFound)
				}

				if body, _ := ioutil.ReadAll(response.Body); string(body) == "Test Body" {
					t.Errorf("got body: %s\n want: %s", body, "Test Body")
				}
			},
		},
	}

	for _, tt := range testSet {
		t.Run(tt.name, func(t *testing.T) {
			mockServerMux := http.NewServeMux()
			MakeHandler(mockServerMux, tt.args.address, tt.args.endpoints)
			mockServer := httptest.NewServer(mockServerMux)
			defer mockServer.Close()

			for _, request := range tt.requests {
				req, err := http.NewRequest(
					request.method,
					httpPrefix+mockServer.Listener.Addr().String()+request.endpoint,
					request.body)
				checkIfNotError(t, err)

				response, err := http.DefaultClient.Do(req)
				checkIfNotError(t, err)

				tt.want(t, response, err)
			}
		})
	}
}

func checkIfNotError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Error(err)
	}
}
