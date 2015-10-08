//package main creates the web server
package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create/{id:[A-Za-z0-9-_]+}", createGraph).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8000", r))
}

func createGraph(w http.ResponseWriter, r *http.Request) {
	filename := "graphs/" + mux.Vars(r)["id"] + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
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
