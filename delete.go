package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func DeleteUser(dbConn *sql.DB, w http.ResponseWriter, req *http.Request) {
	fmt.Println("/delte user")

	if req.URL.Path != "/delete" {
		http.NotFound(w, req)
		return
	}

	var user AuthUser

	err := json.NewDecoder(req.Body).Decode(&user)
	hadErr := HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	err = IsSessionValid(dbConn, w, user.Id, user.Session)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	deleteUserStmt, err := dbConn.Prepare("DELETE FROM user WHERE id = ?")
	defer deleteUserStmt.Close()
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	_, err = deleteUserStmt.Exec(user.Id)
	hadErr = HandleErr(w, err, http.StatusBadRequest)
	if hadErr {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
