// Package main creates the web server.
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	a "github.com/hzck/speedroute/algorithm"
	p "github.com/hzck/speedroute/parser"
)

var mutex = &sync.Mutex{}

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
	// Creating logfile
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			panic(closeErr)
		}
	}()
	log.SetOutput(f)
	// Reading environment variables
	err = godotenv.Load("config.env")
	if err != nil {
		log.Println("Error loading config.env file")
		panic(err)
	}
	// Connecting to the DB
	dbConnString := "postgresql://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_URL") + ":" +
		os.Getenv("DB_PORT") + "/speedroute?pool_max_conns=10"
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		log.Println("Unable to connect to database: ", err)
		panic(err)
	}
	defer dbpool.Close()
	log.Println("### Server started ###")
	fs := osFS{}
	io := ioutilFS{}
	r := mux.NewRouter()
	r.HandleFunc("/create/{id:[A-Za-z0-9-_]+}/{pw:.*}", createGraphHandler(fs)).Methods("POST")
	r.HandleFunc("/graphs", listGraphsHandler(io)).Methods("GET")
	r.HandleFunc("/graph/{id:[A-Za-z0-9-_]+}", getGraphHandler(fs)).Methods("GET")
	r.HandleFunc("/graph/{id:[A-Za-z0-9-_]+}/{pw:.*}", saveGraphHandler(fs)).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Println(http.ListenAndServe(":8001", r))
}

func createFile(fs iFileSystem, filename, fileContent string) (httpStatus int, err error) {
	file, err := fs.Create(filename)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			log.Println(closeErr)
			httpStatus = http.StatusInternalServerError
			err = closeErr
			return
		}
	}()
	_, err = file.WriteString(fileContent)
	if err != nil {
		log.Println(err)
		removeErr := fs.Remove(filename)
		if removeErr != nil {
			log.Println(removeErr)
		}
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
					removeErr := fs.Remove(filename)
					if removeErr != nil {
						log.Println(removeErr)
					}
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
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
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
		err = ioutil.WriteFile(filename, newGraph, 0600)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		graph, err := p.CreateGraphFromFile(filename)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Potential bottleneck draining memory, making sure only one graph is routed at any moment.
		// Should add timing in log if it takes > 1s? 10s?
		mutex.Lock()
		start := time.Now()
		result := a.Route(graph)
		log.Printf("%s taking %s\n", filename, time.Since(start))
		mutex.Unlock()
		js, err := p.CreateJSONFromRoutedPath(result)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(js)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func listGraphsHandler(io iIOUtil) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		graphs, err := io.ReadDir("graphs")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(graphs)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
