package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	model "github.com/hzck/speedroute/model"
)

var app App

func TestMain(m *testing.M) {
	app.InitConfigFile()
	app.InitDB()
	app.InitRoutes()
	code := m.Run()
	os.Exit(code)
}

func TestCreateAccount(t *testing.T) {
	displayname := "Valid_Username_7"
	password := "val!dP@s5word"

	req := createPostRequestForCreateAccount(displayname, password)

	response := executeRequest(req)
	defer clearAccountInDB(displayname)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var a model.Account
	query := "SELECT * FROM account WHERE displayname=$1"
	err := app.Dbpool.QueryRow(context.Background(), query, displayname).Scan(&a.ID, &a.Username, &a.DisplayName, &a.Password, &a.Created, &a.LastUpdated)
	if err != nil {
		panic(err)
	}
	if a.ID <= 0 {
		t.Errorf("ID field is not valid")
	}
	username := strings.ToLower(displayname)
	if a.Username != username {
		t.Errorf("Username '%s' is not the expected '%s'", a.Username, username)
	}
	if a.DisplayName != displayname {
		t.Errorf("DisplayName '%s' is not the expected '%s'", a.DisplayName, displayname)
	}
	match, _ := regexp.MatchString("^\\$argon2id\\$v=19\\$m=65536,t=8,p=1\\$.{22}\\$.{43}$", a.Password)
	if !match {
		t.Errorf("Password field is not set correctly")
	}
	if a.Created.IsZero() {
		t.Errorf("Created field is not set")
	}
	if a.LastUpdated.IsZero() {
		t.Errorf("LastUpdated field is not set")
	}
}

func TestCreateAccountInvalidJSON(t *testing.T) {
	req := createPostRequestForCreateAccount("\",{invalidjson", "val!dP@s5word")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateAccountUsernameNotPopulated(t *testing.T) {
	req := createPostRequestForCreateAccount("", "val!dP@s5word")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateAccountPasswordNotPopulated(t *testing.T) {
	req := createPostRequestForCreateAccount("passwordnotpopulated", "")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateAccountDuplicateUsername(t *testing.T) {
	username := "duplicate"
	createAccountInDB(username)
	defer clearAccountInDB(username)
	req := createPostRequestForCreateAccount(username, "val!dP@s5word")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusConflict, response.Code)
}

func TestCreateAccountUsernameTooShort(t *testing.T) {
	req := createPostRequestForCreateAccount("1", "val!dP@s5word")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateAccountUsernameTooLong(t *testing.T) {
	req := createPostRequestForCreateAccount("this_user_is_31_characters_long", "val!dP@s5word")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateAccountPasswordTooShort(t *testing.T) {
	req := createPostRequestForCreateAccount("shortpassword", "tooshrt")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateAccountInvalidCharacters(t *testing.T) {
	// Note: Not all invalid characters are tested for obvious reasons
	invalidChars := `!"#¤%&/()=?<>|[]{}+-.,;:^*`
	for _, ch := range invalidChars {
		req := createPostRequestForCreateAccount(fmt.Sprintf("username_invalid%c", ch), "val!dP@s5word")
		response := executeRequest(req)
		checkResponseCode(t, http.StatusBadRequest, response.Code)
	}
}

func createPostRequestForCreateAccount(username, password string) *http.Request {
	usernameJSON := ""
	if username != "" {
		usernameJSON = `"username":"` + username + `"`
	}
	passwordJSON := ""
	if password != "" {
		passwordJSON = `"password":"` + password + `"`
	}
	commaJSON := ""
	if usernameJSON != "" && passwordJSON != "" {
		commaJSON = ","
	}
	var jsonStr = []byte("{" + usernameJSON + commaJSON + passwordJSON + "}")
	req, err := http.NewRequest("POST", "/account", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func createAccountInDB(displayname string) {
	query := "INSERT INTO account (username, displayname, password, created, last_updated) VALUES ($1, $2, $3, $4, $4)"
	_, err := app.Dbpool.Exec(context.Background(), query, strings.ToLower(displayname), displayname, "password", time.Now())
	if err != nil {
		panic(err)
	}
}

func clearAccountInDB(displayname string) {
	query := "DELETE FROM account WHERE displayname=$1"
	_, err := app.Dbpool.Exec(context.Background(), query, displayname)
	if err != nil {
		panic(err)
	}
}
