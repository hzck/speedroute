// Package main creates the web server.
package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	a "github.com/hzck/speedroute/algorithm"
	p "github.com/hzck/speedroute/parser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type iIOUtil interface {
	ReadDir(dirname string) ([]string, error)
}

type iFileSystem interface {
	Create(name string) (iFile, error)
	IsNotExist(err error) bool
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
}

type iFile interface {
	Close() error
	WriteString(s string) (int, error)
}

type iFileInfo interface {
	Name() string
}

type ioutilFS struct{}

func (ioutilFS) ReadDir(dirname string) ([]string, error) {
	readDir, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]iFileInfo, len(readDir))
	for i, v := range readDir {
		fileInfos[i] = v
	}
	return extractNames(fileInfos), nil
}

func extractNames(readDir []iFileInfo) []string {
	list := make([]string, len(readDir))
	for i := range list {
		file := readDir[i].Name()
		list[i] = file[0 : len(file)-5]
	}
	return list
}

type osFile struct {
	file *os.File
}

func (ref *osFile) WriteString(s string) (int, error) {
	return ref.file.WriteString(s)
}

func (ref *osFile) Close() error {
	return ref.file.Close()
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

func (osFS) Remove(name string) error {
	return os.Remove(name)
}

func (osFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func main() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("### Server started ###")
	fs := osFS{}
	ioutil := ioutilFS{}
	r := mux.NewRouter()
	r.HandleFunc("/create/{id:[A-Za-z0-9-_]+}/{pw:.*}", createGraphHandler(fs)).Methods("POST")
	r.HandleFunc("/graphs", listGraphsHandler(ioutil)).Methods("GET")
	r.HandleFunc("/graph/{id:[A-Za-z0-9-_]+}", getGraphHandler(fs)).Methods("GET")
	r.HandleFunc("/graph/{id:[A-Za-z0-9-_]+}/{pw:.*}", saveGraphHandler(fs)).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8001", r))
}

func createFile(fs iFileSystem, filename, fileContent string) (int, error) {
	file, err := fs.Create(filename)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	defer file.Close()
	_, err = file.WriteString(fileContent)
	if err != nil {
		log.Println(err)
		fs.Remove(filename)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func createGraphHandler(fs iFileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := "graphs/" + mux.Vars(r)["id"] + ".json"
		if _, err := fs.Stat(filename); fs.IsNotExist(err) {
			fileContent := "{}"
			newGraph, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if len(newGraph) > 0 {
				fileContent, err = p.LivesplitXMLtoJSON(string(newGraph))
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			httpStatus, err := createFile(fs, filename, fileContent)
			if err == nil {
				pwFilename := "passwords/" + mux.Vars(r)["id"] + ".txt"
				httpStatus, err = createFile(fs, pwFilename, mux.Vars(r)["pw"])
				if err != nil {
					fs.Remove(filename)
				}
			}
			w.WriteHeader(httpStatus)
			return
		}
		w.WriteHeader(http.StatusConflict)
	}
}

func getGraphHandler(fs iFileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := "graphs/" + mux.Vars(r)["id"] + ".json"
		if _, err := fs.Stat(filename); fs.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
		return
	}
}

func saveGraphHandler(fs iFileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := "graphs/" + mux.Vars(r)["id"] + ".json"
		if _, err := fs.Stat(filename); fs.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		pwFilename := "passwords/" + mux.Vars(r)["id"] + ".txt"
		if _, err := fs.Stat(pwFilename); fs.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		pw, _ := ioutil.ReadFile(pwFilename)
		if strings.Compare(string(pw), mux.Vars(r)["pw"]) != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newGraph, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ioutil.WriteFile(filename, newGraph, 0644)
		result := a.Route(p.CreateGraphFromFile(filename))
		js, err := p.CreateJSONFromRoutedPath(result)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(js)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func listGraphsHandler(io iIOUtil) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		graphs, err := io.ReadDir("graphs")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(graphs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
