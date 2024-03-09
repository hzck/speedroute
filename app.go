package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	model "github.com/hzck/speedroute/model"
)

// App holds information about the running application.
type App struct {
	Router     *mux.Router
	Dbpool     *pgxpool.Pool
	HashParams *argon2id.Params
}

// InitLogFile initializes app output to "logfile"
func (a *App) InitLogFile() func() {
	// Creating logfile
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	return func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			panic(closeErr)
		}
	}
}

// InitConfigFile reads application configuration from "config.env"
func (a *App) InitConfigFile() {
	// Reading environment variables
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println("Error loading config.env file")
		panic(err)
	}
	a.HashParams = &argon2id.Params{
		Memory:      parseUInt32FromOsEnv("ARGON2ID_MEMORY"),
		Iterations:  parseUInt32FromOsEnv("ARGON2ID_ITERATIONS"),
		Parallelism: parseUInt8FromOsEnv("ARGON2ID_PARALLELISM"),
		SaltLength:  parseUInt32FromOsEnv("ARGON2ID_SALT_LENGTH"),
		KeyLength:   parseUInt32FromOsEnv("ARGON2ID_KEY_LENGTH"),
	}
}

func parseUInt8FromOsEnv(key string) uint8 {
	val, err := strconv.ParseUint(os.Getenv(key), 10, 8)
	if err != nil {
		log.Printf("Error parsing uint8 from config.env file with key %s\n", key)
		panic(err)
	}
	return uint8(val)
}

func parseUInt32FromOsEnv(key string) uint32 {
	val, err := strconv.ParseUint(os.Getenv(key), 10, 32)
	if err != nil {
		log.Printf("Error parsing uint32 from config.env file with key %s\n", key)
		panic(err)
	}
	return uint32(val)
}

// InitDB connects to the database.
func (a *App) InitDB() func() {
	dbConnString := "postgresql://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_URL") + ":" +
		os.Getenv("DB_PORT") + "/speedroute?pool_max_conns=10"
	var err error
	a.Dbpool, err = pgxpool.New(context.Background(), dbConnString)
	if err != nil {
		log.Println("Unable to connect to database: ", err)
		panic(err)
	}
	log.Println("Successfully connected to the DB")
	return func() {
		a.Dbpool.Close()
	}
}

// InitRoutes initializes the routes used by the application
func (a *App) InitRoutes() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/signup", a.createAccount(a.Dbpool)).Methods("POST")
	fs := http.FileServer(http.Dir("frontend/build"))
	a.Router.PathPrefix("/").Handler(http.StripPrefix("/", fs))
}

// Run starts the application.
func (a *App) Run() {
	log.Println("### Server started ###")
	log.Println(http.ListenAndServe(":8001", a.Router))
}

func (a *App) createAccount(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		newAccount := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := decoder.Decode(&newAccount)
		match, _ := regexp.MatchString("^[\\w]{2,30}$", newAccount.Username)
		if err != nil || !match || len(newAccount.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var account model.Account
		account.Username = strings.ToLower(newAccount.Username)
		account.DisplayName = newAccount.Username

		//TODO: If username already taken, no need to create hash
		hash, err := argon2id.CreateHash(newAccount.Password, a.HashParams)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		account.Password = hash
		err = account.CreateAccount(dbpool)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		ws, err := model.CreateWebsession(dbpool, account.ID, parseUInt32FromOsEnv("WEBSESSION_HOURS_LOGGED_IN"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "speedroute_websession", Expires: ws.ExpireAt, Value: ws.Token.String()})
		w.WriteHeader(http.StatusCreated)
	}
}
