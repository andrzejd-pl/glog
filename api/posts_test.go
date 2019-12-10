package api

import (
	"bytes"
	"github.com/andrzejd-pl/glog/repository"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
			name: "Nil repository test",
			args: struct{ repository repository.PostRepository }{repository: nil},
			want: func(writer http.ResponseWriter, request *http.Request) {
				http.Error(writer, ResourceNotFound, http.StatusInternalServerError)
			},
		},
		{
			name: "Empty resource test",
			args: struct{ repository repository.PostRepository }{repository: mockPostRepository(repository.Posts{}, nil)},
			want: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte("[]"))
			},
		},
		{
			name: "Example data test",
			args: struct{ repository repository.PostRepository }{repository: mockPostRepository(repository.Posts{
				stubPost("123456", 10, "Test", "Test Content", "2019-12-31 00:00:00", "2006-01-02 15:04:05"),
			}, nil)},
			want: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte("[{\"Uuid\":123456,\"Title\":\"Test\",\"Content\":\"Test Content\",\"Timestamp\":\"2019-12-31T00:00:00Z\"}]"))
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

type PostRepositoryMock struct {
	posts repository.Posts
	err   error
}

func (repo PostRepositoryMock) GetAllPosts() (repository.Posts, error) {
	return repo.posts, repo.err
}

func mockPostRepository(posts repository.Posts, err error) repository.PostRepository {
	return PostRepositoryMock{
		posts: posts,
		err:   err,
	}
}

func stubPost(uuidValue string, uuidBase int, title, content, timestampValue, timestampFormat string) repository.Post {
	uuid := big.Int{}
	uuid.SetString(uuidValue, uuidBase)
	timestamp, _ := time.Parse(timestampFormat, timestampValue)
	return repository.Post{
		Uuid:      uuid,
		Title:     title,
		Content:   content,
		Timestamp: timestamp,
	}
}
