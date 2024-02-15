package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func getUser(w http.ResponseWriter, req *http.Request, user User) {
    var (
        addId int
        addSession string
    )

    addedRow := db.QueryRow("SELECT id, session FROM user WHERE email = ?", user.Email)
    
    err := addedRow.Scan(&addId, &addSession)
    hadErr := HandleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    retVals := AuthUser{
        string(addId),
        addSession,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(retVals)
}

func Login(w http.ResponseWriter, req *http.Request, user User) {
    var (
        dbPass string
    )

    userRow := db.QueryRow("SELECT password FROM user WHERE email = ?", user.Email)

    err := userRow.Scan(&dbPass)
    hadErr := HandleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    if dbPass != user.Password {
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

    assignSessionStmt, err := db.Prepare("UPDATE user SET session = ? WHERE email = ?")
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

    getUser(w, req, user)
}

func Register(w http.ResponseWriter, req *http.Request, user User) {
    sessionBytes, err := uuid.NewUUID()
    hadErr := HandleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    sessionId := sessionBytes.String() 

    // adding session here in case of error so don't have to reverse change
    createUserStmt, err := db.Prepare("INSERT INTO user (email, password, session) VALUES (?, ?, ?)")
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

    getUser(w, req, user)
}

