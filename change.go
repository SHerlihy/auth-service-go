package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type RequestChangeEmail struct {
	RequestUserAuth
	Email string `json:"email"`
}

type RequestChangePassword struct {
	RequestUserAuth
	Password string `json:"password"`
}

func ChangeEmail(dbConn *sql.DB, w http.ResponseWriter, req *http.Request) {
	fmt.Println("/change email")

	if req.URL.Path != "/email" {
		http.NotFound(w, req)
		return
	}

	var user RequestChangeEmail

	err := json.NewDecoder(req.Body).Decode(&user)
	hadErr := HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	err = IsSessionValid(dbConn, w, user.Id, user.CurrentSession)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	assignEmailStmt, err := dbConn.Prepare(fmt.Sprintf("UPDATE user SET email = '%s' WHERE id = ?", user.Email))
	defer assignEmailStmt.Close()
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	_, err = assignEmailStmt.Exec(user.Id)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func ChangePassword(dbConn *sql.DB, w http.ResponseWriter, req *http.Request) {
	fmt.Println("/change pass")

	if req.URL.Path != "/password" {
		http.NotFound(w, req)
		return
	}

	var user RequestChangePassword

	err := json.NewDecoder(req.Body).Decode(&user)
	hadErr := HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	err = IsSessionValid(dbConn, w, user.Id, user.CurrentSession)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	assignPassStmt, err := dbConn.Prepare(fmt.Sprintf("UPDATE user SET password = '%s' WHERE id = ?", user.Password))
	defer assignPassStmt.Close()
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	_, err = assignPassStmt.Exec(user.Id)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
