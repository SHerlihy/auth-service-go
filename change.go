package main

import (
	"encoding/json"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ChangeUser(w http.ResponseWriter, req *http.Request, id string, session string, changeStmt *sql.Stmt) {
        var (
            dbSession string
        )

        hasEmailErr := db.QueryRow("SELECT session FROM user WHERE id = ?", id).Scan(&dbSession)

        if dbSession == "" {
            hadErr := HandleErr(w, hasEmailErr, http.StatusBadRequest)
            if hadErr {
                return
            }
            return
        }

        if dbSession != session {
            log.Fatal("session doesn't match")
            http.Error(w, hasEmailErr.Error(), http.StatusBadRequest)
            return
        }
    
        _, err := changeStmt.Exec(id)
        hadErr := HandleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }

        retVals := AuthUser{
            dbSession,
            id,
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(retVals)
}

