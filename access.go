package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type RequestAccessUser struct {
	Email    string
	Password string
}

func AccessUser(dbConn *sql.DB, w http.ResponseWriter, req *http.Request) {
	fmt.Println("/access request")

	if req.URL.Path != "/access" {
		http.NotFound(w, req)
		return
	}

	var user RequestAccessUser

	err := json.NewDecoder(req.Body).Decode(&user)
	hadErr := HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	var (
		dbConnId int
	)
	hasEmailErr := dbConn.QueryRow("SELECT id FROM user WHERE email = ?", user.Email).Scan(&dbConnId)

	if hasEmailErr == sql.ErrNoRows {
		register(dbConn, w, req, user)
		return
	}

	if dbConnId > 0 {
		login(dbConn, w, req, user)
		return
	}

	//Scan always has err != nil
	HandleErr(w, hasEmailErr, http.StatusBadRequest)
}

func login(dbConn *sql.DB, w http.ResponseWriter, req *http.Request, user RequestAccessUser) {
	var (
		dbConnPass string
	)

	userRow := dbConn.QueryRow("SELECT password FROM user WHERE email = ?", user.Email)

	err := userRow.Scan(&dbConnPass)
	hadErr := HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	if dbConnPass != user.Password {
		log.Fatal("password doesn't match")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionBytes, err := uuid.NewUUID()
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	sessionId := sessionBytes.String()

    prevAuth, err := getCurrentAuth(dbConn, w, req, user)

	assignSessionStmt, err := dbConn.Prepare("UPDATE user SET session = ? WHERE email = ?")
	defer assignSessionStmt.Close()
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	_, err = assignSessionStmt.Exec(sessionId, user.Email)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

    retVals := ResponseUserAuth{
        RequestUserAuth: RequestUserAuth{
            Id: prevAuth.Id,
            CurrentSession: sessionId,
        },
        PreviousSession: prevAuth.CurrentSession,
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(retVals)
}

func register(dbConn *sql.DB, w http.ResponseWriter, req *http.Request, user RequestAccessUser) {
	sessionBytes, err := uuid.NewUUID()
	hadErr := HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	sessionId := sessionBytes.String()

	// adding session here in case of error so don't have to reverse change
	createUserStmt, err := dbConn.Prepare("INSERT INTO user (email, password, session) VALUES (?, ?, ?)")
	defer createUserStmt.Close()
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	_, err = createUserStmt.Exec(user.Email, user.Password, sessionId)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

    currentAuth, err := getCurrentAuth(dbConn, w, req, user)

    retVals := ResponseUserAuth{
        RequestUserAuth: RequestUserAuth{
            Id: currentAuth.Id,
            CurrentSession: currentAuth.CurrentSession,
        },
        PreviousSession: "",
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(retVals)
}

func getCurrentAuth(dbConn *sql.DB, w http.ResponseWriter, req *http.Request, user RequestAccessUser) (RequestUserAuth, error) {
	var (
        retVals RequestUserAuth
		addId      int
		addSession string
	)

	addedRow := dbConn.QueryRow("SELECT id, session FROM user WHERE email = ?", user.Email)

	err := addedRow.Scan(&addId, &addSession)
	if err != nil {
		return RequestUserAuth{}, err
	}

    retVals = RequestUserAuth{
		fmt.Sprint(addId),
		addSession,
	}

    return retVals, nil
}
