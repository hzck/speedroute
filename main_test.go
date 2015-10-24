package main

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var fileExists = "file_exists"
var errorOnCreate = "error_on_Create"
var errorOnWriteString = "error_on_WriteString"
var errorOnClose = "error_on_Close"

func fullPath(name string) string {
	return "graphs/" + name + ".json"
}

func equals(input, fileName string) bool {
	return strings.Compare(input, fullPath(fileName)) == 0
}

type osMock struct{}

func (osMock) Create(name string) (iFile, error) {
	if equals(name, errorOnCreate) {
		return nil, errors.New("I failed to create your file, master")
	}
	return &fileMock{name}, nil
}

func (osMock) IsNotExist(err error) bool {
	return err == nil
}

func (osMock) Stat(name string) (os.FileInfo, error) {
	if equals(name, fileExists) {
		return nil, errors.New("I sense there is another file with the same name, hm")
	}
	return nil, nil
}

type fileMock struct {
	name string
}

func (f *fileMock) Close() error {
	if equals(f.name, errorOnClose) {
		return errors.New("I can't believe I'm unable to close the file")
	}
	return nil
}
func (f *fileMock) WriteString(s string) (int, error) {
	if equals(f.name, errorOnWriteString) {
		return 0, errors.New("There's something wrong when writing this text")
	}
	return 0, nil
}

// TestCreateGraph checks that createGraph creates a file on the file system if it doesn't exist.
func TestCreateGraph(t *testing.T) {
	var fileNames = []struct {
		in  string
		out int
	}{
		{"unique", http.StatusCreated},
		{fileExists, http.StatusConflict},
		{errorOnCreate, http.StatusInternalServerError},
		{errorOnWriteString, http.StatusInternalServerError},
		{errorOnClose, http.StatusInternalServerError},
	}
	fs := osMock{}
	r := mux.NewRouter()
	r.HandleFunc("/create/{id:[A-Za-z0-9-_]+}", createGraphHandler(fs)).Methods("POST")
	for _, val := range fileNames {
		req, _ := http.NewRequest("POST", "/create/"+val.in, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != val.out {
			t.Errorf("%s didn't return http code %v", val.in, val.out)
		} else {
			t.Logf("OOOOOOOK! %s did in fact return http code %v", val.in, val.out)
		}
	}
}
