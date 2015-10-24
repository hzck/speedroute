// Package main creates the web server.
package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type iFileSystem interface {
	Create(name string) (iFile, error)
	IsNotExist(err error) bool
	Stat(name string) (os.FileInfo, error)
}

type iFile interface {
	Close() error
	WriteString(s string) (int, error)
}

type osFile struct {
	*os.File
}

func (ref *osFile) WriteString(s string) (int, error) {
	return ref.WriteString(s)
}

func (ref *osFile) Close() error {
	return ref.Close()
}

type osFS struct{}

func (osFS) Create(name string) (iFile, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &osFile{f}, nil
}

func (osFS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (osFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func main() {
	fs := osFS{}
	r := mux.NewRouter()
	r.HandleFunc("/create/{id:[A-Za-z0-9-_]+}", createGraphHandler(fs)).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8000", r))
}

func createGraphHandler(fs iFileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := "graphs/" + mux.Vars(r)["id"] + ".json"
		if _, err := fs.Stat(filename); fs.IsNotExist(err) {
			file, err := fs.Create(filename)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = file.WriteString("{}")
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = file.Close()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}
}
