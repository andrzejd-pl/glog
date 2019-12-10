package api

import (
	"bytes"
	"github.com/andrzejd-pl/glog/repository"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllPosts(t *testing.T) {
	type args struct {
		repository repository.PostRepository
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		{
			name: "Empty test",
			args: struct{ repository repository.PostRepository }{repository: nil},
			want: func(writer http.ResponseWriter, request *http.Request) {
				http.Error(writer, ResourceNotFound, http.StatusInternalServerError)
			},
		},
	}
	for _, tt := range tests {
		buff := bytes.NewBufferString("")
		log.SetOutput(buff)
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)

			buff := httptest.NewRecorder()
			GetAllPosts(tt.args.repository).ServeHTTP(buff, request)
			got := buff.Result()

			buff = httptest.NewRecorder()
			tt.want.ServeHTTP(buff, request)
			want := buff.Result()

			if got.StatusCode != want.StatusCode {
				t.Errorf("got status code %v want %v", got.StatusCode, want.StatusCode)
			}

			gotBody, err := ioutil.ReadAll(got.Body)
			checkIfNotError(t, err)
			wantBody, err := ioutil.ReadAll(want.Body)
			checkIfNotError(t, err)

			if string(gotBody) != string(wantBody) {
				t.Errorf("got body %s want %s", gotBody, wantBody)
			}
		})
	}
}
