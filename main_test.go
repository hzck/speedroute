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

var okString = "ok"
var fileExists = "file_exists"
var errorOnCreate = "error_on_Create"
var errorOnWriteString = "error_on_WriteString"
var errorOnPWCreate = "error_on_pw_Create"
var errorOnPWWriteString = "error_on_pw_WriteString"
var superMario64 = "Super_Mario_64"
var closeCount int
var removeCount int

func fullPathJSON(name string) string {
	return "graphs/" + name + ".json"
}

func fullPathPW(name string) string {
	return "passwords/" + name + ".txt"
}

func equalsJSON(input, fileName string) bool {
	return strings.Compare(input, fullPathJSON(fileName)) == 0
}

func equalsPW(input, fileName string) bool {
	return strings.Compare(input, fullPathPW(fileName)) == 0
}

type ioutilMock struct{}

func (ioutilMock) ReadDir(dirname string) ([]string, error) {
	return []string{superMario64}, nil
}

type osMock struct{}

func (osMock) Create(name string) (iFile, error) {
	if equalsJSON(name, errorOnCreate) || equalsPW(name, errorOnPWCreate) {
		return nil, errors.New("I failed to create your file, master")
	}
	return &fileMock{name}, nil
}

func (osMock) IsNotExist(err error) bool {
	return err == nil
}

func (osMock) Remove(name string) error {
	removeCount++
	return nil
}

func (osMock) Stat(name string) (os.FileInfo, error) {
	if equalsJSON(name, fileExists) {
		return nil, errors.New("I sense there is another file with the same name, hm")
	}
	return nil, nil
}

type fileMock struct {
	name string
}

func (f fileMock) Name() string {
	return f.name
}

func (f *fileMock) Close() error {
	closeCount++
	return nil
}

func (f *fileMock) WriteString(s string) (int, error) {
	if equalsJSON(f.name, errorOnWriteString) || equalsPW(f.name, errorOnPWWriteString) {
		return 0, errors.New("There's something wrong when writing this text")
	}
	if strings.HasSuffix(f.name, ".json") {
		if strings.Compare(s, "{}") != 0 {
			return 0, errors.New("Basic file content not specified correctly")
		}
	} else if strings.Compare(s, okString) != 0 {
		return 0, errors.New("Password file content not specified correctly")
	}
	return 0, nil
}

// TestCreateGraph checks that createGraph creates a file on the file system if it doesn't exist.
func TestCreateGraph(t *testing.T) {
	var fileNames = []struct {
		in          string
		pw          string
		out         int
		closeCount  int
		removeCount int
	}{
		{okString, okString, http.StatusCreated, 2, 0},
		{fileExists, okString, http.StatusConflict, 0, 0},
		{errorOnCreate, okString, http.StatusInternalServerError, 0, 0},
		{errorOnWriteString, okString, http.StatusInternalServerError, 1, 1},
		{errorOnPWCreate, okString, http.StatusInternalServerError, 1, 1},
		{errorOnPWWriteString, okString, http.StatusInternalServerError, 2, 2},
	}
	fs := osMock{}
	r := mux.NewRouter()
	r.HandleFunc("/create/{id:[A-Za-z0-9-_]+}/{pw:.*}", createGraphHandler(fs)).Methods("POST")
	for _, val := range fileNames {
		removeCount = 0
		req, _ := http.NewRequest("POST", "/create/"+val.in+"/"+val.pw, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != val.out || val.removeCount != removeCount {
			t.Errorf("id: %s, pw: %s - wanted http status: %v, was: %v - wanted removeCount: %v, was: %v",
				val.in, val.pw, val.out, w.Code, val.removeCount, removeCount)
		}
	}
}

// TestListGraphs lists all graphs stored on the file system.
func TestListGraphs(t *testing.T) {
	ioutil := ioutilMock{}
	r := mux.NewRouter()
	r.HandleFunc("/graphs", listGraphsHandler(ioutil)).Methods("GET")
	removeCount = 0
	req, _ := http.NewRequest("GET", "/graphs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if strings.Compare(w.Body.String(), "[\""+superMario64+"\"]") != 0 {
		t.Errorf("w.Body: %s != %s", w.Body.String(), "[\""+superMario64+"\"]")
	}
}

// TestExtractNames tests that a list of os.FileInfo becomes a list of strings without file extension.
func TestExtractNames(t *testing.T) {
	testData := []string{superMario64, "anothertest", "random"}
	fileMocks := make([]iFileInfo, len(testData))
	for i, v := range testData {
		fileMocks[i] = fileMock{v + ".json"}
	}
	for i, file := range extractNames(fileMocks) {
		if strings.Compare(testData[i], file) != 0 {
			t.Errorf("%s != %s", testData[i], file)
		}
	}
}

/*:= mux.NewRouter()
	r.HandleFunc("/graph/{id:[A-Za-z0-9-_]+}/{pw:.*}", listGraphsHandler(ioutil)).Methods("GET")
	removeCount = 0
	req, _ := http.NewRequest("GET", "/graph/" + superMario64 + "/" + , nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if strings.Compare(w.Body.String(), "[\""+superMario64+"\"]") != 0 {
		t.Errorf("w.Body: %s != %s", w.Body.String(), "[\""+superMario64+"\"]")
	}

}*/

// TestRouting tests the routing functionality in main.
func TestRouting(t *testing.T) {

}
