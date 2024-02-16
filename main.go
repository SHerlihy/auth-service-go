package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type AuthUser struct {
	Id      string `json:"id"`
	Session string `json:"session"`
}

func main() {
	var DB *sql.DB

	//    env := flag.String("e", "dev", "environment: dev|prod")
	//    flag.Parse()
	//
	//    if *env == "dev" {
	//        env_vars.DevEnv()
	//    }
	//
	//    if *env == "prod" {
	//        env_vars.ProdEnv()
	//    }

	DevEnv()

	accessStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBADDR"),
		os.Getenv("DBPORT"),
		os.Getenv("DBDATABASE"),
	)

	fmt.Println("accessStr")
	fmt.Println(accessStr)

	var err error
	DB, err = sql.Open("mysql", accessStr)
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	routeHandler := http.NewServeMux()

	handleAccessUser := func(w http.ResponseWriter, req *http.Request) {
		AccessUser(DB, w, req)
	}
	handleChangeEmail := func(w http.ResponseWriter, req *http.Request) {
		ChangeEmail(DB, w, req)
	}
	handleChangePassword := func(w http.ResponseWriter, req *http.Request) {
		ChangePassword(DB, w, req)
	}
	handleDeleteUser := func(w http.ResponseWriter, req *http.Request) {
		DeleteUser(DB, w, req)
	}

	routeHandler.HandleFunc("POST /access", handleAccessUser)

	routeHandler.HandleFunc("PUT /email", handleChangeEmail)

	routeHandler.HandleFunc("PUT /password", handleChangePassword)

	routeHandler.HandleFunc("DELETE /delete", handleDeleteUser)

	server := &http.Server{
		Addr:    ":8080",
		Handler: routeHandler,
	}

	log.Fatal(server.ListenAndServe())
}

func IsSessionValid(dbConn *sql.DB, w http.ResponseWriter, id string, session string) error {
	var (
		DBSession string
	)

	hasEmailErr := dbConn.QueryRow("SELECT session FROM user WHERE id = ?", id).Scan(&DBSession)

	if hasEmailErr == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	if DBSession != session {
		log.Fatal("session doesn't match")
		return errors.New("sessions do not match")
	}

	return nil
}

func HandleErr(w http.ResponseWriter, err error, status int) bool {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), status)
		return true
	}

	return false
}
